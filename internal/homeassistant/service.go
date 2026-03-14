package homeassistant

import (
	"context"
	"errors"
	"go-github/internal/models"
)

// DeviceStore defines the interface for accessing device data
type DeviceStore interface {
	GetDevice(id string) (*models.Device, error)
}

// HomeAssistantService provides command execution for HomeAssistant devices
type HomeAssistantService struct {
	store DeviceStore
}

// NewHomeAssistantService creates a new HomeAssistantService
func NewHomeAssistantService(store DeviceStore) *HomeAssistantService {
	return &HomeAssistantService{store: store}
}

// validActionsMap maps device types to their supported actions
var validActionsMap = map[string][]string{
	"light":  {"turn_on", "turn_off", "set_brightness"},
	"switch": {"turn_on", "turn_off"},
	"sensor": {},
}

// ExecuteCommand validates and executes a command on the specified device.
// It validates: 1) command fields, 2) device exists, 3) device is controllable,
// 4) action is supported for the device type. Returns nil on success (mock execution).
// ctx is accepted for future use when real HomeAssistant API calls are introduced.
func (s *HomeAssistantService) ExecuteCommand(ctx context.Context, id string, cmd *Command) error {
	if cmd == nil {
		return errors.New("command is required")
	}

	// Validate command fields
	if err := cmd.Validate(); err != nil {
		return err
	}

	// Check if device exists
	device, err := s.store.GetDevice(id)
	if err != nil {
		return errors.New("device not found")
	}

	// Check if device is controllable
	if !device.Controllable {
		return errors.New("device is read-only")
	}

	// Validate action based on device type
	validActions, ok := validActionsMap[device.Type]
	if !ok || !containsAction(validActions, cmd.Action) {
		return errors.New("invalid action for device type")
	}

	// Mock execution
	return nil
}

// containsAction checks whether action is in the list of valid actions
func containsAction(validActions []string, action string) bool {
	for _, a := range validActions {
		if a == action {
			return true
		}
	}
	return false
}
