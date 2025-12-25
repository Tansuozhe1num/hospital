package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"hospital-system/models"
	"os"
	"sync"

	"github.com/google/uuid"
)

type DoctorService struct {
	filename string
	mu       sync.RWMutex
}

func InitDoctorService(c *DoctorService) *DoctorService {
	if c == nil || c.filename == "" {
		return &DoctorService{
			filename: "static/doctors.json",
		}
	}
	return c
}

func (s *DoctorService) readAll() ([]models.Doctor, error) {
	data, err := os.ReadFile(s.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.Doctor{}, nil
		}
		return nil, err
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return []models.Doctor{}, nil
	}

	var doctors []models.Doctor
	if err := json.Unmarshal(data, &doctors); err != nil {
		return nil, err
	}

	return doctors, nil
}

func (s *DoctorService) GetAll(ctx context.Context) ([]models.Doctor, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.readAll()
}

func (s *DoctorService) GetByID(ctx context.Context, id string) (*models.Doctor, error) {
	doctors, err := s.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, doctor := range doctors {
		if doctor.ID == id {
			return &doctor, nil
		}
	}

	return nil, errors.New("doctor not found")
}

func (s *DoctorService) Create(ctx context.Context, doctor *models.Doctor) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	doctors, err := s.readAll()
	if err != nil {
		return err
	}

	if doctor.Name == "" {
		return errors.New("name cannot be empty")
	}
	if doctor.Department == "" {
		return errors.New("department cannot be empty")
	}
	if doctor.Title == "" {
		return errors.New("title cannot be empty")
	}
	if len(doctor.Diseases) < 1 || len(doctor.Diseases) > 3 {
		return errors.New("diseases length must be between 1 and 3")
	}
	if doctor.MaxPatients < 1 {
		doctor.MaxPatients = 30
	}
	if doctor.Fee < 0 {
		return errors.New("fee must be >= 0")
	}
	if len(doctor.WorkSchedule) == 0 {
		doctor.WorkSchedule = defaultDoctorWorkSchedule()
	}

	doctor.ID = uuid.New().String()
	doctors = append(doctors, *doctor)

	data, err := json.MarshalIndent(doctors, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filename, data, 0644)
}

func (s *DoctorService) Update(ctx context.Context, id string, updatedDoctor *models.Doctor) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	doctors, err := s.readAll()
	if err != nil {
		return err
	}

	if updatedDoctor.Name == "" {
		return errors.New("name cannot be empty")
	}
	if updatedDoctor.Department == "" {
		return errors.New("department cannot be empty")
	}
	if updatedDoctor.Title == "" {
		return errors.New("title cannot be empty")
	}
	if len(updatedDoctor.Diseases) < 1 || len(updatedDoctor.Diseases) > 3 {
		return errors.New("diseases length must be between 1 and 3")
	}
	if updatedDoctor.MaxPatients < 1 {
		updatedDoctor.MaxPatients = 30
	}
	if updatedDoctor.Fee < 0 {
		return errors.New("fee must be >= 0")
	}
	if len(updatedDoctor.WorkSchedule) == 0 {
		updatedDoctor.WorkSchedule = defaultDoctorWorkSchedule()
	}

	found := false
	for i, doctor := range doctors {
		if doctor.ID == id {
			updatedDoctor.ID = id
			doctors[i] = *updatedDoctor
			found = true
			break
		}
	}
	if !found {
		return errors.New("doctor not found")
	}

	data, err := json.MarshalIndent(doctors, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filename, data, 0644)
}

func (s *DoctorService) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	doctors, err := s.readAll()
	if err != nil {
		return err
	}

	newDoctors := make([]models.Doctor, 0, len(doctors))
	for _, doctor := range doctors {
		if doctor.ID != id {
			newDoctors = append(newDoctors, doctor)
		}
	}

	data, err := json.MarshalIndent(newDoctors, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filename, data, 0644)
}

func defaultDoctorWorkSchedule() []models.WorkSchedule {
	return []models.WorkSchedule{
		{DayOfWeek: "周一", StartTime: "09:00", EndTime: "17:00", IsAvailable: true},
		{DayOfWeek: "周二", StartTime: "09:00", EndTime: "17:00", IsAvailable: true},
		{DayOfWeek: "周三", StartTime: "09:00", EndTime: "17:00", IsAvailable: true},
		{DayOfWeek: "周四", StartTime: "09:00", EndTime: "17:00", IsAvailable: true},
		{DayOfWeek: "周五", StartTime: "09:00", EndTime: "17:00", IsAvailable: true},
		{DayOfWeek: "周六", StartTime: "09:00", EndTime: "17:00", IsAvailable: false},
		{DayOfWeek: "周日", StartTime: "09:00", EndTime: "17:00", IsAvailable: false},
	}
}
