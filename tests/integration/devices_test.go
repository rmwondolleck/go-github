package integration

import (
	"bytes"
	"encoding/json"
	"go-github/internal/models"
	"go-github/internal/server"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestExecuteCommand_Returns200ForValidCommand tests successful command execution
// This test verifies that a valid command to a controllable device returns HTTP 200
func TestExecuteCommand_Returns200ForValidCommand(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Arrange
	srv := server.New()
	
	// Valid command payload
	command := map[string]interface{}{
		"action": "turn_on",
		"parameters": map[string]interface{}{
			"brightness": 80,
		},
	}
	
	body, err := json.Marshal(command)
	assert.NoError(t, err)
	
	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/homeassistant/devices/device-001/command", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	srv.Router().ServeHTTP(w, req)
	
	// Assert
	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP 200 for valid command")
	
	// Verify JSON response structure
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Response should be valid JSON")
	assert.NotNil(t, response, "Response body should not be nil")
	
	// Response should contain success information
	// The exact structure will be defined when implementing the endpoint
	// For now, we just verify it's valid JSON
}

// TestExecuteCommand_Returns400ForInvalidAction tests invalid action handling
// This test verifies that a command with an invalid or missing action returns HTTP 400
func TestExecuteCommand_Returns400ForInvalidAction(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Arrange
	srv := server.New()
	
	tests := []struct {
		name    string
		command map[string]interface{}
		reason  string
	}{
		{
			name: "missing action",
			command: map[string]interface{}{
				"parameters": map[string]interface{}{
					"brightness": 80,
				},
			},
			reason: "action is required",
		},
		{
			name: "empty action",
			command: map[string]interface{}{
				"action": "",
				"parameters": map[string]interface{}{
					"brightness": 80,
				},
			},
			reason: "action cannot be empty",
		},
		{
			name: "missing parameters",
			command: map[string]interface{}{
				"action": "turn_on",
			},
			reason: "parameters field is required (can be empty map, but not nil)",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			body, err := json.Marshal(tt.command)
			assert.NoError(t, err)
			
			// Act
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/v1/homeassistant/devices/device-001/command", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			srv.Router().ServeHTTP(w, req)
			
			// Assert
			assert.Equal(t, http.StatusBadRequest, w.Code, "Expected HTTP 400 for invalid command: %s", tt.reason)
			
			// Verify error response structure
			var errorResponse models.ErrorResponse
			err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
			assert.NoError(t, err, "Error response should be valid JSON")
			assert.Equal(t, "bad_request", errorResponse.Error, "Error type should be 'bad_request'")
			assert.Equal(t, http.StatusBadRequest, errorResponse.Code, "Error code should be 400")
			assert.NotEmpty(t, errorResponse.Message, "Error message should not be empty")
		})
	}
}

// TestExecuteCommand_Returns405ForReadOnlyDevice tests read-only device handling
// This test verifies that attempting to send a command to a non-controllable device returns HTTP 405
func TestExecuteCommand_Returns405ForReadOnlyDevice(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Arrange
	srv := server.New()
	
	// Valid command payload, but device is read-only
	command := map[string]interface{}{
		"action": "turn_on",
		"parameters": map[string]interface{}{
			"brightness": 80,
		},
	}
	
	body, err := json.Marshal(command)
	assert.NoError(t, err)
	
	// Act
	// Using a device ID that should be read-only (non-controllable)
	// The actual implementation will need to check the device's Controllable field
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/homeassistant/devices/readonly-sensor-001/command", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	srv.Router().ServeHTTP(w, req)
	
	// Assert
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code, "Expected HTTP 405 for read-only device")
	
	// Verify error response structure
	var errorResponse models.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err, "Error response should be valid JSON")
	assert.Equal(t, "method_not_allowed", errorResponse.Error, "Error type should be 'method_not_allowed'")
	assert.Equal(t, http.StatusMethodNotAllowed, errorResponse.Code, "Error code should be 405")
	assert.Contains(t, errorResponse.Message, "not controllable", "Error message should mention device is not controllable")
}
