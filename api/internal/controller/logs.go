package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/stwalsh4118/mirageapi/internal/auth"
	"github.com/stwalsh4118/mirageapi/internal/logutil"
	"github.com/stwalsh4118/mirageapi/internal/railway"
	"github.com/stwalsh4118/mirageapi/internal/store"
	"gorm.io/gorm"
)

// RailwayLogsClient defines the interface for Railway log operations needed by the controller
type RailwayLogsClient interface {
	GetDeploymentLogs(ctx context.Context, input railway.GetDeploymentLogsInput) (railway.GetDeploymentLogsResult, error)
	GetLatestDeploymentID(ctx context.Context, serviceID string) (string, error)
	SubscribeToEnvironmentLogs(ctx context.Context, environmentID string, serviceFilter string) (*websocket.Conn, error)
	SubscribeToDeploymentLogs(ctx context.Context, deploymentID string, filter string) (*websocket.Conn, error)
}

// LogsController handles log retrieval and export endpoints
type LogsController struct {
	DB               *gorm.DB
	Railway          RailwayLogsClient
	AllowedOrigins   []string
	serviceNameCache sync.Map // railwayServiceID (string) -> serviceName (string)
}

// RegisterRoutes registers log-related routes under the provided router group
func (c *LogsController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/services/:id/logs", c.GetServiceLogs)
	r.GET("/services/:id/logs/stream", c.StreamServiceLogs)
	r.GET("/logs/export", c.ExportLogs)
	r.GET("/environments/:id/logs/stream", c.StreamEnvironmentLogs)
}

// ParsedLogDTO represents a parsed log entry for API response
type ParsedLogDTO struct {
	Timestamp   string `json:"timestamp"`
	ServiceName string `json:"serviceName"`
	Severity    string `json:"severity"`
	Message     string `json:"message"`
	RawLine     string `json:"rawLine"`
}

// LogsResponse is the standard response structure for log endpoints
type LogsResponse struct {
	Logs  []ParsedLogDTO `json:"logs"`
	Count int            `json:"count"`
}

