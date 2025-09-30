package railway

import "context"

// CreateServiceInput contains fields to create a service in an environment.
type CreateServiceInput struct {
	ProjectID     string
	EnvironmentID string
	Name          string
	Repo          *string
	Branch        *string
}

// CreateServiceResult captures the created service identifier.
type CreateServiceResult struct {
	ServiceID string
}

// CreateService executes the serviceCreate mutation.
func (c *Client) CreateService(ctx context.Context, in CreateServiceInput) (CreateServiceResult, error) {
	mutation := `mutation ServiceCreate($projectId: String!, $environmentId: String!, $name: String!, $repo: String, $branch: String) {
  serviceCreate(input: { branch: $branch, environmentId: $environmentId, source: { repo: $repo }, projectId: $projectId, name: $name }) {
    id
    name
    projectId
    updatedAt
  }
}`
	vars := map[string]any{
		"projectId":     in.ProjectID,
		"environmentId": in.EnvironmentID,
		"name":          in.Name,
		"repo":          in.Repo,
		"branch":        in.Branch,
	}
	var resp struct {
		ServiceCreate struct {
			ID string `json:"id"`
		} `json:"serviceCreate"`
	}
	if err := c.execute(ctx, mutation, vars, &resp); err != nil {
		return CreateServiceResult{}, err
	}
	return CreateServiceResult{ServiceID: resp.ServiceCreate.ID}, nil
}
