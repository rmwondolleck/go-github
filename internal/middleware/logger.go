package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger logs HTTP requests with request ID, method, path, status, and duration
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()
		requestID, _ := c.Get(RequestIDKey)

		slog.Info("request completed",
			"request_id", requestID,
			"method", method,
			"path", path,
			"status", status,
			"duration", duration,
		)
	}
}
