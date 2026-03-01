package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(RequestID())

	router.GET("/test", func(c *gin.Context) {
		requestID, exists := c.Get(RequestIDKey)
		assert.True(t, exists, "Request ID should be set in context")
		assert.NotEmpty(t, requestID, "Request ID should not be empty")
		c.String(http.StatusOK, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Header().Get("X-Request-ID"), "X-Request-ID header should be set")
}

func TestRequestIDIsUnique(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(RequestID())

	var id1, id2 string
	router.GET("/test", func(c *gin.Context) {
		requestID, _ := c.Get(RequestIDKey)
		c.String(http.StatusOK, requestID.(string))
	})

	// First request
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w1, req1)
	id1 = w1.Body.String()

	// Second request
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w2, req2)
	id2 = w2.Body.String()

	assert.NotEqual(t, id1, id2, "Request IDs should be unique")
}
