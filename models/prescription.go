package models

import "time"

// Prescription represents a prescription issued to a patient
type Prescription struct {
	ID                uint               `json:"id" gorm:"primaryKey"`
	VisitID           uint               `json:"visit_id"`   // Foreign key for the associated visit
	VisitType         string             `json:"visit_type"` // Can be "inpatient" or "outpatient"
	PatientID         uint               `json:"patient_id"`
	Patient           Patient            `json:"patient" gorm:"foreignKey:PatientID"`
	DoctorID          uint               `json:"doctor_id"`
	Doctor            Employee           `json:"doctor" gorm:"foreignKey:DoctorID"`
	DateIssued        time.Time          `json:"date_issued"`
	Notes             string             `json:"notes"`
	PrescriptionItems []PrescriptionItem `json:"items" gorm:"foreignKey:PrescriptionID"`
	AuditFields       `gorm:"embedded"`
}

func (Prescription) TableName() string {
	return "prescription"
}

// PrescriptionItem represents the medicines in a prescription
type PrescriptionItem struct {
	ID                  uint         `json:"id" gorm:"primaryKey"`
	PrescriptionID      uint         `json:"prescription_id"`
	Prescription        Prescription `json:"prescription" gorm:"foreignKey:PrescriptionID"`
	PrescriptionDetails string       `json:"prescription_details" gorm:"type:text"` // Added prescription details
}

func (PrescriptionItem) TableName() string {
	return "prescription_item"
}
