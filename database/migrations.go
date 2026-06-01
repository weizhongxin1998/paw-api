package database

import "database/sql"

func runMigrations(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS projects (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS collections (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			parent_id TEXT,
			name TEXT NOT NULL,
			sort_order INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS requests (
			id TEXT PRIMARY KEY,
			collection_id TEXT NOT NULL,
			name TEXT NOT NULL,
			method TEXT DEFAULT 'GET',
			url TEXT DEFAULT '',
			headers TEXT DEFAULT '{}',
			params TEXT DEFAULT '{}',
			body TEXT DEFAULT '{}',
			auth TEXT DEFAULT '{}',
			script TEXT DEFAULT '',
			sort_order INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS environments (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			name TEXT NOT NULL,
			variables TEXT DEFAULT '[]',
			is_active INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS history (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			request_id TEXT,
			method TEXT NOT NULL,
			url TEXT DEFAULT '',
			headers TEXT DEFAULT '{}',
			body TEXT DEFAULT '{}',
			response_status INTEGER DEFAULT 0,
			response_body TEXT DEFAULT '',
			response_headers TEXT DEFAULT '{}',
			duration_ms INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
		)`,
	}
	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}
