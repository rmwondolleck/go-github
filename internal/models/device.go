package models

import "time"

// Device represents a smart home device with its state and attributes
type Device struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Type         string                 `json:"type"`
	State        string                 `json:"state"`
	Attributes   map[string]interface{} `json:"attributes"`
	LastUpdated  time.Time              `json:"last_updated"`
	Controllable bool                   `json:"controllable"`
}
