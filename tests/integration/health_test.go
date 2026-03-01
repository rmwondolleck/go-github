package integration

import (
	"encoding/json"
	"go-github/internal/models"
	"go-github/internal/server"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthEndpointIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create server instance
	srv := server.New()

	// Test: GET /health returns 200 with JSON health status
	t.Run("returns 200 with health status", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health", nil)
		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	})

	// Test: Response includes status, components, uptime
	t.Run("includes status, components, and uptime", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health", nil)
		srv.Router().ServeHTTP(w, req)

		var response models.HealthStatus
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "healthy", response.Status)
		assert.NotEmpty(t, response.Uptime)
		assert.NotNil(t, response.Components)
		assert.Contains(t, response.Components, "api_server")
		assert.Equal(t, "healthy", response.Components["api_server"])
	})

	// Test: Response time <50ms (performance validation)
	t.Run("response time under 50ms", func(t *testing.T) {
		// Run multiple iterations to ensure consistent performance
		for i := 0; i < 10; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/health", nil)

			start := time.Now()
			srv.Router().ServeHTTP(w, req)
			duration := time.Since(start)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Less(t, duration.Milliseconds(), int64(50), 
				"Response time should be under 50ms (iteration %d: %dms)", i, duration.Milliseconds())
		}
	})

	// Test: Health endpoint is idempotent
	t.Run("is idempotent", func(t *testing.T) {
		// Make multiple requests and ensure they all succeed
		for i := 0; i < 5; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/health", nil)
			srv.Router().ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			
			var response models.HealthStatus
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "healthy", response.Status)
		}
	})
}
