package routes

import (
	appointmentControllers "apps90-hms/controllers/appointment"
	entityController "apps90-hms/controllers/entity"

	"github.com/gin-gonic/gin"
)

func EntityRoutes(r *gin.Engine) {
	entity := r.Group("/entity")
	{
		entity.POST("/", entityController.CreateEntity)
		entity.POST("/user", entityController.CreateUserEntity)
		entity.POST("/employee", entityController.AddEmployee)
		entity.GET("/employee", entityController.GetEmployeeList)
		entity.POST("/patient", entityController.AddPatient)
		entity.GET("/patient", entityController.GetPatientList)
		entity.PUT("/patients/:id", entityController.EditPatient)
		entity.POST("/appointment", appointmentControllers.CreateAppointment)
		entity.GET("/appointment", appointmentControllers.GetAppointments)
		entity.POST("/visit", appointmentControllers.CreateVisit)
		entity.PUT("/visit", appointmentControllers.EditVisit)
		entity.GET("/medicine", entityController.GetMedicines)
		entity.POST("/medicine", entityController.AddMedicine)
		entity.POST("/category", entityController.AddMedicineCategory)

	}
}
