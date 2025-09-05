package railway

import "context"

// CreateEnvironmentInput contains minimal fields to create an environment.
type CreateEnvironmentInput struct {
	ProjectID string
	Name      string
}

// CreateEnvironmentResult captures IDs returned by Railway.
type CreateEnvironmentResult struct {
	EnvironmentID string
}

// CreateEnvironment executes the create environment mutation.
func (c *Client) CreateEnvironment(ctx context.Context, in CreateEnvironmentInput) (CreateEnvironmentResult, error) {
	// TODO: Replace with actual mutation and variables once confirmed.
	mutation := `mutation CreateEnv($projectId: ID!, $name: String!) {\n  createEnvironment(input: { projectId: $projectId, name: $name }) {\n    environment { id }\n  }\n}`
	vars := map[string]any{
		"projectId": in.ProjectID,
		"name":      in.Name,
	}
	var resp struct {
		CreateEnvironment struct {
			Environment struct{ ID string } `json:"environment"`
		} `json:"createEnvironment"`
	}
	if err := c.execute(ctx, mutation, vars, &resp); err != nil {
		return CreateEnvironmentResult{}, err
	}
	return CreateEnvironmentResult{EnvironmentID: resp.CreateEnvironment.Environment.ID}, nil
}

// DestroyEnvironmentInput carries the environment identifier.
type DestroyEnvironmentInput struct {
	EnvironmentID string
}

// DestroyEnvironment removes an environment.
func (c *Client) DestroyEnvironment(ctx context.Context, in DestroyEnvironmentInput) error {
	// TODO: Replace with actual mutation and variables once confirmed.
	mutation := `mutation DeleteEnv($environmentId: ID!) {\n  deleteEnvironment(id: $environmentId) {\n    id\n  }\n}`
	vars := map[string]any{
		"environmentId": in.EnvironmentID,
	}
	var resp struct {
		DeleteEnvironment struct{ ID string } `json:"deleteEnvironment"`
	}
	return c.execute(ctx, mutation, vars, &resp)
}
