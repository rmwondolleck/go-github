package models

import (
	"encoding/json"
	"testing"
)

func TestErrorResponse(t *testing.T) {
	t.Run("JSON marshaling", func(t *testing.T) {
		errResp := ErrorResponse{
			Error:   "test_error",
			Message: "This is a test error",
			Code:    400,
		}

		data, err := json.Marshal(errResp)
		if err != nil {
			t.Fatalf("Failed to marshal ErrorResponse: %v", err)
		}

		expected := `{"error":"test_error","message":"This is a test error","code":400}`
		if string(data) != expected {
			t.Errorf("Expected %s, got %s", expected, string(data))
		}
	})

	t.Run("JSON unmarshaling", func(t *testing.T) {
		jsonData := `{"error":"test_error","message":"This is a test error","code":400}`
		var errResp ErrorResponse

		err := json.Unmarshal([]byte(jsonData), &errResp)
		if err != nil {
			t.Fatalf("Failed to unmarshal ErrorResponse: %v", err)
		}

		if errResp.Error != "test_error" {
			t.Errorf("Expected error 'test_error', got '%s'", errResp.Error)
		}
		if errResp.Message != "This is a test error" {
			t.Errorf("Expected message 'This is a test error', got '%s'", errResp.Message)
		}
		if errResp.Code != 400 {
			t.Errorf("Expected code 400, got %d", errResp.Code)
		}
	})
}
