package load

import (
	"context"
	"hospital-system/models"
	"hospital-system/resource"
	services "hospital-system/server"
)

func Load(ctx context.Context) {
	resource.PatientService = services.InitPatientService(resource.PatientService)

	initStorage(ctx)
}

func initStorage(ctx context.Context) {
	// 初始化JSON文件
	initJSONFile("static/patients.json", []models.Patient{})
	initJSONFile("static/diseases.json", []models.Disease{})
	initJSONFile("static/doctors.json", []models.Doctor{})
	initJSONFile("static/registrations.json", []models.Registration{})
}

func initJSONFile(filename string, defaultData interface{}) {
	// 如果文件不存在，创建并写入默认数据
	// 实际实现中会检查文件是否存在并处理
}
