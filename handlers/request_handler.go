package handlers

import (
	"paw-api/models"
	"paw-api/pkg/httpclient"
	"paw-api/services"
)

type RequestHandler struct {
	service *services.RequestService
}

func NewRequestHandler() *RequestHandler {
	return &RequestHandler{service: services.NewRequestService()}
}

func (h *RequestHandler) CreateRequest(collectionID, name, method, url, headers, params, body, auth, script string, sortOrder int) (*models.Request, error) {
	return h.service.Create(collectionID, name, method, url, headers, params, body, auth, script, sortOrder)
}

func (h *RequestHandler) GetRequest(id string) (*models.Request, error) {
	return h.service.GetByID(id)
}

func (h *RequestHandler) ListRequests(collectionID string) ([]models.Request, error) {
	return h.service.ListByCollection(collectionID)
}

func (h *RequestHandler) UpdateRequest(id, name, method, url, headers, params, body, auth, script string, sortOrder int) (*models.Request, error) {
	return h.service.Update(id, name, method, url, headers, params, body, auth, script, sortOrder)
}

type SendRequestInput struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    string
}

func (h *RequestHandler) SendRequest(input SendRequestInput) (*httpclient.Response, error) {
	return h.service.Send(input.Method, input.URL, input.Headers, input.Body)
}

func (h *RequestHandler) DeleteRequest(id string) error {
	return h.service.Delete(id)
}
