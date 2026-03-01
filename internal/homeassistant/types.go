package homeassistant

import (
	"errors"
	"strings"
)

// Command represents a device control command for HomeAssistant
type Command struct {
	Action     string                 `json:"action"`
	Parameters map[string]interface{} `json:"parameters"`
}

// Validate checks if the Command has valid data
func (c *Command) Validate() error {
	if c.Action == "" {
		return errors.New("action is required")
	}

	if strings.TrimSpace(c.Action) == "" {
		return errors.New("action cannot be empty or whitespace")
	}

	if c.Parameters == nil {
		return errors.New("parameters is required")
	}

	return nil
}
