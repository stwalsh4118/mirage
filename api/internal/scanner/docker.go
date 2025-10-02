package scanner

import (
	"context"
	"time"
)

// DockerfileInfo represents a discovered Dockerfile with parsed metadata.
type DockerfileInfo struct {
	Path         string   `json:"path"`         // Relative path from repo root (e.g., "services/api/Dockerfile")
	ServiceName  string   `json:"serviceName"`  // Inferred service name (e.g., "api")
	BuildContext string   `json:"buildContext"` // Relative path for build context (e.g., "services/api")
	ExposedPorts []int    `json:"exposedPorts"` // Ports from EXPOSE directives
	BuildArgs    []string `json:"buildArgs"`    // Argument names from ARG directives
	BaseImage    string   `json:"baseImage"`    // Base image from FROM directive
}

// Scanner defines the interface for scanning repositories for Dockerfiles.
type Scanner interface {
	// ScanRepository scans a GitHub repository for Dockerfiles.
	// owner: GitHub username or organization
	// repo: Repository name
	// branch: Branch name to scan (e.g., "main", "master")
	// userToken: Optional GitHub PAT (uses service token if empty)
	ScanRepository(ctx context.Context, owner, repo, branch, userToken string) ([]DockerfileInfo, error)
}

// ScanResult wraps the scan output with metadata.
type ScanResult struct {
	Dockerfiles []DockerfileInfo `json:"dockerfiles"`
	ScannedAt   time.Time        `json:"scannedAt"`
	CacheHit    bool             `json:"cacheHit"`
	Owner       string           `json:"owner"`
	Repo        string           `json:"repo"`
	Branch      string           `json:"branch"`
}
