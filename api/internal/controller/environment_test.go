package controller

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stwalsh4118/mirageapi/internal/railway"
	"github.com/stwalsh4118/mirageapi/internal/store"
)

// mockRailwayClientForEnvironment implements RailwayEnvironmentClient for testing
type mockRailwayClientForEnvironment struct {
	getEnvironmentVariablesFunc              func(ctx context.Context, in railway.GetEnvironmentVariablesInput) (railway.GetEnvironmentVariablesResult, error)
	getAllEnvironmentAndServiceVariablesFunc func(ctx context.Context, in railway.GetAllEnvironmentAndServiceVariablesInput) (railway.GetAllEnvironmentAndServiceVariablesResult, error)
	createEnvironmentFunc                    func(ctx context.Context, in railway.CreateEnvironmentInput) (railway.CreateEnvironmentResult, error)
}

func (m *mockRailwayClientForEnvironment) GetEnvironmentVariables(ctx context.Context, in railway.GetEnvironmentVariablesInput) (railway.GetEnvironmentVariablesResult, error) {
	if m.getEnvironmentVariablesFunc != nil {
		return m.getEnvironmentVariablesFunc(ctx, in)
	}
	return railway.GetEnvironmentVariablesResult{Variables: make(map[string]string)}, nil
}

func (m *mockRailwayClientForEnvironment) GetAllEnvironmentAndServiceVariables(ctx context.Context, in railway.GetAllEnvironmentAndServiceVariablesInput) (railway.GetAllEnvironmentAndServiceVariablesResult, error) {
	if m.getAllEnvironmentAndServiceVariablesFunc != nil {
		return m.getAllEnvironmentAndServiceVariablesFunc(ctx, in)
	}
	return railway.GetAllEnvironmentAndServiceVariablesResult{
		EnvironmentVariables: make(map[string]string),
		ServiceVariables:     []railway.ServiceVariables{},
	}, nil
}

func (m *mockRailwayClientForEnvironment) CreateEnvironment(ctx context.Context, in railway.CreateEnvironmentInput) (railway.CreateEnvironmentResult, error) {
	if m.createEnvironmentFunc != nil {
		return m.createEnvironmentFunc(ctx, in)
	}
	return railway.CreateEnvironmentResult{}, nil
}

// Stub methods to satisfy RailwayEnvironmentClient interface
func (m *mockRailwayClientForEnvironment) DestroyEnvironment(ctx context.Context, in railway.DestroyEnvironmentInput) error {
	return nil
}

func (m *mockRailwayClientForEnvironment) CreateProject(ctx context.Context, in railway.CreateProjectInput) (railway.CreateProjectResult, error) {
	return railway.CreateProjectResult{}, nil
}

func (m *mockRailwayClientForEnvironment) DestroyProject(ctx context.Context, in railway.DestroyProjectInput) error {
	return nil
}

func (m *mockRailwayClientForEnvironment) GetProject(ctx context.Context, id string) (railway.Project, error) {
	return railway.Project{}, nil
}

func (m *mockRailwayClientForEnvironment) GetProjectWithDetailsByID(ctx context.Context, id string) (railway.ProjectDetails, error) {
	return railway.ProjectDetails{}, nil
}

func (m *mockRailwayClientForEnvironment) ListProjects(ctx context.Context, limit int) ([]railway.Project, error) {
	return nil, nil
}

func (m *mockRailwayClientForEnvironment) ListProjectsWithDetails(ctx context.Context, limit int) ([]railway.ProjectDetails, error) {
	return nil, nil
}

func (m *mockRailwayClientForEnvironment) GetLatestDeploymentID(ctx context.Context, serviceID string) (string, error) {
	return "", nil
}

func (m *mockRailwayClientForEnvironment) GetDeploymentLogs(ctx context.Context, in railway.GetDeploymentLogsInput) (railway.GetDeploymentLogsResult, error) {
	return railway.GetDeploymentLogsResult{}, nil
}

func (m *mockRailwayClientForEnvironment) SubscribeToEnvironmentLogs(ctx context.Context, environmentID string, serviceFilter string) (*websocket.Conn, error) {
	return nil, nil
}

func (m *mockRailwayClientForEnvironment) SubscribeToDeploymentLogs(ctx context.Context, deploymentID string, filter string) (*websocket.Conn, error) {
	return nil, nil
}

