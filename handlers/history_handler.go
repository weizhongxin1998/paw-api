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

type RecordHistoryInput struct {
	ProjectID       string
	RequestID       string
	Method          string
	URL             string
	Headers         string
	Body            string
	ResponseStatus  int
	ResponseBody    string
	ResponseHeaders string
	DurationMs      int
}

func (h *HistoryHandler) RecordHistory(input RecordHistoryInput) (*models.History, error) {
	return h.service.Record(input.ProjectID, input.RequestID, input.Method, input.URL,
		input.Headers, input.Body, input.ResponseStatus, input.ResponseBody,
		input.ResponseHeaders, input.DurationMs)
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
