package middleware

import (
	"fmt"
	"sync"
	"time"

	"go-api-starterkit/internal/httpx"

	"github.com/gin-gonic/gin"
)

type rateLimitEntry struct {
	count     int
	expiresAt time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	limit   int
	window  time.Duration
	entries map[string]rateLimitEntry
	now     func() time.Time
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	if limit < 1 {
		limit = 1
	}
	if window <= 0 {
		window = time.Minute
	}

	return &RateLimiter{
		limit:   limit,
		window:  window,
		entries: make(map[string]rateLimitEntry),
		now:     time.Now,
	}
}

func (r *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("%s:%s", c.FullPath(), c.ClientIP())
		if !r.allow(key) {
			httpx.Error(c, 429, "too many requests")
			c.Abort()
			return
		}

		c.Next()
	}
}

func (r *RateLimiter) allow(key string) bool {
	allowed, _ := r.allowAt(key)
	return allowed
}

func (r *RateLimiter) allowAt(key string) (bool, time.Time) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := r.now()
	r.cleanupExpired(now)

	entry, ok := r.entries[key]
	if !ok || now.After(entry.expiresAt) {
		r.entries[key] = rateLimitEntry{
			count:     1,
			expiresAt: now.Add(r.window),
		}
		return true, r.entries[key].expiresAt
	}

	if entry.count >= r.limit {
		return false, entry.expiresAt
	}

	entry.count++
	r.entries[key] = entry
	return true, entry.expiresAt
}

func (r *RateLimiter) cleanupExpired(now time.Time) {
	for key, entry := range r.entries {
		if now.After(entry.expiresAt) {
			delete(r.entries, key)
		}
	}
}
