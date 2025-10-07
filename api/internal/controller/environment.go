package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stwalsh4118/mirageapi/internal/railway"
	"github.com/stwalsh4118/mirageapi/internal/status"
	"github.com/stwalsh4118/mirageapi/internal/store"
	"gorm.io/gorm"
)

// RailwayEnvironmentClient defines the interface for Railway operations needed by the environment controller.
// This includes environment, project, and log operations used across environment.go, projects.go, and logs.go.
type RailwayEnvironmentClient interface {
	// Environment operations
	CreateEnvironment(ctx context.Context, in railway.CreateEnvironmentInput) (railway.CreateEnvironmentResult, error)
	DestroyEnvironment(ctx context.Context, in railway.DestroyEnvironmentInput) error
	GetEnvironmentVariables(ctx context.Context, in railway.GetEnvironmentVariablesInput) (railway.GetEnvironmentVariablesResult, error)
	GetAllEnvironmentAndServiceVariables(ctx context.Context, in railway.GetAllEnvironmentAndServiceVariablesInput) (railway.GetAllEnvironmentAndServiceVariablesResult, error)

	// Project operations
	CreateProject(ctx context.Context, in railway.CreateProjectInput) (railway.CreateProjectResult, error)
	DestroyProject(ctx context.Context, in railway.DestroyProjectInput) error
	GetProject(ctx context.Context, id string) (railway.Project, error)
	GetProjectWithDetailsByID(ctx context.Context, id string) (railway.ProjectDetails, error)
	ListProjects(ctx context.Context, limit int) ([]railway.Project, error)
	ListProjectsWithDetails(ctx context.Context, limit int) ([]railway.ProjectDetails, error)

	// Log operations (using websocket.Conn from github.com/coder/websocket)
	GetLatestDeploymentID(ctx context.Context, serviceID string) (string, error)
	GetDeploymentLogs(ctx context.Context, in railway.GetDeploymentLogsInput) (railway.GetDeploymentLogsResult, error)
	SubscribeToEnvironmentLogs(ctx context.Context, environmentID string, serviceFilter string) (*websocket.Conn, error)
	SubscribeToDeploymentLogs(ctx context.Context, deploymentID string, filter string) (*websocket.Conn, error)
}

type EnvironmentController struct {
	DB      *gorm.DB
	Railway RailwayEnvironmentClient
}

func (c *EnvironmentController) RegisterRoutes(r *gin.RouterGroup) {
	// railway proxy helpers
	r.GET("/railway/projects", c.ListRailwayProjects)
	r.GET("/railway/project/:id", c.GetRailwayProject)
	r.DELETE("/railway/project/:id", c.DeleteRailwayProject)
	r.DELETE("/railway/environment/:id", c.DeleteRailwayEnvironment)
	// provisioning endpoints
	r.POST("/provision/project", c.ProvisionProject)
	r.POST("/provision/environment", c.ProvisionEnvironment)
	// metadata retrieval endpoints
	r.GET("/environments/:id/metadata", c.GetEnvironmentMetadata)
	r.GET("/environments/:id/services", c.ListEnvironmentServices)
	r.GET("/environments/:id/snapshot", c.GetEnvironmentSnapshot)
	r.GET("/templates", c.ListTemplates)
}

// ProvisionEnvironmentRequest is the payload to create a Railway environment in an existing project.
type ProvisionEnvironmentRequest struct {
	ProjectID    string                 `json:"projectId"`
	Name         string                 `json:"name"`
	RequestID    string                 `json:"requestId"`
	EnvType      *store.EnvironmentType `json:"envType,omitempty"`      // Optional: dev, staging, prod, ephemeral
	WizardInputs map[string]interface{} `json:"wizardInputs,omitempty"` // Optional: full wizard state
}

type ProvisionEnvironmentResponse struct {
	EnvironmentID        string `json:"environmentId"`        // Mirage internal environment ID
	RailwayEnvironmentID string `json:"railwayEnvironmentId"` // Railway's environment ID
}

