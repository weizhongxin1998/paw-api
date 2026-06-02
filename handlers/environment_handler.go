package handlers

import (
	"paw-api/models"
	"paw-api/services"
)

type EnvironmentHandler struct {
	svc *services.EnvironmentService
}

func NewEnvironmentHandler(svc *services.EnvironmentService) *EnvironmentHandler {
	return &EnvironmentHandler{svc: svc}
}

func (h *EnvironmentHandler) List(projectID int64) ([]models.Environment, error) {
	return h.svc.List(projectID)
}

func (h *EnvironmentHandler) Create(projectID int64, name string, cloneFromID *int64) (*models.Environment, error) {
	return h.svc.Create(projectID, name, cloneFromID)
}

func (h *EnvironmentHandler) Rename(id int64, name string) error {
	return h.svc.Rename(id, name)
}

func (h *EnvironmentHandler) Delete(id int64) error {
	return h.svc.Delete(id)
}

func (h *EnvironmentHandler) Activate(id int64) error {
	return h.svc.Activate(id)
}

func (h *EnvironmentHandler) GetActive(projectID int64) (*models.Environment, error) {
	return h.svc.GetActive(projectID)
}

func (h *EnvironmentHandler) ListVariables(envID int64) ([]models.EnvVariable, error) {
	return h.svc.ListVariables(envID)
}

func (h *EnvironmentHandler) SaveVariables(envID int64, variables []models.EnvVariable) error {
	return h.svc.SaveVariables(envID, variables)
}
