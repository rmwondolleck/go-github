package handlers

import (
	"errors"
	"go-github/internal/homeassistant"
	"go-github/internal/models"
	"net/http"
	"sync"

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

	devices := homeassistant.GetDevices()
	for _, d := range devices {
		resp.Devices = append(resp.Devices, *d)
	}

	resp.Count = len(resp.Devices)
	JSONSuccess(c, http.StatusOK, resp)
}

// ExecuteCommandHandler godoc
// @Summary Execute a device command
// @Description Execute a control command on a HomeAssistant device
// @Tags homeassistant
// @Accept json
// @Produce json
// @Param id path string true "Device ID"
// @Param command body homeassistant.Command true "Command to execute"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 405 {object} models.ErrorResponse
// @Router /api/v1/homeassistant/devices/{id}/command [post]
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

	// Execute command via homeassistant package
	result, err := homeassistant.ExecuteCommand(deviceID, cmd)
	if err != nil {
		switch {
		case errors.Is(err, homeassistant.ErrDeviceNotFound):
			JSONError(c, http.StatusNotFound, "not_found", "device not found: "+deviceID)
		case errors.Is(err, homeassistant.ErrDeviceNotControllable):
			JSONError(c, http.StatusMethodNotAllowed, "method_not_allowed", "device is not controllable: "+deviceID)
		default:
			JSONError(c, http.StatusInternalServerError, "internal_error", err.Error())
		}
		return
	}

	// Command executed successfully
	JSONSuccess(c, http.StatusOK, gin.H{
		"status":    result.Status,
		"device_id": result.DeviceID,
		"action":    result.Action,
	})
}
