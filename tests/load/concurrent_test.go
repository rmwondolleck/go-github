package load

import (
	"net/http"
	"net/http/httptest"
	"runtime"
	"sort"
	"sync"
	"testing"
	"time"

	"go-github/internal/server"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestConcurrentHealthRequests tests 100 concurrent requests to the health endpoint
func TestConcurrentHealthRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Arrange
	srv := server.New()
	concurrentRequests := 100
	var wg sync.WaitGroup
	responseTimes := make([]time.Duration, concurrentRequests)
	var mu sync.Mutex
	successCount := 0

	// Act - Send 100 concurrent requests
	startTime := time.Now()
	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			reqStart := time.Now()
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/health", nil)
			srv.Router().ServeHTTP(w, req)
			reqDuration := time.Since(reqStart)

			mu.Lock()
			responseTimes[index] = reqDuration
			if w.Code == http.StatusOK {
				successCount++
			}
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	totalDuration := time.Since(startTime)

	// Assert - Verify all requests succeeded
	assert.Equal(t, concurrentRequests, successCount, "All requests should succeed")

	// Calculate and validate p99 response time
	p99 := calculatePercentile(responseTimes, 99)
	t.Logf("Health endpoint - Total duration: %v, P99 response time: %v", totalDuration, p99)
	assert.Less(t, p99, 200*time.Millisecond, "P99 response time should be less than 200ms")

	// Log additional stats
	logResponseStats(t, "Health endpoint", responseTimes)
}

// TestConcurrentAPIv1Requests tests 100 concurrent requests to the API v1 endpoint
func TestConcurrentAPIv1Requests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Arrange
	srv := server.New()
	concurrentRequests := 100
	var wg sync.WaitGroup
	responseTimes := make([]time.Duration, concurrentRequests)
	var mu sync.Mutex
	successCount := 0

	// Act - Send 100 concurrent requests
	startTime := time.Now()
	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			reqStart := time.Now()
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1", nil)
			srv.Router().ServeHTTP(w, req)
			reqDuration := time.Since(reqStart)

			mu.Lock()
			responseTimes[index] = reqDuration
			if w.Code == http.StatusOK {
				successCount++
			}
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	totalDuration := time.Since(startTime)

	// Assert - Verify all requests succeeded
	assert.Equal(t, concurrentRequests, successCount, "All requests should succeed")

	// Calculate and validate p99 response time
	p99 := calculatePercentile(responseTimes, 99)
	t.Logf("API v1 endpoint - Total duration: %v, P99 response time: %v", totalDuration, p99)
	assert.Less(t, p99, 200*time.Millisecond, "P99 response time should be less than 200ms")

	// Log additional stats
	logResponseStats(t, "API v1 endpoint", responseTimes)
}

// TestMemoryLeakDetection tests for memory leaks during concurrent requests
func TestMemoryLeakDetection(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Arrange
	srv := server.New()
	runtime.GC() // Force garbage collection before measuring

	// Measure initial memory
	var initialMemStats runtime.MemStats
	runtime.ReadMemStats(&initialMemStats)
	initialAlloc := initialMemStats.Alloc

	// Act - Send multiple batches of concurrent requests
	batches := 5
	requestsPerBatch := 100

	for batch := 0; batch < batches; batch++ {
		var wg sync.WaitGroup
		for i := 0; i < requestsPerBatch; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				w := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/health", nil)
				srv.Router().ServeHTTP(w, req)
			}()
		}
		wg.Wait()

		// Brief pause between batches
		time.Sleep(10 * time.Millisecond)
	}

	// Force garbage collection to clean up temporary allocations
	runtime.GC()
	time.Sleep(100 * time.Millisecond) // Give GC time to complete

	// Measure final memory
	var finalMemStats runtime.MemStats
	runtime.ReadMemStats(&finalMemStats)
	finalAlloc := finalMemStats.Alloc

	// Assert - Check for excessive memory growth
	memoryGrowth := int64(finalAlloc - initialAlloc)
	memoryGrowthMB := float64(memoryGrowth) / 1024 / 1024

	t.Logf("Memory stats: Initial: %d bytes, Final: %d bytes, Growth: %.2f MB",
		initialAlloc, finalAlloc, memoryGrowthMB)
	t.Logf("Total requests: %d, Memory per request: %.2f bytes",
		batches*requestsPerBatch, float64(memoryGrowth)/float64(batches*requestsPerBatch))

	// Assert that memory growth is reasonable (less than 10 MB for 500 requests)
	// This is a reasonable threshold for detecting significant memory leaks
	assert.Less(t, memoryGrowthMB, 10.0, "Memory growth should be less than 10 MB")
}

// calculatePercentile calculates the specified percentile from a slice of durations
func calculatePercentile(durations []time.Duration, percentile int) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	// Create a copy to avoid modifying the original slice
	sorted := make([]time.Duration, len(durations))
	copy(sorted, durations)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	index := (len(sorted) * percentile) / 100
	if index >= len(sorted) {
		index = len(sorted) - 1
	}

	return sorted[index]
}

// logResponseStats logs detailed response time statistics
func logResponseStats(t *testing.T, endpoint string, durations []time.Duration) {
	if len(durations) == 0 {
		return
	}

	// Calculate various percentiles
	p50 := calculatePercentile(durations, 50)
	p95 := calculatePercentile(durations, 95)
	p99 := calculatePercentile(durations, 99)

	// Calculate min and max
	var min, max time.Duration
	var total time.Duration
	min = durations[0]
	max = durations[0]

	for _, d := range durations {
		total += d
		if d < min {
			min = d
		}
		if d > max {
			max = d
		}
	}

	avg := total / time.Duration(len(durations))

	t.Logf("%s - Response time stats:", endpoint)
	t.Logf("  Min: %v, Max: %v, Avg: %v", min, max, avg)
	t.Logf("  P50: %v, P95: %v, P99: %v", p50, p95, p99)
	t.Logf("  Total requests: %d", len(durations))
}

// BenchmarkConcurrentHealthRequests benchmarks concurrent health endpoint requests
func BenchmarkConcurrentHealthRequests(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)
	srv := server.New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for j := 0; j < 100; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				w := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/health", nil)
				srv.Router().ServeHTTP(w, req)
			}()
		}
		wg.Wait()
	}
}

// BenchmarkConcurrentAPIv1Requests benchmarks concurrent API v1 endpoint requests
func BenchmarkConcurrentAPIv1Requests(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)
	srv := server.New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for j := 0; j < 100; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				w := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/api/v1", nil)
				srv.Router().ServeHTTP(w, req)
			}()
		}
		wg.Wait()
	}
}
