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

type Visit struct {
	ID            uint              `json:"id" gorm:"primaryKey"`
	AppointmentID *uint             `json:"appointment_id"` // Nullable for walk-ins
	Appointment   *Appointment      `json:"appointment" gorm:"foreignKey:AppointmentID"`
	VisitDate     time.Time         `json:"visit_date"`
	AdmissionDate time.Time         `json:"admission_date"`
	DischargeDate *time.Time        `json:"discharge_date"` // Nullable for ongoing admissions
	RoomNumber    string            `json:"room_number"`
	Diagnosis     string            `json:"diagnosis"`
	TreatmentPlan string            `json:"treatment_plan"`
	Notes         string            `json:"notes"`
	PatientID     uint              `json:"patient_id"`
	Patient       Patient           `json:"patient" gorm:"foreignKey:PatientID"`
	DoctorID      uint              `json:"doctor_id"`
	Doctor        Employee          `json:"doctor" gorm:"foreignKey:DoctorID"`
	Prescriptions []Prescription    `json:"prescriptions" gorm:"foreignKey:VisitID;references:ID"`
	VisitType     string            `json:"visit_type" `
	AuditFields   `gorm:"embedded"` // Embedding AuditFields
}

func (Visit) TableName() string {
	return "visit"
}
