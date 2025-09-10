package railway

import (
	"context"
	"strings"

	"github.com/rs/zerolog/log"
)

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProjectItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProjectDetails struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Services     []ProjectItem `json:"services"`
	Plugins      []ProjectItem `json:"plugins"`
	Environments []ProjectItem `json:"environments"`
}

const (
	defaultProjectListPageSize = 100

	gqlProjectByID = `
query Project($id: ID!) {
  project(id: $id) {
    id
    name
  }
}
`

	gqlProjectDetailsByID = `
query ProjectDetails($id: ID!) {
  project(id: $id) {
    id
    name
    services { edges { node { id name } } }
    plugins { edges { node { id name } } }
    environments { edges { node { id name } } }
  }
}
`

	gqlListProjectsRoot = `
query ListProjects_root($first: Int!) {
  projects(first: $first) {
    edges {
      node { id name }
    }
  }
}
`

	gqlProjectsDetailsRoot = `
query ProjectsDetails_root($first: Int!) {
  projects(first: $first) {
    edges {
      node {
        id
        name
        services { edges { node { id name } } }
        plugins { edges { node { id name } } }
        environments { edges { node { id name } } }
      }
    }
  }
}
`
)

// GetProject fetches a single project by ID (id, name only).
func (c *Client) GetProject(ctx context.Context, id string) (Project, error) {
	var out struct {
		Project struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"project"`
	}
	gql := gqlProjectByID
	if err := c.execute(ctx, gql, map[string]any{"id": id}, &out); err != nil {
		return Project{}, err
	}
	return Project{ID: out.Project.ID, Name: out.Project.Name}, nil
}

// GetProjectWithDetailsByID fetches one project with services/plugins/environments.
func (c *Client) GetProjectWithDetailsByID(ctx context.Context, id string) (ProjectDetails, error) {
	var out struct {
		Project struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Services struct {
				Edges []struct {
					Node ProjectItem `json:"node"`
				} `json:"edges"`
			} `json:"services"`
			Plugins struct {
				Edges []struct {
					Node ProjectItem `json:"node"`
				} `json:"edges"`
			} `json:"plugins"`
			Environments struct {
				Edges []struct {
					Node ProjectItem `json:"node"`
				} `json:"edges"`
			} `json:"environments"`
		} `json:"project"`
	}
	gql := gqlProjectDetailsByID
	if err := c.execute(ctx, gql, map[string]any{"id": id}, &out); err != nil {
		return ProjectDetails{}, err
	}
	pd := ProjectDetails{ID: out.Project.ID, Name: out.Project.Name}
	for _, se := range out.Project.Services.Edges {
		pd.Services = append(pd.Services, se.Node)
	}
	for _, pe := range out.Project.Plugins.Edges {
		pd.Plugins = append(pd.Plugins, pe.Node)
	}
	for _, ee := range out.Project.Environments.Edges {
		pd.Environments = append(pd.Environments, ee.Node)
	}
	return pd, nil
}

// ListProjects fetches up to N projects for the current token using multiple discovery paths
// (viewer.projects, root projects, and viewer.teams[].projects). Results are merged.
func (c *Client) ListProjects(ctx context.Context, first int) ([]Project, error) {
	if first <= 0 {
		first = defaultProjectListPageSize
	}
	acc := make(map[string]Project)

	// removed me.projects usage (workspace tokens don't expose `me`)

	// removed viewer.projects usage per workspace token model

	// root projects connection (if available)
	var outRoot struct {
		Projects struct {
			Edges []struct {
				Node Project `json:"node"`
			} `json:"edges"`
		} `json:"projects"`
	}
	_ = c.execute(ctx, gqlListProjectsRoot, map[string]any{"first": first}, &outRoot)
	rootCount := 0
	rootNames := make([]string, 0)
	for _, e := range outRoot.Projects.Edges {
		acc[e.Node.ID] = e.Node
		rootCount++
		if len(rootNames) < 5 {
			rootNames = append(rootNames, e.Node.Name)
		}
	}
	log.Debug().Int("count", rootCount).Str("sample", strings.Join(rootNames, ", ")).Msg("railway root.projects")

	// removed viewer.teams[].projects usage per workspace token model

	projects := make([]Project, 0, len(acc))
	namePeek := make([]string, 0)
	for _, p := range acc {
		projects = append(projects, p)
		if len(namePeek) < 5 {
			namePeek = append(namePeek, p.Name)
		}
	}
	log.Info().Int("total", len(projects)).Str("sample", strings.Join(namePeek, ", ")).Msg("railway projects merged")
	return projects, nil
}

// ListProjectsWithDetails returns projects visible to the token along with
// services, plugins, and environments (ids and names only). It queries across
// viewer root and team projects and merges results.
func (c *Client) ListProjectsWithDetails(ctx context.Context, first int) ([]ProjectDetails, error) {
	if first <= 0 {
		first = defaultProjectListPageSize
	}
	acc := make(map[string]ProjectDetails)

	merge := func(p struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Services struct {
			Edges []struct {
				Node ProjectItem `json:"node"`
			} `json:"edges"`
		} `json:"services"`
		Plugins struct {
			Edges []struct {
				Node ProjectItem `json:"node"`
			} `json:"edges"`
		} `json:"plugins"`
		Environments struct {
			Edges []struct {
				Node ProjectItem `json:"node"`
			} `json:"edges"`
		} `json:"environments"`
	}) {
		pd := acc[p.ID]
		pd.ID = p.ID
		pd.Name = p.Name
		pd.Services = []ProjectItem{}
		for _, se := range p.Services.Edges {
			pd.Services = append(pd.Services, se.Node)
		}
		pd.Plugins = []ProjectItem{}
		for _, pe := range p.Plugins.Edges {
			pd.Plugins = append(pd.Plugins, pe.Node)
		}
		pd.Environments = []ProjectItem{}
		for _, ee := range p.Environments.Edges {
			pd.Environments = append(pd.Environments, ee.Node)
		}
		acc[p.ID] = pd
	}

	// removed me.projects with details (workspace tokens don't expose `me`)

	// root projects with details
	var outRoot struct {
		Projects struct {
			Edges []struct {
				Node struct {
					ID       string `json:"id"`
					Name     string `json:"name"`
					Services struct {
						Edges []struct {
							Node ProjectItem `json:"node"`
						} `json:"edges"`
					} `json:"services"`
					Plugins struct {
						Edges []struct {
							Node ProjectItem `json:"node"`
						} `json:"edges"`
					} `json:"plugins"`
					Environments struct {
						Edges []struct {
							Node ProjectItem `json:"node"`
						} `json:"edges"`
					} `json:"environments"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"projects"`
	}
	_ = c.execute(ctx, gqlProjectsDetailsRoot, map[string]any{"first": first}, &outRoot)
	rootNames := make([]string, 0)
	for _, e := range outRoot.Projects.Edges {
		merge(e.Node)
		if len(rootNames) < 5 {
			rootNames = append(rootNames, e.Node.Name)
		}
	}
	log.Debug().Int("count", len(outRoot.Projects.Edges)).Str("sample", strings.Join(rootNames, ", ")).Msg("railway root.projects details")

	result := make([]ProjectDetails, 0, len(acc))
	peek := make([]string, 0)
	for _, p := range acc {
		result = append(result, p)
		if len(peek) < 5 {
			peek = append(peek, p.Name)
		}
	}
	log.Info().Int("total", len(result)).Str("sample", strings.Join(peek, ", ")).Msg("railway projects merged (details)")
	return result, nil
}
