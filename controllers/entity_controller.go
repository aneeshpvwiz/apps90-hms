package controllers

import (
	"apps90-hms/errors"
	"apps90-hms/initializers"
	"apps90-hms/loggers"
	"apps90-hms/models"
	"apps90-hms/schemas"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateEntity(c *gin.Context) {

	var entityInput schemas.EntityInput

	logger := loggers.InitializeLogger()

	if err := c.ShouldBindJSON(&entityInput); err != nil {
		logger.Error("Error binding JSON for Create Entity", "error", err.Error())
		c.Error(models.WrapError(http.StatusBadRequest, errors.ErrBindingJSON, "Invalid request format"))
		return
	}

	var entityFound models.Entity
	initializers.DB.Where("name=?", entityInput.Name).Find(&entityFound)

	if entityFound.ID != 0 {
		logger.Warn("Entity with this name already exists", "name", entityInput.Name)
		c.Error(models.WrapError(http.StatusBadRequest, errors.ErrObjectExists, "Entity with this name already exist"))
		return
	}

	entity := models.Entity{
		Name:    entityInput.Name,
		Address: entityInput.Address,
	}

	initializers.DB.Create(&entity)

	logger.Info("Entity created successfully", "Name", entityInput.Name, "entity_id", entity.ID)

	c.JSON(http.StatusOK, gin.H{"data": entity})

}

func CreateUserEntity(c *gin.Context) {
	var userEntityInput schemas.UserEntityInput

	if err := c.ShouldBindJSON(&userEntityInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userEntity := models.UserEntity{
		UserID:   userEntityInput.UserID,
		EntityID: userEntityInput.EntityID,
	}

	initializers.DB.Create(&userEntity)

	c.JSON(http.StatusOK, gin.H{"data": userEntity})

}
