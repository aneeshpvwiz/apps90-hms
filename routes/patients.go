package routes

import (
	patientController "apps90-hms/controllers/patient"

	"github.com/gin-gonic/gin"
)

func PatientRoutes(r *gin.Engine) {
	patient := r.Group("/patient")
	{
		patient.GET("/details", patientController.GetPatientDetails)
		patient.GET("/visits", patientController.GetPatientVisitHistory)
		patient.POST("/prescription", patientController.CreatePrescription)
		patient.GET("/prescription", patientController.GetPatientPrescriptions)

	}
}
