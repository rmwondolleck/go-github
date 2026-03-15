package mcp

import (
	"context"
	"encoding/json"
	"sync"
	"testing"

	mcpgo "github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// buildToolRequest creates a CallToolRequest with the given arguments map.
func buildToolRequest(args map[string]interface{}) mcpgo.CallToolRequest {
	req := mcpgo.CallToolRequest{}
	req.Params.Name = "execute_command"
	req.Params.Arguments = args
	return req
}

func TestExecuteCommandHandler_TableDriven(t *testing.T) {
	tests := []struct {
		name             string
		args             map[string]interface{}
		wantIsError      bool
		wantTextContains string
	}{
		{
			name: "success: known controllable device",
			args: map[string]interface{}{
				"device_id": "device-001",
				"action":    "turn_on",
			},
			wantIsError:      false,
			wantTextContains: "success",
		},
		{
			name: "not-found: unknown device ID",
			args: map[string]interface{}{
				"device_id": "bad-device-id",
				"action":    "turn_on",
			},
			wantIsError:      true,
			wantTextContains: "device not found",
		},
		{
			name: "non-controllable: read-only sensor",
			args: map[string]interface{}{
				"device_id": "readonly-sensor-001",
				"action":    "turn_on",
			},
			wantIsError:      true,
			wantTextContains: "not controllable",
		},
		{
			name: "missing device_id",
			args: map[string]interface{}{
				"action": "turn_on",
			},
			wantIsError:      true,
			wantTextContains: "device_id is required",
		},
		{
			name: "missing action",
			args: map[string]interface{}{
				"device_id": "device-001",
			},
			wantIsError:      true,
			wantTextContains: "action is required",
		},
		{
			name:             "empty args",
			args:             map[string]interface{}{},
			wantIsError:      true,
			wantTextContains: "device_id is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			req := buildToolRequest(tc.args)

			result, err := ExecuteCommandHandler(ctx, req)
			require.NoError(t, err, "handler itself should not return an error")
			require.NotNil(t, result)

			assert.Equal(t, tc.wantIsError, result.IsError,
				"IsError mismatch for test %q", tc.name)

			require.NotEmpty(t, result.Content, "result content should not be empty")
			tc2, ok := result.Content[0].(mcpgo.TextContent)
			require.True(t, ok, "content should be TextContent")
			assert.Contains(t, tc2.Text, tc.wantTextContains,
				"text should contain %q", tc.wantTextContains)
		})
	}
}

// TestExecuteCommandHandler_SuccessPayload verifies the JSON structure on success.
func TestExecuteCommandHandler_SuccessPayload(t *testing.T) {
	ctx := context.Background()
	req := buildToolRequest(map[string]interface{}{
		"device_id": "device-001",
		"action":    "turn_on",
	})

	result, err := ExecuteCommandHandler(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.False(t, result.IsError)

	tc, ok := result.Content[0].(mcpgo.TextContent)
	require.True(t, ok)

	var payload map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(tc.Text), &payload))
	assert.Equal(t, "success", payload["status"])
	assert.Equal(t, "device-001", payload["device_id"])
	assert.Equal(t, "turn_on", payload["action"])
}

// TestExecuteCommandHandler_ConcurrentRace verifies thread-safety with 10 concurrent calls.
// (T028: race-condition test)
func TestExecuteCommandHandler_ConcurrentRace(t *testing.T) {
	const goroutines = 10
	ctx := context.Background()

	var wg sync.WaitGroup
	wg.Add(goroutines)

	results := make([]*mcpgo.CallToolResult, goroutines)
	errs := make([]error, goroutines)

	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			req := buildToolRequest(map[string]interface{}{
				"device_id": "device-001",
				"action":    "turn_on",
			})
			results[idx], errs[idx] = ExecuteCommandHandler(ctx, req)
		}(i)
	}

	wg.Wait()

	for i := 0; i < goroutines; i++ {
		assert.NoError(t, errs[i], "goroutine %d should not return an error", i)
		assert.NotNil(t, results[i], "goroutine %d should return a non-nil result", i)
	}
}
