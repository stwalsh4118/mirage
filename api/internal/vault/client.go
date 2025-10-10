package vault

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// Client is an HTTP client wrapper for interacting with Vault's HTTP API
type Client struct {
	address        string
	token          string
	tokenMu        sync.RWMutex
	httpClient     *http.Client
	namespace      string
	mountPath      string
	renewalStop    chan struct{}
	stopOnce       sync.Once
	circuitBreaker *CircuitBreaker
	healthStop     chan struct{}
	healthStopOnce sync.Once
	shutdownOnce   sync.Once
}

// NewClient creates a new Vault HTTP client with the provided configuration
func NewClient(cfg Config) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid vault config: %w", err)
	}

	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: DefaultTimeout,
	}

	// Configure TLS if needed (dev only)
	if cfg.SkipVerify {
		log.Warn().Msg("Vault TLS verification disabled - only use in development!")
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				MinVersion:         tls.VersionTLS12,
			},
		}
	}

	client := &Client{
		address:        cfg.Address,
		token:          cfg.Token,
		httpClient:     httpClient,
		namespace:      cfg.Namespace,
		mountPath:      cfg.MountPath,
		renewalStop:    make(chan struct{}),
		circuitBreaker: NewCircuitBreaker(),
		healthStop:     make(chan struct{}),
	}

	// Authenticate with Vault
	if err := client.authenticate(cfg); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Ensure KV v2 secrets engine is mounted at the configured path
	if err := client.ensureMountExists(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ensure mount exists: %w", err)
	}

	// Start periodic health checking
	go client.startHealthChecker(context.Background(), DefaultHealthCheckInterval)

	log.Info().
		Str("address", cfg.Address).
		Str("mount_path", cfg.MountPath).
		Bool("skip_verify", cfg.SkipVerify).
		Msg("Vault HTTP client initialized")

	return client, nil
}

// makeRequest is a helper method for making HTTP requests to Vault with circuit breaker protection
func (c *Client) makeRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	// Wrap with circuit breaker
	return c.circuitBreaker.Call(func() error {
		return c.doRequest(ctx, method, path, body, result)
	})
}

