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

	// Prepare visit responses
	var inpatientVisits []schemas.VisitResponse
	var outpatientVisits []schemas.VisitResponse

	// Iterate through visits
	for _, visit := range visits {
		var prescriptions []models.Prescription

		// Fetch only prescription IDs for the visit
		if err := initializers.DB.Select("id").Where("visit_id = ?", visit.ID).
			Find(&prescriptions).Error; err != nil {
			logger.Error("Error fetching prescriptions for visit", "visit_id", visit.ID, "error", err)
			continue
		}

		// Format prescriptions
		var prescriptionList []schemas.PrescriptionResponse
		for _, prescription := range prescriptions {
			prescriptionList = append(prescriptionList, schemas.PrescriptionResponse{
				ID: prescription.ID,
			})
		}

		// Ensure prescriptions is an empty array if none exist
		if len(prescriptionList) == 0 {
			prescriptionList = []schemas.PrescriptionResponse{}
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
			AdmissionDate: &visit.AdmissionDate,
			DischargeDate: visit.DischargeDate, // Added discharge date
			DoctorName:    visit.Doctor.FirstName + " " + visit.Doctor.LastName,
			Prescriptions: prescriptionList, // Ensures it always returns an array
		}

		// Sort into inpatient or outpatient visits
		if visit.VisitType == "IP" {
			inpatientVisits = append(inpatientVisits, visitResponse)
		} else if visit.VisitType == "OP" {
			outpatientVisits = append(outpatientVisits, visitResponse)
		}
	}

	// Ensure empty lists instead of null if no visits are found
	if len(inpatientVisits) == 0 {
		inpatientVisits = []schemas.VisitResponse{}
	}
	if len(outpatientVisits) == 0 {
		outpatientVisits = []schemas.VisitResponse{}
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
	var input schemas.CreatePrescriptionInput
	logger := loggers.InitializeLogger()

	// Bind request body
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("Error binding JSON for Create Prescription", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request format", "status": "Error"})
		return
	}

	// Validate patient
	var patient models.Patient
	if err := initializers.DB.First(&patient, input.PatientID).Error; err != nil {
		logger.Warn("Patient not found", "patient_id", input.PatientID)
		c.JSON(http.StatusNotFound, gin.H{"message": "Patient not found", "status": "Error"})
		return
	}

	// Validate doctor
	var doctor models.Employee
	if err := initializers.DB.First(&doctor, input.DoctorID).Error; err != nil {
		logger.Warn("Doctor not found", "doctor_id", input.DoctorID)
		c.JSON(http.StatusNotFound, gin.H{"message": "Doctor not found", "status": "Error"})
		return
	}

	// Validate visit
	var visit models.Visit
	if err := initializers.DB.First(&visit, input.VisitID).Error; err != nil {
		logger.Warn("Visit not found", "visit_id", input.VisitID)
		c.JSON(http.StatusNotFound, gin.H{"message": "Visit not found", "status": "Error"})
		return
	}

	// Create prescription
	prescription := models.Prescription{
		VisitID:    input.VisitID,
		VisitType:  input.VisitType,
		PatientID:  input.PatientID,
		DoctorID:   input.DoctorID,
		DateIssued: time.Now(),
		Notes:      input.Notes,
	}

	if err := initializers.DB.Create(&prescription).Error; err != nil {
		logger.Error("Failed to create prescription", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create prescription", "status": "Error"})
		return
	}

	// Create prescription items
	var prescriptionItems []models.PrescriptionItem
	for _, item := range input.PrescriptionDetails {
		prescriptionItems = append(prescriptionItems, models.PrescriptionItem{
			PrescriptionID:      prescription.ID,
			PrescriptionDetails: item,
		})
	}

	if len(prescriptionItems) > 0 {
		if err := initializers.DB.Create(&prescriptionItems).Error; err != nil {
			logger.Error("Failed to create prescription items", "error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create prescription items", "status": "Error"})
			return
		}
	}

	logger.Info("Prescription created successfully", "prescription_id", prescription.ID)
	c.JSON(http.StatusOK, gin.H{
		"data":    prescription.ID,
		"message": "Successfully created prescription",
		"status":  "Success",
	})
}

func GetPrescriptionDetails(c *gin.Context) {
	prescriptionID := c.Query("prescription_id")
	logger := loggers.InitializeLogger()

	// Validate prescription_id
	if prescriptionID == "" {
		logger.Error("Missing prescription_id in GetPrescriptionDetails")
		c.JSON(http.StatusBadRequest, gin.H{"data": nil, "message": "Missing prescription_id", "status": "Error"})
		return
	}

	// Fetch prescription
	var prescription models.Prescription
	if err := initializers.DB.Preload("Doctor").Preload("PrescriptionItems").First(&prescription, prescriptionID).Error; err != nil {
		logger.Error("Prescription not found", "prescription_id", prescriptionID, "error", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"data": nil, "message": "Prescription not found", "status": "Error"})
		return
	}

	// Extract prescription details
	var prescriptionDetails []string
	for _, item := range prescription.PrescriptionItems {
		prescriptionDetails = append(prescriptionDetails, item.PrescriptionDetails)
	}

	// Response format
	response := schemas.PrescriptionDetailsResponse{
		ID:                prescription.ID,
		DoctorName:        prescription.Doctor.FirstName + " " + prescription.Doctor.LastName,
		DateIssued:        prescription.DateIssued,
		Notes:             prescription.Notes,
		PrescriptionItems: prescriptionDetails,
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": "Successfully fetched prescription details",
		"status":  "Success",
	})
}
