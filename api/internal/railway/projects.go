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

// ServiceSource represents the source configuration for a service (repo or image)
type ServiceSource struct {
	Image *string `json:"image,omitempty"`
	Repo  *string `json:"repo,omitempty"`
}

// LatestDeployment represents the latest deployment information for a service
type LatestDeployment struct {
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

// ServiceInstance represents a service instance within an environment with detailed configuration
type ServiceInstance struct {
	ID                      string            `json:"id"`
	ServiceID               string            `json:"serviceId"`
	ServiceName             string            `json:"serviceName"`
	EnvironmentID           string            `json:"environmentId"`
	BuildCommand            *string           `json:"buildCommand,omitempty"`
	Builder                 *string           `json:"builder,omitempty"`
	CreatedAt               *string           `json:"createdAt,omitempty"`
	CronSchedule            *string           `json:"cronSchedule,omitempty"`
	DeletedAt               *string           `json:"deletedAt,omitempty"`
	DrainingSeconds         *int              `json:"drainingSeconds,omitempty"`
	HealthcheckPath         *string           `json:"healthcheckPath,omitempty"`
	HealthcheckTimeout      *int              `json:"healthcheckTimeout,omitempty"`
	IsUpdatable             *bool             `json:"isUpdatable,omitempty"`
	NextCronRunAt           *string           `json:"nextCronRunAt,omitempty"`
	NixpacksPlan            *string           `json:"nixpacksPlan,omitempty"`
	NumReplicas             *int              `json:"numReplicas,omitempty"`
	OverlapSeconds          *int              `json:"overlapSeconds,omitempty"`
	PreDeployCommand        *string           `json:"preDeployCommand,omitempty"`
	RailpackInfo            *string           `json:"railpackInfo,omitempty"`
	RailwayConfigFile       *string           `json:"railwayConfigFile,omitempty"`
	Region                  *string           `json:"region,omitempty"`
	RestartPolicyMaxRetries *int              `json:"restartPolicyMaxRetries,omitempty"`
	RestartPolicyType       *string           `json:"restartPolicyType,omitempty"`
	RootDirectory           *string           `json:"rootDirectory,omitempty"`
	SleepApplication        *bool             `json:"sleepApplication,omitempty"`
	StartCommand            *string           `json:"startCommand,omitempty"`
	UpdatedAt               *string           `json:"updatedAt,omitempty"`
	UpstreamURL             *string           `json:"upstreamUrl,omitempty"`
	WatchPatterns           []string          `json:"watchPatterns,omitempty"`
	Source                  *ServiceSource    `json:"source,omitempty"`
	LatestDeployment        *LatestDeployment `json:"latestDeployment,omitempty"`
}

// ProjectEnvironment represents an environment within a project with its services
type ProjectEnvironment struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Services []ServiceInstance `json:"services"`
}

type ProjectDetails struct {
	ID           string               `json:"id"`
	Name         string               `json:"name"`
	Services     []ProjectItem        `json:"services"`
	Plugins      []ProjectItem        `json:"plugins"`
	Environments []ProjectEnvironment `json:"environments"`
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
    environments {
      edges {
        node {
          id
          name
          serviceInstances {
            edges {
              node {
                buildCommand
                builder
                createdAt
                cronSchedule
                deletedAt
                drainingSeconds
                environmentId
                healthcheckPath
                healthcheckTimeout
                id
                isUpdatable
                nextCronRunAt
                nixpacksPlan
                numReplicas
                overlapSeconds
                preDeployCommand
                railpackInfo
                railwayConfigFile
                region
                restartPolicyMaxRetries
                restartPolicyType
                rootDirectory
                serviceId
                serviceName
                sleepApplication
                startCommand
                updatedAt
                upstreamUrl
                watchPatterns
                source {
                  image
                  repo
                }
                latestDeployment {
                  canRedeploy
                  canRollback
                  createdAt
                  deploymentStopped
                  environmentId
                  id
                  meta
                  projectId
                  serviceId
                  snapshotId
                  staticUrl
                  status
                  statusUpdatedAt
                  suggestAddServiceDomain
                  updatedAt
                  url
                }
              }
            }
          }
        }
      }
    }
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
        environments {
          edges {
            node {
              id
              name
              serviceInstances {
                edges {
                  node {
                    buildCommand
					builder
					createdAt
					cronSchedule
					deletedAt
					drainingSeconds
					environmentId
					healthcheckPath
					healthcheckTimeout
					id
					isUpdatable
					nextCronRunAt
					nixpacksPlan
					numReplicas
					overlapSeconds
					preDeployCommand
					railpackInfo
					railwayConfigFile
					region
					restartPolicyMaxRetries
					restartPolicyType
					rootDirectory
					serviceId
					serviceName
					sleepApplication
					startCommand
					updatedAt
					upstreamUrl
					watchPatterns
					source {
						image
						repo
					}
					latestDeployment {
						canRedeploy
						canRollback
						createdAt
						deploymentStopped
						environmentId
						id
						meta
						projectId
						serviceId
						snapshotId
						staticUrl
						status
						statusUpdatedAt
						suggestAddServiceDomain
						updatedAt
						url
					}
                  }
                }
              }
            }
          }
        }
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
					Node struct {
						ID               string `json:"id"`
						Name             string `json:"name"`
						ServiceInstances struct {
							Edges []struct {
								Node ServiceInstance `json:"node"`
							} `json:"edges"`
						} `json:"serviceInstances"`
					} `json:"node"`
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
		env := ProjectEnvironment{ID: ee.Node.ID, Name: ee.Node.Name}
		for _, sie := range ee.Node.ServiceInstances.Edges {
			env.Services = append(env.Services, sie.Node)
		}
		pd.Environments = append(pd.Environments, env)
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
	if err := c.execute(ctx, gqlListProjectsRoot, map[string]any{"first": first}, &outRoot); err != nil {
		log.Error().Err(err).Str("query", "ListProjects_root").Msg("railway root.projects query failed")
		return nil, err
	}
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
				Node struct {
					ID               string `json:"id"`
					Name             string `json:"name"`
					ServiceInstances struct {
						Edges []struct {
							Node ServiceInstance `json:"node"`
						} `json:"edges"`
					} `json:"serviceInstances"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"environments"`
	}) {
		pd := acc[p.ID]
		pd.ID = p.ID
		pd.Name = p.Name
		if pd.Services == nil {
			pd.Services = []ProjectItem{}
		}
		for _, se := range p.Services.Edges {
			pd.Services = append(pd.Services, se.Node)
		}
		if pd.Plugins == nil {
			pd.Plugins = []ProjectItem{}
		}
		for _, pe := range p.Plugins.Edges {
			pd.Plugins = append(pd.Plugins, pe.Node)
		}
		if pd.Environments == nil {
			pd.Environments = []ProjectEnvironment{}
		}
		for _, ee := range p.Environments.Edges {
			env := ProjectEnvironment{ID: ee.Node.ID, Name: ee.Node.Name}
			for _, sie := range ee.Node.ServiceInstances.Edges {
				env.Services = append(env.Services, sie.Node)
			}
			pd.Environments = append(pd.Environments, env)
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
							Node struct {
								ID               string `json:"id"`
								Name             string `json:"name"`
								ServiceInstances struct {
									Edges []struct {
										Node ServiceInstance `json:"node"`
									} `json:"edges"`
								} `json:"serviceInstances"`
							} `json:"node"`
						} `json:"edges"`
					} `json:"environments"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"projects"`
	}
	if err := c.execute(ctx, gqlProjectsDetailsRoot, map[string]any{"first": first}, &outRoot); err != nil {
		log.Error().Err(err).Str("query", "ProjectsDetails_root").Msg("railway root.projects details query failed")
		return nil, err
	}
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

