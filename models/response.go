package models

type HTTPResponse struct {
	Status      int               `json:"status"`
	StatusText  string            `json:"status_text"`
	Time        int64             `json:"time"`
	Size        int64             `json:"size"`
	Headers     map[string]string `json:"headers"`
	Cookies     []Cookie          `json:"cookies"`
	Body        string            `json:"body"`
	RawRequest  string            `json:"raw_request"`
	CurlCommand string            `json:"curl_command"`
}

type Cookie struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Domain string `json:"domain"`
	Path   string `json:"path"`
}
