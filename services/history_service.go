package services

import (
	"time"
	"paw-api/models"
	"paw-api/repositories"

	"github.com/google/uuid"
)

type HistoryService struct {
	repo *repositories.HistoryRepo
}

func NewHistoryService() *HistoryService {
	return &HistoryService{repo: &repositories.HistoryRepo{}}
}

func (s *HistoryService) Record(projectID, requestID, method, url, headers, body string, responseStatus int, responseBody, responseHeaders string, durationMs int) (*models.History, error) {
	now := time.Now()
	var rid *string
	if requestID != "" {
		rid = &requestID
	}
	h := &models.History{
		ID:              uuid.New().String(),
		ProjectID:       projectID,
		RequestID:       rid,
		Method:          method,
		URL:             url,
		Headers:         headers,
		Body:            body,
		ResponseStatus:  responseStatus,
		ResponseBody:    responseBody,
		ResponseHeaders: responseHeaders,
		DurationMs:      durationMs,
		CreatedAt:       now,
	}
	return h, s.repo.Create(h)
}

func (s *HistoryService) ListByProject(projectID string, limit int) ([]models.History, error) {
	return s.repo.ListByProject(projectID, limit)
}

func (s *HistoryService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *HistoryService) ClearByProject(projectID string) error {
	return s.repo.ClearByProject(projectID)
}
