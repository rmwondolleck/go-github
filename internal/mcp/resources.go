package mcp

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"

	"go-github/internal/cluster"
	"go-github/internal/health"
	"go-github/internal/homeassistant"
	"go-github/internal/services"
)

// DevicesResourceHandler returns all smart home devices as a JSON resource.
func DevicesResourceHandler(_ context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	devices := homeassistant.GetDevices()

	// Collect map values into a slice for stable JSON representation.
	deviceSlice := make([]interface{}, 0, len(devices))
	for _, d := range devices {
		deviceSlice = append(deviceSlice, d)
	}

	data, err := json.Marshal(deviceSlice)
	if err != nil {
		return nil, err
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      "homelab://devices",
			MIMEType: "application/json",
			Text:     string(data),
		},
	}, nil
}

// ServicesResourceHandler returns all homelab services as a JSON resource.
func ServicesResourceHandler(_ context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	svcs := services.GetServices()

	data, err := json.Marshal(svcs)
	if err != nil {
		return nil, err
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      "homelab://services",
			MIMEType: "application/json",
			Text:     string(data),
		},
	}, nil
}

// ClusterServicesResourceHandler returns Kubernetes cluster services as a JSON resource.
func ClusterServicesResourceHandler(_ context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	clusterSvc := cluster.NewService()
	clusterServices, err := clusterSvc.ListServices("")
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(clusterServices)
	if err != nil {
		return nil, err
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      "homelab://cluster/services",
			MIMEType: "application/json",
			Text:     string(data),
		},
	}, nil
}

// HealthResourceHandler returns the current health status as a JSON resource.
func HealthResourceHandler(_ context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	checker := health.NewChecker()
	status := checker.Check()

	data, err := json.Marshal(status)
	if err != nil {
		return nil, err
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      "homelab://health",
			MIMEType: "application/json",
			Text:     string(data),
		},
	}, nil
}
