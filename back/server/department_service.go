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

type DepartmentService struct {
	filename string
	mu       sync.RWMutex
}

func InitDepartmentService(c *DepartmentService) *DepartmentService {
	if c == nil || c.filename == "" {
		return &DepartmentService{
			filename: "static/departments.json",
		}
	}
	return c
}

func (s *DepartmentService) readAll() ([]models.Department, error) {
	data, err := os.ReadFile(s.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.Department{}, nil
		}
		return nil, err
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return []models.Department{}, nil
	}

	var departments []models.Department
	if err := json.Unmarshal(data, &departments); err != nil {
		return nil, err
	}
	return departments, nil
}

func (s *DepartmentService) writeAll(departments []models.Department) error {
	data, err := json.MarshalIndent(departments, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filename, data, 0644)
}

func (s *DepartmentService) GetAll(ctx context.Context) ([]models.Department, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.readAll()
}

func (s *DepartmentService) GetByID(ctx context.Context, id string) (*models.Department, error) {
	departments, err := s.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	for _, d := range departments {
		if d.ID == id {
			return &d, nil
		}
	}
	return nil, errors.New("department not found")
}

func (s *DepartmentService) Create(ctx context.Context, d *models.Department) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	departments, err := s.readAll()
	if err != nil {
		return err
	}
	if d.Name == "" {
		return errors.New("name cannot be empty")
	}
	for _, existing := range departments {
		if existing.Name == d.Name {
			return errors.New("department name already exists")
		}
	}
	d.ID = uuid.New().String()
	departments = append(departments, *d)
	return s.writeAll(departments)
}

func (s *DepartmentService) Update(ctx context.Context, id string, updated *models.Department) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	departments, err := s.readAll()
	if err != nil {
		return err
	}
	if updated.Name == "" {
		return errors.New("name cannot be empty")
	}
	// name uniqueness check excluding current id
	for _, existing := range departments {
		if existing.ID != id && existing.Name == updated.Name {
			return errors.New("department name already exists")
		}
	}
	found := false
	for i, d := range departments {
		if d.ID == id {
			updated.ID = id
			departments[i] = *updated
			found = true
			break
		}
	}
	if !found {
		return errors.New("department not found")
	}
	return s.writeAll(departments)
}

func (s *DepartmentService) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	departments, err := s.readAll()
	if err != nil {
		return err
	}
	newDepartments := make([]models.Department, 0, len(departments))
	for _, d := range departments {
		if d.ID != id {
			newDepartments = append(newDepartments, d)
		}
	}
	return s.writeAll(newDepartments)
}

