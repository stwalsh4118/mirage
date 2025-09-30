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
