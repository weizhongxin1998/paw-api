package handlers

import "paw-api/services"

type SettingsHandler struct {
	svc *services.SettingsService
}

func NewSettingsHandler(svc *services.SettingsService) *SettingsHandler {
	return &SettingsHandler{svc: svc}
}

func (h *SettingsHandler) Get(key string) (string, error) {
	return h.svc.Get(key)
}

func (h *SettingsHandler) Set(key, value string) error {
	return h.svc.Set(key, value)
}

func (h *SettingsHandler) GetAll() (map[string]string, error) {
	return h.svc.GetAll()
}
