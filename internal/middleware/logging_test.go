package middleware

import (
"bytes"
"log/slog"
"net/http"
"net/http/httptest"
"testing"
"time"

"github.com/gin-gonic/gin"
)

func TestLogger(t *testing.T) {
// Set Gin to test mode
gin.SetMode(gin.TestMode)

tests := []struct {
name          string
method        string
path          string
requestID     interface{}
setRequestID  bool
handlerDelay  time.Duration
expectedCode  int
wantMethod    string
wantPath      string
wantRequestID bool
}{
{
name:          "GET request with request ID",
method:        "GET",
path:          "/api/test",
requestID:     "test-request-id-123",
setRequestID:  true,
handlerDelay:  10 * time.Millisecond,
expectedCode:  200,
wantMethod:    "GET",
wantPath:      "/api/test",
wantRequestID: true,
},
{
name:          "POST request without request ID",
method:        "POST",
path:          "/api/create",
requestID:     nil,
setRequestID:  false,
handlerDelay:  5 * time.Millisecond,
expectedCode:  201,
wantMethod:    "POST",
wantPath:      "/api/create",
wantRequestID: false,
},
{
name:          "DELETE request with numeric request ID",
method:        "DELETE",
path:          "/api/delete/1",
requestID:     12345,
setRequestID:  true,
handlerDelay:  0,
expectedCode:  204,
wantMethod:    "DELETE",
wantPath:      "/api/delete/1",
wantRequestID: true,
},
{
name:          "Error response with request ID",
method:        "GET",
path:          "/api/error",
requestID:     "error-request-id",
setRequestID:  true,
handlerDelay:  0,
expectedCode:  500,
wantMethod:    "GET",
wantPath:      "/api/error",
wantRequestID: true,
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
// Create a buffer to capture log output
var logBuf bytes.Buffer
logger := slog.New(slog.NewJSONHandler(&logBuf, &slog.HandlerOptions{
Level: slog.LevelInfo,
}))
slog.SetDefault(logger)

// Create router with logging middleware
router := gin.New()
// Simulate T014 request ID middleware if needed
if tt.setRequestID {
router.Use(func(c *gin.Context) {
c.Set("request_id", tt.requestID)
c.Next()
})
}
router.Use(Logger())

// Add test handler
router.Handle(tt.method, tt.path, func(c *gin.Context) {
if tt.handlerDelay > 0 {
time.Sleep(tt.handlerDelay)
}
c.Status(tt.expectedCode)
})

// Create test request
req := httptest.NewRequest(tt.method, tt.path, nil)
w := httptest.NewRecorder()

// Serve the request
router.ServeHTTP(w, req)

// Verify response code
if w.Code != tt.expectedCode {
t.Errorf("expected status code %d, got %d", tt.expectedCode, w.Code)
}

// Verify log output
logOutput := logBuf.String()
if logOutput == "" {
t.Fatal("expected log output, got empty string")
}

// Check for required fields in log
if !bytes.Contains(logBuf.Bytes(), []byte("request completed")) {
t.Error("log should contain 'request completed' message")
}
if !bytes.Contains(logBuf.Bytes(), []byte(tt.wantMethod)) {
t.Errorf("log should contain method '%s'", tt.wantMethod)
}
if !bytes.Contains(logBuf.Bytes(), []byte(tt.wantPath)) {
t.Errorf("log should contain path '%s'", tt.wantPath)
}
if !bytes.Contains(logBuf.Bytes(), []byte("duration_ms")) {
t.Error("log should contain duration_ms field")
}
if !bytes.Contains(logBuf.Bytes(), []byte("status")) {
t.Error("log should contain status field")
}

// Verify request ID presence or absence
if tt.wantRequestID {
if !bytes.Contains(logBuf.Bytes(), []byte("request_id")) {
t.Error("log should contain request_id when set")
}
}
})
}
}

func TestLogger_DurationTracking(t *testing.T) {
// Set Gin to test mode
gin.SetMode(gin.TestMode)

// Create a buffer to capture log output
var logBuf bytes.Buffer
logger := slog.New(slog.NewJSONHandler(&logBuf, &slog.HandlerOptions{
Level: slog.LevelInfo,
}))
slog.SetDefault(logger)

// Create router with logging middleware
router := gin.New()
router.Use(Logger())

// Add test handler with known delay
delay := 50 * time.Millisecond
router.GET("/test", func(c *gin.Context) {
time.Sleep(delay)
c.Status(http.StatusOK)
})

// Create and serve request
req := httptest.NewRequest("GET", "/test", nil)
w := httptest.NewRecorder()
router.ServeHTTP(w, req)

// Verify log contains duration tracking
if !bytes.Contains(logBuf.Bytes(), []byte("duration_ms")) {
t.Error("log should contain duration_ms field")
}

// Note: We can't precisely verify the duration value due to timing variations,
// but we've confirmed the field exists and the handler delayed as expected
if w.Code != http.StatusOK {
t.Errorf("expected status 200, got %d", w.Code)
}
}

func TestLogger_StructuredFormat(t *testing.T) {
// Set Gin to test mode
gin.SetMode(gin.TestMode)

// Create a buffer to capture log output
var logBuf bytes.Buffer
logger := slog.New(slog.NewJSONHandler(&logBuf, &slog.HandlerOptions{
Level: slog.LevelInfo,
}))
slog.SetDefault(logger)

// Create router with logging middleware
router := gin.New()
// Simulate T014 request ID middleware setting the ID before Logger
router.Use(func(c *gin.Context) {
c.Set("request_id", "structured-id-456")
c.Next()
})
router.Use(Logger())

// Add test handler
router.GET("/structured-test", func(c *gin.Context) {
c.Status(http.StatusOK)
})

// Create and serve request
req := httptest.NewRequest("GET", "/structured-test", nil)
w := httptest.NewRecorder()
router.ServeHTTP(w, req)

// Verify structured format (JSON)
logOutput := logBuf.String()

// Check for JSON-like structure with key fields
requiredFields := []string{
"\"msg\":\"request completed\"",
"\"method\":\"GET\"",
"\"path\":\"/structured-test\"",
"\"status\":200",
"\"duration_ms\":",
"\"request_id\":\"structured-id-456\"",
}

for _, field := range requiredFields {
if !bytes.Contains(logBuf.Bytes(), []byte(field)) {
t.Errorf("log should contain structured field: %s\nGot: %s", field, logOutput)
}
}
}

func TestLogger_CallsNext(t *testing.T) {
// Set Gin to test mode
gin.SetMode(gin.TestMode)

// Suppress log output for this test
var logBuf bytes.Buffer
logger := slog.New(slog.NewJSONHandler(&logBuf, nil))
slog.SetDefault(logger)

// Create router with logging middleware
router := gin.New()
router.Use(Logger())

handlerCalled := false
router.GET("/test", func(c *gin.Context) {
handlerCalled = true
c.Status(http.StatusOK)
})

// Create and serve request
req := httptest.NewRequest("GET", "/test", nil)
w := httptest.NewRecorder()
router.ServeHTTP(w, req)

// Verify handler was called
if !handlerCalled {
t.Error("Logger middleware should call c.Next() to execute handler")
}
}
