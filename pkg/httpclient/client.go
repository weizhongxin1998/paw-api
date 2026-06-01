package httpclient

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Client struct {
	mu sync.Mutex
	hc *http.Client
}

type BodyFile struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	FilePath string `json:"file_path"`
	Enabled  bool   `json:"enabled"`
}

type Request struct {
	Method        string
	URL           string
	Headers       map[string]string
	Body          string
	BodyType      string
	BodyFiles     []BodyFile
	TimeoutMs     int
	FollowRedirect bool
}

type Response struct {
	Status     int                 `json:"status"`
	StatusText string              `json:"status_text"`
	Headers    map[string][]string `json:"headers"`
	Body       string              `json:"body"`
	DurationMs int64               `json:"duration_ms"`
}

func NewClient() *Client {
	jar, _ := cookiejar.New(nil)
	return NewClientWithJar(jar)
}

func NewClientWithJar(jar http.CookieJar) *Client {
	if jar == nil {
		jar, _ = cookiejar.New(nil)
	}
	return &Client{
		hc: &http.Client{
			Timeout: 30 * time.Second,
			Jar:     jar,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 10 {
					return http.ErrUseLastResponse
				}
				return nil
			},
		},
	}
}

func (c *Client) buildBody(req *Request) (io.Reader, string, error) {
	switch req.BodyType {
	case "none", "":
		return nil, "", nil

	case "json":
		return bytes.NewBufferString(req.Body), "application/json", nil

	case "text":
		return bytes.NewBufferString(req.Body), "text/plain", nil

	case "urlencoded":
		vals := url.Values{}
		var pairs []KeyValuePair
		if err := json.Unmarshal([]byte(req.Body), &pairs); err == nil {
			for _, p := range pairs {
				if p.Enabled {
					vals.Set(p.Key, p.Value)
				}
			}
		} else {
			lines := strings.Split(req.Body, "&")
			for _, line := range lines {
				if parts := strings.SplitN(line, "=", 2); len(parts) == 2 {
					vals.Set(parts[0], parts[1])
				}
			}
		}
		return strings.NewReader(vals.Encode()), "application/x-www-form-urlencoded", nil

	case "form-data":
		var buf bytes.Buffer
		w := multipart.NewWriter(&buf)
		for _, f := range req.BodyFiles {
			if !f.Enabled {
				continue
			}
			if f.FilePath != "" {
				file, err := os.Open(f.FilePath)
				if err != nil {
					return nil, "", err
				}
				part, err := w.CreateFormFile(f.Key, filepath.Base(f.FilePath))
				if err != nil {
					file.Close()
					return nil, "", err
				}
				io.Copy(part, file)
				file.Close()
			} else {
				w.WriteField(f.Key, f.Value)
			}
		}
		w.Close()
		return &buf, w.FormDataContentType(), nil

	case "binary":
		data, err := base64.StdEncoding.DecodeString(req.Body)
		if err != nil {
			return nil, "", err
		}
		return bytes.NewReader(data), "application/octet-stream", nil

	default:
		return bytes.NewBufferString(req.Body), "", nil
	}
}

type KeyValuePair struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Enabled bool   `json:"enabled"`
}

func (c *Client) Do(req *Request) (*Response, error) {
	bodyReader, contentType, err := c.buildBody(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(req.Method, req.URL, bodyReader)
	if err != nil {
		return nil, err
	}

	if contentType != "" {
		if _, has := req.Headers["Content-Type"]; !has {
			httpReq.Header.Set("Content-Type", contentType)
		}
	}

	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	c.mu.Lock()
	if req.TimeoutMs > 0 {
		c.hc.Timeout = time.Duration(req.TimeoutMs) * time.Millisecond
	} else {
		c.hc.Timeout = 30 * time.Second
	}

	if !req.FollowRedirect {
		c.hc.CheckRedirect = func(r *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	} else {
		c.hc.CheckRedirect = func(r *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return http.ErrUseLastResponse
			}
			return nil
		}
	}
	c.mu.Unlock()

	start := time.Now()
	httpResp, err := c.hc.Do(httpReq)
	duration := time.Since(start).Milliseconds()
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		Status:     httpResp.StatusCode,
		StatusText: http.StatusText(httpResp.StatusCode),
		Headers:    httpResp.Header,
		Body:       string(respBody),
		DurationMs: duration,
	}, nil
}

func (c *Client) GetCookies() []*http.Cookie {
	if c.hc.Jar == nil {
		return nil
	}
	// Return all cookies — iterate a dummy URL to extract
	return nil
}

func (c *Client) SetCookieJar(jar http.CookieJar) {
	c.hc.Jar = jar
}

func (c *Client) ClearCookies() {
	jar, _ := cookiejar.New(nil)
	c.hc.Jar = jar
}
