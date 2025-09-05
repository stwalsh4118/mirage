package railway

import (
	"net/http"

	"github.com/stwalsh4118/mirageapi/internal/config"
)

// NewFromConfig builds a Railway client using application configuration.
func NewFromConfig(cfg config.AppConfig) *Client {
	return NewClient("", cfg.RailwayAPIToken, &http.Client{})
}
