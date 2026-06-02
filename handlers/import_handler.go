package handlers

import (
	"paw-api/services"
)

type ImportHandler struct {
	svc *services.ImportService
}

func NewImportHandler(svc *services.ImportService) *ImportHandler {
	return &ImportHandler{svc: svc}
}

func (h *ImportHandler) ImportPostman(projectID int64, filePath string) (*services.ImportResult, error) {
	return h.svc.ImportPostman(projectID, filePath)
}
