package research

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthResponse represents the JSON response structure for the health check endpoint
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// Standard library net/http handler
func stdlibHealthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Gin framework handler
func ginHealthHandler(c *gin.Context) {
	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	c.JSON(http.StatusOK, response)
}

// BenchmarkStdlibSimpleRoute benchmarks a simple GET request using stdlib net/http
func BenchmarkStdlibSimpleRoute(b *testing.B) {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", stdlibHealthHandler)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
	}
}

// BenchmarkGinSimpleRoute benchmarks a simple GET request using Gin framework
func BenchmarkGinSimpleRoute(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.GET("/health", ginHealthHandler)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

// BenchmarkStdlibWithMiddleware benchmarks stdlib with a simple logging middleware
func BenchmarkStdlibWithMiddleware(b *testing.B) {
	loggingMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Simple middleware that just calls next
			next(w, r)
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", loggingMiddleware(stdlibHealthHandler))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
	}
}

// BenchmarkGinWithMiddleware benchmarks Gin with a simple logging middleware
func BenchmarkGinWithMiddleware(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	
	// Simple middleware
	router.Use(func(c *gin.Context) {
		c.Next()
	})
	
	router.GET("/health", ginHealthHandler)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

// BenchmarkStdlibMultipleRoutes benchmarks stdlib with 5 different routes
func BenchmarkStdlibMultipleRoutes(b *testing.B) {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", stdlibHealthHandler)
	mux.HandleFunc("/api/v1/users", stdlibHealthHandler)
	mux.HandleFunc("/api/v1/posts", stdlibHealthHandler)
	mux.HandleFunc("/api/v1/comments", stdlibHealthHandler)
	mux.HandleFunc("/api/v1/status", stdlibHealthHandler)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
	}
}

// BenchmarkGinMultipleRoutes benchmarks Gin with 5 different routes
func BenchmarkGinMultipleRoutes(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.GET("/health", ginHealthHandler)
	router.GET("/api/v1/users", ginHealthHandler)
	router.GET("/api/v1/posts", ginHealthHandler)
	router.GET("/api/v1/comments", ginHealthHandler)
	router.GET("/api/v1/status", ginHealthHandler)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

// BenchmarkStdlibParameterizedRoute benchmarks stdlib with URL parameters
func BenchmarkStdlibParameterizedRoute(b *testing.B) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/users/", func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from path (simple version)
		response := HealthResponse{
			Status:    "ok",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/123", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
	}
}

// BenchmarkGinParameterizedRoute benchmarks Gin with URL parameters
func BenchmarkGinParameterizedRoute(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.GET("/api/v1/users/:id", func(c *gin.Context) {
		// Extract user ID from path
		_ = c.Param("id")
		response := HealthResponse{
			Status:    "ok",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
		c.JSON(http.StatusOK, response)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/123", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
