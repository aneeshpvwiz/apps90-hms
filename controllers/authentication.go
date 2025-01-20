package controllers

import (
	"apps90-hms/errors"
	"apps90-hms/initializers"
	"apps90-hms/loggers"
	"apps90-hms/models"
	"apps90-hms/schemas"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {

	var authInput schemas.AuthInput

	// Log incoming request
	logger := loggers.InitializeLogger()

	if err := c.ShouldBindJSON(&authInput); err != nil {
		logger.Error("Error binding JSON for CreateUser", "error", err.Error())
		c.Error(models.WrapError(http.StatusBadRequest, errors.ErrBindingJSON, "Invalid request format"))
		return
	}

	var userFound models.User
	initializers.DB.Where("email=?", authInput.Email).Find(&userFound)

	if userFound.ID != 0 {
		logger.Warn("User with this email already exists", "email", authInput.Email)
		c.Error(models.WrapError(http.StatusBadRequest, errors.ErrUserExists, "Email already registered: "+authInput.Email))
		return

	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Error generating password hash", "error", err.Error())
		c.Error(models.WrapError(http.StatusBadRequest, errors.ErrHashingPassword, "Password hashing error"))
		return
	}

	user := models.User{
		Email:    authInput.Email,
		Password: string(passwordHash),
	}

	initializers.DB.Create(&user)

	logger.Info("User created successfully", "email", authInput.Email, "user_id", user.ID)

	// Sanitize user data: Remove sensitive and unnecessary fields
	userResponse := gin.H{
		"id":    user.ID,
		"email": user.Email,
	}

	// Return success response with sanitized user data
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User created successfully",
		"data":    userResponse,
	})

}

func Login(c *gin.Context) {

	var authInput schemas.AuthInput

	// Log incoming request
	logger := loggers.InitializeLogger()

	if err := c.ShouldBindJSON(&authInput); err != nil {
		logger.Error("Error binding JSON for Login", "error", err.Error())
		c.Error(models.WrapError(http.StatusBadRequest, errors.ErrBindingJSON, "Invalid request format"))
		return
	}

	var userFound models.User
	initializers.DB.Where("email=?", authInput.Email).Find(&userFound)

	if userFound.ID == 0 {
		logger.Warn("User not found", "email", authInput.Email)
		c.Error(models.WrapError(http.StatusNotFound, errors.ErrUserNotFound, "User not found"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(authInput.Password)); err != nil {
		logger.Warn("Invalid password attempt", "email", authInput.Email)
		c.Error(models.WrapError(http.StatusBadRequest, errors.ErrInvalidPassword, "Incorrect password"))
		return
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userFound.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := generateToken.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		logger.Error("Failed to generate token", "error", err.Error())
		c.Error(models.WrapError(http.StatusInternalServerError, errors.ErrGeneratingToken, "Failed to generate token"))
	}
	logger.Info("User logged in successfully", "user_id", userFound.ID)

	// Return success response with sanitized user data
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Succesfully validated the user",
		"data":    token,
	})
}

func GetUserProfile(c *gin.Context) {

	// Initialize the logger
	logger := loggers.InitializeLogger()

	// Retrieve the current user from the context
	user, ok := c.Get("currentUser")
	if !ok {
		logger.Warn("User not found in context", "context_key", "currentUser")
		c.Error(models.WrapError(http.StatusNotFound, errors.ErrUserNotFound, "User not found"))
		return
	}

	// Type assertion to your user struct
	userData, ok := user.(models.User)
	if !ok {
		logger.Error("Failed to assert user from context", "context_key", "currentUser", "user_data", user)
		c.Error(models.WrapError(http.StatusInternalServerError, errors.InternalServerError, "Internal server error - user type assertion failed"))
		return
	}

	// Log successful user retrieval
	logger.Info("User profile found", "user_id", userData.ID, "email", userData.Email)

	// preload the related Entities
	var userWithEntities models.User
	if err := initializers.DB.Preload("Entities").First(&userWithEntities, userData.ID).Error; err != nil {
		logger.Error("Failed to load user and entities", "user_id", userData.ID, "error", err.Error())
		c.Error(models.WrapError(http.StatusInternalServerError, errors.ErrDatabaseFailed, "Failed to load user and entities"))
		return
	}

	// Prepare the response object with the user and entities data
	var entities []gin.H
	for _, entity := range userWithEntities.Entities {
		entities = append(entities, gin.H{
			"id":      entity.ID,
			"name":    entity.Name,
			"address": entity.Address,
		})
	}

	// Log success when user profile and entities are loaded
	logger.Info("User profile and entities loaded", "user_id", userData.ID, "entities_count", len(entities))

	// Create a response object with only non-sensitive fields
	response := gin.H{
		"id":       userData.ID,
		"email":    userData.Email,
		"entities": entities,
	}

	c.JSON(http.StatusOK, response)
}
