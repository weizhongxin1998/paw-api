package handlers

import "paw-api/services"

type VariableHandler struct {
	svc *services.VariableService
}

func NewVariableHandler(svc *services.VariableService) *VariableHandler {
	return &VariableHandler{svc: svc}
}

func (h *VariableHandler) Resolve(text string, envID int64) (string, error) {
	return h.svc.Resolve(text, envID)
}

func (h *VariableHandler) ResolveMap(m map[string]string, envID int64) (map[string]string, error) {
	return h.svc.ResolveMap(m, envID)
}
