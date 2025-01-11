package controllers

import (
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userFound models.User
	initializers.DB.Where("email=?", authInput.Email).Find(&userFound)

	if userFound.ID != 0 {
		logger.Warn("User with this email already exists", "email", authInput.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with this email already exist"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Error generating password hash", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Email:    authInput.Email,
		Password: string(passwordHash),
	}

	initializers.DB.Create(&user)

	logger.Info("User created successfully", "email", authInput.Email, "user_id", user.ID)

	c.JSON(http.StatusOK, gin.H{"data": user})

}

func Login(c *gin.Context) {

	var authInput schemas.AuthInput

	// Log incoming request
	logger := loggers.InitializeLogger()

	if err := c.ShouldBindJSON(&authInput); err != nil {
		logger.Error("Error binding JSON for Login", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userFound models.User
	initializers.DB.Where("email=?", authInput.Email).Find(&userFound)

	if userFound.ID == 0 {
		logger.Warn("User not found", "email", authInput.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(authInput.Password)); err != nil {
		logger.Warn("Invalid password attempt", "email", authInput.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userFound.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := generateToken.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		logger.Error("Failed to generate token", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to generate token"})
	}
	logger.Info("User logged in successfully", "user_id", userFound.ID)

	c.JSON(200, gin.H{
		"token": token,
	})
}

func GetUserProfile(c *gin.Context) {

	user, _ := c.Get("currentUser")

	c.JSON(200, gin.H{
		"user": user,
	})
}
