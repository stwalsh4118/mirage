package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/stwalsh4118/mirageapi/internal/railway"
)

// RailwayServiceClient defines the interface for Railway service operations needed by the controller.
type RailwayServiceClient interface {
	CreateService(ctx context.Context, in railway.CreateServiceInput) (railway.CreateServiceResult, error)
}

// ServicesController handles Railway service provisioning endpoints.
type ServicesController struct {
	Railway RailwayServiceClient
}

// RegisterRoutes registers service-related routes under the provided router group.
func (c *ServicesController) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/provision/services", c.ProvisionServices)
}

// ServiceSpec represents a single service to provision, supporting both
// repository-based and Docker image-based deployments.
type ServiceSpec struct {
	Name string `json:"name"`

	// Repository-based deployment fields
	Repo   *string `json:"repo"`
	Branch *string `json:"branch"`

	// Docker image-based deployment fields
	ImageRegistry *string           `json:"imageRegistry"` // Optional, defaults to Docker Hub
	ImageName     *string           `json:"imageName"`     // Required for image deployment
	ImageTag      *string           `json:"imageTag"`      // Optional, defaults to "latest"
	EnvVars       map[string]string `json:"environmentVariables,omitempty"`

	// Registry authentication (optional, for private images)
	RegistryUsername *string `json:"registryUsername"`
	RegistryPassword *string `json:"registryPassword"`

	// Dockerfile path for monorepo builds (optional, relative to repo root)
	DockerfilePath *string `json:"dockerfilePath,omitempty"`
}

// ProvisionServicesRequest creates one or more services in a given environment.
type ProvisionServicesRequest struct {
	ProjectID     string        `json:"projectId"`
	EnvironmentID string        `json:"environmentId"`
	Services      []ServiceSpec `json:"services"`
	RequestID     string        `json:"requestId"`
}

type ProvisionServicesResponse struct {
	ServiceIDs []string `json:"serviceIds"`
}

// ProvisionServices creates services sequentially and returns their IDs.
// Supports both repository-based and Docker image-based deployments.
func (c *ServicesController) ProvisionServices(ctx *gin.Context) {
	if c.Railway == nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "railway client not configured"})
		return
	}
	var req ProvisionServicesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate each service has either repo OR image fields
	for i, s := range req.Services {
		if err := validateServiceSpec(s); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"service": s.Name,
				"index":   i,
			})
			return
		}
	}

	ids := make([]string, 0, len(req.Services))
	for _, s := range req.Services {
		input := railway.CreateServiceInput{
			ProjectID:     req.ProjectID,
			EnvironmentID: req.EnvironmentID,
			Name:          s.Name,
		}

		// Configure based on deployment type
		if s.ImageName != nil {
			// Docker image deployment
			imageRef := buildImageReference(s)
			input.Image = &imageRef

			// Add registry credentials if provided
			if s.RegistryUsername != nil && s.RegistryPassword != nil {
				input.RegistryCredentials = &railway.RegistryCredentials{
					Username: *s.RegistryUsername,
					Password: *s.RegistryPassword,
				}
			}
		} else {
			// Repository-based deployment
			input.Repo = s.Repo
			input.Branch = s.Branch
		}

		// Merge user-specified environment variables
		// These are added first so system variables can override them if needed
		if s.EnvVars != nil && len(s.EnvVars) > 0 {
			if input.Variables == nil {
				input.Variables = make(map[string]string)
			}
			for k, v := range s.EnvVars {
				input.Variables[k] = v
			}
			log.Debug().
				Str("service", s.Name).
				Int("env_var_count", len(s.EnvVars)).
				Msg("adding user-specified environment variables")
		}

		// Set system variables (these override user variables if there's a conflict)
		// Set Dockerfile path if specified (repository deployments only)
		if s.DockerfilePath != nil && *s.DockerfilePath != "" {
			if input.Variables == nil {
				input.Variables = make(map[string]string)
			}
			input.Variables["RAILWAY_DOCKERFILE_PATH"] = *s.DockerfilePath
			log.Info().
				Str("service", s.Name).
				Str("dockerfile_path", *s.DockerfilePath).
				Msg("setting RAILWAY_DOCKERFILE_PATH system variable for service")
		}

		// Log final variable set if any variables are present
		if len(input.Variables) > 0 {
			log.Debug().
				Str("service", s.Name).
				Interface("variables", input.Variables).
				Msg("creating service with merged variables")
		}

		out, err := c.Railway.CreateService(ctx, input)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error(), "service": s.Name, "partial": ids})
			return
		}
		ids = append(ids, out.ServiceID)
	}
	ctx.JSON(http.StatusOK, ProvisionServicesResponse{ServiceIDs: ids})
}

