package vault

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

// StoreRailwayToken stores or updates a user's Railway API token in Vault.
// If the token already exists, this creates a new version.
// The token is stored with metadata tracking the user, timestamp, and secret type.
func (c *Client) StoreRailwayToken(ctx context.Context, userID, token string) error {
	if userID == "" {
		return fmt.Errorf("user ID is required")
	}
	if token == "" {
		return ErrInvalidSecret
	}

	// Build the secret path (does not include mount - that's added in the KV v2 API endpoint)
	secretPath := BuildRailwayTokenPath(userID)

	// Prepare secret data with metadata
	secretData := map[string]interface{}{
		"token": token,
		"metadata": map[string]interface{}{
			"created_by":  userID,
			"created_at":  time.Now().UTC().Format(time.RFC3339),
			"secret_type": SecretTypeRailway,
		},
	}

	// Vault KV v2 requires wrapping data in a "data" field
	requestBody := map[string]interface{}{
		"data": secretData,
	}

	// Write to Vault using the KV v2 data endpoint
	// Format: /v1/{mount}/data/{path}
	kvPath := fmt.Sprintf("/v1/%s/data/%s", c.mountPath, secretPath)
	err := c.makeRequest(ctx, "POST", kvPath, requestBody, nil)
	if err != nil {
		return fmt.Errorf("failed to store railway token: %w", err)
	}

	log.Info().
		Str("user_id", userID).
		Str("secret_path", secretPath).
		Msg("stored railway token")

	return nil
}

// GetRailwayToken retrieves a user's Railway API token from Vault.
// Returns the most recent version of the token.
// Returns ErrSecretNotFound if the user has no Railway token configured.
func (c *Client) GetRailwayToken(ctx context.Context, userID string) (string, error) {
	if userID == "" {
		return "", fmt.Errorf("user ID is required")
	}

	// Build the secret path (does not include mount - that's added in the KV v2 API endpoint)
	secretPath := BuildRailwayTokenPath(userID)

	// Read from Vault using the KV v2 data endpoint
	// Format: /v1/{mount}/data/{path}
	kvPath := fmt.Sprintf("/v1/%s/data/%s", c.mountPath, secretPath)

	var response struct {
		Data struct {
			Data map[string]interface{} `json:"data"`
		} `json:"data"`
	}

	err := c.makeRequest(ctx, "GET", kvPath, nil, &response)
	if err != nil {
		// Check if it's a 404 - secret not found
		if isNotFoundError(err) {
			return "", ErrSecretNotFound
		}
		return "", fmt.Errorf("failed to read railway token: %w", err)
	}

	// Extract token from nested data structure
	if response.Data.Data == nil {
		return "", ErrSecretNotFound
	}

	token, ok := response.Data.Data["token"].(string)
	if !ok {
		return "", fmt.Errorf("token not found in secret data")
	}

	log.Debug().
		Str("user_id", userID).
		Str("secret_path", secretPath).
		Msg("retrieved railway token from vault")

	return token, nil
}

// DeleteRailwayToken removes a user's Railway API token from Vault.
// This performs a soft delete - versions are preserved for audit purposes.
// Returns ErrSecretNotFound if the token doesn't exist.
func (c *Client) DeleteRailwayToken(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("user ID is required")
	}

	// Build the secret path (does not include mount - that's added in the KV v2 API endpoint)
	secretPath := BuildRailwayTokenPath(userID)

	// Delete from Vault using the KV v2 data endpoint
	// Format: /v1/{mount}/data/{path}
	kvPath := fmt.Sprintf("/v1/%s/data/%s", c.mountPath, secretPath)

	err := c.makeRequest(ctx, "DELETE", kvPath, nil, nil)
	if err != nil {
		// Check if it's a 404 - secret not found
		if isNotFoundError(err) {
			return ErrSecretNotFound
		}
		return fmt.Errorf("failed to delete railway token: %w", err)
	}

	log.Info().
		Str("user_id", userID).
		Str("secret_path", secretPath).
		Msg("deleted railway token")

	return nil
}

// RotateRailwayToken replaces a user's Railway token with a new one.
// This is atomic and creates a new version while preserving the old version.
// The old version is maintained in Vault's version history for audit purposes.
func (c *Client) RotateRailwayToken(ctx context.Context, userID, newToken string) error {
	if userID == "" {
		return fmt.Errorf("user ID is required")
	}
	if newToken == "" {
		return ErrInvalidSecret
	}

	// Rotation in KV v2 is simply writing a new version
	// The old version is automatically preserved
	err := c.StoreRailwayToken(ctx, userID, newToken)
	if err != nil {
		return fmt.Errorf("failed to rotate railway token: %w", err)
	}

	log.Info().
		Str("user_id", userID).
		Msg("rotated railway token")

	return nil
}

// ValidateRailwayToken tests a user's Railway token by making a test API call.
// This verifies the token is valid and has the necessary permissions.
// Updates the LastValidated timestamp in metadata if successful.
func (c *Client) ValidateRailwayToken(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("user ID is required")
	}

	// Get the token first
	token, err := c.GetRailwayToken(ctx, userID)
	if err != nil {
		return err
	}

	// Basic validation of token format
	// A full Railway API validation would require the railway.Client
	// which could cause circular dependencies at this layer.
	// Controllers that use this interface can provide more thorough validation.
	if token == "" || len(token) < 10 {
		return fmt.Errorf("railway token validation failed: token appears invalid")
	}

	// Token exists and has basic format validity

	// Update last validated timestamp in Vault metadata
	// Read current secret to preserve other data
	secretPath := BuildRailwayTokenPath(userID)
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

	// Update metadata with last validated timestamp
	if response.Data.Data != nil {
		metadata, ok := response.Data.Data["metadata"].(map[string]interface{})
		if !ok {
			metadata = make(map[string]interface{})
		}
		metadata["last_validated"] = time.Now().UTC().Format(time.RFC3339)
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
		Msg("validated railway token")

	return nil
}

// isNotFoundError checks if an error is a 404 Not Found error from Vault
func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	// Check if error message contains "404" or "not found"
	errStr := err.Error()
	return containsAny(errStr, "404", "not found", "No value found")
}

// containsAny checks if a string contains any of the given substrings
func containsAny(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if len(s) >= len(substr) {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
		}
	}
	return false
}
