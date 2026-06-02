package httpclient

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

type Config struct {
	Timeout        time.Duration
	FollowRedirect bool
	MaxRedirects   int
	VerifySSL      bool
}

type Client struct {
	client  *http.Client
	proxyURL *url.URL
}

func New() *Client {
	return &Client{
		client: &http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				TLSHandshakeTimeout: 10 * time.Second,
			},
		},
	}
}

func (c *Client) Execute(ctx context.Context, method, urlStr string, headers map[string]string, body []byte) (*Response, error) {
	start := time.Now()

	var reader io.Reader
	if len(body) > 0 {
		reader = io.NopCloser(io.Reader(nil))
		reader = bodyReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, urlStr, reader)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respHeaders := make(map[string]string)
	for k := range resp.Header {
		respHeaders[k] = resp.Header.Get(k)
	}

	elapsed := time.Since(start).Milliseconds()

	return &Response{
		Status:     resp.StatusCode,
		StatusText: resp.Status,
		Headers:    respHeaders,
		Body:       respBody,
		TimeMs:     elapsed,
		Size:       int64(len(respBody)),
	}, nil
}

func (c *Client) ApplySettings(settings map[string]string) {
	// To be implemented when settings management is ready
	_ = settings
}

func (c *Client) EnableSSLVerify(verify bool) {
	if c.client.Transport != nil {
		transport := c.client.Transport.(*http.Transport).Clone()
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: !verify,
		}
		c.client.Transport = transport
	}
}

type Response struct {
	Status     int
	StatusText string
	Headers    map[string]string
	Body       []byte
	TimeMs     int64
	Size       int64
}

type bodyReader []byte

func (b bodyReader) Read(p []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, io.EOF
	}
	n = copy(p, b)
	return n, io.EOF
}
