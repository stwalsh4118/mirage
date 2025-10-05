package server

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stwalsh4118/mirageapi/internal/config"
	"github.com/stwalsh4118/mirageapi/internal/controller"
	"github.com/stwalsh4118/mirageapi/internal/logging"
	"github.com/stwalsh4118/mirageapi/internal/railway"
	"github.com/stwalsh4118/mirageapi/internal/scanner"
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
	v1 := api.Group("/v1")
	if db != nil && rw != nil {
		ec := &controller.EnvironmentController{DB: db, Railway: rw}
		ec.RegisterRoutes(v1)
		sc := &controller.ServicesController{Railway: rw, DB: db}
		sc.RegisterRoutes(v1)
		lc := &controller.LogsController{DB: db, Railway: rw, AllowedOrigins: cfg.AllowedOrigins}
		lc.RegisterRoutes(v1)

		// Initialize Dockerfile scanner for discovery
		githubToken := os.Getenv("GITHUB_SERVICE_TOKEN")
		scanCache := scanner.NewScanCache(scanner.DefaultCacheTTL)
		githubScanner := scanner.NewGitHubScanner(githubToken, scanCache)
		dc := &controller.DiscoveryController{Scanner: githubScanner}
		dc.RegisterRoutes(v1)
	}

	v1.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"env":    cfg.Environment,
		})
	})

	return r
}
