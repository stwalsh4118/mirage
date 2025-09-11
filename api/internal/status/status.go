package status

// Canonical UI status constants used across the system.
const (
	StatusActive     = "active"
	StatusCreating   = "creating"
	StatusDestroying = "destroying"
	StatusError      = "error"
	StatusUnknown    = "unknown"
)

// NormalizeRemoteToUI maps various remote (Railway) status strings into the
// canonical UI statuses. Any unrecognized value maps to "unknown".
func NormalizeRemoteToUI(s string) string {
	switch s {
	case "ready", "healthy", "running", "active":
		return StatusActive
	case "creating", "provisioning", "deploying", "queued":
		return StatusCreating
	case "destroying", "deleting", "terminating":
		return StatusDestroying
	case "error", "failed", "degraded", "crashed":
		return StatusError
	case "", "unknown":
		return StatusUnknown
	default:
		return StatusUnknown
	}
}

// NormalizeLocalToUI ensures any locally stored status is expressed in the
// canonical UI set. It accepts already-canonical values and common synonyms.
func NormalizeLocalToUI(s string) string {
	switch s {
	case StatusActive, StatusCreating, StatusDestroying, StatusError, StatusUnknown:
		return s
	case "ready", "healthy", "running":
		return StatusActive
	case "provisioning", "deploying":
		return StatusCreating
	case "deleting", "terminating":
		return StatusDestroying
	case "failed", "degraded", "crashed":
		return StatusError
	case "":
		return StatusUnknown
	default:
		return StatusUnknown
	}
}
