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
	initializers.DB.AutoMigrate(&models.User{})
	initializers.DB.AutoMigrate(&models.Entity{})
	initializers.DB.AutoMigrate(&models.UserEntity{})
	initializers.DB.AutoMigrate(&models.EmployeeCategory{})
	initializers.DB.AutoMigrate(&models.Employee{})
	initializers.DB.AutoMigrate(&models.Patient{})
	initializers.DB.AutoMigrate(&models.Appointment{})
	initializers.DB.AutoMigrate(&models.InpatientVisit{})
	initializers.DB.AutoMigrate(&models.OutpatientVisit{})
	initializers.DB.AutoMigrate(&models.MedicineCategory{})
	initializers.DB.AutoMigrate(&models.Medicine{})
}
