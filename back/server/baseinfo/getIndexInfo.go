package baseinfo

import (
	"context"
	"hospital-system/models"
	"hospital-system/resource"
)

type IndexInfo struct {
	Patients      []models.Patient      `json:"patients"`
	Doctors       []models.Doctor       `json:"doctors"`
	Registrations []models.Registration `json:"registrations"`
}

func GetIndexInfo(ctx context.Context) ([]IndexInfo, error) {
	var indexInfo []IndexInfo

	// 获取所有患者信息
	patients, err := resource.PatientService.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// 获取所有医生信息
	doctors, err := resource.DoctorService.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// 获取所有挂号信息
	registrations, err := resource.RegistrationService.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	indexInfo = append(indexInfo, IndexInfo{
		Patients:      patients,
		Doctors:       doctors,
		Registrations: registrations,
	})
	return indexInfo, nil
}
