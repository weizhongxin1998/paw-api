package services

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type CodegenRequest struct {
	Method   string
	URL      string
	Headers  map[string]string
	Body     string
	BodyType string
}

type CodegenService struct{}

func NewCodegenService() *CodegenService {
	return &CodegenService{}
}

func (s *CodegenService) Generate(req CodegenRequest, lang string) (string, error) {
	switch lang {
	case "curl":
		return s.generateCurl(req), nil
	case "javascript":
		return s.generateJavaScript(req), nil
	case "python":
		return s.generatePython(req), nil
	case "go":
		return s.generateGo(req), nil
	default:
		return "", fmt.Errorf("unsupported language: %s", lang)
	}
}

func (s *CodegenService) generateCurl(req CodegenRequest) string {
	var b strings.Builder
	b.WriteString("curl -X " + req.Method + " '" + req.URL + "'")

	var keys []string
	for k := range req.Headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		b.WriteString(fmt.Sprintf(" \\\n  -H '%s: %s'", k, req.Headers[k]))
	}

	if req.Body != "" && req.BodyType != "none" {
		ct := contentTypeForBodyType(req.BodyType)
		if ct != "" {
			hasCT := false
			for k := range req.Headers {
				if strings.ToLower(k) == "content-type" {
					hasCT = true
					break
				}
			}
			if !hasCT {
				b.WriteString(fmt.Sprintf(" \\\n  -H 'Content-Type: %s'", ct))
			}
		}
		b.WriteString(fmt.Sprintf(" \\\n  -d '%s'", req.Body))
	}

	return b.String()
}

func contentTypeForBodyType(bodyType string) string {
	switch bodyType {
	case "json":
		return "application/json"
	case "text":
		return "text/plain"
	case "urlencoded":
		return "application/x-www-form-urlencoded"
	case "form-data":
		return "multipart/form-data"
	case "binary":
		return "application/octet-stream"
	default:
		return ""
	}
}

func (s *CodegenService) generateJavaScript(req CodegenRequest) string {
	var b strings.Builder
	b.WriteString("fetch('" + req.URL + "', {\n")
	b.WriteString("  method: '" + req.Method + "',\n")

	headers := make(map[string]string)
	for k, v := range req.Headers {
		headers[k] = v
	}
	if req.Body != "" && req.BodyType != "none" {
		ct := contentTypeForBodyType(req.BodyType)
		if ct != "" {
			hasCT := false
			for k := range headers {
				if strings.ToLower(k) == "content-type" {
					hasCT = true
					break
				}
			}
			if !hasCT {
				headers["Content-Type"] = ct
			}
		}
	}
	if len(headers) > 0 {
		headersJSON, _ := json.MarshalIndent(headers, "  ", "  ")
		b.WriteString("  headers: " + string(headersJSON) + ",\n")
	}

	if req.Body != "" && req.BodyType != "none" {
		switch req.BodyType {
		case "urlencoded":
			b.WriteString("  body: new URLSearchParams(" + jsonString(req.Body) + ").toString(),\n")
		default:
			b.WriteString("  body: " + jsonString(req.Body) + ",\n")
		}
	}

	b.WriteString("})")
	if req.BodyType == "json" {
		b.WriteString("\n  .then(res => res.json())")
	} else {
		b.WriteString("\n  .then(res => res.text())")
	}
	b.WriteString("\n  .then(data => console.log(data))")
	b.WriteString("\n  .catch(err => console.error(err))")

	return b.String()
}

