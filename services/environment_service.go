package services

import (
	"errors"
	"time"
	"paw-api/models"
	"paw-api/repositories"

	"github.com/google/uuid"
)

type EnvironmentService struct {
	repo *repositories.EnvironmentRepo
}

func NewEnvironmentService() *EnvironmentService {
	return &EnvironmentService{repo: &repositories.EnvironmentRepo{}}
}

func (s *EnvironmentService) Create(projectID, name, variables string, isActive bool) (*models.Environment, error) {
	if name == "" {
		return nil, errors.New("environment name is required")
	}
	now := time.Now()
	e := &models.Environment{
		ID:        uuid.New().String(),
		ProjectID: projectID,
		Name:      name,
		Variables: variables,
		IsActive:  isActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return e, s.repo.Create(e)
}

func (s *EnvironmentService) GetByID(id string) (*models.Environment, error) {
	return s.repo.GetByID(id)
}

func (s *EnvironmentService) ListByProject(projectID string) ([]models.Environment, error) {
	return s.repo.ListByProject(projectID)
}

func (s *EnvironmentService) Update(id, name, variables string) (*models.Environment, error) {
	if name == "" {
		return nil, errors.New("environment name is required")
	}
	e, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, errors.New("environment not found")
	}
	e.Name = name
	e.Variables = variables
	e.UpdatedAt = time.Now()
	return e, s.repo.Update(e)
}

func (s *EnvironmentService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *EnvironmentService) SetActive(id, projectID string) (*models.Environment, error) {
	if err := s.repo.DeactivateAll(projectID); err != nil {
		return nil, err
	}
	e, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, errors.New("environment not found")
	}
	e.IsActive = true
	e.UpdatedAt = time.Now()
	return e, s.repo.Update(e)
}
