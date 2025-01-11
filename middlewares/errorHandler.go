package middlewares

import (
	"apps90-hms/models"
	"log"

	"github.com/gin-gonic/gin"
)

// APIErrorMiddleware handles API errors consistently
func APIErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Process the request

		// Check for errors in the Gin context
		for _, ginErr := range c.Errors {
			if apiErr, ok := ginErr.Err.(models.APIError); ok {
				// Send a structured error response
				c.JSON(apiErr.StatusCode, gin.H{
					"code":    apiErr.StatusCode,
					"error":   apiErr.ErrorType,
					"details": apiErr.Message,
				})
				// Stop further processing
				return
			}
		}

		// Debug log if no errors
		log.Println("No errors captured by middleware")
	}
}
