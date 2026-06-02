package repositories

import (
	"database/sql"
	"time"
)

type SettingsRepository struct {
	db *sql.DB
}

func NewSettingsRepo(db *sql.DB) *SettingsRepository {
	return &SettingsRepository{db: db}
}

func (r *SettingsRepository) Get(key string) (string, error) {
	var value string
	err := r.db.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (r *SettingsRepository) Set(key, value string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := r.db.Exec(
		"INSERT OR REPLACE INTO settings (key, value, updated_at) VALUES (?, ?, ?)",
		key, value, now,
	)
	return err
}

func (r *SettingsRepository) GetAll() (map[string]string, error) {
	rows, err := r.db.Query("SELECT key, value FROM settings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		result[key] = value
	}
	return result, nil
}

func (r *SettingsRepository) SetAll(settings map[string]string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now().UTC().Format(time.RFC3339)
	for k, v := range settings {
		if _, err := tx.Exec(
			"INSERT OR REPLACE INTO settings (key, value, updated_at) VALUES (?, ?, ?)",
			k, v, now,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}
