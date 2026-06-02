package models

type Environment struct {
	ID        int64  `json:"id"`
	ProjectID int64  `json:"project_id"`
	Name      string `json:"name"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type EnvVariable struct {
	ID            int64  `json:"id"`
	EnvironmentID int64  `json:"environment_id"`
	Key           string `json:"key"`
	Value         string `json:"value"`
	Enabled       bool   `json:"enabled"`
	SortOrder     int    `json:"sort_order"`
	CreatedAt     string `json:"created_at"`
}
