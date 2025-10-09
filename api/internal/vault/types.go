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
