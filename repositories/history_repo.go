package repositories

import (
	"database/sql"
	"time"

	"paw-api/models"
	"paw-api/pkg/snowflake"
)

type HistoryRepository struct {
	db *sql.DB
	sf *snowflake.Generator
}

func NewHistoryRepo(db *sql.DB, sf *snowflake.Generator) *HistoryRepository {
	return &HistoryRepository{db: db, sf: sf}
}

func (r *HistoryRepository) ListByProject(projectID int64, limit, offset int) ([]models.History, error) {
	rows, err := r.db.Query(
		`SELECT id, project_id, request_id, method, url, request_headers, request_body,
		 response_status, response_headers, response_body, duration_ms, created_at
		 FROM history WHERE project_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		projectID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanHistories(rows)
}

func (r *HistoryRepository) Search(projectID int64, keyword, method string, statusMin, statusMax int, limit, offset int) ([]models.History, int, error) {
	where := "WHERE project_id = ?"
	args := []interface{}{projectID}

	if keyword != "" {
		where += " AND url LIKE ?"
		args = append(args, "%"+keyword+"%")
	}
	if method != "" {
		where += " AND method = ?"
		args = append(args, method)
	}
	if statusMin > 0 {
		where += " AND response_status >= ?"
		args = append(args, statusMin)
	}
	if statusMax > 0 {
		where += " AND response_status < ?"
		args = append(args, statusMax)
	}

	var total int
	countQuery := "SELECT COUNT(*) FROM history " + where
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `SELECT id, project_id, request_id, method, url, request_headers, request_body,
		 response_status, response_headers, response_body, duration_ms, created_at
		 FROM history ` + where + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items, err := scanHistories(rows)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *HistoryRepository) GetByID(id int64) (*models.History, error) {
	var h models.History
	var requestID sql.NullInt64
	err := r.db.QueryRow(
		`SELECT id, project_id, request_id, method, url, request_headers, request_body,
		 response_status, response_headers, response_body, duration_ms, created_at
		 FROM history WHERE id = ?`, id,
	).Scan(
		&h.ID, &h.ProjectID, &requestID, &h.Method, &h.URL,
		&h.RequestHeaders, &h.RequestBody, &h.ResponseStatus,
		&h.ResponseHeaders, &h.ResponseBody, &h.DurationMs, &h.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if requestID.Valid {
		h.RequestID = &requestID.Int64
	}
	return &h, nil
}

func (r *HistoryRepository) Create(history *models.History) error {
	now := time.Now().UTC().Format(time.RFC3339)
	history.ID = r.sf.Next()
	history.CreatedAt = now

	_, err := r.db.Exec(
		`INSERT INTO history (id, project_id, request_id, method, url, request_headers, request_body,
		 response_status, response_headers, response_body, duration_ms, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		history.ID, history.ProjectID, history.RequestID, history.Method, history.URL,
		history.RequestHeaders, history.RequestBody, history.ResponseStatus,
		history.ResponseHeaders, history.ResponseBody, history.DurationMs, history.CreatedAt,
	)
	return err
}

func (r *HistoryRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM history WHERE id = ?", id)
	return err
}

func (r *HistoryRepository) DeleteByProject(projectID int64) error {
	_, err := r.db.Exec("DELETE FROM history WHERE project_id = ?", projectID)
	return err
}

func (r *HistoryRepository) DeleteOlderThan(projectID int64, days int) error {
	_, err := r.db.Exec(
		"DELETE FROM history WHERE project_id = ? AND created_at < datetime('now', ? || ' days')",
		projectID, -days,
	)
	return err
}

func scanHistories(rows *sql.Rows) ([]models.History, error) {
	var items []models.History
	for rows.Next() {
		var h models.History
		var requestID sql.NullInt64
		if err := rows.Scan(
			&h.ID, &h.ProjectID, &requestID, &h.Method, &h.URL,
			&h.RequestHeaders, &h.RequestBody, &h.ResponseStatus,
			&h.ResponseHeaders, &h.ResponseBody, &h.DurationMs, &h.CreatedAt,
		); err != nil {
			return nil, err
		}
		if requestID.Valid {
			h.RequestID = &requestID.Int64
		}
		items = append(items, h)
	}
	return items, nil
}
