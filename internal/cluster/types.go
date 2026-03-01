package cluster

// ServiceInfo represents information about a Kubernetes cluster service
type ServiceInfo struct {
	Name      string   `json:"name"`
	Namespace string   `json:"namespace"`
	Status    string   `json:"status"`
	Endpoints []string `json:"endpoints"`
}