// ProvisionEnvironment creates a new environment under an existing Railway project.
// After successful creation, persists the environment and optional metadata to our database.
func (c *EnvironmentController) ProvisionEnvironment(ctx *gin.Context) {
	if c.Railway == nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "railway client not configured"})
		return
	}
	var req ProvisionEnvironmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := c.Railway.CreateEnvironment(ctx, railway.CreateEnvironmentInput{ProjectID: req.ProjectID, Name: req.Name})
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// Persist environment to database
	var env store.Environment
	var txErr error
	if c.DB != nil {
		envType := store.EnvironmentTypeDev // Default to dev
		if req.EnvType != nil {
			envType = *req.EnvType
		}

		env = store.Environment{
			ID:                   uuid.New().String(),
			Name:                 req.Name,
			Type:                 envType,
			Status:               status.StatusCreating,
			RailwayProjectID:     req.ProjectID,
			RailwayEnvironmentID: res.EnvironmentID,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}

		// Use transaction to ensure Environment and EnvironmentMetadata are created atomically
		txErr = c.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&env).Error; err != nil {
				return err
			}

			// Create EnvironmentMetadata if wizard inputs provided
			if len(req.WizardInputs) > 0 {
				wizardInputsJSON, err := json.Marshal(req.WizardInputs)
				if err != nil {
					log.Warn().Err(err).Msg("failed to marshal wizard inputs, skipping metadata")
					return nil // Don't fail transaction for metadata marshaling errors
				}

				// Create provision outputs JSON with Railway IDs
				provisionOutputs := map[string]interface{}{
					"projectId":     req.ProjectID,
					"environmentId": res.EnvironmentID,
				}
				provisionOutputsJSON, _ := json.Marshal(provisionOutputs)

				metadata := store.EnvironmentMetadata{
					ID:                   uuid.New().String(),
					EnvironmentID:        env.ID,
					WizardInputsJSON:     wizardInputsJSON,
					ProvisionOutputsJSON: provisionOutputsJSON,
					CreatedAt:            time.Now(),
					UpdatedAt:            time.Now(),
				}

				if err := tx.Create(&metadata).Error; err != nil {
					return err
				}

				log.Info().
					Str("env_id", env.ID).
					Str("metadata_id", metadata.ID).
					Msg("persisted environment metadata to database")
			}

			return nil
		})

		if txErr != nil {
			log.Error().Err(txErr).
				Str("project_id", req.ProjectID).
				Str("railway_env_id", res.EnvironmentID).
				Msg("failed to persist environment to database after Railway environment creation")
			// Don't fail the request - Railway resource was created successfully
		} else {
			log.Info().
				Str("env_id", env.ID).
				Str("railway_project_id", req.ProjectID).
				Str("railway_env_id", res.EnvironmentID).
				Msg("persisted environment to database")
		}
	}

	// If DB is available and transaction succeeded, return Mirage environment ID
	// Otherwise fall back to Railway ID
	mirageEnvID := res.EnvironmentID
	if c.DB != nil && txErr == nil {
		mirageEnvID = env.ID
	}

	ctx.JSON(http.StatusOK, ProvisionEnvironmentResponse{
		EnvironmentID:        mirageEnvID,       // Mirage ID for foreign keys (or Railway ID if DB unavailable/failed)
		RailwayEnvironmentID: res.EnvironmentID, // Railway ID for Railway API calls
	})
}

// EnvironmentMetadataDTO represents the complete metadata for an environment
type EnvironmentMetadataDTO struct {
	ID                  string                 `json:"id"`
	EnvironmentID       string                 `json:"environmentId"`
	IsTemplate          bool                   `json:"isTemplate"`
	TemplateName        *string                `json:"templateName,omitempty"`
	TemplateDescription *string                `json:"templateDescription,omitempty"`
	ClonedFromEnvID     *string                `json:"clonedFromEnvId,omitempty"`
	WizardInputs        map[string]interface{} `json:"wizardInputs,omitempty"`
	ProvisionOutputs    map[string]interface{} `json:"provisionOutputs,omitempty"`
	CreatedAt           string                 `json:"createdAt"`
	UpdatedAt           string                 `json:"updatedAt"`
}

// ServiceDetailDTO represents a service with its complete build configuration
type ServiceDetailDTO struct {
	ID               string  `json:"id"`
	EnvironmentID    string  `json:"environmentId"`
	Name             string  `json:"name"`
	Path             string  `json:"path,omitempty"`
	Status           string  `json:"status"`
	RailwayServiceID string  `json:"railwayServiceId,omitempty"`
	DeploymentType   string  `json:"deploymentType"`
	SourceRepo       string  `json:"sourceRepo,omitempty"`
	SourceBranch     string  `json:"sourceBranch,omitempty"`
	DockerfilePath   *string `json:"dockerfilePath,omitempty"`
	BuildContext     *string `json:"buildContext,omitempty"`
	RootDirectory    *string `json:"rootDirectory,omitempty"`
	TargetStage      *string `json:"targetStage,omitempty"`
	DockerImage      string  `json:"dockerImage,omitempty"`
	ImageRegistry    string  `json:"imageRegistry,omitempty"`
	ImageName        string  `json:"imageName,omitempty"`
	ImageTag         string  `json:"imageTag,omitempty"`
	ImageDigest      string  `json:"imageDigest,omitempty"`
	ExposedPorts     []int   `json:"exposedPorts,omitempty"`
	HealthCheckPath  *string `json:"healthCheckPath,omitempty"`
	StartCommand     *string `json:"startCommand,omitempty"`
	CreatedAt        string  `json:"createdAt"`
	UpdatedAt        string  `json:"updatedAt"`
}

