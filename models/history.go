package models

import "time"

type History struct {
	ID              string    `json:"id"`
	ProjectID       string    `json:"project_id"`
	RequestID       *string   `json:"request_id"`
	Method          string    `json:"method"`
	URL             string    `json:"url"`
	Headers         string    `json:"headers"`
	Body            string    `json:"body"`
	ResponseStatus  int       `json:"response_status"`
	ResponseBody    string    `json:"response_body"`
	ResponseHeaders string    `json:"response_headers"`
	DurationMs      int       `json:"duration_ms"`
	CreatedAt       time.Time `json:"created_at"`
}
