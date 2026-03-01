package health

import (
	"context"
	"fmt"
	"time"

	"go-github/internal/models"
)

// Checker handles health check operations
type Checker struct {
	startTime time.Time
}

// NewChecker creates a new health checker instance
func NewChecker() *Checker {
	return &Checker{
		startTime: time.Now(),
	}
}

// Check performs a health check and returns the current health status
func (h *Checker) Check(ctx context.Context) models.HealthStatus {
	uptime := time.Since(h.startTime)
	
	return models.HealthStatus{
		Status: "healthy",
		Uptime: formatUptime(uptime),
		Components: map[string]string{
			"api_server": "healthy",
		},
	}
}

// formatUptime formats a duration into a human-readable string
func formatUptime(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}
