package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stwalsh4118/mirageapi/internal/railway"
)

// mockRailwayClient implements RailwayServiceClient for testing
type mockRailwayClient struct {
	createServiceFunc  func(ctx context.Context, in railway.CreateServiceInput) (railway.CreateServiceResult, error)
	destroyServiceFunc func(ctx context.Context, in railway.DestroyServiceInput) error
}

func (m *mockRailwayClient) CreateService(ctx context.Context, in railway.CreateServiceInput) (railway.CreateServiceResult, error) {
	if m.createServiceFunc != nil {
		return m.createServiceFunc(ctx, in)
	}
	return railway.CreateServiceResult{ServiceID: "test-service-id"}, nil
}

func (m *mockRailwayClient) DestroyService(ctx context.Context, in railway.DestroyServiceInput) error {
	if m.destroyServiceFunc != nil {
		return m.destroyServiceFunc(ctx, in)
	}
	return nil
}

func TestProvisionServices_RepoBasedDeployment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := "https://github.com/user/repo"
	branch := "main"

	mockClient := &mockRailwayClient{
		createServiceFunc: func(ctx context.Context, in railway.CreateServiceInput) (railway.CreateServiceResult, error) {
			// Verify repo-based fields are set
			assert.NotNil(t, in.Repo)
			assert.Equal(t, repo, *in.Repo)
			assert.NotNil(t, in.Branch)
			assert.Equal(t, branch, *in.Branch)
			assert.Nil(t, in.Image)
			return railway.CreateServiceResult{ServiceID: "service-123"}, nil
		},
	}

	controller := &ServicesController{Railway: mockClient}

	reqBody := ProvisionServicesRequest{
		ProjectID:     "proj-123",
		EnvironmentID: "env-123",
		Services: []ServiceSpec{
			{
				Name:   "test-service",
				Repo:   &repo,
				Branch: &branch,
			},
		},
		RequestID: "req-123",
	}

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/provision/services", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.ProvisionServices(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp ProvisionServicesResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, []string{"service-123"}, resp.ServiceIDs)
}

func TestProvisionServices_ImageBasedDeployment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	imageName := "nginx"
	imageTag := "alpine"

	mockClient := &mockRailwayClient{
		createServiceFunc: func(ctx context.Context, in railway.CreateServiceInput) (railway.CreateServiceResult, error) {
			// Verify image-based fields are set
			assert.Nil(t, in.Repo)
			assert.Nil(t, in.Branch)
			assert.NotNil(t, in.Image)
			assert.Equal(t, "nginx:alpine", *in.Image)
			return railway.CreateServiceResult{ServiceID: "service-456"}, nil
		},
	}

	controller := &ServicesController{Railway: mockClient}

	reqBody := ProvisionServicesRequest{
		ProjectID:     "proj-123",
		EnvironmentID: "env-123",
		Services: []ServiceSpec{
			{
				Name:      "nginx-service",
				ImageName: &imageName,
				ImageTag:  &imageTag,
			},
		},
		RequestID: "req-456",
	}

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/provision/services", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.ProvisionServices(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp ProvisionServicesResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, []string{"service-456"}, resp.ServiceIDs)
}

func TestProvisionServices_ImageWithRegistry(t *testing.T) {
	gin.SetMode(gin.TestMode)

	registry := "ghcr.io"
	imageName := "user/app"
	imageTag := "v1.0.0"

	mockClient := &mockRailwayClient{
		createServiceFunc: func(ctx context.Context, in railway.CreateServiceInput) (railway.CreateServiceResult, error) {
			assert.NotNil(t, in.Image)
			assert.Equal(t, "ghcr.io/user/app:v1.0.0", *in.Image)
			return railway.CreateServiceResult{ServiceID: "service-789"}, nil
		},
	}

	controller := &ServicesController{Railway: mockClient}

	reqBody := ProvisionServicesRequest{
		ProjectID:     "proj-123",
		EnvironmentID: "env-123",
		Services: []ServiceSpec{
			{
				Name:          "custom-service",
				ImageRegistry: &registry,
				ImageName:     &imageName,
				ImageTag:      &imageTag,
			},
		},
		RequestID: "req-789",
	}

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/provision/services", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.ProvisionServices(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp ProvisionServicesResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, []string{"service-789"}, resp.ServiceIDs)
}

