package handlers

import "paw-api/services"

type ImporterHandler struct {
	importerSvc    *services.ImporterService
	collectionSvc  *services.CollectionService
	requestSvc     *services.RequestService
}

func NewImporterHandler() *ImporterHandler {
	return &ImporterHandler{
		importerSvc:   services.NewImporterService(),
		collectionSvc: services.NewCollectionService(),
		requestSvc:    services.NewRequestService(),
	}
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

func (h *ImporterHandler) ImportPostman(projectID, content string) (*ImportResponse, error) {
	result, err := h.importerSvc.ImportPostman(content)
	if err != nil {
		return nil, err
	}
	return h.persistImport(projectID, result)
}

func (h *ImporterHandler) ImportSwagger(projectID, content string) (*ImportResponse, error) {
	result, err := h.importerSvc.ImportSwagger(content)
	if err != nil {
		return nil, err
	}
	return h.persistImport(projectID, result)
}

func (h *ImporterHandler) persistImport(projectID string, result *services.ImportResult) (*ImportResponse, error) {
	colMap := make(map[string]string)
	for _, c := range result.Collections {
		created, err := h.collectionSvc.Create(projectID, "", c.Name, 0)
		if err != nil {
			return nil, err
		}
		colMap[c.Name] = created.ID
	}

	for _, r := range result.Requests {
		colID := ""
		if mapped, ok := colMap[r.CollectionID]; ok {
			colID = mapped
		}
		_, err := h.requestSvc.Create(colID, r.Name, r.Method, r.URL, r.Headers, r.Params, r.Body, r.Auth, "", 0)
		if err != nil {
			return nil, err
		}
	}

	resp := &ImportResponse{}
	for _, c := range result.Collections {
		resp.Collections = append(resp.Collections, ImportCollection{Name: c.Name})
	}
	for _, r := range result.Requests {
		resp.Requests = append(resp.Requests, ImportRequest{
			Name: r.Name, CollectionID: r.CollectionID, Method: r.Method,
			URL: r.URL, Headers: r.Headers, Params: r.Params, Body: r.Body, Auth: r.Auth,
		})
	}
	return resp, nil
}

func (h *ImporterHandler) ImportCurl(curlStr string) (*ImportRequest, error) {
	req, err := h.importerSvc.ImportCurl(curlStr)
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
