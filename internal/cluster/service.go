package cluster

import (
	"strings"
)

// Service provides access to cluster service information
type Service struct{}

// NewService creates a new cluster Service instance
func NewService() *Service {
	return &Service{}
}

// ListServices returns a list of cluster services, optionally filtered by name.
// Filter is case-insensitive substring matching; whitespace is trimmed.
// An empty or whitespace-only filter returns all services.
func (s *Service) ListServices(filter string) ([]ServiceInfo, error) {
	filter = strings.TrimSpace(filter)

	// Mock data representing running cluster services
	services := []ServiceInfo{
		{
			Name:      "api-service",
			Namespace: "default",
			Status:    "Running",
			Endpoints: []string{"10.0.0.1:8080"},
		},
		{
			Name:      "database-service",
			Namespace: "default",
			Status:    "Running",
			Endpoints: []string{"10.0.0.2:5432"},
		},
		{
			Name:      "cache-service",
			Namespace: "default",
			Status:    "Running",
			Endpoints: []string{"10.0.0.3:6379"},
		},
	}

	if filter == "" {
		return services, nil
	}

	// Filter by name using case-insensitive substring matching
	lowerFilter := strings.ToLower(filter)
	filtered := make([]ServiceInfo, 0)
	for _, svc := range services {
		if strings.Contains(strings.ToLower(svc.Name), lowerFilter) {
			filtered = append(filtered, svc)
		}
	}

	return filtered, nil
}
