package repositories

import (
	"database/sql"
	"paw-api/database"
	"paw-api/models"
)

type HistoryRepo struct{}

func (r *HistoryRepo) Create(h *models.History) error {
	_, err := database.DB.Exec(
		`INSERT INTO history (id, project_id, request_id, method, url, headers, body, response_status, response_body, response_headers, duration_ms, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		h.ID, h.ProjectID, h.RequestID, h.Method, h.URL, h.Headers, h.Body, h.ResponseStatus, h.ResponseBody, h.ResponseHeaders, h.DurationMs, h.CreatedAt,
	)
	return err
}

func (r *HistoryRepo) ListByProject(projectID string, limit int) ([]models.History, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := database.DB.Query(
		`SELECT id, project_id, request_id, method, url, headers, body, response_status, response_body, response_headers, duration_ms, created_at FROM history WHERE project_id = ? ORDER BY created_at DESC LIMIT ?`, projectID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	history := make([]models.History, 0)
	for rows.Next() {
		var h models.History
		var requestID sql.NullString
		if err := rows.Scan(&h.ID, &h.ProjectID, &requestID, &h.Method, &h.URL, &h.Headers, &h.Body, &h.ResponseStatus, &h.ResponseBody, &h.ResponseHeaders, &h.DurationMs, &h.CreatedAt); err != nil {
			return nil, err
		}
		if requestID.Valid {
			h.RequestID = &requestID.String
		}
		history = append(history, h)
	}
	return history, rows.Err()
}

func (r *HistoryRepo) Delete(id string) error {
	_, err := database.DB.Exec(`DELETE FROM history WHERE id = ?`, id)
	return err
}

func (r *HistoryRepo) ClearByProject(projectID string) error {
	_, err := database.DB.Exec(`DELETE FROM history WHERE project_id = ?`, projectID)
	return err
}
