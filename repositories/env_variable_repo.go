package repositories

import (
	"database/sql"
	"time"

	"paw-api/models"
	"paw-api/pkg/snowflake"
)

type EnvVariableRepository struct {
	db *sql.DB
	sf *snowflake.Generator
}

func NewEnvVariableRepo(db *sql.DB, sf *snowflake.Generator) *EnvVariableRepository {
	return &EnvVariableRepository{db: db, sf: sf}
}

func (r *EnvVariableRepository) ListByEnvironment(envID int64) ([]models.EnvVariable, error) {
	rows, err := r.db.Query(
		"SELECT id, environment_id, key, value, enabled, sort_order, created_at FROM env_variables WHERE environment_id = ? ORDER BY sort_order ASC",
		envID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vars []models.EnvVariable
	for rows.Next() {
		var v models.EnvVariable
		var enabled int
		if err := rows.Scan(&v.ID, &v.EnvironmentID, &v.Key, &v.Value, &enabled, &v.SortOrder, &v.CreatedAt); err != nil {
			return nil, err
		}
		v.Enabled = enabled == 1
		vars = append(vars, v)
	}
	return vars, nil
}

func (r *EnvVariableRepository) Create(variable *models.EnvVariable) error {
	now := time.Now().UTC().Format(time.RFC3339)
	variable.ID = r.sf.Next()
	variable.CreatedAt = now

	enabled := 0
	if variable.Enabled {
		enabled = 1
	}

	_, err := r.db.Exec(
		"INSERT INTO env_variables (id, environment_id, key, value, enabled, sort_order, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		variable.ID, variable.EnvironmentID, variable.Key, variable.Value, enabled, variable.SortOrder, variable.CreatedAt,
	)
	return err
}

func (r *EnvVariableRepository) Update(variable *models.EnvVariable) error {
	enabled := 0
	if variable.Enabled {
		enabled = 1
	}
	_, err := r.db.Exec(
		"UPDATE env_variables SET key = ?, value = ?, enabled = ?, sort_order = ? WHERE id = ?",
		variable.Key, variable.Value, enabled, variable.SortOrder, variable.ID,
	)
	return err
}

func (r *EnvVariableRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM env_variables WHERE id = ?", id)
	return err
}

func (r *EnvVariableRepository) BatchReplace(envID int64, variables []models.EnvVariable) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM env_variables WHERE environment_id = ?", envID); err != nil {
		return err
	}

	now := time.Now().UTC().Format(time.RFC3339)
	for i, v := range variables {
		v.ID = r.sf.Next()
		v.EnvironmentID = envID
		v.SortOrder = i
		v.CreatedAt = now

		enabled := 0
		if v.Enabled {
			enabled = 1
		}

		if _, err := tx.Exec(
			"INSERT INTO env_variables (id, environment_id, key, value, enabled, sort_order, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
			v.ID, v.EnvironmentID, v.Key, v.Value, enabled, v.SortOrder, v.CreatedAt,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}
