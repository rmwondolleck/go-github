package mcp

import (
	"context"
	"testing"

	mcpclient "github.com/mark3labs/mcp-go/client"
	mcpgo "github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestClient creates an in-process MCP client connected to the server
// returned by NewMCPServer. The caller must call Close() on the client when done.
func newTestClient(t *testing.T) (*mcpclient.Client, func()) {
	t.Helper()
	ctx := context.Background()

	s := NewMCPServer()
	c, err := mcpclient.NewInProcessClient(s)
	require.NoError(t, err)

	err = c.Start(ctx)
	require.NoError(t, err)

	initReq := mcpgo.InitializeRequest{}
	initReq.Params.ProtocolVersion = mcpgo.LATEST_PROTOCOL_VERSION
	initReq.Params.ClientInfo = mcpgo.Implementation{Name: "test-client", Version: "0.1"}

	_, err = c.Initialize(ctx, initReq)
	require.NoError(t, err)

	cleanup := func() { _ = c.Close() }
	return c, cleanup
}

// TestNewMCPServer_ReturnsNonNil verifies the constructor returns a usable server.
func TestNewMCPServer_ReturnsNonNil(t *testing.T) {
	s := NewMCPServer()
	assert.NotNil(t, s)
}

// TestNewMCPServer_ReturnsIndependentInstances verifies two calls produce distinct servers.
func TestNewMCPServer_ReturnsIndependentInstances(t *testing.T) {
	s1 := NewMCPServer()
	s2 := NewMCPServer()
	assert.NotNil(t, s1)
	assert.NotNil(t, s2)
	assert.False(t, s1 == s2)
}

// TestInitializeHandshake_ServerName verifies the initialize result carries the right server name.
func TestInitializeHandshake_ServerName(t *testing.T) {
	ctx := context.Background()
	s := NewMCPServer()
	c, err := mcpclient.NewInProcessClient(s)
	require.NoError(t, err)
	require.NoError(t, c.Start(ctx))
	defer c.Close()

	initReq := mcpgo.InitializeRequest{}
	initReq.Params.ProtocolVersion = mcpgo.LATEST_PROTOCOL_VERSION
	initReq.Params.ClientInfo = mcpgo.Implementation{Name: "test", Version: "0.1"}

	result, err := c.Initialize(ctx, initReq)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, serverName, result.ServerInfo.Name)
	assert.Equal(t, serverVersion, result.ServerInfo.Version)

	// Capabilities should be present for resources, tools, and prompts.
	assert.NotNil(t, result.Capabilities.Resources)
	assert.NotNil(t, result.Capabilities.Tools)
	assert.NotNil(t, result.Capabilities.Prompts)
}

// TestResourcesList_ReturnsFourResources verifies resources/list returns exactly 4 resources.
func TestResourcesList_ReturnsFourResources(t *testing.T) {
	ctx := context.Background()
	c, cleanup := newTestClient(t)
	defer cleanup()

	result, err := c.ListResources(ctx, mcpgo.ListResourcesRequest{})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.Resources, 4)

	uris := make([]string, len(result.Resources))
	for i, r := range result.Resources {
		uris[i] = r.URI
	}
	assert.Contains(t, uris, "homelab://devices")
	assert.Contains(t, uris, "homelab://services")
	assert.Contains(t, uris, "homelab://cluster/services")
	assert.Contains(t, uris, "homelab://health")
}

// TestToolsList_ContainsExecuteCommand verifies execute_command tool is registered.
func TestToolsList_ContainsExecuteCommand(t *testing.T) {
	ctx := context.Background()
	c, cleanup := newTestClient(t)
	defer cleanup()

	result, err := c.ListTools(ctx, mcpgo.ListToolsRequest{})
	require.NoError(t, err)
	require.NotNil(t, result)

	var found *mcpgo.Tool
	for i := range result.Tools {
		if result.Tools[i].Name == "execute_command" {
			found = &result.Tools[i]
			break
		}
	}
	require.NotNil(t, found, "execute_command tool should be registered")
	assert.Contains(t, found.InputSchema.Required, "device_id")
	assert.Contains(t, found.InputSchema.Required, "action")
}

// TestPromptsList_ContainsBothPrompts verifies both prompt templates are registered.
func TestPromptsList_ContainsBothPrompts(t *testing.T) {
	ctx := context.Background()
	c, cleanup := newTestClient(t)
	defer cleanup()

	result, err := c.ListPrompts(ctx, mcpgo.ListPromptsRequest{})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.Prompts, 2)

	names := make([]string, len(result.Prompts))
	for i, p := range result.Prompts {
		names[i] = p.Name
	}
	assert.Contains(t, names, "device_control")
	assert.Contains(t, names, "service_status")
}

// TestPromptsList_DeviceControlHasRequiredArg verifies device_control has device_name arg.
func TestPromptsList_DeviceControlHasRequiredArg(t *testing.T) {
	ctx := context.Background()
	c, cleanup := newTestClient(t)
	defer cleanup()

	result, err := c.ListPrompts(ctx, mcpgo.ListPromptsRequest{})
	require.NoError(t, err)

	for _, p := range result.Prompts {
		if p.Name == "device_control" {
			require.NotEmpty(t, p.Arguments)
			assert.Equal(t, "device_name", p.Arguments[0].Name)
			assert.True(t, p.Arguments[0].Required)
			return
		}
	}
	t.Fatal("device_control prompt not found")
}

// TestPromptsList_ServiceStatusHasRequiredArg verifies service_status has service_name arg.
func TestPromptsList_ServiceStatusHasRequiredArg(t *testing.T) {
	ctx := context.Background()
	c, cleanup := newTestClient(t)
	defer cleanup()

	result, err := c.ListPrompts(ctx, mcpgo.ListPromptsRequest{})
	require.NoError(t, err)

	for _, p := range result.Prompts {
		if p.Name == "service_status" {
			require.NotEmpty(t, p.Arguments)
			assert.Equal(t, "service_name", p.Arguments[0].Name)
			assert.True(t, p.Arguments[0].Required)
			return
		}
	}
	t.Fatal("service_status prompt not found")
}

// TestGetPrompt_UnknownPromptReturnsError verifies unknown prompt names produce an error.
// (T033: unknown-prompt test)
func TestGetPrompt_UnknownPromptReturnsError(t *testing.T) {
	ctx := context.Background()
	c, cleanup := newTestClient(t)
	defer cleanup()

	req := mcpgo.GetPromptRequest{}
	req.Params.Name = "nonexistent_prompt"

	_, err := c.GetPrompt(ctx, req)
	assert.Error(t, err, "requesting a nonexistent prompt should return an error")
}