func TestProvisionServices_ImageWithCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)

	imageName := "private/app"
	username := "user"
	password := "pass"

	mockClient := &mockRailwayClient{
		createServiceFunc: func(ctx context.Context, in railway.CreateServiceInput) (railway.CreateServiceResult, error) {
			assert.NotNil(t, in.Image)
			assert.NotNil(t, in.RegistryCredentials)
			assert.Equal(t, username, in.RegistryCredentials.Username)
			assert.Equal(t, password, in.RegistryCredentials.Password)
			return railway.CreateServiceResult{ServiceID: "service-private"}, nil
		},
	}

	controller := &ServicesController{Railway: mockClient}

	reqBody := ProvisionServicesRequest{
		ProjectID:     "proj-123",
		EnvironmentID: "env-123",
		Services: []ServiceSpec{
			{
				Name:             "private-service",
				ImageName:        &imageName,
				RegistryUsername: &username,
				RegistryPassword: &password,
			},
		},
		RequestID: "req-private",
	}

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/provision/services", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.ProvisionServices(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProvisionServices_ValidationBothRepoAndImage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := "https://github.com/user/repo"
	branch := "main"
	imageName := "nginx"

	controller := &ServicesController{Railway: &mockRailwayClient{}}

	reqBody := ProvisionServicesRequest{
		ProjectID:     "proj-123",
		EnvironmentID: "env-123",
		Services: []ServiceSpec{
			{
				Name:      "invalid-service",
				Repo:      &repo,
				Branch:    &branch,
				ImageName: &imageName,
			},
		},
		RequestID: "req-invalid",
	}

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/provision/services", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.ProvisionServices(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "cannot specify both")
}

func TestProvisionServices_ValidationNeither(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := &ServicesController{Railway: &mockRailwayClient{}}

	reqBody := ProvisionServicesRequest{
		ProjectID:     "proj-123",
		EnvironmentID: "env-123",
		Services: []ServiceSpec{
			{
				Name: "invalid-service",
			},
		},
		RequestID: "req-invalid",
	}

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/provision/services", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.ProvisionServices(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "must specify either")
}

func TestProvisionServices_ValidationMissingBranch(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := "https://github.com/user/repo"

	controller := &ServicesController{Railway: &mockRailwayClient{}}

	reqBody := ProvisionServicesRequest{
		ProjectID:     "proj-123",
		EnvironmentID: "env-123",
		Services: []ServiceSpec{
			{
				Name: "invalid-service",
				Repo: &repo,
			},
		},
		RequestID: "req-invalid",
	}

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/provision/services", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.ProvisionServices(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "branch is required")
}

func TestBuildImageReference_DefaultTag(t *testing.T) {
	imageName := "nginx"
	spec := ServiceSpec{
		ImageName: &imageName,
	}

	result := buildImageReference(spec)
	assert.Equal(t, "nginx:latest", result)
}

func TestBuildImageReference_WithTag(t *testing.T) {
	imageName := "nginx"
	tag := "alpine"
	spec := ServiceSpec{
		ImageName: &imageName,
		ImageTag:  &tag,
	}

	result := buildImageReference(spec)
	assert.Equal(t, "nginx:alpine", result)
}

func TestBuildImageReference_WithRegistry(t *testing.T) {
	registry := "ghcr.io"
	imageName := "user/app"
	tag := "v1.0"
	spec := ServiceSpec{
		ImageRegistry: &registry,
		ImageName:     &imageName,
		ImageTag:      &tag,
	}

	result := buildImageReference(spec)
	assert.Equal(t, "ghcr.io/user/app:v1.0", result)
}

func TestBuildImageReference_FullReference(t *testing.T) {
	registry := "docker.io"
	imageName := "library/postgres"
	tag := "15-alpine"
	spec := ServiceSpec{
		ImageRegistry: &registry,
		ImageName:     &imageName,
		ImageTag:      &tag,
	}

	result := buildImageReference(spec)
	assert.Equal(t, "docker.io/library/postgres:15-alpine", result)
}

func TestValidateDockerfilePath_ValidPaths(t *testing.T) {
	validPaths := []string{
		"Dockerfile",
		"services/api/Dockerfile",
		"apps/backend/prod.dockerfile",
		"packages/auth/Dockerfile",
		"path/to/nested/service/Dockerfile",
	}

	for _, path := range validPaths {
		t.Run(path, func(t *testing.T) {
			err := validateDockerfilePath(path)
			assert.NoError(t, err, "path %q should be valid", path)
		})
	}
}

func TestValidateDockerfilePath_InvalidPaths(t *testing.T) {
	tests := []struct {
		path        string
		description string
	}{
		{"", "empty path"},
		{"/services/api/Dockerfile", "absolute path (Unix)"},
		{"C:\\apps\\Dockerfile", "absolute path (Windows)"},
		{"../Dockerfile", "parent directory traversal"},
		{"../../../Dockerfile", "multiple parent traversal"},
		{"services/../../../etc/passwd", "path with traversal in middle"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := validateDockerfilePath(tt.path)
			assert.Error(t, err, "path %q should be invalid: %s", tt.path, tt.description)
		})
	}
}

