package models

type History struct {
	ID              int64  `json:"id"`
	ProjectID       int64  `json:"project_id"`
	RequestID       *int64 `json:"request_id"`
	Method          string `json:"method"`
	URL             string `json:"url"`
	RequestHeaders  string `json:"request_headers"`
	RequestBody     string `json:"request_body"`
	ResponseStatus  int    `json:"response_status"`
	ResponseHeaders string `json:"response_headers"`
	ResponseBody    string `json:"response_body"`
	DurationMs      int    `json:"duration_ms"`
	CreatedAt       string `json:"created_at"`
}
