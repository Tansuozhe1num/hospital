package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"hospital-system/models"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AccountService struct {
	filename string
	mu       sync.RWMutex
}

func InitAccountService(c *AccountService) *AccountService {
	if c == nil || c.filename == "" {
		return &AccountService{
			filename: "static/accounts.json",
		}
	}
	return c
}

func (s *AccountService) readAll() ([]models.Account, error) {
	data, err := os.ReadFile(s.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.Account{}, nil
		}
		return nil, err
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return []models.Account{}, nil
	}

	var accounts []models.Account
	if err := json.Unmarshal(data, &accounts); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *AccountService) GetAll(ctx context.Context) ([]models.Account, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.readAll()
}

func (s *AccountService) GetByID(ctx context.Context, id string) (*models.Account, error) {
	accounts, err := s.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	for _, a := range accounts {
		if a.ID == id {
			return &a, nil
		}
	}
	return nil, errors.New("account not found")
}

func (s *AccountService) GetByUsername(ctx context.Context, username string) (*models.Account, error) {
	username = normalizeUsername(username)
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	accounts, err := s.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	for _, a := range accounts {
		if normalizeUsername(a.Username) == username {
			return &a, nil
		}
	}
	return nil, errors.New("account not found")
}

func (s *AccountService) EnsureAccount(ctx context.Context, username string, password string, role string) error {
	username = normalizeUsername(username)
	password = strings.TrimSpace(password)
	role = strings.TrimSpace(role)
	if username == "" {
		return errors.New("username cannot be empty")
	}
	if password == "" {
		return errors.New("password cannot be empty")
	}
	switch role {
	case "patient", "doctor", "admin":
	default:
		return errors.New("invalid role")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	accounts, err := s.readAll()
	if err != nil {
		return err
	}

	for _, a := range accounts {
		if normalizeUsername(a.Username) == username {
			roleMatches := strings.TrimSpace(a.Role) == role
			passwordMatches := bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(password)) == nil
			if roleMatches && passwordMatches {
				return nil
			}

			hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}

			now := time.Now()
			for i := range accounts {
				if normalizeUsername(accounts[i].Username) == username {
					accounts[i].PasswordHash = string(hashBytes)
					accounts[i].Role = role
					accounts[i].UpdatedAt = now
					if accounts[i].CreatedAt.IsZero() {
						accounts[i].CreatedAt = now
					}
					break
				}
			}

			return s.writeAll(accounts)
		}
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	now := time.Now()
	accounts = append(accounts, models.Account{
		ID:           uuid.New().String(),
		Username:     username,
		PasswordHash: string(hashBytes),
		Role:         role,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	return s.writeAll(accounts)
}

func (s *AccountService) UpsertDoctorAccount(ctx context.Context, doctorID string, username string, password string) (*models.Account, error) {
	doctorID = strings.TrimSpace(doctorID)
	username = normalizeUsername(username)
	password = strings.TrimSpace(password)
	if doctorID == "" {
		return nil, errors.New("doctorId cannot be empty")
	}
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	if password == "" {
		return nil, errors.New("password cannot be empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	accounts, err := s.readAll()
	if err != nil {
		return nil, err
	}

	existingIndex := -1
	for i := range accounts {
		if strings.TrimSpace(accounts[i].LinkedID) == doctorID && strings.TrimSpace(accounts[i].Role) == "doctor" {
			existingIndex = i
			break
		}
	}

	for i := range accounts {
		if normalizeUsername(accounts[i].Username) != username {
			continue
		}
		if existingIndex >= 0 && accounts[i].ID == accounts[existingIndex].ID {
			continue
		}
		return nil, errors.New("username already exists")
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	now := time.Now()

	if existingIndex >= 0 {
		accounts[existingIndex].Username = username
		accounts[existingIndex].PasswordHash = string(hashBytes)
		accounts[existingIndex].Role = "doctor"
		accounts[existingIndex].LinkedID = doctorID
		accounts[existingIndex].UpdatedAt = now
		if accounts[existingIndex].CreatedAt.IsZero() {
			accounts[existingIndex].CreatedAt = now
		}
		if err := s.writeAll(accounts); err != nil {
			return nil, err
		}
		return &accounts[existingIndex], nil
	}

	account := models.Account{
		ID:           uuid.New().String(),
		Username:     username,
		PasswordHash: string(hashBytes),
		Role:         "doctor",
		LinkedID:     doctorID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	accounts = append(accounts, account)
	if err := s.writeAll(accounts); err != nil {
		return nil, err
	}
	return &account, nil
}

func (s *AccountService) LoginOrRegister(ctx context.Context, username string, password string) (*models.Account, bool, error) {
	username = normalizeUsername(username)
	password = strings.TrimSpace(password)
	if username == "" {
		return nil, false, errors.New("username cannot be empty")
	}
	if password == "" {
		return nil, false, errors.New("password cannot be empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	accounts, err := s.readAll()
	if err != nil {
		return nil, false, err
	}

	for i := range accounts {
		if normalizeUsername(accounts[i].Username) != username {
			continue
		}
		if err := bcrypt.CompareHashAndPassword([]byte(accounts[i].PasswordHash), []byte(password)); err != nil {
			return nil, false, errors.New("invalid credentials")
		}
		now := time.Now()
		accounts[i].UpdatedAt = now
		if err := s.writeAll(accounts); err != nil {
			return nil, false, err
		}
		return &accounts[i], false, nil
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, false, err
	}
	now := time.Now()
	account := models.Account{
		ID:           uuid.New().String(),
		Username:     username,
		PasswordHash: string(hashBytes),
		Role:         "patient",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	accounts = append(accounts, account)
	if err := s.writeAll(accounts); err != nil {
		return nil, false, err
	}
	return &account, true, nil
}

func (s *AccountService) writeAll(accounts []models.Account) error {
	data, err := json.MarshalIndent(accounts, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filename, data, 0644)
}

func normalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}
