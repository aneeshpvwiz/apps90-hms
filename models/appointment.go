package models

import "time"

// Appointment represents an appointment between a patient and a doctor
type Appointment struct {
	ID              uint              `json:"id" gorm:"primaryKey"`
	AppointmentTime time.Time         `json:"appointment_time"`
	Reason          string            `json:"reason"`
	Notes           string            `json:"notes"`
	PatientID       uint              `json:"patient_id"`
	Patient         Patient           `json:"patient" gorm:"foreignKey:PatientID"`
	EmployeeID      uint              `json:"employee_id"`
	Employee        Employee          `json:"employee" gorm:"foreignKey:EmployeeID"`
	EntityID        uint              `json:"entity_id"`
	Entity          Entity            `json:"enity" gorm:"foreignKey:EntityID"`
	AuditFields     `gorm:"embedded"` // Embedding AuditFields
}

func (Appointment) TableName() string {
	return "appointment"
}
