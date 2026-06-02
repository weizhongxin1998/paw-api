package services

import (
	"paw-api/models"
	"paw-api/repositories"
)

type HistoryService struct {
	repo *repositories.HistoryRepository
}

func NewHistoryService(repo *repositories.HistoryRepository) *HistoryService {
	return &HistoryService{repo: repo}
}

func (s *HistoryService) List(projectID int64, page, pageSize int) ([]models.History, int, error) {
	offset := (page - 1) * pageSize
	items, err := s.repo.ListByProject(projectID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	allItems, err := s.repo.ListByProject(projectID, 0, 0)
	if err != nil {
		return nil, 0, err
	}

	return items, len(allItems), nil
}

func (s *HistoryService) Search(projectID int64, keyword, method string, statusMin, statusMax int, page, pageSize int) ([]models.History, int, error) {
	offset := (page - 1) * pageSize
	return s.repo.Search(projectID, keyword, method, statusMin, statusMax, pageSize, offset)
}

func (s *HistoryService) Get(id int64) (*models.History, error) {
	return s.repo.GetByID(id)
}

func (s *HistoryService) Clear(projectID int64) error {
	return s.repo.DeleteByProject(projectID)
}

func (s *HistoryService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *HistoryService) CleanOld(projectID int64, days int) (int, error) {
	if err := s.repo.DeleteOlderThan(projectID, days); err != nil {
		return 0, err
	}
	return 0, nil
}
