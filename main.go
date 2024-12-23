package main

import (
	"apps90-hms/controllers"
	"apps90-hms/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	r := gin.Default()
	r.POST("/posts", controllers.PostsCreate)
	r.GET("/posts", controllers.PostsList)
	r.GET("/posts/:id", controllers.PostDetails)
	r.PUT("/posts/:id", controllers.PostUpdate)
	r.DELETE("/posts/:id", controllers.PostDelete)
	r.Run() // listen and serve on 0.0.0.0:3000
}
