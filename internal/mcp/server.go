// Package mcp implements the Model Context Protocol server for the homelab API.
// It exposes home assistant devices, services, cluster state, and health data
// to AI assistants such as VS Code Copilot and JetBrains AI via the MCP stdio
// transport.
package mcp

import (
	"context"
	"log/slog"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	serverName    = "go-github-homelab"
	serverVersion = "1.0.0"
)

// NewMCPServer constructs and returns a fully-configured *server.MCPServer.
// All resources, tools, and prompt stubs are registered here.
func NewMCPServer() *server.MCPServer {
	s := server.NewMCPServer(
		serverName,
		serverVersion,
		server.WithResourceCapabilities(true, true),
		server.WithToolCapabilities(true),
		server.WithPromptCapabilities(true),
	)

	registerResources(s)
	registerTools(s)
	registerPrompts(s)

	return s
}

// Run starts the MCP stdio server and blocks until ctx is cancelled or an I/O
// error occurs. It is designed to be launched as a goroutine alongside the HTTP
// server.
func Run(ctx context.Context) error {
	slog.Info("mcp server started", "transport", "stdio")

	mcpServer := NewMCPServer()
	stdioServer := server.NewStdioServer(mcpServer)

	return stdioServer.Listen(ctx, os.Stdin, os.Stdout)
}

// registerResources registers the four homelab resource endpoints.
func registerResources(s *server.MCPServer) {
	s.AddResource(
		mcp.NewResource("homelab://devices", "Homelab Devices",
			mcp.WithResourceDescription("All smart home devices managed by Home Assistant"),
			mcp.WithMIMEType("application/json"),
		),
		DevicesResourceHandler,
	)
	s.AddResource(
		mcp.NewResource("homelab://services", "Homelab Services",
			mcp.WithResourceDescription("All homelab services (prometheus, grafana, etc.)"),
			mcp.WithMIMEType("application/json"),
		),
		ServicesResourceHandler,
	)
	s.AddResource(
		mcp.NewResource("homelab://cluster/services", "Cluster Services",
			mcp.WithResourceDescription("Kubernetes cluster services and their endpoints"),
			mcp.WithMIMEType("application/json"),
		),
		ClusterServicesResourceHandler,
	)
	s.AddResource(
		mcp.NewResource("homelab://health", "Health Status",
			mcp.WithResourceDescription("Current health status and uptime of the homelab API"),
			mcp.WithMIMEType("application/json"),
		),
		HealthResourceHandler,
	)
}

// registerTools registers the execute_command tool.
func registerTools(s *server.MCPServer) {
	executeCommandTool := mcp.NewTool(
		"execute_command",
		mcp.WithDescription("Execute a control command on a Home Assistant device"),
		mcp.WithString("device_id",
			mcp.Required(),
			mcp.Description("The unique identifier of the target device"),
		),
		mcp.WithString("action",
			mcp.Required(),
			mcp.Description("The action to perform (e.g. turn_on, turn_off, set_brightness)"),
		),
		mcp.WithObject("parameters",
			mcp.Description("Optional parameters for the action"),
		),
	)
	s.AddTool(executeCommandTool, ExecuteCommandHandler)
}

// registerPrompts registers the device_control and service_status prompt templates.
func registerPrompts(s *server.MCPServer) {
	deviceControlPrompt := mcp.NewPrompt(
		"device_control",
		mcp.WithPromptDescription("Generate a prompt to control a specific smart home device"),
		mcp.WithArgument("device_name",
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("The human-readable name of the device to control"),
		),
	)
	s.AddPrompt(deviceControlPrompt, DeviceControlPromptHandler)

	serviceStatusPrompt := mcp.NewPrompt(
		"service_status",
		mcp.WithPromptDescription("Generate a prompt to check the status of a homelab service"),
		mcp.WithArgument("service_name",
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("The name of the service to check (e.g. prometheus, grafana)"),
		),
	)
	s.AddPrompt(serviceStatusPrompt, ServiceStatusPromptHandler)
}
