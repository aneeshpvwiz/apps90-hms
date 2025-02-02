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

	// Fetch visits (both inpatient and outpatient)
	var visits []models.Visit
	if err := initializers.DB.Where("patient_id = ?", patientID).Preload("Doctor").Find(&visits).Error; err != nil {
		logger.Error("Error fetching visits for patient", "patient_id", patientID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"data": nil, "message": "Failed to fetch visit history", "status": "Error"})
		return
	}

	// Prepare visit responses
	var inpatientVisits []schemas.VisitResponse
	var outpatientVisits []schemas.VisitResponse

	// Iterate through visits
	for _, visit := range visits {
		var prescriptions []models.Prescription

		// Fetch only prescription IDs for the visit
		if err := initializers.DB.Select("id").Where("visit_id = ? AND visit_type = ?", visit.ID, visit.VisitType).
			Find(&prescriptions).Error; err != nil {
			logger.Error("Error fetching prescriptions for visit", "visit_id", visit.ID, "error", err)
			continue
		}

		// Extract prescription IDs
		var prescriptionData []schemas.PrescriptionResponse
		for _, prescription := range prescriptions {
			prescriptionData = append(prescriptionData, schemas.PrescriptionResponse{
				ID: prescription.ID,
			})
		}

		// Create a VisitResponse object and populate necessary fields
		visitResponse := schemas.VisitResponse{
			ID:            visit.ID,
			VisitDate:     visit.VisitDate,
			RoomNumber:    visit.RoomNumber,
			Diagnosis:     visit.Diagnosis,
			TreatmentPlan: visit.TreatmentPlan,
			Notes:         visit.Notes,
			VisitType:     visit.VisitType,
			DoctorName:    visit.Doctor.FirstName + " " + visit.Doctor.LastName, // Fetch doctor's name
			Prescriptions: prescriptionData,
		}

		// Sort into inpatient or outpatient visits
		if visit.VisitType == "IP" {
			inpatientVisits = append(inpatientVisits, visitResponse)
		} else if visit.VisitType == "OP" {
			outpatientVisits = append(outpatientVisits, visitResponse)
		}
	}

	// Send the formatted response
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"inpatient_visits":  inpatientVisits,
			"outpatient_visits": outpatientVisits,
		},
		"message": "Successfully fetched patient visit history",
		"status":  "Success",
	})
}

func CreatePrescription(c *gin.Context) {
	var prescriptionInput schemas.PrescriptionInput
	logger := loggers.InitializeLogger()

	if err := c.ShouldBindJSON(&prescriptionInput); err != nil {
		logger.Error("Error binding JSON for CreatePrescription", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Failure",
			"message": "Invalid input format",
		})
		return
	}

	if prescriptionInput.VisitID == 0 || prescriptionInput.VisitType == "" {
		logger.Error("Visit ID and Visit Type are required")
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Failure",
			"message": "Visit ID and Visit Type are required",
		})
		return
	}

	// Validate doctor
	var doctor models.Employee
	if err := initializers.DB.Where("id = ? AND is_active = ?", prescriptionInput.DoctorID, true).First(&doctor).Error; err != nil {
		logger.Error("Invalid doctor_id", "doctor_id", prescriptionInput.DoctorID, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Failure",
			"message": "Invalid doctor ID",
		})
		return
	}

	// Validate patient
	var patient models.Patient
	if err := initializers.DB.Where("id = ? AND is_active = ?", prescriptionInput.PatientID, true).First(&patient).Error; err != nil {
		logger.Error("Invalid patient_id", "patient_id", prescriptionInput.PatientID, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Failure",
			"message": "Invalid patient ID",
		})
		return
	}

	// Validate visit exists based on visit_type
	var visitExists bool
	if prescriptionInput.VisitType == "OP" {
		var outpatientVisit models.Visit
		if err := initializers.DB.Where("id = ?", prescriptionInput.VisitID).First(&outpatientVisit).Error; err != nil {
			logger.Error("Outpatient visit not found", "visit_id", prescriptionInput.VisitID, "error", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "Failure",
				"message": "Visit ID does not exist for the given visit type",
			})
			return
		}
		visitExists = true
	} else if prescriptionInput.VisitType == "IP" {
		var inpatientVisit models.Visit
		if err := initializers.DB.Where("id = ?", prescriptionInput.VisitID).First(&inpatientVisit).Error; err != nil {
			logger.Error("Inpatient visit not found", "visit_id", prescriptionInput.VisitID, "error", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "Failure",
				"message": "Visit ID does not exist for the given visit type",
			})
			return
		}
		visitExists = true
	}

	if !visitExists {
		logger.Error("Visit ID does not exist for the given visit type", "visit_id", prescriptionInput.VisitID, "visit_type", prescriptionInput.VisitType)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Failure",
			"message": "Visit ID does not exist for the given visit type",
		})
		return
	}

	// Create prescription
	prescription := models.Prescription{
		PatientID:  prescriptionInput.PatientID,
		DoctorID:   prescriptionInput.DoctorID,
		VisitID:    prescriptionInput.VisitID,
		VisitType:  prescriptionInput.VisitType,
		DateIssued: time.Now(),
		Notes:      prescriptionInput.Notes,
	}

	if err := initializers.DB.Create(&prescription).Error; err != nil {
		logger.Error("Error creating prescription", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Failure",
			"message": "Error creating prescription",
		})
		return
	}

	// Add prescription items
	for _, item := range prescriptionInput.Items {
		prescriptionItem := models.PrescriptionItem{
			PrescriptionID: prescription.ID,
			MedicineID:     item.MedicineID,
			Quantity:       item.Quantity,
			Instructions:   item.Instructions,
		}

		if err := initializers.DB.Create(&prescriptionItem).Error; err != nil {
			logger.Error("Error creating prescription item", "prescription_id", prescription.ID, "error", err)
			continue
		}
	}

	logger.Info("Prescription created successfully", "prescription_id", prescription.ID)
	c.JSON(http.StatusOK, gin.H{
		"data":    prescription.ID,
		"message": "Successfully created a prescription",
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
