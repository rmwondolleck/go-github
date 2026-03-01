package health

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewChecker(t *testing.T) {
	checker := NewChecker()
	assert.NotNil(t, checker)
	assert.False(t, checker.startTime.IsZero())
}

func TestCheck_ReturnsHealthyStatus(t *testing.T) {
	checker := NewChecker()
	ctx := context.Background()

	status := checker.Check(ctx)

	assert.Equal(t, "healthy", status.Status)
	assert.NotEmpty(t, status.Uptime)
	assert.NotNil(t, status.Components)
}

func TestCheck_IncludesAPIServerComponent(t *testing.T) {
	checker := NewChecker()
	ctx := context.Background()

	status := checker.Check(ctx)

	assert.Contains(t, status.Components, "api_server")
	assert.Equal(t, "healthy", status.Components["api_server"])
}

func TestCheck_IncludesUptime(t *testing.T) {
	checker := NewChecker()
	
	// Wait a small amount of time to ensure uptime is measurable
	time.Sleep(10 * time.Millisecond)
	
	ctx := context.Background()
	status := checker.Check(ctx)

	assert.NotEmpty(t, status.Uptime)
	assert.Contains(t, status.Uptime, "s") // Should contain seconds
}

func TestFormatUptime_Seconds(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "zero duration",
			duration: 0,
			expected: "0s",
		},
		{
			name:     "5 seconds",
			duration: 5 * time.Second,
			expected: "5s",
		},
		{
			name:     "30 seconds",
			duration: 30 * time.Second,
			expected: "30s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatUptime(tt.duration)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatUptime_Minutes(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "1 minute",
			duration: 1 * time.Minute,
			expected: "1m 0s",
		},
		{
			name:     "5 minutes 30 seconds",
			duration: 5*time.Minute + 30*time.Second,
			expected: "5m 30s",
		},
		{
			name:     "30 minutes",
			duration: 30 * time.Minute,
			expected: "30m 0s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatUptime(tt.duration)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatUptime_Hours(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "1 hour",
			duration: 1 * time.Hour,
			expected: "1h 0m 0s",
		},
		{
			name:     "2 hours 30 minutes 45 seconds",
			duration: 2*time.Hour + 30*time.Minute + 45*time.Second,
			expected: "2h 30m 45s",
		},
		{
			name:     "23 hours",
			duration: 23 * time.Hour,
			expected: "23h 0m 0s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatUptime(tt.duration)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatUptime_Days(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "1 day",
			duration: 24 * time.Hour,
			expected: "1d 0h 0m 0s",
		},
		{
			name:     "2 days 5 hours 30 minutes 15 seconds",
			duration: 2*24*time.Hour + 5*time.Hour + 30*time.Minute + 15*time.Second,
			expected: "2d 5h 30m 15s",
		},
		{
			name:     "7 days",
			duration: 7 * 24 * time.Hour,
			expected: "7d 0h 0m 0s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatUptime(tt.duration)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCheck_ContextHandling(t *testing.T) {
	checker := NewChecker()
	
	// Test with background context
	status := checker.Check(context.Background())
	assert.Equal(t, "healthy", status.Status)
	
	// Test with custom context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	status = checker.Check(ctx)
	assert.Equal(t, "healthy", status.Status)
}

func TestCheck_ConsistentComponents(t *testing.T) {
	checker := NewChecker()
	ctx := context.Background()

	// Call Check multiple times
	for i := 0; i < 5; i++ {
		status := checker.Check(ctx)
		assert.Equal(t, "healthy", status.Status)
		assert.Contains(t, status.Components, "api_server")
		assert.Equal(t, "healthy", status.Components["api_server"])
	}
}

func TestCheck_UptimeIncreases(t *testing.T) {
	checker := NewChecker()
	ctx := context.Background()

	// First check
	status1 := checker.Check(ctx)

	// Wait enough time to ensure uptime increases
	time.Sleep(1100 * time.Millisecond)

	// Second check
	status2 := checker.Check(ctx)

	// Both should be healthy but have different uptimes
	assert.Equal(t, "healthy", status1.Status)
	assert.Equal(t, "healthy", status2.Status)
	assert.NotEqual(t, status1.Uptime, status2.Uptime)
}
