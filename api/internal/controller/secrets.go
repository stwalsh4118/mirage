package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/stwalsh4118/mirageapi/internal/auth"
	"github.com/stwalsh4118/mirageapi/internal/railway"
	"github.com/stwalsh4118/mirageapi/internal/vault"
)

// SecretsController handles secret management operations for Railway tokens,
// GitHub tokens, Docker credentials, and other user secrets stored in Vault.
type SecretsController struct {
	Vault   *vault.Client
	Railway *railway.Client // For validating Railway tokens
}

// RegisterRoutes registers all secret management routes
func (c *SecretsController) RegisterRoutes(r *gin.RouterGroup) {
	secrets := r.Group("/secrets")
	{
		// Railway token management
		secrets.POST("/railway", c.StoreRailwayToken)
		secrets.GET("/railway/status", c.GetRailwayTokenStatus)
		secrets.POST("/railway/validate", c.ValidateRailwayToken)
		secrets.DELETE("/railway", c.DeleteRailwayToken)
		secrets.POST("/railway/rotate", c.RotateRailwayToken)

		// GitHub token management
		secrets.POST("/github", c.StoreGitHubToken)
		secrets.GET("/github/status", c.GetGitHubTokenStatus)
		secrets.POST("/github/validate", c.ValidateGitHubToken)
		secrets.DELETE("/github", c.DeleteGitHubToken)
	}
}

// StoreRailwayTokenRequest is the request payload for storing a Railway token
type StoreRailwayTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// StoreRailwayTokenResponse is the response for successful token storage
type StoreRailwayTokenResponse struct {
	Success   bool   `json:"success"`
	Validated bool   `json:"validated"`
	StoredAt  string `json:"stored_at"`
	Message   string `json:"message"`
}

// RailwayTokenStatusResponse is the response for token status check
type RailwayTokenStatusResponse struct {
	Configured    bool    `json:"configured"`
	LastValidated *string `json:"last_validated,omitempty"`
	NeedsRotation bool    `json:"needs_rotation"`
	Message       string  `json:"message,omitempty"`
}

// ValidateRailwayTokenResponse is the response for token validation
type ValidateRailwayTokenResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
}

// DeleteRailwayTokenResponse is the response for token deletion
type DeleteRailwayTokenResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// RotateRailwayTokenRequest is the request payload for rotating a Railway token
type RotateRailwayTokenRequest struct {
	NewToken string `json:"new_token" binding:"required"`
}

// RotateRailwayTokenResponse is the response for successful token rotation
type RotateRailwayTokenResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// StoreRailwayToken stores or updates a user's Railway API token in Vault.
// Before storing, it validates the token by attempting to create a Railway client
// and making a test API call.
//
// POST /api/v1/secrets/railway
// Request body: {"token": "railway-token-here"}
// Response: {"success": true, "validated": true, "stored_at": "2024-01-01T00:00:00Z"}
func (c *SecretsController) StoreRailwayToken(ctx *gin.Context) {
	if c.Vault == nil {
		log.Error().Msg("vault client not configured in secrets controller")
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "secret management not available",
		})
		return
	}

	// Get authenticated user
	user, err := auth.GetCurrentUser(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to get current user")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "authentication required",
		})
		return
	}

	// Parse request body
	var req StoreRailwayTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Str("user_id", user.ID).Msg("invalid request body for store railway token")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate token is not empty
	if req.Token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "token cannot be empty",
		})
		return
	}

	// Validate the Railway token by making a test API call
	// Create a temporary Railway client with the provided token
	testClient := railway.NewClient("", req.Token, nil)

	// Try to list projects as a lightweight validation
	// This will fail if the token is invalid
	_, err = testClient.ListProjects(ctx, 1)
	if err != nil {
		log.Warn().
			Err(err).
			Str("user_id", user.ID).
			Msg("railway token validation failed")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid railway token",
			"message": "The provided Railway token could not be validated. Please check that your token is correct and has the necessary permissions.",
			"details": err.Error(),
		})
		return
	}

	// Store the token in Vault
	err = c.Vault.StoreRailwayToken(ctx, user.ID, req.Token)
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", user.ID).
			Msg("failed to store railway token in vault")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to store railway token",
			"details": err.Error(),
		})
		return
	}

	log.Info().
		Str("user_id", user.ID).
		Msg("successfully stored railway token")

	ctx.JSON(http.StatusOK, StoreRailwayTokenResponse{
		Success:   true,
		Validated: true,
		StoredAt:  getCurrentTimestamp(),
		Message:   "Railway token stored successfully",
	})
}

