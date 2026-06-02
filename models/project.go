package models

type Project struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type ProjectStats struct {
	RequestCount    int `json:"request_count"`
	CollectionCount int `json:"collection_count"`
}
