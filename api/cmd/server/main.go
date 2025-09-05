package main

import (
	"github.com/rs/zerolog/log"

	"github.com/stwalsh4118/mirageapi/internal/config"
	"github.com/stwalsh4118/mirageapi/internal/logging"
	"github.com/stwalsh4118/mirageapi/internal/server"
)

func main() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	logging.Setup(cfg.Environment)

	engine := server.NewHTTPServer(cfg)

	port := cfg.HTTPPort
	if port == "" {
		port = config.DefaultHTTPPort
	}

	addr := ":" + port
	log.Info().Str("addr", addr).Str("env", cfg.Environment).Msg("starting api")
	if err := engine.Run(addr); err != nil {
		log.Error().Err(err).Msg("server exited with error")
	}
}
