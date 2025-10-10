package railway

import "errors"

// Error types for Railway client operations
var (
	// ErrNoRailwayToken is returned when a user has no Railway token configured in Vault
	ErrNoRailwayToken = errors.New("no railway token configured for user")

	// ErrVaultDisabled is returned when attempting to get user-specific client with Vault disabled
	ErrVaultDisabled = errors.New("vault is disabled, cannot fetch user token")

	// ErrVaultRequired is returned when Vault is required but not provided
	ErrVaultRequired = errors.New("vault not configured, cannot get user railway client")
)
