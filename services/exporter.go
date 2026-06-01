package services

import (
	"encoding/json"
	"sort"
	"paw-api/models"
)

type ExporterService struct {
	requestService *RequestService
	collectionSvc  *CollectionService
}

func NewExporterService() *ExporterService {
	return &ExporterService{
		requestService: NewRequestService(),
		collectionSvc:  NewCollectionService(),
	}
}

type pmExportCollection struct {
	Info struct {
		Name   string `json:"name"`
		Schema string `json:"schema"`
	} `json:"info"`
	Item []pmExportItem `json:"item"`
}

type pmExportItem struct {
	Name    string          `json:"name"`
	Request *pmExportReq    `json:"request,omitempty"`
	Item    []pmExportItem  `json:"item,omitempty"`
}

type pmExportReq struct {
	Method  string            `json:"method"`
	Header  []pmExportHeader `json:"header"`
	URL     pmExportURL       `json:"url"`
	Body    *pmExportBody     `json:"body,omitempty"`
}

type pmExportHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type pmExportURL struct {
	Raw   string        `json:"raw"`
}

type pmExportBody struct {
	Mode string `json:"mode"`
	Raw  string `json:"raw,omitempty"`
}

func (s *ExporterService) ExportPostman(collections []models.Collection, requests []models.Request, colName string) (string, error) {
	rootItems := s.buildPMTree(collections, requests, "")
	col := pmExportCollection{}
	col.Info.Name = colName
	if colName == "" {
		col.Info.Name = "Exported"
	}
	col.Info.Schema = "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
	col.Item = rootItems

	b, err := json.MarshalIndent(col, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (s *ExporterService) buildPMTree(collections []models.Collection, requests []models.Request, parentID string) []pmExportItem {
	var items []pmExportItem
	childCols := filterCollections(collections, parentID)
	sort.Slice(childCols, func(i, j int) bool { return childCols[i].SortOrder < childCols[j].SortOrder })

	for _, col := range childCols {
		childItems := s.buildPMTree(collections, requests, col.ID)
		colReqs := filterRequestsByCollection(requests, col.ID)
		for _, req := range colReqs {
			childItems = append(childItems, s.requestToPMItem(req))
		}
		if len(childItems) > 0 {
			items = append(items, pmExportItem{
				Name: col.Name,
				Item: childItems,
			})
		}
	}

	// Root level requests (no parent collection)
	if parentID == "" {
		rootReqs := filterRequestsByCollection(requests, "")
		for _, req := range rootReqs {
			items = append(items, s.requestToPMItem(req))
		}
	}

	return items
}

func (s *ExporterService) requestToPMItem(req models.Request) pmExportItem {
	item := pmExportItem{
		Name: req.Name,
		Request: &pmExportReq{
			Method: req.Method,
			URL:    pmExportURL{Raw: req.URL},
		},
	}

	var headers []pmExportHeader
	h := parseKVJSON(req.Headers)
	for _, pair := range h {
		headers = append(headers, pmExportHeader{Key: pair.Key, Value: pair.Value, Type: "text"})
	}
	if len(headers) > 0 {
		item.Request.Header = headers
	}

	if req.Body != "" {
		item.Request.Body = &pmExportBody{
			Mode: "raw",
			Raw:  req.Body,
		}
	}

	return item
}

func filterCollections(cols []models.Collection, parentID string) []models.Collection {
	var result []models.Collection
	for _, c := range cols {
		pID := ""
		if c.ParentID != nil {
			pID = *c.ParentID
		}
		if pID == parentID {
			result = append(result, c)
		}
	}
	return result
}

func filterRequestsByCollection(requests []models.Request, colID string) []models.Request {
	var result []models.Request
	for _, r := range requests {
		if r.CollectionID == colID {
			result = append(result, r)
		}
	}
	return result
}

func parseKVJSON(s string) []struct {
	Key   string `json:"key"`
	Value string `json:"value"`
} {
	var result []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if s == "" {
		return result
	}
	json.Unmarshal([]byte(s), &result)
	return result
}

func (s *ExporterService) ExportSwagger(collections []models.Collection, requests []models.Request, apiName string) (string, error) {
	swagger := map[string]interface{}{
		"openapi": "3.0.3",
		"info": map[string]interface{}{
			"title":   apiName,
			"version": "1.0.0",
		},
		"paths": map[string]interface{}{},
	}

	paths := swagger["paths"].(map[string]interface{})
	for _, req := range requests {
		path := req.URL
		if path == "" {
			path = "/"
		}
		method := stringsToLower(req.Method)
		if _, ok := paths[path]; !ok {
			paths[path] = map[string]interface{}{}
		}
		pathObj := paths[path].(map[string]interface{})
		pathObj[method] = map[string]interface{}{
			"summary": req.Name,
			"operationId": req.Name,
			"responses": map[string]interface{}{
				"200": map[string]interface{}{
					"description": "Successful response",
				},
			},
		}
	}

	b, err := json.MarshalIndent(swagger, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func stringsToLower(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			b[i] = s[i] + 32
		} else {
			b[i] = s[i]
		}
	}
	return string(b)
}