// CreateProjectInput contains optional parameters for creating a project.
type CreateProjectInput struct {
	DefaultEnvironmentName *string
	Name                   *string
}

// CreateProjectResult captures identifiers returned by Railway when creating a project.
type CreateProjectResult struct {
	ProjectID         string
	BaseEnvironmentID string
	Name              string
}

// CreateProject executes the projectCreate mutation.
func (c *Client) CreateProject(ctx context.Context, in CreateProjectInput) (CreateProjectResult, error) {
	mutation := `mutation ProjectCreate($defaultEnvironmentName: String, $name: String) {
  projectCreate(input: { defaultEnvironmentName: $defaultEnvironmentName, name: $name }) {
    id
    name
    environments {
      edges {
        cursor
        node { id name }
      }
    }
  }
}`
	vars := map[string]any{
		"defaultEnvironmentName": in.DefaultEnvironmentName,
		"name":                   in.Name,
	}
	var resp struct {
		ProjectCreate struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			Environments struct {
				Edges []struct {
					Cursor string `json:"cursor"`
					Node   struct {
						ID   string `json:"id"`
						Name string `json:"name"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"environments"`
		} `json:"projectCreate"`
	}
	if err := c.execute(ctx, mutation, vars, &resp); err != nil {
		return CreateProjectResult{}, err
	}
	log.Info().Interface("resp", resp).Msg("creating project")

	envID := ""
	if len(resp.ProjectCreate.Environments.Edges) > 0 {
		envID = resp.ProjectCreate.Environments.Edges[0].Node.ID
	}
	// Fallback: query project environments if still empty
	if envID == "" {
		pd, err := c.GetProjectWithDetailsByID(ctx, resp.ProjectCreate.ID)
		if err != nil {
			log.Warn().Err(err).Str("project_id", resp.ProjectCreate.ID).Msg("fallback fetch project details failed")
		} else if len(pd.Environments) > 0 {
			envID = pd.Environments[0].ID
		}
	}
	return CreateProjectResult{ProjectID: resp.ProjectCreate.ID, BaseEnvironmentID: envID, Name: resp.ProjectCreate.Name}, nil
}

// DestroyProjectInput carries the project identifier.
type DestroyProjectInput struct {
	ProjectID string
}

// DestroyProject removes a project and all its associated resources.
// WARNING: This is a destructive operation that cannot be undone.
// All environments, services, and data within the project will be deleted.
func (c *Client) DestroyProject(ctx context.Context, in DestroyProjectInput) error {
	mutation := `mutation ProjectDelete($projectId: String!) {
  projectDelete(id: $projectId)
}`
	vars := map[string]any{
		"projectId": in.ProjectID,
	}

	log.Warn().
		Str("project_id", in.ProjectID).
		Msg("deleting Railway project - this operation is irreversible")

	var resp struct {
		ProjectDelete bool `json:"projectDelete"`
	}
	return c.execute(ctx, mutation, vars, &resp)
}
