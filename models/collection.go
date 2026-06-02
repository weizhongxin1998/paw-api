package models

type Collection struct {
	ID        int64  `json:"id"`
	ProjectID int64  `json:"project_id"`
	ParentID  *int64 `json:"parent_id"`
	Name      string `json:"name"`
	SortOrder int    `json:"sort_order"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type TreeItem struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Type      string     `json:"type"`
	Method    string     `json:"method,omitempty"`
	URL       string     `json:"url,omitempty"`
	Children  []TreeItem `json:"children"`
	SortOrder int        `json:"sort_order"`
}
