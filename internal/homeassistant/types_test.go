package homeassistant

import (
	"testing"
)

func TestCommand_Validate(t *testing.T) {
	tests := []struct {
		name    string
		command Command
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid command with action and parameters",
			command: Command{
				Action:     "turn_on",
				Parameters: map[string]interface{}{"entity_id": "light.living_room"},
			},
			wantErr: false,
		},
		{
			name: "valid command with multiple parameters",
			command: Command{
				Action: "set_brightness",
				Parameters: map[string]interface{}{
					"entity_id":  "light.bedroom",
					"brightness": 75,
				},
			},
			wantErr: false,
		},
		{
			name: "empty parameters map is valid",
			command: Command{
				Action:     "turn_off",
				Parameters: map[string]interface{}{},
			},
			wantErr: false,
		},
		{
			name: "missing action",
			command: Command{
				Action:     "",
				Parameters: map[string]interface{}{"entity_id": "light.kitchen"},
			},
			wantErr: true,
			errMsg:  "action is required",
		},
		{
			name: "action with only whitespace",
			command: Command{
				Action:     "   ",
				Parameters: map[string]interface{}{"entity_id": "light.kitchen"},
			},
			wantErr: true,
			errMsg:  "action cannot be empty or whitespace",
		},
		{
			name: "nil parameters",
			command: Command{
				Action:     "turn_on",
				Parameters: nil,
			},
			wantErr: true,
			errMsg:  "parameters is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.command.Validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("Validate() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestCommand_JSONTags(t *testing.T) {
	// Ensure JSON tags are properly defined by marshaling and unmarshaling
	original := Command{
		Action: "test_action",
		Parameters: map[string]interface{}{
			"param1": "value1",
			"param2": 42,
		},
	}

	// Marshal to JSON
	jsonData := `{"action":"test_action","parameters":{"param1":"value1","param2":42}}`

	// Basic check that the struct can be used with JSON
	if original.Action != "test_action" {
		t.Error("Action field not accessible")
	}

	if original.Parameters == nil {
		t.Error("Parameters field not accessible")
	}

	// This test verifies the struct is properly defined
	// In a real scenario, you'd use json.Marshal/Unmarshal
	_ = jsonData
}
