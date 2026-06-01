package handlers

import (
	"paw-api/models"
	"paw-api/services"
)

type CollectionHandler struct {
	service *services.CollectionService
}

func NewCollectionHandler() *CollectionHandler {
	return &CollectionHandler{service: services.NewCollectionService()}
}

func (h *CollectionHandler) CreateCollection(projectID, parentID, name string, sortOrder int) (*models.Collection, error) {
	return h.service.Create(projectID, parentID, name, sortOrder)
}

func (h *CollectionHandler) GetCollection(id string) (*models.Collection, error) {
	return h.service.GetByID(id)
}

func (h *CollectionHandler) ListCollections(projectID string) ([]models.Collection, error) {
	return h.service.ListByProject(projectID)
}

func (h *CollectionHandler) UpdateCollection(id, name string, parentID *string, sortOrder int) (*models.Collection, error) {
	return h.service.Update(id, name, parentID, sortOrder)
}

func (h *CollectionHandler) DeleteCollection(id string) error {
	return h.service.Delete(id)
}
