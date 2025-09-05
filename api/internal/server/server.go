package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stwalsh4118/mirageapi/internal/config"
	"github.com/stwalsh4118/mirageapi/internal/controller"
	"github.com/stwalsh4118/mirageapi/internal/logging"
	"github.com/stwalsh4118/mirageapi/internal/railway"
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
	if db != nil && rw != nil {
		ec := &controller.EnvironmentController{DB: db, Railway: rw}
		ec.RegisterRoutes(r)
	}

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"env":    cfg.Environment,
		})
	})

	return r
}
