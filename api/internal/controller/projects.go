package controller

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stwalsh4118/mirageapi/internal/railway"
	"github.com/stwalsh4118/mirageapi/internal/status"
	"github.com/stwalsh4118/mirageapi/internal/store"
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
	ProjectID         string `json:"projectId"`
	BaseEnvironmentID string `json:"baseEnvironmentId"`
	Name              string `json:"name"`
}

// ProvisionProject creates a new Railway project by delegating to the railway client.
// After successful creation, persists the base environment to our database.
func (c *EnvironmentController) ProvisionProject(ctx *gin.Context) {
	if c.Railway == nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "railway client not configured"})
		return
	}
	var req ProvisionProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := c.Railway.CreateProject(ctx, railway.CreateProjectInput{DefaultEnvironmentName: req.DefaultEnvironmentName, Name: req.Name})
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// Persist the base environment to database
	if c.DB != nil {
		envName := "production" // Default Railway environment name
		if req.DefaultEnvironmentName != nil {
			envName = *req.DefaultEnvironmentName
		}

		env := store.Environment{
			ID:                   uuid.New().String(),
			Name:                 envName,
			Type:                 store.EnvironmentTypeProd, // Base environment defaults to prod
			Status:               status.StatusCreating,
			RailwayProjectID:     res.ProjectID,
			RailwayEnvironmentID: res.BaseEnvironmentID,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}

		if err := c.DB.Create(&env).Error; err != nil {
			log.Error().Err(err).
				Str("project_id", res.ProjectID).
				Str("env_id", res.BaseEnvironmentID).
				Msg("failed to persist environment to database after Railway project creation")
			// Don't fail the request - Railway resource was created successfully
		} else {
			log.Info().
				Str("env_id", env.ID).
				Str("railway_project_id", res.ProjectID).
				Str("railway_env_id", res.BaseEnvironmentID).
				Msg("persisted base environment to database")
		}
	}

	ctx.JSON(http.StatusOK, ProvisionProjectResponse{ProjectID: res.ProjectID, BaseEnvironmentID: res.BaseEnvironmentID, Name: res.Name})
}

// DeleteRailwayEnvironment deletes a Railway environment by its Railway environment ID.
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

	log.Info().Str("railway_env_id", railwayEnvID).Msg("deleting railway environment")
	if err := c.Railway.DestroyEnvironment(ctx, railway.DestroyEnvironmentInput{EnvironmentID: railwayEnvID}); err != nil {
		log.Error().Err(err).Str("railway_env_id", railwayEnvID).Msg("railway delete environment failed")
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// DeleteRailwayProject deletes a Railway project and all its associated resources.
// WARNING: This is a destructive operation that cannot be undone.
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

	log.Warn().Str("project_id", projectID).Msg("deleting railway project - irreversible operation")
	if err := c.Railway.DestroyProject(ctx, railway.DestroyProjectInput{ProjectID: projectID}); err != nil {
		log.Error().Err(err).Str("project_id", projectID).Msg("railway delete project failed")
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
