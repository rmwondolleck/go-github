package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	defaultCORSOrigins = "http://localhost:3000"
)

// CORS returns a gin.HandlerFunc middleware that adds CORS headers to responses.
// It reads allowed origins from the CORS_ORIGINS environment variable (default: http://localhost:3000).
// Multiple origins can be specified separated by commas.
// Allowed methods: GET, POST, OPTIONS
// Allowed headers: Content-Type, Authorization
func CORS() gin.HandlerFunc {
	// Get allowed origins from environment or use default
	originsStr := os.Getenv("CORS_ORIGINS")
	if originsStr == "" {
		originsStr = defaultCORSOrigins
	}

	// Parse origins into a slice
	allowedOrigins := parseOrigins(originsStr)

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		if origin != "" && isOriginAllowed(origin, allowedOrigins) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// Handle preflight OPTIONS requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// parseOrigins splits a comma-separated string of origins into a slice
func parseOrigins(originsStr string) []string {
	origins := strings.Split(originsStr, ",")
	var result []string
	for _, origin := range origins {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// isOriginAllowed checks if the given origin is in the allowed list
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if origin == allowed {
			return true
		}
	}
	return false
}
