package main

import (
	"apps90-hms/initializers"
	"apps90-hms/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()

}

func main() {
	// Migrate the schema
	initializers.DB.AutoMigrate(&models.Post{})
}
