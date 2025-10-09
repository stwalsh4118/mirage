package vault

import "strings"

// Path components for Vault KV v2 secret paths
const (
	// PathUsers is the top-level path component for user secrets
	PathUsers = "users"

	// PathRailway is the path component for Railway API tokens
	PathRailway = "railway"

	// PathGitHub is the path component for GitHub personal access tokens
	PathGitHub = "github"

	// PathDocker is the path component for Docker registry credentials
	PathDocker = "docker"

	// PathEnvVars is the path component for environment-specific secrets
	PathEnvVars = "env_vars"

	// PathCustom is the path component for custom user-defined secrets
	PathCustom = "custom"
)

// Secret type constants for metadata
const (
	// SecretTypeRailway identifies Railway API tokens
	SecretTypeRailway = "railway_token"

	// SecretTypeGitHub identifies GitHub personal access tokens
	SecretTypeGitHub = "github_pat"

	// SecretTypeDocker identifies Docker registry credentials
	SecretTypeDocker = "docker_credentials"

	// SecretTypeEnvironment identifies environment-specific secrets
	SecretTypeEnvironment = "environment_secret"

	// SecretTypeCustom identifies custom user-defined secrets
	SecretTypeCustom = "custom"
)

// BuildSecretPath constructs a Vault KV v2 secret path from components.
// Note: This does NOT include the mount path, as the KV v2 API format is /v1/{mount}/data/{path}
// Example: BuildSecretPath("users", "user-123", "railway") -> "users/user-123/railway"
func BuildSecretPath(components ...string) string {
	return strings.Join(components, "/")
}

// BuildUserSecretPath constructs a path for a user-specific secret.
// Note: Does not include mount path - use with KV v2 API endpoint /v1/{mount}/data/{path}
// Example: BuildUserSecretPath("user-123", "railway") -> "users/user-123/railway"
func BuildUserSecretPath(userID string, secretPath ...string) string {
	components := append([]string{PathUsers, userID}, secretPath...)
	return BuildSecretPath(components...)
}

// BuildRailwayTokenPath constructs the path for a user's Railway API token.
// Note: Does not include mount path - use with KV v2 API endpoint /v1/{mount}/data/{path}
// Example: BuildRailwayTokenPath("user-123") -> "users/user-123/railway"
func BuildRailwayTokenPath(userID string) string {
	return BuildUserSecretPath(userID, PathRailway)
}

// BuildGitHubTokenPath constructs the path for a user's GitHub PAT.
// Note: Does not include mount path - use with KV v2 API endpoint /v1/{mount}/data/{path}
// Example: BuildGitHubTokenPath("user-123") -> "users/user-123/github"
func BuildGitHubTokenPath(userID string) string {
	return BuildUserSecretPath(userID, PathGitHub)
}

// BuildDockerCredentialsPath constructs the path for Docker registry credentials.
// Note: Does not include mount path - use with KV v2 API endpoint /v1/{mount}/data/{path}
// Example: BuildDockerCredentialsPath("user-123", "docker.io") -> "users/user-123/docker/docker.io"
func BuildDockerCredentialsPath(userID, registry string) string {
	return BuildUserSecretPath(userID, PathDocker, registry)
}

// BuildEnvironmentSecretPath constructs the path for an environment-specific secret.
// Note: Does not include mount path - use with KV v2 API endpoint /v1/{mount}/data/{path}
// Example: BuildEnvironmentSecretPath("user-123", "env-456") -> "users/user-123/env_vars/env-456"
func BuildEnvironmentSecretPath(userID, envID string) string {
	return BuildUserSecretPath(userID, PathEnvVars, envID)
}

// BuildCustomSecretPath constructs the path for a custom user secret.
// Note: Does not include mount path - use with KV v2 API endpoint /v1/{mount}/data/{path}
// Example: BuildCustomSecretPath("user-123", "my-secret") -> "users/user-123/custom/my-secret"
func BuildCustomSecretPath(userID, key string) string {
	return BuildUserSecretPath(userID, PathCustom, key)
}
