package repositories

import (
	"database/sql"
	"time"

	"paw-api/models"
	"paw-api/pkg/snowflake"
)

type ProjectRepository struct {
	db  *sql.DB
	sf  *snowflake.Generator
}

func NewProjectRepo(db *sql.DB, sf *snowflake.Generator) *ProjectRepository {
	return &ProjectRepository{db: db, sf: sf}
}

func (r *ProjectRepository) List() ([]models.Project, error) {
	rows, err := r.db.Query(
		"SELECT id, name, description, created_at, updated_at FROM projects ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func (r *ProjectRepository) GetByID(id int64) (*models.Project, error) {
	var p models.Project
	err := r.db.QueryRow(
		"SELECT id, name, description, created_at, updated_at FROM projects WHERE id = ?", id,
	).Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProjectRepository) Create(project *models.Project) error {
	now := time.Now().UTC().Format(time.RFC3339)
	project.ID = r.sf.Next()
	project.CreatedAt = now
	project.UpdatedAt = now

	_, err := r.db.Exec(
		"INSERT INTO projects (id, name, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		project.ID, project.Name, project.Description, project.CreatedAt, project.UpdatedAt,
	)
	return err
}

func (r *ProjectRepository) Update(project *models.Project) error {
	project.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	_, err := r.db.Exec(
		"UPDATE projects SET name = ?, description = ?, updated_at = ? WHERE id = ?",
		project.Name, project.Description, project.UpdatedAt, project.ID,
	)
	return err
}

func (r *ProjectRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM projects WHERE id = ?", id)
	return err
}
