package services

import "paw-api/repositories"

type SettingsService struct {
	repo *repositories.SettingsRepository
}

func NewSettingsService(repo *repositories.SettingsRepository) *SettingsService {
	return &SettingsService{repo: repo}
}

func (s *SettingsService) Get(key string) (string, error) {
	return s.repo.Get(key)
}

func (s *SettingsService) Set(key, value string) error {
	return s.repo.Set(key, value)
}

func (s *SettingsService) GetAll() (map[string]string, error) {
	return s.repo.GetAll()
}
