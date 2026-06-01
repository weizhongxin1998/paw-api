package repositories

import (
	"database/sql"
	"paw-api/database"
	"paw-api/models"
)

type RequestRepo struct{}

func (r *RequestRepo) Create(req *models.Request) error {
	_, err := database.DB.Exec(
		`INSERT INTO requests (id, collection_id, name, method, url, headers, params, body, auth, script, sort_order, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		req.ID, req.CollectionID, req.Name, req.Method, req.URL, req.Headers, req.Params, req.Body, req.Auth, req.Script, req.SortOrder, req.CreatedAt, req.UpdatedAt,
	)
	return err
}

func (r *RequestRepo) GetByID(id string) (*models.Request, error) {
	req := &models.Request{}
	err := database.DB.QueryRow(
		`SELECT id, collection_id, name, method, url, headers, params, body, auth, script, sort_order, created_at, updated_at FROM requests WHERE id = ?`, id,
	).Scan(&req.ID, &req.CollectionID, &req.Name, &req.Method, &req.URL, &req.Headers, &req.Params, &req.Body, &req.Auth, &req.Script, &req.SortOrder, &req.CreatedAt, &req.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return req, err
}

func (r *RequestRepo) ListByCollection(collectionID string) ([]models.Request, error) {
	rows, err := database.DB.Query(
		`SELECT id, collection_id, name, method, url, headers, params, body, auth, script, sort_order, created_at, updated_at FROM requests WHERE collection_id = ? ORDER BY sort_order`, collectionID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var requests []models.Request
	for rows.Next() {
		var req models.Request
		if err := rows.Scan(&req.ID, &req.CollectionID, &req.Name, &req.Method, &req.URL, &req.Headers, &req.Params, &req.Body, &req.Auth, &req.Script, &req.SortOrder, &req.CreatedAt, &req.UpdatedAt); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, rows.Err()
}

func (r *RequestRepo) Update(req *models.Request) error {
	_, err := database.DB.Exec(
		`UPDATE requests SET name = ?, method = ?, url = ?, headers = ?, params = ?, body = ?, auth = ?, script = ?, sort_order = ?, updated_at = ? WHERE id = ?`,
		req.Name, req.Method, req.URL, req.Headers, req.Params, req.Body, req.Auth, req.Script, req.SortOrder, req.UpdatedAt, req.ID,
	)
	return err
}

func (r *RequestRepo) Delete(id string) error {
	_, err := database.DB.Exec(`DELETE FROM requests WHERE id = ?`, id)
	return err
}
