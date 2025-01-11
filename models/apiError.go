package models

import "fmt"

// APIError represents a structured error for API responses
type APIError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	ErrorType  string `json:"details,omitempty"`
}

// Implement the error interface for APIError
func (e APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Message, e.Message)
}

// WrapError wraps an error with additional context
func WrapError(statusCode int, err error, message string, details string) APIError {
	return APIError{
		StatusCode: statusCode,
		Message:    fmt.Sprintf("%s: %s", message, details),
		ErrorType:  err.Error(),
	}
}
