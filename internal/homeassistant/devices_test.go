package homeassistant

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDevices(t *testing.T) {
	devices := GetDevices()
	assert.NotNil(t, devices, "GetDevices should not return nil")
	assert.Greater(t, len(devices), 0, "GetDevices should return at least one device")
}

func TestGetDevices_ContainsExpectedDevices(t *testing.T) {
	devices := GetDevices()
	_, ok := devices["device-001"]
	assert.True(t, ok, "should contain device-001")
	_, ok = devices["readonly-sensor-001"]
	assert.True(t, ok, "should contain readonly-sensor-001")
}

func TestGetDevice(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		wantOk   bool
		wantName string
	}{
		{
			name:     "existing controllable device",
			id:       "device-001",
			wantOk:   true,
			wantName: "Living Room Light",
		},
		{
			name:     "existing read-only device",
			id:       "readonly-sensor-001",
			wantOk:   true,
			wantName: "Temperature Sensor",
		},
		{
			name:   "unknown device ID",
			id:     "unknown-device-999",
			wantOk: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			device, ok := GetDevice(tc.id)
			assert.Equal(t, tc.wantOk, ok)
			if tc.wantOk {
				require.NotNil(t, device)
				assert.Equal(t, tc.wantName, device.Name)
			} else {
				assert.Nil(t, device)
			}
		})
	}
}

func TestExecuteCommand(t *testing.T) {
	tests := []struct {
		name       string
		deviceID   string
		cmd        Command
		wantErr    error
		wantStatus string
	}{
		{
			name:       "success for controllable device",
			deviceID:   "device-001",
			cmd:        Command{Action: "turn_on", Parameters: map[string]interface{}{}},
			wantErr:    nil,
			wantStatus: "success",
		},
		{
			name:     "error for unknown device",
			deviceID: "bad-device-id",
			cmd:      Command{Action: "turn_on", Parameters: map[string]interface{}{}},
			wantErr:  ErrDeviceNotFound,
		},
		{
			name:     "error for non-controllable device",
			deviceID: "readonly-sensor-001",
			cmd:      Command{Action: "turn_on", Parameters: map[string]interface{}{}},
			wantErr:  ErrDeviceNotControllable,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ExecuteCommand(tc.deviceID, tc.cmd)
			if tc.wantErr != nil {
				assert.True(t, errors.Is(err, tc.wantErr), "expected error %v, got %v", tc.wantErr, err)
				assert.Empty(t, result.Status)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.wantStatus, result.Status)
				assert.Equal(t, tc.deviceID, result.DeviceID)
				assert.Equal(t, tc.cmd.Action, result.Action)
			}
		})
	}
}
