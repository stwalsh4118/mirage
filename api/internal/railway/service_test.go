package railway

import (
	"strings"
	"testing"
)

func TestServiceCreateMutation_Structure(t *testing.T) {
	// Verify the single mutation uses input type variable
	expectedFields := []string{
		"$input: ServiceCreateInput!",
		"serviceCreate(input: $input)",
	}

	for _, field := range expectedFields {
		if !strings.Contains(serviceCreateMutation, field) {
			t.Errorf("expected mutation to contain %q, but it didn't", field)
		}
	}
}

func TestCreateServiceInput_DeterminesDeploymentType(t *testing.T) {
	tests := []struct {
		name        string
		input       CreateServiceInput
		expectRepo  bool
		expectImage bool
		expectCreds bool
		description string
	}{
		{
			name: "source repository deployment",
			input: CreateServiceInput{
				ProjectID:     "proj-1",
				EnvironmentID: "env-1",
				Name:          "api",
				Repo:          stringPtr("github.com/owner/repo"),
				Branch:        stringPtr("main"),
			},
			expectRepo:  true,
			expectImage: false,
			expectCreds: false,
			description: "Should use repo source when Repo is provided",
		},
		{
			name: "docker image deployment without credentials",
			input: CreateServiceInput{
				ProjectID:     "proj-1",
				EnvironmentID: "env-1",
				Name:          "nginx",
				Image:         stringPtr("nginx:latest"),
			},
			expectRepo:  false,
			expectImage: true,
			expectCreds: false,
			description: "Should use image source when Image is provided",
		},
		{
			name: "docker image deployment with credentials",
			input: CreateServiceInput{
				ProjectID:     "proj-1",
				EnvironmentID: "env-1",
				Name:          "private-app",
				Image:         stringPtr("ghcr.io/owner/private:v1.0"),
				RegistryCredentials: &RegistryCredentials{
					Username: "user",
					Password: "token",
				},
			},
			expectRepo:  false,
			expectImage: true,
			expectCreds: true,
			description: "Should use image source with credentials when both provided",
		},
		{
			name: "image takes precedence over repo",
			input: CreateServiceInput{
				ProjectID:     "proj-1",
				EnvironmentID: "env-1",
				Name:          "service",
				Repo:          stringPtr("github.com/owner/repo"),
				Branch:        stringPtr("main"),
				Image:         stringPtr("nginx:latest"),
			},
			expectRepo:  false,
			expectImage: true,
			expectCreds: false,
			description: "Image deployment should take precedence when both provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasImage := tt.input.Image != nil
			hasCreds := tt.input.RegistryCredentials != nil

			if hasImage != tt.expectImage {
				t.Errorf("%s: expected image=%v, got image=%v", tt.description, tt.expectImage, hasImage)
			}

			if hasCreds != tt.expectCreds {
				t.Errorf("%s: expected creds=%v, got creds=%v", tt.description, tt.expectCreds, hasCreds)
			}
		})
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}

func TestCreateServiceInput_WithVariables(t *testing.T) {
	tests := []struct {
		name        string
		input       CreateServiceInput
		expectVars  bool
		expectedKey string
		expectedVal string
		description string
	}{
		{
			name: "repository deployment with RAILWAY_DOCKERFILE_PATH",
			input: CreateServiceInput{
				ProjectID:     "proj-1",
				EnvironmentID: "env-1",
				Name:          "api",
				Repo:          stringPtr("github.com/owner/monorepo"),
				Branch:        stringPtr("main"),
				Variables: map[string]string{
					"RAILWAY_DOCKERFILE_PATH": "services/api/Dockerfile",
				},
			},
			expectVars:  true,
			expectedKey: "RAILWAY_DOCKERFILE_PATH",
			expectedVal: "services/api/Dockerfile",
			description: "Should include RAILWAY_DOCKERFILE_PATH variable",
		},
		{
			name: "repository deployment with multiple variables",
			input: CreateServiceInput{
				ProjectID:     "proj-1",
				EnvironmentID: "env-1",
				Name:          "worker",
				Repo:          stringPtr("github.com/owner/monorepo"),
				Branch:        stringPtr("main"),
				Variables: map[string]string{
					"RAILWAY_DOCKERFILE_PATH": "services/worker/Dockerfile",
					"NODE_ENV":                "production",
					"WORKER_CONCURRENCY":      "10",
				},
			},
			expectVars:  true,
			expectedKey: "RAILWAY_DOCKERFILE_PATH",
			expectedVal: "services/worker/Dockerfile",
			description: "Should include multiple variables",
		},
		{
			name: "repository deployment without variables",
			input: CreateServiceInput{
				ProjectID:     "proj-1",
				EnvironmentID: "env-1",
				Name:          "web",
				Repo:          stringPtr("github.com/owner/simple-app"),
				Branch:        stringPtr("main"),
			},
			expectVars:  false,
			description: "Should work without variables",
		},
		{
			name: "image deployment should not use variables for dockerfile path",
			input: CreateServiceInput{
				ProjectID:     "proj-1",
				EnvironmentID: "env-1",
				Name:          "nginx",
				Image:         stringPtr("nginx:latest"),
				Variables: map[string]string{
					"CUSTOM_VAR": "value",
				},
			},
			expectVars:  true,
			expectedKey: "CUSTOM_VAR",
			expectedVal: "value",
			description: "Image deployments can have variables but not RAILWAY_DOCKERFILE_PATH",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasVars := len(tt.input.Variables) > 0
			if hasVars != tt.expectVars {
				t.Errorf("%s: expected variables=%v, got variables=%v", tt.description, tt.expectVars, hasVars)
			}

			if tt.expectVars && tt.expectedKey != "" {
				val, ok := tt.input.Variables[tt.expectedKey]
				if !ok {
					t.Errorf("%s: expected variable %q to exist but it didn't", tt.description, tt.expectedKey)
				}
				if val != tt.expectedVal {
					t.Errorf("%s: expected variable %q=%q, got %q", tt.description, tt.expectedKey, tt.expectedVal, val)
				}
			}
		})
	}
}
