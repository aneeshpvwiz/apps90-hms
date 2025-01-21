package routes

import (
	authControllers "apps90-hms/controllers/auth"
	"apps90-hms/middlewares"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", authControllers.CreateUser)
		auth.POST("/login", authControllers.Login)
		auth.GET("/profile", middlewares.CheckAuth, authControllers.GetUserProfile)
	}
}
