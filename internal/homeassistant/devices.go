package homeassistant

import (
	"errors"
	"go-github/internal/models"
	"time"
)

// ErrDeviceNotFound is returned when the requested device ID does not exist.
var ErrDeviceNotFound = errors.New("device not found")

// ErrDeviceNotControllable is returned when the requested device cannot be controlled.
var ErrDeviceNotControllable = errors.New("device is not controllable")

// CommandResult represents the result of executing a command on a device.
type CommandResult struct {
	Status   string `json:"status"`
	DeviceID string `json:"device_id"`
	Action   string `json:"action"`
}

// mockDevices is the in-memory device store shared by the homeassistant package.
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

// GetDevices returns all devices from the mock device store.
func GetDevices() map[string]*models.Device {
	return mockDevices
}

// GetDevice returns the device with the given ID, or false if not found.
func GetDevice(id string) (*models.Device, bool) {
	device, ok := mockDevices[id]
	return device, ok
}

// ExecuteCommand executes a command on the specified device.
// Returns ErrDeviceNotFound if the device ID is unknown.
// Returns ErrDeviceNotControllable if the device is read-only.
func ExecuteCommand(deviceID string, cmd Command) (CommandResult, error) {
	device, ok := mockDevices[deviceID]
	if !ok {
		return CommandResult{}, ErrDeviceNotFound
	}

	if !device.Controllable {
		return CommandResult{}, ErrDeviceNotControllable
	}

	return CommandResult{
		Status:   "success",
		DeviceID: deviceID,
		Action:   cmd.Action,
	}, nil
}
