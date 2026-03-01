package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTokenBucket_TryConsume(t *testing.T) {
	// Arrange
	bucket := NewTokenBucket(10, 1) // 10 tokens, refill 1 per second

	// Act & Assert - consume all tokens
	for i := 0; i < 10; i++ {
		assert.True(t, bucket.TryConsume(), "Should be able to consume token %d", i+1)
	}

	// Should fail on 11th attempt
	assert.False(t, bucket.TryConsume(), "Should not be able to consume 11th token")
}

func TestTokenBucket_Refill(t *testing.T) {
	// Arrange
	bucket := NewTokenBucket(10, 10) // 10 tokens, refill 10 per second

	// Act - consume all tokens
	for i := 0; i < 10; i++ {
		bucket.TryConsume()
	}
	assert.False(t, bucket.TryConsume(), "Should be empty")

	// Wait for refill
	time.Sleep(1100 * time.Millisecond)

	// Assert - should have refilled
	assert.True(t, bucket.TryConsume(), "Should have refilled after waiting")
}

func TestTokenBucket_GetTokens(t *testing.T) {
	// Arrange
	bucket := NewTokenBucket(100, 10)

	// Act - consume some tokens
	for i := 0; i < 5; i++ {
		bucket.TryConsume()
	}

	// Assert
	tokens := bucket.GetTokens()
	assert.Greater(t, tokens, 94.0, "Should have ~95 tokens remaining")
	assert.LessOrEqual(t, tokens, 100.0, "Should not exceed max tokens")
}

func TestRateLimit_AllowsRequestsUnderLimit(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RateLimitWithConfig(10, 1)) // 10 requests per minute

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Act & Assert - make 10 requests (all should succeed)
	for i := 0; i < 10; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Request %d should succeed", i+1)
		assert.NotEmpty(t, w.Header().Get("X-RateLimit-Limit"), "Should have X-RateLimit-Limit header")
		assert.NotEmpty(t, w.Header().Get("X-RateLimit-Remaining"), "Should have X-RateLimit-Remaining header")
		assert.NotEmpty(t, w.Header().Get("X-RateLimit-Reset"), "Should have X-RateLimit-Reset header")
	}
}

func TestRateLimit_BlocksRequestsOverLimit(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RateLimitWithConfig(5, 1)) // 5 requests per minute

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Act - make 5 successful requests
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "Request %d should succeed", i+1)
	}

	// Assert - 6th request should be blocked
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code, "6th request should be rate limited")
	assert.Equal(t, "500", w.Header().Get("X-RateLimit-Limit"))
	assert.Equal(t, "0", w.Header().Get("X-RateLimit-Remaining"))
	assert.Equal(t, "60", w.Header().Get("X-RateLimit-Reset"))
}

func TestRateLimit_PerIPIsolation(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RateLimitWithConfig(3, 1)) // 3 requests per minute

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Act & Assert - IP1 makes 3 requests (all succeed)
	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "IP1 request %d should succeed", i+1)
	}

	// IP1's 4th request should fail
	w1 := httptest.NewRecorder()
	req1 := httptest.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "192.168.1.1:12345"
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusTooManyRequests, w1.Code, "IP1's 4th request should be blocked")

	// IP2 should still be able to make requests
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "192.168.1.2:12345"
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code, "IP2's first request should succeed")
}

func TestRateLimit_HeaderValues(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RateLimit()) // Default: 500 req/min

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "500", w.Header().Get("X-RateLimit-Limit"), "Limit should be 500")
	assert.NotEmpty(t, w.Header().Get("X-RateLimit-Remaining"), "Remaining should be set")
	assert.Equal(t, "60", w.Header().Get("X-RateLimit-Reset"), "Reset should be 60 seconds")
}

func TestRateLimit_ConcurrentRequests(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RateLimitWithConfig(50, 1)) // 50 requests per minute

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Act - make 100 concurrent requests
	var wg sync.WaitGroup
	successCount := 0
	rateLimitCount := 0
	var mu sync.Mutex

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)

			mu.Lock()
			if w.Code == http.StatusOK {
				successCount++
			} else if w.Code == http.StatusTooManyRequests {
				rateLimitCount++
			}
			mu.Unlock()
		}()
	}
	wg.Wait()

	// Assert
	assert.Equal(t, 50, successCount, "Should have exactly 50 successful requests")
	assert.Equal(t, 50, rateLimitCount, "Should have exactly 50 rate-limited requests")
}

func TestRateLimit_RefillOverTime(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	// Use a faster refill rate for testing: 10 requests per minute = refill every 6 seconds
	router.Use(RateLimitWithConfig(10, 1))

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Act - consume all tokens
	for i := 0; i < 10; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// Next request should fail
	w1 := httptest.NewRecorder()
	req1 := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusTooManyRequests, w1.Code)

	// Wait for refill (6 seconds should give us 1 token)
	time.Sleep(6500 * time.Millisecond)

	// Assert - should be able to make one more request
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code, "Should succeed after refill")
}

func TestRateLimit_AbortsMiddlewareChain(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RateLimitWithConfig(1, 1)) // 1 request per minute

	handlerCalled := false
	router.GET("/test", func(c *gin.Context) {
		handlerCalled = true
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Act - first request succeeds
	w1 := httptest.NewRecorder()
	req1 := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w1, req1)
	assert.True(t, handlerCalled, "Handler should be called on first request")

	// Reset flag
	handlerCalled = false

	// Second request should be blocked
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w2, req2)

	// Assert
	assert.Equal(t, http.StatusTooManyRequests, w2.Code)
	assert.False(t, handlerCalled, "Handler should not be called when rate limited")
}

func TestNewRateLimiter(t *testing.T) {
	// Arrange & Act
	limiter := NewRateLimiter(100, 2) // 100 requests per 2 minutes

	// Assert
	assert.Equal(t, 100.0, limiter.maxTokens)
	assert.InDelta(t, 100.0/(2.0*60.0), limiter.refillRate, 0.001, "Refill rate should be ~0.833 tokens/sec")
}

func TestRateLimiter_GetBucket(t *testing.T) {
	// Arrange
	limiter := NewRateLimiter(100, 1)

	// Act
	bucket1 := limiter.GetBucket("192.168.1.1")
	bucket2 := limiter.GetBucket("192.168.1.1")
	bucket3 := limiter.GetBucket("192.168.1.2")

	// Assert
	assert.Same(t, bucket1, bucket2, "Same IP should return same bucket")
	assert.NotSame(t, bucket1, bucket3, "Different IPs should have different buckets")
}

func TestRateLimit_ErrorResponse(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RateLimitWithConfig(1, 1))

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Act - exhaust the limit
	w1 := httptest.NewRecorder()
	req1 := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w1, req1)

	// Make rate-limited request
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w2, req2)

	// Assert
	assert.Equal(t, http.StatusTooManyRequests, w2.Code)
	assert.Contains(t, w2.Body.String(), "Rate limit exceeded", "Response should contain error message")
}

func BenchmarkRateLimit(b *testing.B) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RateLimit())

	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = fmt.Sprintf("192.168.1.%d:12345", i%256)
		router.ServeHTTP(w, req)
	}
}
