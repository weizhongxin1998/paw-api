package httpclient

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

type Client struct {
	hc *http.Client
}

type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    string
}

type Response struct {
	Status     int                 `json:"status"`
	StatusText string              `json:"status_text"`
	Headers    map[string][]string `json:"headers"`
	Body       string              `json:"body"`
	DurationMs int64               `json:"duration_ms"`
}

func NewClient() *Client {
	return &Client{
		hc: &http.Client{
			Timeout: 30 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 10 {
					return http.ErrUseLastResponse
				}
				return nil
			},
		},
	}
}

func (c *Client) Do(req *Request) (*Response, error) {
	var bodyReader io.Reader
	if req.Body != "" {
		bodyReader = bytes.NewBufferString(req.Body)
	}

	httpReq, err := http.NewRequest(req.Method, req.URL, bodyReader)
	if err != nil {
		return nil, err
	}

	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

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
