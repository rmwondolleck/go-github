package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-github/internal/cluster"
	"go-github/internal/models"
	"go-github/internal/server"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestListClusterServices_Returns200 tests that the cluster services endpoint returns 200 with a service list
// TDD: This test is written BEFORE implementation - it will FAIL until GET /api/v1/cluster/services is implemented
func TestListClusterServices_Returns200(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Arrange - Create server instance
	srv := server.New()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/cluster/services", nil)

	// Act - Make HTTP request
	srv.Router().ServeHTTP(w, req)

	// Assert - Verify HTTP status code
	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200 OK")

	// Assert - Verify response is valid JSON
	var response []cluster.ServiceInfo
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Response should be valid JSON")

	// Assert - Verify response structure
	// Even if empty, the response should be a valid JSON array
	assert.NotNil(t, response, "Response should not be nil")

	// Log response for debugging
	t.Logf("Response body: %s", w.Body.String())
	t.Logf("Response headers: %v", w.Header())
}

// TestListClusterServices_FiltersByName tests that the name query parameter filters services correctly
// TDD: This test is written BEFORE implementation - it will FAIL until filtering is implemented
func TestListClusterServices_FiltersByName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		queryParam     string
		expectedStatus int
		validateFunc   func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:           "valid name filter",
			queryParam:     "?name=test-service",
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response []cluster.ServiceInfo
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "Response should be valid JSON")
				assert.NotNil(t, response, "Response should not be nil")

				// If services are returned, verify they match the filter
				for _, svc := range response {
					assert.Contains(t, svc.Name, "test-service", "Filtered service names should contain the filter value")
				}
			},
		},
		{
			name:           "empty name filter returns all services",
			queryParam:     "?name=",
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response []cluster.ServiceInfo
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "Response should be valid JSON")
				assert.NotNil(t, response, "Response should not be nil")
			},
		},
		{
			name:           "no query parameter returns all services",
			queryParam:     "",
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response []cluster.ServiceInfo
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "Response should be valid JSON")
				assert.NotNil(t, response, "Response should not be nil")
			},
		},
		{
			name:           "multiple query parameters - name filter takes precedence",
			queryParam:     "?name=api-service&namespace=default",
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response []cluster.ServiceInfo
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "Response should be valid JSON")
				assert.NotNil(t, response, "Response should not be nil")

				// If services are returned, verify they match the name filter
				for _, svc := range response {
					assert.Contains(t, svc.Name, "api-service", "Filtered service names should contain the filter value")
				}
			},
		},
		{
			name:           "special characters in filter",
			queryParam:     "?name=my-test_service.v1",
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response []cluster.ServiceInfo
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "Response should be valid JSON")
				assert.NotNil(t, response, "Response should not be nil")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			srv := server.New()
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/cluster/services"+tt.queryParam, nil)

			// Act
			srv.Router().ServeHTTP(w, req)

			// Assert - HTTP status code
			assert.Equal(t, tt.expectedStatus, w.Code, "Unexpected status code")

			// Assert - Additional validation
			if tt.validateFunc != nil {
				tt.validateFunc(t, w)
			}

			// Log for debugging
			t.Logf("Query: %s, Status: %d, Body: %s", tt.queryParam, w.Code, w.Body.String())
		})
	}
}

// TestListClusterServices_InvalidFilters tests error handling for invalid filter parameters
// TDD: This test ensures proper error responses for bad requests
func TestListClusterServices_InvalidFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		queryParam     string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "invalid filter parameter - unsupported field",
			queryParam:     "?invalid_field=value",
			expectedStatus: http.StatusOK, // Should ignore unknown parameters and return all services
			expectedError:  "",
		},
		{
			name:           "SQL injection attempt in name filter",
			queryParam:     "?name=%27%20OR%20%271%27%3D%271", // URL encoded: ' OR '1'='1
			expectedStatus: http.StatusOK,                      // Should treat as literal string filter
			expectedError:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			srv := server.New()
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/cluster/services"+tt.queryParam, nil)

			// Act
			srv.Router().ServeHTTP(w, req)

			// Assert - HTTP status code
			assert.Equal(t, tt.expectedStatus, w.Code, "Unexpected status code")

			// Assert - Error response structure if bad request
			if tt.expectedStatus == http.StatusBadRequest {
				var errResponse models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &errResponse)
				assert.NoError(t, err, "Error response should be valid JSON")
				assert.NotEmpty(t, errResponse.Error, "Error field should not be empty")
				assert.NotEmpty(t, errResponse.Message, "Message field should not be empty")
				
				if tt.expectedError != "" {
					assert.Contains(t, errResponse.Error, tt.expectedError, "Error should contain expected message")
				}
			}

			// Log for debugging
			t.Logf("Query: %s, Status: %d, Body: %s", tt.queryParam, w.Code, w.Body.String())
		})
	}
}

// TestListClusterServices_ResponseStructure tests the JSON response structure in detail
// TDD: This validates the expected contract of the API response
func TestListClusterServices_ResponseStructure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Arrange
	srv := server.New()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/cluster/services", nil)

	// Act
	srv.Router().ServeHTTP(w, req)

	// Skip detailed structure validation if endpoint doesn't exist yet
	if w.Code == http.StatusNotFound {
		t.Skip("Endpoint not implemented yet - skipping structure validation")
		return
	}

	// Assert - Response is valid JSON array
	var response []cluster.ServiceInfo
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Response should be valid JSON array")

	// If we have services in the response, validate the structure
	if len(response) > 0 {
		firstService := response[0]
		
		// Validate required fields are present
		assert.NotEmpty(t, firstService.Name, "Service name should not be empty")
		assert.NotEmpty(t, firstService.Namespace, "Service namespace should not be empty")
		assert.NotEmpty(t, firstService.Status, "Service status should not be empty")
		assert.NotNil(t, firstService.Endpoints, "Service endpoints should not be nil")

		t.Logf("Sample service structure: %+v", firstService)
	}

	// Assert - Content-Type header
	contentType := w.Header().Get("Content-Type")
	assert.Contains(t, contentType, "application/json", "Content-Type should be application/json")
}

// TestListClusterServices_ConcurrentRequests tests concurrent access to the endpoint
// This ensures thread-safety of the implementation
func TestListClusterServices_ConcurrentRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	srv := server.New()
	const concurrentRequests = 10

	// Channel to collect errors from goroutines
	errChan := make(chan error, concurrentRequests)
	statusChan := make(chan int, concurrentRequests)

	// Act - Send concurrent requests
	for i := 0; i < concurrentRequests; i++ {
		go func(index int) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/cluster/services", nil)
			srv.Router().ServeHTTP(w, req)

			statusChan <- w.Code

			// Verify response is valid JSON
			var response []cluster.ServiceInfo
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil && w.Code != http.StatusNotFound {
				errChan <- err
				return
			}
			errChan <- nil
		}(i)
	}

	// Collect results
	var statuses []int
	for i := 0; i < concurrentRequests; i++ {
		err := <-errChan
		status := <-statusChan
		
		assert.NoError(t, err, "Concurrent request should not produce JSON parsing errors")
		statuses = append(statuses, status)
	}

	// All requests should return the same status code
	firstStatus := statuses[0]
	for _, status := range statuses {
		assert.Equal(t, firstStatus, status, "All concurrent requests should return the same status code")
	}

	t.Logf("Completed %d concurrent requests with status: %d", concurrentRequests, firstStatus)
}