func TestValidateDockerfilePath_MaxLength(t *testing.T) {
	// Create a path longer than 512 characters
	longPath := "services/"
	for len(longPath) < 520 {
		longPath += "very/long/path/"
	}
	longPath += "Dockerfile"

	err := validateDockerfilePath(longPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds maximum length")
}

func TestProvisionServices_WithDockerfilePath(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := "https://github.com/user/monorepo"
	branch := "main"
	dockerfilePath := "services/api/Dockerfile"

	mockClient := &mockRailwayClient{
		createServiceFunc: func(ctx context.Context, in railway.CreateServiceInput) (railway.CreateServiceResult, error) {
			// Verify dockerfile path is set as variable
			assert.NotNil(t, in.Repo)
			assert.NotNil(t, in.Variables)
			assert.Equal(t, dockerfilePath, in.Variables["RAILWAY_DOCKERFILE_PATH"])
			return railway.CreateServiceResult{ServiceID: "service-monorepo"}, nil
		},
	}

	controller := &ServicesController{Railway: mockClient}

	reqBody := ProvisionServicesRequest{
		ProjectID:     "proj-123",
		EnvironmentID: "env-123",
		Services: []ServiceSpec{
			{
				Name:           "api-service",
				Repo:           &repo,
				Branch:         &branch,
				DockerfilePath: &dockerfilePath,
			},
		},
		RequestID: "req-monorepo",
	}

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/provision/services", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.ProvisionServices(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp ProvisionServicesResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, []string{"service-monorepo"}, resp.ServiceIDs)
}

func TestProvisionServices_DockerfilePathOnlyForRepoDeployment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	imageName := "nginx"
	dockerfilePath := "services/api/Dockerfile"

	controller := &ServicesController{Railway: &mockRailwayClient{}}

	reqBody := ProvisionServicesRequest{
		ProjectID:     "proj-123",
		EnvironmentID: "env-123",
		Services: []ServiceSpec{
			{
				Name:           "invalid-service",
				ImageName:      &imageName,
				DockerfilePath: &dockerfilePath,
			},
		},
		RequestID: "req-invalid",
	}

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/provision/services", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.ProvisionServices(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "dockerfilePath can only be specified for repository-based deployments")
}

func TestProvisionServices_InvalidDockerfilePath(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := "https://github.com/user/monorepo"
	branch := "main"
	dockerfilePath := "../../../etc/passwd"

	controller := &ServicesController{Railway: &mockRailwayClient{}}

	reqBody := ProvisionServicesRequest{
		ProjectID:     "proj-123",
		EnvironmentID: "env-123",
		Services: []ServiceSpec{
			{
				Name:           "invalid-service",
				Repo:           &repo,
				Branch:         &branch,
				DockerfilePath: &dockerfilePath,
			},
		},
		RequestID: "req-invalid",
	}

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/provision/services", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.ProvisionServices(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "cannot traverse parent directories")
}

func TestProvisionServices_WithEnvVars(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		envVars        map[string]string
		dockerfilePath *string
		expectedVars   map[string]string
		expectVarsSet  bool
		description    string
	}{
		{
			name:          "env vars only",
			envVars:       map[string]string{"PORT": "8080", "NODE_ENV": "production"},
			expectedVars:  map[string]string{"PORT": "8080", "NODE_ENV": "production"},
			expectVarsSet: true,
			description:   "Should set user environment variables",
		},
		{
			name:           "dockerfile path only",
			dockerfilePath: ptrString("services/api/Dockerfile"),
			expectedVars:   map[string]string{"RAILWAY_DOCKERFILE_PATH": "services/api/Dockerfile"},
			expectVarsSet:  true,
			description:    "Should set RAILWAY_DOCKERFILE_PATH system variable",
		},
		{
			name:           "both env vars and dockerfile path",
			envVars:        map[string]string{"PORT": "8080", "LOG_LEVEL": "info"},
			dockerfilePath: ptrString("services/api/Dockerfile"),
			expectedVars: map[string]string{
				"PORT":                    "8080",
				"LOG_LEVEL":               "info",
				"RAILWAY_DOCKERFILE_PATH": "services/api/Dockerfile",
			},
			expectVarsSet: true,
			description:   "Should merge user vars and system vars",
		},
		{
			name:           "system var overrides user var",
			envVars:        map[string]string{"RAILWAY_DOCKERFILE_PATH": "wrong/path", "PORT": "8080"},
			dockerfilePath: ptrString("correct/path/Dockerfile"),
			expectedVars: map[string]string{
				"RAILWAY_DOCKERFILE_PATH": "correct/path/Dockerfile", // System var wins
				"PORT":                    "8080",
			},
			expectVarsSet: true,
			description:   "System RAILWAY_DOCKERFILE_PATH should override user-specified value",
		},
		{
			name:          "no variables",
			expectVarsSet: false,
			description:   "Should work with no variables set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := "https://github.com/user/repo"
			branch := "main"

			mockClient := &mockRailwayClient{
				createServiceFunc: func(ctx context.Context, in railway.CreateServiceInput) (railway.CreateServiceResult, error) {
					// Verify variables are set correctly
					if tt.expectVarsSet {
						assert.NotNil(t, in.Variables, tt.description)
						assert.Equal(t, len(tt.expectedVars), len(in.Variables), "variable count should match for %s", tt.description)
						for k, v := range tt.expectedVars {
							actualValue, exists := in.Variables[k]
							assert.True(t, exists, "variable %q should exist for %s", k, tt.description)
							assert.Equal(t, v, actualValue, "variable %q value mismatch for %s", k, tt.description)
						}
					} else {
						// Either nil or empty map
						assert.True(t, in.Variables == nil || len(in.Variables) == 0, "variables should be empty for %s", tt.description)
					}
					return railway.CreateServiceResult{ServiceID: "service-123"}, nil
				},
			}

			controller := &ServicesController{Railway: mockClient}

			reqBody := ProvisionServicesRequest{
				ProjectID:     "proj-123",
				EnvironmentID: "env-123",
				Services: []ServiceSpec{
					{
						Name:           "test-service",
						Repo:           &repo,
						Branch:         &branch,
						EnvVars:        tt.envVars,
						DockerfilePath: tt.dockerfilePath,
					},
				},
				RequestID: "req-test",
			}

			body, _ := json.Marshal(reqBody)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/provision/services", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")

			controller.ProvisionServices(c)

			assert.Equal(t, http.StatusOK, w.Code, "should return 200 OK for %s", tt.description)
			var resp ProvisionServicesResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)
			assert.Equal(t, []string{"service-123"}, resp.ServiceIDs)
		})
	}
}