// GetRailwayTokenStatus checks if the user has a Railway token configured
// and returns its status.
//
// GET /api/v1/secrets/railway/status
// Response: {"configured": true, "last_validated": "2024-01-01T00:00:00Z", "needs_rotation": false}
func (c *SecretsController) GetRailwayTokenStatus(ctx *gin.Context) {
	if c.Vault == nil {
		log.Error().Msg("vault client not configured in secrets controller")
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "secret management not available",
		})
		return
	}

	// Get authenticated user
	user, err := auth.GetCurrentUser(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to get current user")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "authentication required",
		})
		return
	}

	// Try to get the Railway token to check if it exists
	_, err = c.Vault.GetRailwayToken(ctx, user.ID)
	if err != nil {
		if err == vault.ErrSecretNotFound {
			log.Debug().
				Str("user_id", user.ID).
				Msg("no railway token configured")
			ctx.JSON(http.StatusOK, RailwayTokenStatusResponse{
				Configured:    false,
				NeedsRotation: false,
				Message:       "No Railway token configured. Please add your Railway token to get started.",
			})
			return
		}

		// Other error (Vault unavailable, etc.)
		log.Error().
			Err(err).
			Str("user_id", user.ID).
			Msg("failed to check railway token status")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to check token status",
			"details": err.Error(),
		})
		return
	}

	// Token exists - get metadata if available
	// For now, we just return that it's configured
	// TODO: In the future, we can fetch metadata (last_validated, etc.) from Vault
	ctx.JSON(http.StatusOK, RailwayTokenStatusResponse{
		Configured:    true,
		LastValidated: nil, // TODO: Fetch from metadata when available
		NeedsRotation: false,
		Message:       "Railway token is configured",
	})
}

// ValidateRailwayToken tests the user's Railway token by making a test API call
// to Railway. This can be used to verify the token is still valid.
//
// POST /api/v1/secrets/railway/validate
// Response: {"valid": true, "message": "Railway token is valid"}
func (c *SecretsController) ValidateRailwayToken(ctx *gin.Context) {
	if c.Vault == nil {
		log.Error().Msg("vault client not configured in secrets controller")
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "secret management not available",
		})
		return
	}

	// Get authenticated user
	user, err := auth.GetCurrentUser(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to get current user")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "authentication required",
		})
		return
	}

	// Get the Railway token from Vault
	token, err := c.Vault.GetRailwayToken(ctx, user.ID)
	if err != nil {
		if err == vault.ErrSecretNotFound {
			log.Debug().
				Str("user_id", user.ID).
				Msg("no railway token configured for validation")
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "no railway token configured",
				"message": "Please add your Railway token first",
			})
			return
		}

		log.Error().
			Err(err).
			Str("user_id", user.ID).
			Msg("failed to get railway token for validation")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to retrieve token",
			"details": err.Error(),
		})
		return
	}

	// Create a Railway client with the user's token and test it
	testClient := railway.NewClient("", token, nil)

	// Try to list projects as a lightweight validation
	_, err = testClient.ListProjects(ctx, 1)
	if err != nil {
		log.Warn().
			Err(err).
			Str("user_id", user.ID).
			Msg("railway token validation failed")
		ctx.JSON(http.StatusOK, ValidateRailwayTokenResponse{
			Valid:   false,
			Message: "Railway token is invalid or expired. Please update your token.",
		})
		return
	}

	// Update the last validated timestamp in Vault
	// This uses the ValidateRailwayToken method which updates metadata
	err = c.Vault.ValidateRailwayToken(ctx, user.ID)
	if err != nil {
		// Log but don't fail - the token is valid even if we can't update metadata
		log.Warn().
			Err(err).
			Str("user_id", user.ID).
			Msg("failed to update last validated timestamp")
	}

	log.Info().
		Str("user_id", user.ID).
		Msg("railway token validation successful")

	ctx.JSON(http.StatusOK, ValidateRailwayTokenResponse{
		Valid:   true,
		Message: "Railway token is valid",
	})
}