// GetServiceLogs fetches historical logs for a specific service
// GET /api/v1/services/:id/logs?limit=500&search=error&minSeverity=WARN
func (c *LogsController) GetServiceLogs(ctx *gin.Context) {
	if c.Railway == nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "railway client not configured"})
		return
	}

	// Get authenticated user
	user, err := auth.GetCurrentUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	// Get Mirage service ID from URL parameter
	serviceID := ctx.Param("id")
	if serviceID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "service id required"})
		return
	}

	// Look up service in database by Mirage ID with ownership check
	var service store.Service
	err = c.DB.Where("id = ? AND user_id = ?", serviceID, user.ID).First(&service).Error
	if err == gorm.ErrRecordNotFound {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
		return
	} else if err != nil {
		log.Error().Err(err).Str("service_id", serviceID).Msg("failed to query service")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve service"})
		return
	}

	// Parse query parameters
	limit := 500
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
			return
		}
	}
	if limit > 1000 {
		limit = 1000 // Enforce maximum
	}
	if limit <= 0 {
		limit = 500 // Default
	}

	searchQuery := ctx.Query("search")
	minSeverity := ctx.Query("minSeverity")

	// Get latest deployment ID for the service from Railway
	if service.RailwayServiceID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "service has no railway service id"})
		return
	}

	deploymentID, err := c.Railway.GetLatestDeploymentID(ctx, service.RailwayServiceID)
	if err != nil {
		log.Error().Err(err).Str("railway_service_id", service.RailwayServiceID).Msg("failed to get latest deployment")
		ctx.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("failed to get deployment: %s", err.Error())})
		return
	}

	log.Info().
		Str("service_id", serviceID).
		Str("railway_service_id", service.RailwayServiceID).
		Str("service_name", service.Name).
		Str("deployment_id", deploymentID).
		Int("limit", limit).
		Str("search", searchQuery).
		Str("min_severity", minSeverity).
		Msg("fetching service logs")

	// Fetch logs from Railway
	railwayInput := railway.GetDeploymentLogsInput{
		DeploymentID: deploymentID,
		Limit:        limit,
		Filter:       searchQuery, // Railway handles text filtering
	}

	railwayResult, err := c.Railway.GetDeploymentLogs(ctx, railwayInput)
	if err != nil {
		log.Error().Err(err).Str("deployment_id", deploymentID).Msg("railway get deployment logs failed")
		ctx.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("failed to fetch logs: %s", err.Error())})
		return
	}

	// Parse and filter logs
	parsedLogs := make([]ParsedLogDTO, 0, len(railwayResult.Logs))
	minPriority := logutil.SeverityPriority(minSeverity)

	for _, railwayLog := range railwayResult.Logs {
		// Parse the log line
		parsed := logutil.ParseLogLine(railwayLog.Message, service.Name)

		// Use Railway's severity if provided, otherwise use detected severity
		if railwayLog.Severity != "" {
			parsed.Severity = logutil.NormalizeSeverity(railwayLog.Severity)
		}

		// Use Railway's timestamp if provided
		if railwayLog.Timestamp != "" {
			if ts, err := time.Parse(time.RFC3339Nano, railwayLog.Timestamp); err == nil {
				parsed.Timestamp = ts
			}
		}

		// Apply severity filter
		if minSeverity != "" && logutil.SeverityPriority(parsed.Severity) < minPriority {
			continue
		}

		parsedLogs = append(parsedLogs, ParsedLogDTO{
			Timestamp:   parsed.Timestamp.Format(time.RFC3339),
			ServiceName: service.Name,
			Severity:    parsed.Severity,
			Message:     parsed.Message,
			RawLine:     parsed.RawLine,
		})
	}

	log.Info().
		Int("fetched", len(railwayResult.Logs)).
		Int("returned", len(parsedLogs)).
		Msg("service logs processed")

	ctx.JSON(http.StatusOK, LogsResponse{
		Logs:  parsedLogs,
		Count: len(parsedLogs),
	})
}

