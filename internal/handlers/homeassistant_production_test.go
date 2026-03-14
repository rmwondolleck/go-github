package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestDeviceListHandler_ProductionScenario tests the handler in a realistic concurrent scenario
// simulating multiple simultaneous API requests, ensuring the pool works correctly under load
func TestDeviceListHandler_ProductionScenario(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Number of concurrent requests to simulate
	numRequests := 100

	var wg sync.WaitGroup
	errors := make(chan error, numRequests)

	// Launch concurrent requests
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(requestNum int) {
			defer wg.Done()

			// Create a new request and response recorder for each goroutine
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/devices", nil)

			// Call the handler
			DeviceListHandler(c)

			// Verify response
			if w.Code != http.StatusOK {
				errors <- assert.AnError
				return
			}

			// Verify response structure
			var response DeviceListResponse
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				errors <- err
				return
			}

			// Verify response fields
			if response.Devices == nil {
				errors <- assert.AnError
				return
			}
		}(i)
	}

	// Wait for all requests to complete
	wg.Wait()
	close(errors)

	// Check for any errors
	errorCount := 0
	for err := range errors {
		if err != nil {
			errorCount++
			t.Errorf("Request failed: %v", err)
		}
	}

	assert.Equal(t, 0, errorCount, "All concurrent requests should succeed")
}

// TestPoolingBehavior verifies that the pool actually reuses objects
func TestPoolingBehavior(t *testing.T) {
	// Get multiple responses and track if we see the same capacity
	// This indirectly verifies pooling since pooled objects maintain capacity
	capacities := make(map[int]int)

	for i := 0; i < 10; i++ {
		resp := getResponseFromPool()
		capacities[cap(resp.Devices)]++
		putResponseInPool(resp)
	}

	// All should have the same capacity (50) from the pool
	assert.Contains(t, capacities, 50, "Pool should provide objects with capacity 50")
	assert.Equal(t, 1, len(capacities), "All pooled objects should have same capacity")
}
