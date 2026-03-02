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

// Service represents the HomeAssistant service interface
// This is a placeholder for the actual implementation
type Service interface {
	ExecuteCommand(deviceID string, command *Command) error
}

// mockService is a mock implementation for testing
// This simulates the expected behavior of the real service
type mockService struct {
	store *mockDeviceStore
}

func newMockService(store *mockDeviceStore) *mockService {
	return &mockService{
		store: store,
	}
}

// ExecuteCommand simulates command execution logic
// This is what the real implementation should do
func (s *mockService) ExecuteCommand(deviceID string, command *Command) error {
	// Validate command
	if err := command.Validate(); err != nil {
		return err
	}

	// Check if device exists
	device, err := s.store.GetDevice(deviceID)
	if err != nil {
		return errors.New("device not found")
	}

	// Check if device is controllable
	if !device.Controllable {
		return errors.New("device is read-only")
	}

	// Validate action based on device type
	validActions := getValidActionsForDeviceType(device.Type)
	if !isValidAction(command.Action, validActions) {
		return errors.New("invalid action for device type")
	}

	// Command execution would happen here in real implementation
	return nil
}

// Helper function to get valid actions for a device type
func getValidActionsForDeviceType(deviceType string) []string {
	validActionsMap := map[string][]string{
		"light": {"turn_on", "turn_off", "set_brightness"},
		"switch": {"turn_on", "turn_off"},
		"sensor": {}, // Sensors have no valid actions as they're read-only
	}
	
	if actions, exists := validActionsMap[deviceType]; exists {
		return actions
	}
	return []string{}
}

// Helper function to check if an action is valid
func isValidAction(action string, validActions []string) bool {
	for _, validAction := range validActions {
		if action == validAction {
			return true
		}
	}
	return false
}

func TestExecuteCommand_SucceedsForValidDevice(t *testing.T) {
	store := newMockDeviceStore()
	service := newMockService(store)

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
			err := service.ExecuteCommand(tt.deviceID, tt.command)

			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExecuteCommand_FailsForInvalidDevice(t *testing.T) {
	store := newMockDeviceStore()
	service := newMockService(store)

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
			err := service.ExecuteCommand(tt.deviceID, tt.command)

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
	service := newMockService(store)

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
			err := service.ExecuteCommand(tt.deviceID, tt.command)

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
	service := newMockService(store)

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
			err := service.ExecuteCommand(tt.deviceID, tt.command)

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
	service := newMockService(store)

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

			err := service.ExecuteCommand(tt.deviceID, tt.command)

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
	service := newMockService(store)

	// Run multiple commands concurrently to test thread safety
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			cmd := &Command{
				Action:     "turn_on",
				Parameters: map[string]interface{}{"entity_id": "light.living_room"},
			}
			err := service.ExecuteCommand("light.living_room", cmd)
			if err != nil {
				t.Errorf("ExecuteCommand() in goroutine failed: %v", err)
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestExecuteCommand_ValidatesCommandBeforeDeviceLookup ensures proper validation order
func TestExecuteCommand_ValidatesCommandBeforeDeviceLookup(t *testing.T) {
	store := newMockDeviceStore()
	service := newMockService(store)

	// Invalid command should fail validation before device lookup
	invalidCommand := &Command{
		Action:     "", // Invalid: empty action
		Parameters: map[string]interface{}{"entity_id": "light.living_room"},
	}

	err := service.ExecuteCommand("light.living_room", invalidCommand)
	
	if err == nil {
		t.Error("ExecuteCommand() should validate command before device lookup")
		return
	}

	if err.Error() != "action is required" {
		t.Errorf("ExecuteCommand() should return validation error first, got: %v", err)
	}
}
