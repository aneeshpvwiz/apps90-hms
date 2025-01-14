package controllers

import (
	"apps90-hms/errors"
	"apps90-hms/initializers"
	"apps90-hms/loggers"
	"apps90-hms/models"
	"apps90-hms/schemas"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateEntity(c *gin.Context) {

	var entityInput schemas.EntityInput

	logger := loggers.InitializeLogger()

	if err := c.ShouldBindJSON(&entityInput); err != nil {
		logger.Error("Error binding JSON for Create Entity", "error", err.Error())
		c.Error(models.WrapError(http.StatusBadRequest, errors.ErrBindingJSON, "Invalid request format"))
		return
	}

	var entityFound models.Entity
	initializers.DB.Where("name=?", entityInput.Name).Find(&entityFound)

	if entityFound.ID != 0 {
		logger.Warn("Entity with this name already exists", "name", entityInput.Name)
		c.Error(models.WrapError(http.StatusBadRequest, errors.ErrObjectExists, "Entity with this name already exist"))
		return
	}

	entity := models.Entity{
		Name:    entityInput.Name,
		Address: entityInput.Address,
	}

	initializers.DB.Create(&entity)

	logger.Info("Entity created successfully", "Name", entityInput.Name, "entity_id", entity.ID)

	c.JSON(http.StatusOK, gin.H{"data": entity})

}

func CreateUserEntity(c *gin.Context) {
	var userEntityInput schemas.UserEntityInput

	if err := c.ShouldBindJSON(&userEntityInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userEntity := models.UserEntity{
		UserID:   userEntityInput.UserID,
		EntityID: userEntityInput.EntityID,
	}

	initializers.DB.Create(&userEntity)

	c.JSON(http.StatusOK, gin.H{"data": userEntity})

}

func AddEmployee(c *gin.Context) {
	var employeeInput schemas.EmployeeInput

	logger := loggers.InitializeLogger()

	// Bind input data
	if err := c.ShouldBindJSON(&employeeInput); err != nil {
		logger.Error("Error binding JSON for Add Employee", "error", err.Error())
		c.Error(models.WrapError(http.StatusBadRequest, errors.ErrBindingJSON, "Invalid request format"))
		return
	}

	var employeeFound models.Employee
	initializers.DB.Where("email=?", employeeInput.Email).Find(&employeeFound)

	if employeeFound.ID != 0 {
		logger.Warn("Employee with this email already exists", "name", employeeFound.FirstName)
		c.Error(models.WrapError(http.StatusBadRequest, errors.ErrObjectExists, "Employee with this email already exist"))
		return
	}

	// Create the employee
	employee := models.Employee{
		FirstName:          employeeInput.FirstName,
		LastName:           employeeInput.LastName,
		Email:              employeeInput.Email,
		PhoneNumber:        employeeInput.PhoneNumber,
		DateOfBirth:        employeeInput.DateOfBirth,
		EntityID:           employeeInput.EntityID,
		EmployeeCategoryID: employeeInput.EmployeeCategoryID, // Use the EmployeeCategory ID to define the role
	}

	initializers.DB.Create(&employee)

	logger.Info("Employee added successfully", "employee ID", employee.ID)

	c.JSON(http.StatusOK, gin.H{"data": employee.ID, "message": "Sucessfully created an employee", "status": "Success"})

}

func AddPatient(c *gin.Context) {
	var patientInput schemas.PatientInput

	logger := loggers.InitializeLogger()

	// Bind input data
	if err := c.ShouldBindJSON(&patientInput); err != nil {
		logger.Error("Error binding JSON for Add Patient", "error", err.Error())
		c.Error(models.WrapError(http.StatusBadRequest, errors.ErrBindingJSON, "Invalid request format"))
		return
	}

	// Parse DateOfBirth
	dateOfBirth, err := time.Parse("2006-01-02", patientInput.DateOfBirth) // Assuming format is YYYY-MM-DD
	if err != nil {
		logger.Error("Error parsing DateOfBirth", "error", err.Error(), "date_of_birth", patientInput.DateOfBirth)
		c.Error(models.WrapError(http.StatusBadRequest, errors.InternalServerError, "Invalid DateOfBirth format"))
		return
	}

	// Check if patient with the same email already exists
	var patientFound models.Patient
	initializers.DB.Where("email = ?", patientInput.Email).Find(&patientFound)

	if patientFound.ID != 0 {
		logger.Warn("Patient with this email already exists", "email", patientInput.Email)
		c.Error(models.WrapError(http.StatusBadRequest, errors.ErrObjectExists, "Patient with this email already exists"))
		return
	}

	// Create the patient
	patient := models.Patient{
		FirstName:     patientInput.FirstName,
		LastName:      patientInput.LastName,
		Gender:        patientInput.Gender,
		DateOfBirth:   dateOfBirth,
		ContactNumber: patientInput.ContactNumber,
		Email:         patientInput.Email,
		Address:       patientInput.Address,
		EntityID:      patientInput.EntityID,
		MaritalStatus: patientInput.MaritalStatus,
		Occupation:    patientInput.Occupation,
		DoctorID:      patientInput.DoctorID,
	}

	initializers.DB.Create(&patient)

	logger.Info("Patient added successfully", "email", patientInput.Email, "patient_id", patient.ID)

	c.JSON(http.StatusOK, gin.H{"data": patient})
}

func GetEmployeeList(c *gin.Context) {
	var employees []models.Employee
	entityID := c.DefaultQuery("entity_id", "0")                      // Entity ID from query parameters
	EmployeeCategoryID := c.DefaultQuery("employee_category_id", "0") // Entity ID from query parameters

	if entityID == "0" {
		c.Error(models.WrapError(http.StatusBadRequest, errors.ErrBadRequest, "Entity ID is required"))
		return
	}

	// Convert entityID to uint
	entityIDUint, err := strconv.ParseUint(entityID, 10, 32)
	if err != nil {
		c.Error(models.WrapError(http.StatusBadRequest, errors.ErrBadRequest, "Invalid Entity ID"))
		return
	}

	// Find employees belonging to the specified entity (entity_id)
	initializers.DB.Where("entity_id = ? AND employee_category_id = ?", entityIDUint, EmployeeCategoryID).Find(&employees)

	if len(employees) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No employees found for the given entity"})
		return
	}

	// Create a simplified response with only the ID and name of the employees
	var employeeList []map[string]interface{}
	for _, employee := range employees {
		employeeList = append(employeeList, map[string]interface{}{
			"id":   employee.ID,
			"name": employee.FirstName + " " + employee.LastName,
		})
	}

	// Return the list of employees
	c.JSON(http.StatusOK, gin.H{"data": employeeList, "status": "success", "message": "Emplyee list returned succefully"})
}

func GetPatientList(c *gin.Context) {
	var patients []models.Patient
	entityID := c.DefaultQuery("entity_id", "0")

	if entityID == "0" {
		c.Error(models.WrapError(http.StatusBadRequest, errors.ErrBadRequest, "entity ID is required"))
		return
	}

	// Find patients assigned to the specified doctor
	initializers.DB.Where("entity_id = ?", entityID).Find(&patients)

	if len(patients) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No patients found for the given doctor"})
		return
	}

	// Create a simplified response with only the ID, name, and doctor name of the patients
	var patientList []map[string]interface{}
	for _, patient := range patients {
		patientList = append(patientList, map[string]interface{}{
			"id":     patient.ID,
			"name":   patient.FirstName + " " + patient.LastName,
			"doctor": patient.Doctor.FirstName + " " + patient.Doctor.LastName, // Doctor's name
		})
	}

	// Return the simplified list of patients
	c.JSON(http.StatusOK, gin.H{"data": patientList})
}
