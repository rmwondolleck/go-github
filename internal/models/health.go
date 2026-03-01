package models

// HealthStatus represents the health status of the API service
type HealthStatus struct {
	Status     string            `json:"status"`
	Uptime     string            `json:"uptime"`
	Components map[string]string `json:"components"`
}