// ExportLogs exports logs in the specified format (JSON, CSV, TXT)
// GET /api/v1/logs/export?serviceId=abc&format=csv&limit=1000
func (c *LogsController) ExportLogs(ctx *gin.Context) {
	if c.Railway == nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "railway client not configured"})
		return
	}

	// Get Railway service ID from query parameter
	railwayServiceID := ctx.Query("serviceId")
	if railwayServiceID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "serviceId query parameter required"})
		return
	}

	// Get format (default: json)
	format := strings.ToLower(ctx.Query("format"))
	if format == "" {
		format = "json"
	}
	if format != "json" && format != "csv" && format != "txt" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "format must be one of: json, csv, txt"})
		return
	}

	// Look up service in database by Railway service ID
	var service store.Service
	if err := c.DB.Where("railway_service_id = ?", railwayServiceID).First(&service).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
			return
		}
		log.Error().Err(err).Str("railway_service_id", railwayServiceID).Msg("failed to query service")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve service"})
		return
	}

	// Parse limit
	limit := 1000
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}
	if limit > 1000 {
		limit = 1000
	}

	// Get latest deployment ID for the service from Railway
	if service.RailwayServiceID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "service has no railway service id"})
		return
	}

	deploymentID, err := c.Railway.GetLatestDeploymentID(ctx, service.RailwayServiceID)
	if err != nil {
		log.Error().Err(err).Str("railway_service_id", service.RailwayServiceID).Msg("failed to get latest deployment")
		ctx.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("failed to get deployment: %s", err.Error())})
		return
	}

	log.Info().
		Str("railway_service_id", railwayServiceID).
		Str("service_name", service.Name).
		Str("format", format).
		Int("limit", limit).
		Msg("exporting service logs")

	// Fetch logs from Railway
	railwayInput := railway.GetDeploymentLogsInput{
		DeploymentID: deploymentID,
		Limit:        limit,
	}

	railwayResult, err := c.Railway.GetDeploymentLogs(ctx, railwayInput)
	if err != nil {
		log.Error().Err(err).Str("deployment_id", deploymentID).Msg("railway get deployment logs failed")
		ctx.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("failed to fetch logs: %s", err.Error())})
		return
	}

	// Parse logs
	parsedLogs := make([]logutil.ParsedLog, 0, len(railwayResult.Logs))
	for _, railwayLog := range railwayResult.Logs {
		parsed := logutil.ParseLogLine(railwayLog.Message, service.Name)

		// Use Railway's severity if provided
		if railwayLog.Severity != "" {
			parsed.Severity = logutil.NormalizeSeverity(railwayLog.Severity)
		}

		// Use Railway's timestamp if provided
		if railwayLog.Timestamp != "" {
			if ts, err := time.Parse(time.RFC3339Nano, railwayLog.Timestamp); err == nil {
				parsed.Timestamp = ts
			}
		}

		parsedLogs = append(parsedLogs, parsed)
	}

	// Format logs based on requested format
	var output []byte
	var contentType string
	var fileExt string

	switch format {
	case "json":
		output, err = logutil.FormatAsJSON(parsedLogs)
		contentType = "application/json"
		fileExt = "json"
	case "csv":
		output, err = logutil.FormatAsCSV(parsedLogs)
		contentType = "text/csv"
		fileExt = "csv"
	case "txt":
		output, err = logutil.FormatAsPlainText(parsedLogs)
		contentType = "text/plain"
		fileExt = "txt"
	}

	if err != nil {
		log.Error().Err(err).Str("format", format).Msg("failed to format logs")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to format logs"})
		return
	}

	// Generate filename: service-name-logs-timestamp.ext
	timestamp := time.Now().UTC().Format("20060102-150405")
	filename := fmt.Sprintf("%s-logs-%s.%s", service.Name, timestamp, fileExt)

	// Set headers for file download
	ctx.Header("Content-Type", contentType)
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Header("Content-Length", strconv.Itoa(len(output)))

	log.Info().
		Str("filename", filename).
		Int("bytes", len(output)).
		Msg("logs exported successfully")

	ctx.Data(http.StatusOK, contentType, output)
}

// WebSocket message types
const (
	messageTypeLog    = "log"
	messageTypeStatus = "status"
	messageTypeError  = "error"
	messageTypePing   = "ping"
)

// WebSocketMessage is the standard message format for WebSocket communication
type WebSocketMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

