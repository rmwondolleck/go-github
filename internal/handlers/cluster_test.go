package handlers

import (
	"encoding/json"
	"go-github/internal/cluster"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestListClusterServicesHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		queryParam         string
		expectedStatusCode int
	}{
		{
			name:               "returns all services",
			queryParam:         "",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "filters services by name",
			queryParam:         "?name=home",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "returns empty list for unknown name",
			queryParam:         "?name=nonexistent-service-xyz",
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/cluster/services"+tt.queryParam, nil)

			ListClusterServicesHandler(c)

			if w.Code != tt.expectedStatusCode {
				t.Errorf("expected status code %d, got %d", tt.expectedStatusCode, w.Code)
			}

			var response []cluster.ServiceInfo
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}
		})
	}
}
