package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequestID(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	t.Run("generates unique request IDs", func(t *testing.T) {
		// Create a new Gin router with the RequestID middleware
		router := gin.New()
		router.Use(RequestID())

		// Track request IDs
		requestIDs := make(map[string]bool)

		// Add a test handler that captures the request ID
		router.GET("/test", func(c *gin.Context) {
			requestID, exists := c.Get("request_id")
			if !exists {
				t.Error("request_id not found in context")
				return
			}

			requestIDStr, ok := requestID.(string)
			if !ok {
				t.Error("request_id is not a string")
				return
			}

			if requestIDStr == "" {
				t.Error("request_id is empty")
				return
			}

			// Check for duplicates
			if requestIDs[requestIDStr] {
				t.Errorf("duplicate request ID: %s", requestIDStr)
			}
			requestIDs[requestIDStr] = true

			c.String(http.StatusOK, "OK")
		})

		// Make multiple requests to ensure unique IDs
		for i := 0; i < 10; i++ {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("expected status 200, got %d", w.Code)
			}
		}

		// Verify we got 10 unique request IDs
		if len(requestIDs) != 10 {
			t.Errorf("expected 10 unique request IDs, got %d", len(requestIDs))
		}
	})

	t.Run("sets X-Request-ID header", func(t *testing.T) {
		router := gin.New()
		router.Use(RequestID())

		router.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "OK")
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Check that X-Request-ID header is set
		requestID := w.Header().Get("X-Request-ID")
		if requestID == "" {
			t.Error("X-Request-ID header not set")
		}

		// Verify it's a valid UUID format (36 characters with dashes)
		if len(requestID) != 36 {
			t.Errorf("invalid UUID format: %s", requestID)
		}
	})

	t.Run("request ID in context matches header", func(t *testing.T) {
		router := gin.New()
		router.Use(RequestID())

		var contextRequestID string

		router.GET("/test", func(c *gin.Context) {
			requestID, exists := c.Get("request_id")
			if !exists {
				t.Error("request_id not found in context")
				return
			}

			contextRequestID = requestID.(string)
			c.String(http.StatusOK, "OK")
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		headerRequestID := w.Header().Get("X-Request-ID")

		if contextRequestID != headerRequestID {
			t.Errorf("context request ID (%s) does not match header (%s)", contextRequestID, headerRequestID)
		}
	})

	t.Run("middleware calls next handler", func(t *testing.T) {
		router := gin.New()
		router.Use(RequestID())

		handlerCalled := false
		router.GET("/test", func(c *gin.Context) {
			handlerCalled = true
			c.String(http.StatusOK, "OK")
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if !handlerCalled {
			t.Error("handler was not called")
		}
	})
}

// BenchmarkRequestID benchmarks the performance of the RequestID middleware
func BenchmarkRequestID(b *testing.B) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(RequestID())

	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
