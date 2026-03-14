package handlers

import (
	"go-github/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

// jsonAPI is the jsoniter API instance configured for fastest performance
var jsonAPI = jsoniter.ConfigFastest

// JSONSuccess sends a successful JSON response using jsoniter for improved performance
func JSONSuccess(c *gin.Context, code int, data interface{}) {
	bytes, err := jsonAPI.Marshal(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "encoding_error",
			Message: "Failed to encode response",
			Code:    http.StatusInternalServerError,
		})
		return
	}
	c.Data(code, "application/json; charset=utf-8", bytes)
}

// JSONError sends an error JSON response using the ErrorResponse model with jsoniter
func JSONError(c *gin.Context, code int, err string, message string) {
	errorResponse := models.ErrorResponse{
		Error:   err,
		Message: message,
		Code:    code,
	}
	bytes, marshalErr := jsonAPI.Marshal(errorResponse)
	if marshalErr != nil {
		// Fallback to gin's JSON if jsoniter fails
		c.JSON(code, errorResponse)
		return
	}
	c.Data(code, "application/json; charset=utf-8", bytes)
}

// NotFound sends a 404 Not Found error response
func NotFound(c *gin.Context, message string) {
JSONError(c, http.StatusNotFound, "not_found", message)
}

// BadRequest sends a 400 Bad Request error response
func BadRequest(c *gin.Context, message string) {
JSONError(c, http.StatusBadRequest, "bad_request", message)
}

// InternalError sends a 500 Internal Server Error response
func InternalError(c *gin.Context, message string) {
JSONError(c, http.StatusInternalServerError, "internal_error", message)
}
