package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/stwalsh4118/mirageapi/internal/railway"
)

type ProjectDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type EnvWithServicesDTO struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Services []ProjectDTO `json:"services"`
}

type ProjectDetailsDTO struct {
	ID           string               `json:"id"`
	Name         string               `json:"name"`
	Services     []ProjectDTO         `json:"services"`
	Plugins      []ProjectDTO         `json:"plugins"`
	Environments []EnvWithServicesDTO `json:"environments"`
}

// ListRailwayProjects returns projects filtered by comma-separated name list (?names=a,b,c)
// If ?details=1 is provided, returns services/plugins/environments for each project.
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
			pd := ProjectDetailsDTO{ID: p.ID, Name: p.Name, Services: []ProjectDTO{}, Plugins: []ProjectDTO{}, Environments: []EnvWithServicesDTO{}}
			for _, s := range p.Services {
				pd.Services = append(pd.Services, ProjectDTO{ID: s.ID, Name: s.Name})
			}
			for _, g := range p.Plugins {
				pd.Plugins = append(pd.Plugins, ProjectDTO{ID: g.ID, Name: g.Name})
			}
			for _, e := range p.Environments {
				env := EnvWithServicesDTO{ID: e.ID, Name: e.Name, Services: []ProjectDTO{}}
				for _, es := range e.Services {
					env.Services = append(env.Services, ProjectDTO{ID: es.ID, Name: es.Name})
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
		pd := ProjectDetailsDTO{ID: p.ID, Name: p.Name, Services: []ProjectDTO{}, Plugins: []ProjectDTO{}, Environments: []EnvWithServicesDTO{}}
		for _, s := range p.Services {
			pd.Services = append(pd.Services, ProjectDTO{ID: s.ID, Name: s.Name})
		}
		for _, g := range p.Plugins {
			pd.Plugins = append(pd.Plugins, ProjectDTO{ID: g.ID, Name: g.Name})
		}
		for _, e := range p.Environments {
			env := EnvWithServicesDTO{ID: e.ID, Name: e.Name, Services: []ProjectDTO{}}
			for _, es := range e.Services {
				env.Services = append(env.Services, ProjectDTO{ID: es.ID, Name: es.Name})
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
