package main

import (
	"apps90-hms/controllers"
	"apps90-hms/initializers"
	"apps90-hms/middlewares"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	r := gin.Default()
	// Auth Urls

	r.POST("/register", controllers.CreateUser)
	r.POST("/login", controllers.Login)
	r.GET("/user/profile", middlewares.CheckAuth, controllers.GetUserProfile)

	// Entity Urls

	r.POST("/entity", controllers.CreateEntity)
	r.POST("/entity/user", controllers.CreateUserEntity)

	r.Run() // listen and serve on 0.0.0.0:3000
}