func TestProvisionServices_MultipleServicesWithDifferentEnvVars(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := "https://github.com/user/monorepo"
	branch := "main"

	// Track which service is being created
	serviceIndex := 0
	expectedServices := []struct {
		name         string
		expectedVars map[string]string
	}{
		{
			name: "api-service",
			expectedVars: map[string]string{
				"PORT":                    "8080",
				"NODE_ENV":                "production",
				"RAILWAY_DOCKERFILE_PATH": "services/api/Dockerfile",
			},
		},
		{
			name: "worker-service",
			expectedVars: map[string]string{
				"PORT":                    "8081",
				"QUEUE_NAME":              "jobs",
				"RAILWAY_DOCKERFILE_PATH": "services/worker/Dockerfile",
			},
		},
		{
			name: "web-service",
			expectedVars: map[string]string{
				"PORT":                    "3000",
				"RAILWAY_DOCKERFILE_PATH": "services/web/Dockerfile",
			},
		},
	}

	mockClient := &mockRailwayClient{
		createServiceFunc: func(ctx context.Context, in railway.CreateServiceInput) (railway.CreateServiceResult, error) {
			expected := expectedServices[serviceIndex]
			assert.Equal(t, expected.name, in.Name, "service name should match")

			assert.NotNil(t, in.Variables, "variables should be set for %s", expected.name)
			for k, v := range expected.expectedVars {
				actualValue, exists := in.Variables[k]
				assert.True(t, exists, "variable %q should exist for service %s", k, expected.name)
				assert.Equal(t, v, actualValue, "variable %q value mismatch for service %s", k, expected.name)
			}

			serviceIndex++
			return railway.CreateServiceResult{ServiceID: "service-" + expected.name}, nil
		},
	}

	controller := &ServicesController{Railway: mockClient}

	reqBody := ProvisionServicesRequest{
		ProjectID:     "proj-123",
		EnvironmentID: "env-123",
		Services: []ServiceSpec{
			{
				Name:           "api-service",
				Repo:           &repo,
				Branch:         &branch,
				EnvVars:        map[string]string{"PORT": "8080", "NODE_ENV": "production"},
				DockerfilePath: ptrString("services/api/Dockerfile"),
			},
			{
				Name:           "worker-service",
				Repo:           &repo,
				Branch:         &branch,
				EnvVars:        map[string]string{"PORT": "8081", "QUEUE_NAME": "jobs"},
				DockerfilePath: ptrString("services/worker/Dockerfile"),
			},
			{
				Name:           "web-service",
				Repo:           &repo,
				Branch:         &branch,
				EnvVars:        map[string]string{"PORT": "3000"},
				DockerfilePath: ptrString("services/web/Dockerfile"),
			},
		},
		RequestID: "req-multi",
	}

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/provision/services", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.ProvisionServices(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp ProvisionServicesResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 3, len(resp.ServiceIDs), "should create 3 services")
}

// Helper function to create string pointer
func ptrString(s string) *string {
	return &s
}
