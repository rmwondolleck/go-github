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

func TestCheckerReturnsHealthyStatus(t *testing.T) {
	checker := NewChecker()
	ctx := context.Background()

	status := checker.Check(ctx)

	assert.Equal(t, "healthy", status.Status)
	assert.NotEmpty(t, status.Uptime)
	assert.NotNil(t, status.Components)
	assert.Equal(t, "healthy", status.Components["api_server"])
}

func TestCheckerIncludesUptime(t *testing.T) {
	checker := NewChecker()
	
	// Wait a bit to ensure uptime is measurable
	time.Sleep(10 * time.Millisecond)
	
	ctx := context.Background()
	status := checker.Check(ctx)

	assert.NotEmpty(t, status.Uptime)
	assert.Contains(t, status.Uptime, "ms")
}

func TestCheckerUptimeIncreases(t *testing.T) {
	checker := NewChecker()
	ctx := context.Background()

	// First check
	status1 := checker.Check(ctx)
	
	// Wait a bit
	time.Sleep(50 * time.Millisecond)
	
	// Second check
	status2 := checker.Check(ctx)

	// Uptime should be different (and increasing)
	assert.NotEqual(t, status1.Uptime, status2.Uptime)
}
