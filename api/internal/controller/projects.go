package controller

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stwalsh4118/mirageapi/internal/auth"
	"github.com/stwalsh4118/mirageapi/internal/railway"
	"github.com/stwalsh4118/mirageapi/internal/status"
	"github.com/stwalsh4118/mirageapi/internal/store"
	"gorm.io/gorm"
)

type ProjectDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ServiceSourceDTO struct {
	Image *string `json:"image,omitempty"`
	Repo  *string `json:"repo,omitempty"`
}

type LatestDeploymentDTO struct {
	CanRedeploy             *bool          `json:"canRedeploy,omitempty"`
	CanRollback             *bool          `json:"canRollback,omitempty"`
	CreatedAt               *string        `json:"createdAt,omitempty"`
	DeploymentStopped       *bool          `json:"deploymentStopped,omitempty"`
	EnvironmentID           *string        `json:"environmentId,omitempty"`
	ID                      *string        `json:"id,omitempty"`
	Meta                    map[string]any `json:"meta,omitempty"`
	ProjectID               *string        `json:"projectId,omitempty"`
	ServiceID               *string        `json:"serviceId,omitempty"`
	SnapshotID              *string        `json:"snapshotId,omitempty"`
	StaticURL               *string        `json:"staticUrl,omitempty"`
	Status                  *string        `json:"status,omitempty"`
	StatusUpdatedAt         *string        `json:"statusUpdatedAt,omitempty"`
	SuggestAddServiceDomain *bool          `json:"suggestAddServiceDomain,omitempty"`
	UpdatedAt               *string        `json:"updatedAt,omitempty"`
	URL                     *string        `json:"url,omitempty"`
}

type ServiceInstanceDTO struct {
	ID                      string               `json:"id"`
	ServiceID               string               `json:"serviceId"`
	ServiceName             string               `json:"serviceName"`
	EnvironmentID           string               `json:"environmentId"`
	BuildCommand            *string              `json:"buildCommand,omitempty"`
	Builder                 *string              `json:"builder,omitempty"`
	CreatedAt               *string              `json:"createdAt,omitempty"`
	CronSchedule            *string              `json:"cronSchedule,omitempty"`
	DeletedAt               *string              `json:"deletedAt,omitempty"`
	DrainingSeconds         *int                 `json:"drainingSeconds,omitempty"`
	HealthcheckPath         *string              `json:"healthcheckPath,omitempty"`
	HealthcheckTimeout      *int                 `json:"healthcheckTimeout,omitempty"`
	IsUpdatable             *bool                `json:"isUpdatable,omitempty"`
	NextCronRunAt           *string              `json:"nextCronRunAt,omitempty"`
	NixpacksPlan            *string              `json:"nixpacksPlan,omitempty"`
	NumReplicas             *int                 `json:"numReplicas,omitempty"`
	OverlapSeconds          *int                 `json:"overlapSeconds,omitempty"`
	PreDeployCommand        *string              `json:"preDeployCommand,omitempty"`
	RailpackInfo            *string              `json:"railpackInfo,omitempty"`
	RailwayConfigFile       *string              `json:"railwayConfigFile,omitempty"`
	Region                  *string              `json:"region,omitempty"`
	RestartPolicyMaxRetries *int                 `json:"restartPolicyMaxRetries,omitempty"`
	RestartPolicyType       *string              `json:"restartPolicyType,omitempty"`
	RootDirectory           *string              `json:"rootDirectory,omitempty"`
	SleepApplication        *bool                `json:"sleepApplication,omitempty"`
	StartCommand            *string              `json:"startCommand,omitempty"`
	UpdatedAt               *string              `json:"updatedAt,omitempty"`
	UpstreamURL             *string              `json:"upstreamUrl,omitempty"`
	WatchPatterns           []string             `json:"watchPatterns,omitempty"`
	Source                  *ServiceSourceDTO    `json:"source,omitempty"`
	LatestDeployment        *LatestDeploymentDTO `json:"latestDeployment,omitempty"`
}

type EnvWithServicesDTO struct {
	ID       string               `json:"id"`
	Name     string               `json:"name"`
	Services []ServiceInstanceDTO `json:"services"`
}

type ProjectDetailsDTO struct {
	ID           string               `json:"id"`
	Name         string               `json:"name"`
	Services     []ProjectDTO         `json:"services"`
	Environments []EnvWithServicesDTO `json:"environments"`
}

