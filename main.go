package main

import (
	"apps90-hms/initializers"
	"apps90-hms/routes"

	"apps90-hms/loggers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()

	gin.SetMode(gin.ReleaseMode)

}

func main() {

	logger := loggers.InitializeLogger()

	logger.Info("Application started")

	router := routes.InitRoutes()

	// Use the custom error-handling middleware
	//router.Use(middlewares.APIErrorMiddleware())

	// Start the server
	logger.Info("Starting server", "address", ":8080")
	if err := router.Run(); err != nil {
		logger.Error("Failed to start server: %v", err)
	}

	//router.Run() // listen and serve on 0.0.0.0:3000
}