func TestGetEnvironmentSnapshot_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test database
	db, err := store.Open(":memory:")
	require.NoError(t, err)

	// Create test environment
	env := store.Environment{
		ID:                   "mirage-env-123",
		Name:                 "test-env",
		Type:                 store.EnvironmentTypeDev,
		RailwayProjectID:     "railway-proj-123",
		RailwayEnvironmentID: "railway-env-123",
		SourceRepo:           "https://github.com/user/repo",
		SourceBranch:         "main",
		Status:               "active",
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
	require.NoError(t, db.Create(&env).Error)

	// Create test services
	service1 := store.Service{
		ID:               "service-1",
		EnvironmentID:    env.ID,
		Name:             "api",
		DeploymentType:   store.DeploymentTypeSourceRepo,
		SourceRepo:       "https://github.com/user/repo",
		SourceBranch:     "main",
		RailwayServiceID: "railway-service-1",
		Status:           "active",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	service2 := store.Service{
		ID:               "service-2",
		EnvironmentID:    env.ID,
		Name:             "worker",
		DeploymentType:   store.DeploymentTypeSourceRepo,
		SourceRepo:       "https://github.com/user/repo",
		SourceBranch:     "main",
		RailwayServiceID: "railway-service-2",
		Status:           "active",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	require.NoError(t, db.Create(&service1).Error)
	require.NoError(t, db.Create(&service2).Error)

	// Mock Railway client
	mockRailway := &mockRailwayClientForEnvironment{
		getAllEnvironmentAndServiceVariablesFunc: func(ctx context.Context, in railway.GetAllEnvironmentAndServiceVariablesInput) (railway.GetAllEnvironmentAndServiceVariablesResult, error) {
			assert.Equal(t, "railway-proj-123", in.ProjectID)
			assert.Equal(t, "railway-env-123", in.EnvironmentID)
			return railway.GetAllEnvironmentAndServiceVariablesResult{
				EnvironmentVariables: map[string]string{
					"NODE_ENV": "development",
					"API_KEY":  "test-key",
				},
				ServiceVariables: []railway.ServiceVariables{
					{
						ServiceID:   "railway-service-1",
						ServiceName: "api",
						Variables: map[string]string{
							"PORT": "8080",
						},
					},
					{
						ServiceID:   "railway-service-2",
						ServiceName: "worker",
						Variables: map[string]string{
							"QUEUE_URL": "redis://localhost",
						},
					},
				},
			}, nil
		},
	}

	controller := &EnvironmentController{
		DB:      db,
		Railway: mockRailway,
	}

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "railway-env-123"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/api/environments/railway-env-123/snapshot", nil)

	// Execute
	controller.GetEnvironmentSnapshot(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var snapshot EnvironmentSnapshot
	err = json.Unmarshal(w.Body.Bytes(), &snapshot)
	require.NoError(t, err)

	// Verify environment
	assert.Equal(t, env.ID, snapshot.Environment.ID)
	assert.Equal(t, env.Name, snapshot.Environment.Name)
	assert.Equal(t, env.RailwayProjectID, snapshot.Environment.RailwayProjectID)
	assert.Equal(t, env.RailwayEnvironmentID, snapshot.Environment.RailwayEnvironmentID)

	// Verify services
	assert.Len(t, snapshot.Services, 2)
	assert.Equal(t, "api", snapshot.Services[0].Name)
	assert.Equal(t, "worker", snapshot.Services[1].Name)

	// Verify environment variables
	assert.Len(t, snapshot.EnvironmentVariables, 2)
	assert.Equal(t, "development", snapshot.EnvironmentVariables["NODE_ENV"])
	assert.Equal(t, "test-key", snapshot.EnvironmentVariables["API_KEY"])

	// Verify service variables
	assert.Len(t, snapshot.ServiceVariables, 2)
	assert.Equal(t, "railway-service-1", snapshot.ServiceVariables[0].ServiceID)
	assert.Equal(t, "api", snapshot.ServiceVariables[0].ServiceName)
	assert.Equal(t, "8080", snapshot.ServiceVariables[0].Variables["PORT"])
	assert.Equal(t, "railway-service-2", snapshot.ServiceVariables[1].ServiceID)
	assert.Equal(t, "worker", snapshot.ServiceVariables[1].ServiceName)
	assert.Equal(t, "redis://localhost", snapshot.ServiceVariables[1].Variables["QUEUE_URL"])
}

func TestGetEnvironmentSnapshot_EnvironmentNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test database
	db, err := store.Open(":memory:")
	require.NoError(t, err)

	controller := &EnvironmentController{
		DB:      db,
		Railway: &mockRailwayClientForEnvironment{},
	}

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "nonexistent-env"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/api/environments/nonexistent-env/snapshot", nil)

	// Execute
	controller.GetEnvironmentSnapshot(c)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "environment not found")
}

func TestGetEnvironmentSnapshot_RailwayAPIFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test database
	db, err := store.Open(":memory:")
	require.NoError(t, err)

	// Create test environment
	env := store.Environment{
		ID:                   "mirage-env-123",
		Name:                 "test-env",
		Type:                 store.EnvironmentTypeDev,
		RailwayProjectID:     "railway-proj-123",
		RailwayEnvironmentID: "railway-env-123",
		Status:               "active",
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
	require.NoError(t, db.Create(&env).Error)

	// Mock Railway client that returns an error
	mockRailway := &mockRailwayClientForEnvironment{
		getAllEnvironmentAndServiceVariablesFunc: func(ctx context.Context, in railway.GetAllEnvironmentAndServiceVariablesInput) (railway.GetAllEnvironmentAndServiceVariablesResult, error) {
			return railway.GetAllEnvironmentAndServiceVariablesResult{}, errors.New("railway api error")
		},
	}

	controller := &EnvironmentController{
		DB:      db,
		Railway: mockRailway,
	}

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "railway-env-123"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/api/environments/railway-env-123/snapshot", nil)

	// Execute
	controller.GetEnvironmentSnapshot(c)

	// Assert - should succeed but with empty variables
	assert.Equal(t, http.StatusOK, w.Code)

	var snapshot EnvironmentSnapshot
	err = json.Unmarshal(w.Body.Bytes(), &snapshot)
	require.NoError(t, err)

	// Verify environment is still returned
	assert.Equal(t, env.ID, snapshot.Environment.ID)

	// Verify variables are empty (not nil)
	assert.NotNil(t, snapshot.EnvironmentVariables)
	assert.Len(t, snapshot.EnvironmentVariables, 0)
	assert.NotNil(t, snapshot.ServiceVariables)
	assert.Len(t, snapshot.ServiceVariables, 0)
}

func TestGetEnvironmentSnapshot_NoServices(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test database
	db, err := store.Open(":memory:")
	require.NoError(t, err)

	// Create test environment without services
	env := store.Environment{
		ID:                   "mirage-env-123",
		Name:                 "test-env",
		Type:                 store.EnvironmentTypeDev,
		RailwayProjectID:     "railway-proj-123",
		RailwayEnvironmentID: "railway-env-123",
		Status:               "active",
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
	require.NoError(t, db.Create(&env).Error)

	// Mock Railway client
	mockRailway := &mockRailwayClientForEnvironment{
		getAllEnvironmentAndServiceVariablesFunc: func(ctx context.Context, in railway.GetAllEnvironmentAndServiceVariablesInput) (railway.GetAllEnvironmentAndServiceVariablesResult, error) {
			return railway.GetAllEnvironmentAndServiceVariablesResult{
				EnvironmentVariables: map[string]string{"TEST": "value"},
				ServiceVariables:     []railway.ServiceVariables{},
			}, nil
		},
	}

	controller := &EnvironmentController{
		DB:      db,
		Railway: mockRailway,
	}

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "railway-env-123"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/api/environments/railway-env-123/snapshot", nil)

	// Execute
	controller.GetEnvironmentSnapshot(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var snapshot EnvironmentSnapshot
	err = json.Unmarshal(w.Body.Bytes(), &snapshot)
	require.NoError(t, err)

	// Verify environment
	assert.Equal(t, env.ID, snapshot.Environment.ID)

	// Verify services array is not nil but empty
	assert.NotNil(t, snapshot.Services)
	assert.Len(t, snapshot.Services, 0)

	// Verify variables are still returned
	assert.Len(t, snapshot.EnvironmentVariables, 1)
	assert.Equal(t, "value", snapshot.EnvironmentVariables["TEST"])
	assert.Len(t, snapshot.ServiceVariables, 0)
}

func TestGetEnvironmentSnapshot_NilRailwayClient(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test database
	db, err := store.Open(":memory:")
	require.NoError(t, err)

	// Create test environment
	env := store.Environment{
		ID:                   "mirage-env-123",
		Name:                 "test-env",
		Type:                 store.EnvironmentTypeDev,
		RailwayProjectID:     "railway-proj-123",
		RailwayEnvironmentID: "railway-env-123",
		Status:               "active",
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
	require.NoError(t, db.Create(&env).Error)

	// Controller with nil Railway client
	controller := &EnvironmentController{
		DB:      db,
		Railway: nil,
	}

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "railway-env-123"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/api/environments/railway-env-123/snapshot", nil)

	// Execute
	controller.GetEnvironmentSnapshot(c)

	// Assert - should succeed but with empty variables
	assert.Equal(t, http.StatusOK, w.Code)

	var snapshot EnvironmentSnapshot
	err = json.Unmarshal(w.Body.Bytes(), &snapshot)
	require.NoError(t, err)

	// Verify environment is returned
	assert.Equal(t, env.ID, snapshot.Environment.ID)

	// Verify variables are empty
	assert.NotNil(t, snapshot.EnvironmentVariables)
	assert.Len(t, snapshot.EnvironmentVariables, 0)
	assert.NotNil(t, snapshot.ServiceVariables)
	assert.Len(t, snapshot.ServiceVariables, 0)
}

func TestGetEnvironmentSnapshot_MissingIDParam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := store.Open(":memory:")
	require.NoError(t, err)

	controller := &EnvironmentController{
		DB:      db,
		Railway: &mockRailwayClientForEnvironment{},
	}

	// Create test request without ID param
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{} // No ID param
	c.Request = httptest.NewRequest(http.MethodGet, "/api/environments//snapshot", nil)

	// Execute
	controller.GetEnvironmentSnapshot(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "railway environment id required")
}

func TestGetEnvironmentSnapshot_DatabaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a database and close it to simulate errors
	db, err := store.Open(":memory:")
	require.NoError(t, err)

	// Get the underlying SQL DB and close it
	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.Close()

	controller := &EnvironmentController{
		DB:      db,
		Railway: &mockRailwayClientForEnvironment{},
	}

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "railway-env-123"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/api/environments/railway-env-123/snapshot", nil)

	// Execute
	controller.GetEnvironmentSnapshot(c)

	// Assert - should return internal server error
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetEnvironmentSnapshot_EmptyVariables(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test database
	db, err := store.Open(":memory:")
	require.NoError(t, err)

	// Create test environment
	env := store.Environment{
		ID:                   "mirage-env-123",
		Name:                 "test-env",
		Type:                 store.EnvironmentTypeDev,
		RailwayProjectID:     "railway-proj-123",
		RailwayEnvironmentID: "railway-env-123",
		Status:               "active",
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
	require.NoError(t, db.Create(&env).Error)

	// Mock Railway client that returns empty variables
	mockRailway := &mockRailwayClientForEnvironment{
		getAllEnvironmentAndServiceVariablesFunc: func(ctx context.Context, in railway.GetAllEnvironmentAndServiceVariablesInput) (railway.GetAllEnvironmentAndServiceVariablesResult, error) {
			return railway.GetAllEnvironmentAndServiceVariablesResult{
				EnvironmentVariables: make(map[string]string), // Empty map
				ServiceVariables:     []railway.ServiceVariables{},
			}, nil
		},
	}

	controller := &EnvironmentController{
		DB:      db,
		Railway: mockRailway,
	}

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "railway-env-123"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/api/environments/railway-env-123/snapshot", nil)

	// Execute
	controller.GetEnvironmentSnapshot(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var snapshot EnvironmentSnapshot
	err = json.Unmarshal(w.Body.Bytes(), &snapshot)
	require.NoError(t, err)

	// Verify environment is returned
	assert.Equal(t, env.ID, snapshot.Environment.ID)

	// Verify variables are empty but not nil
	assert.NotNil(t, snapshot.EnvironmentVariables)
	assert.Len(t, snapshot.EnvironmentVariables, 0)
	assert.NotNil(t, snapshot.ServiceVariables)
	assert.Len(t, snapshot.ServiceVariables, 0)
}