// StreamEnvironmentLogs streams real-time logs from Railway to frontend clients via WebSocket
// GET /api/v1/environments/:id/logs/stream?services=svc1,svc2
// Auth is handled via first message after connection (token sent encrypted in WebSocket payload)
func (c *LogsController) StreamEnvironmentLogs(ginCtx *gin.Context) {
	if c.Railway == nil {
		ginCtx.JSON(http.StatusServiceUnavailable, gin.H{"error": "railway client not configured"})
		return
	}

	// Get environment ID from URL parameter
	environmentID := ginCtx.Param("id")
	if environmentID == "" {
		ginCtx.JSON(http.StatusBadRequest, gin.H{"error": "environment id required"})
		return
	}

	// Get allowed origins from controller config (defaults to wildcard for development)
	allowedOrigins := c.AllowedOrigins
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"*"} // Fallback for development only
		log.Warn().Msg("no allowed origins configured for websocket, using wildcard (not recommended for production)")
	}

	// Upgrade HTTP connection to WebSocket FIRST (before auth)
	conn, err := websocket.Accept(ginCtx.Writer, ginCtx.Request, &websocket.AcceptOptions{
		OriginPatterns: allowedOrigins,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to upgrade to websocket")
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "connection closed")

	// Create context with timeout for auth message
	authCtx, authCancel := context.WithTimeout(ginCtx.Request.Context(), 5*time.Second)
	defer authCancel()

	// Read first message - must be auth message with JWT token
	_, authMsgBytes, err := conn.Read(authCtx)
	if err != nil {
		log.Error().Err(err).Msg("failed to read auth message")
		c.sendWebSocketMessage(ginCtx.Request.Context(), conn, messageTypeError, "authentication required")
		conn.Close(websocket.StatusPolicyViolation, "no auth message received")
		return
	}

	// Parse auth message
	var authMsg struct {
		Type  string `json:"type"`
		Token string `json:"token"`
	}
	if err := json.Unmarshal(authMsgBytes, &authMsg); err != nil {
		log.Error().Err(err).Msg("failed to parse auth message")
		c.sendWebSocketMessage(ginCtx.Request.Context(), conn, messageTypeError, "invalid auth message format")
		conn.Close(websocket.StatusPolicyViolation, "invalid auth message")
		return
	}

	if authMsg.Type != "auth" || authMsg.Token == "" {
		log.Error().Msg("auth message missing type or token")
		c.sendWebSocketMessage(ginCtx.Request.Context(), conn, messageTypeError, "invalid auth message")
		conn.Close(websocket.StatusPolicyViolation, "invalid auth message")
		return
	}

	// Verify JWT token and get user
	user, err := auth.VerifyAndLoadUser(ginCtx.Request.Context(), c.DB, authMsg.Token)
	if err != nil {
		log.Error().Err(err).Msg("failed to verify auth token")
		c.sendWebSocketMessage(ginCtx.Request.Context(), conn, messageTypeError, "authentication failed")
		conn.Close(websocket.StatusPolicyViolation, "authentication failed")
		return
	}

	log.Info().
		Str("user_id", user.ID).
		Str("environment_id", environmentID).
		Msg("websocket authenticated successfully")

	// Look up environment to verify it exists and user owns it
	var env store.Environment
	if err := c.DB.Where("railway_environment_id = ? AND user_id = ?", environmentID, user.ID).First(&env).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.sendWebSocketMessage(ginCtx.Request.Context(), conn, messageTypeError, "environment not found")
			conn.Close(websocket.StatusPolicyViolation, "environment not found")
			return
		}
		log.Error().Err(err).Str("environment_id", environmentID).Msg("failed to query environment")
		c.sendWebSocketMessage(ginCtx.Request.Context(), conn, messageTypeError, "failed to retrieve environment")
		conn.Close(websocket.StatusInternalError, "database error")
		return
	}

	// Parse optional service filter from query params
	serviceFilter := ginCtx.Query("services")

	log.Info().
		Str("environment_id", environmentID).
		Str("service_filter", serviceFilter).
		Str("user_id", user.ID).
		Msg("client authenticated and connecting to environment log stream")

	// Send initial status message
	if err := c.sendWebSocketMessage(ginCtx, conn, messageTypeStatus, "connected"); err != nil {
		log.Error().Err(err).Msg("failed to send status message")
		return
	}

	// Create context for this connection
	ctx, cancel := context.WithCancel(ginCtx.Request.Context())
	defer cancel()

	// Subscribe to Railway logs
	railwayConn, err := c.Railway.SubscribeToEnvironmentLogs(ctx, environmentID, serviceFilter)
	if err != nil {
		log.Error().Err(err).
			Str("environment_id", environmentID).
			Msg("failed to subscribe to railway logs")
		c.sendWebSocketMessage(ginCtx, conn, messageTypeError, fmt.Sprintf("failed to subscribe: %s", err.Error()))
		return
	}
	defer railwayConn.Close(websocket.StatusNormalClosure, "unsubscribing")

	log.Info().
		Str("environment_id", environmentID).
		Msg("railway subscription established")

	// Start goroutines for reading from Railway and handling client pings
	errChan := make(chan error, 2)

	// Goroutine 1: Read from Railway and relay to frontend
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				railwayLog, err := railway.ReadLogMessage(ctx, railwayConn)
				if err != nil {
					errChan <- fmt.Errorf("railway read error: %w", err)
					return
				}

				// Skip nil messages (non-data messages)
				if railwayLog == nil {
					continue
				}

				// Parse and format the log
				parsed := logutil.ParseLogLine(railwayLog.Message, "")

				// Use Railway's severity if provided
				if railwayLog.Severity != "" {
					parsed.Severity = logutil.NormalizeSeverity(railwayLog.Severity)
				}

				// Use Railway's timestamp if provided
				if railwayLog.Timestamp != "" {
					if ts, err := time.Parse(time.RFC3339Nano, railwayLog.Timestamp); err == nil {
						parsed.Timestamp = ts
					}
				}

				// Get service name from tags if available (with caching)
				serviceName := "unknown"
				if railwayLog.Tags != nil {
					if svcID, ok := railwayLog.Tags["serviceId"]; ok {
						serviceName = c.getServiceName(svcID)
					}
				}
				parsed.ServiceName = serviceName

				// Send log to frontend client
				logDTO := ParsedLogDTO{
					Timestamp:   parsed.Timestamp.Format(time.RFC3339),
					ServiceName: serviceName,
					Severity:    parsed.Severity,
					Message:     parsed.Message,
					RawLine:     parsed.RawLine,
				}

				if err := c.sendWebSocketMessage(ctx, conn, messageTypeLog, logDTO); err != nil {
					errChan <- fmt.Errorf("frontend write error: %w", err)
					return
				}
			}
		}
	}()

	// Goroutine 2: Read from frontend (handle pings/pongs and disconnects)
	go func() {
		for {
			_, _, err := conn.Read(ctx)
			if err != nil {
				errChan <- fmt.Errorf("frontend read error: %w", err)
				return
			}
			// Client sent a message (probably ping), just acknowledge by continuing
		}
	}()

	// Wait for error or context cancellation
	select {
	case err := <-errChan:
		log.Info().Err(err).
			Str("environment_id", environmentID).
			Msg("websocket stream ended")
	case <-ctx.Done():
		log.Info().
			Str("environment_id", environmentID).
			Msg("websocket stream cancelled")
	}
}

