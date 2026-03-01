package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// TokenBucket represents a token bucket for rate limiting
type TokenBucket struct {
	tokens         float64
	maxTokens      float64
	refillRate     float64 // tokens per second
	lastRefillTime time.Time
	mu             sync.Mutex
}

// NewTokenBucket creates a new token bucket
func NewTokenBucket(maxTokens, refillRate float64) *TokenBucket {
	return &TokenBucket{
		tokens:         maxTokens,
		maxTokens:      maxTokens,
		refillRate:     refillRate,
		lastRefillTime: time.Now(),
	}
}

// TryConsume attempts to consume one token from the bucket
// Returns true if successful, false if insufficient tokens
func (tb *TokenBucket) TryConsume() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// Refill tokens based on elapsed time
	now := time.Now()
	elapsed := now.Sub(tb.lastRefillTime).Seconds()
	tb.tokens += elapsed * tb.refillRate
	if tb.tokens > tb.maxTokens {
		tb.tokens = tb.maxTokens
	}
	tb.lastRefillTime = now

	// Try to consume a token
	if tb.tokens >= 1.0 {
		tb.tokens -= 1.0
		return true
	}
	return false
}

// GetTokens returns the current number of tokens available
func (tb *TokenBucket) GetTokens() float64 {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// Refill tokens before returning
	now := time.Now()
	elapsed := now.Sub(tb.lastRefillTime).Seconds()
	tb.tokens += elapsed * tb.refillRate
	if tb.tokens > tb.maxTokens {
		tb.tokens = tb.maxTokens
	}
	tb.lastRefillTime = now

	return tb.tokens
}

// RateLimiter manages rate limiting for multiple IPs
type RateLimiter struct {
	buckets    sync.Map // map[string]*TokenBucket
	maxTokens  float64
	refillRate float64
}

// NewRateLimiter creates a new rate limiter
// maxRequests: maximum number of requests allowed in the time window
// perMinutes: time window in minutes
func NewRateLimiter(maxRequests int, perMinutes int) *RateLimiter {
	maxTokens := float64(maxRequests)
	refillRate := maxTokens / (float64(perMinutes) * 60.0) // tokens per second

	return &RateLimiter{
		maxTokens:  maxTokens,
		refillRate: refillRate,
	}
}

// GetBucket gets or creates a token bucket for the given IP
func (rl *RateLimiter) GetBucket(ip string) *TokenBucket {
	if bucket, ok := rl.buckets.Load(ip); ok {
		return bucket.(*TokenBucket)
	}

	bucket := NewTokenBucket(rl.maxTokens, rl.refillRate)
	actual, _ := rl.buckets.LoadOrStore(ip, bucket)
	return actual.(*TokenBucket)
}

// RateLimit returns a gin middleware that implements rate limiting
// Default: 500 requests per minute per IP
func RateLimit() gin.HandlerFunc {
	return RateLimitWithConfig(500, 1)
}

// RateLimitWithConfig returns a gin middleware with custom rate limiting configuration
func RateLimitWithConfig(maxRequests int, perMinutes int) gin.HandlerFunc {
	limiter := NewRateLimiter(maxRequests, perMinutes)

	return func(c *gin.Context) {
		// Get client IP
		clientIP := c.ClientIP()

		// Get or create bucket for this IP
		bucket := limiter.GetBucket(clientIP)

		// Try to consume a token
		if !bucket.TryConsume() {
			// Rate limit exceeded
			c.Header("X-RateLimit-Limit", "500")
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", "60")
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			c.Abort()
			return
		}

		// Get current token count for headers
		remaining := int(bucket.GetTokens())
		if remaining < 0 {
			remaining = 0
		}

		// Add rate limit headers
		c.Header("X-RateLimit-Limit", "500")
		c.Header("X-RateLimit-Remaining", string(rune(remaining+'0')))
		c.Header("X-RateLimit-Reset", "60")

		c.Next()
	}
}
