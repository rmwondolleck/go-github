package server

import (
	"net/http"
	"sync"

	"go-github/internal/handlers"
	"go-github/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Server represents the HTTP server
type Server struct {
	router     *gin.Engine
	httpServer *http.Server
	mu         sync.RWMutex
}

// New creates a new server instance with middleware chain
func New() *Server {
	router := gin.New()
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())

	// Swagger documentation
	router.GET("/api/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health endpoint
	router.GET("/health", healthHandler)

	// API v1 routes group — rate limiting applied here only
	v1 := router.Group("/api/v1")
	v1.Use(middleware.RateLimit())
	{
		// Placeholder for API routes
		v1.GET("", apiRootHandler)
		v1.GET("/services", handlers.ListServicesHandler)

		// Cluster services endpoint
		v1.GET("/cluster/services", handlers.ListClusterServicesHandler)

		// HomeAssistant device endpoints
		v1.POST("/homeassistant/devices/:id/command", handlers.ExecuteCommandHandler)
	}

	return &Server{router: router}
}

// Run starts the HTTP server on the specified port
func (s *Server) Run(port string) error {
	s.mu.Lock()
	s.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: s.router,
	}
	s.mu.Unlock()
	return s.httpServer.ListenAndServe()
}

// Router returns the gin router (useful for testing)
func (s *Server) Router() *gin.Engine {
	return s.router
}

// healthHandler godoc
// @Summary Health check
// @Description Get the health status of the API
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// apiRootHandler godoc
// @Summary API root
// @Description Get API version information
// @Tags api
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/v1 [get]
func apiRootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "API v1",
	})
}
