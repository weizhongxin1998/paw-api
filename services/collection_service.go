package services

import (
	"paw-api/models"
	"paw-api/pkg/snowflake"
	"paw-api/repositories"
)

type CollectionService struct {
	collectionRepo *repositories.CollectionRepository
	requestRepo    *repositories.RequestRepository
	sf             *snowflake.Generator
}

func NewCollectionService(
	collectionRepo *repositories.CollectionRepository,
	requestRepo *repositories.RequestRepository,
	sf *snowflake.Generator,
) *CollectionService {
	return &CollectionService{
		collectionRepo: collectionRepo,
		requestRepo:    requestRepo,
		sf:             sf,
	}
}

func (s *CollectionService) GetTree(projectID int64) ([]models.TreeItem, error) {
	collections, err := s.collectionRepo.ListByProject(projectID)
	if err != nil {
		return nil, err
	}

	collectionsByParent := make(map[int64][]models.Collection)
	var roots []models.Collection
	for _, c := range collections {
		if c.ParentID == nil {
			roots = append(roots, c)
		} else {
			collectionsByParent[*c.ParentID] = append(collectionsByParent[*c.ParentID], c)
		}
	}

	var buildTree func(parent *models.Collection) models.TreeItem
	buildTree = func(parent *models.Collection) models.TreeItem {
		item := models.TreeItem{
			ID:        parent.ID,
			Name:      parent.Name,
			Type:      "folder",
			SortOrder: parent.SortOrder,
		}

		for _, child := range collectionsByParent[parent.ID] {
			item.Children = append(item.Children, buildTree(&child))
		}

		requests, err := s.requestRepo.ListByCollection(parent.ID)
		if err == nil {
			for _, req := range requests {
				item.Children = append(item.Children, models.TreeItem{
					ID:        req.ID,
					Name:      req.Name,
					Type:      "request",
					Method:    req.Method,
					URL:       req.URL,
					SortOrder: req.SortOrder,
				})
			}
		}

		return item
	}

	var tree []models.TreeItem
	for _, root := range roots {
		tree = append(tree, buildTree(&root))
	}

	return tree, nil
}

func (s *CollectionService) Get(id int64) (*models.Collection, error) {
	return s.collectionRepo.GetByID(id)
}

func (s *CollectionService) Create(projectID int64, parentID *int64, name string) (*models.Collection, error) {
	collection := &models.Collection{
		ProjectID: projectID,
		ParentID:  parentID,
		Name:      name,
	}
	if err := s.collectionRepo.Create(collection); err != nil {
		return nil, err
	}
	return collection, nil
}

func (s *CollectionService) Rename(id int64, name string) error {
	collection, err := s.collectionRepo.GetByID(id)
	if err != nil {
		return err
	}
	collection.Name = name
	return s.collectionRepo.Update(collection)
}

func (s *CollectionService) Move(id int64, parentID *int64, sortOrder int) error {
	if err := s.collectionRepo.MoveToParent(id, parentID); err != nil {
		return err
	}
	return s.collectionRepo.UpdateSortOrder(id, sortOrder)
}

func (s *CollectionService) Delete(id int64) error {
	return s.collectionRepo.Delete(id)
}
