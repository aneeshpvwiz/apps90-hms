package controllers

import (
	"apps90-hms/errors"
	"apps90-hms/initializers"
	"apps90-hms/loggers"
	"apps90-hms/models"
	"apps90-hms/schemas"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateAppointment(c *gin.Context) {

	var appointmentInput schemas.AppointmentInput

	logger := loggers.InitializeLogger()

	// Bind the input JSON to the struct
	if err := c.ShouldBindJSON(&appointmentInput); err != nil {
		logger.Error("Error binding JSON for Create Appointment", "error", err.Error())
		c.Error(models.WrapError(http.StatusBadRequest, errors.ErrBindingJSON, "Invalid request format"))
		return
	}

	// Validate if patient exists
	var patient models.Patient
	initializers.DB.First(&patient, appointmentInput.PatientID)
	if patient.ID == 0 {
		logger.Warn("Patient not found", "patient_id", appointmentInput.PatientID)
		c.Error(models.WrapError(http.StatusNotFound, errors.ErrObjectNotFound, "Patient not found"))
		return
	}

	// Validate if employee (doctor) exists
	var doctor models.Employee
	initializers.DB.First(&doctor, appointmentInput.DoctorID)
	if doctor.ID == 0 {
		logger.Warn("Doctor not found", "employee_id", appointmentInput.DoctorID)
		c.Error(models.WrapError(http.StatusNotFound, errors.ErrObjectNotFound, "Doctor not found"))
		return
	}

	// Validate if entity exists
	var entity models.Entity
	initializers.DB.First(&entity, appointmentInput.EntityID)
	if entity.ID == 0 {
		logger.Warn("Entity not found", "entity_id", appointmentInput.EntityID)
		c.Error(models.WrapError(http.StatusNotFound, errors.ErrObjectNotFound, "Entity not found"))
		return
	}

	// Create the appointment
	appointment := models.Appointment{
		AppointmentTime: appointmentInput.AppointmentTime,
		Reason:          appointmentInput.Reason,
		Notes:           appointmentInput.Notes,
		PatientID:       appointmentInput.PatientID,
		EmployeeID:      appointmentInput.DoctorID,
		EntityID:        appointmentInput.EntityID,
	}

	initializers.DB.Create(&appointment)

	logger.Info("Appointment created successfully", "appointment_id", appointment.ID)

	// Return the created appointment details
	// Return the created appointment details with the custom response format
	c.JSON(http.StatusOK, gin.H{
		"data":    appointment.ID,
		"message": "Successfully created an appointment",
		"status":  "Success",
	})
}

func GetAppointments(c *gin.Context) {
	var appointments []models.Appointment

	// Initialize logger
	logger := loggers.InitializeLogger()

	entityID := c.DefaultQuery("entity_id", "")

	// Build query to filter appointments
	query := initializers.DB.Preload("Patient").Preload("Employee").Preload("Entity")

	// Filter by entity
	if entityID != "" {
		var entity models.Entity
		if err := initializers.DB.First(&entity, entityID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Entity not found", "status": "Failure"})
			return
		}
		query = query.Where("entity_id = ?", entityID)
	}

	// Retrieve appointments from the database
	if err := query.Find(&appointments).Error; err != nil {
		logger.Error("Error fetching appointments", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve appointments", "status": "Failure"})
		return
	}

	// Check if appointments exist
	if len(appointments) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No appointments found", "status": "Success", "data": []string{}})
		return
	}

	// Return the list of appointments with the desired response format
	c.JSON(http.StatusOK, gin.H{
		"data":    appointments,
		"message": "Successfully retrieved appointments",
		"status":  "Success",
	})
}
