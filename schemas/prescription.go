package schemas

import "time"

type CreatePrescriptionInput struct {
	VisitID             uint     `json:"visit_id" binding:"required"`
	VisitType           string   `json:"visit_type" binding:"required"` // "IP" or "OP"
	PatientID           uint     `json:"patient_id" binding:"required"`
	DoctorID            uint     `json:"doctor_id" binding:"required"`
	Notes               string   `json:"notes"`
	PrescriptionDetails []string `json:"prescription_details" binding:"required"`
}

type PrescriptionDetailsResponse struct {
	ID                uint      `json:"id"`
	DoctorName        string    `json:"doctor_name"`
	DateIssued        time.Time `json:"date_issued"`
	Notes             string    `json:"notes"`
	PrescriptionItems []string  `json:"prescription_items"`
}

type EditPrescriptionRequest struct {
	PrescriptionID    uint     `json:"prescription_id" binding:"required"`
	PrescriptionItems []string `json:"prescription_items" binding:"required"`
}

type EditPrescriptionItem struct {
	ID                  uint   `json:"id" binding:"required"` // Prescription item ID
	PrescriptionDetails string `json:"prescription_details"`  // New prescription details
}
