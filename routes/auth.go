package routes

import (
	"apps90-hms/controllers"
	"apps90-hms/middlewares"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", controllers.CreateUser)
		auth.POST("/login", controllers.Login)
		auth.GET("/profile", middlewares.CheckAuth, controllers.GetUserProfile)
	}
}
