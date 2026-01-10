package resource

import (
	"hospital-system/models"
	services "hospital-system/server"
)

var (
	Patients      = []models.Patient{}
	Doctors       = []models.Doctor{}
	Registrations = []models.Registration{}
	Diseases      = []models.Disease{}
	Departments   = []models.Department{}
)

var (
	PatientService      *services.PatientService
	DiseaseService      *services.DiseaseService
	DoctorService       *services.DoctorService
	RegistrationService *services.RegistrationService
	AccountService      *services.AccountService
	DepartmentService   *services.DepartmentService
)
