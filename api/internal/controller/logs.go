package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/stwalsh4118/mirageapi/internal/logutil"
	"github.com/stwalsh4118/mirageapi/internal/railway"
	"github.com/stwalsh4118/mirageapi/internal/store"
	"gorm.io/gorm"
)

// RailwayLogsClient defines the interface for Railway log operations needed by the controller
type RailwayLogsClient interface {
	GetDeploymentLogs(ctx context.Context, input railway.GetDeploymentLogsInput) (railway.GetDeploymentLogsResult, error)
	GetLatestDeploymentID(ctx context.Context, serviceID string) (string, error)
}

// LogsController handles log retrieval and export endpoints
type LogsController struct {
	DB      *gorm.DB
	Railway RailwayLogsClient
}

// RegisterRoutes registers log-related routes under the provided router group
func (c *LogsController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/services/:id/logs", c.GetServiceLogs)
	r.GET("/logs/export", c.ExportLogs)
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

	// Get Railway service ID from URL parameter
	railwayServiceID := ctx.Param("id")
	if railwayServiceID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "service id required"})
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
		Str("railway_service_id", railwayServiceID).
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
