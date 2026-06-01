package handlers

import (
	"paw-api/models"
	"paw-api/services"
)

type HistoryHandler struct {
	service *services.HistoryService
}

func NewHistoryHandler() *HistoryHandler {
	return &HistoryHandler{service: services.NewHistoryService()}
}

func (h *HistoryHandler) RecordHistory(projectID, method, url, headers, body string) (*models.History, error) {
	return h.service.Record(projectID, method, url, headers, body)
}

func (h *HistoryHandler) ListHistory(projectID string, limit int) ([]models.History, error) {
	return h.service.ListByProject(projectID, limit)
}

func (h *HistoryHandler) DeleteHistory(id string) error {
	return h.service.Delete(id)
}

func (h *HistoryHandler) ClearHistory(projectID string) error {
	return h.service.ClearByProject(projectID)
}
