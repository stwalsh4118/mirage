package railway

import (
	"context"

	"github.com/rs/zerolog/log"
)

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
