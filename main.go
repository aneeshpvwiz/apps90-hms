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

	r.POST("/register", controllers.CreateUser)
	r.POST("/login", controllers.Login)
	r.GET("/user/profile", middlewares.CheckAuth, controllers.GetUserProfile)
	r.POST("/posts", controllers.PostsCreate)
	r.GET("/posts", controllers.PostsList)
	r.GET("/posts/:id", controllers.PostDetails)
	r.PUT("/posts/:id", controllers.PostUpdate)
	r.DELETE("/posts/:id", controllers.PostDelete)
	r.Run() // listen and serve on 0.0.0.0:3000
}
