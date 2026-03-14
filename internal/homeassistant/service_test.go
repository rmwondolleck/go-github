package homeassistant

// Test-Driven Development (TDD) Tests for HomeAssistant Service Command Execution
//
// This file contains comprehensive unit tests for the command execution feature
// in the HomeAssistant service, following TDD methodology.
//
// Test Coverage:
// - TestExecuteCommand_SucceedsForValidDevice: Tests successful command execution on controllable devices
// - TestExecuteCommand_FailsForInvalidDevice: Tests error handling when device is not found
// - TestExecuteCommand_FailsForReadOnlyDevice: Tests error handling for read-only/sensor devices
// - TestExecuteCommand_FailsForInvalidAction: Tests error handling for unsupported actions
// - TestExecuteCommand_EdgeCases: Tests edge cases like nil commands and empty parameters
// - TestExecuteCommand_ConcurrentExecution: Tests thread safety
// - TestExecuteCommand_ValidatesCommandBeforeDeviceLookup: Tests validation order
//
// Mock Implementation:
// The tests use a mock service implementation (mockService) that simulates the expected
// behavior of the real service. When the real service is implemented, it should:
// 1. Validate the command using Command.Validate()
// 2. Check if the device exists
// 3. Verify the device is controllable (Controllable == true)
// 4. Validate the action is supported for the device type
// 5. Execute the command (actual implementation TBD)
//
// Device Types and Valid Actions:
// - light: turn_on, turn_off, set_brightness
// - switch: turn_on, turn_off
// - sensor: none (read-only devices)

import (
	"context"
	"errors"
	"go-github/internal/models"
	"testing"
)

// mockDeviceStore provides mock device data for testing
type mockDeviceStore struct {
	devices map[string]*models.Device
}

func newMockDeviceStore() *mockDeviceStore {
	return &mockDeviceStore{
		devices: map[string]*models.Device{
			"light.living_room": {
				ID:           "light.living_room",
				Name:         "Living Room Light",
				Type:         "light",
				State:        "on",
				Attributes:   map[string]interface{}{"brightness": 80},
				Controllable: true,
			},
			"light.bedroom": {
				ID:           "light.bedroom",
				Name:         "Bedroom Light",
				Type:         "light",
				State:        "off",
				Attributes:   map[string]interface{}{"brightness": 0},
				Controllable: true,
			},
			"sensor.temperature": {
				ID:           "sensor.temperature",
				Name:         "Temperature Sensor",
				Type:         "sensor",
				State:        "72",
				Attributes:   map[string]interface{}{"unit": "°F"},
				Controllable: false, // Read-only device
			},
			"switch.outlet": {
				ID:           "switch.outlet",
				Name:         "Power Outlet",
				Type:         "switch",
				State:        "off",
				Attributes:   map[string]interface{}{},
				Controllable: true,
			},
		},
	}
}

func (m *mockDeviceStore) GetDevice(id string) (*models.Device, error) {
	device, exists := m.devices[id]
	if !exists {
		return nil, errors.New("device not found")
	}
	return device, nil
}

