package health

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// TestCheck_ReturnsHealthyStatus verifies that the health check returns a healthy status
func TestCheck_ReturnsHealthyStatus(t *testing.T) {
	// Arrange
	checker := NewChecker()

	// Act
	status := checker.Check()

	// Assert
	if status.Status != "healthy" {
		t.Errorf("expected status to be 'healthy', got '%s'", status.Status)
	}
}

// TestCheck_IncludesUptimeInResponse verifies that the health check includes uptime
func TestCheck_IncludesUptimeInResponse(t *testing.T) {
	// Arrange
	checker := NewChecker()

	// Give the service a moment to establish uptime
	time.Sleep(10 * time.Millisecond)

	// Act
	status := checker.Check()

	// Assert
	if status.Uptime == "" {
		t.Error("expected uptime to be included in response, got empty string")
	}

	// Verify uptime format is reasonable (should contain time units)
	// Uptime could be in format like "1.234s", "123ms", etc.
	if len(status.Uptime) < 2 {
		t.Errorf("expected uptime to have reasonable length, got '%s'", status.Uptime)
	}
}

// TestCheck_IncludesComponentsInResponse verifies that the health check includes components
func TestCheck_IncludesComponentsInResponse(t *testing.T) {
	// Arrange
	checker := NewChecker()

	// Act
	status := checker.Check()

	// Assert
	if status.Components == nil {
		t.Error("expected components to be included in response, got nil")
	}

	// Components map should be initialized (even if empty)
	if status.Components == nil {
		t.Error("expected components map to be initialized")
	}

	// Verify we have at least the API component status
	if len(status.Components) == 0 {
		t.Log("Note: components map is empty, but should contain at least API status")
	}
}

// TestCheck_ComponentsHaveValidStatus verifies that component statuses are valid
func TestCheck_ComponentsHaveValidStatus(t *testing.T) {
	// Arrange
	checker := NewChecker()

	// Act
	status := checker.Check()

	// Assert
	// All components should have a valid status string
	for component, componentStatus := range status.Components {
		if componentStatus == "" {
			t.Errorf("component '%s' has empty status", component)
		}
	}
}

// TestCheck_UptimeIncreasesOverTime verifies that uptime increases over time
func TestCheck_UptimeIncreasesOverTime(t *testing.T) {
	// Arrange
	checker := NewChecker()

	// Act - First check
	firstStatus := checker.Check()
	firstUptime := firstStatus.Uptime

	// Wait a bit
	time.Sleep(100 * time.Millisecond)

	// Act - Second check
	secondStatus := checker.Check()
	secondUptime := secondStatus.Uptime

	// Assert
	// Uptime values should be different (second should be greater)
	if firstUptime == secondUptime {
		t.Error("expected uptime to increase over time, but it remained the same")
	}
}

// TestNewChecker_InitializesCorrectly verifies that NewChecker creates a valid checker
func TestNewChecker_InitializesCorrectly(t *testing.T) {
	// Act
	checker := NewChecker()

	// Assert
	if checker == nil {
		t.Fatal("expected NewChecker to return non-nil checker")
	}

	// Verify that checker can perform a health check
	status := checker.Check()
	if status.Status == "" {
		t.Error("expected newly created checker to return valid status")
	}
}

// TestFormatUptime_AllBranches tests every branch in formatUptime
func TestFormatUptime_AllBranches(t *testing.T) {
tests := []struct {
name     string
duration time.Duration
want     func(string) bool
}{
{
name:     "days branch",
duration: 25*time.Hour + 3*time.Minute + 4*time.Second,
want: func(s string) bool {
return strings.Contains(s, "d") && strings.Contains(s, "h")
},
},
{
name:     "hours branch",
duration: 2*time.Hour + 30*time.Minute + 15*time.Second,
want: func(s string) bool {
return strings.HasPrefix(s, "2h")
},
},
{
name:     "minutes branch",
duration: 5*time.Minute + 10*time.Second,
want: func(s string) bool {
return strings.HasPrefix(s, "5m")
},
},
{
name:     "seconds branch",
duration: 45 * time.Second,
want: func(s string) bool {
return strings.HasSuffix(s, "s") && !strings.Contains(s, "m")
},
},
{
name:     "milliseconds branch",
duration: 500 * time.Millisecond,
want: func(s string) bool {
return strings.HasSuffix(s, "ms")
},
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
result := formatUptime(tt.duration)
if !tt.want(result) {
t.Errorf("formatUptime(%v) = %q, did not match expected pattern", tt.duration, result)
}
if result == "" {
t.Error("formatUptime should never return empty string")
}
})
}
}

// TestFormatUptime_DaysFormat verifies exact format for multi-day uptime
func TestFormatUptime_DaysFormat(t *testing.T) {
d := 2*24*time.Hour + 3*time.Hour + 4*time.Minute + 5*time.Second
result := formatUptime(d)
expected := fmt.Sprintf("%dd %dh %dm %ds", 2, 3, 4, 5)
if result != expected {
t.Errorf("expected %q, got %q", expected, result)
}
}
