package controller

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/stwalsh4118/mirageapi/internal/auth"
	"github.com/stwalsh4118/mirageapi/internal/scanner"
	"github.com/stwalsh4118/mirageapi/internal/vault"
)

// DockerfileScanner defines the interface for scanning repositories for Dockerfiles.
type DockerfileScanner interface {
	ScanRepository(ctx context.Context, owner, repo, branch, userToken string) ([]scanner.DockerfileInfo, error)
}

// DiscoveryController handles Dockerfile discovery endpoints.
type DiscoveryController struct {
	Scanner DockerfileScanner
	Vault   VaultClient
}

// VaultClient defines the interface for fetching GitHub tokens from Vault
type VaultClient interface {
	GetGitHubToken(ctx context.Context, userID string) (string, error)
}

// RegisterRoutes registers discovery-related routes under the provided router group.
func (c *DiscoveryController) RegisterRoutes(r *gin.RouterGroup) {
	discovery := r.Group("/discovery")
	{
		discovery.POST("/dockerfiles", c.DiscoverDockerfiles)
	}
}

// DiscoverDockerfilesRequest represents the request to scan a repository for Dockerfiles.
type DiscoverDockerfilesRequest struct {
	// GitHub repository owner (username or organization)
	Owner string `json:"owner" binding:"required"`

	// GitHub repository name
	Repo string `json:"repo" binding:"required"`

	// Branch to scan (e.g., "main", "master")
	Branch string `json:"branch" binding:"required"`

	// Optional user GitHub token for private repositories or higher rate limits
	UserToken string `json:"userToken,omitempty"`

	// Optional request tracking ID
	RequestID string `json:"requestId,omitempty"`
}

// ServiceDTO represents a discovered service with its Dockerfile metadata.
type ServiceDTO struct {
	Name           string   `json:"name"`
	DockerfilePath string   `json:"dockerfilePath"`
	BuildContext   string   `json:"buildContext"`
	ExposedPorts   []int    `json:"exposedPorts"`
	BuildArgs      []string `json:"buildArgs"`
	BaseImage      string   `json:"baseImage,omitempty"`
}

// DiscoverDockerfilesResponse represents the response containing discovered services.
type DiscoverDockerfilesResponse struct {
	Services []ServiceDTO `json:"services"`
	Owner    string       `json:"owner"`
	Repo     string       `json:"repo"`
	Branch   string       `json:"branch"`
	CacheHit bool         `json:"cacheHit,omitempty"`
}

