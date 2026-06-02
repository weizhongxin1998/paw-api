package handlers

import (
	"paw-api/models"
	"paw-api/services"
)

type ProjectHandler struct {
	svc *services.ProjectService
}

func NewProjectHandler(svc *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{svc: svc}
}

func (h *ProjectHandler) List() ([]models.Project, error) {
	return h.svc.List()
}

func (h *ProjectHandler) Get(id int64) (*models.Project, error) {
	return h.svc.Get(id)
}

func (h *ProjectHandler) Create(name, description string) (*models.Project, error) {
	return h.svc.Create(name, description)
}

func (h *ProjectHandler) Update(id int64, name, description string) (*models.Project, error) {
	return h.svc.Update(id, name, description)
}

func (h *ProjectHandler) Delete(id int64) error {
	return h.svc.Delete(id)
}

func (h *ProjectHandler) GetStats(id int64) (models.ProjectStats, error) {
	return h.svc.GetStats(id)
}
