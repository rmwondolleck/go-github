package handlers

import (
	"go-github/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDeviceListHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		expectedStatus int
	}{
		{
			name:           "returns 200 OK with device list",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/devices", nil)

			DeviceListHandler(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestGetResponseFromPool(t *testing.T) {
	resp := getResponseFromPool()

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Devices)
	assert.Equal(t, 0, len(resp.Devices))
	assert.GreaterOrEqual(t, cap(resp.Devices), 50, "should have capacity for at least 50 devices")

	// Return to pool for cleanup
	putResponseInPool(resp)
}

func TestPutResponseInPool(t *testing.T) {
	resp := getResponseFromPool()

	// Simulate usage
	resp.Count = 10
	resp.Devices = append(resp.Devices, make([]models.Device, 10)...)

	// Return to pool
	putResponseInPool(resp)

	// Get again and verify it was reset
	resp2 := getResponseFromPool()
	assert.Equal(t, 0, resp2.Count)
	assert.Equal(t, 0, len(resp2.Devices))

	// Cleanup
	putResponseInPool(resp2)
}

func TestPoolReuse(t *testing.T) {
	// Get a response from pool
	resp1 := getResponseFromPool()

	// Use and return it
	putResponseInPool(resp1)

	// Get another response - might be the same object (pool reuse)
	resp2 := getResponseFromPool()

	// We can't guarantee it's the same object due to pool internals,
	// but we can verify both are valid
	assert.NotNil(t, resp1)
	assert.NotNil(t, resp2)
	assert.NotNil(t, resp2.Devices)

	// Cleanup
	putResponseInPool(resp2)
}
