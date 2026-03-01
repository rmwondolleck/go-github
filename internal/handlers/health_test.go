package handlers

import (
	"encoding/json"
	"go-github/internal/health"
	"go-github/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	checker := health.NewChecker()
	router := gin.New()
	router.GET("/health", HealthHandler(checker))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	
	// Measure response time
	start := time.Now()
	router.ServeHTTP(w, req)
	duration := time.Since(start)

	// Assert status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Assert response structure
	var response models.HealthStatus
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", response.Status)
	assert.NotEmpty(t, response.Uptime)
	assert.NotNil(t, response.Components)
	assert.Equal(t, "healthy", response.Components["api_server"])

	// Assert response time is under 50ms
	assert.Less(t, duration.Milliseconds(), int64(50), "Response time should be less than 50ms")
}

func TestHealthHandlerResponseTime(t *testing.T) {
	gin.SetMode(gin.TestMode)

	checker := health.NewChecker()
	router := gin.New()
	router.GET("/health", HealthHandler(checker))

	// Run multiple times to ensure consistent performance
	for i := 0; i < 10; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health", nil)
		
		start := time.Now()
		router.ServeHTTP(w, req)
		duration := time.Since(start)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Less(t, duration.Milliseconds(), int64(50), "Response time should be less than 50ms in iteration %d", i)
	}
}
