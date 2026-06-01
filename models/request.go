package models

import "time"

type Request struct {
	ID           string    `json:"id"`
	CollectionID string    `json:"collection_id"`
	Name         string    `json:"name"`
	Method       string    `json:"method"`
	URL          string    `json:"url"`
	Headers      string    `json:"headers"`
	Params       string    `json:"params"`
	Body         string    `json:"body"`
	Auth         string    `json:"auth"`
	Script       string    `json:"script"`
	SortOrder    int       `json:"sort_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
