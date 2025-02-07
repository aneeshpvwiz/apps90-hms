package schemas

import "time"

type EntityInput struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
}

type UserEntityInput struct {
	UserID   uint `json:"user_id" binding:"required"`
	EntityID uint `json:"entity_id" binding:"required"`
}

type PatientInput struct {
	FirstName     string `json:"first_name" binding:"required"`
	LastName      string `json:"last_name" binding:"required"`
	Gender        string `json:"gender" binding:"required"`
	DateOfBirth   string `json:"date_of_birth" binding:"required"`
	ContactNumber string `json:"contact_number" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	Address       string `json:"address" binding:"required"`
	EntityID      uint   `json:"entity_id" binding:"required"`
	MaritalStatus string `json:"marital_status"`
	Occupation    string `json:"occupation"`
	DoctorID      uint   `json:"doctor_id" binding:"required"`
}

type EmployeeInput struct {
	FirstName          string `json:"first_name" binding:"required"`
	LastName           string `json:"last_name" binding:"required"`
	Email              string `json:"email" binding:"required,email"`
	PhoneNumber        string `json:"phone_number" binding:"required"`
	DateOfBirth        string `json:"date_of_birth" binding:"required"`
	EntityID           uint   `json:"entity_id" binding:"required"`
	EmployeeCategoryID uint   `json:"employee_category_id" binding:"required"`
}

// AppointmentInput represents the structure of appointment input in the request body.
type AppointmentInput struct {
	AppointmentTime time.Time `json:"appointment_time"`
	Reason          string    `json:"reason"`
	Notes           string    `json:"notes"`
	PatientID       uint      `json:"patient_id"`
	DoctorID        uint      `json:"doctor_id"`
	EntityID        uint      `json:"entity_id"`
}

type VisitInput struct {
	PatientID     uint       `json:"patient_id"`
	DoctorID      uint       `json:"doctor_id"`
	AppointmentID *uint      `json:"appointment_id,omitempty"` // Nullable for walk-ins
	VisitDate     time.Time  `json:"visit_date"`
	AdmissionDate time.Time  `json:"admission_date"`
	DischargeDate *time.Time `json:"discharge_date,omitempty"` // Nullable for ongoing admissions
	RoomNumber    string     `json:"room_number"`
	Diagnosis     string     `json:"diagnosis"`
	TreatmentPlan string     `json:"treatment_plan"`
	Notes         string     `json:"notes"`
	VisitType     string     `json:"visit_type"` // IP or OP
}

type PrescriptionResponse struct {
	ID uint `json:"id"` // Only returning prescription ID
}

type VisitResponse struct {
	ID            uint                   `json:"id"`
	VisitDate     time.Time              `json:"visit_date"`
	AdmissionDate *time.Time             `json:"admission_date,omitempty"`
	DischargeDate *time.Time             `json:"discharge_date,omitempty"`
	RoomNumber    string                 `json:"room_number"`
	Diagnosis     string                 `json:"diagnosis"`
	TreatmentPlan string                 `json:"treatment_plan"`
	Notes         string                 `json:"notes"`
	VisitType     string                 `json:"visit_type"`
	DoctorName    string                 `json:"doctor_name"` // Added doctor name
	Prescriptions []PrescriptionResponse `json:"prescriptions,omitempty"`
}
