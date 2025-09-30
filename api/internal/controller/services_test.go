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
	createServiceFunc func(ctx context.Context, in railway.CreateServiceInput) (railway.CreateServiceResult, error)
}

func (m *mockRailwayClient) CreateService(ctx context.Context, in railway.CreateServiceInput) (railway.CreateServiceResult, error) {
	if m.createServiceFunc != nil {
		return m.createServiceFunc(ctx, in)
	}
	return railway.CreateServiceResult{ServiceID: "test-service-id"}, nil
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
