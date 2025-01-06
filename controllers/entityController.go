package controllers

import (
	"apps90-hms/initializers"
	"apps90-hms/models"

	"github.com/gin-gonic/gin"
)

func CreateEntity(c *gin.Context) {
	//Get data off request body

	var body struct {
		name    string
		address string
	}

	c.Bind(&body)

	// Create Post

	post := models.Entity{Name: body.name, Address: body.address}

	result := initializers.DB.Create(&post)

	if result.Error != nil {
		c.Status(400)
		return
	}

	// Return the service

	c.JSON(200, gin.H{
		"message": post,
	})
}
