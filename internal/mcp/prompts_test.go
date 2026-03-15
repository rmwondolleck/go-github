package mcp

import (
	"context"
	"testing"

	mcpgo "github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildPromptRequest(name string, args map[string]string) mcpgo.GetPromptRequest {
	req := mcpgo.GetPromptRequest{}
	req.Params.Name = name
	req.Params.Arguments = args
	return req
}

func TestDeviceControlPromptHandler_TableDriven(t *testing.T) {
	tests := []struct {
		name         string
		args         map[string]string
		wantErr      bool
		wantContains []string
	}{
		{
			name: "success with device name",
			args: map[string]string{"device_name": "Living Room Light"},
			wantContains: []string{
				"Living Room Light",
				"homelab://devices",
				"execute_command",
			},
		},
		{
			name:    "empty device_name returns error",
			args:    map[string]string{"device_name": ""},
			wantErr: true,
		},
		{
			name:    "missing device_name returns error",
			args:    map[string]string{},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			req := buildPromptRequest("device_control", tc.args)

			result, err := DeviceControlPromptHandler(ctx, req)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)
			require.NotEmpty(t, result.Messages)

			msg := result.Messages[0]
			assert.Equal(t, mcpgo.RoleUser, msg.Role)

			tc2, ok := mcpgo.AsTextContent(msg.Content)
			require.True(t, ok, "message content should be TextContent")
			for _, want := range tc.wantContains {
				assert.Contains(t, tc2.Text, want,
					"message text should contain %q", want)
			}
		})
	}
}

func TestServiceStatusPromptHandler_TableDriven(t *testing.T) {
	tests := []struct {
		name         string
		args         map[string]string
		wantErr      bool
		wantContains []string
	}{
		{
			name: "success with service name",
			args: map[string]string{"service_name": "prometheus"},
			wantContains: []string{
				"prometheus",
				"homelab://services",
				"homelab://cluster/services",
			},
		},
		{
			name:    "empty service_name returns error",
			args:    map[string]string{"service_name": ""},
			wantErr: true,
		},
		{
			name:    "missing service_name returns error",
			args:    map[string]string{},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			req := buildPromptRequest("service_status", tc.args)

			result, err := ServiceStatusPromptHandler(ctx, req)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)
			require.NotEmpty(t, result.Messages)

			msg := result.Messages[0]
			assert.Equal(t, mcpgo.RoleUser, msg.Role)

			tc2, ok := mcpgo.AsTextContent(msg.Content)
			require.True(t, ok, "message content should be TextContent")
			for _, want := range tc.wantContains {
				assert.Contains(t, tc2.Text, want,
					"message text should contain %q", want)
			}
		})
	}
}
