package handlers

import (
	"go-github/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListServicesHandler godoc
// @Summary List available services
// @Description Returns list of all services in the homelab
// @Tags services
// @Produce json
// @Success 200 {object} models.ServicesResponse
// @Router /api/v1/services [get]
func ListServicesHandler(c *gin.Context) {
	services := []models.Service{
		{
			Name:     "homeassistant",
			Type:     "home-automation",
			Status:   "running",
			Endpoint: "http://homeassistant.local:8123",
		},
		{
			Name:     "prometheus",
			Type:     "monitoring",
			Status:   "running",
			Endpoint: "http://prometheus.local:9090",
		},
		{
			Name:     "grafana",
			Type:     "visualization",
			Status:   "running",
			Endpoint: "http://grafana.local:3000",
		},
		{
			Name:     "node-exporter",
			Type:     "metrics",
			Status:   "running",
			Endpoint: "http://node-exporter.local:9100",
		},
		{
			Name:     "alertmanager",
			Type:     "alerting",
			Status:   "running",
			Endpoint: "http://alertmanager.local:9093",
		},
	}

	response := models.ServicesResponse{
		Services: services,
	}

	JSONSuccess(c, http.StatusOK, response)
}
