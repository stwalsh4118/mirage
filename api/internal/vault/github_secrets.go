package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// StoreGitHubToken stores or updates a user's GitHub Personal Access Token in Vault.
// If the token already exists, this creates a new version.
// The token is stored with metadata tracking the user, timestamp, and secret type.
func (c *Client) StoreGitHubToken(ctx context.Context, userID, token string) error {
	if userID == "" {
		return fmt.Errorf("user ID is required")
	}
	if token == "" {
		return ErrInvalidSecret
	}

	// Build the secret path (does not include mount - that's added in the KV v2 API endpoint)
	secretPath := BuildGitHubTokenPath(userID)

	// Use generic helper to store the token
	err := c.storeTokenSecret(ctx, secretPath, SecretTypeGitHub, token, userID)
	if err != nil {
		return fmt.Errorf("failed to store github token: %w", err)
	}

	log.Info().
		Str("user_id", userID).
		Str("secret_path", secretPath).
		Msg("stored github token")

	return nil
}

// GetGitHubToken retrieves a user's GitHub Personal Access Token from Vault.
// Returns the most recent version of the token.
// Returns ErrSecretNotFound if the user has no GitHub token configured.
func (c *Client) GetGitHubToken(ctx context.Context, userID string) (string, error) {
	if userID == "" {
		return "", fmt.Errorf("user ID is required")
	}

	// Build the secret path (does not include mount - that's added in the KV v2 API endpoint)
	secretPath := BuildGitHubTokenPath(userID)

	// Use generic helper to retrieve the token
	token, err := c.getTokenSecret(ctx, secretPath)
	if err != nil {
		if err == ErrSecretNotFound {
			return "", ErrSecretNotFound
		}
		return "", fmt.Errorf("failed to read github token: %w", err)
	}

	log.Debug().
		Str("user_id", userID).
		Str("secret_path", secretPath).
		Msg("retrieved github token from vault")

	return token, nil
}

// DeleteGitHubToken removes a user's GitHub Personal Access Token from Vault.
// This performs a soft delete - versions are preserved for audit purposes.
// Returns ErrSecretNotFound if the token doesn't exist.
func (c *Client) DeleteGitHubToken(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("user ID is required")
	}

	// Build the secret path (does not include mount - that's added in the KV v2 API endpoint)
	secretPath := BuildGitHubTokenPath(userID)

	// Use generic helper to delete the token
	err := c.deleteTokenSecret(ctx, secretPath)
	if err != nil {
		if err == ErrSecretNotFound {
			return ErrSecretNotFound
		}
		return fmt.Errorf("failed to delete github token: %w", err)
	}

	log.Info().
		Str("user_id", userID).
		Str("secret_path", secretPath).
		Msg("deleted github token")

	return nil
}

// ValidateGitHubToken validates a GitHub token by calling the GitHub API.
// It returns the GitHub username and token scopes if valid.
// This is exposed for use by controllers that need to validate tokens before storing.
func (c *Client) ValidateGitHubToken(ctx context.Context, token string) (username string, scopes []string, err error) {
	if token == "" {
		return "", nil, fmt.Errorf("token is required")
	}

	// Create request to GitHub user API
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create github api request: %w", err)
	}

	// GitHub API requires Bearer token authentication
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	// Create HTTP client with timeout to prevent indefinite hangs
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Execute request
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("github api request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("github token invalid: status %d", resp.StatusCode)
	}

	// Parse user response
	var user struct {
		Login string `json:"login"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return "", nil, fmt.Errorf("failed to decode github user response: %w", err)
	}

	// Extract scopes from response headers
	// GitHub returns scopes in X-OAuth-Scopes header as comma-separated values
	// Handle both "repo, read:org" and "repo,read:org" formats
	scopesHeader := resp.Header.Get("X-OAuth-Scopes")
	if scopesHeader != "" {
		for _, scope := range strings.Split(scopesHeader, ",") {
			scope = strings.TrimSpace(scope)
			if scope != "" {
				scopes = append(scopes, scope)
			}
		}
	}

	log.Debug().
		Str("username", user.Login).
		Strs("scopes", scopes).
		Msg("validated github token")

	return user.Login, scopes, nil
}

// ValidateGitHubTokenAndUpdateMetadata validates a GitHub token with the GitHub API and updates its metadata in Vault.
// This is similar to ValidateRailwayToken - validates the token exists and updates last_validated timestamp.
func (c *Client) ValidateGitHubTokenAndUpdateMetadata(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("user ID is required")
	}

	// Get the token first
	token, err := c.GetGitHubToken(ctx, userID)
	if err != nil {
		return err
	}

	// Validate with GitHub API
	username, scopes, err := c.ValidateGitHubToken(ctx, token)
	if err != nil {
		return fmt.Errorf("github token validation failed: %w", err)
	}

	// Update last validated timestamp in Vault metadata
	// Read current secret to preserve other data
	secretPath := BuildGitHubTokenPath(userID)
	kvPath := fmt.Sprintf("/v1/%s/data/%s", c.mountPath, secretPath)

	var response struct {
		Data struct {
			Data map[string]interface{} `json:"data"`
		} `json:"data"`
	}

	err = c.makeRequest(ctx, "GET", kvPath, nil, &response)
	if err != nil {
		return fmt.Errorf("failed to read secret for validation update: %w", err)
	}

	// Update metadata with last validated timestamp and validation info
	if response.Data.Data != nil {
		metadata, ok := response.Data.Data["metadata"].(map[string]interface{})
		if !ok {
			metadata = make(map[string]interface{})
		}
		metadata["last_validated"] = time.Now().UTC().Format(time.RFC3339)
		metadata["github_username"] = username
		metadata["github_scopes"] = scopes
		response.Data.Data["metadata"] = metadata

		// Write back the updated secret
		requestBody := map[string]interface{}{
			"data": response.Data.Data,
		}

		err = c.makeRequest(ctx, "POST", kvPath, requestBody, nil)
		if err != nil {
			log.Warn().
				Err(err).
				Str("user_id", userID).
				Msg("failed to update last_validated timestamp")
			// Don't fail the validation if we can't update the timestamp
		}
	}

	log.Info().
		Str("user_id", userID).
		Str("github_username", username).
		Msg("validated github token")

	return nil
}