func (s *CodegenService) generatePython(req CodegenRequest) string {
	var b strings.Builder

	imports := []string{"import requests"}
	if req.BodyType == "json" {
		imports = append(imports, "import json")
	}
	b.WriteString(strings.Join(imports, "\n") + "\n\n")

	b.WriteString(fmt.Sprintf("url = '%s'\n", req.URL))
	b.WriteString(fmt.Sprintf("method = '%s'\n", strings.ToLower(req.Method)))

	if len(req.Headers) > 0 {
		b.WriteString("\nheaders = {\n")
		var keys []string
		for k := range req.Headers {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			b.WriteString(fmt.Sprintf("    '%s': '%s',\n", k, req.Headers[k]))
		}
		b.WriteString("}\n")
	}

	switch req.BodyType {
	case "json":
		b.WriteString("\ndata = " + req.Body + "\n")
		b.WriteString(fmt.Sprintf("\nresponse = requests.request(method, url, headers=headers, json=data)\n"))
	case "urlencoded", "form-data":
		b.WriteString("\ndata = " + req.Body + "\n")
		b.WriteString(fmt.Sprintf("\nresponse = requests.request(method, url, headers=headers, data=data)\n"))
	default:
		if req.Body != "" {
			b.WriteString(fmt.Sprintf("\ndata = '''%s'''\n", req.Body))
			b.WriteString(fmt.Sprintf("\nresponse = requests.request(method, url, headers=headers, data=data)\n"))
		} else {
			b.WriteString(fmt.Sprintf("\nresponse = requests.request(method, url, headers=headers)\n"))
		}
	}

	b.WriteString("print(response.status_code)")
	b.WriteString("\nprint(response.text)")

	return b.String()
}

func (s *CodegenService) generateGo(req CodegenRequest) string {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("  \"bytes\"\n")
	b.WriteString("  \"fmt\"\n")
	b.WriteString("  \"io\"\n")
	b.WriteString("  \"net/http\"\n")
	b.WriteString("  \"net/url\"\n")
	if req.BodyType == "json" {
		b.WriteString("  \"encoding/json\"\n")
	}
	if req.BodyType == "form-data" {
		b.WriteString("  \"mime/multipart\"\n")
	}
	b.WriteString(")\n\n")
	b.WriteString("func main() {\n")

	bodyVar := "nil"
	if req.Body != "" && req.BodyType != "none" {
		switch req.BodyType {
		case "json":
			b.WriteString(fmt.Sprintf("  body := %s\n", req.Body))
			b.WriteString("  bodyBytes, _ := json.Marshal(body)\n")
			bodyVar = "bytes.NewReader(bodyBytes)"
		case "urlencoded":
			b.WriteString("  data := url.Values{}\n")
			b.WriteString(fmt.Sprintf("  // Parse body as key=value pairs: %s\n", req.Body))
			bodyVar = "strings.NewReader(data.Encode())"
		default:
			b.WriteString(fmt.Sprintf("  body := bytes.NewReader([]byte(`%s`))\n", req.Body))
			bodyVar = "body"
		}
		if bodyVar != "nil" {
			b.WriteString(fmt.Sprintf("  req, _ := http.NewRequest(\"%s\", \"%s\", %s)\n", req.Method, req.URL, bodyVar))
		}
	} else {
		b.WriteString(fmt.Sprintf("  req, _ := http.NewRequest(\"%s\", \"%s\", nil)\n", req.Method, req.URL))
	}

	var keys []string
	for k := range req.Headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		b.WriteString(fmt.Sprintf("  req.Header.Set(\"%s\", \"%s\")\n", k, req.Headers[k]))
	}

	if req.Body != "" && req.BodyType != "none" {
		ct := contentTypeForBodyType(req.BodyType)
		if ct != "" {
			hasCT := false
			for k := range req.Headers {
				if strings.ToLower(k) == "content-type" {
					hasCT = true
					break
				}
			}
			if !hasCT {
				b.WriteString(fmt.Sprintf("  req.Header.Set(\"Content-Type\", \"%s\")\n", ct))
			}
		}
	}

	b.WriteString("\n  client := &http.Client{}\n")
	b.WriteString("  resp, err := client.Do(req)\n")
	b.WriteString("  if err != nil {\n")
	b.WriteString("    fmt.Println(\"Error:\", err)\n")
	b.WriteString("    return\n")
	b.WriteString("  }\n")
	b.WriteString("  defer resp.Body.Close()\n\n")
	b.WriteString("  bodyBytes, _ := io.ReadAll(resp.Body)\n")
	b.WriteString("  fmt.Println(string(bodyBytes))\n")
	b.WriteString("}\n")

	return b.String()
}

func jsonString(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}
