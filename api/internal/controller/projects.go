package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type ProjectDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProjectDetailsDTO struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Services     []ProjectDTO `json:"services"`
	Plugins      []ProjectDTO `json:"plugins"`
	Environments []ProjectDTO `json:"environments"`
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
			pd := ProjectDetailsDTO{ID: p.ID, Name: p.Name}
			if pd.Services == nil {
				pd.Services = []ProjectDTO{}
			}
			if pd.Plugins == nil {
				pd.Plugins = []ProjectDTO{}
			}
			if pd.Environments == nil {
				pd.Environments = []ProjectDTO{}
			}
			for _, s := range p.Services {
				pd.Services = append(pd.Services, ProjectDTO{ID: s.ID, Name: s.Name})
			}
			for _, g := range p.Plugins {
				pd.Plugins = append(pd.Plugins, ProjectDTO{ID: g.ID, Name: g.Name})
			}
			for _, e := range p.Environments {
				pd.Environments = append(pd.Environments, ProjectDTO{ID: e.ID, Name: e.Name})
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
		pd := ProjectDetailsDTO{ID: p.ID, Name: p.Name, Services: []ProjectDTO{}, Plugins: []ProjectDTO{}, Environments: []ProjectDTO{}}
		for _, s := range p.Services {
			pd.Services = append(pd.Services, ProjectDTO{ID: s.ID, Name: s.Name})
		}
		for _, g := range p.Plugins {
			pd.Plugins = append(pd.Plugins, ProjectDTO{ID: g.ID, Name: g.Name})
		}
		for _, e := range p.Environments {
			pd.Environments = append(pd.Environments, ProjectDTO{ID: e.ID, Name: e.Name})
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
