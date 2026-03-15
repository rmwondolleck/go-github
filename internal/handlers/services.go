package handlers

import (
	"go-github/internal/models"
	"go-github/internal/services"
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
	response := models.ServicesResponse{
		Services: services.GetServices(),
	}

	JSONSuccess(c, http.StatusOK, response)
}
