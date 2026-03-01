package handlers

import (
"go-github/internal/models"
"net/http"

"github.com/gin-gonic/gin"
)

// JSONSuccess sends a successful JSON response
func JSONSuccess(c *gin.Context, code int, data interface{}) {
c.JSON(code, data)
}

// JSONError sends an error JSON response using the ErrorResponse model
func JSONError(c *gin.Context, code int, err string, message string) {
c.JSON(code, models.ErrorResponse{
Error:   err,
Message: message,
Code:    code,
})
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
