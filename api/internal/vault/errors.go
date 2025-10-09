package vault

import "errors"

// Standard errors for secret operations
var (
	// ErrSecretNotFound is returned when a secret does not exist
	ErrSecretNotFound = errors.New("secret not found")

	// ErrVersionNotFound is returned when a specific secret version does not exist
	ErrVersionNotFound = errors.New("secret version not found")

	// ErrInvalidSecret is returned when a secret value is invalid or malformed
	ErrInvalidSecret = errors.New("invalid secret value")

	// ErrVaultUnavailable is returned when Vault is unreachable or unhealthy
	ErrVaultUnavailable = errors.New("vault is unavailable")

	// ErrUnauthorized is returned when access to a secret is denied
	ErrUnauthorized = errors.New("unauthorized access to secret")

	// ErrSecretLocked is returned when a secret is locked due to CAS conflict
	ErrSecretLocked = errors.New("secret is locked (CAS conflict)")
)
