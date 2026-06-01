package handlers

import "paw-api/services"

type ImporterHandler struct {
	service *services.ImporterService
}

func NewImporterHandler() *ImporterHandler {
	return &ImporterHandler{service: services.NewImporterService()}
}

type ImportCollection struct {
	Name string `json:"name"`
}

type ImportRequest struct {
	Name         string `json:"name"`
	CollectionID string `json:"collection_id"`
	Method       string `json:"method"`
	URL          string `json:"url"`
	Headers      string `json:"headers"`
	Params       string `json:"params"`
	Body         string `json:"body"`
	Auth         string `json:"auth"`
}

type ImportResponse struct {
	Collections []ImportCollection `json:"collections"`
	Requests    []ImportRequest    `json:"requests"`
}

func (h *ImporterHandler) ImportPostman(content string) (*ImportResponse, error) {
	result, err := h.service.ImportPostman(content)
	if err != nil {
		return nil, err
	}
	return toImportResponse(result), nil
}

func (h *ImporterHandler) ImportSwagger(content string) (*ImportResponse, error) {
	result, err := h.service.ImportSwagger(content)
	if err != nil {
		return nil, err
	}
	return toImportResponse(result), nil
}

func (h *ImporterHandler) ImportCurl(curlStr string) (*ImportRequest, error) {
	req, err := h.service.ImportCurl(curlStr)
	if err != nil {
		return nil, err
	}
	return &ImportRequest{
		Name:   req.Name,
		Method: req.Method,
		URL:    req.URL,
		Headers: req.Headers,
		Params:  req.Params,
		Body:    req.Body,
	}, nil
}

func toImportResponse(r *services.ImportResult) *ImportResponse {
	resp := &ImportResponse{}
	for _, c := range r.Collections {
		resp.Collections = append(resp.Collections, ImportCollection{
			Name: c.Name,
		})
	}
	for _, r := range r.Requests {
		resp.Requests = append(resp.Requests, ImportRequest{
			Name:         r.Name,
			CollectionID: r.CollectionID,
			Method:       r.Method,
			URL:          r.URL,
			Headers:      r.Headers,
			Params:       r.Params,
			Body:         r.Body,
			Auth:         r.Auth,
		})
	}
	return resp
}
