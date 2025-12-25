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

type DiseaseService struct {
	filename string
	mu       sync.RWMutex
}

func InitDiseaseService(c *DiseaseService) *DiseaseService {
	if c == nil || c.filename == "" {
		return &DiseaseService{
			filename: "static/diseases.json",
		}
	}
	return c
}

func (s *DiseaseService) readAll() ([]models.Disease, error) {
	data, err := os.ReadFile(s.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.Disease{}, nil
		}
		return nil, err
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return []models.Disease{}, nil
	}

	var diseases []models.Disease
	if err := json.Unmarshal(data, &diseases); err != nil {
		return nil, err
	}

	return diseases, nil
}

func (s *DiseaseService) GetAll(ctx context.Context) ([]models.Disease, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.readAll()
}

func (s *DiseaseService) GetByID(ctx context.Context, id string) (*models.Disease, error) {
	diseases, err := s.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, disease := range diseases {
		if disease.ID == id {
			return &disease, nil
		}
	}

	return nil, errors.New("disease not found")
}

func (s *DiseaseService) Create(ctx context.Context, disease *models.Disease) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	diseases, err := s.readAll()
	if err != nil {
		return err
	}

	if disease.Name == "" {
		return errors.New("name cannot be empty")
	}
	if disease.Category == "" {
		return errors.New("category cannot be empty")
	}

	disease.ID = uuid.New().String()
	diseases = append(diseases, *disease)

	data, err := json.MarshalIndent(diseases, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filename, data, 0644)
}

func (s *DiseaseService) Update(ctx context.Context, id string, updatedDisease *models.Disease) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	diseases, err := s.readAll()
	if err != nil {
		return err
	}

	if updatedDisease.Name == "" {
		return errors.New("name cannot be empty")
	}
	if updatedDisease.Category == "" {
		return errors.New("category cannot be empty")
	}

	found := false
	for i, disease := range diseases {
		if disease.ID == id {
			updatedDisease.ID = id
			diseases[i] = *updatedDisease
			found = true
			break
		}
	}
	if !found {
		return errors.New("disease not found")
	}

	data, err := json.MarshalIndent(diseases, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filename, data, 0644)
}

func (s *DiseaseService) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	diseases, err := s.readAll()
	if err != nil {
		return err
	}

	newDiseases := make([]models.Disease, 0, len(diseases))
	for _, disease := range diseases {
		if disease.ID != id {
			newDiseases = append(newDiseases, disease)
		}
	}

	data, err := json.MarshalIndent(newDiseases, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filename, data, 0644)
}