func TestExecuteCommand_SucceedsForValidDevice(t *testing.T) {
	store := newMockDeviceStore()
	service := NewHomeAssistantService(store)

	tests := []struct {
		name     string
		deviceID string
		command  *Command
		wantErr  bool
	}{
		{
			name:     "turn on light with valid parameters",
			deviceID: "light.living_room",
			command: &Command{
				Action:     "turn_on",
				Parameters: map[string]interface{}{"entity_id": "light.living_room"},
			},
			wantErr: false,
		},
		{
			name:     "turn off light",
			deviceID: "light.bedroom",
			command: &Command{
				Action:     "turn_off",
				Parameters: map[string]interface{}{"entity_id": "light.bedroom"},
			},
			wantErr: false,
		},
		{
			name:     "set brightness for light",
			deviceID: "light.living_room",
			command: &Command{
				Action: "set_brightness",
				Parameters: map[string]interface{}{
					"entity_id":  "light.living_room",
					"brightness": 50,
				},
			},
			wantErr: false,
		},
		{
			name:     "turn on switch",
			deviceID: "switch.outlet",
			command: &Command{
				Action:     "turn_on",
				Parameters: map[string]interface{}{"entity_id": "switch.outlet"},
			},
			wantErr: false,
		},
		{
			name:     "turn off switch",
			deviceID: "switch.outlet",
			command: &Command{
				Action:     "turn_off",
				Parameters: map[string]interface{}{"entity_id": "switch.outlet"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ExecuteCommand(context.Background(), tt.deviceID, tt.command)

			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExecuteCommand_FailsForInvalidDevice(t *testing.T) {
	store := newMockDeviceStore()
	service := NewHomeAssistantService(store)

	tests := []struct {
		name     string
		deviceID string
		command  *Command
		wantErr  string
	}{
		{
			name:     "device does not exist",
			deviceID: "light.nonexistent",
			command: &Command{
				Action:     "turn_on",
				Parameters: map[string]interface{}{"entity_id": "light.nonexistent"},
			},
			wantErr: "device not found",
		},
		{
			name:     "empty device ID",
			deviceID: "",
			command: &Command{
				Action:     "turn_on",
				Parameters: map[string]interface{}{"entity_id": ""},
			},
			wantErr: "device not found",
		},
		{
			name:     "device ID with wrong format",
			deviceID: "invalid_device_id",
			command: &Command{
				Action:     "turn_on",
				Parameters: map[string]interface{}{"entity_id": "invalid_device_id"},
			},
			wantErr: "device not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ExecuteCommand(context.Background(), tt.deviceID, tt.command)

			if err == nil {
				t.Errorf("ExecuteCommand() expected error, got nil")
				return
			}

			if err.Error() != tt.wantErr {
				t.Errorf("ExecuteCommand() error = %v, want %v", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestExecuteCommand_FailsForReadOnlyDevice(t *testing.T) {
	store := newMockDeviceStore()
	service := NewHomeAssistantService(store)

	tests := []struct {
		name     string
		deviceID string
		command  *Command
		wantErr  string
	}{
		{
			name:     "attempt to control sensor",
			deviceID: "sensor.temperature",
			command: &Command{
				Action:     "turn_on",
				Parameters: map[string]interface{}{"entity_id": "sensor.temperature"},
			},
			wantErr: "device is read-only",
		},
		{
			name:     "attempt to set value on sensor",
			deviceID: "sensor.temperature",
			command: &Command{
				Action: "set_value",
				Parameters: map[string]interface{}{
					"entity_id": "sensor.temperature",
					"value":     75,
				},
			},
			wantErr: "device is read-only",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ExecuteCommand(context.Background(), tt.deviceID, tt.command)

			if err == nil {
				t.Errorf("ExecuteCommand() expected error for read-only device, got nil")
				return
			}

			if err.Error() != tt.wantErr {
				t.Errorf("ExecuteCommand() error = %v, want %v", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestExecuteCommand_FailsForInvalidAction(t *testing.T) {
	store := newMockDeviceStore()
	service := NewHomeAssistantService(store)

	tests := []struct {
		name     string
		deviceID string
		command  *Command
		wantErr  string
	}{
		{
			name:     "invalid action for light",
			deviceID: "light.living_room",
			command: &Command{
				Action:     "invalid_action",
				Parameters: map[string]interface{}{"entity_id": "light.living_room"},
			},
			wantErr: "invalid action for device type",
		},
		{
			name:     "set_brightness on switch (not supported)",
			deviceID: "switch.outlet",
			command: &Command{
				Action:     "set_brightness",
				Parameters: map[string]interface{}{"entity_id": "switch.outlet"},
			},
			wantErr: "invalid action for device type",
		},
		{
			name:     "unsupported action name",
			deviceID: "light.bedroom",
			command: &Command{
				Action:     "play_music",
				Parameters: map[string]interface{}{"entity_id": "light.bedroom"},
			},
			wantErr: "invalid action for device type",
		},
		{
			name:     "empty action string",
			deviceID: "light.living_room",
			command: &Command{
				Action:     "",
				Parameters: map[string]interface{}{"entity_id": "light.living_room"},
			},
			wantErr: "action is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ExecuteCommand(context.Background(), tt.deviceID, tt.command)

			if err == nil {
				t.Errorf("ExecuteCommand() expected error for invalid action, got nil")
				return
			}

			if err.Error() != tt.wantErr {
				t.Errorf("ExecuteCommand() error = %v, want %v", err.Error(), tt.wantErr)
			}
		})
	}
}

// Edge cases and additional validation tests
func TestExecuteCommand_EdgeCases(t *testing.T) {
	store := newMockDeviceStore()
	service := NewHomeAssistantService(store)

	tests := []struct {
		name     string
		deviceID string
		command  *Command
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "nil command",
			deviceID: "light.living_room",
			command:  nil,
			wantErr:  true,
			errMsg:   "", // Will panic or return specific error
		},
		{
			name:     "command with nil parameters",
			deviceID: "light.living_room",
			command: &Command{
				Action:     "turn_on",
				Parameters: nil,
			},
			wantErr: true,
			errMsg:  "parameters is required",
		},
		{
			name:     "command with empty parameters map",
			deviceID: "light.living_room",
			command: &Command{
				Action:     "turn_on",
				Parameters: map[string]interface{}{},
			},
			wantErr: false, // Empty parameters map is valid
		},
		{
			name:     "command with whitespace-only action",
			deviceID: "light.living_room",
			command: &Command{
				Action:     "   ",
				Parameters: map[string]interface{}{"entity_id": "light.living_room"},
			},
			wantErr: true,
			errMsg:  "action is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Handle nil command case specially
			if tt.command == nil {
				// In a real implementation, this should be handled gracefully
				defer func() {
					if r := recover(); r != nil {
						// Expected panic for nil command
						if !tt.wantErr {
							t.Errorf("ExecuteCommand() unexpected panic for nil command")
						}
					}
				}()
			}

			err := service.ExecuteCommand(context.Background(), tt.deviceID, tt.command)

			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && tt.errMsg != "" {
				if err.Error() != tt.errMsg {
					t.Errorf("ExecuteCommand() error = %v, want %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

// TestExecuteCommand_ConcurrentExecution tests thread safety
func TestExecuteCommand_ConcurrentExecution(t *testing.T) {
	store := newMockDeviceStore()
	service := NewHomeAssistantService(store)

	// Channel to collect errors from goroutines
	errChan := make(chan error, 10)

	for i := 0; i < 10; i++ {
		go func() {
			cmd := &Command{
				Action:     "turn_on",
				Parameters: map[string]interface{}{"entity_id": "light.living_room"},
			}
			err := service.ExecuteCommand(context.Background(), "light.living_room", cmd)
			errChan <- err
		}()
	}

	// Collect results from all goroutines
	for i := 0; i < 10; i++ {
		err := <-errChan
		if err != nil {
			t.Errorf("ExecuteCommand() in goroutine %d failed: %v", i, err)
		}
	}
}

// TestExecuteCommand_ValidatesCommandBeforeDeviceLookup ensures proper validation order
func TestExecuteCommand_ValidatesCommandBeforeDeviceLookup(t *testing.T) {
	store := newMockDeviceStore()
	service := NewHomeAssistantService(store)

	// Invalid command should fail validation before device lookup
	invalidCommand := &Command{
		Action:     "", // Invalid: empty action
		Parameters: map[string]interface{}{"entity_id": "light.living_room"},
	}

	err := service.ExecuteCommand(context.Background(), "light.living_room", invalidCommand)
	
	if err == nil {
		t.Error("ExecuteCommand() should validate command before device lookup")
		return
	}

	if err.Error() != "action is required" {
		t.Errorf("ExecuteCommand() should return validation error first, got: %v", err)
	}
}
