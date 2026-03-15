package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"

	"go-github/internal/homeassistant"
)

// ExecuteCommandHandler handles the execute_command MCP tool call.
// It extracts device_id and action from the request arguments, executes the
// command via the homeassistant package, and returns a structured result.
func ExecuteCommandHandler(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()

	deviceID, _ := args["device_id"].(string)
	if deviceID == "" {
		return mcp.NewToolResultError("device_id is required"), nil
	}

	action, _ := args["action"].(string)
	if action == "" {
		return mcp.NewToolResultError("action is required"), nil
	}

	// Extract optional parameters map.
	var params map[string]interface{}
	if raw, ok := args["parameters"]; ok && raw != nil {
		if m, ok := raw.(map[string]interface{}); ok {
			params = m
		}
	}
	if params == nil {
		params = map[string]interface{}{}
	}

	cmd := homeassistant.Command{
		Action:     action,
		Parameters: params,
	}

	result, err := homeassistant.ExecuteCommand(deviceID, cmd)
	if err != nil {
		switch {
		case errors.Is(err, homeassistant.ErrDeviceNotFound):
			return mcp.NewToolResultError(fmt.Sprintf("device not found: %s", deviceID)), nil
		case errors.Is(err, homeassistant.ErrDeviceNotControllable):
			return mcp.NewToolResultError(fmt.Sprintf("device is not controllable: %s", deviceID)), nil
		default:
			return mcp.NewToolResultError(err.Error()), nil
		}
	}

	// Marshal CommandResult to JSON text.
	data, err := json.Marshal(result)
	if err != nil {
		return mcp.NewToolResultError("failed to marshal result: " + err.Error()), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}
