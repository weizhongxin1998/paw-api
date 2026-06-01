package services

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"paw-api/pkg/httpclient"
)

type ImporterService struct{}

func NewImporterService() *ImporterService {
	return &ImporterService{}
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

type ImportResult struct {
	Collections []ImportCollection
	Requests    []ImportRequest
}

func (s *ImporterService) ImportPostman(content string) (*ImportResult, error) {
	var pm struct {
		Info struct {
			Name string `json:"name"`
		} `json:"info"`
		Item []pmItem `json:"item"`
	}
	if err := json.Unmarshal([]byte(content), &pm); err != nil {
		return nil, fmt.Errorf("invalid Postman collection: %w", err)
	}
	result := &ImportResult{}
	colName := pm.Info.Name
	if colName == "" {
		colName = "Imported"
	}
	s.flattenPMItems(pm.Item, "", &result.Collections, &result.Requests, colName)
	return result, nil
}

type pmItem struct {
	Name    string    `json:"name"`
	Request *pmReq   `json:"request,omitempty"`
	Item    []pmItem `json:"item,omitempty"`
}

type pmReq struct {
	Method  string            `json:"method"`
	URL     pmURL             `json:"url"`
	Header  []pmHeader       `json:"header"`
	Body    *pmBody          `json:"body,omitempty"`
	Auth    *pmAuth          `json:"auth,omitempty"`
}

type pmURL struct {
	Raw   string   `json:"raw"`
	Path  []string `json:"path"`
	Query []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"query"`
}

type pmHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type pmBody struct {
	Mode     string          `json:"mode"`
	Raw      string          `json:"raw,omitempty"`
	URLEncoded []pmHeader   `json:"urlencoded,omitempty"`
	FormData  []pmFormField `json:"formdata,omitempty"`
}

type pmFormField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
	Src   string `json:"src,omitempty"`
}

type pmAuth struct {
	Type   string       `json:"type"`
	Bearer []pmAuthPair `json:"bearer,omitempty"`
	Basic  []pmAuthPair `json:"basic,omitempty"`
}

type pmAuthPair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

func (s *ImporterService) flattenPMItems(items []pmItem, parentID string, cols *[]ImportCollection, reqs *[]ImportRequest, colName string) string {
	colID := ""
	if parentID == "" {
		colID = "pm-import-" + colName
	} else {
		colID = "pm-folder-" + fmt.Sprintf("%p", &items)
	}
	for _, item := range items {
		if item.Request != nil {
			r := s.convertPMRequest(item, colID)
			if r != nil {
				*reqs = append(*reqs, *r)
			}
		}
		if item.Item != nil {
			s.flattenPMItems(item.Item, colID, cols, reqs, item.Name)
		}
	}
	return colID
}

func (s *ImporterService) convertPMRequest(item pmItem, colID string) *ImportRequest {
	req := item.Request
	if req == nil {
		return nil
	}

	rawURL := req.URL.Raw
	if rawURL == "" {
		rawURL = strings.Join(req.URL.Path, "/")
	}

	headers := pmHeadersToJSON(req.Header)
	params := "[]"
	if len(req.URL.Query) > 0 {
		var pairs []httpclient.KeyValuePair
		for _, q := range req.URL.Query {
			pairs = append(pairs, httpclient.KeyValuePair{Key: q.Key, Value: q.Value, Enabled: true})
		}
		b, _ := json.Marshal(pairs)
		params = string(b)
	}

	body := ""
	if req.Body != nil {
		switch req.Body.Mode {
		case "raw":
			body = req.Body.Raw
		case "urlencoded":
			vals := url.Values{}
			for _, p := range req.Body.URLEncoded {
				vals.Set(p.Key, p.Value)
			}
			body = vals.Encode()
		case "formdata":
			var pairs []httpclient.KeyValuePair
			for _, f := range req.Body.FormData {
				pairs = append(pairs, httpclient.KeyValuePair{Key: f.Key, Value: f.Value, Enabled: true})
			}
			b, _ := json.Marshal(pairs)
			body = string(b)
		}
	}

	auth := "{}"
	if req.Auth != nil {
		authMap := map[string]string{"type": req.Auth.Type}
		allPairs := append(req.Auth.Bearer, req.Auth.Basic...)
		for _, p := range allPairs {
			authMap[p.Key] = p.Value
		}
		b, _ := json.Marshal(authMap)
		auth = string(b)
	}

	return &ImportRequest{
		Name:         item.Name,
		CollectionID: colID,
		Method:       req.Method,
		URL:          rawURL,
		Headers:      headers,
		Params:       params,
		Body:         body,
		Auth:         auth,
	}
}

func pmHeadersToJSON(headers []pmHeader) string {
	if len(headers) == 0 {
		return "[]"
	}
	var pairs []httpclient.KeyValuePair
	for _, h := range headers {
		pairs = append(pairs, httpclient.KeyValuePair{Key: h.Key, Value: h.Value, Enabled: true})
	}
	b, _ := json.Marshal(pairs)
	return string(b)
}

func (s *ImporterService) ImportSwagger(content string) (*ImportResult, error) {
	var doc struct {
		Info struct {
			Title string `json:"title"`
		} `json:"info"`
		Paths map[string]map[string]struct {
			Summary     string `json:"summary"`
			OperationID string `json:"operationId"`
			Parameters  []struct {
				Name     string `json:"name"`
				In       string `json:"in"`
				Required bool   `json:"required"`
				Schema   struct {
					Type string `json:"type"`
				} `json:"schema"`
			} `json:"parameters"`
			RequestBody *struct {
				Content map[string]struct {
					Schema *struct {
						Type string `json:"type"`
					} `json:"schema"`
				} `json:"content"`
			} `json:"requestBody,omitempty"`
		} `json:"paths"`
	}
	if err := json.Unmarshal([]byte(content), &doc); err != nil {
		return nil, fmt.Errorf("invalid OpenAPI/Swagger JSON: %w", err)
	}

	result := &ImportResult{}
	apiName := doc.Info.Title
	if apiName == "" {
		apiName = "API Import"
	}

	for path, methods := range doc.Paths {
		for method, op := range methods {
			u, _ := url.Parse(path)
			headers := "[]"
			params := "[]"
			if op.Parameters != nil {
				var hPairs, pPairs []httpclient.KeyValuePair
				for _, p := range op.Parameters {
					switch p.In {
					case "header":
						hPairs = append(hPairs, httpclient.KeyValuePair{Key: p.Name, Value: "", Enabled: true})
					case "query":
						pPairs = append(pPairs, httpclient.KeyValuePair{Key: p.Name, Value: "", Enabled: true})
					}
				}
				if len(hPairs) > 0 {
					b, _ := json.Marshal(hPairs)
					headers = string(b)
				}
				if len(pPairs) > 0 {
					b, _ := json.Marshal(pPairs)
					params = string(b)
				}
			}

			name := op.Summary
			if name == "" {
				name = op.OperationID
			}
			if name == "" {
				name = method + " " + path
			}

			r := &ImportRequest{
				Name:         name,
				CollectionID: apiName,
				Method:       strings.ToUpper(method),
				URL:          u.String(),
				Headers:      headers,
				Params:       params,
			}
			result.Requests = append(result.Requests, *r)
		}
	}

	return result, nil
}

func (s *ImporterService) ImportCurl(curlStr string) (*ImportRequest, error) {
	curlStr = strings.TrimSpace(curlStr)
	curlStr = strings.TrimPrefix(curlStr, "curl ")
	curlStr = strings.TrimPrefix(curlStr, "curl ")

	ir := &ImportRequest{
		Method:  "GET",
		Headers: "[]",
		Params:  "[]",
		Body:    "",
	}

	args := tokenizeCurl(curlStr)
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "-X" || arg == "--request":
			if i+1 < len(args) {
				i++
				ir.Method = strings.ToUpper(args[i])
			}
		case arg == "-H" || arg == "--header":
			if i+1 < len(args) {
				i++
				parts := strings.SplitN(args[i], ":", 2)
				if len(parts) == 2 {
					var pairs []httpclient.KeyValuePair
					json.Unmarshal([]byte(ir.Headers), &pairs)
					pairs = append(pairs, httpclient.KeyValuePair{Key: strings.TrimSpace(parts[0]), Value: strings.TrimSpace(parts[1]), Enabled: true})
					b, _ := json.Marshal(pairs)
					ir.Headers = string(b)
				}
			}
		case arg == "-d" || arg == "--data" || arg == "--data-raw" || arg == "--data-binary":
			if i+1 < len(args) {
				i++
				ir.Body = args[i]
			}
		case !strings.HasPrefix(arg, "-") && ir.URL == "":
			ir.URL = arg
		}
	}

	if strings.HasPrefix(ir.Body, "{") || strings.HasPrefix(ir.Body, "[") {
		bodyStruct := map[string]interface{}{}
		if json.Unmarshal([]byte(ir.Body), &bodyStruct) == nil {
			// JSON body, keep as is
		}
	}

	return ir, nil
}

func tokenizeCurl(s string) []string {
	var tokens []string
	current := ""
	inQuote := false
	quoteChar := byte(0)
	escape := false

	for i := 0; i < len(s); i++ {
		c := s[i]
		if escape {
			current += string(c)
			escape = false
			continue
		}
		if c == '\\' && inQuote {
			escape = true
			continue
		}
		if c == '"' || c == '\'' {
			if !inQuote {
				inQuote = true
				quoteChar = c
			} else if c == quoteChar {
				inQuote = false
			} else {
				current += string(c)
			}
			continue
		}
		if c == ' ' && !inQuote {
			if current != "" {
				tokens = append(tokens, current)
				current = ""
			}
			continue
		}
		current += string(c)
	}
	if current != "" {
		tokens = append(tokens, current)
	}
	return tokens
}
