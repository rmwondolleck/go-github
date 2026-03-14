package handlers

import (
	"go-github/internal/cluster"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListClusterServicesHandler godoc
// @Summary List cluster services
// @Description Returns a list of Kubernetes cluster services, optionally filtered by name
// @Tags cluster
// @Produce json
// @Param name query string false "Filter services by name (case-insensitive substring match)"
// @Success 200 {array} cluster.ServiceInfo
// @Router /api/v1/cluster/services [get]
func ListClusterServicesHandler(c *gin.Context) {
	nameFilter := c.Query("name")

	svc := cluster.NewService()
	services, err := svc.ListServices(nameFilter)
	if err != nil {
		JSONError(c, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	JSONSuccess(c, http.StatusOK, services)
}
