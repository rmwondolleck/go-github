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
	assert.Contains(t, w.Body.String(), "healthy")
	assert.Contains(t, w.Body.String(), "uptime")
	assert.Contains(t, w.Body.String(), "components")
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
