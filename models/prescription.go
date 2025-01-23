package models

import "time"

// Prescription represents a prescription issued to a patient
type Prescription struct {
	ID                uint               `json:"id" gorm:"primaryKey"`
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
	ID             uint         `json:"id" gorm:"primaryKey"`
	PrescriptionID uint         `json:"prescription_id"`
	Prescription   Prescription `json:"prescription" gorm:"foreignKey:PrescriptionID"`
	MedicineID     uint         `json:"medicine_id"`
	Medicine       Medicine     `json:"medicine" gorm:"foreignKey:MedicineID"`
	Quantity       int          `json:"quantity"`
	Instructions   string       `json:"instructions"`
}

func (PrescriptionItem) TableName() string {
	return "prescription_item"
}
