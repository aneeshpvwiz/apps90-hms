package middlewares

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

// Custom error definitions
var (
	ErrUserNotFound   = errors.New("user not found")
	ErrDatabaseFailed = errors.New("database connection failed")
)

// APIError represents a structured error for API responses
type APIError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
}

// Implement the error interface for APIError
func (e APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Message, e.Details)
}

// APIErrorMiddleware handles API errors consistently
func APIErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Process the request

		// Handle errors captured during request processing
		for _, ginErr := range c.Errors {
			if apiErr, ok := ginErr.Err.(APIError); ok {
				// Send a structured error response
				c.JSON(apiErr.StatusCode, gin.H{
					"error":   apiErr.Message,
					"details": apiErr.Details,
				})
				return
			}
		}
	}
}
