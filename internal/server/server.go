package server

import (
	"net/http"
	"sync"

	"go-github/internal/middleware"

	"github.com/gin-gonic/gin"
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

	// Health endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API v1 routes group
	v1 := router.Group("/api/v1")
	{
		// Placeholder for API routes
		v1.GET("", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "API v1",
			})
		})
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