// doRequest performs the actual HTTP request to Vault
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	// Create request body if needed
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, method, c.address+path, reqBody)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	// Set headers
	req.Header.Set("X-Vault-Token", c.getToken())
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.namespace != "" {
		req.Header.Set("X-Vault-Namespace", c.namespace)
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	// Handle error responses
	if resp.StatusCode >= 400 {
		var errResp struct {
			Errors []string `json:"errors"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			// Failed to decode JSON error response, read raw body as fallback
			rawBody, readErr := io.ReadAll(resp.Body)
			if readErr != nil || len(rawBody) == 0 {
				return fmt.Errorf("vault error (%d): failed to decode error response: %w", resp.StatusCode, err)
			}
			return fmt.Errorf("vault error (%d): failed to decode error response: %w (raw body: %s)", resp.StatusCode, err, string(rawBody))
		}
		return fmt.Errorf("vault error (%d): %v", resp.StatusCode, errResp.Errors)
	}

	// Decode response if needed
	if result != nil && resp.StatusCode != 204 {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}

// HealthCheck performs a health check against the Vault server
func (c *Client) HealthCheck(ctx context.Context) error {
	status, err := c.GetStatus(ctx)
	if err != nil {
		return err
	}

	if status.Sealed {
		return fmt.Errorf("vault is sealed")
	}

	if !status.Initialized {
		return fmt.Errorf("vault is not initialized")
	}

	if !status.Available {
		return fmt.Errorf("vault is not available")
	}

	return nil
}

// GetStatus retrieves detailed status information from Vault
func (c *Client) GetStatus(ctx context.Context) (*VaultStatus, error) {
	// Health check doesn't require authentication, so use direct HTTP
	req, err := http.NewRequestWithContext(ctx, "GET", c.address+"/v1/sys/health", nil)
	if err != nil {
		return &VaultStatus{Available: false}, fmt.Errorf("create health check request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &VaultStatus{Available: false}, fmt.Errorf("execute health check: %w", err)
	}
	defer resp.Body.Close()

	// Parse health response
	var healthResp struct {
		Initialized bool   `json:"initialized"`
		Sealed      bool   `json:"sealed"`
		Standby     bool   `json:"standby"`
		Version     string `json:"version"`
		ClusterName string `json:"cluster_name"`
		ClusterID   string `json:"cluster_id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&healthResp); err != nil {
		return &VaultStatus{Available: false}, fmt.Errorf("decode health response: %w", err)
	}

	// Vault health endpoint returns different status codes based on state
	// 200 = healthy, 503 = sealed, 501 = not initialized
	available := resp.StatusCode == 200

	return &VaultStatus{
		Available:   available,
		Initialized: healthResp.Initialized,
		Sealed:      healthResp.Sealed,
		Version:     healthResp.Version,
		ClusterID:   healthResp.ClusterID,
		ClusterName: healthResp.ClusterName,
	}, nil
}

// ensureMountExists checks if the KV v2 secrets engine is mounted at the configured path
// and creates it if it doesn't exist. This ensures the mount is ready for use.
func (c *Client) ensureMountExists(ctx context.Context) error {
	// Check if mount already exists
	mountPath := fmt.Sprintf("/v1/sys/mounts/%s", c.mountPath)

	var mountInfo map[string]interface{}
	err := c.doRequest(ctx, "GET", mountPath, nil, &mountInfo)

	if err == nil {
		// Mount exists, verify it's KV v2
		if mountType, ok := mountInfo["type"].(string); ok && mountType == "kv" {
			if options, ok := mountInfo["options"].(map[string]interface{}); ok {
				if version, ok := options["version"].(string); ok && version == "2" {
					log.Debug().
						Str("mount_path", c.mountPath).
						Msg("KV v2 secrets engine already mounted")
					return nil
				}
			}
		}
		// Mount exists but wrong type/version
		log.Warn().
			Str("mount_path", c.mountPath).
			Msg("Mount exists but is not KV v2 - will not modify existing mount")
		return fmt.Errorf("mount %s exists but is not configured as KV v2", c.mountPath)
	}

	// Check if error is 404 (mount doesn't exist) or other error
	if !isNotFoundError(err) {
		return fmt.Errorf("failed to check mount status: %w", err)
	}

	// Mount doesn't exist, create it
	log.Info().
		Str("mount_path", c.mountPath).
		Msg("KV v2 secrets engine not found, creating mount")

	mountConfig := map[string]interface{}{
		"type": "kv",
		"options": map[string]interface{}{
			"version": "2",
		},
		"description": "Mirage secrets storage (KV v2)",
	}

	err = c.doRequest(ctx, "POST", mountPath, mountConfig, nil)
	if err != nil {
		return fmt.Errorf("failed to create KV v2 mount at %s: %w", c.mountPath, err)
	}

	log.Info().
		Str("mount_path", c.mountPath).
		Msg("Successfully created KV v2 secrets engine mount")

	return nil
}

// authenticate selects and executes the appropriate authentication method
func (c *Client) authenticate(cfg Config) error {
	if cfg.Token != "" {
		// Use token authentication (development)
		if err := c.authenticateWithToken(context.Background()); err != nil {
			return fmt.Errorf("token authentication failed: %w", err)
		}
		log.Info().Msg("authenticated with Vault using token")
	} else if cfg.RoleID != "" && cfg.SecretID != "" {
		// Use AppRole authentication (production)
		if err := c.authenticateWithAppRole(context.Background(), cfg.RoleID, cfg.SecretID); err != nil {
			return fmt.Errorf("AppRole authentication failed: %w", err)
		}
		log.Info().Msg("authenticated with Vault using AppRole")
	} else {
		return fmt.Errorf("no valid authentication method configured (need token or role_id+secret_id)")
	}

	// Start automatic token renewal in background
	go c.startTokenRenewal(context.Background())

	return nil
}

// authenticateWithToken validates a token by looking it up
func (c *Client) authenticateWithToken(ctx context.Context) error {
	// Verify token by looking up self
	_, err := c.lookupToken(ctx)
	if err != nil {
		return fmt.Errorf("token lookup failed: %w", err)
	}
	return nil
}

// authenticateWithAppRole performs AppRole authentication and sets the client token
func (c *Client) authenticateWithAppRole(ctx context.Context, roleID, secretID string) error {
	reqBody := AppRoleLoginRequest{
		RoleID:   roleID,
		SecretID: secretID,
	}

	var resp AppRoleLoginResponse
	if err := c.makeRequest(ctx, "POST", "/v1/auth/approle/login", reqBody, &resp); err != nil {
		return fmt.Errorf("AppRole login request failed: %w", err)
	}

	// Update client token with the one received from AppRole login
	c.setToken(resp.Auth.ClientToken)

	log.Info().
		Int("lease_duration", resp.Auth.LeaseDuration).
		Bool("renewable", resp.Auth.Renewable).
		Strs("policies", resp.Auth.Policies).
		Msg("AppRole authentication successful")

	return nil
}

// lookupToken retrieves information about the current token
func (c *Client) lookupToken(ctx context.Context) (*TokenInfo, error) {
	var resp TokenLookupResponse
	if err := c.makeRequest(ctx, "GET", "/v1/auth/token/lookup-self", nil, &resp); err != nil {
		return nil, fmt.Errorf("token lookup failed: %w", err)
	}
	return &resp.Data, nil
}

// renewToken renews the current token
func (c *Client) renewToken(ctx context.Context) (*AuthInfo, error) {
	var resp TokenRenewalResponse
	if err := c.makeRequest(ctx, "POST", "/v1/auth/token/renew-self", nil, &resp); err != nil {
		return nil, fmt.Errorf("token renewal failed: %w", err)
	}
	return &resp.Auth, nil
}

// startTokenRenewal starts a goroutine that automatically renews the token before it expires
func (c *Client) startTokenRenewal(ctx context.Context) {
	// Get token information
	tokenInfo, err := c.lookupToken(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to lookup token for renewal setup")
		return
	}

	if !tokenInfo.Renewable {
		log.Warn().Msg("Vault token is not renewable - automatic renewal disabled")
		return
	}

	if tokenInfo.TTL <= 0 {
		log.Warn().Msg("Vault token has no TTL - automatic renewal disabled")
		return
	}

	log.Info().
		Int("ttl", tokenInfo.TTL).
		Int("creation_ttl", tokenInfo.CreationTTL).
		Msg("starting automatic token renewal")

	// Calculate when to renew (at 50% of TTL)
	renewalInterval := time.Duration(float64(tokenInfo.TTL) * DefaultRenewalBuffer * float64(time.Second))
	ticker := time.NewTicker(renewalInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.renewalStop:
			log.Info().Msg("stopping token renewal")
			return
		case <-ctx.Done():
			log.Info().Msg("context cancelled, stopping token renewal")
			return
		case <-ticker.C:
			authInfo, err := c.renewToken(ctx)
			if err != nil {
				log.Error().Err(err).Msg("token renewal failed")
				// Continue trying - don't exit the renewal loop
				continue
			}

			log.Info().
				Int("new_lease_duration", authInfo.LeaseDuration).
				Msg("Vault token renewed successfully")

			// Update ticker interval based on new lease duration
			newInterval := time.Duration(float64(authInfo.LeaseDuration) * DefaultRenewalBuffer * float64(time.Second))
			ticker.Reset(newInterval)
		}
	}
}