// DeleteRailwayToken removes the user's Railway token from Vault.
//
// DELETE /api/v1/secrets/railway
// Response: {"success": true, "message": "Railway token deleted successfully"}
func (c *SecretsController) DeleteRailwayToken(ctx *gin.Context) {
	if c.Vault == nil {
		log.Error().Msg("vault client not configured in secrets controller")
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "secret management not available",
		})
		return
	}

	// Get authenticated user
	user, err := auth.GetCurrentUser(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to get current user")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "authentication required",
		})
		return
	}

	// Delete the token from Vault
	err = c.Vault.DeleteRailwayToken(ctx, user.ID)
	if err != nil {
		if err == vault.ErrSecretNotFound {
			log.Debug().
				Str("user_id", user.ID).
				Msg("no railway token to delete")
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "no railway token configured",
				"message": "There is no Railway token to delete",
			})
			return
		}

		log.Error().
			Err(err).
			Str("user_id", user.ID).
			Msg("failed to delete railway token")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to delete railway token",
			"details": err.Error(),
		})
		return
	}

	log.Info().
		Str("user_id", user.ID).
		Msg("successfully deleted railway token")

	ctx.JSON(http.StatusOK, DeleteRailwayTokenResponse{
		Success: true,
		Message: "Railway token deleted successfully",
	})
}

// RotateRailwayToken replaces the user's Railway token with a new one.
// This validates the new token before storing it and creates a new version
// in Vault while preserving the old version for audit purposes.
//
// POST /api/v1/secrets/railway/rotate
// Request body: {"new_token": "new-railway-token"}
// Response: {"success": true, "message": "Railway token rotated successfully"}
func (c *SecretsController) RotateRailwayToken(ctx *gin.Context) {
	if c.Vault == nil {
		log.Error().Msg("vault client not configured in secrets controller")
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "secret management not available",
		})
		return
	}

	// Get authenticated user
	user, err := auth.GetCurrentUser(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to get current user")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "authentication required",
		})
		return
	}

	// Parse request body
	var req RotateRailwayTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Str("user_id", user.ID).Msg("invalid request body for rotate railway token")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate new token is not empty
	if req.NewToken == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "new token cannot be empty",
		})
		return
	}

	// Validate the new Railway token by making a test API call
	testClient := railway.NewClient("", req.NewToken, nil)

	// Try to list projects as a lightweight validation
	_, err = testClient.ListProjects(ctx, 1)
	if err != nil {
		log.Warn().
			Err(err).
			Str("user_id", user.ID).
			Msg("new railway token validation failed during rotation")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid railway token",
			"message": "The new Railway token could not be validated. Please check that your token is correct and has the necessary permissions.",
			"details": err.Error(),
		})
		return
	}

	// Rotate the token in Vault
	// This creates a new version while preserving the old one
	err = c.Vault.RotateRailwayToken(ctx, user.ID, req.NewToken)
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", user.ID).
			Msg("failed to rotate railway token")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to rotate railway token",
			"details": err.Error(),
		})
		return
	}

	log.Info().
		Str("user_id", user.ID).
		Msg("successfully rotated railway token")

	ctx.JSON(http.StatusOK, RotateRailwayTokenResponse{
		Success: true,
		Message: "Railway token rotated successfully",
	})
}

// StoreGitHubTokenRequest is the request payload for storing a GitHub token
type StoreGitHubTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// StoreGitHubTokenResponse is the response for successful GitHub token storage
type StoreGitHubTokenResponse struct {
	Success   bool     `json:"success"`
	Validated bool     `json:"validated"`
	StoredAt  string   `json:"stored_at"`
	Username  string   `json:"username"`
	Scopes    []string `json:"scopes"`
	Message   string   `json:"message"`
}

