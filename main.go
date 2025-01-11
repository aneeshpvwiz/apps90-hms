package main

import (
	"apps90-hms/initializers"
	"apps90-hms/loggers"
	"apps90-hms/routes"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()

}

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

// WrapError wraps an error with additional context
func WrapError(statusCode int, err error, message string) APIError {
	return APIError{
		StatusCode: statusCode,
		Message:    message,
		Details:    err.Error(),
	}
}

// Middleware to handle API errors
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

// Mock function simulating a database operation
func getUserFromDB(userID string) (string, error) {
	if userID == "1" {
		return "John Doe", nil
	}
	if userID == "2" {
		// Wrapping database error with additional context
		return "", fmt.Errorf("querying user ID %s: %w", userID, ErrDatabaseFailed)
	}
	// Wrapping "user not found" error
	return "", fmt.Errorf("querying user ID %s: %w", userID, ErrUserNotFound)
}

func main() {
	// Initialize the logger
	loggers.InitLogger()

	// Get the logger instance
	log := loggers.GetLogger()

	log.Info("Application started")

	router := routes.InitRoutes()

	router.Use(APIErrorMiddleware())

	router.GET("/user/:id", func(c *gin.Context) {
		userID := c.Param("id")
		user, err := getUserFromDB(userID)

		if err != nil {
			// Register the error using the context
			if errors.Is(err, ErrUserNotFound) {
				c.Error(WrapError(http.StatusNotFound, err, "User not found"))
			} else {
				c.Error(WrapError(http.StatusInternalServerError, err, "Database error"))
			}
			return
		}

		// Respond with the user data
		c.JSON(http.StatusOK, gin.H{"user": user})
	})

	router.Run() // listen and serve on 0.0.0.0:3000
}
