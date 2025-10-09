package server

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stwalsh4118/mirageapi/internal/auth"
	"github.com/stwalsh4118/mirageapi/internal/config"
	"github.com/stwalsh4118/mirageapi/internal/controller"
	"github.com/stwalsh4118/mirageapi/internal/logging"
	"github.com/stwalsh4118/mirageapi/internal/railway"
	"github.com/stwalsh4118/mirageapi/internal/scanner"
	"github.com/stwalsh4118/mirageapi/internal/webhooks"
	"gorm.io/gorm"
)

// NewHTTPServer configures and returns a Gin engine.
func NewHTTPServer(cfg config.AppConfig, deps ...any) *gin.Engine {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(logging.GinLogger())

	// CORS configuration from environment
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Inject Railway project ID into request context for downstream controller usage
	r.Use(func(c *gin.Context) {
		if cfg.RailwayProjectID != "" {
			c.Set("railway_project_id", cfg.RailwayProjectID)
		}
		c.Next()
	})

	var db *gorm.DB
	var rw *railway.Client
	for _, d := range deps {
		switch v := d.(type) {
		case *gorm.DB:
			db = v
		case *railway.Client:
			rw = v
		}
	}

	api := r.Group("/api")

	// Public webhook endpoints (no auth)
	if db != nil && cfg.ClerkWebhookSecret != "" {
		webhookHandler := webhooks.NewClerkWebhookHandler(db, cfg.ClerkWebhookSecret)
		api.POST("/webhooks/clerk", webhookHandler.HandleWebhook)
	}

	// v1 API routes
	v1 := api.Group("/v1")

	// Public health endpoint
	v1.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"env":    cfg.Environment,
		})
	})

	// Apply authentication middleware to all other v1 routes
	if db != nil {
		// Create authenticated route group for regular HTTP routes
		authed := v1.Group("")
		authed.Use(auth.RequireAuth(db))
		{
			if rw != nil {
				ec := &controller.EnvironmentController{DB: db, Railway: rw}
				ec.RegisterRoutes(authed)
				sc := &controller.ServicesController{Railway: rw, DB: db}
				sc.RegisterRoutes(authed)

				// Register non-WebSocket log routes with regular auth
				lc := &controller.LogsController{DB: db, Railway: rw, AllowedOrigins: cfg.AllowedOrigins}
				authed.GET("/services/:id/logs", lc.GetServiceLogs)
				authed.GET("/logs/export", lc.ExportLogs)
			}

			// Initialize Dockerfile scanner for discovery
			githubToken := os.Getenv("GITHUB_SERVICE_TOKEN")
			scanCache := scanner.NewScanCache(scanner.DefaultCacheTTL)
			githubScanner := scanner.NewGitHubScanner(githubToken, scanCache)
			dc := &controller.DiscoveryController{Scanner: githubScanner}
			dc.RegisterRoutes(authed)

			// TODO: Add UserController when task 16-9 is implemented
			// uc := &controller.UserController{DB: db}
			// uc.RegisterRoutes(authed)
		}

		// WebSocket routes - auth happens via first message after connection (not middleware)
		if rw != nil {
			// Register WebSocket log streaming routes
			// Note: Auth is handled inside the handler by reading first message
			lc := &controller.LogsController{DB: db, Railway: rw, AllowedOrigins: cfg.AllowedOrigins}
			v1.GET("/services/:id/logs/stream", lc.StreamServiceLogs)
			v1.GET("/environments/:id/logs/stream", lc.StreamEnvironmentLogs)
		}
	}

	return r
}
