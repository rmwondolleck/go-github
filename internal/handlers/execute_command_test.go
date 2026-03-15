package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestExecuteCommandHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		deviceID           string
		body               map[string]interface{}
		expectedStatusCode int
	}{
		{
			name:     "executes command on controllable device",
			deviceID: "device-001",
			body: map[string]interface{}{
				"action":     "turn_on",
				"parameters": map[string]interface{}{"brightness": 100},
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:     "returns 404 for unknown device",
			deviceID: "unknown-device",
			body: map[string]interface{}{
				"action":     "turn_on",
				"parameters": map[string]interface{}{},
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:     "returns 405 for non-controllable device",
			deviceID: "readonly-sensor-001",
			body: map[string]interface{}{
				"action":     "turn_on",
				"parameters": map[string]interface{}{},
			},
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:               "returns 400 for invalid JSON body",
			deviceID:           "device-001",
			body:               nil,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:     "returns 400 when action is missing",
			deviceID: "device-001",
			body: map[string]interface{}{
				"action":     "",
				"parameters": map[string]interface{}{},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:     "returns 400 when parameters are missing",
			deviceID: "device-001",
			body: map[string]interface{}{
				"action": "turn_on",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			var bodyBytes []byte
			if tt.body != nil {
				var err error
				bodyBytes, err = json.Marshal(tt.body)
				if err != nil {
					t.Fatalf("failed to marshal request body: %v", err)
				}
			} else {
				bodyBytes = []byte("invalid-json{")
			}

			c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/homeassistant/devices/"+tt.deviceID+"/command", bytes.NewReader(bodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Params = gin.Params{{Key: "id", Value: tt.deviceID}}

			ExecuteCommandHandler(c)

			if w.Code != tt.expectedStatusCode {
				t.Errorf("expected status code %d, got %d (body: %s)", tt.expectedStatusCode, w.Code, w.Body.String())
			}
		})
	}
}
