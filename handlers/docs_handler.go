package handlers

import (
	"paw-api/services"
)

type DocsHandler struct {
	svc *services.DocsService
}

func NewDocsHandler(svc *services.DocsService) *DocsHandler {
	return &DocsHandler{svc: svc}
}

func (h *DocsHandler) GenerateMarkdown(projectID int64) (string, error) {
	return h.svc.GenerateMarkdown(projectID)
}

func (h *DocsHandler) GenerateHTML(projectID int64) (string, error) {
	return h.svc.GenerateHTML(projectID)
}

func (h *DocsHandler) GenerateRequestMarkdown(requestID int64) (string, error) {
	return h.svc.GenerateRequestMarkdown(requestID)
}

func (h *DocsHandler) GenerateRequestHTML(requestID int64) (string, error) {
	return h.svc.GenerateRequestHTML(requestID)
}
