package vault

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

// Client is an HTTP client wrapper for interacting with Vault's HTTP API
type Client struct {
	address    string
	token      string
	httpClient *http.Client
	namespace  string
	mountPath  string
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
		address:    cfg.Address,
		token:      cfg.Token,
		httpClient: httpClient,
		namespace:  cfg.Namespace,
		mountPath:  cfg.MountPath,
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
	req.Header.Set("X-Vault-Token", c.token)
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
