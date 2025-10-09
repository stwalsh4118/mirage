package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stwalsh4118/mirageapi/internal/auth"
	"github.com/stwalsh4118/mirageapi/internal/railway"
	"github.com/stwalsh4118/mirageapi/internal/store"
	"gorm.io/gorm"
)

// RailwayServiceClient defines the interface for Railway service operations needed by the controller.
type RailwayServiceClient interface {
	CreateService(ctx context.Context, in railway.CreateServiceInput) (railway.CreateServiceResult, error)
	DestroyService(ctx context.Context, in railway.DestroyServiceInput) error
}

// ServicesController handles Railway service provisioning endpoints.
type ServicesController struct {
	Railway RailwayServiceClient
	DB      *gorm.DB
}

// RegisterRoutes registers service-related routes under the provided router group.
func (c *ServicesController) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/provision/services", c.ProvisionServices)
	r.GET("/services/:id", c.GetService)
	r.DELETE("/railway/service/:id", c.DeleteRailwayService)
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
	ProjectID            string        `json:"projectId"`
	EnvironmentID        string        `json:"environmentId"`        // Mirage internal ID for database FK
	RailwayEnvironmentID string        `json:"railwayEnvironmentId"` // Railway ID for Railway API calls
	Services             []ServiceSpec `json:"services"`
	RequestID            string        `json:"requestId"`
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

	// Get authenticated user
	user, err := auth.GetCurrentUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	var req ProvisionServicesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify user owns the parent environment before creating services
	if c.DB != nil {
		var env store.Environment
		err = c.DB.Where("id = ? AND user_id = ?", req.EnvironmentID, user.ID).First(&env).Error
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "environment not found or access denied"})
			return
		} else if err != nil {
			log.Error().Err(err).Msg("failed to verify environment ownership")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify environment ownership"})
			return
		}
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

	// Determine which environment ID to use for Railway API
	// If RailwayEnvironmentID is provided, use it; otherwise fall back to EnvironmentID for backward compatibility
	railwayEnvID := req.RailwayEnvironmentID
	if railwayEnvID == "" {
		railwayEnvID = req.EnvironmentID
		log.Warn().
			Str("environment_id", req.EnvironmentID).
			Msg("RailwayEnvironmentID not provided, using EnvironmentID for Railway API (may fail if it's a Mirage ID)")
	}

	ids := make([]string, 0, len(req.Services))
	for _, s := range req.Services {
		input := railway.CreateServiceInput{
			ProjectID:     req.ProjectID,
			EnvironmentID: railwayEnvID, // Use Railway environment ID for Railway API
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
		if len(s.EnvVars) > 0 {
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

		// Persist service to database with UserID
		if c.DB != nil {
			serviceModel, err := serviceSpecToModel(s, req.EnvironmentID, out.ServiceID)
			if err != nil {
				log.Error().Err(err).
					Str("service_name", s.Name).
					Str("railway_service_id", out.ServiceID).
					Msg("failed to convert service spec to model")
				// Continue - Railway service was created successfully
				continue
			}

			// Set UserID from authenticated user
			serviceModel.UserID = user.ID

			if err := c.DB.Create(&serviceModel).Error; err != nil {
				log.Error().Err(err).
					Str("service_name", s.Name).
					Str("railway_service_id", out.ServiceID).
					Msg("failed to persist service to database after Railway service creation")
				// Don't fail the request - Railway service was created successfully
			} else {
				log.Info().
					Str("service_id", serviceModel.ID).
					Str("service_name", serviceModel.Name).
					Str("railway_service_id", out.ServiceID).
					Str("user_id", user.ID).
					Str("deployment_type", string(serviceModel.DeploymentType)).
					Msg("persisted service to database with ownership")
			}
		}
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

// serviceSpecToModel converts a ServiceSpec and Railway service ID to a store.Service model.
func serviceSpecToModel(spec ServiceSpec, environmentID string, railwayServiceID string) (store.Service, error) {
	service := store.Service{
		ID:               uuid.New().String(),
		EnvironmentID:    environmentID,
		Name:             spec.Name,
		Status:           "provisioning",
		RailwayServiceID: railwayServiceID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Determine deployment type and set fields accordingly
	if spec.ImageName != nil {
		// Docker image deployment
		service.DeploymentType = store.DeploymentTypeDockerImage
		service.DockerImage = buildImageReference(spec)

		if spec.ImageRegistry != nil {
			service.ImageRegistry = *spec.ImageRegistry
		}
		if spec.ImageName != nil {
			service.ImageName = *spec.ImageName
		}
		if spec.ImageTag != nil {
			service.ImageTag = *spec.ImageTag
		}
		// ImageDigest and ImageAuthStored would be set later if needed
	} else {
		// Repository-based deployment
		service.DeploymentType = store.DeploymentTypeSourceRepo

		if spec.Repo != nil {
			service.SourceRepo = *spec.Repo
		}
		if spec.Branch != nil {
			service.SourceBranch = *spec.Branch
		}
		if spec.DockerfilePath != nil {
			service.DockerfilePath = spec.DockerfilePath
		}
	}

	// ExposedPortsJSON - initialize as empty array for now
	// In the future, this could be extracted from service configuration
	service.ExposedPortsJSON = "[]"

	return service, nil
}

// GetService retrieves a single service by ID with its complete build configuration
func (c *ServicesController) GetService(ctx *gin.Context) {
	serviceID := ctx.Param("id")
	if serviceID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "service id required"})
		return
	}

	// Get authenticated user
	user, err := auth.GetCurrentUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	// Query service with ownership check
	var service store.Service
	err = c.DB.Where("id = ? AND user_id = ?", serviceID, user.ID).First(&service).Error
	if err == gorm.ErrRecordNotFound {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
		return
	} else if err != nil {
		log.Error().Err(err).Str("service_id", serviceID).Msg("failed to query service")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve service"})
		return
	}

	// Convert to ServiceDetailDTO
	dto := serviceModelToServiceDetailDTO(service)
	ctx.JSON(http.StatusOK, dto)
}

// serviceModelToServiceDetailDTO converts a store.Service model to a ServiceDetailDTO
// Note: ServiceDetailDTO is defined in environment.go to avoid duplication
func serviceModelToServiceDetailDTO(svc store.Service) ServiceDetailDTO {
	dto := ServiceDetailDTO{
		ID:               svc.ID,
		EnvironmentID:    svc.EnvironmentID,
		Name:             svc.Name,
		Path:             svc.Path,
		Status:           svc.Status,
		RailwayServiceID: svc.RailwayServiceID,
		DeploymentType:   string(svc.DeploymentType),
		SourceRepo:       svc.SourceRepo,
		SourceBranch:     svc.SourceBranch,
		DockerfilePath:   svc.DockerfilePath,
		BuildContext:     svc.BuildContext,
		RootDirectory:    svc.RootDirectory,
		TargetStage:      svc.TargetStage,
		DockerImage:      svc.DockerImage,
		ImageRegistry:    svc.ImageRegistry,
		ImageName:        svc.ImageName,
		ImageTag:         svc.ImageTag,
		ImageDigest:      svc.ImageDigest,
		HealthCheckPath:  svc.HealthCheckPath,
		StartCommand:     svc.StartCommand,
		CreatedAt:        svc.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:        svc.UpdatedAt.UTC().Format(time.RFC3339),
	}

	// Parse ExposedPortsJSON
	if svc.ExposedPortsJSON != "" {
		var ports []int
		if err := json.Unmarshal([]byte(svc.ExposedPortsJSON), &ports); err == nil {
			dto.ExposedPorts = ports
		}
	}

	return dto
}

// DeleteRailwayService deletes a Railway service by its Railway service ID.
// After successful Railway deletion, cleans up the service record from the database.
func (c *ServicesController) DeleteRailwayService(ctx *gin.Context) {
	if c.Railway == nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "railway client not configured"})
		return
	}
	railwayServiceID := ctx.Param("id")
	if railwayServiceID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "service id required"})
		return
	}

	// Get authenticated user
	user, err := auth.GetCurrentUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	// Verify ownership before deleting
	var service store.Service
	err = c.DB.Where("railway_service_id = ? AND user_id = ?", railwayServiceID, user.ID).First(&service).Error
	if err == gorm.ErrRecordNotFound {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "service not found or access denied"})
		return
	} else if err != nil {
		log.Error().Err(err).Str("railway_service_id", railwayServiceID).Msg("failed to verify service ownership")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify ownership"})
		return
	}

	// Step 1: Delete from Railway first (fail fast if Railway API fails)
	log.Info().
		Str("railway_service_id", railwayServiceID).
		Str("user_id", user.ID).
		Msg("deleting railway service")
	if err := c.Railway.DestroyService(ctx, railway.DestroyServiceInput{ServiceID: railwayServiceID}); err != nil {
		log.Error().Err(err).Str("railway_service_id", railwayServiceID).Msg("railway delete service failed")
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// Step 2: Clean up database (Railway deletion succeeded, so clean up our record)
	// Note: service was already fetched above with ownership verification
	if c.DB != nil {
		// Delete the service from the database
		result := c.DB.Where("railway_service_id = ?", railwayServiceID).Delete(&store.Service{})
		if result.Error != nil {
			log.Error().Err(result.Error).
				Str("railway_service_id", railwayServiceID).
				Msg("failed to delete service from database")
			// Don't fail the request - Railway resource was deleted successfully
		} else {
			log.Info().
				Str("railway_service_id", railwayServiceID).
				Int64("rows_deleted", result.RowsAffected).
				Msg("deleted service from database")
		}
	}

	ctx.Status(http.StatusNoContent)
}
