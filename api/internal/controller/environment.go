package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stwalsh4118/mirageapi/internal/railway"
	"github.com/stwalsh4118/mirageapi/internal/status"
	"github.com/stwalsh4118/mirageapi/internal/store"
	"gorm.io/gorm"
)

type EnvironmentController struct {
	DB      *gorm.DB
	Railway *railway.Client
}

type createEnvRequest struct {
	Name             string                `json:"name" binding:"required"`
	Type             store.EnvironmentType `json:"type" binding:"required"`
	SourceRepo       string                `json:"sourceRepo"`
	SourceBranch     string                `json:"sourceBranch"`
	SourceCommit     string                `json:"sourceCommit"`
	TTLSeconds       *int64                `json:"ttlSeconds"`
	RailwayProjectID *string               `json:"railwayProjectId"`
}

type envResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt,omitempty"`
}

func (c *EnvironmentController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/environments", c.ListEnvironments)
	r.POST("/environments", c.CreateEnvironment)
	r.GET("/environments/:id", c.GetEnvironment)
	r.DELETE("/environments/:id", c.DestroyEnvironment)
	// railway proxy helpers
	r.GET("/railway/projects", c.ListRailwayProjects)
	r.GET("/railway/project/:id", c.GetRailwayProject)
	r.DELETE("/railway/project/:id", c.DeleteRailwayProject)
	r.DELETE("/railway/environment/:id", c.DeleteRailwayEnvironment)
	// provisioning endpoints
	r.POST("/provision/project", c.ProvisionProject)
	r.POST("/provision/environment", c.ProvisionEnvironment)
}

func (c *EnvironmentController) ListEnvironments(ctx *gin.Context) {
	var envs []store.Environment
	if err := c.DB.Find(&envs).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list environments"})
		return
	}
	out := make([]envResponse, 0, len(envs))
	for _, e := range envs {
		out = append(out, envResponse{ID: e.ID, Name: e.Name, Type: string(e.Type), Status: status.NormalizeLocalToUI(e.Status), CreatedAt: e.CreatedAt.UTC().Format(time.RFC3339)})
	}
	ctx.JSON(http.StatusOK, out)
}

func (c *EnvironmentController) CreateEnvironment(ctx *gin.Context) {
	var req createEnvRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := uuid.NewString()
	env := store.Environment{
		ID:           id,
		Name:         req.Name,
		Type:         req.Type,
		SourceRepo:   req.SourceRepo,
		SourceBranch: req.SourceBranch,
		SourceCommit: req.SourceCommit,
		Status:       status.StatusCreating,
		TTLSeconds:   req.TTLSeconds,
		CreatedAt:    time.Now(),
	}
	if err := c.DB.Create(&env).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to persist environment"})
		return
	}

	if c.Railway != nil {
		// Determine Railway project id: prefer request override, otherwise use server-configured
		pid := ""
		if req.RailwayProjectID != nil && *req.RailwayProjectID != "" {
			pid = *req.RailwayProjectID
		} else {
			projectID, _ := ctx.Get("railway_project_id")
			pid, _ = projectID.(string)
		}
		go func(e store.Environment, projectID string) {
			res, err := c.Railway.CreateEnvironment(ctx, railway.CreateEnvironmentInput{ProjectID: projectID, Name: e.Name})
			if err != nil {
				log.Error().Err(err).Str("env_id", e.ID).Msg("railway create env failed")
				_ = c.DB.Model(&e).Update("Status", status.StatusError).Error
				return
			}
			updates := map[string]any{"Status": status.StatusActive, "RailwayEnvironmentID": res.EnvironmentID}
			_ = c.DB.Model(&e).Updates(updates).Error
		}(env, pid)
	}

	ctx.JSON(http.StatusAccepted, envResponse{ID: env.ID, Name: env.Name, Type: string(env.Type), Status: status.NormalizeLocalToUI(env.Status), CreatedAt: env.CreatedAt.UTC().Format(time.RFC3339)})
}

func (c *EnvironmentController) GetEnvironment(ctx *gin.Context) {
	id := ctx.Param("id")
	var env store.Environment
	if err := c.DB.First(&env, "id = ?", id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	ctx.JSON(http.StatusOK, envResponse{ID: env.ID, Name: env.Name, Type: string(env.Type), Status: status.NormalizeLocalToUI(env.Status), CreatedAt: env.CreatedAt.UTC().Format(time.RFC3339)})
}

func (c *EnvironmentController) DestroyEnvironment(ctx *gin.Context) {
	id := ctx.Param("id")
	var env store.Environment
	if err := c.DB.First(&env, "id = ?", id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if c.Railway != nil && env.RailwayEnvironmentID != "" {
		go func(e store.Environment) {
			if err := c.Railway.DestroyEnvironment(ctx, railway.DestroyEnvironmentInput{EnvironmentID: e.RailwayEnvironmentID}); err != nil {
				log.Error().Err(err).Str("env_id", e.ID).Msg("railway destroy env failed")
			}
		}(env)
	}
	if err := c.DB.Delete(&env).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete"})
		return
	}
	ctx.Status(http.StatusNoContent)
}

// status normalization moved to internal/status

// ProvisionEnvironmentRequest is the payload to create a Railway environment in an existing project.
type ProvisionEnvironmentRequest struct {
	ProjectID    string                 `json:"projectId"`
	Name         string                 `json:"name"`
	RequestID    string                 `json:"requestId"`
	EnvType      *store.EnvironmentType `json:"envType,omitempty"`      // Optional: dev, staging, prod, ephemeral
	WizardInputs map[string]interface{} `json:"wizardInputs,omitempty"` // Optional: full wizard state
}

type ProvisionEnvironmentResponse struct {
	EnvironmentID string `json:"environmentId"`
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
	if c.DB != nil {
		envType := store.EnvironmentTypeDev // Default to dev
		if req.EnvType != nil {
			envType = *req.EnvType
		}

		env := store.Environment{
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
		txErr := c.DB.Transaction(func(tx *gorm.DB) error {
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

	ctx.JSON(http.StatusOK, ProvisionEnvironmentResponse{EnvironmentID: res.EnvironmentID})
}
