package integration

import (
"net/http"
"net/http/httptest"
"os"
"testing"

"go-github/internal/middleware"
"github.com/gin-gonic/gin"
"github.com/stretchr/testify/assert"
)

func TestCORS_Integration(t *testing.T) {
gin.SetMode(gin.TestMode)

t.Run("default origin integration", func(t *testing.T) {
os.Unsetenv("CORS_ORIGINS")
defer os.Unsetenv("CORS_ORIGINS")

router := gin.New()
router.Use(middleware.CORS())
router.GET("/api/test", func(c *gin.Context) {
c.JSON(http.StatusOK, gin.H{"message": "ok"})
})

req := httptest.NewRequest("GET", "/api/test", nil)
req.Header.Set("Origin", "http://localhost:3000")
w := httptest.NewRecorder()
router.ServeHTTP(w, req)

assert.Equal(t, http.StatusOK, w.Code)
assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
assert.Equal(t, "GET, POST, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
assert.Equal(t, "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
})

t.Run("custom origin integration", func(t *testing.T) {
os.Setenv("CORS_ORIGINS", "https://example.com")
defer os.Unsetenv("CORS_ORIGINS")

router := gin.New()
router.Use(middleware.CORS())
router.GET("/api/test", func(c *gin.Context) {
c.JSON(http.StatusOK, gin.H{"message": "ok"})
})

req := httptest.NewRequest("GET", "/api/test", nil)
req.Header.Set("Origin", "https://example.com")
w := httptest.NewRecorder()
router.ServeHTTP(w, req)

assert.Equal(t, http.StatusOK, w.Code)
assert.Equal(t, "https://example.com", w.Header().Get("Access-Control-Allow-Origin"))
assert.Equal(t, "GET, POST, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
assert.Equal(t, "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
})

t.Run("multiple origins integration", func(t *testing.T) {
os.Setenv("CORS_ORIGINS", "https://app.example.com,https://admin.example.com,http://localhost:8080")
defer os.Unsetenv("CORS_ORIGINS")

router := gin.New()
router.Use(middleware.CORS())
router.GET("/api/test", func(c *gin.Context) {
c.JSON(http.StatusOK, gin.H{"message": "ok"})
})

// Test first origin
req1 := httptest.NewRequest("GET", "/api/test", nil)
req1.Header.Set("Origin", "https://app.example.com")
w1 := httptest.NewRecorder()
router.ServeHTTP(w1, req1)

assert.Equal(t, http.StatusOK, w1.Code)
assert.Equal(t, "https://app.example.com", w1.Header().Get("Access-Control-Allow-Origin"))

// Test second origin
req2 := httptest.NewRequest("GET", "/api/test", nil)
req2.Header.Set("Origin", "https://admin.example.com")
w2 := httptest.NewRecorder()
router.ServeHTTP(w2, req2)

assert.Equal(t, http.StatusOK, w2.Code)
assert.Equal(t, "https://admin.example.com", w2.Header().Get("Access-Control-Allow-Origin"))
})

t.Run("preflight OPTIONS request integration", func(t *testing.T) {
os.Setenv("CORS_ORIGINS", "https://example.com")
defer os.Unsetenv("CORS_ORIGINS")

router := gin.New()
router.Use(middleware.CORS())
router.POST("/api/test", func(c *gin.Context) {
c.JSON(http.StatusOK, gin.H{"message": "ok"})
})

req := httptest.NewRequest("OPTIONS", "/api/test", nil)
req.Header.Set("Origin", "https://example.com")
req.Header.Set("Access-Control-Request-Method", "POST")
w := httptest.NewRecorder()
router.ServeHTTP(w, req)

assert.Equal(t, http.StatusNoContent, w.Code)
assert.Equal(t, "https://example.com", w.Header().Get("Access-Control-Allow-Origin"))
assert.Equal(t, "GET, POST, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
assert.Equal(t, "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
assert.Empty(t, w.Body.String())
})

t.Run("forbidden origin integration", func(t *testing.T) {
os.Setenv("CORS_ORIGINS", "https://example.com")
defer os.Unsetenv("CORS_ORIGINS")

router := gin.New()
router.Use(middleware.CORS())
router.GET("/api/test", func(c *gin.Context) {
c.JSON(http.StatusOK, gin.H{"message": "ok"})
})

req := httptest.NewRequest("GET", "/api/test", nil)
req.Header.Set("Origin", "https://evil.com")
w := httptest.NewRecorder()
router.ServeHTTP(w, req)

assert.Equal(t, http.StatusOK, w.Code)
assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
assert.Empty(t, w.Header().Get("Access-Control-Allow-Methods"))
assert.Empty(t, w.Header().Get("Access-Control-Allow-Headers"))
})
}
