package scanner

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCacheKey_String(t *testing.T) {
	key := CacheKey{
		Owner:  "octocat",
		Repo:   "Hello-World",
		Branch: "main",
	}

	expected := "octocat/Hello-World@main"
	assert.Equal(t, expected, key.String())
}

func TestCachedScan_IsExpired(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		expiresAt time.Time
		expected  bool
	}{
		{
			name:      "not expired",
			expiresAt: now.Add(5 * time.Minute),
			expected:  false,
		},
		{
			name:      "expired",
			expiresAt: now.Add(-1 * time.Minute),
			expected:  true,
		},
		{
			name:      "just expired",
			expiresAt: now.Add(-1 * time.Second),
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cached := &CachedScan{
				ExpiresAt: tt.expiresAt,
			}
			assert.Equal(t, tt.expected, cached.IsExpired())
		})
	}
}

func TestScanCache_GetSet(t *testing.T) {
	cache := NewScanCache(5 * time.Minute)

	key := CacheKey{Owner: "test", Repo: "repo", Branch: "main"}
	result := []DockerfileInfo{
		{Path: "Dockerfile", ServiceName: "root"},
	}

	// Cache miss
	cached, ok := cache.Get(key)
	assert.False(t, ok)
	assert.Nil(t, cached)

	// Set value
	cache.Set(key, result)

	// Cache hit
	cached, ok = cache.Get(key)
	assert.True(t, ok)
	assert.Equal(t, result, cached)
}

func TestScanCache_Expiration(t *testing.T) {
	// Use very short TTL for testing
	cache := NewScanCache(100 * time.Millisecond)

	key := CacheKey{Owner: "test", Repo: "repo", Branch: "main"}
	result := []DockerfileInfo{
		{Path: "Dockerfile", ServiceName: "root"},
	}

	// Set value
	cache.Set(key, result)

	// Immediate cache hit
	cached, ok := cache.Get(key)
	assert.True(t, ok)
	assert.Equal(t, result, cached)

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Cache miss after expiration
	cached, ok = cache.Get(key)
	assert.False(t, ok)
	assert.Nil(t, cached)
}

func TestScanCache_Clear(t *testing.T) {
	cache := NewScanCache(100 * time.Millisecond)

	// Add multiple entries with unique keys
	for i := 0; i < 5; i++ {
		key := CacheKey{
			Owner:  "owner",
			Repo:   "repo" + string(rune('a'+i)),
			Branch: "main",
		}
		cache.Set(key, []DockerfileInfo{})
	}

	assert.Equal(t, 5, cache.Size())

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Clear expired entries
	cache.Clear(context.Background())

	assert.Equal(t, 0, cache.Size())
}

func TestScanCache_Size(t *testing.T) {
	cache := NewScanCache(5 * time.Minute)

	assert.Equal(t, 0, cache.Size())

	// Add entries
	cache.Set(CacheKey{Owner: "a", Repo: "1", Branch: "main"}, []DockerfileInfo{})
	assert.Equal(t, 1, cache.Size())

	cache.Set(CacheKey{Owner: "b", Repo: "2", Branch: "main"}, []DockerfileInfo{})
	assert.Equal(t, 2, cache.Size())

	cache.Set(CacheKey{Owner: "c", Repo: "3", Branch: "main"}, []DockerfileInfo{})
	assert.Equal(t, 3, cache.Size())
}

func TestScanCache_DefaultTTL(t *testing.T) {
	// Create cache with zero TTL, should use default
	cache := NewScanCache(0)

	key := CacheKey{Owner: "test", Repo: "repo", Branch: "main"}
	cache.Set(key, []DockerfileInfo{})

	// Check that cached entry has correct expiration (around DefaultCacheTTL)
	cache.mu.RLock()
	cached := cache.cache[key.String()]
	cache.mu.RUnlock()

	assert.NotNil(t, cached)
	expectedExpiration := time.Now().Add(DefaultCacheTTL)
	// Allow 1 second tolerance
	assert.WithinDuration(t, expectedExpiration, cached.ExpiresAt, 1*time.Second)
}
