package handlers

import (
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

	// In a real implementation, this would fetch devices from a data store
	// For now, we'll return an empty list or sample data
	// The pool will still be used when this handler is called

	resp.Count = len(resp.Devices)
	JSONSuccess(c, http.StatusOK, resp)
}
