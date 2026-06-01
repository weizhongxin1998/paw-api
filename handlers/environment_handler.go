package handlers

import (
	"paw-api/models"
	"paw-api/services"
)

type EnvironmentHandler struct {
	service *services.EnvironmentService
}

func NewEnvironmentHandler() *EnvironmentHandler {
	return &EnvironmentHandler{service: services.NewEnvironmentService()}
}

func (h *EnvironmentHandler) CreateEnvironment(projectID, name, variables string, isActive bool) (*models.Environment, error) {
	return h.service.Create(projectID, name, variables, isActive)
}

func (h *EnvironmentHandler) GetEnvironment(id string) (*models.Environment, error) {
	return h.service.GetByID(id)
}

func (h *EnvironmentHandler) ListEnvironments(projectID string) ([]models.Environment, error) {
	return h.service.ListByProject(projectID)
}

func (h *EnvironmentHandler) UpdateEnvironment(id, name, variables string) (*models.Environment, error) {
	return h.service.Update(id, name, variables)
}

func (h *EnvironmentHandler) DeleteEnvironment(id string) error {
	return h.service.Delete(id)
}

func (h *EnvironmentHandler) SetActiveEnvironment(id, projectID string) (*models.Environment, error) {
	return h.service.SetActive(id, projectID)
}
