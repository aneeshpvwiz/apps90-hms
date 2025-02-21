package models

import "time"

type Entity struct {
	ID          uint              `json:"id" gorm:"primary_key"`
	Name        string            `json:"name" gorm:"type:varchar(254);unique"`
	Address     string            `json:"address" gorm:"type:text"`
	Users       []User            `gorm:"many2many:user_entity;"`
	Employees   []Employee        `json:"employees" gorm:"foreignKey:EntityID"`
	Patients    []Patient         `json:"patients" gorm:"foreignKey:EntityID"` // One-to-many relationship with Patient
	AuditFields `gorm:"embedded"` // Embedding AuditFields
}

// TableName specifies the table name for the Entity model.
func (Entity) TableName() string {
	return "entity"
}

type UserEntity struct {
	ID          uint              `json:"id" gorm:"primaryKey"`
	UserID      uint              `json:"user_id" gorm:"index"`   // Foreign key for User
	EntityID    uint              `json:"entity_id" gorm:"index"` // Foreign key for Entity
	AuditFields `gorm:"embedded"` // Embedding AuditFields
}

// TableName specifies the table name for the Entity model.
func (UserEntity) TableName() string {
	return "user_entity"
}

type EmployeeCategory struct {
	ID          uint              `json:"id" gorm:"primary_key"`
	Name        string            `json:"name" gorm:"type:varchar(20);unique"`
	Description string            `json:"description"`
	Employees   []Employee        `json:"employees" gorm:"foreignKey:EmployeeCategoryID"`
	AuditFields `gorm:"embedded"` // Embedding AuditFields
}

// TableName specifies the table name for the Entity model.
func (EmployeeCategory) TableName() string {
	return "employee_category"
}

type Employee struct {
	ID                 uint              `json:"id" gorm:"primary_key"`
	FirstName          string            `json:"first_name"`
	LastName           string            `json:"last_name"`
	Email              string            `json:"email" gorm:"unique"`
	PhoneNumber        string            `json:"phone_number"`
	DateOfBirth        string            `json:"date_of_birth"`
	EntityID           uint              `json:"entity_id"` // Foreign key to Entity
	Entity             Entity            `json:"entity" gorm:"foreignKey:EntityID"`
	EmployeeCategoryID uint              `json:"employee_category_id"` // Foreign key to EmployeeCategory
	EmployeeCategory   EmployeeCategory  `json:"employee_category" gorm:"foreignKey:EmployeeCategoryID"`
	Patients           []Patient         `json:"patients" gorm:"foreignKey:DoctorID"` // List of patients assigned to this doctor
	IsActive           bool               `json:"is_active" gorm:"default:true"` 
	AuditFields        `gorm:"embedded"` // Embedding AuditFields
}

// TableName specifies the table name for the Entity model.
func (Employee) TableName() string {
	return "employee"
}

// Patient represents a patient in the hospital
type Patient struct {
	ID            uint              `json:"id" gorm:"primaryKey"`
	FirstName     string            `json:"first_name"`
	LastName      string            `json:"last_name"`
	Gender        string            `json:"gender"`
	DateOfBirth   time.Time         `json:"date_of_birth"`
	ContactNumber string            `json:"contact_number"`
	Email         string            `json:"email" gorm:"unique"`
	Address       string            `json:"address" gorm:"type:text"`
	EntityID      uint              `json:"entity_id"`
	Entity        Entity            `json:"entity" gorm:"foreignKey:EntityID"`
	MaritalStatus string            `json:"marital_status"`
	Occupation    string            `json:"occupation" gorm:"default:null"`
	DoctorID      uint              `json:"doctor_id"`                         // Foreign key to Employee (Doctor)
	Doctor        Employee          `json:"doctor" gorm:"foreignKey:DoctorID"` // Reference to the doctor
	IsActive      bool              `json:"is_active" gorm:"default:true"`
	AuditFields   `gorm:"embedded"` // Embedding AuditFields
}

func (Patient) TableName() string {
	return "patient"
}
