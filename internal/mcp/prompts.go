package mcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

// DeviceControlPromptHandler returns a rendered prompt template for controlling
// a specific device. The caller must supply a non-empty device_name argument.
func DeviceControlPromptHandler(_ context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	deviceName := req.Params.Arguments["device_name"]
	if deviceName == "" {
		return nil, errors.New("device_name is required")
	}

	text := fmt.Sprintf(
		"You are controlling the smart home device '%s'.\n"+
			"Use the resource homelab://devices to get the current state of all devices.\n"+
			"Use the execute_command tool to perform actions on '%s'.\n"+
			"Available actions include: turn_on, turn_off, set_brightness.\n"+
			"Always confirm the current state before executing a command.",
		deviceName, deviceName,
	)

	return mcp.NewGetPromptResult(
		"Device control prompt for "+deviceName,
		[]mcp.PromptMessage{
			mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(text)),
		},
	), nil
}

// ServiceStatusPromptHandler returns a rendered prompt template for checking
// the status of a specific homelab service.
func ServiceStatusPromptHandler(_ context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	serviceName := req.Params.Arguments["service_name"]
	if serviceName == "" {
		return nil, errors.New("service_name is required")
	}

	text := fmt.Sprintf(
		"You are checking the status of the homelab service '%s'.\n"+
			"Use the resource homelab://services to get the list of all homelab services.\n"+
			"Use the resource homelab://cluster/services to get the Kubernetes cluster service details.\n"+
			"Report whether '%s' is running, its endpoint, and any relevant health information.",
		serviceName, serviceName,
	)

	return mcp.NewGetPromptResult(
		"Service status prompt for "+serviceName,
		[]mcp.PromptMessage{
			mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(text)),
		},
	), nil
}
