package railway

import (
	"context"
	_ "embed"
)

// RegistryCredentials holds authentication for private container registries.
type RegistryCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// CreateServiceInput contains fields to create a service in an environment.
// Supports both source repository and Docker image deployments.
type CreateServiceInput struct {
	ProjectID     string `json:"projectId"`
	EnvironmentID string `json:"environmentId"`
	Name          string

	// Source repository deployment (legacy/existing behavior)
	Repo   *string `json:"repo"`
	Branch *string `json:"branch"`

	// Docker image deployment (new)
	Image               *string              `json:"image"`               // e.g., "nginx:latest", "ghcr.io/owner/repo:v1.0"
	RegistryCredentials *RegistryCredentials `json:"registryCredentials"` // Optional, for private images

	// Service variables (e.g., RAILWAY_DOCKERFILE_PATH, custom env vars)
	Variables map[string]string `json:"variables,omitempty"`
}

// CreateServiceResult captures the created service identifier.
type CreateServiceResult struct {
	ServiceID string `json:"serviceId"`
}

// Embedded GraphQL mutations
var (
	//go:embed queries/mutations/service-create.graphql
	serviceCreateMutation string

	//go:embed queries/mutations/service-delete.graphql
	gqlServiceDelete string
)

// CreateService executes the serviceCreate mutation.
// Supports both source repository and Docker image deployments.
// Conditionally builds the input map structure to only include relevant fields.
func (c *Client) CreateService(ctx context.Context, in CreateServiceInput) (CreateServiceResult, error) {
	// Build the input structure conditionally based on deployment type
	input := map[string]any{
		"projectId":     in.ProjectID,
		"environmentId": in.EnvironmentID,
		"name":          in.Name,
	}

	if in.Image != nil {
		// Docker image deployment - only include image in source
		input["source"] = map[string]any{
			"image": *in.Image,
		}

		// Add credentials if provided
		if in.RegistryCredentials != nil {
			input["registryCredentials"] = map[string]any{
				"username": in.RegistryCredentials.Username,
				"password": in.RegistryCredentials.Password,
			}
		}
	} else {
		// Source repository deployment - only include repo in source
		input["source"] = map[string]any{
			"repo": in.Repo,
		}
		input["branch"] = in.Branch
	}

	// Add service variables if provided
	if len(in.Variables) > 0 {
		input["variables"] = in.Variables
	}

	vars := map[string]any{
		"input": input,
	}

	var resp struct {
		ServiceCreate struct {
			ID string `json:"id"`
		} `json:"serviceCreate"`
	}
	if err := c.execute(ctx, serviceCreateMutation, vars, &resp); err != nil {
		return CreateServiceResult{}, err
	}
	return CreateServiceResult{ServiceID: resp.ServiceCreate.ID}, nil
}

// DestroyServiceInput carries the service identifier.
type DestroyServiceInput struct {
	ServiceID string
}

// DestroyService removes a service from Railway.
func (c *Client) DestroyService(ctx context.Context, in DestroyServiceInput) error {
	mutation := gqlServiceDelete
	vars := map[string]any{
		"serviceId": in.ServiceID,
	}
	var resp struct {
		ServiceDelete bool `json:"serviceDelete"`
	}
	return c.execute(ctx, mutation, vars, &resp)
}
