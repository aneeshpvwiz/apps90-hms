package controllers

import (
	"apps90-hms/initializers"
	"apps90-hms/models"
	"apps90-hms/schemas"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateEntity(c *gin.Context) {

	var entityInput schemas.EntityInput

	if err := c.ShouldBindJSON(&entityInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var entityFound models.Entity
	initializers.DB.Where("name=?", entityInput.Name).Find(&entityFound)

	if entityFound.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Entity with this name already exist"})
		return
	}

	entity := models.Entity{
		Name:    entityInput.Name,
		Address: entityInput.Address,
	}

	initializers.DB.Create(&entity)

	c.JSON(http.StatusOK, gin.H{"data": entity})

}

func CreateUserEntity(c *gin.Context) {
	var userEntityInput schemas.UserEntityInput

	if err := c.ShouldBindJSON(&userEntityInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("#######", userEntityInput)

	userEntity := models.UserEntity{
		UserID:   userEntityInput.UserID,
		EntityID: userEntityInput.EntityID,
	}

	initializers.DB.Create(&userEntity)

	c.JSON(http.StatusOK, gin.H{"data": userEntity})

}
