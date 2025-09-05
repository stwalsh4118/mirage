package jobs

import (
	"context"
	"math/rand"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stwalsh4118/mirageapi/internal/railway"
	"github.com/stwalsh4118/mirageapi/internal/store"
	"gorm.io/gorm"
)

// EnvironmentPublisher allows emitting updates to downstream consumers (e.g., websockets).
type EnvironmentPublisher interface {
	PublishEnvironmentUpdated(environmentID string, newStatus string)
}

// LogPublisher is a minimal publisher that logs updates.
type LogPublisher struct{}

func (LogPublisher) PublishEnvironmentUpdated(environmentID string, newStatus string) {
	log.Info().Str("env_id", environmentID).Str("status", newStatus).Msg("environment status updated")
}

// StartStatusPoller starts a background loop that periodically polls Railway for environment statuses
// and reconciles them with the local database. It returns a stop function to halt the loop.
func StartStatusPoller(
	ctx context.Context,
	db *gorm.DB,
	rw *railway.Client,
	interval time.Duration,
	jitterFraction float64,
	publisher EnvironmentPublisher,
) (stop func()) {
	if publisher == nil {
		publisher = LogPublisher{}
	}

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		log.Info().Dur("interval", interval).Float64("jitter_fraction", jitterFraction).Msg("status poller started")
		defer log.Info().Msg("status poller stopped")
		for {
			// Sleep for interval with jitter
			d := addJitter(interval, jitterFraction)
			select {
			case <-time.After(d):
				if err := pollOnce(ctx, db, rw, publisher); err != nil {
					log.Error().Err(err).Msg("status poller iteration failed")
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return cancel
}

func addJitter(base time.Duration, fraction float64) time.Duration {
	if fraction <= 0 {
		return base
	}
	j := base.Seconds() * fraction
	// random value in [-j/2, +j/2]
	delta := (rand.Float64() - 0.5) * j
	return time.Duration((base.Seconds() + delta) * float64(time.Second))
}

func pollOnce(ctx context.Context, db *gorm.DB, rw *railway.Client, publisher EnvironmentPublisher) error {
	var envs []store.Environment
	if err := db.WithContext(ctx).Where("railway_environment_id <> ''").Find(&envs).Error; err != nil {
		return err
	}
	for _, e := range envs {
		status, err := rw.GetEnvironmentStatus(ctx, e.RailwayEnvironmentID)
		if err != nil {
			log.Error().Err(err).Str("env_id", e.ID).Msg("failed to fetch remote status")
			continue
		}
		changed, newStatus := reconcileEnvironmentStatus(e.Status, status)
		if !changed {
			continue
		}
		if err := db.WithContext(ctx).Model(&e).Update("status", newStatus).Error; err != nil {
			log.Error().Err(err).Str("env_id", e.ID).Str("status", newStatus).Msg("failed to update environment status")
			continue
		}
		publisher.PublishEnvironmentUpdated(e.ID, newStatus)
	}
	return nil
}

// reconcileEnvironmentStatus compares local and remote status strings and returns whether
// a change should be applied and the normalized target status.
func reconcileEnvironmentStatus(local string, remote string) (bool, string) {
	// For MVP, treat remote as source of truth and normalize a few known states.
	n := normalizeStatus(remote)
	if n == "" {
		return false, local
	}
	if n == local {
		return false, local
	}
	return true, n
}

func normalizeStatus(s string) string {
	switch s {
	case "creating", "provisioning", "deploying":
		return "creating"
	case "ready", "healthy", "running":
		return "ready"
	case "error", "failed", "degraded":
		return "error"
	default:
		return s
	}
}
