package logging

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Setup configures zerolog's global settings based on environment.
func Setup(environment string) {
	// Human-friendly time while keeping JSON output compact.
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.TimestampFieldName = "ts"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "msg"

	// Default to info; use debug in non-production
	level := zerolog.InfoLevel
	if environment != "production" {
		level = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(level)

	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

// GinLogger returns a Gin middleware that logs requests using zerolog.
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		size := c.Writer.Size()

		evt := log.Info()
		if len(c.Errors) > 0 || status >= 500 {
			evt = log.Error()
		}

		evt.Str("method", c.Request.Method)
		evt.Str("path", path)
		evt.Str("query", query)
		evt.Int("status", status)
		evt.Dur("latency_ms", latency)
		evt.Int("bytes", size)
		evt.Str("client_ip", c.ClientIP())
		evt.Str("user_agent", c.Request.UserAgent())
		evt.Msg("http_request")
	}
}
