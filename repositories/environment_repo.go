package repositories

import (
	"database/sql"
	"paw-api/database"
	"paw-api/models"
)

type EnvironmentRepo struct{}

func (r *EnvironmentRepo) Create(env *models.Environment) error {
	_, err := database.DB.Exec(
		`INSERT INTO environments (id, project_id, name, variables, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		env.ID, env.ProjectID, env.Name, env.Variables, env.IsActive, env.CreatedAt, env.UpdatedAt,
	)
	return err
}

func (r *EnvironmentRepo) GetByID(id string) (*models.Environment, error) {
	e := &models.Environment{}
	err := database.DB.QueryRow(
		`SELECT id, project_id, name, variables, is_active, created_at, updated_at FROM environments WHERE id = ?`, id,
	).Scan(&e.ID, &e.ProjectID, &e.Name, &e.Variables, &e.IsActive, &e.CreatedAt, &e.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return e, err
}

func (r *EnvironmentRepo) ListByProject(projectID string) ([]models.Environment, error) {
	rows, err := database.DB.Query(
		`SELECT id, project_id, name, variables, is_active, created_at, updated_at FROM environments WHERE project_id = ? ORDER BY name`, projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	envs := make([]models.Environment, 0)
	for rows.Next() {
		var e models.Environment
		if err := rows.Scan(&e.ID, &e.ProjectID, &e.Name, &e.Variables, &e.IsActive, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		envs = append(envs, e)
	}
	return envs, rows.Err()
}

func (r *EnvironmentRepo) Update(env *models.Environment) error {
	_, err := database.DB.Exec(
		`UPDATE environments SET name = ?, variables = ?, is_active = ?, updated_at = ? WHERE id = ?`,
		env.Name, env.Variables, env.IsActive, env.UpdatedAt, env.ID,
	)
	return err
}

func (r *EnvironmentRepo) DeactivateAll(projectID string) error {
	_, err := database.DB.Exec(
		`UPDATE environments SET is_active = 0 WHERE project_id = ?`, projectID,
	)
	return err
}

func (r *EnvironmentRepo) Delete(id string) error {
	_, err := database.DB.Exec(`DELETE FROM environments WHERE id = ?`, id)
	return err
}
