package homeassistant

import (
	"encoding/json"
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
			errMsg:  "action is required",
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
	// Test JSON marshaling
	original := Command{
		Action: "test_action",
		Parameters: map[string]interface{}{
			"param1": "value1",
			"param2": float64(42),
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal Command: %v", err)
	}

	// Verify JSON structure contains expected fields
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		t.Fatalf("failed to unmarshal JSON to map: %v", err)
	}

	if jsonMap["action"] != "test_action" {
		t.Errorf("expected action 'test_action', got %v", jsonMap["action"])
	}

	params, ok := jsonMap["parameters"].(map[string]interface{})
	if !ok {
		t.Fatal("parameters not found or not a map")
	}

	if params["param1"] != "value1" {
		t.Errorf("expected param1 'value1', got %v", params["param1"])
	}

	if params["param2"] != float64(42) {
		t.Errorf("expected param2 42, got %v", params["param2"])
	}

	// Test JSON unmarshaling
	jsonInput := `{"action":"turn_on","parameters":{"entity_id":"light.living_room","brightness":80}}`
	var cmd Command
	if err := json.Unmarshal([]byte(jsonInput), &cmd); err != nil {
		t.Fatalf("failed to unmarshal JSON to Command: %v", err)
	}

	if cmd.Action != "turn_on" {
		t.Errorf("expected action 'turn_on', got %s", cmd.Action)
	}

	if cmd.Parameters["entity_id"] != "light.living_room" {
		t.Errorf("expected entity_id 'light.living_room', got %v", cmd.Parameters["entity_id"])
	}

	if cmd.Parameters["brightness"] != float64(80) {
		t.Errorf("expected brightness 80, got %v", cmd.Parameters["brightness"])
	}
}
