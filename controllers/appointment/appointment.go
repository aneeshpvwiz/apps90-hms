package appointmentControllers

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
		query = query.Where("entity_id = ? AND is_active = ?", entityID, true)
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

	// Prepare the filtered response
	var appointmentResponses []map[string]interface{}
	for _, appointment := range appointments {
		appointmentResponses = append(appointmentResponses, map[string]interface{}{
			"appointment_id":    appointment.ID,
			"appointment_time":  appointment.AppointmentTime,
			"reason":            appointment.Reason,
			"notes":             appointment.Notes,
			"patient_firstname": appointment.Patient.FirstName,
			"patient_lastname":  appointment.Patient.LastName,
			"patient_gender":    appointment.Patient.Gender,
			"patient_dob":       appointment.Patient.DateOfBirth,
			"doctor_firstname":  appointment.Employee.FirstName,
			"doctor_lastname":   appointment.Employee.LastName,
		})
	}

	// Return the list of appointments with the desired response format
	c.JSON(http.StatusOK, gin.H{
		"data":    appointmentResponses,
		"message": "Successfully retrieved appointments",
		"status":  "Success",
	})
}

func CreateVisit(c *gin.Context) {
	var input schemas.VisitInput // Use your appropriate input schema
	logger := loggers.InitializeLogger()

	// Bind the request body
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("Error binding JSON for Visit", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request format", "status": "Error"})
		return
	}

	// Verify patient exists
	var patient models.Patient
	if err := initializers.DB.First(&patient, input.PatientID).Error; err != nil {
		logger.Warn("Patient not found", "patient_id", input.PatientID)
		c.JSON(http.StatusNotFound, gin.H{"message": "Patient not found", "status": "Error"})
		return
	}

	// Verify doctor exists
	var doctor models.Employee
	if err := initializers.DB.First(&doctor, input.DoctorID).Error; err != nil {
		logger.Warn("Doctor not found", "doctor_id", input.DoctorID)
		c.JSON(http.StatusNotFound, gin.H{"message": "Doctor not found", "status": "Error"})
		return
	}

	// Verify optional appointment if provided
	var appointment *models.Appointment
	if input.AppointmentID != nil {
		if err := initializers.DB.First(&appointment, *input.AppointmentID).Error; err != nil {
			logger.Warn("Appointment not found", "appointment_id", *input.AppointmentID)
			c.JSON(http.StatusNotFound, gin.H{"message": "Appointment not found", "status": "Error"})
			return
		}
	}

	// Create the visit (Inpatient or Outpatient)
	visit := models.Visit{
		AppointmentID: input.AppointmentID,
		VisitDate:     input.VisitDate,
		AdmissionDate: input.AdmissionDate,
		DischargeDate: input.DischargeDate,
		RoomNumber:    input.RoomNumber,
		Diagnosis:     input.Diagnosis,
		TreatmentPlan: input.TreatmentPlan,
		Notes:         input.Notes,
		PatientID:     input.PatientID,
		DoctorID:      input.DoctorID,
		VisitType:     input.VisitType, // Use the 'visit_type' provided in input
	}

	// Check if visit type is correct (inpatient or outpatient)
	if visit.VisitType != "IP" && visit.VisitType != "OP" {
		logger.Warn("Invalid visit type", "visit_type", visit.VisitType)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid visit type. Should be 'IP' or 'OP'.", "status": "Error"})
		return
	}

	// Save the visit
	if err := initializers.DB.Create(&visit).Error; err != nil {
		logger.Error("Failed to create visit", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create visit", "status": "Error"})
		return
	}

	logger.Info("Visit created successfully", "visit_id", visit.ID)
	c.JSON(http.StatusOK, gin.H{
		"data":    visit.ID,
		"message": "Successfully created visit",
		"status":  "Success",
	})
}

func EditVisit(c *gin.Context) {
	var req schemas.EditVisitRequest
	logger := loggers.InitializeLogger()

	// Bind JSON input
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Invalid request payload in EditVisit", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload", "status": "Error"})
		return
	}

	// Validate visit_id
	if req.VisitID == 0 {
		logger.Error("Missing visit_id in EditVisit request")
		c.JSON(http.StatusBadRequest, gin.H{"message": "visit_id is required", "status": "Error"})
		return
	}

	// Fetch the visit from the database
	var visit models.Visit
	if err := initializers.DB.First(&visit, req.VisitID).Error; err != nil {
		logger.Error("Visit not found", "visit_id", req.VisitID, "error", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"message": "Visit not found", "status": "Error"})
		return
	}

	// Prepare update fields
	updates := map[string]interface{}{}
	if req.RoomNumber != nil {
		updates["room_number"] = *req.RoomNumber
	}
	if req.Diagnosis != nil {
		updates["diagnosis"] = *req.Diagnosis
	}
	if req.TreatmentPlan != nil {
		updates["treatment_plan"] = *req.TreatmentPlan
	}
	if req.Notes != nil {
		updates["notes"] = *req.Notes
	}
	if req.AdmissionDate != nil {
		updates["admission_date"] = *req.AdmissionDate
	}
	if req.DischargeDate != nil {
		updates["discharge_date"] = *req.DischargeDate
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	// Perform update
	if err := initializers.DB.Model(&visit).Updates(updates).Error; err != nil {
		logger.Error("Failed to update visit", "visit_id", req.VisitID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update visit", "status": "Error"})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{"message": "Visit updated successfully", "status": "Success"})
}
