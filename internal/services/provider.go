// Package services provides shared service data for the homelab API.
package services

import "go-github/internal/models"

// GetServices returns the list of mock homelab services.
// Both the HTTP handler and the MCP server consume this shared data.
func GetServices() []models.Service {
	return []models.Service{
		{
			Name:     "homeassistant",
			Type:     "home-automation",
			Status:   "running",
			Endpoint: "http://homeassistant.local:8123",
		},
		{
			Name:     "prometheus",
			Type:     "monitoring",
			Status:   "running",
			Endpoint: "http://prometheus.local:9090",
		},
		{
			Name:     "grafana",
			Type:     "visualization",
			Status:   "running",
			Endpoint: "http://grafana.local:3000",
		},
		{
			Name:     "node-exporter",
			Type:     "metrics",
			Status:   "running",
			Endpoint: "http://node-exporter.local:9100",
		},
		{
			Name:     "alertmanager",
			Type:     "alerting",
			Status:   "running",
			Endpoint: "http://alertmanager.local:9093",
		},
	}
}
