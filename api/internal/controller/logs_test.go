package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stwalsh4118/mirageapi/internal/railway"
	"github.com/stwalsh4118/mirageapi/internal/store"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// MockRailwayClient implements railway.Client for testing
type MockRailwayClient struct {
	GetDeploymentLogsFunc     func(ctx context.Context, input railway.GetDeploymentLogsInput) (railway.GetDeploymentLogsResult, error)
	GetLatestDeploymentIDFunc func(ctx context.Context, serviceID string) (string, error)
}

func (m *MockRailwayClient) GetDeploymentLogs(ctx context.Context, input railway.GetDeploymentLogsInput) (railway.GetDeploymentLogsResult, error) {
	if m.GetDeploymentLogsFunc != nil {
		return m.GetDeploymentLogsFunc(ctx, input)
	}
	return railway.GetDeploymentLogsResult{}, nil
}

func (m *MockRailwayClient) GetLatestDeploymentID(ctx context.Context, serviceID string) (string, error) {
	if m.GetLatestDeploymentIDFunc != nil {
		return m.GetLatestDeploymentIDFunc(ctx, serviceID)
	}
	return "mock-deployment-id", nil
}

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&store.Service{}, &store.Environment{}); err != nil {
		panic(err)
	}
	return db
}

func TestGetServiceLogs_Success(t *testing.T) {
	// Setup test database
	db := setupTestDB()

	// Create test service
	service := store.Service{
		ID:               "test-service-id",
		Name:             "test-api",
		RailwayServiceID: "railway-service-123",
		EnvironmentID:    "env-123",
	}
	require.NoError(t, db.Create(&service).Error)

	// Mock Railway client
	mockRailway := &MockRailwayClient{
		GetLatestDeploymentIDFunc: func(ctx context.Context, serviceID string) (string, error) {
			assert.Equal(t, "railway-service-123", serviceID)
			return "deployment-456", nil
		},
		GetDeploymentLogsFunc: func(ctx context.Context, input railway.GetDeploymentLogsInput) (railway.GetDeploymentLogsResult, error) {
			assert.Equal(t, "deployment-456", input.DeploymentID)
			return railway.GetDeploymentLogsResult{
				Logs: []railway.DeploymentLog{
					{
						Timestamp: "2024-01-01T12:00:00Z",
						Message:   "Server started",
						Severity:  "INFO",
					},
					{
						Timestamp: "2024-01-01T12:01:00Z",
						Message:   "Request received",
						Severity:  "INFO",
					},
				},
			}, nil
		},
	}

	// Create controller
	controller := &LogsController{
		DB:      db,
		Railway: mockRailway,
	}

	// Setup Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()
	controller.RegisterRoutes(router.Group("/api/v1"))

	// Create request - use Railway service ID in URL
	req := httptest.NewRequest("GET", "/api/v1/services/railway-service-123/logs?limit=100", nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response LogsResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, 2, response.Count)
	assert.Len(t, response.Logs, 2)
	assert.Equal(t, "test-api", response.Logs[0].ServiceName)
	assert.Contains(t, response.Logs[0].Message, "Server started")
}

