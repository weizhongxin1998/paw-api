package services

import (
	"paw-api/models"
	"paw-api/pkg/snowflake"
	"paw-api/repositories"
)

type EnvironmentService struct {
	envRepo  *repositories.EnvironmentRepository
	varRepo  *repositories.EnvVariableRepository
	sf       *snowflake.Generator
}

func NewEnvironmentService(
	envRepo *repositories.EnvironmentRepository,
	varRepo *repositories.EnvVariableRepository,
	sf *snowflake.Generator,
) *EnvironmentService {
	return &EnvironmentService{envRepo: envRepo, varRepo: varRepo, sf: sf}
}

func (s *EnvironmentService) List(projectID int64) ([]models.Environment, error) {
	return s.envRepo.ListByProject(projectID)
}

func (s *EnvironmentService) Create(projectID int64, name string, baseURL string, cloneFromID *int64) (*models.Environment, error) {
	env := &models.Environment{
		ProjectID: projectID,
		Name:      name,
		BaseURL:   baseURL,
	}
	if err := s.envRepo.Create(env); err != nil {
		return nil, err
	}

	if cloneFromID != nil {
		vars, err := s.varRepo.ListByEnvironment(*cloneFromID)
		if err != nil {
			return env, nil
		}
		cloned := make([]models.EnvVariable, len(vars))
		for i, v := range vars {
			cloned[i] = v
			cloned[i].ID = 0
		}
		_ = s.varRepo.BatchReplace(env.ID, cloned)
	}

	return env, nil
}

func (s *EnvironmentService) Rename(id int64, name string) error {
	env, err := s.envRepo.GetByID(id)
	if err != nil {
		return err
	}
	env.Name = name
	return s.envRepo.Update(env)
}

func (s *EnvironmentService) Delete(id int64) error {
	return s.envRepo.Delete(id)
}

func (s *EnvironmentService) Activate(id int64) error {
	env, err := s.envRepo.GetByID(id)
	if err != nil {
		return err
	}
	return s.envRepo.SetActive(env.ProjectID, id)
}

func (s *EnvironmentService) GetActive(projectID int64) (*models.Environment, error) {
	return s.envRepo.GetActive(projectID)
}

func (s *EnvironmentService) ListVariables(envID int64) ([]models.EnvVariable, error) {
	return s.varRepo.ListByEnvironment(envID)
}

func (s *EnvironmentService) SaveVariables(envID int64, variables []models.EnvVariable) error {
	return s.varRepo.BatchReplace(envID, variables)
}

func (s *EnvironmentService) SaveBaseURL(envID int64, baseURL string) error {
	return s.envRepo.SaveBaseURL(envID, baseURL)
}
