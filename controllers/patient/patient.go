package patientController

import (
	"apps90-hms/errors"
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

	// Validate patient_id
	if patientID == "" {
		logger.Error("Missing patient_id in GetPatientVisitHistory")
		c.JSON(http.StatusBadRequest, gin.H{"data": nil, "message": "Missing patient_id", "status": "Error"})
		return
	}

	// Fetch inpatient visits
	var inpatientVisits []models.InpatientVisit
	if err := initializers.DB.Where("patient_id = ?", patientID).Find(&inpatientVisits).Error; err != nil {
		logger.Error("Error fetching inpatient visits for patient", "patient_id", patientID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"data": nil, "message": "Failed to fetch visit history", "status": "Error"})
		return
	}

	// Fetch outpatient visits
	var outpatientVisits []models.OutpatientVisit
	if err := initializers.DB.Where("patient_id = ?", patientID).Find(&outpatientVisits).Error; err != nil {
		logger.Error("Error fetching outpatient visits for patient", "patient_id", patientID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"data": nil, "message": "Failed to fetch visit history", "status": "Error"})
		return
	}

	// Now, fetch prescriptions for each outpatient visit explicitly
	for i, visit := range outpatientVisits {
		var prescriptions []models.Prescription
		// Fetch prescriptions using the visit ID
		if err := initializers.DB.Where("patient_id = ? AND visit_id = ? AND visit_type = ?", patientID, visit.ID, "OP").Find(&prescriptions).Error; err != nil {
			logger.Error("Error fetching prescriptions for visit", "visit_id", visit.ID, "error", err)
			continue
		}

		// Assign the fetched prescriptions to the visit
		outpatientVisits[i].Prescriptions = prescriptions
	}

	// Now, fetch prescriptions for each visit (inpatient or outpatient) explicitly
	for i, visit := range inpatientVisits {
		var prescriptions []models.Prescription

		// Fetch prescriptions based on visit_id and visit_type from the Prescription model
		if err := initializers.DB.Where("patient_id = ? AND visit_id = ? AND visit_type = ?", patientID, visit.ID, "IP").Find(&prescriptions).Error; err != nil {
			logger.Error("Error fetching prescriptions for visit", "visit_id", visit.ID, "error", err)
			continue
		}

		// Assign the fetched prescriptions to the visit
		inpatientVisits[i].Prescriptions = prescriptions
	}

	// Format responses
	formattedInpatientVisits := formatInpatientVisits(inpatientVisits)
	formattedOutpatientVisits := formatOutpatientVisits(outpatientVisits)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"inpatient_visits":  formattedInpatientVisits,
			"outpatient_visits": formattedOutpatientVisits,
		},
		"message": "Successfully fetched patient visit history",
		"status":  "Success",
	})
}

func formatInpatientVisits(visits []models.InpatientVisit) []gin.H {
	var formatted []gin.H

	for _, visit := range visits {
		prescriptionData := []gin.H{}
		if visit.Prescriptions != nil {
			for _, prescription := range visit.Prescriptions {
				prescriptionData = append(prescriptionData, gin.H{
					"id":          prescription.ID,
					"doctor_name": prescription.Doctor.FirstName,
					"date_issued": prescription.DateIssued.Format("2006-01-02"),
					"notes":       prescription.Notes,
				})
			}
		}

		// Check if DischargeDate is nil (pointer is nil)
		dischargeDate := ""
		if visit.DischargeDate != nil {
			dischargeDate = visit.DischargeDate.Format("2006-01-02")
		}

		formatted = append(formatted, gin.H{
			"id":             visit.ID,
			"admission_date": visit.AdmissionDate.Format("2006-01-02"),
			"discharge_date": dischargeDate,
			"notes":          visit.Notes,
			"prescriptions":  prescriptionData,
		})
	}

	return formatted
}

func formatOutpatientVisits(visits []models.OutpatientVisit) []gin.H {
	var formatted []gin.H

	for _, visit := range visits {
		prescriptionData := []gin.H{}
		if visit.Prescriptions != nil {
			for _, prescription := range visit.Prescriptions {
				prescriptionData = append(prescriptionData, gin.H{
					"id":          prescription.ID,
					"doctor_name": prescription.Doctor.FirstName,
					"date_issued": prescription.DateIssued.Format("2006-01-02"),
					"notes":       prescription.Notes,
				})
			}
		}

		formatted = append(formatted, gin.H{
			"id":            visit.ID,
			"visit_date":    visit.VisitDate.Format("2006-01-02"),
			"notes":         visit.Notes,
			"prescriptions": prescriptionData,
		})
	}

	return formatted
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

func GetPrescriptionDetails(c *gin.Context) {
	logger := loggers.InitializeLogger()

	// Fetch the prescription_id from query parameters
	prescriptionID := c.Query("prescription_id")
	if prescriptionID == "" {
		logger.Error("Missing prescription_id in query parameters")
		c.Error(models.WrapError(http.StatusBadRequest, errors.ErrBindingJSON, "Missing prescription_id in query"))
		return
	}

	var prescription models.Prescription

	// Fetch the prescription details, including prescription items
	if err := initializers.DB.Preload("PrescriptionItems.Medicine").
		Where("id = ?", prescriptionID).
		First(&prescription).Error; err != nil {
		logger.Error("Error fetching prescription details", "prescription_id", prescriptionID, "error", err)
		c.Error(models.WrapError(http.StatusInternalServerError, errors.ErrDatabaseFailed, "Error fetching prescription details"))
		return
	}

	// Format response structure
	response := gin.H{
		"prescription_id": prescription.ID,
		"date_issued":     prescription.DateIssued.Format("2006-01-02"),
		"notes":           prescription.Notes,
		"items": func(items []models.PrescriptionItem) []gin.H {
			formattedItems := make([]gin.H, 0)
			for _, item := range items {
				formattedItems = append(formattedItems, gin.H{
					"medicine_id":   item.MedicineID,
					"medicine_name": item.Medicine.Name,
					"quantity":      item.Quantity,
					"instructions":  item.Instructions,
				})
			}
			return formattedItems
		}(prescription.PrescriptionItems),
	}

	logger.Info("Prescription details fetched successfully", "prescription_id", prescription.ID)

	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": "Successfully fetched prescription details",
		"status":  "Success",
	})
}
