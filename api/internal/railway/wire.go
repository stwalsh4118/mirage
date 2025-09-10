package railway

import (
	"github.com/stwalsh4118/mirageapi/internal/config"
)

// NewFromConfig builds a Railway client using application configuration.
func NewFromConfig(cfg config.AppConfig) *Client {
	// Use the default 30s timeout defined in NewClient.
	return NewClient(cfg.RailwayEndpoint, cfg.RailwayAPIToken, nil)
}
