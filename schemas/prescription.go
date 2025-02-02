package schemas

type PrescriptionInput struct {
	VisitID    uint   `json:"visit_id" binding:"required"`
	VisitType  string `json:"visit_type" binding:"required,oneof=IP OP"`
	PatientID  uint   `json:"patient_id" binding:"required"`
	DoctorID   uint   `json:"doctor_id" binding:"required"`
	DateIssued string `json:"date_issued" binding:"required"`
	Notes      string `json:"notes"`
	Items      []struct {
		MedicineID   uint   `json:"medicine_id" binding:"required"`
		Quantity     int    `json:"quantity" binding:"required,min=1"`
		Instructions string `json:"instructions"`
	} `json:"items" binding:"required,dive"`
}
