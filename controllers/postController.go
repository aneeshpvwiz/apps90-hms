package controllers

import (
	"apps90-hms/initializers"
	"apps90-hms/models"

	"github.com/gin-gonic/gin"
)

func PostsCreate(c *gin.Context) {
	//Get data off request body

	var body struct {
		Body  string
		Title string
	}

	c.Bind(&body)

	// Create Post

	post := models.Post{Name: body.Title, Body: body.Body}

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

func PostsList(c *gin.Context) {
	var posts []models.Post

	initializers.DB.Find(&posts)

	c.JSON(200, gin.H{
		"posts": posts,
	})
}

func PostDetails(c *gin.Context) {

	id := c.Param("id")
	var post models.Post

	initializers.DB.Find(&post, id)

	c.JSON(200, gin.H{
		"post": post,
	})
}

func PostUpdate(c *gin.Context) {

	id := c.Param("id")

	var body struct {
		Body  string
		Title string
	}

	c.Bind(&body)

	var post models.Post

	initializers.DB.Find(&post, id)

	initializers.DB.Model(&post).Updates(models.Post{
		Name: body.Title,
		Body: body.Body,
	})

	c.JSON(200, gin.H{
		"post": post,
	})
}

func PostDelete(c *gin.Context) {

	id := c.Param("id")

	initializers.DB.Delete(&models.Post{}, id)

	c.Status(200)
}
