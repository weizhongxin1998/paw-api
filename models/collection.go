package models

import "time"

type Collection struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	ParentID  *string   `json:"parent_id"`
	Name      string    `json:"name"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
