package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stwalsh4118/mirageapi/internal/config"
	"github.com/stwalsh4118/mirageapi/internal/logging"
)

// NewHTTPServer configures and returns a Gin engine.
func NewHTTPServer(cfg config.AppConfig) *gin.Engine {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(logging.GinLogger())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"env":    cfg.Environment,
		})
	})

	return r
}
