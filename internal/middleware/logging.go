package middleware

import (
"log/slog"
"time"

"github.com/gin-gonic/gin"
)

// Logger returns a gin.HandlerFunc middleware that logs HTTP requests using slog.
// It logs the request method, path, status code, duration in milliseconds,
// and includes the request_id from the context if available.
func Logger() gin.HandlerFunc {
return func(c *gin.Context) {
start := time.Now()
requestID, _ := c.Get("request_id")

c.Next()

duration := time.Since(start)
slog.Info("request completed",
"request_id", requestID,
"method", c.Request.Method,
"path", c.Request.URL.Path,
"status", c.Writer.Status(),
"duration_ms", duration.Milliseconds(),
)
}
}
