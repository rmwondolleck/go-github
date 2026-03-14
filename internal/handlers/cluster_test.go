package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-github/internal/cluster"

	"github.com/gin-gonic/gin"
)

func TestListClusterServicesHandler_ReturnsOK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/cluster/services", nil)

	ListClusterServicesHandler(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestListClusterServicesHandler_ValidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/cluster/services", nil)

	ListClusterServicesHandler(c)

	var response []cluster.ServiceInfo
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	if response == nil {
		t.Error("expected non-nil response array")
	}
}

func TestListClusterServicesHandler_ReturnsServices(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		queryParam    string
		expectedCount int
	}{
		{
			name:          "no filter returns all services",
			queryParam:    "",
			expectedCount: 3,
		},
		{
			name:          "filter by api returns one service",
			queryParam:    "api",
			expectedCount: 1,
		},
		{
			name:          "filter by database returns one service",
			queryParam:    "database",
			expectedCount: 1,
		},
		{
			name:          "filter with no match returns empty list",
			queryParam:    "nonexistent",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest("GET", "/api/v1/cluster/services", nil)
			if tt.queryParam != "" {
				c.Request.URL.RawQuery = "name=" + tt.queryParam
			}

			ListClusterServicesHandler(c)

			if w.Code != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
			}

			var response []cluster.ServiceInfo
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}

			if len(response) != tt.expectedCount {
				t.Errorf("expected %d services, got %d", tt.expectedCount, len(response))
			}
		})
	}
}

func TestListClusterServicesHandler_FilterCaseInsensitive(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/cluster/services", nil)
	c.Request.URL.RawQuery = "name=API"

	ListClusterServicesHandler(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response []cluster.ServiceInfo
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(response) != 1 {
		t.Errorf("expected 1 service for case-insensitive filter 'API', got %d", len(response))
	}

	if len(response) > 0 && response[0].Name != "api-service" {
		t.Errorf("expected service name 'api-service', got '%s'", response[0].Name)
	}
}

func TestListClusterServicesHandler_EmptyFilterReturnsAll(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		rawQuery   string
	}{
		{name: "no query param", rawQuery: ""},
		{name: "empty name param", rawQuery: "name="},
		{name: "whitespace name param", rawQuery: "name=   "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/api/v1/cluster/services", nil)
			if tt.rawQuery != "" {
				c.Request.URL.RawQuery = tt.rawQuery
			}

			ListClusterServicesHandler(c)

			if w.Code != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
			}

			var response []cluster.ServiceInfo
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}

			if len(response) != 3 {
				t.Errorf("expected 3 services, got %d", len(response))
			}
		})
	}
}

func TestListClusterServicesHandler_ServiceStructure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/cluster/services", nil)

	ListClusterServicesHandler(c)

	var response []cluster.ServiceInfo
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	for i, svc := range response {
		if svc.Name == "" {
			t.Errorf("service at index %d has empty name", i)
		}
		if svc.Namespace == "" {
			t.Errorf("service '%s' has empty namespace", svc.Name)
		}
		if svc.Status == "" {
			t.Errorf("service '%s' has empty status", svc.Name)
		}
		if svc.Endpoints == nil {
			t.Errorf("service '%s' has nil endpoints", svc.Name)
		}
	}
}

func TestListClusterServicesHandler_ContentType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/cluster/services", nil)

	ListClusterServicesHandler(c)

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json; charset=utf-8" {
		t.Errorf("expected content type 'application/json; charset=utf-8', got '%s'", contentType)
	}
}
