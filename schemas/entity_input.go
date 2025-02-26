package schemas

import (
	"encoding/json"
	"time"
)

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
	IsActive      *bool  `json:"is_active,omitempty"`
}

type EmployeeInput struct {
	FirstName          string `json:"first_name" binding:"required"`
	LastName           string `json:"last_name" binding:"required"`
	Email              string `json:"email" binding:"required,email"`
	PhoneNumber        string `json:"phone_number" binding:"required"`
	DateOfBirth        time.Time `json:"date_of_birth" binding:"required"`
	EntityID           uint   `json:"entity_id" binding:"required"`
	EmployeeCategoryID uint   `json:"employee_category_id" binding:"required"`
}

// AppointmentInput represents the structure of appointment input in the request body.
type AppointmentInput struct {
	StartTime time.Time `json:"start_time"`
	EndTime time.Time `json:"end_time"`
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
	ID         uint      `json:"id"`
	DateIssued time.Time `json:"date_issued"`
}

type VisitResponse struct {
	ID            uint                   `json:"id"`
	VisitDate     time.Time              `json:"visit_date"`
	AdmissionDate *time.Time             `json:"admission_date,omitempty"`
	DischargeDate *time.Time             `json:"discharge_date,omitempty"`
	Diagnosis     string                 `json:"diagnosis"`
	TreatmentPlan string                 `json:"treatment_plan"`
	Notes         string                 `json:"notes"`
	VisitType     string                 `json:"visit_type"`
	DoctorName    string                 `json:"doctor_name"` // Added doctor name
	Prescriptions []PrescriptionResponse `json:"prescriptions"`
}

type MedicineCategoryResponse struct {
	ID        uint               `json:"id"`
	Name      string             `json:"name"`
	Medicines []MedicineResponse `json:"medicines"`
}

type MedicineResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type AddMedicineCategoryRequest struct {
	NameTranslations json.RawMessage `json:"name_translations" binding:"required"`
	EntityID    uint   `json:"entity_id" binding:"required"`
}

type AddMedicineRequest struct {
	EntityID         uint            `json:"entity_id" binding:"required"`
	CategoryID       uint            `json:"category_id" binding:"required"`
	NameTranslations json.RawMessage `json:"name_translations" binding:"required"`
	Description      string          `json:"description"`
	Price           float64          `json:"price"`
	Quantity        string            `json:"quantity"`
}

type EditVisitRequest struct {
	VisitID       uint       `json:"visit_id" binding:"required"`
	RoomNumber    *string    `json:"room_number,omitempty"`
	Diagnosis     *string    `json:"diagnosis,omitempty"`
	TreatmentPlan *string    `json:"treatment_plan,omitempty"`
	Notes         *string    `json:"notes,omitempty"`
	AdmissionDate *time.Time `json:"admission_date,omitempty"`
	DischargeDate *time.Time `json:"discharge_date,omitempty"`
	IsActive      *bool      `json:"is_active,omitempty"`
}

