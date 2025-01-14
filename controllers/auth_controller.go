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

	user, _ := c.Get("currentUser")

	c.JSON(200, gin.H{
		"user": user,
	})
}
