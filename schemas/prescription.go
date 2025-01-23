package schemas

type PrescriptionInput struct {
	PatientID uint                    `json:"patient_id" binding:"required"`
	DoctorID  uint                    `json:"doctor_id" binding:"required"`
	Notes     string                  `json:"notes"`
	Items     []PrescriptionItemInput `json:"items" binding:"required"`
}

type PrescriptionItemInput struct {
	MedicineID   uint   `json:"medicine_id" binding:"required"`
	Quantity     int    `json:"quantity" binding:"required"`
	Instructions string `json:"instructions"`
}
