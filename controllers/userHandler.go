package controllers

import (
	"apps90-hms/models"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Mocked errors
var (
	ErrUserNotFound   = errors.New("user not found")
	ErrDatabaseFailed = errors.New("database connection failed")
)

// Mock function simulating a database operation
func getUserFromDB(userID string) (string, error) {
	if userID == "1" {
		return "John Doe", nil
	}
	if userID == "2" {
		return "", models.WrapError(http.StatusInternalServerError, ErrDatabaseFailed, "Database error", "Database Error")
	}
	return "", models.WrapError(http.StatusNotFound, ErrUserNotFound, "User not found", fmt.Sprintf("User with ID %s does not exist", userID))
}

// GetUserHandler handles GET requests to fetch a user by ID
func GetUserHandler(c *gin.Context) {
	userID := c.Param("id")
	user, err := getUserFromDB(userID)

	if err != nil {
		// Register the error in the context
		c.Error(err)
		return
	}

	// Return the user data if no errors
	c.JSON(http.StatusOK, gin.H{"user": user})
}
