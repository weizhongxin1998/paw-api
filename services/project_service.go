package services

import (
	"errors"
	"time"
	"paw-api/models"
	"paw-api/repositories"

	"github.com/google/uuid"
)

type ProjectService struct {
	repo *repositories.ProjectRepo
}

func NewProjectService() *ProjectService {
	return &ProjectService{repo: &repositories.ProjectRepo{}}
}

func (s *ProjectService) Create(name, description string) (*models.Project, error) {
	if name == "" {
		return nil, errors.New("project name is required")
	}
	now := time.Now()
	p := &models.Project{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	return p, s.repo.Create(p)
}

func (s *ProjectService) GetByID(id string) (*models.Project, error) {
	return s.repo.GetByID(id)
}

func (s *ProjectService) List() ([]models.Project, error) {
	return s.repo.List()
}

func (s *ProjectService) Update(id, name, description string) (*models.Project, error) {
	if name == "" {
		return nil, errors.New("project name is required")
	}
	p, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errors.New("project not found")
	}
	p.Name = name
	p.Description = description
	p.UpdatedAt = time.Now()
	return p, s.repo.Update(p)
}

func (s *ProjectService) Delete(id string) error {
	return s.repo.Delete(id)
}
