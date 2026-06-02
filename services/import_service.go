package services

import (
	"encoding/json"
	"os"

	"paw-api/models"
	"paw-api/pkg/snowflake"
	"paw-api/repositories"
)

type ImportService struct {
	collectionRepo *repositories.CollectionRepository
	requestRepo    *repositories.RequestRepository
	sf             *snowflake.Generator
}

type ImportResult struct {
	Collections int `json:"collections"`
	Requests    int `json:"requests"`
}

func NewImportService(
	collectionRepo *repositories.CollectionRepository,
	requestRepo *repositories.RequestRepository,
	sf *snowflake.Generator,
) *ImportService {
	return &ImportService{
		collectionRepo: collectionRepo,
		requestRepo:    requestRepo,
		sf:             sf,
	}
}

// ---------- Postman JSON structures ----------

type postmanCollection struct {
	Info postmanInfo   `json:"info"`
	Item []postmanItem `json:"item"`
}

type postmanInfo struct {
	Name string `json:"name"`
}

type postmanItem struct {
	Name    string           `json:"name"`
	Item    []postmanItem    `json:"item,omitempty"`
	Request *postmanRequest  `json:"request,omitempty"`
}

type postmanRequest struct {
	Method string       `json:"method"`
	URL    postmanURL   `json:"url"`
	Header []postmanKV  `json:"header,omitempty"`
	Body   *postmanBody `json:"body,omitempty"`
}

type postmanURL struct {
	Raw string `json:"raw"`
}

type postmanKV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type postmanBody struct {
	Mode       string      `json:"mode"`
	Raw        string      `json:"raw,omitempty"`
	Formdata   []postmanKV `json:"formdata,omitempty"`
	Urlencoded []postmanKV `json:"urlencoded,omitempty"`
}

// ---------- Public API ----------

func (s *ImportService) ImportPostman(projectID int64, filePath string) (*ImportResult, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var pc postmanCollection
	if err := json.Unmarshal(data, &pc); err != nil {
		return nil, err
	}

	result := &ImportResult{}

	root := &models.Collection{
		ProjectID: projectID,
		Name:      pc.Info.Name,
	}
	if err := s.collectionRepo.Create(root); err != nil {
		return nil, err
	}
	result.Collections++

	if err := s.processItems(projectID, root.ID, pc.Item, result); err != nil {
		return nil, err
	}

	return result, nil
}

// ---------- Internal ----------

func (s *ImportService) processItems(projectID int64, parentID int64, items []postmanItem, result *ImportResult) error {
	for _, item := range items {
		if len(item.Item) > 0 {
			coll := &models.Collection{
				ProjectID: projectID,
				ParentID:  &parentID,
				Name:      item.Name,
			}
			if err := s.collectionRepo.Create(coll); err != nil {
				return err
			}
			result.Collections++
			if err := s.processItems(projectID, coll.ID, item.Item, result); err != nil {
				return err
			}
		} else if item.Request != nil {
			headers := convertHeaders(item.Request.Header)
			bodyType, body := convertBody(item.Request.Body)

			req := &models.Request{
				CollectionID: parentID,
				Name:         item.Name,
				Method:       item.Request.Method,
				URL:          item.Request.URL.Raw,
				Headers:      headers,
				Params:       "[]",
				BodyType:     bodyType,
				Body:         body,
				Auth:         `{"type":"none"}`,
			}
			if err := s.requestRepo.Create(req); err != nil {
				return err
			}
			result.Requests++
		}
	}
	return nil
}

func convertHeaders(headers []postmanKV) string {
	out := make([]map[string]interface{}, 0, len(headers))
	for _, h := range headers {
		out = append(out, map[string]interface{}{
			"key":         h.Key,
			"value":       h.Value,
			"description": "",
			"enabled":     true,
		})
	}
	b, _ := json.Marshal(out)
	return string(b)
}

func convertBody(body *postmanBody) (string, string) {
	if body == nil {
		return "none", "{}"
	}

	switch body.Mode {
	case "raw":
		content, _ := json.Marshal(map[string]string{"content": body.Raw})
		return "raw", string(content)
	case "formdata":
		b, _ := json.Marshal(body.Formdata)
		return "formdata", string(b)
	case "urlencoded":
		b, _ := json.Marshal(body.Urlencoded)
		return "urlencoded", string(b)
	default:
		return "none", "{}"
	}
}
