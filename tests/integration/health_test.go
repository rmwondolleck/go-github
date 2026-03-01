package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-github/internal/models"
	"go-github/internal/server"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHealthEndpoint_Returns200 verifies that the health endpoint returns HTTP 200 status
func TestHealthEndpoint_Returns200(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create server instance
	srv := server.New()

	// Create test request
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/health", nil)
	require.NoError(t, err)

	// Execute request
	srv.Router().ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code, "Health endpoint should return 200 OK")
}

// TestHealthEndpoint_IncludesStatusAndUptime verifies the response includes required fields
func TestHealthEndpoint_IncludesStatusAndUptime(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create server instance
	srv := server.New()

	// Create test request
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/health", nil)
	require.NoError(t, err)

	// Execute request
	srv.Router().ServeHTTP(w, req)

	// Assert response status
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response body
	var health models.HealthStatus
	err = json.Unmarshal(w.Body.Bytes(), &health)
	require.NoError(t, err, "Response should be valid JSON matching HealthStatus model")

	// Verify required fields
	assert.NotEmpty(t, health.Status, "Health response must include 'status' field")
	assert.NotEmpty(t, health.Uptime, "Health response must include 'uptime' field")

	// Verify status value is valid
	assert.Equal(t, "ok", health.Status, "Status should be 'ok' when service is healthy")

	// Verify uptime is a valid duration string and non-zero
	duration, err := time.ParseDuration(health.Uptime)
	require.NoError(t, err, "Uptime should be a valid duration string")
	assert.Greater(t, duration, time.Duration(0), "Uptime should be greater than zero")
}

// TestHealthEndpoint_ResponseUnder50ms verifies response time is under 50ms
func TestHealthEndpoint_ResponseUnder50ms(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create server instance
	srv := server.New()

	// Create test request
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/health", nil)
	require.NoError(t, err)

	// Measure response time
	start := time.Now()
	srv.Router().ServeHTTP(w, req)
	duration := time.Since(start)

	// Assert response status
	assert.Equal(t, http.StatusOK, w.Code)

	// Assert response time
	assert.Less(t, duration.Milliseconds(), int64(50), 
		"Health endpoint should respond in under 50ms, got %v", duration)
}
