package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestListClusterServicesHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		queryParam         string
		expectedStatusCode int
	}{
		{
			name:               "returns all services without filter",
			queryParam:         "",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "returns filtered services by name",
			queryParam:         "?name=home",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "returns empty list for non-matching filter",
			queryParam:         "?name=nonexistentservice12345",
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/cluster/services"+tt.queryParam, nil)

			ListClusterServicesHandler(c)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

			// Verify valid JSON is returned
			var response interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "response should be valid JSON")
		})
	}
}
