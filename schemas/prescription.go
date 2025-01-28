package schemas

// PrescriptionInput represents the input structure for creating a prescription
type PrescriptionInput struct {
	PatientID uint                    `json:"patient_id"` // ID of the patient
	DoctorID  uint                    `json:"doctor_id"`  // ID of the doctor
	VisitID   uint                    `json:"visit_id"`   // ID of the visit (inpatient or outpatient)
	VisitType string                  `json:"visit_type"` // Type of visit: "inpatient" or "outpatient"
	Notes     string                  `json:"notes"`      // Additional notes for the prescription
	Items     []PrescriptionItemInput `json:"items"`      // List of prescription items
}

// PrescriptionItemInput represents the structure for a single prescription item
type PrescriptionItemInput struct {
	MedicineID   uint   `json:"medicine_id"`  // ID of the medicine
	Quantity     int    `json:"quantity"`     // Quantity of the medicine
	Instructions string `json:"instructions"` // Instructions for the medicine usage
}