// DiscoverDockerfiles handles POST /api/v1/discovery/dockerfiles requests.
// It scans a GitHub repository for Dockerfiles and returns metadata about discovered services.
func (c *DiscoveryController) DiscoverDockerfiles(ctx *gin.Context) {
	if c.Scanner == nil {
		log.Error().Msg("scanner not configured")
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "dockerfile scanner not configured"})
		return
	}

	var req DiscoverDockerfilesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("invalid request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate and sanitize inputs
	if err := validateDiscoveryRequest(&req); err != nil {
		log.Warn().Err(err).Interface("request", req).Msg("request validation failed")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get authenticated user
	user, err := auth.GetCurrentUser(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to get current user for discovery")
		// Continue without authentication - will use provided token or no token
	}

	// Determine which GitHub token to use (priority order):
	// 1. Token provided in request body (for override/testing)
	// 2. User's token from Vault (if authenticated and configured)
	// 3. No token (public repositories only)
	githubToken := req.UserToken
	tokenSource := "request"

	if githubToken == "" && user != nil && c.Vault != nil {
		// Try to fetch user's GitHub token from Vault
		vaultToken, err := c.Vault.GetGitHubToken(ctx, user.ID)
		if err == nil {
			githubToken = vaultToken
			tokenSource = "vault"
			log.Debug().
				Str("user_id", user.ID).
				Msg("using github token from vault for discovery")
		} else if err != vault.ErrSecretNotFound {
			// Log non-404 errors but continue
			log.Warn().
				Err(err).
				Str("user_id", user.ID).
				Msg("failed to fetch github token from vault, continuing without token")
		}
	}

	if githubToken == "" {
		tokenSource = "none"
	}

	log.Info().
		Str("owner", req.Owner).
		Str("repo", req.Repo).
		Str("branch", req.Branch).
		Str("request_id", req.RequestID).
		Str("token_source", tokenSource).
		Bool("authenticated", user != nil).
		Msg("dockerfile discovery requested")

	// Scan repository
	dockerfiles, err := c.Scanner.ScanRepository(ctx, req.Owner, req.Repo, req.Branch, githubToken)
	if err != nil {
		log.Error().
			Err(err).
			Str("owner", req.Owner).
			Str("repo", req.Repo).
			Str("branch", req.Branch).
			Msg("dockerfile scan failed")
		ctx.JSON(http.StatusBadGateway, gin.H{
			"error": fmt.Sprintf("failed to scan repository: %s", err.Error()),
		})
		return
	}

	// Transform scanner results to API response
	services := make([]ServiceDTO, 0, len(dockerfiles))
	for _, df := range dockerfiles {
		services = append(services, ServiceDTO{
			Name:           df.ServiceName,
			DockerfilePath: df.Path,
			BuildContext:   df.BuildContext,
			ExposedPorts:   df.ExposedPorts,
			BuildArgs:      df.BuildArgs,
			BaseImage:      df.BaseImage,
		})
	}

	log.Info().
		Str("owner", req.Owner).
		Str("repo", req.Repo).
		Str("branch", req.Branch).
		Int("services_found", len(services)).
		Msg("dockerfile discovery completed")

	resp := DiscoverDockerfilesResponse{
		Services: services,
		Owner:    req.Owner,
		Repo:     req.Repo,
		Branch:   req.Branch,
	}

	ctx.JSON(http.StatusOK, resp)
}

// validateDiscoveryRequest validates and sanitizes the discovery request.
func validateDiscoveryRequest(req *DiscoverDockerfilesRequest) error {
	// Trim whitespace
	req.Owner = strings.TrimSpace(req.Owner)
	req.Repo = strings.TrimSpace(req.Repo)
	req.Branch = strings.TrimSpace(req.Branch)

	// Validate owner format (alphanumeric, hyphens, underscores)
	if !isValidGitHubIdentifier(req.Owner) {
		return fmt.Errorf("invalid owner format: %s (must be alphanumeric with hyphens/underscores)", req.Owner)
	}

	// Validate repo format
	if !isValidGitHubIdentifier(req.Repo) {
		return fmt.Errorf("invalid repo format: %s (must be alphanumeric with hyphens/underscores/dots)", req.Repo)
	}

	// Validate branch format (allow slashes for refs like 'feature/my-branch')
	if !isValidBranchName(req.Branch) {
		return fmt.Errorf("invalid branch format: %s", req.Branch)
	}

	// Empty strings after trimming
	if req.Owner == "" || req.Repo == "" || req.Branch == "" {
		return fmt.Errorf("owner, repo, and branch are required")
	}

	return nil
}

// isValidGitHubIdentifier checks if a string is a valid GitHub username/org/repo name.
// Allows alphanumeric characters, hyphens, underscores, and dots.
func isValidGitHubIdentifier(s string) bool {
	if len(s) == 0 || len(s) > 100 {
		return false
	}

	for _, r := range s {
		if !((r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '-' || r == '_' || r == '.') {
			return false
		}
	}

	return true
}

// isValidBranchName checks if a string is a valid Git branch name.
// Allows alphanumeric characters, hyphens, underscores, dots, and forward slashes.
func isValidBranchName(s string) bool {
	if len(s) == 0 || len(s) > 255 {
		return false
	}

	// Branch names can't start with a dot or slash
	if s[0] == '.' || s[0] == '/' {
		return false
	}

	for _, r := range s {
		if !((r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '-' || r == '_' || r == '.' || r == '/') {
			return false
		}
	}

	return true
}