// getServiceName retrieves service name from cache or database
// Uses sync.Map for concurrent access without explicit locking
func (c *LogsController) getServiceName(railwayServiceID string) string {
	// Check cache first
	if cachedName, ok := c.serviceNameCache.Load(railwayServiceID); ok {
		return cachedName.(string)
	}

	// Query database on cache miss
	var service store.Service
	if err := c.DB.Where("railway_service_id = ?", railwayServiceID).First(&service).Error; err != nil {
		// Don't cache "unknown" to allow retries if service is added later
		return "unknown"
	}

	// Store in cache for future lookups
	c.serviceNameCache.Store(railwayServiceID, service.Name)

	return service.Name
}

// StreamServiceLogs streams real-time logs from a specific service's deployment to frontend clients via WebSocket
// GET /api/v1/services/:id/logs/stream?search=error
// Auth is handled via first message after connection (token sent encrypted in WebSocket payload)
func (c *LogsController) StreamServiceLogs(ginCtx *gin.Context) {
	if c.Railway == nil {
		ginCtx.JSON(http.StatusServiceUnavailable, gin.H{"error": "railway client not configured"})
		return
	}

	// Get Mirage service ID from URL parameter
	serviceID := ginCtx.Param("id")
	if serviceID == "" {
		ginCtx.JSON(http.StatusBadRequest, gin.H{"error": "service id required"})
		return
	}

	// Get allowed origins from controller config (defaults to wildcard for development)
	allowedOrigins := c.AllowedOrigins
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"*"} // Fallback for development only
		log.Warn().Msg("no allowed origins configured for websocket, using wildcard (not recommended for production)")
	}

	// Upgrade HTTP connection to WebSocket FIRST (before auth)
	conn, err := websocket.Accept(ginCtx.Writer, ginCtx.Request, &websocket.AcceptOptions{
		OriginPatterns: allowedOrigins,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to upgrade to websocket")
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "connection closed")

	// Create context with timeout for auth message
	authCtx, authCancel := context.WithTimeout(ginCtx.Request.Context(), 5*time.Second)
	defer authCancel()

	// Read first message - must be auth message with JWT token
	_, authMsgBytes, err := conn.Read(authCtx)
	if err != nil {
		log.Error().Err(err).Msg("failed to read auth message")
		c.sendWebSocketMessage(ginCtx.Request.Context(), conn, messageTypeError, "authentication required")
		conn.Close(websocket.StatusPolicyViolation, "no auth message received")
		return
	}

	// Parse auth message
	var authMsg struct {
		Type  string `json:"type"`
		Token string `json:"token"`
	}
	if err := json.Unmarshal(authMsgBytes, &authMsg); err != nil {
		log.Error().Err(err).Msg("failed to parse auth message")
		c.sendWebSocketMessage(ginCtx.Request.Context(), conn, messageTypeError, "invalid auth message format")
		conn.Close(websocket.StatusPolicyViolation, "invalid auth message")
		return
	}

	if authMsg.Type != "auth" || authMsg.Token == "" {
		log.Error().Msg("auth message missing type or token")
		c.sendWebSocketMessage(ginCtx.Request.Context(), conn, messageTypeError, "invalid auth message")
		conn.Close(websocket.StatusPolicyViolation, "invalid auth message")
		return
	}

	// Verify JWT token and get user
	user, err := auth.VerifyAndLoadUser(ginCtx.Request.Context(), c.DB, authMsg.Token)
	if err != nil {
		log.Error().Err(err).Msg("failed to verify auth token")
		c.sendWebSocketMessage(ginCtx.Request.Context(), conn, messageTypeError, "authentication failed")
		conn.Close(websocket.StatusPolicyViolation, "authentication failed")
		return
	}

	log.Info().
		Str("user_id", user.ID).
		Str("service_id", serviceID).
		Msg("websocket authenticated successfully")

	// Look up service in database by Mirage ID with ownership check
	var service store.Service
	err = c.DB.Where("id = ? AND user_id = ?", serviceID, user.ID).First(&service).Error
	if err == gorm.ErrRecordNotFound {
		c.sendWebSocketMessage(ginCtx.Request.Context(), conn, messageTypeError, "service not found")
		conn.Close(websocket.StatusPolicyViolation, "service not found")
		return
	} else if err != nil {
		log.Error().Err(err).Str("service_id", serviceID).Msg("failed to query service")
		c.sendWebSocketMessage(ginCtx.Request.Context(), conn, messageTypeError, "failed to retrieve service")
		conn.Close(websocket.StatusInternalError, "database error")
		return
	}

	// Get latest deployment ID for the service
	deploymentID, err := c.Railway.GetLatestDeploymentID(ginCtx, service.RailwayServiceID)
	if err != nil {
		log.Error().Err(err).Str("railway_service_id", service.RailwayServiceID).Msg("failed to get latest deployment")
		c.sendWebSocketMessage(ginCtx.Request.Context(), conn, messageTypeError, fmt.Sprintf("failed to get deployment: %s", err.Error()))
		conn.Close(websocket.StatusInternalError, "deployment lookup failed")
		return
	}

	// Parse optional search filter from query params
	searchFilter := ginCtx.Query("search")

	log.Info().
		Str("mirage_service_id", serviceID).
		Str("railway_service_id", service.RailwayServiceID).
		Str("service_name", service.Name).
		Str("deployment_id", deploymentID).
		Str("search_filter", searchFilter).
		Str("user_id", user.ID).
		Msg("client authenticated and connecting to service log stream")

	// Send initial status message
	if err := c.sendWebSocketMessage(ginCtx, conn, messageTypeStatus, "connected"); err != nil {
		log.Error().Err(err).Msg("failed to send status message")
		return
	}

	// Create context for this connection
	ctx, cancel := context.WithCancel(ginCtx.Request.Context())
	defer cancel()

	// Subscribe to Railway deployment logs
	railwayConn, err := c.Railway.SubscribeToDeploymentLogs(ctx, deploymentID, searchFilter)
	if err != nil {
		log.Error().Err(err).
			Str("deployment_id", deploymentID).
			Msg("failed to subscribe to railway deployment logs")
		c.sendWebSocketMessage(ginCtx, conn, messageTypeError, fmt.Sprintf("failed to subscribe: %s", err.Error()))
		return
	}
	defer railwayConn.Close(websocket.StatusNormalClosure, "unsubscribing")

	log.Info().
		Str("deployment_id", deploymentID).
		Str("service_name", service.Name).
		Msg("railway deployment logs subscription established")

	// Start goroutines for reading from Railway and handling client pings
	errChan := make(chan error, 2)

	// Goroutine 1: Read from Railway and relay to frontend
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// Read log messages from Railway (returns an array)
				railwayLogs, err := railway.ReadDeploymentLogMessage(ctx, railwayConn)
				if err != nil {
					errChan <- fmt.Errorf("read railway message: %w", err)
					return
				}

				// Skip non-log messages (status, ack, etc.)
				if railwayLogs == nil {
					continue
				}

				// Process each log in the batch
				for _, railwayLog := range railwayLogs {
					// Parse the log line
					parsed := logutil.ParseLogLine(railwayLog.Message, service.Name)

					// Use Railway's severity if provided, otherwise use detected severity
					if railwayLog.Severity != "" {
						parsed.Severity = logutil.NormalizeSeverity(railwayLog.Severity)
					}

					// Use Railway's timestamp if provided
					if railwayLog.Timestamp != "" {
						if ts, err := time.Parse(time.RFC3339Nano, railwayLog.Timestamp); err == nil {
							parsed.Timestamp = ts
						}
					}

					// Create log DTO for frontend
					logDTO := ParsedLogDTO{
						Timestamp:   parsed.Timestamp.Format(time.RFC3339),
						ServiceName: service.Name,
						Severity:    parsed.Severity,
						Message:     parsed.Message,
						RawLine:     parsed.RawLine,
					}

					// Send log to frontend client
					if err := c.sendWebSocketMessage(ctx, conn, messageTypeLog, logDTO); err != nil {
						errChan <- fmt.Errorf("send to frontend: %w", err)
						return
					}
				}
			}
		}
	}()

	// Goroutine 2: Read from frontend (for ping/disconnect detection)
	go func() {
		for {
			_, _, err := conn.Read(ctx)
			if err != nil {
				errChan <- fmt.Errorf("read from frontend: %w", err)
				return
			}
			// Client sent a message (likely a ping) - no action needed
		}
	}()

	// Wait for error or context cancellation
	select {
	case err := <-errChan:
		// Check if this is a normal client disconnect or actual error
		if err != nil && (err.Error() == "read from frontend: failed to get reader: received close frame: status = StatusNormalClosure and reason = \"Client disconnect\"" ||
			strings.Contains(err.Error(), "StatusNormalClosure")) {
			log.Info().
				Str("deployment_id", deploymentID).
				Str("service_name", service.Name).
				Msg("websocket stream closed normally by client")
		} else {
			log.Info().Err(err).
				Str("deployment_id", deploymentID).
				Str("service_name", service.Name).
				Msg("websocket stream ended with error")
		}
	case <-ctx.Done():
		log.Info().
			Str("deployment_id", deploymentID).
			Str("service_name", service.Name).
			Msg("websocket stream cancelled")
	}
}

// sendWebSocketMessage sends a typed message to the WebSocket client
func (c *LogsController) sendWebSocketMessage(ctx context.Context, conn *websocket.Conn, msgType string, data interface{}) error {
	msg := WebSocketMessage{
		Type: msgType,
		Data: data,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	if err := conn.Write(ctx, websocket.MessageText, msgBytes); err != nil {
		return fmt.Errorf("write message: %w", err)
	}

	return nil
}