// ListRailwayProjects returns projects filtered by comma-separated name list (?names=a,b,c)
// If ?details=1 is provided, returns services/environments for each project.
func (c *EnvironmentController) ListRailwayProjects(ctx *gin.Context) {
	if c.Railway == nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "railway client not configured"})
		return
	}
	details := ctx.Query("details") == "1"
	namesParam := ctx.Query("names")
	nameSet := map[string]struct{}{}
	if namesParam != "" {
		for _, n := range strings.Split(namesParam, ",") {
			clean := strings.ToLower(strings.TrimSpace(n))
			if clean != "" {
				nameSet[clean] = struct{}{}
			}
		}
	}
	if details {
		projects, err := c.Railway.ListProjectsWithDetails(ctx, 200)
		if err != nil {
			log.Error().Err(err).Msg("railway list projects (details) failed")
			ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		log.Debug().Int("pre_filter", len(projects)).Str("names_param", namesParam).Msg("projects details fetched")
		out := make([]ProjectDetailsDTO, 0)
		for _, p := range projects {
			if len(nameSet) > 0 {
				if _, ok := nameSet[strings.ToLower(strings.TrimSpace(p.Name))]; !ok {
					continue
				}
			}
			pd := ProjectDetailsDTO{ID: p.ID, Name: p.Name, Services: []ProjectDTO{}, Environments: []EnvWithServicesDTO{}}
			for _, s := range p.Services {
				pd.Services = append(pd.Services, ProjectDTO{ID: s.ID, Name: s.Name})
			}
			for _, e := range p.Environments {
				env := EnvWithServicesDTO{ID: e.ID, Name: e.Name, Services: []ServiceInstanceDTO{}}
				for _, es := range e.Services {
					dto := ServiceInstanceDTO{
						ID:                      es.ID,
						ServiceID:               es.ServiceID,
						ServiceName:             es.ServiceName,
						EnvironmentID:           es.EnvironmentID,
						BuildCommand:            es.BuildCommand,
						Builder:                 es.Builder,
						CreatedAt:               es.CreatedAt,
						CronSchedule:            es.CronSchedule,
						DeletedAt:               es.DeletedAt,
						DrainingSeconds:         es.DrainingSeconds,
						HealthcheckPath:         es.HealthcheckPath,
						HealthcheckTimeout:      es.HealthcheckTimeout,
						IsUpdatable:             es.IsUpdatable,
						NextCronRunAt:           es.NextCronRunAt,
						NixpacksPlan:            es.NixpacksPlan,
						NumReplicas:             es.NumReplicas,
						OverlapSeconds:          es.OverlapSeconds,
						PreDeployCommand:        es.PreDeployCommand,
						RailpackInfo:            es.RailpackInfo,
						RailwayConfigFile:       es.RailwayConfigFile,
						Region:                  es.Region,
						RestartPolicyMaxRetries: es.RestartPolicyMaxRetries,
						RestartPolicyType:       es.RestartPolicyType,
						RootDirectory:           es.RootDirectory,
						SleepApplication:        es.SleepApplication,
						StartCommand:            es.StartCommand,
						UpdatedAt:               es.UpdatedAt,
						UpstreamURL:             es.UpstreamURL,
						WatchPatterns:           es.WatchPatterns,
					}
					if es.Source != nil {
						dto.Source = &ServiceSourceDTO{
							Image: es.Source.Image,
							Repo:  es.Source.Repo,
						}
					}
					if es.LatestDeployment != nil {
						dto.LatestDeployment = &LatestDeploymentDTO{
							CanRedeploy:             es.LatestDeployment.CanRedeploy,
							CanRollback:             es.LatestDeployment.CanRollback,
							CreatedAt:               es.LatestDeployment.CreatedAt,
							DeploymentStopped:       es.LatestDeployment.DeploymentStopped,
							EnvironmentID:           es.LatestDeployment.EnvironmentID,
							ID:                      es.LatestDeployment.ID,
							Meta:                    es.LatestDeployment.Meta,
							ProjectID:               es.LatestDeployment.ProjectID,
							ServiceID:               es.LatestDeployment.ServiceID,
							SnapshotID:              es.LatestDeployment.SnapshotID,
							StaticURL:               es.LatestDeployment.StaticURL,
							Status:                  es.LatestDeployment.Status,
							StatusUpdatedAt:         es.LatestDeployment.StatusUpdatedAt,
							SuggestAddServiceDomain: es.LatestDeployment.SuggestAddServiceDomain,
							UpdatedAt:               es.LatestDeployment.UpdatedAt,
							URL:                     es.LatestDeployment.URL,
						}
					}
					env.Services = append(env.Services, dto)
				}
				pd.Environments = append(pd.Environments, env)
			}
			out = append(out, pd)
		}
		log.Info().Int("post_filter", len(out)).Msg("projects details returned")
		ctx.JSON(http.StatusOK, out)
		return
	}

	projects, err := c.Railway.ListProjects(ctx, 200)
	if err != nil {
		log.Error().Err(err).Msg("railway list projects failed")
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	log.Debug().Int("pre_filter", len(projects)).Str("names_param", namesParam).Msg("projects fetched")
	out := make([]ProjectDTO, 0)
	for _, p := range projects {
		if len(nameSet) > 0 {
			if _, ok := nameSet[strings.ToLower(strings.TrimSpace(p.Name))]; !ok {
				continue
			}
		}
		out = append(out, ProjectDTO{ID: p.ID, Name: p.Name})
	}
	log.Info().Int("post_filter", len(out)).Msg("projects returned")
	ctx.JSON(http.StatusOK, out)
}

// GetRailwayProject returns a single project by id; if details=1, includes relations.
func (c *EnvironmentController) GetRailwayProject(ctx *gin.Context) {
	if c.Railway == nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "railway client not configured"})
		return
	}
	id := ctx.Param("id")

	// Get authenticated user
	user, err := auth.GetCurrentUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	// Verify user owns at least one environment in this project
	var count int64
	err = c.DB.Model(&store.Environment{}).
		Where("railway_project_id = ? AND user_id = ?", id, user.ID).
		Count(&count).Error
	if err != nil {
		log.Error().Err(err).Str("project_id", id).Msg("failed to verify project ownership")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify ownership"})
		return
	}
	if count == 0 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "project not found or access denied"})
		return
	}

	details := ctx.Query("details") == "1"
	if details {
		p, err := c.Railway.GetProjectWithDetailsByID(ctx, id)
		if err != nil {
			log.Error().Err(err).Str("id", id).Msg("railway get project (details) failed")
			ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		pd := ProjectDetailsDTO{ID: p.ID, Name: p.Name, Services: []ProjectDTO{}, Environments: []EnvWithServicesDTO{}}
		for _, s := range p.Services {
			pd.Services = append(pd.Services, ProjectDTO{ID: s.ID, Name: s.Name})
		}
		for _, e := range p.Environments {
			env := EnvWithServicesDTO{ID: e.ID, Name: e.Name, Services: []ServiceInstanceDTO{}}
			for _, es := range e.Services {
				dto := ServiceInstanceDTO{
					ID:                      es.ID,
					ServiceID:               es.ServiceID,
					ServiceName:             es.ServiceName,
					EnvironmentID:           es.EnvironmentID,
					BuildCommand:            es.BuildCommand,
					Builder:                 es.Builder,
					CreatedAt:               es.CreatedAt,
					CronSchedule:            es.CronSchedule,
					DeletedAt:               es.DeletedAt,
					DrainingSeconds:         es.DrainingSeconds,
					HealthcheckPath:         es.HealthcheckPath,
					HealthcheckTimeout:      es.HealthcheckTimeout,
					IsUpdatable:             es.IsUpdatable,
					NextCronRunAt:           es.NextCronRunAt,
					NixpacksPlan:            es.NixpacksPlan,
					NumReplicas:             es.NumReplicas,
					OverlapSeconds:          es.OverlapSeconds,
					PreDeployCommand:        es.PreDeployCommand,
					RailpackInfo:            es.RailpackInfo,
					RailwayConfigFile:       es.RailwayConfigFile,
					Region:                  es.Region,
					RestartPolicyMaxRetries: es.RestartPolicyMaxRetries,
					RestartPolicyType:       es.RestartPolicyType,
					RootDirectory:           es.RootDirectory,
					SleepApplication:        es.SleepApplication,
					StartCommand:            es.StartCommand,
					UpdatedAt:               es.UpdatedAt,
					UpstreamURL:             es.UpstreamURL,
					WatchPatterns:           es.WatchPatterns,
				}
				if es.Source != nil {
					dto.Source = &ServiceSourceDTO{
						Image: es.Source.Image,
						Repo:  es.Source.Repo,
					}
				}
				if es.LatestDeployment != nil {
					dto.LatestDeployment = &LatestDeploymentDTO{
						CanRedeploy:             es.LatestDeployment.CanRedeploy,
						CanRollback:             es.LatestDeployment.CanRollback,
						CreatedAt:               es.LatestDeployment.CreatedAt,
						DeploymentStopped:       es.LatestDeployment.DeploymentStopped,
						EnvironmentID:           es.LatestDeployment.EnvironmentID,
						ID:                      es.LatestDeployment.ID,
						Meta:                    es.LatestDeployment.Meta,
						ProjectID:               es.LatestDeployment.ProjectID,
						ServiceID:               es.LatestDeployment.ServiceID,
						SnapshotID:              es.LatestDeployment.SnapshotID,
						StaticURL:               es.LatestDeployment.StaticURL,
						Status:                  es.LatestDeployment.Status,
						StatusUpdatedAt:         es.LatestDeployment.StatusUpdatedAt,
						SuggestAddServiceDomain: es.LatestDeployment.SuggestAddServiceDomain,
						UpdatedAt:               es.LatestDeployment.UpdatedAt,
						URL:                     es.LatestDeployment.URL,
					}
				}
				env.Services = append(env.Services, dto)
			}
			pd.Environments = append(pd.Environments, env)
		}
		ctx.JSON(http.StatusOK, pd)
		return
	}
	p, err := c.Railway.GetProject(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("railway get project failed")
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, ProjectDTO{ID: p.ID, Name: p.Name})
}