// StopRenewal stops the automatic token renewal goroutine
func (c *Client) StopRenewal() {
	c.stopOnce.Do(func() {
		close(c.renewalStop)
	})
}

// getToken retrieves the client token with synchronization
func (c *Client) getToken() string {
	c.tokenMu.RLock()
	defer c.tokenMu.RUnlock()
	return c.token
}

// setToken updates the client token with synchronization
func (c *Client) setToken(token string) {
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()
	c.token = token
}

// GetCircuitState returns the current state of the circuit breaker
func (c *Client) GetCircuitState() CircuitState {
	return c.circuitBreaker.GetState()
}

// GetCircuitMetrics returns metrics for the circuit breaker
func (c *Client) GetCircuitMetrics() map[string]interface{} {
	return c.circuitBreaker.GetMetrics()
}

// ResetCircuit manually resets the circuit breaker to closed state
func (c *Client) ResetCircuit() {
	c.circuitBreaker.Reset()
}

// startHealthChecker starts a background goroutine that periodically checks Vault health
func (c *Client) startHealthChecker(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Info().
		Dur("interval", interval).
		Msg("starting Vault health checker")

	for {
		select {
		case <-c.healthStop:
			log.Info().Msg("stopping Vault health checker")
			return
		case <-ctx.Done():
			log.Info().Msg("context cancelled, stopping Vault health checker")
			return
		case <-ticker.C:
			if err := c.HealthCheck(ctx); err != nil {
				log.Warn().
					Err(err).
					Str("circuit_state", c.GetCircuitState().String()).
					Msg("Vault health check failed")
			} else {
				log.Debug().
					Str("circuit_state", c.GetCircuitState().String()).
					Msg("Vault health check passed")
			}
		}
	}
}

// StopHealth stops the periodic health checker
func (c *Client) StopHealth() {
	c.healthStopOnce.Do(func() {
		close(c.healthStop)
	})
}

// Shutdown gracefully stops all background goroutines
func (c *Client) Shutdown() {
	c.shutdownOnce.Do(func() {
		c.StopRenewal()
		c.StopHealth()
		log.Info().Msg("Vault client shutdown complete")
	})
}

// GetMountPath returns the KV v2 secrets engine mount path
func (c *Client) GetMountPath() string {
	return c.mountPath
}

// DeletePathRecursive deletes all secrets under a given path recursively
// This is useful for cleaning up all secrets for a user
func (c *Client) DeletePathRecursive(ctx context.Context, path string) error {
	// For KV v2, we need to delete the metadata to permanently remove the secret
	// Deleting metadata also deletes all versions
	err := c.doRequest(ctx, "DELETE", path, nil, nil)
	if err != nil {
		// If the path doesn't exist, that's not an error
		if isNotFoundError(err) {
			log.Debug().Str("path", path).Msg("path not found, nothing to delete")
			return nil
		}
		return fmt.Errorf("failed to delete path %s: %w", path, err)
	}

	log.Debug().Str("path", path).Msg("deleted path from vault")
	return nil
}
