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
)

// PatientService 病人管理的单例
var PatientService *services.PatientService