// ProvisionProjectRequest is the payload for creating a new Railway project
// through our API. Name is optional based on GraphQL support; defaultEnvironmentName
// allows specifying the initial environment name.
type ProvisionProjectRequest struct {
	DefaultEnvironmentName *string `json:"defaultEnvironmentName"`
	Name                   *string `json:"name"`
	RequestID              string  `json:"requestId"`
}

type ProvisionProjectResponse struct {
	ProjectID            string `json:"projectId"`
	BaseEnvironmentID    string `json:"baseEnvironmentId"`    // Mirage internal environment ID
	RailwayEnvironmentID string `json:"railwayEnvironmentId"` // Railway's environment ID
	Name                 string `json:"name"`
}

// ProvisionProject creates a new Railway project by delegating to the railway client.
// After successful creation, explicitly fetches the default environment and persists
// both the environment and metadata to our database atomically.
func (c *EnvironmentController) ProvisionProject(ctx *gin.Context) {
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

	var req ProvisionProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Step 1: Create the Railway project
	res, err := c.Railway.CreateProject(ctx, railway.CreateProjectInput{DefaultEnvironmentName: req.DefaultEnvironmentName, Name: req.Name})
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// Step 2: Explicitly fetch the default environment from Railway
	// Railway mutation responses can be unreliable, so we explicitly query for the environment
	railwayEnvID := res.BaseEnvironmentID
	envName := "production" // Default Railway environment name
	if req.DefaultEnvironmentName != nil {
		envName = *req.DefaultEnvironmentName
	}

	// If mutation didn't return environment ID, explicitly fetch it
	if railwayEnvID == "" {
		log.Info().
			Str("project_id", res.ProjectID).
			Msg("base environment ID not in mutation response, fetching explicitly")

		pd, err := c.Railway.GetProjectWithDetailsByID(ctx, res.ProjectID)
		if err != nil {
			log.Error().Err(err).Str("project_id", res.ProjectID).Msg("failed to fetch project details for default environment")
			ctx.JSON(http.StatusBadGateway, gin.H{"error": "project created but failed to retrieve default environment: " + err.Error()})
			return
		}

		if len(pd.Environments) == 0 {
			log.Error().Str("project_id", res.ProjectID).Msg("project created but has no environments")
			ctx.JSON(http.StatusBadGateway, gin.H{"error": "project created but no default environment found"})
			return
		}

		// Find the environment by name (or take first one)
		for _, env := range pd.Environments {
			if env.Name == envName {
				railwayEnvID = env.ID
				break
			}
		}
		if railwayEnvID == "" {
			railwayEnvID = pd.Environments[0].ID
			envName = pd.Environments[0].Name
			log.Warn().
				Str("project_id", res.ProjectID).
				Str("expected_name", envName).
				Str("actual_name", pd.Environments[0].Name).
				Msg("default environment name mismatch, using first environment")
		}
	}

	// Step 3: Persist the base environment and metadata to database
	if c.DB != nil {
		env := store.Environment{
			ID:                   uuid.New().String(),
			UserID:               user.ID, // Set from authenticated user
			Name:                 envName,
			Type:                 store.EnvironmentTypeProd, // Base environment defaults to prod
			Status:               status.StatusCreating,
			RailwayProjectID:     res.ProjectID,
			RailwayEnvironmentID: railwayEnvID,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}

		// Use transaction to ensure Environment and EnvironmentMetadata are created atomically
		txErr := c.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&env).Error; err != nil {
				return err
			}

			// Always create EnvironmentMetadata with provision outputs
			// This ensures we can track and clone environments even without wizard inputs
			provisionOutputs := map[string]interface{}{
				"projectId":            res.ProjectID,
				"railwayEnvironmentId": railwayEnvID,
				"environmentName":      envName,
			}
			provisionOutputsJSON, _ := json.Marshal(provisionOutputs)

			metadata := store.EnvironmentMetadata{
				ID:                   uuid.New().String(),
				UserID:               user.ID, // Set from authenticated user
				EnvironmentID:        env.ID,
				WizardInputsJSON:     nil, // No wizard inputs for project creation
				ProvisionOutputsJSON: provisionOutputsJSON,
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			}

			if err := tx.Create(&metadata).Error; err != nil {
				return err
			}

			log.Info().
				Str("mirage_env_id", env.ID).
				Str("metadata_id", metadata.ID).
				Str("railway_project_id", res.ProjectID).
				Str("railway_env_id", railwayEnvID).
				Msg("persisted base environment and metadata to database")

			return nil
		})

		if txErr != nil {
			log.Error().Err(txErr).
				Str("project_id", res.ProjectID).
				Str("railway_env_id", railwayEnvID).
				Msg("failed to persist environment to database after Railway project creation")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "project created but failed to persist to database: " + txErr.Error()})
			return
		}

		// Return Mirage environment ID for frontend use (foreign keys)
		ctx.JSON(http.StatusOK, ProvisionProjectResponse{
			ProjectID:            res.ProjectID,
			BaseEnvironmentID:    env.ID,       // Mirage ID for foreign keys
			RailwayEnvironmentID: railwayEnvID, // Railway ID for Railway API calls
			Name:                 res.Name,
		})
		return
	}

	// If DB is nil, return Railway ID for backward compatibility
	ctx.JSON(http.StatusOK, ProvisionProjectResponse{
		ProjectID:            res.ProjectID,
		BaseEnvironmentID:    railwayEnvID,
		RailwayEnvironmentID: railwayEnvID,
		Name:                 res.Name,
	})
}

