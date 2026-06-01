package services

import (
	"errors"
	"time"
	"paw-api/models"
	"paw-api/repositories"

	"github.com/google/uuid"
)

type CollectionService struct {
	repo *repositories.CollectionRepo
}

func NewCollectionService() *CollectionService {
	return &CollectionService{repo: &repositories.CollectionRepo{}}
}

func (s *CollectionService) Create(projectID, parentID, name string, sortOrder int) (*models.Collection, error) {
	if name == "" {
		return nil, errors.New("collection name is required")
	}
	if sortOrder <= 0 {
		maxOrder, err := s.repo.GetMaxSortOrder(projectID, parentID)
		if err != nil {
			return nil, err
		}
		sortOrder = maxOrder + 1
	}
	now := time.Now()
	var pid *string
	if parentID != "" {
		pid = &parentID
	}
	c := &models.Collection{
		ID:        uuid.New().String(),
		ProjectID: projectID,
		ParentID:  pid,
		Name:      name,
		SortOrder: sortOrder,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return c, s.repo.Create(c)
}

func (s *CollectionService) GetByID(id string) (*models.Collection, error) {
	return s.repo.GetByID(id)
}

func (s *CollectionService) ListByProject(projectID string) ([]models.Collection, error) {
	return s.repo.ListByProject(projectID)
}

func (s *CollectionService) Update(id, name string, parentID *string, sortOrder int) (*models.Collection, error) {
	c, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, errors.New("collection not found")
	}
	if name != "" {
		c.Name = name
	}
	if parentID != nil {
		c.ParentID = parentID
	}
	c.SortOrder = sortOrder
	c.UpdatedAt = time.Now()
	return c, s.repo.Update(c)
}

func (s *CollectionService) Delete(id string) error {
	return s.repo.Delete(id)
}
