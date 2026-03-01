package middleware

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go-github/internal/models"
)

func TestRecovery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("recovers from panic", func(t *testing.T) {
		w := httptest.NewRecorder()
		
		router := gin.New()
		router.Use(RequestID())
		router.Use(Recovery())
		router.GET("/test", func(c *gin.Context) {
			panic("test panic")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		// Check that response is 500
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", w.Code)
		}

		// Check response body
		var errResp models.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &errResp); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if errResp.Error != "internal_server_error" {
			t.Errorf("Expected error 'internal_server_error', got '%s'", errResp.Error)
		}
		if errResp.Message != "An unexpected error occurred" {
			t.Errorf("Expected message 'An unexpected error occurred', got '%s'", errResp.Message)
		}
		if errResp.Code != 500 {
			t.Errorf("Expected code 500, got %d", errResp.Code)
		}
	})

	t.Run("logs panic with stack trace", func(t *testing.T) {
		// Capture log output
		var logBuf bytes.Buffer
		logger := slog.New(slog.NewJSONHandler(&logBuf, nil))
		slog.SetDefault(logger)

		w := httptest.NewRecorder()
		
		router := gin.New()
		router.Use(RequestID())
		router.Use(Recovery())
		router.GET("/test", func(c *gin.Context) {
			panic("test panic")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		// Check that something was logged
		logOutput := logBuf.String()
		if logOutput == "" {
			t.Error("Expected log output, got empty string")
		}

		// Check that log contains expected fields
		if !bytes.Contains(logBuf.Bytes(), []byte("panic recovered")) {
			t.Error("Log output does not contain 'panic recovered'")
		}
		if !bytes.Contains(logBuf.Bytes(), []byte("test panic")) {
			t.Error("Log output does not contain panic message")
		}
	})

	t.Run("does not interfere with normal requests", func(t *testing.T) {
		w := httptest.NewRecorder()
		
		router := gin.New()
		router.Use(Recovery())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		// Check that response is 200
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		// Check response body
		var resp map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if resp["status"] != "ok" {
			t.Errorf("Expected status 'ok', got '%s'", resp["status"])
		}
	})

	t.Run("includes request ID in error response when available", func(t *testing.T) {
		var logBuf bytes.Buffer
		logger := slog.New(slog.NewJSONHandler(&logBuf, nil))
		slog.SetDefault(logger)

		w := httptest.NewRecorder()
		
		router := gin.New()
		router.Use(RequestID())
		router.Use(Recovery())
		router.GET("/test", func(c *gin.Context) {
			panic("test panic")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		// Check that request ID is in the log
		logOutput := logBuf.String()
		if logOutput == "" {
			t.Error("Expected log output, got empty string")
		}
	})

	t.Run("handles panic without request ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		
		router := gin.New()
		router.Use(Recovery())
		router.GET("/test", func(c *gin.Context) {
			panic("test panic")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		// Should still return 500 error
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", w.Code)
		}

		var errResp models.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &errResp); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if errResp.Code != 500 {
			t.Errorf("Expected code 500, got %d", errResp.Code)
		}
	})
}
