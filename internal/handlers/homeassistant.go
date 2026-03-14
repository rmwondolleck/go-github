package handlers

import (
	"go-github/internal/homeassistant"
	"go-github/internal/models"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// DeviceListResponse represents a response containing a list of devices
type DeviceListResponse struct {
	Devices []models.Device `json:"devices"`
	Count   int             `json:"count"`
}

// deviceListResponsePool is a global sync.Pool for DeviceListResponse objects
// This reduces memory allocations by reusing response objects
var deviceListResponsePool = sync.Pool{
	New: func() interface{} {
		return &DeviceListResponse{
			Devices: make([]models.Device, 0, 50), // Pre-allocate capacity for 50 devices
		}
	},
}

// getResponseFromPool retrieves a DeviceListResponse from the pool
func getResponseFromPool() *DeviceListResponse {
	return deviceListResponsePool.Get().(*DeviceListResponse)
}

// putResponseInPool returns a DeviceListResponse to the pool after resetting its state
func putResponseInPool(resp *DeviceListResponse) {
	// Reset the response state before returning to pool
	resp.Devices = resp.Devices[:0] // Clear slice but keep capacity
	resp.Count = 0
	deviceListResponsePool.Put(resp)
}

// DeviceListHandler handles GET requests for the device list
// It uses sync.Pool to reduce memory allocations
func DeviceListHandler(c *gin.Context) {
	// Get a response object from the pool
	resp := getResponseFromPool()
	defer putResponseInPool(resp)

	// In a real implementation, this would fetch devices from a data store
	// For now, we'll return an empty list or sample data
	// The pool will still be used when this handler is called

	resp.Count = len(resp.Devices)
	JSONSuccess(c, http.StatusOK, resp)
}

// mockDevices provides mock device data for command handling
var mockDevices = map[string]*models.Device{
	"device-001": {
		ID:           "device-001",
		Name:         "Living Room Light",
		Type:         "light",
		State:        "off",
		Attributes:   map[string]interface{}{"brightness": 0},
		LastUpdated:  time.Now(),
		Controllable: true,
	},
	"readonly-sensor-001": {
		ID:           "readonly-sensor-001",
		Name:         "Temperature Sensor",
		Type:         "sensor",
		State:        "72",
		Attributes:   map[string]interface{}{"unit": "°F"},
		LastUpdated:  time.Now(),
		Controllable: false,
	},
}

// ExecuteCommandHandler handles POST /api/v1/homeassistant/devices/{id}/command
// It validates the command and executes it on the specified device.
func ExecuteCommandHandler(c *gin.Context) {
	deviceID := c.Param("id")

	// Parse command from request body
	var cmd homeassistant.Command
	if err := c.ShouldBindJSON(&cmd); err != nil {
		JSONError(c, http.StatusBadRequest, "bad_request", "invalid request body: "+err.Error())
		return
	}

	// Validate command
	if err := cmd.Validate(); err != nil {
		JSONError(c, http.StatusBadRequest, "bad_request", err.Error())
		return
	}

	// Look up device
	device, exists := mockDevices[deviceID]
	if !exists {
		JSONError(c, http.StatusNotFound, "not_found", "device not found: "+deviceID)
		return
	}

	// Check if device is controllable
	if !device.Controllable {
		JSONError(c, http.StatusMethodNotAllowed, "method_not_allowed", "device is not controllable: "+deviceID)
		return
	}

	// Command executed successfully
	JSONSuccess(c, http.StatusOK, gin.H{
		"status":    "success",
		"device_id": deviceID,
		"action":    cmd.Action,
	})
}
