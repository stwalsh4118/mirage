package webhooks

import (
	"context"
	"sync"
	"time"
)

// WebhookIDCache tracks processed webhook IDs to prevent duplicates
type WebhookIDCache struct {
	mu     sync.RWMutex
	ids    map[string]time.Time
	ttl    time.Duration
	cancel context.CancelFunc
}

// NewWebhookIDCache creates a new webhook ID cache with the specified TTL
func NewWebhookIDCache(ttl time.Duration) *WebhookIDCache {
	ctx, cancel := context.WithCancel(context.Background())

	cache := &WebhookIDCache{
		ids:    make(map[string]time.Time),
		ttl:    ttl,
		cancel: cancel,
	}

	// Cleanup expired entries periodically
	go cache.cleanup(ctx)

	return cache
}

// MarkProcessed marks a webhook ID as processed
func (c *WebhookIDCache) MarkProcessed(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ids[id] = time.Now().Add(c.ttl)
}

// IsProcessed checks if a webhook ID has been processed
func (c *WebhookIDCache) IsProcessed(id string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	expiry, exists := c.ids[id]
	if !exists {
		return false
	}

	// Check if expired
	return time.Now().Before(expiry)
}

// Stop stops the cleanup goroutine and releases resources
func (c *WebhookIDCache) Stop() {
	if c.cancel != nil {
		c.cancel()
	}
}

func (c *WebhookIDCache) cleanup(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Context cancelled, stop cleanup
			return
		case <-ticker.C:
			c.mu.Lock()
			now := time.Now()
			for id, expiry := range c.ids {
				if now.After(expiry) {
					delete(c.ids, id)
				}
			}
			c.mu.Unlock()
		}
	}
}
