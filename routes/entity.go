package routes

import (
	"apps90-hms/controllers"

	"github.com/gin-gonic/gin"
)

func EntityRoutes(r *gin.Engine) {
	entity := r.Group("/entity")
	{
		entity.POST("/", controllers.CreateEntity)
		entity.POST("/user", controllers.CreateUserEntity)
		entity.POST("/employee", controllers.AddEmployee)
		entity.GET("/employee", controllers.GetEmployeeList)
		entity.POST("/patient", controllers.AddPatient)
		entity.GET("/patient", controllers.GetPatientList)
		entity.POST("/appointment", controllers.CreateAppointment)
		entity.GET("/appointment", controllers.GetAppointments)

	}
}
