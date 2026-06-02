package database

import "database/sql"

func Run(db *sql.DB) error {
	for _, ddl := range []string{
		createProjects,
		createCollections,
		createRequests,
		createEnvironments,
		createEnvVariables,
		createHistory,
		createSettings,
	} {
		if _, err := db.Exec(ddl); err != nil {
			return err
		}
	}

	return insertDefaults(db)
}

const createProjects = `
CREATE TABLE IF NOT EXISTS projects (
    id          INTEGER PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at  TEXT NOT NULL,
    updated_at  TEXT NOT NULL
);`

const createCollections = `
CREATE TABLE IF NOT EXISTS collections (
    id          INTEGER PRIMARY KEY,
    project_id  INTEGER NOT NULL,
    parent_id   INTEGER DEFAULT NULL,
    name        TEXT NOT NULL,
    sort_order  INTEGER NOT NULL DEFAULT 0,
    created_at  TEXT NOT NULL,
    updated_at  TEXT NOT NULL,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES collections(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_collections_project ON collections(project_id);
CREATE INDEX IF NOT EXISTS idx_collections_parent ON collections(parent_id);`

const createRequests = `
CREATE TABLE IF NOT EXISTS requests (
    id              INTEGER PRIMARY KEY,
    collection_id   INTEGER NOT NULL,
    name            TEXT NOT NULL,
    description     TEXT NOT NULL DEFAULT '',
    method          TEXT NOT NULL,
    url             TEXT NOT NULL DEFAULT '',
    headers         TEXT NOT NULL DEFAULT '[]',
    params          TEXT NOT NULL DEFAULT '[]',
    body_type       TEXT NOT NULL DEFAULT 'none',
    body            TEXT NOT NULL DEFAULT '{}',
    auth            TEXT NOT NULL DEFAULT '{"type":"none"}',
    sort_order      INTEGER NOT NULL DEFAULT 0,
    created_at      TEXT NOT NULL,
    updated_at      TEXT NOT NULL,
    FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_requests_collection ON requests(collection_id);`

const createEnvironments = `
CREATE TABLE IF NOT EXISTS environments (
    id          INTEGER PRIMARY KEY,
    project_id  INTEGER NOT NULL,
    name        TEXT NOT NULL,
    is_active   INTEGER NOT NULL DEFAULT 0,
    created_at  TEXT NOT NULL,
    updated_at  TEXT NOT NULL,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_environments_project ON environments(project_id);`

const createEnvVariables = `
CREATE TABLE IF NOT EXISTS env_variables (
    id              INTEGER PRIMARY KEY,
    environment_id  INTEGER NOT NULL,
    key             TEXT NOT NULL,
    value           TEXT NOT NULL DEFAULT '',
    enabled         INTEGER NOT NULL DEFAULT 1,
    sort_order      INTEGER NOT NULL DEFAULT 0,
    created_at      TEXT NOT NULL,
    FOREIGN KEY (environment_id) REFERENCES environments(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_env_vars_env ON env_variables(environment_id);`

const createHistory = `
CREATE TABLE IF NOT EXISTS history (
    id                INTEGER PRIMARY KEY,
    project_id        INTEGER NOT NULL,
    request_id        INTEGER DEFAULT NULL,
    method            TEXT NOT NULL,
    url               TEXT NOT NULL,
    request_headers   TEXT NOT NULL DEFAULT '[]',
    request_body      TEXT NOT NULL DEFAULT '',
    response_status   INTEGER NOT NULL,
    response_headers  TEXT NOT NULL DEFAULT '[]',
    response_body     TEXT NOT NULL DEFAULT '',
    duration_ms       INTEGER NOT NULL,
    created_at        TEXT NOT NULL,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (request_id) REFERENCES requests(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_history_project ON history(project_id);
CREATE INDEX IF NOT EXISTS idx_history_created ON history(created_at);`

const createSettings = `
CREATE TABLE IF NOT EXISTS settings (
    key         TEXT PRIMARY KEY,
    value       TEXT NOT NULL,
    updated_at  TEXT NOT NULL
);`

func insertDefaults(db *sql.DB) error {
	defaults := map[string]string{
		"general.timeout":           "30",
		"general.follow_redirects":  "true",
		"general.max_redirects":     "10",
		"general.ssl_verify":        "true",
		"proxy.mode":                "\"none\"",
		"proxy.http":                "\"\"",
		"proxy.https":               "\"\"",
		"theme.mode":                "\"light\"",
		"theme.color":               "\"green\"",
		"theme.font_size":           "14",
		"history.retention_days":    "30",
		"app.last_project_id":       "0",
	}

	for k, v := range defaults {
		_, err := db.Exec(
			"INSERT OR IGNORE INTO settings (key, value, updated_at) VALUES (?, ?, datetime('now'))",
			k, v,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
