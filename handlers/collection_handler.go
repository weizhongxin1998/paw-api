package handlers

import (
	"paw-api/models"
	"paw-api/services"
)

type CollectionHandler struct {
	svc *services.CollectionService
}

func NewCollectionHandler(svc *services.CollectionService) *CollectionHandler {
	return &CollectionHandler{svc: svc}
}

func (h *CollectionHandler) GetTree(projectID int64) ([]models.TreeItem, error) {
	return h.svc.GetTree(projectID)
}

func (h *CollectionHandler) Get(id int64) (*models.Collection, error) {
	return h.svc.Get(id)
}

func (h *CollectionHandler) Create(projectID int64, parentID *int64, name string) (*models.Collection, error) {
	return h.svc.Create(projectID, parentID, name)
}

func (h *CollectionHandler) Rename(id int64, name string) error {
	return h.svc.Rename(id, name)
}

func (h *CollectionHandler) Move(id int64, parentID *int64, sortOrder int) error {
	return h.svc.Move(id, parentID, sortOrder)
}

func (h *CollectionHandler) Delete(id int64) error {
	return h.svc.Delete(id)
}
