package handlers

import (
	"context"
	"sync"

	"paw-api/models"
	"paw-api/services"
)

type RequestHandler struct {
	svc     *services.RequestService
	cancels sync.Map // map[int64]context.CancelFunc
}

func NewRequestHandler(svc *services.RequestService) *RequestHandler {
	return &RequestHandler{svc: svc}
}

func (h *RequestHandler) Get(id int64) (*models.Request, error) {
	return h.svc.Get(id)
}

func (h *RequestHandler) Create(collectionID int64, name, method string) (*models.Request, error) {
	return h.svc.Create(collectionID, name, method)
}

func (h *RequestHandler) Update(req *models.Request) error {
	return h.svc.Update(req)
}

func (h *RequestHandler) Clone(id int64) (*models.Request, error) {
	return h.svc.Clone(id)
}

func (h *RequestHandler) Move(id int64, collectionID int64, sortOrder int) error {
	return h.svc.Move(id, collectionID, sortOrder)
}

func (h *RequestHandler) Delete(id int64) error {
	return h.svc.Delete(id)
}

func (h *RequestHandler) SendRequest(sessionID int64, req *models.Request, envID int64) (*models.HTTPResponse, error) {
	ctx, cancel := context.WithCancel(context.Background())
	h.cancels.Store(sessionID, cancel)
	defer func() {
		h.cancels.Delete(sessionID)
		cancel()
	}()

	return h.svc.Send(ctx, req, envID)
}

func (h *RequestHandler) SendQuick(sessionID int64, method, url, headersJSON, body string, envID int64) (*models.HTTPResponse, error) {
	ctx, cancel := context.WithCancel(context.Background())
	h.cancels.Store(sessionID, cancel)
	defer func() {
		h.cancels.Delete(sessionID)
		cancel()
	}()

	return h.svc.SendQuick(ctx, method, url, headersJSON, body, envID)
}

func (h *RequestHandler) CancelRequest(sessionID int64) {
	if cancel, ok := h.cancels.Load(sessionID); ok {
		cancel.(context.CancelFunc)()
		h.cancels.Delete(sessionID)
	}
}
