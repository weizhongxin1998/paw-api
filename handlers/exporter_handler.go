package handlers

import (
	"encoding/json"
	"paw-api/models"
	"paw-api/services"
)

type ExporterHandler struct {
	service *services.ExporterService
}

func NewExporterHandler() *ExporterHandler {
	return &ExporterHandler{service: services.NewExporterService()}
}

func (h *ExporterHandler) ExportPostman(collectionsJSON, requestsJSON, colName string) (string, error) {
	cols, reqs := parseExportInput(collectionsJSON, requestsJSON)
	return h.service.ExportPostman(cols, reqs, colName)
}

func (h *ExporterHandler) ExportSwagger(collectionsJSON, requestsJSON, apiName string) (string, error) {
	cols, reqs := parseExportInput(collectionsJSON, requestsJSON)
	return h.service.ExportSwagger(cols, reqs, apiName)
}

func parseExportInput(collectionsJSON, requestsJSON string) ([]models.Collection, []models.Request) {
	var cols []models.Collection
	var reqs []models.Request
	json.Unmarshal([]byte(collectionsJSON), &cols)
	json.Unmarshal([]byte(requestsJSON), &reqs)
	if cols == nil {
		cols = []models.Collection{}
	}
	if reqs == nil {
		reqs = []models.Request{}
	}
	return cols, reqs
}
