package handlers

import (
	"paw-api/models"
	"paw-api/pkg/httpclient"
	"paw-api/services"
)

type RequestHandler struct {
	service      *services.RequestService
	assertSvc    *services.AssertService
	codegenSvc   *services.CodegenService
}

func NewRequestHandler() *RequestHandler {
	return &RequestHandler{
		service:    services.NewRequestService(),
		assertSvc:  services.NewAssertService(),
		codegenSvc: services.NewCodegenService(),
	}
}

func (h *RequestHandler) CreateRequest(collectionID, name, method, url, headers, params, body, auth, script string, sortOrder int) (*models.Request, error) {
	return h.service.Create(collectionID, name, method, url, headers, params, body, auth, script, sortOrder)
}

func (h *RequestHandler) GetRequest(id string) (*models.Request, error) {
	return h.service.GetByID(id)
}

func (h *RequestHandler) ListRequests(collectionID string) ([]models.Request, error) {
	return h.service.ListByCollection(collectionID)
}

func (h *RequestHandler) ListRequestsByProject(projectID string) ([]models.Request, error) {
	return h.service.ListByProjectID(projectID)
}

func (h *RequestHandler) UpdateRequest(id, collectionID, name, method, url, headers, params, body, auth, script string, sortOrder int) (*models.Request, error) {
	return h.service.Update(id, collectionID, name, method, url, headers, params, body, auth, script, sortOrder)
}

type SendRequestInput struct {
	Method         string
	URL            string
	Headers        map[string]string
	Body           string
	BodyType       string
	BodyFiles      []httpclient.BodyFile
	AuthType       string
	AuthData       map[string]string
	TimeoutMs      int
	FollowRedirect bool
}

func (h *RequestHandler) SendRequest(input SendRequestInput) (*httpclient.Response, error) {
	return h.service.Send(input.Method, input.URL, input.Headers, input.Body, input.BodyType, input.BodyFiles, input.AuthType, input.AuthData, input.TimeoutMs, input.FollowRedirect)
}

func (h *RequestHandler) DeleteRequest(id string) error {
	return h.service.Delete(id)
}

// --- Assertions ---

type AssertRule struct {
	Type   string `json:"type"`
	Target string `json:"target"`
	Value  string `json:"value"`
}

type AssertResult struct {
	Rule   AssertRule `json:"rule"`
	Passed bool       `json:"passed"`
	Actual string     `json:"actual"`
	Error  string     `json:"error"`
}

type RunAssertsInput struct {
	Method         string
	URL            string
	Headers        map[string]string
	Body           string
	BodyType       string
	BodyFiles      []httpclient.BodyFile
	AuthType       string
	AuthData       map[string]string
	TimeoutMs      int
	FollowRedirect bool
	Asserts        []AssertRule
}

type RunAssertsResponse struct {
	Response *httpclient.Response `json:"response"`
	Asserts  []AssertResult       `json:"asserts"`
}

func (h *RequestHandler) RunAsserts(input RunAssertsInput) (*RunAssertsResponse, error) {
	resp, err := h.service.Send(input.Method, input.URL, input.Headers, input.Body, input.BodyType, input.BodyFiles, input.AuthType, input.AuthData, input.TimeoutMs, input.FollowRedirect)
	if err != nil {
		return nil, err
	}
	var rules []services.AssertRule
	for _, a := range input.Asserts {
		rules = append(rules, services.AssertRule{
			Type:   a.Type,
			Target: a.Target,
			Value:  a.Value,
		})
	}
	results := h.assertSvc.Run(resp, rules)
	var out []AssertResult
	for _, r := range results {
		out = append(out, AssertResult{
			Rule: AssertRule{
				Type:   r.Rule.Type,
				Target: r.Rule.Target,
				Value:  r.Rule.Value,
			},
			Passed: r.Passed,
			Actual: r.Actual,
			Error:  r.Error,
		})
	}
	return &RunAssertsResponse{Response: resp, Asserts: out}, nil
}

// --- Code Generation ---

type CodegenInput struct {
	Method   string
	URL      string
	Headers  map[string]string
	Body     string
	BodyType string
}

func (h *RequestHandler) GenerateCode(input CodegenInput, lang string) (string, error) {
	return h.codegenSvc.Generate(services.CodegenRequest{
		Method:   input.Method,
		URL:      input.URL,
		Headers:  input.Headers,
		Body:     input.Body,
		BodyType: input.BodyType,
	}, lang)
}
