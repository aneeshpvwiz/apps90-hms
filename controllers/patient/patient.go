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

	// Prepare response lists
	inpatientVisits := []schemas.VisitResponse{}
	outpatientVisits := []schemas.VisitResponse{}

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
	var input schemas.PrescriptionInput
	logger := loggers.InitializeLogger()

	// Bind request body
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("Invalid request format for Create Prescription", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request format", "status": "Error"})
		return
	}

	// Validate Patient
	var patient models.Patient
	if err := initializers.DB.First(&patient, input.PatientID).Error; err != nil {
		logger.Warn("Patient not found", "patient_id", input.PatientID)
		c.JSON(http.StatusNotFound, gin.H{"message": "Patient not found", "status": "Error"})
		return
	}

	// Validate Doctor
	var doctor models.Employee
	if err := initializers.DB.First(&doctor, input.DoctorID).Error; err != nil {
		logger.Warn("Doctor not found", "doctor_id", input.DoctorID)
		c.JSON(http.StatusNotFound, gin.H{"message": "Doctor not found", "status": "Error"})
		return
	}

	// Validate Visit
	var visit models.Visit
	if err := initializers.DB.First(&visit, input.VisitID).Error; err != nil {
		logger.Warn("Visit not found", "visit_id", input.VisitID)
		c.JSON(http.StatusNotFound, gin.H{"message": "Visit not found", "status": "Error"})
		return
	}

	// Convert DateIssued string to time.Time
	dateIssued, err := time.Parse("2006-01-02", input.DateIssued)
	if err != nil {
		logger.Error("Invalid date format", "date_issued", input.DateIssued)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid date format (YYYY-MM-DD)", "status": "Error"})
		return
	}

	// Create Prescription
	prescription := models.Prescription{
		VisitID:    input.VisitID,
		VisitType:  input.VisitType,
		PatientID:  input.PatientID,
		DoctorID:   input.DoctorID,
		DateIssued: dateIssued,
		Notes:      input.Notes,
	}

	// Save Prescription
	if err := initializers.DB.Create(&prescription).Error; err != nil {
		logger.Error("Failed to create prescription", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create prescription", "status": "Error"})
		return
	}

	// Insert Prescription Items
	var prescriptionItems []models.PrescriptionItem
	for _, item := range input.Items {
		prescriptionItems = append(prescriptionItems, models.PrescriptionItem{
			PrescriptionID: prescription.ID,
			MedicineID:     item.MedicineID,
			Quantity:       item.Quantity,
			Instructions:   item.Instructions,
		})
	}

	// Save Prescription Items
	if len(prescriptionItems) > 0 {
		if err := initializers.DB.Create(&prescriptionItems).Error; err != nil {
			logger.Error("Failed to save prescription items", "error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save prescription items", "status": "Error"})
			return
		}
	}

	logger.Info("Prescription created successfully", "prescription_id", prescription.ID)
	c.JSON(http.StatusOK, gin.H{
		"data":    gin.H{"prescription_id": prescription.ID},
		"message": "Prescription created successfully",
		"status":  "Success",
	})
}

func GetPrescriptionDetails(c *gin.Context) {
	prescriptionID := c.Query("prescription_id")
	logger := loggers.InitializeLogger()

	// Validate prescription_id
	if prescriptionID == "" {
		logger.Error("Missing prescription_id in request")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing prescription ID", "status": "Error"})
		return
	}

	// Fetch Prescription with related data
	var prescription models.Prescription
	if err := initializers.DB.
		Preload("Doctor").
		Preload("PrescriptionItems.Medicine").
		First(&prescription, prescriptionID).Error; err != nil {
		logger.Error("Prescription not found", "prescription_id", prescriptionID)
		c.JSON(http.StatusNotFound, gin.H{"message": "Prescription not found", "status": "Error"})
		return
	}

	// Prepare Response
	response := gin.H{
		"prescription_id": prescription.ID,
		"doctor_name":     prescription.Doctor.FirstName + " " + prescription.Doctor.LastName,
		"date_issued":     prescription.DateIssued.Format("2006-01-02"),
		"notes":           prescription.Notes,
		"items":           []gin.H{},
	}

	// Add Prescription Items with Medicine Names
	for _, item := range prescription.PrescriptionItems {
		response["items"] = append(response["items"].([]gin.H), gin.H{
			"medicine_name": item.Medicine.Name, // Fetching medicine name
			"quantity":      item.Quantity,
			"instructions":  item.Instructions,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": "Successfully fetched prescription details",
		"status":  "Success",
	})
}
