package handlers

import (
	"paw-api/models"
	"paw-api/services"
)

type HistoryHandler struct {
	svc *services.HistoryService
}

func NewHistoryHandler(svc *services.HistoryService) *HistoryHandler {
	return &HistoryHandler{svc: svc}
}

func (h *HistoryHandler) List(projectID int64, page, pageSize int) ([]models.History, error) {
	items, _, err := h.svc.List(projectID, page, pageSize)
	return items, err
}

func (h *HistoryHandler) Search(projectID int64, keyword, method string, statusMin, statusMax int, page, pageSize int) ([]models.History, int, error) {
	return h.svc.Search(projectID, keyword, method, statusMin, statusMax, page, pageSize)
}

func (h *HistoryHandler) Get(id int64) (*models.History, error) {
	return h.svc.Get(id)
}

func (h *HistoryHandler) Clear(projectID int64) error {
	return h.svc.Clear(projectID)
}

func (h *HistoryHandler) Delete(id int64) error {
	return h.svc.Delete(id)
}

func (h *HistoryHandler) CleanOld(projectID int64, days int) error {
	_, err := h.svc.CleanOld(projectID, days)
	return err
}
