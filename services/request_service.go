package services

import (
	"context"
	"encoding/json"

	"paw-api/models"
	"paw-api/pkg/httpclient"
	"paw-api/pkg/snowflake"
	"paw-api/repositories"
)

type RequestService struct {
	requestRepo *repositories.RequestRepository
	historyRepo *repositories.HistoryRepository
	sf          *snowflake.Generator
	httpClient  *httpclient.Client
}

func NewRequestService(
	requestRepo *repositories.RequestRepository,
	historyRepo *repositories.HistoryRepository,
	sf *snowflake.Generator,
	httpClient *httpclient.Client,
) *RequestService {
	return &RequestService{
		requestRepo: requestRepo,
		historyRepo: historyRepo,
		sf:          sf,
		httpClient:  httpClient,
	}
}

func (s *RequestService) Get(id int64) (*models.Request, error) {
	return s.requestRepo.GetByID(id)
}

func (s *RequestService) Create(collectionID int64, name, method string) (*models.Request, error) {
	req := &models.Request{
		CollectionID: collectionID,
		Name:         name,
		Method:       method,
		Headers:      "[]",
		Params:       "[]",
		BodyType:     "none",
		Body:         "{}",
		Auth:         `{"type":"none"}`,
	}
	if err := s.requestRepo.Create(req); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *RequestService) Update(req *models.Request) error {
	return s.requestRepo.Update(req)
}

func (s *RequestService) Clone(id int64) (*models.Request, error) {
	original, err := s.requestRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	original.ID = 0
	original.Name = original.Name + " (Copy)"
	if err := s.requestRepo.Create(original); err != nil {
		return nil, err
	}
	return original, nil
}

func (s *RequestService) Move(id int64, collectionID int64, sortOrder int) error {
	if err := s.requestRepo.MoveToCollection(id, collectionID); err != nil {
		return err
	}
	return s.requestRepo.UpdateSortOrder(id, sortOrder)
}

func (s *RequestService) Delete(id int64) error {
	return s.requestRepo.Delete(id)
}

func (s *RequestService) Send(ctx context.Context, req *models.Request, envID int64) (*models.HTTPResponse, error) {
	headers := parseKvJSON(req.Headers)

	resp, err := s.httpClient.Execute(ctx, req.Method, req.URL, headers, []byte(req.Body))
	if err != nil {
		return nil, err
	}

	result := &models.HTTPResponse{
		Status:     resp.Status,
		StatusText: resp.StatusText,
		Time:       resp.TimeMs,
		Size:       resp.Size,
		Headers:    resp.Headers,
		Body:       string(resp.Body),
	}

	// 自动保存历史记录
	history := &models.History{
		ProjectID:       0,
		Method:          req.Method,
		URL:             req.URL,
		RequestHeaders:  req.Headers,
		RequestBody:     extractBodyContent(req.Body, req.BodyType),
		ResponseStatus:  resp.Status,
		ResponseHeaders: toJSON(resp.Headers),
		ResponseBody:    string(resp.Body),
		DurationMs:      int(resp.TimeMs),
	}
	if req.ID != 0 {
		history.RequestID = &req.ID
	}
	_ = s.historyRepo.Create(history)

	return result, nil
}

func (s *RequestService) SendQuick(ctx context.Context, method, url, headersJSON, body string, envID int64) (*models.HTTPResponse, error) {
	headers := parseKvJSON(headersJSON)

	resp, err := s.httpClient.Execute(ctx, method, url, headers, []byte(body))
	if err != nil {
		return nil, err
	}

	result := &models.HTTPResponse{
		Status:     resp.Status,
		StatusText: resp.StatusText,
		Time:       resp.TimeMs,
		Size:       resp.Size,
		Headers:    resp.Headers,
		Body:       string(resp.Body),
	}

	return result, nil
}

func parseKvJSON(raw string) map[string]string {
	var kvs []struct {
		Key     string `json:"key"`
		Value   string `json:"value"`
		Enabled bool   `json:"enabled"`
	}
	if err := json.Unmarshal([]byte(raw), &kvs); err != nil {
		return nil
	}

	result := make(map[string]string)
	for _, kv := range kvs {
		if kv.Enabled {
			result[kv.Key] = kv.Value
		}
	}
	return result
}

func extractBodyContent(bodyJSON, bodyType string) string {
	if bodyType == "raw" {
		var raw struct {
			Content string `json:"content"`
		}
		if json.Unmarshal([]byte(bodyJSON), &raw) == nil {
			return raw.Content
		}
	}
	return bodyJSON
}

func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
