package models

import "time"

// Prescription represents a prescription issued to a patient
type Prescription struct {
	ID                uint               `json:"id" gorm:"primaryKey"`
	VisitID           uint               `json:"visit_id"`   // Foreign key for the associated visit
	VisitType         string             `json:"visit_type" gorm:"type:varchar(20)"` // Can be "inpatient" or "outpatient"
	PatientID         uint               `json:"patient_id"`
	Patient           Patient            `json:"patient" gorm:"foreignKey:PatientID"`
	DoctorID          uint               `json:"doctor_id"`
	Doctor            Employee           `json:"doctor" gorm:"foreignKey:DoctorID"`
	DateIssued        time.Time          `json:"date_issued"`
	Notes             string             `json:"notes"`
	PrescriptionItems []PrescriptionItem `json:"items" gorm:"foreignKey:PrescriptionID"`
	IsActive          bool               `json:"is_active" gorm:"default:true"`
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
	MedicineID         uint       `json:"medicine_id"`
	Medicine           Medicine   `gorm:"foreignKey:MedicineID"`
	MedicineIntervals  string     `json:"medicine_intervals" gorm:"type:varchar(50)"`
	IsActive            bool         `json:"is_active" gorm:"default:true"`
	AuditFields       `gorm:"embedded"`
}

func (PrescriptionItem) TableName() string {
	return "prescription_item"
}


// PrescriptionTemplate represents a predefined prescription structure for diseases
type PrescriptionTemplate struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	EntityID    uint      `json:"entity_id"` // Foreign key reference to Entity
	Entity      Entity    `json:"entity" gorm:"foreignKey:EntityID"`
	DiseaseName string    `json:"disease_name" gorm:"type:varchar(200);unique;not null"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	Items       []PrescriptionTemplateItems `json:"items" gorm:"foreignKey:TemplateID;constraint:OnDelete:CASCADE"`
	AuditFields
}

func (PrescriptionTemplate) TableName() string {
	return "prescription_templates"
}

// PrescriptionTemplateItem represents individual medicines in a prescription template
type PrescriptionTemplateItems struct {
	ID         uint                `json:"id" gorm:"primaryKey"`
	TemplateID uint                `json:"template_id" gorm:"not null"` // Foreign key referencing PrescriptionTemplate
	Template   PrescriptionTemplate `json:"template" gorm:"foreignKey:TemplateID"`
	MedicineID uint                `json:"medicine_id"` // Foreign key referencing Medicine table
	Medicine   Medicine            `json:"medicine" gorm:"foreignKey:MedicineID"`
	PrescriptionDetails string       `json:"prescription_details" gorm:"type:text"` // Added prescription details
	MedicineIntervals  string     `json:"medicine_intervals" gorm:"type:varchar(50)"`
	IsActive   bool                `json:"is_active" gorm:"default:true"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
}

func (PrescriptionTemplateItems) TableName() string {
	return "prescription_template_items"
}
