package services

import (
	"encoding/json"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

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

// ---------- OpenAPI structures ----------

type openAPIDoc struct {
	OpenAPI  string                     `yaml:"openapi"`
	Swagger  string                     `yaml:"swagger"`
	Info     struct {
		Title   string `yaml:"title"`
		Version string `yaml:"version"`
	} `yaml:"info"`
	Servers  []struct {
		URL string `yaml:"url"`
	} `yaml:"servers"`
	Host     string                      `yaml:"host"`
	BasePath string                      `yaml:"basePath"`
	Schemes  []string                    `yaml:"schemes"`
	Paths    map[string]openAPIPathItem  `yaml:"paths"`
}

type openAPIPathItem struct {
	Parameters []openAPIParameter `yaml:"parameters"`
	Get        *openAPIOperation  `yaml:"get"`
	Post       *openAPIOperation  `yaml:"post"`
	Put        *openAPIOperation  `yaml:"put"`
	Patch      *openAPIOperation  `yaml:"patch"`
	Delete     *openAPIOperation  `yaml:"delete"`
	Options    *openAPIOperation  `yaml:"options"`
	Head       *openAPIOperation  `yaml:"head"`
}

type openAPIOperation struct {
	Tags        []string           `yaml:"tags"`
	Summary     string             `yaml:"summary"`
	Description string             `yaml:"description"`
	OperationID string             `yaml:"operationId"`
	Parameters  []openAPIParameter `yaml:"parameters"`
}

type openAPIParameter struct {
	Name        string `yaml:"name"`
	In          string `yaml:"in"`
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
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

func (s *ImportService) ImportOpenAPI(projectID int64, filePath string) (*ImportResult, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var doc openAPIDoc
	// Auto-detect format: try JSON first, fall back to YAML
	if err := json.Unmarshal(data, &doc); err != nil {
		if err := yaml.Unmarshal(data, &doc); err != nil {
			return nil, err
		}
	}

	return s.processOpenAPI(projectID, &doc)
}

// ---------- Internal ----------

func (s *ImportService) processOpenAPI(projectID int64, doc *openAPIDoc) (*ImportResult, error) {
	baseURL := resolveBaseURL(doc)
	result := &ImportResult{}
	collMap := make(map[string]int64)

	for path, pathItem := range doc.Paths {
		ops := map[string]*openAPIOperation{
			"GET":    pathItem.Get,
			"POST":   pathItem.Post,
			"PUT":    pathItem.Put,
			"PATCH":  pathItem.Patch,
			"DELETE": pathItem.Delete,
			"OPTIONS": pathItem.Options,
			"HEAD":   pathItem.Head,
		}

		for method, op := range ops {
			if op == nil {
				continue
			}

			tag := resolveTag(op, path)
			collID, exists := collMap[tag]
			if !exists {
				coll := &models.Collection{
					ProjectID: projectID,
					Name:      tag,
				}
				if err := s.collectionRepo.Create(coll); err != nil {
					return nil, err
				}
				collMap[tag] = coll.ID
				collID = coll.ID
				result.Collections++
			}

			fullURL := baseURL + path
			params := convertOpenAPIParams(pathItem.Parameters, op.Parameters)
			name := op.Summary
			if name == "" {
				name = op.OperationID
			}
			if name == "" {
				name = method + " " + path
			}

			req := &models.Request{
				CollectionID: collID,
				Name:         name,
				Description:  op.Description,
				Method:       method,
				URL:          fullURL,
				Headers:      "[]",
				Params:       params,
				BodyType:     "none",
				Body:         "{}",
				Auth:         `{"type":"none"}`,
			}
			if err := s.requestRepo.Create(req); err != nil {
				return nil, err
			}
			result.Requests++
		}
	}

	return result, nil
}

func resolveBaseURL(doc *openAPIDoc) string {
	if len(doc.Servers) > 0 {
		return strings.TrimRight(doc.Servers[0].URL, "/")
	}
	scheme := "https"
	if len(doc.Schemes) > 0 {
		scheme = doc.Schemes[0]
	}
	return strings.TrimRight(scheme+"://"+doc.Host+doc.BasePath, "/")
}

func resolveTag(op *openAPIOperation, path string) string {
	if len(op.Tags) > 0 {
		return op.Tags[0]
	}
	seg := strings.Trim(path, "/")
	if idx := strings.Index(seg, "/"); idx > 0 {
		seg = seg[:idx]
	}
	if seg == "" {
		return "Default"
	}
	return seg
}

func convertOpenAPIParams(pathParams, opParams []openAPIParameter) string {
	out := make([]map[string]interface{}, 0, len(pathParams)+len(opParams))
	for _, p := range append(pathParams, opParams...) {
		out = append(out, map[string]interface{}{
			"key":         p.Name,
			"value":       "",
			"description": p.Description,
			"enabled":     p.Required,
		})
	}
	if out == nil {
		return "[]"
	}
	b, _ := json.Marshal(out)
	return string(b)
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
