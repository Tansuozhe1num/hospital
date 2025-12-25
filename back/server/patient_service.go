package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"hospital-system/models"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/google/uuid"
)

type PatientService struct {
	filename string
	mu       sync.RWMutex
}

func InitPatientService(c *PatientService) *PatientService {
	if c == nil || c.filename == "" {
		return &PatientService{
			filename: "static/patients.json",
		}
	}
	return c
}

func (s *PatientService) readAll() ([]models.Patient, error) {
	data, err := os.ReadFile(s.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.Patient{}, nil
		}
		return nil, err
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return []models.Patient{}, nil
	}

	var patients []models.Patient
	if err := json.Unmarshal(data, &patients); err != nil {
		return nil, err
	}

	return patients, nil
}

func (s *PatientService) GetAll(ctx context.Context) ([]models.Patient, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.readAll()
}

func (s *PatientService) GetByID(ctx context.Context, id string) (*models.Patient, error) {
	patients, err := s.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, patient := range patients {
		if patient.ID == id {
			return &patient, nil
		}
	}

	return nil, errors.New("patient not found")
}

func (s *PatientService) Create(ctx context.Context, patient *models.Patient) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	patients, err := s.readAll()
	if err != nil {
		return err
	}

	if ok, err := checkPatientValidity(patient); !ok {
		return err
	}

	patient.ID = uuid.New().String()
	now := time.Now()
	if patient.CreatedAt.IsZero() {
		patient.CreatedAt = now
	}
	patient.UpdatedAt = now

	patients = append(patients, *patient)

	data, err := json.MarshalIndent(patients, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filename, data, 0644)
}

func (s *PatientService) Update(ctx context.Context, id string, updatedPatient *models.Patient) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	patients, err := s.readAll()
	if err != nil {
		return err
	}

	// 小系统，用最笨的办法写的，直接全部取出来遍历一遍，实际系统要是过1e7的话要卡死，，没链接数据库就这么弄了。。
	found := false
	for i, patient := range patients {
		if patient.ID == id {
			updatedPatient.ID = id
			now := time.Now()
			if updatedPatient.CreatedAt.IsZero() {
				updatedPatient.CreatedAt = patient.CreatedAt
			}
			updatedPatient.UpdatedAt = now
			patients[i] = *updatedPatient
			found = true
			break
		}
	}

	if !found {
		return errors.New("patient not found")
	}

	data, err := json.MarshalIndent(patients, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filename, data, 0644)
}

func (s *PatientService) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	patients, err := s.readAll()
	if err != nil {
		return err
	}

	var newPatients = []models.Patient{}
	for _, patient := range patients {
		if patient.ID != id {
			newPatients = append(newPatients, patient)
		}
	}

	data, err := json.MarshalIndent(newPatients, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filename, data, 0644)
}

func checkPatientValidity(patient *models.Patient) (ok bool, err error) {
	// 检查病人的输入信息是否合规，是才可以录入，否则返回false
	if patient.Name == "" {
		return false, errors.New("name cannot be empty")
	}

	if patient.Age < 1 || patient.Age > 150 {
		return false, errors.New("age must be between 1 and 150")
	}

	if patient.Gender != "男" && patient.Gender != "女" {
		return false, errors.New("gender must be either '男' or '女'")
	}

	if len(patient.Phone) != 11 {
		return false, errors.New("phone number must be 11 digits")
	}

	if !regexp.MustCompile(`^[1-9]\d{5}(18|19|([23]\d))\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9xX]$`).MatchString(patient.IDCard) {
		return false, errors.New("ID card number must be 18 digits")
	}

	if patient.Address == "" {
		return false, errors.New("address cannot be empty")
	}

	if patient.EmergencyContact == "" {
		return false, errors.New("emergency contact cannot be empty")
	}

	return true, nil
}
