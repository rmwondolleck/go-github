package models

// Service represents a service in the homelab
type Service struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Status   string `json:"status"`
	Endpoint string `json:"endpoint"`
}

// ServicesResponse represents the response for listing services
type ServicesResponse struct {
	Services []Service `json:"services"`
}
