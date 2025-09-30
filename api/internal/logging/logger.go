package logging

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	if environment != "production" {
		cw := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}
		log.Logger = zerolog.New(cw).With().Timestamp().Logger()
	} else {
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
}

// GinLogger returns a Gin middleware that logs requests using zerolog.
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Ensure request id exists and is propagated
		reqID := c.GetHeader("X-Request-Id")
		if reqID == "" {
			reqID = uuid.NewString()
			c.Request.Header.Set("X-Request-Id", reqID)
		}
		c.Writer.Header().Set("X-Request-Id", reqID)

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		size := c.Writer.Size()

		evt := log.Info()
		if len(c.Errors) > 0 || status >= 500 {
			evt = log.Error()
		}

		evt.Str("request_id", reqID)
		evt.Str("method", c.Request.Method)
		evt.Str("path", path)
		evt.Str("query", query)
		evt.Int("status", status)
		evt.Float64("latency_ms", float64(latency)/float64(time.Millisecond))
		evt.Int("bytes", size)
		evt.Str("client_ip", c.ClientIP())
		evt.Str("user_agent", c.Request.UserAgent())
		evt.Msg("http_request")
	}
}