// DeleteRailwayEnvironment deletes a Railway environment by its Railway environment ID.
// After successful Railway deletion, cleans up the database: services, metadata, and environment records.
func (c *EnvironmentController) DeleteRailwayEnvironment(ctx *gin.Context) {
	if c.Railway == nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "railway client not configured"})
		return
	}
	railwayEnvID := ctx.Param("id")
	if railwayEnvID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "environment id required"})
		return
	}

	// Get authenticated user
	user, err := auth.GetCurrentUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	// Verify ownership before deleting
	var env store.Environment
	err = c.DB.Where("railway_environment_id = ? AND user_id = ?", railwayEnvID, user.ID).First(&env).Error
	if err == gorm.ErrRecordNotFound {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "environment not found or access denied"})
		return
	} else if err != nil {
		log.Error().Err(err).Str("railway_env_id", railwayEnvID).Msg("failed to verify environment ownership")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify ownership"})
		return
	}

	// Step 1: Delete from Railway first (fail fast if Railway API fails)
	log.Info().
		Str("railway_env_id", railwayEnvID).
		Str("user_id", user.ID).
		Msg("deleting railway environment")
	if err := c.Railway.DestroyEnvironment(ctx, railway.DestroyEnvironmentInput{EnvironmentID: railwayEnvID}); err != nil {
		log.Error().Err(err).Str("railway_env_id", railwayEnvID).Msg("railway delete environment failed")
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// Step 2: Clean up database (Railway deletion succeeded, so clean up our records)
	// Note: env was already fetched above with ownership verification
	if c.DB != nil {

		// Use transaction to ensure atomic cleanup
		txErr := c.DB.Transaction(func(tx *gorm.DB) error {
			// Delete all services for this environment
			result := tx.Where("environment_id = ?", env.ID).Delete(&store.Service{})
			if result.Error != nil {
				return result.Error
			}
			log.Info().
				Str("mirage_env_id", env.ID).
				Int64("services_deleted", result.RowsAffected).
				Msg("deleted services from database")

			// Delete metadata for this environment (may not exist for old environments)
			result = tx.Where("environment_id = ?", env.ID).Delete(&store.EnvironmentMetadata{})
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected > 0 {
				log.Info().
					Str("mirage_env_id", env.ID).
					Msg("deleted environment metadata from database")
			}

			// Delete the environment itself
			if err := tx.Delete(&env).Error; err != nil {
				return err
			}
			log.Info().
				Str("mirage_env_id", env.ID).
				Str("railway_env_id", railwayEnvID).
				Msg("deleted environment from database")

			return nil
		})

		if txErr != nil {
			log.Error().Err(txErr).
				Str("railway_env_id", railwayEnvID).
				Str("mirage_env_id", env.ID).
				Msg("failed to clean up database after Railway environment deletion")
			// Don't fail the request - Railway resource was deleted successfully
		}
	}

	ctx.Status(http.StatusNoContent)
}

