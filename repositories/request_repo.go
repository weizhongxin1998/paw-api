package repositories

import (
	"database/sql"
	"time"

	"paw-api/models"
	"paw-api/pkg/snowflake"
)

type RequestRepository struct {
	db *sql.DB
	sf *snowflake.Generator
}

func NewRequestRepo(db *sql.DB, sf *snowflake.Generator) *RequestRepository {
	return &RequestRepository{db: db, sf: sf}
}

func (r *RequestRepository) ListByCollection(collectionID int64) ([]models.Request, error) {
	rows, err := r.db.Query(
		`SELECT id, collection_id, name, description, method, url, headers, params, body_type, body, auth, sort_order, created_at, updated_at
		 FROM requests WHERE collection_id = ? ORDER BY sort_order ASC`, collectionID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.Request
	for rows.Next() {
		var req models.Request
		if err := rows.Scan(
			&req.ID, &req.CollectionID, &req.Name, &req.Description, &req.Method, &req.URL,
			&req.Headers, &req.Params, &req.BodyType, &req.Body, &req.Auth, &req.SortOrder,
			&req.CreatedAt, &req.UpdatedAt,
		); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, nil
}

func (r *RequestRepository) GetByID(id int64) (*models.Request, error) {
	var req models.Request
	err := r.db.QueryRow(
		`SELECT id, collection_id, name, description, method, url, headers, params, body_type, body, auth, sort_order, created_at, updated_at
		 FROM requests WHERE id = ?`, id,
	).Scan(
		&req.ID, &req.CollectionID, &req.Name, &req.Description, &req.Method, &req.URL,
		&req.Headers, &req.Params, &req.BodyType, &req.Body, &req.Auth, &req.SortOrder,
		&req.CreatedAt, &req.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *RequestRepository) Create(request *models.Request) error {
	now := time.Now().UTC().Format(time.RFC3339)
	request.ID = r.sf.Next()
	request.CreatedAt = now
	request.UpdatedAt = now

	_, err := r.db.Exec(
		`INSERT INTO requests (id, collection_id, name, description, method, url, headers, params, body_type, body, auth, sort_order, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		request.ID, request.CollectionID, request.Name, request.Description, request.Method, request.URL,
		request.Headers, request.Params, request.BodyType, request.Body, request.Auth, request.SortOrder,
		request.CreatedAt, request.UpdatedAt,
	)
	return err
}

func (r *RequestRepository) Update(request *models.Request) error {
	request.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	_, err := r.db.Exec(
		`UPDATE requests SET name = ?, description = ?, method = ?, url = ?, headers = ?, params = ?,
		 body_type = ?, body = ?, auth = ?, sort_order = ?, updated_at = ? WHERE id = ?`,
		request.Name, request.Description, request.Method, request.URL,
		request.Headers, request.Params, request.BodyType, request.Body, request.Auth, request.SortOrder,
		request.UpdatedAt, request.ID,
	)
	return err
}

func (r *RequestRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM requests WHERE id = ?", id)
	return err
}

func (r *RequestRepository) UpdateSortOrder(id int64, sortOrder int) error {
	_, err := r.db.Exec("UPDATE requests SET sort_order = ? WHERE id = ?", sortOrder, id)
	return err
}

func (r *RequestRepository) MoveToCollection(id int64, collectionID int64) error {
	_, err := r.db.Exec("UPDATE requests SET collection_id = ? WHERE id = ?", collectionID, id)
	return err
}
