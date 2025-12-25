package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"hospital-system/models"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

type RegistrationService struct {
	filename string
	mu       sync.RWMutex
}

func InitRegistrationService(c *RegistrationService) *RegistrationService {
	if c == nil || c.filename == "" {
		return &RegistrationService{
			filename: "static/registrations.json",
		}
	}
	return c
}

func normalizeDepartments(r *models.Registration) {
	if r == nil {
		return
	}
	if len(r.Departments) == 0 && r.Department != "" {
		r.Departments = []string{r.Department}
	}
	if r.Department == "" && len(r.Departments) > 0 {
		r.Department = r.Departments[0]
	}
}

func (s *RegistrationService) readAll() ([]models.Registration, error) {
	data, err := os.ReadFile(s.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.Registration{}, nil
		}
		return nil, err
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return []models.Registration{}, nil
	}

	var registrations []models.Registration
	if err := json.Unmarshal(data, &registrations); err != nil {
		return nil, err
	}

	for i := range registrations {
		normalizeDepartments(&registrations[i])
	}

	return registrations, nil
}

func (s *RegistrationService) GetAll(ctx context.Context) ([]models.Registration, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.readAll()
}

func (s *RegistrationService) GetByID(ctx context.Context, id string) (*models.Registration, error) {
	registrations, err := s.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, registration := range registrations {
		if registration.ID == id {
			return &registration, nil
		}
	}

	return nil, errors.New("registration not found")
}

func (s *RegistrationService) Create(ctx context.Context, registration *models.Registration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	registrations, err := s.readAll()
	if err != nil {
		return err
	}

	normalizeDepartments(registration)

	if registration.PatientID == "" {
		return errors.New("patientId cannot be empty")
	}
	if registration.DoctorID == "" {
		return errors.New("doctorId cannot be empty")
	}
	if len(registration.Departments) == 0 {
		return errors.New("departments cannot be empty")
	}
	if registration.VisitDate.IsZero() {
		return errors.New("visitDate cannot be empty")
	}
	if registration.TimeSlot == "" {
		return errors.New("timeSlot cannot be empty")
	}
	if registration.Status == "" {
		registration.Status = "pending"
	}
	if registration.Status != "pending" && registration.Status != "confirmed" && registration.Status != "completed" && registration.Status != "cancelled" {
		return errors.New("invalid status")
	}

	registration.ID = uuid.New().String()
	now := time.Now()
	if registration.RegistrationDate.IsZero() {
		registration.RegistrationDate = now
	}
	if registration.CreatedAt.IsZero() {
		registration.CreatedAt = now
	}

	registrations = append(registrations, *registration)

	data, err := json.MarshalIndent(registrations, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filename, data, 0644)
}

func (s *RegistrationService) Update(ctx context.Context, id string, updatedRegistration *models.Registration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	registrations, err := s.readAll()
	if err != nil {
		return err
	}

	normalizeDepartments(updatedRegistration)

	if updatedRegistration.PatientID == "" {
		return errors.New("patientId cannot be empty")
	}
	if updatedRegistration.DoctorID == "" {
		return errors.New("doctorId cannot be empty")
	}
	if len(updatedRegistration.Departments) == 0 {
		return errors.New("departments cannot be empty")
	}
	if updatedRegistration.VisitDate.IsZero() {
		return errors.New("visitDate cannot be empty")
	}
	if updatedRegistration.TimeSlot == "" {
		return errors.New("timeSlot cannot be empty")
	}
	if updatedRegistration.Status == "" {
		updatedRegistration.Status = "pending"
	}
	if updatedRegistration.Status != "pending" && updatedRegistration.Status != "confirmed" && updatedRegistration.Status != "completed" && updatedRegistration.Status != "cancelled" {
		return errors.New("invalid status")
	}

	found := false
	for i, registration := range registrations {
		if registration.ID == id {
			updatedRegistration.ID = id
			if updatedRegistration.RegistrationDate.IsZero() {
				updatedRegistration.RegistrationDate = registration.RegistrationDate
			}
			if updatedRegistration.CreatedAt.IsZero() {
				updatedRegistration.CreatedAt = registration.CreatedAt
			}
			registrations[i] = *updatedRegistration
			found = true
			break
		}
	}
	if !found {
		return errors.New("registration not found")
	}

	data, err := json.MarshalIndent(registrations, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filename, data, 0644)
}

func (s *RegistrationService) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	registrations, err := s.readAll()
	if err != nil {
		return err
	}

	newRegistrations := make([]models.Registration, 0, len(registrations))
	for _, registration := range registrations {
		if registration.ID != id {
			newRegistrations = append(newRegistrations, registration)
		}
	}

	data, err := json.MarshalIndent(newRegistrations, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filename, data, 0644)
}
