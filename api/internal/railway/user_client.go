package railway

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/stwalsh4118/mirageapi/internal/vault"
)

// GetUserRailwayClient creates a Railway client for a specific user
// by fetching their Railway token from Vault.
//
// This function is used when you want to make Railway API calls on behalf
// of a specific user using their personal Railway API token stored in Vault.
//
// Parameters:
//   - ctx: Context for the operation
//   - userID: The ID of the user whose Railway token to fetch
//   - vaultClient: The Vault client to use for fetching the token
//   - endpoint: Railway API endpoint (use empty string for default)
//   - httpc: HTTP client to use (use nil for default)
//
// Returns:
//   - *Client: A Railway client configured with the user's token
//   - error: An error if the token cannot be fetched or is invalid
//
// Errors:
//   - ErrVaultRequired: If vaultClient is nil
//   - ErrNoRailwayToken: If the user has no Railway token in Vault
//   - Other errors from Vault operations
//
// Example:
//
//	client, err := GetUserRailwayClient(ctx, user.ID, vaultClient, "", nil)
//	if err != nil {
//	    return fmt.Errorf("failed to get user railway client: %w", err)
//	}
//	// Use client for Railway API calls...
func GetUserRailwayClient(
	ctx context.Context,
	userID string,
	vaultClient *vault.Client,
	endpoint string,
	httpc *http.Client,
) (*Client, error) {
	if vaultClient == nil {
		return nil, ErrVaultRequired
	}

	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	// Fetch the user's Railway token from Vault
	token, err := vaultClient.GetRailwayToken(ctx, userID)
	if err != nil {
		// Check if the error is because the secret wasn't found
		if err == vault.ErrSecretNotFound {
			return nil, ErrNoRailwayToken
		}
		return nil, fmt.Errorf("failed to get railway token for user %s: %w", userID, err)
	}

	if token == "" {
		return nil, ErrNoRailwayToken
	}

	// Create a new Railway client with the user's token
	client := NewClient(endpoint, token, httpc)

	log.Debug().
		Str("user_id", userID).
		Msg("created user-specific railway client")

	return client, nil
}

// GetRailwayClientForUser returns a user-specific Railway client if Vault
// is enabled, otherwise returns the global Railway client as a fallback.
//
// This function provides a graceful degradation path for environments where
// Vault is not configured or during migration from global token to per-user tokens.
//
// Parameters:
//   - ctx: Context for the operation
//   - userID: The ID of the user whose Railway token to fetch
//   - vaultClient: The Vault client (can be nil if Vault is disabled)
//   - globalClient: The global Railway client to use as fallback (accepts interface or concrete type)
//
// Returns:
//   - *Client: A user-specific client if Vault is enabled, otherwise the global client
//   - error: An error if Vault is enabled but the token cannot be fetched
//
// Behavior:
//   - If vaultClient is nil: Returns globalClient with a debug log
//   - If vaultClient is not nil: Attempts to get user-specific client
//   - If user has no token: Returns error (no fallback when Vault is enabled)
//
// Example:
//
//	client, err := GetRailwayClientForUser(ctx, user.ID, vaultClient, globalClient)
//	if err != nil {
//	    return fmt.Errorf("failed to get railway client: %w", err)
//	}
//	// Use client for Railway API calls...
func GetRailwayClientForUser(
	ctx context.Context,
	userID string,
	vaultClient *vault.Client,
	globalClient interface{},
) (*Client, error) {
	if vaultClient != nil {
		// Vault is enabled, try to get user-specific client
		// Extract endpoint from global client to maintain consistency
		var endpoint string

		// Type assert the global client to get endpoint
		concreteClient, ok := globalClient.(*Client)
		if !ok {
			// Fail fast if globalClient is not a valid *Client
			return nil, fmt.Errorf("invalid global railway client: expected *Client")
		}
		endpoint = concreteClient.endpoint

		client, err := GetUserRailwayClient(ctx, userID, vaultClient, endpoint, nil)
		if err != nil {
			return nil, err
		}

		return client, nil
	}

	// Vault is not enabled, fall back to global client
	log.Debug().
		Str("user_id", userID).
		Msg("vault disabled, using global railway client")

	// Type assert the interface to concrete Client type
	if concreteClient, ok := globalClient.(*Client); ok {
		return concreteClient, nil
	}

	// If we can't assert to concrete type, return error
	return nil, fmt.Errorf("global client is not a valid Railway client")
}
