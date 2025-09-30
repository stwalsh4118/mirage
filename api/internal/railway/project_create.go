package railway

import "context"

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
    baseEnvironmentId
  }
}`
	vars := map[string]any{
		"defaultEnvironmentName": in.DefaultEnvironmentName,
		"name":                   in.Name,
	}
	var resp struct {
		ProjectCreate struct {
			ID                string `json:"id"`
			Name              string `json:"name"`
			BaseEnvironmentID string `json:"baseEnvironmentId"`
		} `json:"projectCreate"`
	}
	if err := c.execute(ctx, mutation, vars, &resp); err != nil {
		return CreateProjectResult{}, err
	}
	return CreateProjectResult{ProjectID: resp.ProjectCreate.ID, BaseEnvironmentID: resp.ProjectCreate.BaseEnvironmentID, Name: resp.ProjectCreate.Name}, nil
}
