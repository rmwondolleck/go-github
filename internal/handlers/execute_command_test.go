package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestExecuteCommandHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := `{"action":"turn_on","parameters":{"brightness":100}}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/homeassistant/devices/device-001/command", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "device-001"}}

	ExecuteCommandHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "device-001", response["device_id"])
	assert.Equal(t, "turn_on", response["action"])
}

func TestExecuteCommandHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/homeassistant/devices/device-001/command", bytes.NewBufferString("not-json"))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "device-001"}}

	ExecuteCommandHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestExecuteCommandHandler_InvalidCommand(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Missing required "action" field
	body := `{"parameters":{"brightness":100}}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/homeassistant/devices/device-001/command", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "device-001"}}

	ExecuteCommandHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestExecuteCommandHandler_DeviceNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := `{"action":"turn_on","parameters":{}}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/homeassistant/devices/nonexistent-device/command", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "nonexistent-device"}}

	ExecuteCommandHandler(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestExecuteCommandHandler_NonControllableDevice(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := `{"action":"turn_on","parameters":{}}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/homeassistant/devices/readonly-sensor-001/command", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "readonly-sensor-001"}}

	ExecuteCommandHandler(c)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}
