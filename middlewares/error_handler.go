package middlewares

import (
	"apps90-hms/models"
	"log"

	"github.com/gin-gonic/gin"
)

// APIErrorMiddleware handles API errors consistently
func APIErrorMiddleware() gin.HandlerFunc {
	log.Printf("Reached here ##############")
	return func(c *gin.Context) {
		c.Next() // Process the request

		log.Printf("Reached here ##############")

		// Check for errors in the Gin context
		for _, ginErr := range c.Errors {
			if apiErr, ok := ginErr.Err.(models.APIError); ok {
				// Send a structured error response
				c.JSON(apiErr.StatusCode, gin.H{
					"code":    apiErr.StatusCode,
					"error":   apiErr.ErrorType,
					"message": apiErr.Message,
				})

				// Log the error for debugging purposes
				log.Printf("API Error ##############: %v", apiErr)
				// Stop further processing
				return
			}
		}

		// Debug log if no errors
		log.Println("No errors captured by middleware")
	}
}
