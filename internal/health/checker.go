package health

import (
	"context"
	"go-github/internal/models"
	"time"
)

// Checker provides health check functionality
type Checker struct {
	startTime time.Time
}

// NewChecker creates a new health checker instance
func NewChecker() *Checker {
	return &Checker{
		startTime: time.Now(),
	}
}

// Check performs health check and returns status
func (h *Checker) Check(ctx context.Context) models.HealthStatus {
	uptime := time.Since(h.startTime)
	
	return models.HealthStatus{
		Status: "healthy",
		Uptime: uptime.String(),
		Components: map[string]string{
			"api_server": "healthy",
		},
	}
}
