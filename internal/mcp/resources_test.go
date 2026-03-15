package mcp

import (
	"context"
	"encoding/json"
	"testing"

	mcpgo "github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDevicesResourceHandler(t *testing.T) {
	ctx := context.Background()
	req := mcpgo.ReadResourceRequest{}
	req.Params.URI = "homelab://devices"

	contents, err := DevicesResourceHandler(ctx, req)
	require.NoError(t, err)
	require.Len(t, contents, 1)

	tc, ok := contents[0].(mcpgo.TextResourceContents)
	require.True(t, ok, "should return TextResourceContents")
	assert.Equal(t, "application/json", tc.MIMEType)
	assert.NotEmpty(t, tc.Text)

	// Verify it's valid JSON.
	var raw json.RawMessage
	assert.NoError(t, json.Unmarshal([]byte(tc.Text), &raw))
}

func TestServicesResourceHandler(t *testing.T) {
	ctx := context.Background()
	req := mcpgo.ReadResourceRequest{}
	req.Params.URI = "homelab://services"

	contents, err := ServicesResourceHandler(ctx, req)
	require.NoError(t, err)
	require.Len(t, contents, 1)

	tc, ok := contents[0].(mcpgo.TextResourceContents)
	require.True(t, ok)
	assert.Equal(t, "application/json", tc.MIMEType)
	assert.NotEmpty(t, tc.Text)

	var raw json.RawMessage
	assert.NoError(t, json.Unmarshal([]byte(tc.Text), &raw))
}

func TestClusterServicesResourceHandler(t *testing.T) {
	ctx := context.Background()
	req := mcpgo.ReadResourceRequest{}
	req.Params.URI = "homelab://cluster/services"

	contents, err := ClusterServicesResourceHandler(ctx, req)
	require.NoError(t, err)
	require.Len(t, contents, 1)

	tc, ok := contents[0].(mcpgo.TextResourceContents)
	require.True(t, ok)
	assert.Equal(t, "application/json", tc.MIMEType)
	assert.NotEmpty(t, tc.Text)

	var raw json.RawMessage
	assert.NoError(t, json.Unmarshal([]byte(tc.Text), &raw))
}

func TestHealthResourceHandler(t *testing.T) {
	ctx := context.Background()
	req := mcpgo.ReadResourceRequest{}
	req.Params.URI = "homelab://health"

	contents, err := HealthResourceHandler(ctx, req)
	require.NoError(t, err)
	require.Len(t, contents, 1)

	tc, ok := contents[0].(mcpgo.TextResourceContents)
	require.True(t, ok)
	assert.Equal(t, "application/json", tc.MIMEType)
	assert.NotEmpty(t, tc.Text)

	// Health payload must contain "status" key.
	var m map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(tc.Text), &m))
	_, hasStatus := m["status"]
	assert.True(t, hasStatus, "health response should contain 'status' key")
}

func TestAllResourceHandlers_TableDriven(t *testing.T) {
	tests := []struct {
		name    string
		handler func(context.Context, mcpgo.ReadResourceRequest) ([]mcpgo.ResourceContents, error)
		uri     string
	}{
		{"devices", DevicesResourceHandler, "homelab://devices"},
		{"services", ServicesResourceHandler, "homelab://services"},
		{"cluster/services", ClusterServicesResourceHandler, "homelab://cluster/services"},
		{"health", HealthResourceHandler, "homelab://health"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			req := mcpgo.ReadResourceRequest{}
			req.Params.URI = tc.uri

			contents, err := tc.handler(ctx, req)
			require.NoError(t, err, "handler should not return an error")
			require.Len(t, contents, 1, "should return exactly one ResourceContents item")

			textContents, ok := contents[0].(mcpgo.TextResourceContents)
			require.True(t, ok, "result should be TextResourceContents")
			assert.Equal(t, "application/json", textContents.MIMEType)
			assert.NotEmpty(t, textContents.Text)

			// Validate it's valid JSON.
			var raw json.RawMessage
			assert.NoError(t, json.Unmarshal([]byte(textContents.Text), &raw))
		})
	}
}

// TestResourceHandlers_CancelledContext verifies handlers don't panic on cancelled context.
func TestResourceHandlers_CancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately.

	handlers := []struct {
		name    string
		handler func(context.Context, mcpgo.ReadResourceRequest) ([]mcpgo.ResourceContents, error)
	}{
		{"devices", DevicesResourceHandler},
		{"services", ServicesResourceHandler},
		{"cluster/services", ClusterServicesResourceHandler},
		{"health", HealthResourceHandler},
	}

	for _, tc := range handlers {
		t.Run(tc.name, func(t *testing.T) {
			req := mcpgo.ReadResourceRequest{}
			// These handlers use mock data and don't depend on context being alive,
			// so they should succeed even with a cancelled context.
			assert.NotPanics(t, func() {
				_, _ = tc.handler(ctx, req)
			})
		})
	}
}
