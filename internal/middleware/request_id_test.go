package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("generates and sets request ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)

		middleware := RequestID()
		middleware(c)

		// Check that request_id was set in context
		requestID, exists := c.Get("request_id")
		if !exists {
			t.Error("request_id not set in context")
		}

		// Check that it's a string
		requestIDStr, ok := requestID.(string)
		if !ok {
			t.Error("request_id is not a string")
		}

		// Check that it's not empty
		if requestIDStr == "" {
			t.Error("request_id is empty")
		}

		// Check that X-Request-ID header was set
		header := w.Header().Get("X-Request-ID")
		if header == "" {
			t.Error("X-Request-ID header not set")
		}

		// Check that header matches context value
		if header != requestIDStr {
			t.Errorf("X-Request-ID header (%s) does not match context value (%s)", header, requestIDStr)
		}
	})

	t.Run("generates unique request IDs", func(t *testing.T) {
		middleware := RequestID()

		// Generate multiple request IDs
		requestIDs := make(map[string]bool)
		for i := 0; i < 100; i++ {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test", nil)

			middleware(c)

			requestID, _ := c.Get("request_id")
			requestIDStr := requestID.(string)

			if requestIDs[requestIDStr] {
				t.Errorf("Duplicate request ID generated: %s", requestIDStr)
			}
			requestIDs[requestIDStr] = true
		}
	})
}
