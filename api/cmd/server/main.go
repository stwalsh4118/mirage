package main

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/stwalsh4118/mirageapi/internal/config"
	"github.com/stwalsh4118/mirageapi/internal/jobs"
	"github.com/stwalsh4118/mirageapi/internal/logging"
	"github.com/stwalsh4118/mirageapi/internal/railway"
	"github.com/stwalsh4118/mirageapi/internal/server"
	"github.com/stwalsh4118/mirageapi/internal/store"
)

func main() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	logging.Setup(cfg.Environment)

	// Initialize DB (SQLite MVP) using env var SQLITE_PATH if provided
	sqlitePath := os.Getenv("SQLITE_PATH")
	db, err := store.Open(sqlitePath)
	if err != nil {
		log.Fatal().Err(err).Str("path", sqlitePath).Msg("failed to init database")
	}

	rw := railway.NewFromConfig(cfg)

	// Start status poller (Phase 1)
	pollInterval := time.Duration(cfg.PollIntervalSeconds) * time.Second
	_ = jobs.StartStatusPoller(
		context.Background(),
		db,
		rw,
		pollInterval,
		cfg.PollJitterFraction,
		nil, // use default log publisher
	)

	engine := server.NewHTTPServer(cfg, db, rw)

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
