package vault

import "context"

// SecretStore defines the interface for all secret management operations in Mirage.
// This interface abstracts Vault operations and provides a clean contract for secret storage,
// retrieval, versioning, and lifecycle management.
//
// Implementations of this interface should:
//   - Ensure all secrets are encrypted at rest in Vault
//   - Namespace secrets by user ID to enforce isolation
//   - Handle Vault unavailability gracefully with circuit breaking
//   - Log all operations for audit purposes (without logging secret values)
//   - Return appropriate errors from the errors.go file
type SecretStore interface {
	// Railway token management
	//
	// Railway API tokens are stored per-user to enable user-specific Railway operations
	// and proper attribution of Railway actions.

	// StoreRailwayToken stores or updates a user's Railway API token.
	// If the token already exists, it creates a new version.
	// Returns an error if the token is invalid or Vault is unavailable.
	StoreRailwayToken(ctx context.Context, userID, token string) error

	// GetRailwayToken retrieves a user's Railway API token.
	// Returns the most recent version of the token.
	// Returns ErrSecretNotFound if the user has no Railway token configured.
	GetRailwayToken(ctx context.Context, userID string) (string, error)

	// DeleteRailwayToken removes a user's Railway API token.
	// This performs a soft delete - versions are preserved for audit purposes.
	// Returns ErrSecretNotFound if the token doesn't exist.
	DeleteRailwayToken(ctx context.Context, userID string) error

	// RotateRailwayToken replaces a user's Railway token with a new one.
	// This is atomic and creates a new version while preserving the old version.
	// Returns an error if the new token is invalid or Vault is unavailable.
	RotateRailwayToken(ctx context.Context, userID, newToken string) error

	// ValidateRailwayToken tests a user's Railway token by making a test API call.
	// Updates the LastValidated timestamp in metadata if successful.
	// Returns an error if the token is invalid or the test call fails.
	ValidateRailwayToken(ctx context.Context, userID string) error

	// GitHub PAT management
	//
	// GitHub Personal Access Tokens are used for accessing private repositories
	// during Dockerfile discovery and repository scanning.

	// StoreGitHubToken stores or updates a user's GitHub Personal Access Token.
	// If the token already exists, it creates a new version.
	// Returns an error if the token is invalid or Vault is unavailable.
	StoreGitHubToken(ctx context.Context, userID, token string) error

	// GetGitHubToken retrieves a user's GitHub Personal Access Token.
	// Returns the most recent version of the token.
	// Returns ErrSecretNotFound if the user has no GitHub token configured.
	GetGitHubToken(ctx context.Context, userID string) (string, error)

	// DeleteGitHubToken removes a user's GitHub Personal Access Token.
	// This performs a soft delete - versions are preserved for audit purposes.
	// Returns ErrSecretNotFound if the token doesn't exist.
	DeleteGitHubToken(ctx context.Context, userID string) error

	// Docker credentials management
	//
	// Docker registry credentials enable pulling from private container registries.
	// Each user can configure multiple registries (docker.io, ghcr.io, etc.).

	// StoreDockerCredentials stores or updates Docker registry credentials.
	// Credentials are stored per-registry, so a user can have credentials for
	// multiple registries (e.g., Docker Hub, GitHub Container Registry).
	// Returns an error if the credentials are invalid or Vault is unavailable.
	StoreDockerCredentials(ctx context.Context, userID string, creds DockerCredentials) error

	// GetDockerCredentials retrieves credentials for a specific registry.
	// Returns ErrSecretNotFound if no credentials exist for the registry.
	GetDockerCredentials(ctx context.Context, userID, registry string) (DockerCredentials, error)

	// ListDockerRegistries returns a list of all registries the user has configured.
	// Returns an empty slice if no registries are configured.
	ListDockerRegistries(ctx context.Context, userID string) ([]string, error)

	// DeleteDockerCredentials removes credentials for a specific registry.
	// Returns ErrSecretNotFound if no credentials exist for the registry.
	DeleteDockerCredentials(ctx context.Context, userID, registry string) error

	// Environment-specific secrets
	//
	// Environment secrets are key-value pairs associated with a specific Railway environment.
	// These can include database URLs, API keys, or any other environment-specific configuration.

	// StoreEnvironmentSecret stores a single environment secret.
	// Environment secrets are namespaced by user ID and environment ID.
	// Returns an error if Vault is unavailable.
	StoreEnvironmentSecret(ctx context.Context, userID, envID, key, value string) error

	// GetEnvironmentSecret retrieves a single environment secret by key.
	// Returns ErrSecretNotFound if the secret doesn't exist.
	GetEnvironmentSecret(ctx context.Context, userID, envID, key string) (string, error)

	// GetAllEnvironmentSecrets retrieves all secrets for an environment.
	// Returns a map of key-value pairs. Returns an empty map if no secrets exist.
	GetAllEnvironmentSecrets(ctx context.Context, userID, envID string) (map[string]string, error)

	// DeleteEnvironmentSecret removes a single environment secret.
	// Returns ErrSecretNotFound if the secret doesn't exist.
	DeleteEnvironmentSecret(ctx context.Context, userID, envID, key string) error

	// BulkStoreEnvironmentSecrets stores multiple environment secrets atomically.
	// This is useful for importing or updating many secrets at once.
	// If any secret fails to store, the entire operation should be rolled back.
	// Returns an error if any secret is invalid or Vault is unavailable.
	BulkStoreEnvironmentSecrets(ctx context.Context, userID, envID string, secrets map[string]string) error

	// Generic secret management
	//
	// Generic secrets provide a flexible way to store any user-defined secrets
	// with custom metadata, tags, and versioning support.

	// StoreSecret stores a generic secret with metadata.
	// This is the most flexible secret storage method, supporting custom keys,
	// values, and metadata for any use case not covered by specific secret types.
	// Returns an error if the secret is invalid or Vault is unavailable.
	StoreSecret(ctx context.Context, userID, key, value string, metadata SecretMetadata) error

	// GetSecret retrieves a complete secret with all metadata.
	// Returns ErrSecretNotFound if the secret doesn't exist.
	GetSecret(ctx context.Context, userID, key string) (Secret, error)

	// GetSecretValue retrieves only the secret value without metadata.
	// This is more efficient when metadata is not needed.
	// Returns ErrSecretNotFound if the secret doesn't exist.
	GetSecretValue(ctx context.Context, userID, key string) (string, error)

	// DeleteSecret removes a generic secret.
	// This performs a soft delete - versions are preserved for audit purposes.
	// Returns ErrSecretNotFound if the secret doesn't exist.
	DeleteSecret(ctx context.Context, userID, key string) error

	// ListSecrets returns metadata for all of a user's secrets.
	// This does not include secret values, only metadata for discovery/listing.
	// Returns an empty slice if the user has no secrets.
	ListSecrets(ctx context.Context, userID string) ([]SecretMetadata, error)

	// Version management
	//
	// Vault KV v2 automatically versions all secrets. These methods provide
	// access to secret history for audit, recovery, and rollback purposes.

	// GetSecretVersion retrieves a specific version of a secret.
	// Returns ErrVersionNotFound if the version doesn't exist.
	GetSecretVersion(ctx context.Context, userID, key string, version int) (Secret, error)

	// ListSecretVersions returns all version numbers for a secret.
	// Versions are returned in descending order (newest first).
	// Returns an empty slice if the secret has no versions.
	ListSecretVersions(ctx context.Context, userID, key string) ([]int, error)

	// RollbackSecret creates a new version with the value from a previous version.
	// This does not delete the current version - it creates a new version
	// with the old value, maintaining a complete audit trail.
	// Returns ErrVersionNotFound if the target version doesn't exist.
	RollbackSecret(ctx context.Context, userID, key string, toVersion int) error

	// Metadata management
	//
	// Metadata operations allow updating secret metadata without changing
	// the secret value itself. This is useful for tagging, categorization,
	// and tracking validation status.

	// UpdateSecretMetadata updates the metadata for a secret without changing its value.
	// Only the metadata fields are updated; the secret value remains unchanged.
	// Returns ErrSecretNotFound if the secret doesn't exist.
	UpdateSecretMetadata(ctx context.Context, userID, key string, metadata SecretMetadata) error

	// GetSecretMetadata retrieves only the metadata for a secret.
	// This is more efficient than GetSecret when only metadata is needed.
	// Returns ErrSecretNotFound if the secret doesn't exist.
	GetSecretMetadata(ctx context.Context, userID, key string) (SecretMetadata, error)

	// Health and status
	//
	// These methods provide health checking and monitoring capabilities
	// for the Vault connection and secret cache.

	// HealthCheck performs a health check on the Vault connection.
	// This should verify Vault is reachable, unsealed, and initialized.
	// Returns an error if Vault is unavailable or unhealthy.
	HealthCheck(ctx context.Context) error

	// GetVaultStatus retrieves detailed status information about Vault.
	// This includes seal status, version, cluster information, etc.
	// Returns an error if Vault is unreachable.
	GetVaultStatus(ctx context.Context) (VaultStatus, error)

	// GetCacheStats returns statistics about the secret cache.
	// This is useful for monitoring cache effectiveness and debugging.
	// Returns empty stats if caching is not enabled.
	GetCacheStats() CacheStats
}
