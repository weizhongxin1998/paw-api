package models

type Request struct {
	ID           int64  `json:"id"`
	CollectionID int64  `json:"collection_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Method       string `json:"method"`
	URL          string `json:"url"`
	Headers      string `json:"headers"`
	Params       string `json:"params"`
	BodyType     string `json:"body_type"`
	Body         string `json:"body"`
	Auth         string `json:"auth"`
	SortOrder    int    `json:"sort_order"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
