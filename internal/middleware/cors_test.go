package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		corsOrigins    string
		requestOrigin  string
		requestMethod  string
		expectedStatus int
		expectCORS     bool
		wantOrigin     string
	}{
		{
			name:           "default origin allowed",
			corsOrigins:    "",
			requestOrigin:  "http://localhost:3000",
			requestMethod:  "GET",
			expectedStatus: 200,
			expectCORS:     true,
			wantOrigin:     "http://localhost:3000",
		},
		{
			name:           "custom origin allowed",
			corsOrigins:    "https://example.com",
			requestOrigin:  "https://example.com",
			requestMethod:  "GET",
			expectedStatus: 200,
			expectCORS:     true,
			wantOrigin:     "https://example.com",
		},
		{
			name:           "multiple origins - first allowed",
			corsOrigins:    "https://example.com,https://app.example.com,http://localhost:8080",
			requestOrigin:  "https://example.com",
			requestMethod:  "GET",
			expectedStatus: 200,
			expectCORS:     true,
			wantOrigin:     "https://example.com",
		},
		{
			name:           "multiple origins - middle allowed",
			corsOrigins:    "https://example.com,https://app.example.com,http://localhost:8080",
			requestOrigin:  "https://app.example.com",
			requestMethod:  "POST",
			expectedStatus: 200,
			expectCORS:     true,
			wantOrigin:     "https://app.example.com",
		},
		{
			name:           "multiple origins - last allowed",
			corsOrigins:    "https://example.com,https://app.example.com,http://localhost:8080",
			requestOrigin:  "http://localhost:8080",
			requestMethod:  "GET",
			expectedStatus: 200,
			expectCORS:     true,
			wantOrigin:     "http://localhost:8080",
		},
		{
			name:           "forbidden origin",
			corsOrigins:    "https://example.com",
			requestOrigin:  "https://evil.com",
			requestMethod:  "GET",
			expectedStatus: 200,
			expectCORS:     false,
			wantOrigin:     "",
		},
		{
			name:           "no origin header",
			corsOrigins:    "https://example.com",
			requestOrigin:  "",
			requestMethod:  "GET",
			expectedStatus: 200,
			expectCORS:     false,
			wantOrigin:     "",
		},
		{
			name:           "multiple origins with spaces",
			corsOrigins:    "https://example.com, https://app.example.com , http://localhost:8080",
			requestOrigin:  "https://app.example.com",
			requestMethod:  "GET",
			expectedStatus: 200,
			expectCORS:     true,
			wantOrigin:     "https://app.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			if tt.corsOrigins != "" {
				os.Setenv("CORS_ORIGINS", tt.corsOrigins)
			} else {
				os.Unsetenv("CORS_ORIGINS")
			}
			defer os.Unsetenv("CORS_ORIGINS")

			// Create router with CORS middleware
			router := gin.New()
			router.Use(CORS())

			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "ok"})
			})
			router.POST("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "ok"})
			})

			// Create test request
			req := httptest.NewRequest(tt.requestMethod, "/test", nil)
			if tt.requestOrigin != "" {
				req.Header.Set("Origin", tt.requestOrigin)
			}
			w := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(w, req)

			// Verify response code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Verify CORS headers
			if tt.expectCORS {
				assert.Equal(t, tt.wantOrigin, w.Header().Get("Access-Control-Allow-Origin"))
				assert.Equal(t, "GET, POST, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
				assert.Equal(t, "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
				assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
			} else {
				assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
			}
		})
	}
}

