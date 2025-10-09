package vault

import "time"

const (
	// DefaultMountPath is the default KV v2 mount path for Mirage secrets
	DefaultMountPath = "mirage"
	// DefaultTimeout is the default HTTP client timeout for Vault requests
	DefaultTimeout = 30 * time.Second
	// DefaultRenewalBuffer is how much time before expiry we should renew the token (50% of TTL)
	DefaultRenewalBuffer = 0.5
)

// TokenLookupResponse represents the response from token lookup-self endpoint
type TokenLookupResponse struct {
	Data TokenInfo `json:"data"`
}

// TokenInfo contains information about a Vault token
type TokenInfo struct {
	Accessor       string   `json:"accessor"`
	CreationTime   int64    `json:"creation_time"`
	CreationTTL    int      `json:"creation_ttl"`
	DisplayName    string   `json:"display_name"`
	EntityID       string   `json:"entity_id"`
	ExpireTime     string   `json:"expire_time"`
	ExplicitMaxTTL int      `json:"explicit_max_ttl"`
	ID             string   `json:"id"`
	IssueTime      string   `json:"issue_time"`
	NumUses        int      `json:"num_uses"`
	Orphan         bool     `json:"orphan"`
	Path           string   `json:"path"`
	Policies       []string `json:"policies"`
	Renewable      bool     `json:"renewable"`
	TTL            int      `json:"ttl"`
}

// AppRoleLoginRequest represents the request body for AppRole login
type AppRoleLoginRequest struct {
	RoleID   string `json:"role_id"`
	SecretID string `json:"secret_id"`
}

// AppRoleLoginResponse represents the response from AppRole login
type AppRoleLoginResponse struct {
	Auth AuthInfo `json:"auth"`
}

// AuthInfo contains authentication information from Vault
type AuthInfo struct {
	ClientToken   string   `json:"client_token"`
	Accessor      string   `json:"accessor"`
	Policies      []string `json:"policies"`
	TokenPolicies []string `json:"token_policies"`
	LeaseDuration int      `json:"lease_duration"`
	Renewable     bool     `json:"renewable"`
}

// TokenRenewalRequest represents the request body for token renewal
type TokenRenewalRequest struct {
	Increment string `json:"increment,omitempty"`
}

// TokenRenewalResponse represents the response from token renewal
type TokenRenewalResponse struct {
	Auth AuthInfo `json:"auth"`
}

// VaultStatus represents the health status of the Vault server
type VaultStatus struct {
	Available   bool   `json:"available"`
	Initialized bool   `json:"initialized"`
	Sealed      bool   `json:"sealed"`
	Version     string `json:"version"`
	ClusterID   string `json:"cluster_id"`
	ClusterName string `json:"cluster_name"`
}

// CircuitState represents the state of the circuit breaker
type CircuitState int

const (
	// CircuitClosed means requests are allowed through
	CircuitClosed CircuitState = iota
	// CircuitOpen means requests are blocked due to failures
	CircuitOpen
	// CircuitHalfOpen means trying to recover, allowing limited requests
	CircuitHalfOpen
)

// String returns the string representation of CircuitState
func (s CircuitState) String() string {
	switch s {
	case CircuitClosed:
		return "closed"
	case CircuitOpen:
		return "open"
	case CircuitHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

const (
	// DefaultFailureThreshold is the number of failures before opening the circuit
	DefaultFailureThreshold = 5
	// DefaultCircuitTimeout is the time before transitioning to half-open
	DefaultCircuitTimeout = 30 * time.Second
	// DefaultSuccessThreshold is the number of successes needed in half-open to close
	DefaultSuccessThreshold = 2
	// DefaultHealthCheckInterval is the interval for periodic health checks
	DefaultHealthCheckInterval = 30 * time.Second
)
