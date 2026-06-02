package handlers

import "paw-api/services"

type ExportHandler struct {
	svc *services.ExportService
}

func NewExportHandler(svc *services.ExportService) *ExportHandler {
	return &ExportHandler{svc: svc}
}

func (h *ExportHandler) ExportPostman(projectID int64) (string, error) {
	return h.svc.ExportPostman(projectID)
}
