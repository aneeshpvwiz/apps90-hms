package patientController

import (
	"apps90-hms/initializers"
	"apps90-hms/loggers"
	"apps90-hms/models"
	"apps90-hms/schemas"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetPatientDetails(c *gin.Context) {
	entityID := c.Query("entity_id")
	patientID := c.Query("patient_id")

	logger := loggers.InitializeLogger()

	// Validate entity_id and patient_id
	if entityID == "" || patientID == "" {
		logger.Error("Missing entity_id or patient_id in GetPatientDetails")
		c.JSON(http.StatusBadRequest, gin.H{"data": nil, "message": "Missing entity_id or patient_id", "status": "Error"})
		return
	}

	var patient models.Patient
	if err := initializers.DB.Where("id = ? AND entity_id = ?", patientID, entityID).First(&patient).Error; err != nil {
		logger.Error("Patient not found", "entity_id", entityID, "patient_id", patientID, "error", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"data": nil, "message": "Patient not found", "status": "Error"})
		return
	}

	logger.Info("Patient details fetched successfully", "patient_id", patientID, "entity_id", entityID)
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":             patient.ID,
			"first_name":     patient.FirstName,
			"last_name":      patient.LastName,
			"gender":         patient.Gender,
			"date_of_birth":  patient.DateOfBirth,
			"contact_number": patient.ContactNumber,
			"email":          patient.Email,
			"address":        patient.Address,
			"marital_status": patient.MaritalStatus,
			"occupation":     patient.Occupation,
		},
		"message": "Successfully fetched patient details",
		"status":  "Success",
	})
}

func GetPatientVisitHistory(c *gin.Context) {
	patientID := c.Query("patient_id")

	logger := loggers.InitializeLogger()

	// Validate query parameters
	if patientID == "" {
		logger.Error("Missing entity_id or patient_id in GetPatientVisitHistory")
		c.JSON(http.StatusBadRequest, gin.H{"data": nil, "message": "Missing entity_id or patient_id", "status": "Error"})
		return
	}

	var outpatientVisits []models.OutpatientVisit
	var inpatientVisits []models.InpatientVisit

	// Fetch outpatient visits
	initializers.DB.Where("patient_id = ?", patientID).Find(&outpatientVisits)

	// Fetch inpatient visits
	initializers.DB.Where("patient_id = ?", patientID).Find(&inpatientVisits)

	// Format the visits
	formattedOutpatientVisits := formatOutpatientVisits(outpatientVisits)
	formattedInpatientVisits := formatInpatientVisits(inpatientVisits)

	logger.Info("Patient visit history fetched successfully", "patient_id", patientID)
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"outpatient_visits": formattedOutpatientVisits,
			"inpatient_visits":  formattedInpatientVisits,
		},
		"message": "Successfully fetched visit history",
		"status":  "Success",
	})
}

func formatInpatientVisits(visits []models.InpatientVisit) []gin.H {
	var formattedVisits []gin.H
	for _, visit := range visits {
		formattedVisits = append(formattedVisits, gin.H{
			"id":             visit.ID,
			"admission_date": visit.AdmissionDate,
			"discharge_date": visit.DischargeDate,
			"room_number":    visit.RoomNumber,
			"diagnosis":      visit.Diagnosis,
			"treatment_plan": visit.TreatmentPlan,
			"notes":          visit.Notes,
			"doctor_id":      visit.DoctorID,
		})
	}
	return formattedVisits
}

func formatOutpatientVisits(visits []models.OutpatientVisit) []gin.H {
	var formattedVisits []gin.H
	for _, visit := range visits {
		formattedVisits = append(formattedVisits, gin.H{
			"id":             visit.ID,
			"visit_date":     visit.VisitDate,
			"diagnosis":      visit.Diagnosis,
			"treatment_plan": visit.TreatmentPlan,
			"notes":          visit.Notes,
			"doctor_id":      visit.DoctorID,
		})
	}
	return formattedVisits
}

// CreatePrescription handles creating a new prescription
func CreatePrescription(c *gin.Context) {
	var input schemas.PrescriptionInput
	logger := loggers.InitializeLogger()

	// Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("Error binding JSON for CreatePrescription", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"data": nil, "message": "Invalid request format", "status": "Error"})
		return
	}

	// Validate patient
	var patient models.Patient
	if err := initializers.DB.First(&patient, input.PatientID).Error; err != nil {
		logger.Error("Patient not found", "patient_id", input.PatientID, "error", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"data": nil, "message": "Patient not found", "status": "Error"})
		return
	}

	// Validate doctor
	var doctor models.Employee
	if err := initializers.DB.First(&doctor, input.DoctorID).Error; err != nil {
		logger.Error("Doctor not found", "doctor_id", input.DoctorID, "error", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"data": nil, "message": "Doctor not found", "status": "Error"})
		return
	}

	// Create prescription
	prescription := models.Prescription{
		PatientID:  input.PatientID,
		DoctorID:   input.DoctorID,
		DateIssued: time.Now(),
		Notes:      input.Notes,
	}

	initializers.DB.Create(&prescription)

	// Add prescription items
	for _, item := range input.Items {
		prescriptionItem := models.PrescriptionItem{
			PrescriptionID: prescription.ID,
			MedicineID:     item.MedicineID,
			Quantity:       item.Quantity,
			Instructions:   item.Instructions,
		}
		initializers.DB.Create(&prescriptionItem)
	}

	logger.Info("Prescription created successfully", "prescription_id", prescription.ID)
	c.JSON(http.StatusOK, gin.H{
		"data":    prescription.ID,
		"message": "Successfully created prescription",
		"status":  "Success",
	})
}

func GetPatientPrescriptions(c *gin.Context) {
	// Fetch patient_id from query params
	patientID := c.Query("patient_id")

	logger := loggers.InitializeLogger()

	// Validate patient_id
	if patientID == "" {
		logger.Error("Missing patient_id in GetPatientPrescriptions")
		c.JSON(http.StatusBadRequest, gin.H{"data": nil, "message": "Missing patient_id", "status": "Error"})
		return
	}

	var prescriptions []models.Prescription

	// Fetch prescriptions for the given patient
	if err := initializers.DB.Where("patient_id = ?", patientID).Preload("PrescriptionItems.Medicine").Find(&prescriptions).Error; err != nil {
		logger.Error("Error fetching prescriptions for patient", "patient_id", patientID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"data": nil, "message": "Failed to fetch prescriptions", "status": "Error"})
		return
	}

	// Format the response
	formattedPrescriptions := formatPrescriptions(prescriptions)

	logger.Info("Fetched prescriptions for patient", "patient_id", patientID)
	c.JSON(http.StatusOK, gin.H{
		"data":    formattedPrescriptions,
		"message": "Successfully fetched prescriptions",
		"status":  "Success",
	})
}

func formatPrescriptions(prescriptions []models.Prescription) []gin.H {
	var formatted []gin.H
	for _, prescription := range prescriptions {
		formatted = append(formatted, gin.H{
			"id":          prescription.ID,
			"date_issued": prescription.DateIssued.Format("2006-01-02"), // Format date
			"notes":       prescription.Notes,
			"items":       formatPrescriptionItems(prescription.PrescriptionItems),
		})
	}
	return formatted
}

func formatPrescriptionItems(items []models.PrescriptionItem) []gin.H {
	var formatted []gin.H
	for _, item := range items {
		formatted = append(formatted, gin.H{
			"medicine_id":   item.MedicineID,
			"medicine_name": item.Medicine.Name,
			"quantity":      item.Quantity,
			"instructions":  item.Instructions,
		})
	}
	return formatted
}
