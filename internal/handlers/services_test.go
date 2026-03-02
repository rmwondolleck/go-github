package handlers

import (
	"encoding/json"
	"go-github/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestListServicesHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		expectedStatusCode int
		expectedCount      int
	}{
		{
			name:               "returns list of services",
			expectedStatusCode: http.StatusOK,
			expectedCount:      5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			ListServicesHandler(c)

			if w.Code != tt.expectedStatusCode {
				t.Errorf("expected status code %d, got %d", tt.expectedStatusCode, w.Code)
			}

			var response models.ServicesResponse
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}

			if len(response.Services) != tt.expectedCount {
				t.Errorf("expected %d services, got %d", tt.expectedCount, len(response.Services))
			}
		})
	}
}

func TestListServicesHandler_ResponseStructure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ListServicesHandler(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response models.ServicesResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response.Services == nil {
		t.Error("expected services to not be nil")
	}

	// Verify at least some services are returned
	if len(response.Services) == 0 {
		t.Error("expected at least one service to be returned")
	}

	// Verify structure of first service
	if len(response.Services) > 0 {
		service := response.Services[0]

		if service.Name == "" {
			t.Error("expected service name to not be empty")
		}

		if service.Type == "" {
			t.Error("expected service type to not be empty")
		}

		if service.Status == "" {
			t.Error("expected service status to not be empty")
		}

		if service.Endpoint == "" {
			t.Error("expected service endpoint to not be empty")
		}
	}
}

func TestListServicesHandler_SpecificServices(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ListServicesHandler(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response models.ServicesResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	// Check for specific expected services
	expectedServices := map[string]bool{
		"homeassistant": false,
		"prometheus":    false,
		"grafana":       false,
	}

	for _, service := range response.Services {
		if _, exists := expectedServices[service.Name]; exists {
			expectedServices[service.Name] = true
		}
	}

	for serviceName, found := range expectedServices {
		if !found {
			t.Errorf("expected service '%s' not found in response", serviceName)
		}
	}
}

func TestListServicesHandler_JSONFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ListServicesHandler(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json; charset=utf-8" {
		t.Errorf("expected content type 'application/json; charset=utf-8', got '%s'", contentType)
	}

	// Verify it's valid JSON
	var response interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
}

func TestListServicesHandler_ServiceFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ListServicesHandler(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response models.ServicesResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	// Test each service has all required fields
	for i, service := range response.Services {
		t.Run("service_"+service.Name, func(t *testing.T) {
			if service.Name == "" {
				t.Errorf("service at index %d has empty name", i)
			}

			if service.Type == "" {
				t.Errorf("service '%s' has empty type", service.Name)
			}

			if service.Status == "" {
				t.Errorf("service '%s' has empty status", service.Name)
			}

			if service.Endpoint == "" {
				t.Errorf("service '%s' has empty endpoint", service.Name)
			}

			// Validate endpoint format (should start with http:// or https://)
			if len(service.Endpoint) >= 4 {
				if service.Endpoint[:4] != "http" {
					t.Errorf("service '%s' endpoint '%s' does not start with http", service.Name, service.Endpoint)
				}
			} else if len(service.Endpoint) > 0 {
				t.Errorf("service '%s' endpoint '%s' is too short to be a valid URL", service.Name, service.Endpoint)
			}
		})
	}
}
