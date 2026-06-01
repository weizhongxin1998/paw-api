package repositories

import (
	"database/sql"
	"paw-api/database"
	"paw-api/models"
)

type ProjectRepo struct{}

func (r *ProjectRepo) Create(project *models.Project) error {
	_, err := database.DB.Exec(
		`INSERT INTO projects (id, name, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		project.ID, project.Name, project.Description, project.CreatedAt, project.UpdatedAt,
	)
	return err
}

func (r *ProjectRepo) GetByID(id string) (*models.Project, error) {
	p := &models.Project{}
	err := database.DB.QueryRow(
		`SELECT id, name, description, created_at, updated_at FROM projects WHERE id = ?`, id,
	).Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return p, err
}

func (r *ProjectRepo) List() ([]models.Project, error) {
	rows, err := database.DB.Query(`SELECT id, name, description, created_at, updated_at FROM projects ORDER BY created_at DESC`)
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
	return projects, rows.Err()
}

func (r *ProjectRepo) Update(project *models.Project) error {
	_, err := database.DB.Exec(
		`UPDATE projects SET name = ?, description = ?, updated_at = ? WHERE id = ?`,
		project.Name, project.Description, project.UpdatedAt, project.ID,
	)
	return err
}

func (r *ProjectRepo) Delete(id string) error {
	_, err := database.DB.Exec(`DELETE FROM projects WHERE id = ?`, id)
	return err
}
