package controller

import (
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
	Name         string                `json:"name" binding:"required"`
	Type         store.EnvironmentType `json:"type" binding:"required"`
	SourceRepo   string                `json:"sourceRepo"`
	SourceBranch string                `json:"sourceBranch"`
	SourceCommit string                `json:"sourceCommit"`
	TTLSeconds   *int64                `json:"ttlSeconds"`
}

type envResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt,omitempty"`
}

func (c *EnvironmentController) RegisterRoutes(r *gin.Engine) {
	r.GET("/environments", c.ListEnvironments)
	r.POST("/environments", c.CreateEnvironment)
	r.GET("/environments/:id", c.GetEnvironment)
	r.DELETE("/environments/:id", c.DestroyEnvironment)
	// railway proxy helpers
	r.GET("/railway/projects", c.ListRailwayProjects)
	r.GET("/railway/project/:id", c.GetRailwayProject)
	// provisioning endpoints
	r.POST("/api/provision/project", c.ProvisionProject)
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
		// Read project id from context (set by server)
		projectID, _ := ctx.Get("railway_project_id")
		pid, _ := projectID.(string)
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
