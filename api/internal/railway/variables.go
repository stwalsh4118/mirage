package railway

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/rs/zerolog/log"
)

// Embedded GraphQL query
var (
	//go:embed queries/queries/variables-get.graphql
	gqlGetVariables string
)

// GetEnvironmentVariablesInput contains parameters for fetching variables
type GetEnvironmentVariablesInput struct {
	ProjectID     string
	EnvironmentID string
	ServiceID     *string // nil for environment-level, set for service-level
}

// GetEnvironmentVariablesResult contains the fetched variables as a flat map
type GetEnvironmentVariablesResult struct {
	Variables map[string]string // Variable name -> value
}

// GetEnvironmentVariables fetches all variables for an environment or service.
// Returns a flat map of variable name -> value.
// Railway's variables query returns a simple JSON object, not relay-style connections.
func (c *Client) GetEnvironmentVariables(
	ctx context.Context,
	in GetEnvironmentVariablesInput,
) (GetEnvironmentVariablesResult, error) {
	query := gqlGetVariables

	vars := map[string]any{
		"projectId":     in.ProjectID,
		"environmentId": in.EnvironmentID,
	}

	// Add serviceId if requesting service-level variables
	if in.ServiceID != nil {
		vars["serviceId"] = *in.ServiceID
	} else {
		vars["serviceId"] = nil
	}

	log.Info().
		Str("project_id", in.ProjectID).
		Str("environment_id", in.EnvironmentID).
		Str("service_id", stringPtrToLogValue(in.ServiceID)).
		Msg("fetching environment variables from Railway")

	var resp struct {
		Variables map[string]string `json:"variables"`
	}

	if err := c.execute(ctx, query, vars, &resp); err != nil {
		return GetEnvironmentVariablesResult{}, fmt.Errorf("fetch environment variables: %w", err)
	}

	// Log count only, never log actual variable values for security
	log.Info().
		Str("project_id", in.ProjectID).
		Str("environment_id", in.EnvironmentID).
		Int("variable_count", len(resp.Variables)).
		Msg("successfully fetched environment variables")

	return GetEnvironmentVariablesResult{
		Variables: resp.Variables,
	}, nil
}

// ServiceVariables represents variables for a specific service
type ServiceVariables struct {
	ServiceID   string
	ServiceName string
	Variables   map[string]string
}

// GetAllEnvironmentAndServiceVariablesInput contains parameters for fetching all variables
type GetAllEnvironmentAndServiceVariablesInput struct {
	ProjectID     string
	EnvironmentID string
}

// GetAllEnvironmentAndServiceVariablesResult contains environment and all service variables
type GetAllEnvironmentAndServiceVariablesResult struct {
	EnvironmentVariables map[string]string  // Environment-level variables
	ServiceVariables     []ServiceVariables // Per-service variables
}

// GetAllEnvironmentAndServiceVariables fetches environment variables plus variables for ALL services
// in the environment. This is more efficient than calling GetEnvironmentVariables multiple times
// for cloning workflows.
func (c *Client) GetAllEnvironmentAndServiceVariables(
	ctx context.Context,
	in GetAllEnvironmentAndServiceVariablesInput,
) (GetAllEnvironmentAndServiceVariablesResult, error) {
	log.Info().
		Str("project_id", in.ProjectID).
		Str("environment_id", in.EnvironmentID).
		Msg("fetching all environment and service variables")

	// Step 1: Get project details to find all services in the environment
	projectDetails, err := c.GetProjectWithDetailsByID(ctx, in.ProjectID)
	if err != nil {
		return GetAllEnvironmentAndServiceVariablesResult{}, fmt.Errorf("fetch project details: %w", err)
	}

	// Find the target environment
	var targetEnv *ProjectEnvironment
	for i := range projectDetails.Environments {
		if projectDetails.Environments[i].ID == in.EnvironmentID {
			targetEnv = &projectDetails.Environments[i]
			break
		}
	}

	if targetEnv == nil {
		return GetAllEnvironmentAndServiceVariablesResult{}, fmt.Errorf("environment %s not found in project %s", in.EnvironmentID, in.ProjectID)
	}

	log.Info().
		Str("environment_id", in.EnvironmentID).
		Int("service_count", len(targetEnv.Services)).
		Msg("found services in environment")

	// Step 2: Fetch environment-level variables
	envVarsResult, err := c.GetEnvironmentVariables(ctx, GetEnvironmentVariablesInput{
		ProjectID:     in.ProjectID,
		EnvironmentID: in.EnvironmentID,
		ServiceID:     nil, // nil = environment-level
	})
	if err != nil {
		return GetAllEnvironmentAndServiceVariablesResult{}, fmt.Errorf("fetch environment variables: %w", err)
	}

	// Step 3: Fetch variables for each service
	serviceVars := make([]ServiceVariables, 0, len(targetEnv.Services))
	for _, svc := range targetEnv.Services {
		svcVarsResult, err := c.GetEnvironmentVariables(ctx, GetEnvironmentVariablesInput{
			ProjectID:     in.ProjectID,
			EnvironmentID: in.EnvironmentID,
			ServiceID:     &svc.ServiceID,
		})
		if err != nil {
			log.Warn().
				Str("service_id", svc.ServiceID).
				Str("service_name", svc.ServiceName).
				Err(err).
				Msg("failed to fetch service variables, skipping")
			continue // Skip services that fail, don't fail entire operation
		}

		serviceVars = append(serviceVars, ServiceVariables{
			ServiceID:   svc.ServiceID,
			ServiceName: svc.ServiceName,
			Variables:   svcVarsResult.Variables,
		})
	}

	log.Info().
		Str("environment_id", in.EnvironmentID).
		Int("env_variable_count", len(envVarsResult.Variables)).
		Int("services_with_variables", len(serviceVars)).
		Msg("successfully fetched all environment and service variables")

	return GetAllEnvironmentAndServiceVariablesResult{
		EnvironmentVariables: envVarsResult.Variables,
		ServiceVariables:     serviceVars,
	}, nil
}

// stringPtrToLogValue converts a string pointer to a log-safe value
func stringPtrToLogValue(s *string) string {
	if s == nil {
		return "(nil)"
	}
	return *s
}
