package vault

import (
	"context"
	"fmt"
	"time"
)

// storeTokenSecret is a generic internal helper for storing token-type secrets.
// This eliminates duplication across Railway, GitHub, and other provider token storage.
// The public methods (StoreRailwayToken, StoreGitHubToken, etc.) wrap this with
// provider-specific validation and logging.
func (c *Client) storeTokenSecret(ctx context.Context, secretPath, secretType, token, userID string) error {
	// Prepare secret data with metadata
	secretData := map[string]interface{}{
		"token": token,
		"metadata": map[string]interface{}{
			"created_by":  userID,
			"created_at":  time.Now().UTC().Format(time.RFC3339),
			"secret_type": secretType,
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
		return err
	}

	return nil
}

// getTokenSecret is a generic internal helper for retrieving token-type secrets.
// This eliminates duplication across Railway, GitHub, and other provider token retrieval.
// The public methods (GetRailwayToken, GetGitHubToken, etc.) wrap this with
// provider-specific logging.
func (c *Client) getTokenSecret(ctx context.Context, secretPath string) (string, error) {
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
		return "", err
	}

	// Extract token from nested data structure
	if response.Data.Data == nil {
		return "", ErrSecretNotFound
	}

	token, ok := response.Data.Data["token"].(string)
	if !ok {
		return "", fmt.Errorf("token not found in secret data")
	}

	return token, nil
}

// deleteTokenSecret is a generic internal helper for deleting token-type secrets.
// This eliminates duplication across Railway, GitHub, and other provider token deletion.
// The public methods (DeleteRailwayToken, DeleteGitHubToken, etc.) wrap this with
// provider-specific logging.
// This performs a soft delete - versions are preserved for audit purposes.
func (c *Client) deleteTokenSecret(ctx context.Context, secretPath string) error {
	// Delete from Vault using the KV v2 data endpoint
	// Format: /v1/{mount}/data/{path}
	kvPath := fmt.Sprintf("/v1/%s/data/%s", c.mountPath, secretPath)

	err := c.makeRequest(ctx, "DELETE", kvPath, nil, nil)
	if err != nil {
		// Check if it's a 404 - secret not found
		if isNotFoundError(err) {
			return ErrSecretNotFound
		}
		return err
	}

	return nil
}
