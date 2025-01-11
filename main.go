package main

import (
	"apps90-hms/controllers"
	"apps90-hms/initializers"
	"apps90-hms/loggers"
	"apps90-hms/middlewares"
	"apps90-hms/routes"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()

}

func main() {
	// Initialize the logger
	loggers.InitLogger()

	// Get the logger instance
	log := loggers.GetLogger()

	log.Info("Application started")

	router := routes.InitRoutes()

	// Use the custom error-handling middleware
	router.Use(middlewares.APIErrorMiddleware())

	// Define routes
	router.GET("/user/:id", controllers.GetUserHandler)

	// Start the server
	log.Println("Starting server on :8080")
	if err := router.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	//router.Run() // listen and serve on 0.0.0.0:3000
}