// GitHubTokenStatusResponse is the response for GitHub token status check
type GitHubTokenStatusResponse struct {
	Configured    bool     `json:"configured"`
	Username      *string  `json:"username,omitempty"`
	Scopes        []string `json:"scopes,omitempty"`
	LastValidated *string  `json:"last_validated,omitempty"`
	Message       string   `json:"message,omitempty"`
}

// ValidateGitHubTokenResponse is the response for GitHub token validation
type ValidateGitHubTokenResponse struct {
	Valid    bool     `json:"valid"`
	Username string   `json:"username,omitempty"`
	Scopes   []string `json:"scopes,omitempty"`
	Message  string   `json:"message"`
}

// DeleteGitHubTokenResponse is the response for GitHub token deletion
type DeleteGitHubTokenResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// StoreGitHubToken stores or updates a user's GitHub Personal Access Token in Vault.
// Before storing, it validates the token by calling the GitHub API to verify it's valid
// and extract the username and scopes.
//
// POST /api/v1/secrets/github
// Request body: {"token": "ghp_xxxxxxxxxxxx"}
// Response: {"success": true, "validated": true, "username": "octocat", "scopes": ["repo", "read:org"]}
func (c *SecretsController) StoreGitHubToken(ctx *gin.Context) {
	if c.Vault == nil {
		log.Error().Msg("vault client not configured in secrets controller")
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "secret management not available",
		})
		return
	}

	// Get authenticated user
	user, err := auth.GetCurrentUser(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to get current user")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "authentication required",
		})
		return
	}

	// Parse request body
	var req StoreGitHubTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Str("user_id", user.ID).Msg("invalid request body for store github token")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate token is not empty
	if req.Token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "token cannot be empty",
		})
		return
	}

	// Validate the GitHub token by calling the GitHub API
	// This will also extract the username and scopes
	username, scopes, err := c.Vault.ValidateGitHubToken(ctx, req.Token)
	if err != nil {
		log.Warn().
			Err(err).
			Str("user_id", user.ID).
			Msg("github token validation failed")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid github token",
			"message": "The provided GitHub token could not be validated. Please check that your token is correct and has the necessary permissions.",
			"details": err.Error(),
		})
		return
	}

	// Store the token in Vault
	err = c.Vault.StoreGitHubToken(ctx, user.ID, req.Token)
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", user.ID).
			Msg("failed to store github token in vault")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to store github token",
			"details": err.Error(),
		})
		return
	}

	// Update metadata with validation info
	err = c.Vault.ValidateGitHubTokenAndUpdateMetadata(ctx, user.ID)
	if err != nil {
		// Log but don't fail - token is stored successfully
		log.Warn().
			Err(err).
			Str("user_id", user.ID).
			Msg("failed to update github token metadata")
	}

	log.Info().
		Str("user_id", user.ID).
		Str("github_username", username).
		Msg("successfully stored github token")

	ctx.JSON(http.StatusOK, StoreGitHubTokenResponse{
		Success:   true,
		Validated: true,
		StoredAt:  time.Now().UTC().Format(time.RFC3339),
		Username:  username,
		Scopes:    scopes,
		Message:   "GitHub token stored successfully",
	})
}

// GetGitHubTokenStatus checks if the user has a GitHub token configured and returns
// its status including username and scopes if available.
//
// GET /api/v1/secrets/github/status
// Response: {"configured": true, "username": "octocat", "scopes": ["repo"]}
func (c *SecretsController) GetGitHubTokenStatus(ctx *gin.Context) {
	if c.Vault == nil {
		log.Error().Msg("vault client not configured in secrets controller")
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "secret management not available",
		})
		return
	}

	// Get authenticated user
	user, err := auth.GetCurrentUser(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to get current user")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "authentication required",
		})
		return
	}

	// Try to get the GitHub token to check if it exists
	_, err = c.Vault.GetGitHubToken(ctx, user.ID)
	if err != nil {
		if err == vault.ErrSecretNotFound {
			log.Debug().
				Str("user_id", user.ID).
				Msg("no github token configured")
			ctx.JSON(http.StatusOK, GitHubTokenStatusResponse{
				Configured: false,
				Message:    "No GitHub token configured. Please add your GitHub token to access private repositories.",
			})
			return
		}

		// Other error (Vault unavailable, etc.)
		log.Error().
			Err(err).
			Str("user_id", user.ID).
			Msg("failed to check github token status")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to check token status",
			"details": err.Error(),
		})
		return
	}

	// Token exists - return basic status
	// Note: In the future, we can fetch metadata (username, scopes, last_validated) from Vault
	ctx.JSON(http.StatusOK, GitHubTokenStatusResponse{
		Configured: true,
		Message:    "GitHub token is configured",
	})
}

