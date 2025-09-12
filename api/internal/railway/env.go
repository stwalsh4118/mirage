package railway

import (
	"context"

	"github.com/rs/zerolog/log"
)

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
	mutation := `mutation EnvironmentCreate($projectId: String!, $name: String!) {
    environmentCreate(input: { projectId: $projectId, name: $name }) {
        createdAt
        id
        name
        projectId
        updatedAt
    }
}`

	vars := map[string]any{
		"projectId": in.ProjectID,
		"name":      in.Name,
	}

	log.Info().Str("project_id", in.ProjectID).Str("name", in.Name).Msg("creating environment")
	var resp struct {
		EnvironmentCreate struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			ProjectID string `json:"projectId"`
		} `json:"environmentCreate"`
	}
	if err := c.execute(ctx, mutation, vars, &resp); err != nil {
		return CreateEnvironmentResult{}, err
	}
	return CreateEnvironmentResult{EnvironmentID: resp.EnvironmentCreate.ID}, nil
}

// DestroyEnvironmentInput carries the environment identifier.
type DestroyEnvironmentInput struct {
	EnvironmentID string
}

// DestroyEnvironment removes an environment.
func (c *Client) DestroyEnvironment(ctx context.Context, in DestroyEnvironmentInput) error {
	// TODO: Replace with actual mutation and variables once confirmed.
	mutation := `mutation DeleteEnv($environmentId: ID!) {
  deleteEnvironment(id: $environmentId) {
    id
  }
}`
	vars := map[string]any{
		"environmentId": in.EnvironmentID,
	}
	var resp struct {
		DeleteEnvironment struct{ ID string } `json:"deleteEnvironment"`
	}
	return c.execute(ctx, mutation, vars, &resp)
}

// GetEnvironmentStatus fetches the current status string for a Railway environment.
func (c *Client) GetEnvironmentStatus(ctx context.Context, environmentID string) (string, error) {
	// NOTE: This query is a placeholder and may need adjustment to match Railway's schema.
	query := `query EnvStatus($environmentId: ID!) {
  environment(id: $environmentId) {
    id
    status
  }
}`
	vars := map[string]any{
		"environmentId": environmentID,
	}
	var resp struct {
		Environment struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		} `json:"environment"`
	}
	if err := c.execute(ctx, query, vars, &resp); err != nil {
		return "", err
	}
	return resp.Environment.Status, nil
}
