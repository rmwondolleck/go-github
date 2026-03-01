package handlers

import (
"encoding/json"
"go-github/internal/models"
"net/http"
"net/http/httptest"
"testing"

"github.com/gin-gonic/gin"
)

func TestJSONSuccess(t *testing.T) {
gin.SetMode(gin.TestMode)

tests := []struct {
name       string
statusCode int
data       interface{}
expected   string
}{
{
name:       "success with map data",
statusCode: http.StatusOK,
data:       map[string]string{"message": "success"},
expected:   `{"message":"success"}`,
},
{
name:       "success with struct data",
statusCode: http.StatusCreated,
data:       struct {
ID   int    `json:"id"`
Name string `json:"name"`
}{ID: 1, Name: "test"},
expected: `{"id":1,"name":"test"}`,
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
w := httptest.NewRecorder()
c, _ := gin.CreateTestContext(w)

JSONSuccess(c, tt.statusCode, tt.data)

if w.Code != tt.statusCode {
t.Errorf("expected status code %d, got %d", tt.statusCode, w.Code)
}

var got map[string]interface{}
if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
t.Fatalf("failed to unmarshal response: %v", err)
}

var expected map[string]interface{}
if err := json.Unmarshal([]byte(tt.expected), &expected); err != nil {
t.Fatalf("failed to unmarshal expected: %v", err)
}

// Compare JSON objects
gotJSON, _ := json.Marshal(got)
expectedJSON, _ := json.Marshal(expected)
if string(gotJSON) != string(expectedJSON) {
t.Errorf("expected body %s, got %s", expectedJSON, gotJSON)
}
})
}
}

func TestJSONError(t *testing.T) {
gin.SetMode(gin.TestMode)

tests := []struct {
name       string
statusCode int
err        string
message    string
}{
{
name:       "not found error",
statusCode: http.StatusNotFound,
err:        "not_found",
message:    "Resource not found",
},
{
name:       "bad request error",
statusCode: http.StatusBadRequest,
err:        "bad_request",
message:    "Invalid input",
},
{
name:       "internal server error",
statusCode: http.StatusInternalServerError,
err:        "internal_error",
message:    "Something went wrong",
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
w := httptest.NewRecorder()
c, _ := gin.CreateTestContext(w)

JSONError(c, tt.statusCode, tt.err, tt.message)

if w.Code != tt.statusCode {
t.Errorf("expected status code %d, got %d", tt.statusCode, w.Code)
}

var response models.ErrorResponse
if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
t.Fatalf("failed to unmarshal response: %v", err)
}

if response.Error != tt.err {
t.Errorf("expected error %s, got %s", tt.err, response.Error)
}

if response.Message != tt.message {
t.Errorf("expected message %s, got %s", tt.message, response.Message)
}

if response.Code != tt.statusCode {
t.Errorf("expected code %d, got %d", tt.statusCode, response.Code)
}
})
}
}

func TestNotFound(t *testing.T) {
gin.SetMode(gin.TestMode)

w := httptest.NewRecorder()
c, _ := gin.CreateTestContext(w)

message := "User not found"
NotFound(c, message)

if w.Code != http.StatusNotFound {
t.Errorf("expected status code %d, got %d", http.StatusNotFound, w.Code)
}

var response models.ErrorResponse
if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
t.Fatalf("failed to unmarshal response: %v", err)
}

if response.Error != "not_found" {
t.Errorf("expected error 'not_found', got %s", response.Error)
}

if response.Message != message {
t.Errorf("expected message %s, got %s", message, response.Message)
}

if response.Code != http.StatusNotFound {
t.Errorf("expected code %d, got %d", http.StatusNotFound, response.Code)
}
}

func TestBadRequest(t *testing.T) {
gin.SetMode(gin.TestMode)

w := httptest.NewRecorder()
c, _ := gin.CreateTestContext(w)

message := "Invalid request parameters"
BadRequest(c, message)

if w.Code != http.StatusBadRequest {
t.Errorf("expected status code %d, got %d", http.StatusBadRequest, w.Code)
}

var response models.ErrorResponse
if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
t.Fatalf("failed to unmarshal response: %v", err)
}

if response.Error != "bad_request" {
t.Errorf("expected error 'bad_request', got %s", response.Error)
}

if response.Message != message {
t.Errorf("expected message %s, got %s", message, response.Message)
}

if response.Code != http.StatusBadRequest {
t.Errorf("expected code %d, got %d", http.StatusBadRequest, response.Code)
}
}

func TestInternalError(t *testing.T) {
gin.SetMode(gin.TestMode)

w := httptest.NewRecorder()
c, _ := gin.CreateTestContext(w)

message := "Database connection failed"
InternalError(c, message)

if w.Code != http.StatusInternalServerError {
t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, w.Code)
}

var response models.ErrorResponse
if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
t.Fatalf("failed to unmarshal response: %v", err)
}

if response.Error != "internal_error" {
t.Errorf("expected error 'internal_error', got %s", response.Error)
}

if response.Message != message {
t.Errorf("expected message %s, got %s", message, response.Message)
}

if response.Code != http.StatusInternalServerError {
t.Errorf("expected code %d, got %d", http.StatusInternalServerError, response.Code)
}
}