// ValidateGitHubToken tests the user's GitHub token by making a test API call
// to GitHub. This verifies the token is still valid and returns username and scopes.
//
// POST /api/v1/secrets/github/validate
// Response: {"valid": true, "username": "octocat", "scopes": ["repo"], "message": "GitHub token is valid"}
func (c *SecretsController) ValidateGitHubToken(ctx *gin.Context) {
	if c.Vault == nil {
		log.Error().Msg("vault client not configured in secrets controller")
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "secret management not available",
		})
		return
	}

	// Get authenticated user
	user, err := auth.GetCurrentUser(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to get current user")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "authentication required",
		})
		return
	}

	// Get the GitHub token from Vault
	token, err := c.Vault.GetGitHubToken(ctx, user.ID)
	if err != nil {
		if err == vault.ErrSecretNotFound {
			log.Debug().
				Str("user_id", user.ID).
				Msg("no github token configured for validation")
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "no github token configured",
				"message": "Please add your GitHub token first",
			})
			return
		}

		log.Error().
			Err(err).
			Str("user_id", user.ID).
			Msg("failed to get github token for validation")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to retrieve token",
			"details": err.Error(),
		})
		return
	}

	// Validate the token with GitHub API
	username, scopes, err := c.Vault.ValidateGitHubToken(ctx, token)
	if err != nil {
		log.Warn().
			Err(err).
			Str("user_id", user.ID).
			Msg("github token validation failed")
		ctx.JSON(http.StatusOK, ValidateGitHubTokenResponse{
			Valid:   false,
			Message: "GitHub token is invalid or expired. Please update your token.",
		})
		return
	}

	// Update the last validated timestamp in Vault
	err = c.Vault.ValidateGitHubTokenAndUpdateMetadata(ctx, user.ID)
	if err != nil {
		// Log but don't fail - the token is valid even if we can't update metadata
		log.Warn().
			Err(err).
			Str("user_id", user.ID).
			Msg("failed to update last validated timestamp")
	}

	log.Info().
		Str("user_id", user.ID).
		Str("github_username", username).
		Msg("github token validation successful")

	ctx.JSON(http.StatusOK, ValidateGitHubTokenResponse{
		Valid:    true,
		Username: username,
		Scopes:   scopes,
		Message:  "GitHub token is valid",
	})
}

// DeleteGitHubToken removes the user's GitHub token from Vault.
//
// DELETE /api/v1/secrets/github
// Response: {"success": true, "message": "GitHub token deleted successfully"}
func (c *SecretsController) DeleteGitHubToken(ctx *gin.Context) {
	if c.Vault == nil {
		log.Error().Msg("vault client not configured in secrets controller")
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "secret management not available",
		})
		return
	}

	// Get authenticated user
	user, err := auth.GetCurrentUser(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to get current user")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "authentication required",
		})
		return
	}

	// Delete the token from Vault
	err = c.Vault.DeleteGitHubToken(ctx, user.ID)
	if err != nil {
		if err == vault.ErrSecretNotFound {
			log.Debug().
				Str("user_id", user.ID).
				Msg("no github token to delete")
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "no github token configured",
				"message": "There is no GitHub token to delete",
			})
			return
		}

		log.Error().
			Err(err).
			Str("user_id", user.ID).
			Msg("failed to delete github token")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to delete github token",
			"details": err.Error(),
		})
		return
	}

	log.Info().
		Str("user_id", user.ID).
		Msg("successfully deleted github token")

	ctx.JSON(http.StatusOK, DeleteGitHubTokenResponse{
		Success: true,
		Message: "GitHub token deleted successfully",
	})
}

// getCurrentTimestamp returns the current time in RFC3339 format
func getCurrentTimestamp() string {
	return time.Now().Format(time.RFC3339)
}
