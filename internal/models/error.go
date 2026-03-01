package models

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
Error   string `json:"error"`
Message string `json:"message"`
Code    int    `json:"code"`
}