func TestCORS_PreflightRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		corsOrigins    string
		requestOrigin  string
		expectCORS     bool
		wantOrigin     string
		expectedStatus int
	}{
		{
			name:           "preflight with allowed origin",
			corsOrigins:    "https://example.com",
			requestOrigin:  "https://example.com",
			expectCORS:     true,
			wantOrigin:     "https://example.com",
			expectedStatus: 204,
		},
		{
			name:           "preflight with forbidden origin",
			corsOrigins:    "https://example.com",
			requestOrigin:  "https://evil.com",
			expectCORS:     false,
			wantOrigin:     "",
			expectedStatus: 404, // No route handler for OPTIONS, returns 404
		},
		{
			name:           "preflight with default origin",
			corsOrigins:    "",
			requestOrigin:  "http://localhost:3000",
			expectCORS:     true,
			wantOrigin:     "http://localhost:3000",
			expectedStatus: 204,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			if tt.corsOrigins != "" {
				os.Setenv("CORS_ORIGINS", tt.corsOrigins)
			} else {
				os.Unsetenv("CORS_ORIGINS")
			}
			defer os.Unsetenv("CORS_ORIGINS")

			// Create router with CORS middleware
			router := gin.New()
			router.Use(CORS())

			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "ok"})
			})

			// Create preflight OPTIONS request
			req := httptest.NewRequest("OPTIONS", "/test", nil)
			req.Header.Set("Origin", tt.requestOrigin)
			req.Header.Set("Access-Control-Request-Method", "POST")
			w := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(w, req)

			// Verify response code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Verify CORS headers
			if tt.expectCORS {
				assert.Equal(t, tt.wantOrigin, w.Header().Get("Access-Control-Allow-Origin"))
				assert.Equal(t, "GET, POST, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
				assert.Equal(t, "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
				assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
				// Verify response body is empty for allowed preflight
				assert.Empty(t, w.Body.String())
			} else {
				assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
			}
		})
	}
}

func TestCORS_CallsNext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	os.Setenv("CORS_ORIGINS", "https://example.com")
	defer os.Unsetenv("CORS_ORIGINS")

	router := gin.New()
	router.Use(CORS())

	handlerCalled := false
	router.GET("/test", func(c *gin.Context) {
		handlerCalled = true
		c.Status(http.StatusOK)
	})

	// Create test request (not OPTIONS)
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Verify handler was called
	assert.True(t, handlerCalled, "CORS middleware should call c.Next() for non-OPTIONS requests")
}

func TestCORS_PreflightDoesNotCallNext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	os.Setenv("CORS_ORIGINS", "https://example.com")
	defer os.Unsetenv("CORS_ORIGINS")

	router := gin.New()
	router.Use(CORS())

	handlerCalled := false
	router.OPTIONS("/test", func(c *gin.Context) {
		handlerCalled = true
		c.Status(http.StatusOK)
	})

	// Create preflight OPTIONS request with allowed origin
	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Verify handler was NOT called (preflight with allowed origin should abort)
	assert.False(t, handlerCalled, "CORS middleware should abort for OPTIONS requests from allowed origins")
	assert.Equal(t, 204, w.Code)
}

func TestCORS_PreflightForbiddenOriginCallsNext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	os.Setenv("CORS_ORIGINS", "https://example.com")
	defer os.Unsetenv("CORS_ORIGINS")

	router := gin.New()
	router.Use(CORS())

	handlerCalled := false
	router.OPTIONS("/test", func(c *gin.Context) {
		handlerCalled = true
		c.Status(http.StatusOK)
	})

	// Create preflight OPTIONS request with forbidden origin
	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "https://evil.com")
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Verify handler WAS called (preflight with forbidden origin should continue to handler)
	assert.True(t, handlerCalled, "CORS middleware should call c.Next() for OPTIONS requests from forbidden origins")
	assert.Equal(t, 200, w.Code)
}

func TestParseOrigins(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single origin",
			input:    "https://example.com",
			expected: []string{"https://example.com"},
		},
		{
			name:     "multiple origins",
			input:    "https://example.com,https://app.example.com,http://localhost:3000",
			expected: []string{"https://example.com", "https://app.example.com", "http://localhost:3000"},
		},
		{
			name:     "origins with spaces",
			input:    "https://example.com, https://app.example.com , http://localhost:3000",
			expected: []string{"https://example.com", "https://app.example.com", "http://localhost:3000"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: nil,
		},
		{
			name:     "origin with trailing comma",
			input:    "https://example.com,",
			expected: []string{"https://example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseOrigins(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsOriginAllowed(t *testing.T) {
	allowedOrigins := []string{"https://example.com", "https://app.example.com", "http://localhost:3000"}

	tests := []struct {
		name     string
		origin   string
		expected bool
	}{
		{
			name:     "allowed origin",
			origin:   "https://example.com",
			expected: true,
		},
		{
			name:     "another allowed origin",
			origin:   "http://localhost:3000",
			expected: true,
		},
		{
			name:     "forbidden origin",
			origin:   "https://evil.com",
			expected: false,
		},
		{
			name:     "similar but not exact origin",
			origin:   "https://example.com.evil.com",
			expected: false,
		},
		{
			name:     "case sensitive",
			origin:   "https://Example.com",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isOriginAllowed(tt.origin, allowedOrigins)
			assert.Equal(t, tt.expected, result)
		})
	}
}
