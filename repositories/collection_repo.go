package repositories

import (
	"database/sql"
	"time"

	"paw-api/models"
	"paw-api/pkg/snowflake"
)

type CollectionRepository struct {
	db *sql.DB
	sf *snowflake.Generator
}

func NewCollectionRepo(db *sql.DB, sf *snowflake.Generator) *CollectionRepository {
	return &CollectionRepository{db: db, sf: sf}
}

func (r *CollectionRepository) ListByProject(projectID int64) ([]models.Collection, error) {
	rows, err := r.db.Query(
		"SELECT id, project_id, parent_id, name, sort_order, created_at, updated_at FROM collections WHERE project_id = ? ORDER BY sort_order ASC",
		projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collections []models.Collection
	for rows.Next() {
		var c models.Collection
		var parentID sql.NullInt64
		if err := rows.Scan(&c.ID, &c.ProjectID, &parentID, &c.Name, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		if parentID.Valid {
			c.ParentID = &parentID.Int64
		}
		collections = append(collections, c)
	}
	return collections, nil
}

func (r *CollectionRepository) GetByID(id int64) (*models.Collection, error) {
	var c models.Collection
	var parentID sql.NullInt64
	err := r.db.QueryRow(
		"SELECT id, project_id, parent_id, name, sort_order, created_at, updated_at FROM collections WHERE id = ?", id,
	).Scan(&c.ID, &c.ProjectID, &parentID, &c.Name, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if parentID.Valid {
		c.ParentID = &parentID.Int64
	}
	return &c, nil
}

func (r *CollectionRepository) Create(collection *models.Collection) error {
	now := time.Now().UTC().Format(time.RFC3339)
	collection.ID = r.sf.Next()
	collection.CreatedAt = now
	collection.UpdatedAt = now

	_, err := r.db.Exec(
		"INSERT INTO collections (id, project_id, parent_id, name, sort_order, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		collection.ID, collection.ProjectID, collection.ParentID, collection.Name, collection.SortOrder, collection.CreatedAt, collection.UpdatedAt,
	)
	return err
}

func (r *CollectionRepository) Update(collection *models.Collection) error {
	collection.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	_, err := r.db.Exec(
		"UPDATE collections SET name = ?, sort_order = ?, updated_at = ? WHERE id = ?",
		collection.Name, collection.SortOrder, collection.UpdatedAt, collection.ID,
	)
	return err
}

func (r *CollectionRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM collections WHERE id = ?", id)
	return err
}

func (r *CollectionRepository) UpdateSortOrder(id int64, sortOrder int) error {
	_, err := r.db.Exec("UPDATE collections SET sort_order = ? WHERE id = ?", sortOrder, id)
	return err
}

func (r *CollectionRepository) MoveToParent(id int64, parentID *int64) error {
	_, err := r.db.Exec("UPDATE collections SET parent_id = ? WHERE id = ?", parentID, id)
	return err
}
