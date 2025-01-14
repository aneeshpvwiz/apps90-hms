package schemas

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
