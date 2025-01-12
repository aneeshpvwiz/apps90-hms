package models

import (
	"fmt"
)

// APIError represents a structured error for API responses
type APIError struct {
	StatusCode int    `json:"-"`
	ErrorType  string `json:"details,omitempty"`
	Message    string `json:"message"`
}

// Implement the error interface for APIError
func (e APIError) Error() string {
	return fmt.Sprintf("%s", e.Message)
}

// WrapError wraps an error with additional context
func WrapError(statusCode int, err error, message string) APIError {
	return APIError{
		StatusCode: statusCode,
		ErrorType:  err.Error(),
		Message:    message,
	}
}
