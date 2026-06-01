package services

import (
	"errors"
	"time"
	"paw-api/models"
	"paw-api/pkg/httpclient"
	"paw-api/repositories"

	"github.com/google/uuid"
)

type RequestService struct {
	repo *repositories.RequestRepo
}

func NewRequestService() *RequestService {
	return &RequestService{repo: &repositories.RequestRepo{}}
}

func (s *RequestService) Create(collectionID, name, method, url, headers, params, body, auth, script string, sortOrder int) (*models.Request, error) {
	if name == "" {
		return nil, errors.New("request name is required")
	}
	now := time.Now()
	r := &models.Request{
		ID:           uuid.New().String(),
		CollectionID: collectionID,
		Name:         name,
		Method:       method,
		URL:          url,
		Headers:      headers,
		Params:       params,
		Body:         body,
		Auth:         auth,
		Script:       script,
		SortOrder:    sortOrder,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	return r, s.repo.Create(r)
}

func (s *RequestService) GetByID(id string) (*models.Request, error) {
	return s.repo.GetByID(id)
}

func (s *RequestService) ListByCollection(collectionID string) ([]models.Request, error) {
	return s.repo.ListByCollection(collectionID)
}

func (s *RequestService) Update(id, name, method, url, headers, params, body, auth, script string, sortOrder int) (*models.Request, error) {
	r, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, errors.New("request not found")
	}
	r.Name = name
	r.Method = method
	r.URL = url
	r.Headers = headers
	r.Params = params
	r.Body = body
	r.Auth = auth
	r.Script = script
	r.SortOrder = sortOrder
	r.UpdatedAt = time.Now()
	return r, s.repo.Update(r)
}

func (s *RequestService) Send(method, url string, headers map[string]string, body string) (*httpclient.Response, error) {
	client := httpclient.NewClient()
	return client.Do(&httpclient.Request{
		Method:  method,
		URL:     url,
		Headers: headers,
		Body:    body,
	})
}

func (s *RequestService) Delete(id string) error {
	return s.repo.Delete(id)
}