// TemplateListItemDTO represents a template in the list view
type TemplateListItemDTO struct {
	ID                  string  `json:"id"`
	EnvironmentID       string  `json:"environmentId"`
	TemplateName        string  `json:"templateName"`
	TemplateDescription *string `json:"templateDescription,omitempty"`
	EnvironmentName     string  `json:"environmentName"`
	EnvironmentType     string  `json:"environmentType"`
	ServiceCount        int     `json:"serviceCount"`
	CreatedAt           string  `json:"createdAt"`
}

// GetEnvironmentMetadata retrieves the metadata for a specific environment
// The :id parameter is the Railway environment ID
func (c *EnvironmentController) GetEnvironmentMetadata(ctx *gin.Context) {
	railwayEnvID := ctx.Param("id")
	if railwayEnvID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "railway environment id required"})
		return
	}

	// Look up the Mirage environment by Railway ID
	var env store.Environment
	if err := c.DB.Where("railway_environment_id = ?", railwayEnvID).First(&env).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "environment not found"})
			return
		}
		log.Error().Err(err).Str("railway_env_id", railwayEnvID).Msg("failed to query environment")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve environment"})
		return
	}

	// Look up metadata using the Mirage environment ID
	var metadata store.EnvironmentMetadata
	log.Info().Str("railway_env_id", railwayEnvID).Str("mirage_env_id", env.ID).Msg("getting environment metadata")
	if err := c.DB.Where("environment_id = ?", env.ID).First(&metadata).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "environment metadata not found"})
			return
		}
		log.Error().Err(err).Str("mirage_env_id", env.ID).Msg("failed to query environment metadata")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve metadata"})
		return
	}

	// Unmarshal JSON fields
	var wizardInputs map[string]interface{}
	var provisionOutputs map[string]interface{}

	if len(metadata.WizardInputsJSON) > 0 {
		if err := json.Unmarshal(metadata.WizardInputsJSON, &wizardInputs); err != nil {
			log.Warn().Err(err).Msg("failed to unmarshal wizard inputs")
		}
	}

	if len(metadata.ProvisionOutputsJSON) > 0 {
		if err := json.Unmarshal(metadata.ProvisionOutputsJSON, &provisionOutputs); err != nil {
			log.Warn().Err(err).Msg("failed to unmarshal provision outputs")
		}
	}

	dto := EnvironmentMetadataDTO{
		ID:                  metadata.ID,
		EnvironmentID:       metadata.EnvironmentID,
		IsTemplate:          metadata.IsTemplate,
		TemplateName:        metadata.TemplateName,
		TemplateDescription: metadata.TemplateDescription,
		ClonedFromEnvID:     metadata.ClonedFromEnvID,
		WizardInputs:        wizardInputs,
		ProvisionOutputs:    provisionOutputs,
		CreatedAt:           metadata.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:           metadata.UpdatedAt.UTC().Format(time.RFC3339),
	}

	ctx.JSON(http.StatusOK, dto)
}

// ListEnvironmentServices retrieves all services for a specific environment
// The :id parameter is the Railway environment ID
func (c *EnvironmentController) ListEnvironmentServices(ctx *gin.Context) {
	railwayEnvID := ctx.Param("id")
	if railwayEnvID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "railway environment id required"})
		return
	}

	// Look up the Mirage environment by Railway ID
	var env store.Environment
	if err := c.DB.Where("railway_environment_id = ?", railwayEnvID).First(&env).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "environment not found"})
			return
		}
		log.Error().Err(err).Str("railway_env_id", railwayEnvID).Msg("failed to query environment")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve environment"})
		return
	}

	// Query services using the Mirage environment ID
	var services []store.Service
	log.Info().Str("railway_env_id", railwayEnvID).Str("mirage_env_id", env.ID).Msg("listing environment services")
	if err := c.DB.Where("environment_id = ?", env.ID).Find(&services).Error; err != nil {
		log.Error().Err(err).Str("mirage_env_id", env.ID).Msg("failed to query services")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve services"})
		return
	}

	dtos := make([]ServiceDetailDTO, 0, len(services))
	for _, svc := range services {
		dto := serviceToDTO(svc)
		dtos = append(dtos, dto)
	}

	ctx.JSON(http.StatusOK, dtos)
}

