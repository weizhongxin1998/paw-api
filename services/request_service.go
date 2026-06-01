package services

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"paw-api/models"
	"paw-api/pkg/httpclient"
	"paw-api/repositories"

	"github.com/google/uuid"
)

type RequestService struct {
	repo   *repositories.RequestRepo
	client *httpclient.Client
}

func NewRequestService() *RequestService {
	svc := &RequestService{repo: &repositories.RequestRepo{}}
	svc.client = httpclient.NewClientWithJar(getSharedJar())
	return svc
}

func (s *RequestService) Create(collectionID, name, method, url, headers, params, body, auth, script string, sortOrder int) (*models.Request, error) {
	if name == "" {
		return nil, errors.New("request name is required")
	}
	if sortOrder <= 0 {
		maxOrder, err := s.repo.GetMaxSortOrder(collectionID)
		if err != nil {
			return nil, err
		}
		sortOrder = maxOrder + 1
	}
	now := time.Now()
	r := &models.Request{
		ID:           uuid.New().String(),
		CollectionID: collectionID,
		Name:         name,
		Method:       method,
		URL:          url,
		Headers:      headers,
		Params:       params,
		Body:         body,
		Auth:         auth,
		Script:       script,
		SortOrder:    sortOrder,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	return r, s.repo.Create(r)
}

func (s *RequestService) GetByID(id string) (*models.Request, error) {
	return s.repo.GetByID(id)
}

func (s *RequestService) ListByCollection(collectionID string) ([]models.Request, error) {
	return s.repo.ListByCollection(collectionID)
}

func (s *RequestService) ListByProjectID(projectID string) ([]models.Request, error) {
	return s.repo.ListByProjectID(projectID)
}

func (s *RequestService) Update(id, collectionID, name, method, url, headers, params, body, auth, script string, sortOrder int) (*models.Request, error) {
	r, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, errors.New("request not found")
	}
	if collectionID != "" {
		r.CollectionID = collectionID
	}
	r.Name = name
	r.Method = method
	r.URL = url
	r.Headers = headers
	r.Params = params
	r.Body = body
	r.Auth = auth
	r.Script = script
	r.SortOrder = sortOrder
	r.UpdatedAt = time.Now()
	return r, s.repo.Update(r)
}

func (s *RequestService) Send(method, url string, headers map[string]string, body string, bodyType string, bodyFiles []httpclient.BodyFile, authType string, authData map[string]string, timeoutMs int, followRedirect bool) (*httpclient.Response, error) {
	headers = s.applyAuth(headers, authType, authData, url, body, bodyType)
	return s.client.Do(&httpclient.Request{
		Method:         method,
		URL:            url,
		Headers:        headers,
		Body:           body,
		BodyType:       bodyType,
		BodyFiles:      bodyFiles,
		TimeoutMs:      timeoutMs,
		FollowRedirect: followRedirect,
	})
}

func (s *RequestService) applyAuth(headers map[string]string, authType string, authData map[string]string, requestURL, body, bodyType string) map[string]string {
	if headers == nil {
		headers = make(map[string]string)
	}
	switch authType {
	case "basic":
		user := authData["username"]
		pass := authData["password"]
		if user != "" || pass != "" {
			raw := user + ":" + pass
			encoded := base64.StdEncoding.EncodeToString([]byte(raw))
			headers["Authorization"] = "Basic " + encoded
		}

	case "bearer":
		if token := authData["token"]; token != "" {
			headers["Authorization"] = "Bearer " + token
		}

	case "digest":
		user := authData["username"]
		pass := authData["password"]
		if user != "" && pass != "" {
			headers = s.applyDigestAuth(headers, requestURL, user, pass, methodFromAuthData(authData), body, bodyType)
		}

	case "oauth2":
		tokenURL := authData["token_url"]
		clientID := authData["client_id"]
		clientSecret := authData["client_secret"]
		if tokenURL != "" && clientID != "" {
			token := s.fetchOAuth2Token(tokenURL, clientID, clientSecret)
			if token != "" {
				headers["Authorization"] = "Bearer " + token
			}
		}
	}
	return headers
}

func methodFromAuthData(authData map[string]string) string {
	if m, ok := authData["method"]; ok && m != "" {
		return m
	}
	return "GET"
}

func (s *RequestService) applyDigestAuth(headers map[string]string, rawURL, user, pass, method, body, bodyType string) map[string]string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return headers
	}
	var probeBody io.Reader
	if body != "" && bodyType != "none" {
		probeBody = strings.NewReader(body)
	}
	req, _ := http.NewRequest(method, rawURL, probeBody)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{Timeout: 10 * time.Second, CheckRedirect: func(r *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }}
	resp, err := client.Do(req)
	if err != nil {
		return headers
	}
	resp.Body.Close()

	wwwAuth := resp.Header.Get("WWW-Authenticate")
	if wwwAuth == "" {
		return headers
	}

	realm := extractDigestParam(wwwAuth, "realm")
	nonce := extractDigestParam(wwwAuth, "nonce")
	opaque := extractDigestParam(wwwAuth, "opaque")
	if realm == "" || nonce == "" {
		return headers
	}

	cnonce := make([]byte, 8)
	rand.Read(cnonce)
	nc := "00000001"
	ha1 := md5Hex(user + ":" + realm + ":" + pass)
	ha2 := md5Hex(method + ":" + u.Path)
	response := md5Hex(ha1 + ":" + nonce + ":" + nc + ":" + hex.EncodeToString(cnonce) + ":" + "auth" + ":" + ha2)

	authHeader := fmt.Sprintf(`Digest username="%s", realm="%s", nonce="%s", uri="%s", response="%s"`, user, realm, nonce, u.Path, response)
	if opaque != "" {
		authHeader += fmt.Sprintf(`, opaque="%s"`, opaque)
	}
	authHeader += `, cnonce="` + hex.EncodeToString(cnonce) + `", nc=` + nc + `, qop=auth`
	headers["Authorization"] = authHeader
	return headers
}

func extractDigestParam(challenge, param string) string {
	marker := param + "="
	idx := strings.Index(challenge, marker)
	if idx == -1 {
		return ""
	}
	val := challenge[idx+len(marker):]
	if strings.HasPrefix(val, "\"") {
		val = val[1:]
		end := strings.Index(val, "\"")
		if end == -1 {
			return ""
		}
		return val[:end]
	}
	end := strings.IndexAny(val, ", \t")
	if end == -1 {
		return val
	}
	return val[:end]
}

func md5Hex(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

func (s *RequestService) fetchOAuth2Token(tokenURL, clientID, clientSecret string) string {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return ""
	}
	if token, ok := result["access_token"].(string); ok {
		return token
	}
	return ""
}

func (s *RequestService) Delete(id string) error {
	return s.repo.Delete(id)
}