func TestGetServiceLogs_ServiceNotFound(t *testing.T) {
	// Setup test database (empty)
	db := setupTestDB()

	// Mock Railway client (should not be called)
	mockRailway := &MockRailwayClient{}

	// Create controller
	controller := &LogsController{
		DB:      db,
		Railway: mockRailway,
	}

	// Setup Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()
	controller.RegisterRoutes(router.Group("/api/v1"))

	// Create request with non-existent service ID
	req := httptest.NewRequest("GET", "/api/v1/services/nonexistent/logs", nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Contains(t, response["error"], "service not found")
}

func TestGetServiceLogs_InvalidLimit(t *testing.T) {
	// Setup test database
	db := setupTestDB()

	// Create test service
	service := store.Service{
		ID:               "test-service-id",
		Name:             "test-api",
		RailwayServiceID: "railway-service-123",
	}
	require.NoError(t, db.Create(&service).Error)

	// Mock Railway client (should not be called)
	mockRailway := &MockRailwayClient{}

	// Create controller
	controller := &LogsController{
		DB:      db,
		Railway: mockRailway,
	}

	// Setup Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()
	controller.RegisterRoutes(router.Group("/api/v1"))

	// Create request with invalid limit - use Railway service ID in URL
	req := httptest.NewRequest("GET", "/api/v1/services/railway-service-123/logs?limit=invalid", nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Contains(t, response["error"], "invalid limit")
}

func TestGetServiceLogs_RailwayClientNotConfigured(t *testing.T) {
	// Setup test database
	db := setupTestDB()

	// Create controller WITHOUT Railway client
	controller := &LogsController{
		DB:      db,
		Railway: nil,
	}

	// Setup Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()
	controller.RegisterRoutes(router.Group("/api/v1"))

	// Create request
	req := httptest.NewRequest("GET", "/api/v1/services/railway-service-123/logs", nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var response map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Contains(t, response["error"], "railway client not configured")
}

func TestExportLogs_CSV(t *testing.T) {
	// Setup test database
	db := setupTestDB()

	// Create test service
	service := store.Service{
		ID:               "test-service-id",
		Name:             "test-api",
		RailwayServiceID: "railway-service-123",
	}
	require.NoError(t, db.Create(&service).Error)

	// Mock Railway client
	mockRailway := &MockRailwayClient{
		GetLatestDeploymentIDFunc: func(ctx context.Context, serviceID string) (string, error) {
			return "deployment-456", nil
		},
		GetDeploymentLogsFunc: func(ctx context.Context, input railway.GetDeploymentLogsInput) (railway.GetDeploymentLogsResult, error) {
			return railway.GetDeploymentLogsResult{
				Logs: []railway.DeploymentLog{
					{
						Timestamp: "2024-01-01T12:00:00Z",
						Message:   "Server started",
						Severity:  "INFO",
					},
				},
			}, nil
		},
	}

	// Create controller
	controller := &LogsController{
		DB:      db,
		Railway: mockRailway,
	}

	// Setup Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()
	controller.RegisterRoutes(router.Group("/api/v1"))

	// Create request - use Railway service ID in query param
	req := httptest.NewRequest("GET", "/api/v1/logs/export?serviceId=railway-service-123&format=csv", nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/csv", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), ".csv")
	assert.Contains(t, w.Body.String(), "Timestamp,Service,Level,Message")
}

func TestExportLogs_JSON(t *testing.T) {
	// Setup test database
	db := setupTestDB()

	// Create test service
	service := store.Service{
		ID:               "test-service-id",
		Name:             "test-api",
		RailwayServiceID: "railway-service-123",
	}
	require.NoError(t, db.Create(&service).Error)

	// Mock Railway client
	mockRailway := &MockRailwayClient{
		GetLatestDeploymentIDFunc: func(ctx context.Context, serviceID string) (string, error) {
			return "deployment-456", nil
		},
		GetDeploymentLogsFunc: func(ctx context.Context, input railway.GetDeploymentLogsInput) (railway.GetDeploymentLogsResult, error) {
			return railway.GetDeploymentLogsResult{
				Logs: []railway.DeploymentLog{
					{
						Timestamp: "2024-01-01T12:00:00Z",
						Message:   "Server started",
						Severity:  "INFO",
					},
				},
			}, nil
		},
	}

	// Create controller
	controller := &LogsController{
		DB:      db,
		Railway: mockRailway,
	}

	// Setup Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()
	controller.RegisterRoutes(router.Group("/api/v1"))

	// Create request - use Railway service ID in query param
	req := httptest.NewRequest("GET", "/api/v1/logs/export?serviceId=railway-service-123&format=json", nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), ".json")

	// Verify valid JSON
	var logs []interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &logs))
}

func TestExportLogs_InvalidFormat(t *testing.T) {
	// Setup test database
	db := setupTestDB()

	// Mock Railway client (should not be called)
	mockRailway := &MockRailwayClient{}

	// Create controller
	controller := &LogsController{
		DB:      db,
		Railway: mockRailway,
	}

	// Setup Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()
	controller.RegisterRoutes(router.Group("/api/v1"))

	// Create request with invalid format - use Railway service ID in query param
	req := httptest.NewRequest("GET", "/api/v1/logs/export?serviceId=railway-service-123&format=xml", nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Contains(t, response["error"], "format must be one of")
}

func TestExportLogs_MissingServiceId(t *testing.T) {
	// Setup test database
	db := setupTestDB()

	// Mock Railway client (should not be called)
	mockRailway := &MockRailwayClient{}

	// Create controller
	controller := &LogsController{
		DB:      db,
		Railway: mockRailway,
	}

	// Setup Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()
	controller.RegisterRoutes(router.Group("/api/v1"))

	// Create request without serviceId
	req := httptest.NewRequest("GET", "/api/v1/logs/export?format=json", nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Contains(t, response["error"], "serviceId")
}
