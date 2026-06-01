package handlers

import (
	"paw-api/models"
	"paw-api/services"
)

type ProjectHandler struct {
	service *services.ProjectService
}

func NewProjectHandler() *ProjectHandler {
	return &ProjectHandler{service: services.NewProjectService()}
}

func (h *ProjectHandler) CreateProject(name, description string) (*models.Project, error) {
	return h.service.Create(name, description)
}

func (h *ProjectHandler) GetProject(id string) (*models.Project, error) {
	return h.service.GetByID(id)
}

func (h *ProjectHandler) ListProjects() ([]models.Project, error) {
	return h.service.List()
}

func (h *ProjectHandler) UpdateProject(id, name, description string) (*models.Project, error) {
	return h.service.Update(id, name, description)
}

func (h *ProjectHandler) DeleteProject(id string) error {
	return h.service.Delete(id)
}
