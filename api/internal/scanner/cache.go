package scanner

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	// DefaultCacheTTL is the default time-to-live for cached scan results.
	DefaultCacheTTL = 5 * time.Minute
)

// CacheKey uniquely identifies a repository scan.
type CacheKey struct {
	Owner  string
	Repo   string
	Branch string
}

// String returns a string representation of the cache key.
func (k CacheKey) String() string {
	return fmt.Sprintf("%s/%s@%s", k.Owner, k.Repo, k.Branch)
}

// CachedScan holds scan results with expiration metadata.
type CachedScan struct {
	Result    []DockerfileInfo
	CachedAt  time.Time
	ExpiresAt time.Time
}

// IsExpired checks if the cached scan has expired.
func (c *CachedScan) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

// ScanCache provides in-memory caching for scan results.
type ScanCache struct {
	cache map[string]*CachedScan
	mu    sync.RWMutex
	ttl   time.Duration
}

// NewScanCache creates a new scan cache with the specified TTL.
func NewScanCache(ttl time.Duration) *ScanCache {
	if ttl == 0 {
		ttl = DefaultCacheTTL
	}
	return &ScanCache{
		cache: make(map[string]*CachedScan),
		ttl:   ttl,
	}
}

// Get retrieves a cached scan result if it exists and hasn't expired.
func (c *ScanCache) Get(key CacheKey) ([]DockerfileInfo, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cached, ok := c.cache[key.String()]
	if !ok {
		return nil, false
	}

	if cached.IsExpired() {
		log.Debug().Str("key", key.String()).Msg("cache entry expired")
		return nil, false
	}

	log.Debug().
		Str("key", key.String()).
		Time("cached_at", cached.CachedAt).
		Msg("cache hit")

	return cached.Result, true
}

// Set stores a scan result in the cache.
func (c *ScanCache) Set(key CacheKey, result []DockerfileInfo) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	c.cache[key.String()] = &CachedScan{
		Result:    result,
		CachedAt:  now,
		ExpiresAt: now.Add(c.ttl),
	}

	log.Debug().
		Str("key", key.String()).
		Int("count", len(result)).
		Time("expires_at", now.Add(c.ttl)).
		Msg("cached scan result")
}

// Clear removes expired entries from the cache.
func (c *ScanCache) Clear(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()

	removed := 0

	for key, cached := range c.cache {
		if cached.IsExpired() {
			delete(c.cache, key)
			removed++
		}
	}

	if removed > 0 {
		log.Debug().Int("removed", removed).Msg("cleared expired cache entries")
	}
}

// Size returns the current number of cached entries.
func (c *ScanCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.cache)
}

// StartCleanup starts a background goroutine that periodically clears expired entries.
func (c *ScanCache) StartCleanup(ctx context.Context, interval time.Duration) {
	if interval == 0 {
		interval = 1 * time.Minute
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Debug().Msg("stopping cache cleanup")
				return
			case <-ticker.C:
				c.Clear(ctx)
			}
		}
	}()

	log.Info().Dur("interval", interval).Msg("started cache cleanup")
}