// DeleteRailwayProject deletes a Railway project and all its associated resources.
// WARNING: This is a destructive operation that cannot be undone.
// After successful Railway deletion, cleans up the database: all environments, metadata, and services.
func (c *EnvironmentController) DeleteRailwayProject(ctx *gin.Context) {
	if c.Railway == nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "railway client not configured"})
		return
	}
	projectID := ctx.Param("id")
	if projectID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "project id required"})
		return
	}

	// Get authenticated user
	user, err := auth.GetCurrentUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	// Verify user owns at least one environment in this project
	var count int64
	err = c.DB.Model(&store.Environment{}).
		Where("railway_project_id = ? AND user_id = ?", projectID, user.ID).
		Count(&count).Error
	if err != nil {
		log.Error().Err(err).Str("project_id", projectID).Msg("failed to verify project ownership")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify ownership"})
		return
	}
	if count == 0 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "project not found or access denied"})
		return
	}

	// Step 1: Delete from Railway first (fail fast if Railway API fails)
	log.Warn().
		Str("project_id", projectID).
		Str("user_id", user.ID).
		Msg("deleting railway project - irreversible operation")
	if err := c.Railway.DestroyProject(ctx, railway.DestroyProjectInput{ProjectID: projectID}); err != nil {
		log.Error().Err(err).Str("project_id", projectID).Msg("railway delete project failed")
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// Step 2: Clean up database (Railway deletion succeeded, so clean up all project resources)
	if c.DB != nil {
		// Find all environments associated with this Railway project owned by this user
		var envs []store.Environment
		if err := c.DB.Where("railway_project_id = ? AND user_id = ?", projectID, user.ID).Find(&envs).Error; err != nil {
			log.Error().Err(err).Str("project_id", projectID).Msg("failed to query environments for cleanup")
			// Don't fail the request - Railway resource was deleted successfully
			ctx.Status(http.StatusNoContent)
			return
		}

		if len(envs) == 0 {
			log.Info().Str("project_id", projectID).Msg("no environments found in database for this project")
			ctx.Status(http.StatusNoContent)
			return
		}

		// Use transaction to ensure atomic cleanup of all environments and their resources
		txErr := c.DB.Transaction(func(tx *gorm.DB) error {
			totalServices := int64(0)
			totalMetadata := int64(0)

			for _, env := range envs {
				// Delete all services for this environment
				result := tx.Where("environment_id = ?", env.ID).Delete(&store.Service{})
				if result.Error != nil {
					return result.Error
				}
				totalServices += result.RowsAffected

				// Delete metadata for this environment (may not exist)
				result = tx.Where("environment_id = ?", env.ID).Delete(&store.EnvironmentMetadata{})
				if result.Error != nil {
					return result.Error
				}
				totalMetadata += result.RowsAffected

				// Delete the environment itself
				if err := tx.Delete(&env).Error; err != nil {
					return err
				}
			}

			log.Info().
				Str("project_id", projectID).
				Int("environments_deleted", len(envs)).
				Int64("services_deleted", totalServices).
				Int64("metadata_deleted", totalMetadata).
				Msg("deleted all project resources from database")

			return nil
		})

		if txErr != nil {
			log.Error().Err(txErr).
				Str("project_id", projectID).
				Int("environments_found", len(envs)).
				Msg("failed to clean up database after Railway project deletion")
			// Don't fail the request - Railway resource was deleted successfully
		}
	}

	ctx.Status(http.StatusNoContent)
}
