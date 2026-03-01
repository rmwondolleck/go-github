package middleware

import (
	"log/slog"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go-github/internal/models"
)

// Recovery middleware recovers from panics and returns appropriate error responses
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID, _ := c.Get("request_id")
				slog.Error("panic recovered",
					"request_id", requestID,
					"error", err,
					"stack", string(debug.Stack()),
				)
				c.AbortWithStatusJSON(500, models.ErrorResponse{
					Error:   "internal_server_error",
					Message: "An unexpected error occurred",
					Code:    500,
				})
			}
		}()
		c.Next()
	}
}
