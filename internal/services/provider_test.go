package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetServices(t *testing.T) {
	services := GetServices()
	assert.NotEmpty(t, services, "GetServices should return a non-empty slice")
}

func TestGetServices_ContainsPrometheus(t *testing.T) {
	services := GetServices()
	found := false
	for _, svc := range services {
		if svc.Name == "prometheus" {
			found = true
			break
		}
	}
	assert.True(t, found, "GetServices should contain an entry with Name == 'prometheus'")
}

func TestGetServices_AllHaveNonEmptyStatus(t *testing.T) {
	services := GetServices()
	for _, svc := range services {
		assert.NotEmpty(t, svc.Status, "service %q should have a non-empty Status", svc.Name)
	}
}

func TestGetServices_ExpectedCount(t *testing.T) {
	services := GetServices()
	assert.Equal(t, 5, len(services), "should return exactly 5 services")
}

func TestGetServices_AllHaveNonEmptyNames(t *testing.T) {
	services := GetServices()
	for _, svc := range services {
		assert.NotEmpty(t, svc.Name, "service should have a non-empty Name")
	}
}
