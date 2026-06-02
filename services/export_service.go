package services

import (
	"encoding/json"
	"net/url"
	"sort"
	"strings"

	"paw-api/models"
	"paw-api/pkg/snowflake"
	"paw-api/repositories"
)

type ExportService struct {
	projectRepo    *repositories.ProjectRepository
	collectionRepo *repositories.CollectionRepository
	requestRepo    *repositories.RequestRepository
	sf             *snowflake.Generator
}

func NewExportService(
	projectRepo *repositories.ProjectRepository,
	collectionRepo *repositories.CollectionRepository,
	requestRepo *repositories.RequestRepository,
	sf *snowflake.Generator,
) *ExportService {
	return &ExportService{
		projectRepo:    projectRepo,
		collectionRepo: collectionRepo,
		requestRepo:    requestRepo,
		sf:             sf,
	}
}

func (s *ExportService) ExportPostman(projectID int64) (string, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return "", err
	}

	collections, err := s.collectionRepo.ListByProject(projectID)
	if err != nil {
		return "", err
	}

	collMap := make(map[int64]*models.Collection)
	childrenMap := make(map[int64][]*models.Collection)
	var roots []*models.Collection

	for i := range collections {
		c := &collections[i]
		collMap[c.ID] = c
		if c.ParentID == nil {
			roots = append(roots, c)
		} else {
			childrenMap[*c.ParentID] = append(childrenMap[*c.ParentID], c)
		}
	}

	sort.Slice(roots, func(i, j int) bool { return roots[i].SortOrder < roots[j].SortOrder })
	for _, children := range childrenMap {
		sort.Slice(children, func(i, j int) bool { return children[i].SortOrder < children[j].SortOrder })
	}

	items := make([]postmanItem, 0, len(roots))
	for _, root := range roots {
		items = append(items, s.buildCollectionItem(root, childrenMap))
	}

	pm := postmanCollection{
		Info: postmanInfo{
			Name:   project.Name,
			Schema: "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		Item: items,
	}

	jsonBytes, err := json.MarshalIndent(pm, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func (s *ExportService) buildCollectionItem(coll *models.Collection, childrenMap map[int64][]*models.Collection) postmanItem {
	item := postmanItem{
		Name: coll.Name,
		Item: make([]postmanItem, 0),
	}

	requests, err := s.requestRepo.ListByCollection(coll.ID)
	if err == nil {
		sort.Slice(requests, func(i, j int) bool { return requests[i].SortOrder < requests[j].SortOrder })
		for i := range requests {
			item.Item = append(item.Item, s.buildRequestItem(&requests[i]))
		}
	}

	children := childrenMap[coll.ID]
	sort.Slice(children, func(i, j int) bool { return children[i].SortOrder < children[j].SortOrder })
	for _, child := range children {
		item.Item = append(item.Item, s.buildCollectionItem(child, childrenMap))
	}

	return item
}

func (s *ExportService) buildRequestItem(req *models.Request) postmanItem {
	return postmanItem{
		Name: req.Name,
		Request: &postmanRequest{
			Method: req.Method,
			URL:    buildPostmanURL(req.URL, req.Params),
			Header: exportHeaders(req.Headers),
			Body:   exportBody(req.BodyType, req.Body),
		},
	}
}

func buildPostmanURL(rawURL string, paramsJSON string) postmanURL {
	pu := postmanURL{Raw: rawURL}

	u, err := url.Parse(rawURL)
	if err != nil || u.Host == "" {
		return pu
	}

	pu.Protocol = u.Scheme
	pu.Port = u.Port()

	hostParts := strings.Split(u.Hostname(), ".")
	pu.Host = make([]string, 0, len(hostParts))
	for _, p := range hostParts {
		if p != "" {
			pu.Host = append(pu.Host, p)
		}
	}

	pathStr := strings.Trim(u.Path, "/")
	if pathStr != "" {
		pu.Path = strings.Split(pathStr, "/")
	}

	queryParams := parseQueryParams(paramsJSON, u.Query())
	pu.Query = queryParams

	return pu
}

func parseQueryParams(paramsJSON string, urlQuery url.Values) []postmanKV {
	var result []postmanKV

	if paramsJSON != "" && paramsJSON != "[]" {
		var params []struct {
			Key     string `json:"key"`
			Value   string `json:"value"`
			Enabled bool   `json:"enabled"`
		}
		if err := json.Unmarshal([]byte(paramsJSON), &params); err == nil {
			for _, p := range params {
				if strings.TrimSpace(p.Key) != "" {
					result = append(result, postmanKV{Key: p.Key, Value: p.Value})
				}
			}
		}
	}

	for key, values := range urlQuery {
		for _, v := range values {
			if !containsQueryKey(result, key) {
				result = append(result, postmanKV{Key: key, Value: v})
			}
		}
	}

	return result
}

func containsQueryKey(params []postmanKV, key string) bool {
	for _, p := range params {
		if p.Key == key {
			return true
		}
	}
	return false
}

func exportHeaders(headersJSON string) []postmanKV {
	if headersJSON == "" || headersJSON == "[]" {
		return nil
	}

	var headers []struct {
		Key     string `json:"key"`
		Value   string `json:"value"`
		Enabled bool   `json:"enabled"`
	}
	if err := json.Unmarshal([]byte(headersJSON), &headers); err != nil {
		return nil
	}

	result := make([]postmanKV, 0, len(headers))
	for _, h := range headers {
		result = append(result, postmanKV{Key: h.Key, Value: h.Value})
	}
	return result
}

func exportBody(bodyType string, bodyJSON string) *postmanBody {
	if bodyType == "" || bodyType == "none" {
		return nil
	}

	body := &postmanBody{Mode: bodyType}

	switch bodyType {
	case "raw":
		if bodyJSON != "" && bodyJSON != "{}" {
			var rawBody struct {
				Content string `json:"content"`
			}
			if err := json.Unmarshal([]byte(bodyJSON), &rawBody); err == nil {
				body.Raw = rawBody.Content
			} else {
				body.Raw = bodyJSON
			}
		}
	case "formdata":
		if bodyJSON != "" && bodyJSON != "[]" {
			body.Formdata = unmarshalKV(bodyJSON)
		}
	case "urlencoded":
		if bodyJSON != "" && bodyJSON != "[]" {
			body.Urlencoded = unmarshalKV(bodyJSON)
		}
	}

	return body
}

func unmarshalKV(jsonStr string) []postmanKV {
	var entries []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := json.Unmarshal([]byte(jsonStr), &entries); err != nil {
		return nil
	}
	result := make([]postmanKV, 0, len(entries))
	for _, e := range entries {
		result = append(result, postmanKV{Key: e.Key, Value: e.Value})
	}
	return result
}
