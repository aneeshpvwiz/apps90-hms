package main

import (
	"apps90-hms/initializers"
	"apps90-hms/loggers"
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

	router.Run() // listen and serve on 0.0.0.0:3000
}
