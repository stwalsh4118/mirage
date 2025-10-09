package vault

import "time"

// SecretMetadata contains metadata about a secret stored in Vault.
// This metadata is used to track secret lifecycle, ownership, and classification.
type SecretMetadata struct {
	// CreatedBy is the user ID who created the secret
	CreatedBy string `json:"created_by"`

	// CreatedAt is the timestamp when the secret was created
	CreatedAt time.Time `json:"created_at"`

	// LastValidated is the timestamp when the secret was last validated/tested
	// Nil if never validated
	LastValidated *time.Time `json:"last_validated,omitempty"`

	// SecretType identifies the type of secret (railway_token, github_pat, etc.)
	// Should be one of the SecretType* constants
	SecretType string `json:"secret_type"`

	// Tags are user-defined tags for organizing and filtering secrets
	Tags []string `json:"tags,omitempty"`

	// Version is the current version number of the secret
	Version int `json:"version"`
}

// Secret represents a complete secret with its value and metadata.
// This is the primary domain object for secret operations.
type Secret struct {
	// Key is the secret identifier (e.g., "railway", "github", "my-api-key")
	Key string `json:"key"`

	// Value is the actual secret data (e.g., token, password)
	Value string `json:"value"`

	// Metadata contains additional information about the secret
	Metadata SecretMetadata `json:"metadata"`
}

// DockerCredentials contains Docker registry authentication credentials.
// These are stored as structured data in Vault for docker registry access.
type DockerCredentials struct {
	// Registry is the Docker registry URL (e.g., "docker.io", "ghcr.io")
	Registry string `json:"registry"`

	// Username is the registry username
	Username string `json:"username"`

	// Password is the registry password or access token
	Password string `json:"password"`

	// Email is the email address associated with the registry account
	// Optional for most registries
	Email string `json:"email,omitempty"`
}

// CacheStats contains cache performance metrics for the secret cache.
// Used for monitoring and debugging cache effectiveness.
type CacheStats struct {
	// Size is the current number of entries in the cache
	Size int `json:"size"`

	// Hits is the total number of cache hits
	Hits int64 `json:"hits"`

	// Misses is the total number of cache misses
	Misses int64 `json:"misses"`

	// HitRate is the cache hit ratio (hits / (hits + misses))
	HitRate float64 `json:"hit_rate"`
}
