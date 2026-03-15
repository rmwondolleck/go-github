package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	gin.SetMode(gin.TestMode)

	srv := New()
	assert.NotNil(t, srv)
	assert.NotNil(t, srv.Router())
}

func TestHealthEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	srv := New()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	srv.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}

func TestAPIv1Endpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	srv := New()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1", nil)
	srv.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "API v1")
}

func TestRateLimitAppliedToAPIv1Routes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	srv := New()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1", nil)
	srv.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Header().Get("X-RateLimit-Limit"), "X-RateLimit-Limit header should be set on /api/v1 routes")
	assert.NotEmpty(t, w.Header().Get("X-RateLimit-Remaining"), "X-RateLimit-Remaining header should be set on /api/v1 routes")
}

func TestRateLimitNotAppliedToHealth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	srv := New()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	srv.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Header().Get("X-RateLimit-Limit"), "X-RateLimit-Limit header should not be set on /health")
}

func TestRateLimitNotAppliedToAPIDocs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	srv := New()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/docs/index.html", nil)
	srv.Router().ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusTooManyRequests, w.Code, "/api/docs routes should not be rate limited")
	assert.Empty(t, w.Header().Get("X-RateLimit-Limit"), "X-RateLimit-Limit header should not be set on /api/docs routes")
}

func TestGracefulShutdown(t *testing.T) {
	gin.SetMode(gin.TestMode)

	srv := New()

	// Start server in goroutine with random port
	go func() {
		_ = srv.Run("0")
	}()

	// Brief pause to ensure httpServer is initialized
	// (not waiting for full server startup, just for srv.httpServer to be set)
	time.Sleep(100 * time.Millisecond)

	// Test shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := srv.GracefulShutdown(ctx)
	assert.NoError(t, err)
}

func TestGracefulShutdownWithoutRun(t *testing.T) {
	gin.SetMode(gin.TestMode)

	srv := New()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := srv.GracefulShutdown(ctx)
	assert.NoError(t, err)
}
