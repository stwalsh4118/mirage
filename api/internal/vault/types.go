package vault

import "time"

const (
	// DefaultMountPath is the default KV v2 mount path for Mirage secrets
	DefaultMountPath = "mirage"
	// DefaultTimeout is the default HTTP client timeout for Vault requests
	DefaultTimeout = 30 * time.Second
)