// ListTemplates retrieves all environments marked as templates
func (c *EnvironmentController) ListTemplates(ctx *gin.Context) {
	var metadataList []store.EnvironmentMetadata
	if err := c.DB.Where("is_template = ?", true).Find(&metadataList).Error; err != nil {
		log.Error().Err(err).Msg("failed to query templates")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve templates"})
		return
	}

	dtos := make([]TemplateListItemDTO, 0, len(metadataList))
	for _, meta := range metadataList {
		// Fetch associated environment for additional info
		var env store.Environment
		if err := c.DB.Preload("Services").First(&env, "id = ?", meta.EnvironmentID).Error; err != nil {
			log.Warn().Err(err).Str("env_id", meta.EnvironmentID).Msg("failed to load environment for template")
			continue
		}

		templateName := "Unnamed Template"
		if meta.TemplateName != nil {
			templateName = *meta.TemplateName
		}

		dto := TemplateListItemDTO{
			ID:                  meta.ID,
			EnvironmentID:       meta.EnvironmentID,
			TemplateName:        templateName,
			TemplateDescription: meta.TemplateDescription,
			EnvironmentName:     env.Name,
			EnvironmentType:     string(env.Type),
			ServiceCount:        len(env.Services),
			CreatedAt:           meta.CreatedAt.UTC().Format(time.RFC3339),
		}
		dtos = append(dtos, dto)
	}

	ctx.JSON(http.StatusOK, dtos)
}

// ServiceVariablesSnapshot represents variables for a specific service
type ServiceVariablesSnapshot struct {
	ServiceID   string            `json:"serviceId"`
	ServiceName string            `json:"serviceName"`
	Variables   map[string]string `json:"variables"`
}

// EnvironmentSnapshot represents the complete data needed for cloning an environment
type EnvironmentSnapshot struct {
	Environment          store.Environment          `json:"environment"`
	Services             []store.Service            `json:"services"`
	EnvironmentVariables map[string]string          `json:"environmentVariables"`
	ServiceVariables     []ServiceVariablesSnapshot `json:"serviceVariables"`
}

// GetEnvironmentSnapshot returns all data needed to clone an environment.
// This includes the environment details, all services, and environment variables from Railway.
// The :id parameter is the Railway environment ID.
func (c *EnvironmentController) GetEnvironmentSnapshot(ctx *gin.Context) {
	railwayEnvID := ctx.Param("id")
	if railwayEnvID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "railway environment id required"})
		return
	}

	// Look up the Mirage environment by Railway ID
	var env store.Environment
	if err := c.DB.Preload("Services").Where("railway_environment_id = ?", railwayEnvID).First(&env).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "environment not found"})
			return
		}
		log.Error().Err(err).Str("railway_env_id", railwayEnvID).Msg("failed to query environment")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve environment"})
		return
	}

	log.Info().
		Str("railway_env_id", railwayEnvID).
		Str("mirage_env_id", env.ID).
		Int("service_count", len(env.Services)).
		Msg("fetching environment snapshot")

	// Fetch environment and service variables from Railway API
	envVars := make(map[string]string)
	serviceVars := make([]ServiceVariablesSnapshot, 0)

	if c.Railway != nil && env.RailwayEnvironmentID != "" && env.RailwayProjectID != "" {
		result, err := c.Railway.GetAllEnvironmentAndServiceVariables(ctx, railway.GetAllEnvironmentAndServiceVariablesInput{
			ProjectID:     env.RailwayProjectID,
			EnvironmentID: env.RailwayEnvironmentID,
		})
		if err != nil {
			log.Warn().
				Err(err).
				Str("railway_env_id", env.RailwayEnvironmentID).
				Str("railway_project_id", env.RailwayProjectID).
				Msg("failed to fetch variables for snapshot, returning empty variables")
			// Don't fail - just return empty vars
		} else {
			envVars = result.EnvironmentVariables

			// Convert Railway service variables to snapshot format
			for _, svcVar := range result.ServiceVariables {
				serviceVars = append(serviceVars, ServiceVariablesSnapshot{
					ServiceID:   svcVar.ServiceID,
					ServiceName: svcVar.ServiceName,
					Variables:   svcVar.Variables,
				})
			}
		}
	}

	snapshot := EnvironmentSnapshot{
		Environment:          env,
		Services:             env.Services,
		EnvironmentVariables: envVars,
		ServiceVariables:     serviceVars,
	}

	log.Info().
		Str("railway_env_id", railwayEnvID).
		Str("mirage_env_id", env.ID).
		Int("service_count", len(env.Services)).
		Int("env_variable_count", len(envVars)).
		Int("services_with_variables", len(serviceVars)).
		Msg("successfully created environment snapshot")

	ctx.JSON(http.StatusOK, snapshot)
}

// serviceToDTO converts a store.Service model to a ServiceDetailDTO
func serviceToDTO(svc store.Service) ServiceDetailDTO {
	dto := ServiceDetailDTO{
		ID:               svc.ID,
		EnvironmentID:    svc.EnvironmentID,
		Name:             svc.Name,
		Path:             svc.Path,
		Status:           status.NormalizeLocalToUI(svc.Status),
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
