package vault

import "fmt"

// Config holds configuration for the Vault client
type Config struct {
	// Address is the Vault server address (e.g., http://localhost:8200)
	Address string
	// Token is the Vault authentication token (for dev mode)
	Token string
	// RoleID is the AppRole role ID (for production AppRole authentication)
	RoleID string
	// SecretID is the AppRole secret ID (for production AppRole authentication)
	SecretID string
	// Namespace is the Vault namespace (optional, for Vault Enterprise)
	Namespace string
	// SkipVerify controls TLS certificate verification (dev only, should be false in production)
	SkipVerify bool
	// MountPath is the KV v2 secrets engine mount path (default: "mirage")
	MountPath string
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Address == "" {
		return fmt.Errorf("vault address is required")
	}
	if c.Token == "" && (c.RoleID == "" || c.SecretID == "") {
		return fmt.Errorf("either token (dev) or role_id+secret_id (prod) must be provided")
	}
	if c.MountPath == "" {
		c.MountPath = DefaultMountPath
	}
	return nil
}
