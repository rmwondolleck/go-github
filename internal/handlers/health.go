package handlers

import (
	"context"
	"go-github/internal/health"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles the health check endpoint
// @Summary Health check
// @Description Get service health status and uptime
// @Tags health
// @Produce json
// @Success 200 {object} models.HealthStatus
// @Failure 503 {object} models.ErrorResponse
// @Router /health [get]
func HealthHandler(checker *health.Checker) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create context with timeout for health check
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Millisecond)
		defer cancel()

		// Perform health check
		status := checker.Check(ctx)

		// Return 200 OK with health status
		JSONSuccess(c, http.StatusOK, status)
	}
}