// validateServiceSpec ensures a service has either repo OR image configuration, not both or neither.
func validateServiceSpec(s ServiceSpec) error {
	hasRepo := s.Repo != nil && *s.Repo != ""
	hasImage := s.ImageName != nil && *s.ImageName != ""

	if hasRepo && hasImage {
		return fmt.Errorf("service '%s': cannot specify both repository and image deployment options", s.Name)
	}

	if !hasRepo && !hasImage {
		return fmt.Errorf("service '%s': must specify either repository (repo+branch) or image (imageName) deployment", s.Name)
	}

	// Additional validation for repo deployment
	if hasRepo && (s.Branch == nil || *s.Branch == "") {
		return fmt.Errorf("service '%s': branch is required when repo is specified", s.Name)
	}

	// Additional validation for registry credentials
	if (s.RegistryUsername != nil && s.RegistryPassword == nil) ||
		(s.RegistryUsername == nil && s.RegistryPassword != nil) {
		return fmt.Errorf("service '%s': both registryUsername and registryPassword must be provided together", s.Name)
	}

	// Validate Dockerfile path if provided (only valid for repo deployments)
	if s.DockerfilePath != nil && *s.DockerfilePath != "" {
		if !hasRepo {
			return fmt.Errorf("service '%s': dockerfilePath can only be specified for repository-based deployments", s.Name)
		}
		if err := validateDockerfilePath(*s.DockerfilePath); err != nil {
			return fmt.Errorf("service '%s': %w", s.Name, err)
		}
	}

	return nil
}

// buildImageReference constructs a full Docker image reference from components.
// Format: [registry/]name[:tag]
func buildImageReference(s ServiceSpec) string {
	const defaultTag = "latest"

	// Start with image name (required)
	imageRef := *s.ImageName

	// Add registry prefix if provided and not already in image name
	if s.ImageRegistry != nil && *s.ImageRegistry != "" {
		// Only prepend registry if the image name doesn't already contain a registry
		// (i.e., doesn't contain a slash or starts with known registries)
		imageRef = fmt.Sprintf("%s/%s", *s.ImageRegistry, imageRef)
	}

	// Add tag suffix if provided, otherwise use default
	if s.ImageTag != nil && *s.ImageTag != "" {
		imageRef = fmt.Sprintf("%s:%s", imageRef, *s.ImageTag)
	} else {
		imageRef = fmt.Sprintf("%s:%s", imageRef, defaultTag)
	}

	return imageRef
}

// validateDockerfilePath validates that a Dockerfile path is safe and relative.
// Rejects absolute paths and parent directory traversal attempts.
func validateDockerfilePath(path string) error {
	if path == "" {
		return fmt.Errorf("dockerfile path cannot be empty")
	}

	// Reject absolute paths (Unix-style)
	if len(path) > 0 && path[0] == '/' {
		return fmt.Errorf("dockerfile path must be relative to repository root, not absolute: %s", path)
	}

	// Reject absolute paths (Windows-style)
	if len(path) > 1 && path[1] == ':' {
		return fmt.Errorf("dockerfile path must be relative to repository root, not absolute: %s", path)
	}

	// Reject parent directory traversal
	if len(path) >= 3 && (path[:3] == "../" || path[:3] == "..\\") {
		return fmt.Errorf("dockerfile path cannot traverse parent directories: %s", path)
	}

	// Check for .. anywhere in path (more thorough check)
	for i := 0; i < len(path)-1; i++ {
		if path[i] == '.' && path[i+1] == '.' {
			// Allow consecutive dots only if they're part of a filename (not directory traversal)
			if i > 0 && (path[i-1] == '/' || path[i-1] == '\\') {
				return fmt.Errorf("dockerfile path cannot traverse parent directories: %s", path)
			}
			if i+2 < len(path) && (path[i+2] == '/' || path[i+2] == '\\') {
				return fmt.Errorf("dockerfile path cannot traverse parent directories: %s", path)
			}
		}
	}

	// Ensure reasonable length
	const maxPathLength = 512
	if len(path) > maxPathLength {
		return fmt.Errorf("dockerfile path exceeds maximum length of %d characters", maxPathLength)
	}

	return nil
}
