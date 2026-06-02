package repositories

import (
	"database/sql"
	"paw-api/database"
	"paw-api/models"
)

type CollectionRepo struct{}

func scanCollection(scanner interface {
	Scan(dest ...interface{}) error
}) (*models.Collection, error) {
	var c models.Collection
	var parentID sql.NullString
	err := scanner.Scan(&c.ID, &c.ProjectID, &parentID, &c.Name, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if parentID.Valid {
		c.ParentID = &parentID.String
	}
	return &c, nil
}

func (r *CollectionRepo) Create(c *models.Collection) error {
	_, err := database.DB.Exec(
		`INSERT INTO collections (id, project_id, parent_id, name, sort_order, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		c.ID, c.ProjectID, c.ParentID, c.Name, c.SortOrder, c.CreatedAt, c.UpdatedAt,
	)
	return err
}

func (r *CollectionRepo) GetByID(id string) (*models.Collection, error) {
	row := database.DB.QueryRow(
		`SELECT id, project_id, parent_id, name, sort_order, created_at, updated_at FROM collections WHERE id = ?`, id,
	)
	c, err := scanCollection(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return c, err
}

func (r *CollectionRepo) ListByProject(projectID string) ([]models.Collection, error) {
	rows, err := database.DB.Query(
		`SELECT id, project_id, parent_id, name, sort_order, created_at, updated_at FROM collections WHERE project_id = ? ORDER BY sort_order`, projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	collections := make([]models.Collection, 0)
	for rows.Next() {
		c, err := scanCollection(rows)
		if err != nil {
			return nil, err
		}
		collections = append(collections, *c)
	}
	return collections, rows.Err()
}

func (r *CollectionRepo) Update(c *models.Collection) error {
	_, err := database.DB.Exec(
		`UPDATE collections SET name = ?, parent_id = ?, sort_order = ?, updated_at = ? WHERE id = ?`,
		c.Name, c.ParentID, c.SortOrder, c.UpdatedAt, c.ID,
	)
	return err
}

func (r *CollectionRepo) GetMaxSortOrder(projectID, parentID string) (int, error) {
	var max sql.NullInt64
	if parentID == "" {
		err := database.DB.QueryRow(`SELECT MAX(sort_order) FROM collections WHERE project_id = ? AND parent_id IS NULL`, projectID).Scan(&max)
		if err != nil {
			return 0, err
		}
	} else {
		err := database.DB.QueryRow(`SELECT MAX(sort_order) FROM collections WHERE project_id = ? AND parent_id = ?`, projectID, parentID).Scan(&max)
		if err != nil {
			return 0, err
		}
	}
	if max.Valid {
		return int(max.Int64), nil
	}
	return 0, nil
}

func (r *CollectionRepo) ListByParent(parentID string) ([]models.Collection, error) {
	rows, err := database.DB.Query(
		`SELECT id, project_id, parent_id, name, sort_order, created_at, updated_at FROM collections WHERE parent_id = ? ORDER BY sort_order`, parentID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	collections := make([]models.Collection, 0)
	for rows.Next() {
		c, err := scanCollection(rows)
		if err != nil {
			return nil, err
		}
		collections = append(collections, *c)
	}
	return collections, rows.Err()
}

func (r *CollectionRepo) Delete(id string) error {
	_, err := database.DB.Exec(`DELETE FROM collections WHERE id = ?`, id)
	return err
}
