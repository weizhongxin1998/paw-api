package repositories

import (
	"database/sql"
	"time"

	"paw-api/models"
	"paw-api/pkg/snowflake"
)

type EnvironmentRepository struct {
	db *sql.DB
	sf *snowflake.Generator
}

func NewEnvironmentRepo(db *sql.DB, sf *snowflake.Generator) *EnvironmentRepository {
	return &EnvironmentRepository{db: db, sf: sf}
}

func (r *EnvironmentRepository) ListByProject(projectID int64) ([]models.Environment, error) {
	rows, err := r.db.Query(
		"SELECT id, project_id, name, base_url, is_active, created_at, updated_at FROM environments WHERE project_id = ? ORDER BY created_at ASC",
		projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var envs []models.Environment
	for rows.Next() {
		var e models.Environment
		var isActive int
		if err := rows.Scan(&e.ID, &e.ProjectID, &e.Name, &e.BaseURL, &isActive, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		e.IsActive = isActive == 1
		envs = append(envs, e)
	}
	return envs, nil
}

func (r *EnvironmentRepository) GetByID(id int64) (*models.Environment, error) {
	var e models.Environment
	var isActive int
	err := r.db.QueryRow(
		"SELECT id, project_id, name, base_url, is_active, created_at, updated_at FROM environments WHERE id = ?", id,
	).Scan(&e.ID, &e.ProjectID, &e.Name, &e.BaseURL, &isActive, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	e.IsActive = isActive == 1
	return &e, nil
}

func (r *EnvironmentRepository) Create(env *models.Environment) error {
	now := time.Now().UTC().Format(time.RFC3339)
	env.ID = r.sf.Next()
	env.CreatedAt = now
	env.UpdatedAt = now

	isActive := 0
	if env.IsActive {
		isActive = 1
	}

	_, err := r.db.Exec(
		"INSERT INTO environments (id, project_id, name, base_url, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		env.ID, env.ProjectID, env.Name, env.BaseURL, isActive, env.CreatedAt, env.UpdatedAt,
	)
	return err
}

func (r *EnvironmentRepository) Update(env *models.Environment) error {
	env.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	isActive := 0
	if env.IsActive {
		isActive = 1
	}
	_, err := r.db.Exec(
		"UPDATE environments SET name = ?, base_url = ?, is_active = ?, updated_at = ? WHERE id = ?",
		env.Name, env.BaseURL, isActive, env.UpdatedAt, env.ID,
	)
	return err
}

func (r *EnvironmentRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM environments WHERE id = ?", id)
	return err
}

func (r *EnvironmentRepository) SetActive(projectID int64, envID int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("UPDATE environments SET is_active = 0 WHERE project_id = ?", projectID); err != nil {
		return err
	}
	if _, err := tx.Exec("UPDATE environments SET is_active = 1 WHERE id = ?", envID); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *EnvironmentRepository) GetActive(projectID int64) (*models.Environment, error) {
	var e models.Environment
	var isActive int
	err := r.db.QueryRow(
		"SELECT id, project_id, name, base_url, is_active, created_at, updated_at FROM environments WHERE project_id = ? AND is_active = 1",
		projectID,
	).Scan(&e.ID, &e.ProjectID, &e.Name, &e.BaseURL, &isActive, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	e.IsActive = isActive == 1
	return &e, nil
}

func (r *EnvironmentRepository) SaveBaseURL(id int64, baseURL string) error {
	_, err := r.db.Exec("UPDATE environments SET base_url = ?, updated_at = ? WHERE id = ?",
		baseURL, time.Now().UTC().Format(time.RFC3339), id)
	return err
}
