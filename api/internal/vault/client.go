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
	address     string
	token       string
	tokenMu     sync.RWMutex
	httpClient  *http.Client
	namespace   string
	mountPath   string
	renewalStop chan struct{}
	stopOnce    sync.Once
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
		address:     cfg.Address,
		token:       cfg.Token,
		httpClient:  httpClient,
		namespace:   cfg.Namespace,
		mountPath:   cfg.MountPath,
		renewalStop: make(chan struct{}),
	}

	// Authenticate with Vault
	if err := client.authenticate(cfg); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	log.Info().
		Str("address", cfg.Address).
		Str("mount_path", cfg.MountPath).
		Bool("skip_verify", cfg.SkipVerify).
		Msg("Vault HTTP client initialized")

	return client, nil
}

// makeRequest is a helper method for making HTTP requests to Vault
func (c *Client) makeRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
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
	req, err := http.NewRequestWithContext(ctx, "GET", c.address+"/v1/sys/health", nil)
	if err != nil {
		return fmt.Errorf("create health check request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute health check: %w", err)
	}
	defer resp.Body.Close()

	// Vault health endpoint returns different status codes based on state
	// 200 = healthy, 503 = sealed, 501 = not initialized
	if resp.StatusCode == 503 {
		return fmt.Errorf("vault is sealed")
	}
	if resp.StatusCode == 501 {
		return fmt.Errorf("vault is not initialized")
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("vault health check failed with status %d", resp.StatusCode)
	}

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
