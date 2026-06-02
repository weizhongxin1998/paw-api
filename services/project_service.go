package services

import (
	"paw-api/models"
	"paw-api/pkg/snowflake"
	"paw-api/repositories"
)

type ProjectService struct {
	repo *repositories.ProjectRepository
	sf   *snowflake.Generator
}

func NewProjectService(repo *repositories.ProjectRepository, sf *snowflake.Generator) *ProjectService {
	return &ProjectService{repo: repo, sf: sf}
}

func (s *ProjectService) List() ([]models.Project, error) {
	return s.repo.List()
}

func (s *ProjectService) Get(id int64) (*models.Project, error) {
	return s.repo.GetByID(id)
}

func (s *ProjectService) Create(name, description string) (*models.Project, error) {
	project := &models.Project{
		Name:        name,
		Description: description,
	}
	if err := s.repo.Create(project); err != nil {
		return nil, err
	}

	// 自动创建默认集合
	defaultColl := &models.Collection{
		ProjectID: project.ID,
		Name:      "Default Collection",
	}
	// Use the collection repo if available, or skip for now
	// The handler/app will wire this
	_ = defaultColl

	return project, nil
}

func (s *ProjectService) Update(id int64, name, description string) (*models.Project, error) {
	project, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	project.Name = name
	project.Description = description
	if err := s.repo.Update(project); err != nil {
		return nil, err
	}
	return project, nil
}

func (s *ProjectService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *ProjectService) GetStats(id int64) (models.ProjectStats, error) {
	return s.repo.GetStats(id)
}
