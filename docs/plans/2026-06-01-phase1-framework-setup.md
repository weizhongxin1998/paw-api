# Phase 1: Framework Setup Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Establish the full project skeleton: Go backend packages, SQLite database, Vue 3 frontend with Naive UI + Pinia + router, three-column layout, and theme system.

**Architecture:** Handlers (Wails Bind) → Services → Repositories → SQLite. Frontend uses Vue 3 SFCs with `<script setup>`, Pinia stores, and Naive UI components.

**Tech Stack:** Go 1.23, Wails v2, Vue 3 + TypeScript, Naive UI, Pinia, vue-router, `@vicons/ionicons5`, `modernc.org/sqlite`, `google/uuid`

---

## File Structure

```
paw-api/
├── database/
│   └── database.go            # SQLite connection + InitDB
│   └── migrations.go          # CREATE TABLE statements + migrations
├── models/                    # Data models
│   ├── project.go
│   ├── collection.go
│   ├── request.go
│   ├── environment.go
│   └── history.go
├── repositories/              # Data access layer
│   ├── project_repo.go
│   ├── collection_repo.go
│   ├── request_repo.go
│   ├── environment_repo.go
│   └── history_repo.go
├── services/                  # Business logic layer
│   ├── project_service.go
│   ├── collection_service.go
│   ├── request_service.go
│   ├── environment_service.go
│   └── history_service.go
├── handlers/                  # Wails binding layer
│   ├── project_handler.go
│   ├── collection_handler.go
│   ├── request_handler.go
│   ├── environment_handler.go
│   └── history_handler.go
├── pkg/
│   └── httpclient/            # placeholder
│       └── client.go
├── app.go                     # Updated: Wire up DB + handlers
├── main.go                    # Updated: Register all handlers
│
├── frontend/
│   ├── package.json           # Updated: Add naive-ui, pinia, vue-router, icons
│   ├── src/
│   │   ├── main.ts            # Updated: Bootstrap Naive UI + Pinia + Router
│   │   ├── App.vue            # Updated: NConfigProvider + layout
│   │   ├── style.css          # Updated: Global styles for desktop app
│   │   ├── router/
│   │   │   └── index.ts       # New: Route definitions
│   │   ├── stores/
│   │   │   ├── project.ts     # New: Project store stub
│   │   │   ├── request.ts     # New: Request store stub
│   │   │   ├── environment.ts # New: Environment store stub
│   │   │   ├── history.ts     # New: History store stub
│   │   │   └── tabs.ts        # New: Tabs store stub
│   │   ├── composables/
│   │   │   ├── useRequest.ts  # New: Request composable stub
│   │   │   ├── useCollection.ts # New: Collection composable stub
│   │   │   └── useTheme.ts    # New: Theme composable (light/dark + color switch)
│   │   ├── types/
│   │   │   ├── project.ts     # New: TS types for Project
│   │   │   ├── request.ts     # New: TS types for Request
│   │   │   └── environment.ts # New: TS types for Environment
│   │   ├── views/
│   │   │   ├── Workspace.vue  # New: Main workspace view
│   │   │   ├── History.vue    # New: History view stub
│   │   │   ├── Docs.vue       # New: Docs view stub
│   │   │   └── TestRunner.vue # New: Test runner view stub
│   │   └── components/
│   │       ├── AppSidebar.vue # New: Left sidebar (collection tree placeholder)
│   │       ├── TabBar.vue     # New: Top tab bar placeholder
│   │       ├── RequestEditor.vue # New: Request editor placeholder
│   │       └── ResponseViewer.vue # New: Response viewer placeholder
```

---

### Task 1: Go — Create database layer and models

**Files:**
- Create: `database/database.go`
- Create: `database/migrations.go`
- Create: `models/project.go`
- Create: `models/collection.go`
- Create: `models/request.go`
- Create: `models/environment.go`
- Create: `models/history.go`

- [ ] **Step 1: Create `models/project.go`**

```go
package models

import "time"

type Project struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
```

- [ ] **Step 2: Create `models/collection.go`**

```go
package models

import "time"

type Collection struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	ParentID  *string   `json:"parent_id"`
	Name      string    `json:"name"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
```

- [ ] **Step 3: Create `models/request.go`**

```go
package models

import "time"

type Request struct {
	ID           string    `json:"id"`
	CollectionID string    `json:"collection_id"`
	Name         string    `json:"name"`
	Method       string    `json:"method"`
	URL          string    `json:"url"`
	Headers      string    `json:"headers"`
	Params       string    `json:"params"`
	Body         string    `json:"body"`
	Auth         string    `json:"auth"`
	Script       string    `json:"script"`
	SortOrder    int       `json:"sort_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
```

- [ ] **Step 4: Create `models/environment.go`**

```go
package models

import "time"

type Environment struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	Name      string    `json:"name"`
	Variables string    `json:"variables"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
```

- [ ] **Step 5: Create `models/history.go`**

```go
package models

import "time"

type History struct {
	ID              string    `json:"id"`
	ProjectID       string    `json:"project_id"`
	RequestID       *string   `json:"request_id"`
	Method          string    `json:"method"`
	URL             string    `json:"url"`
	Headers         string    `json:"headers"`
	Body            string    `json:"body"`
	ResponseStatus  int       `json:"response_status"`
	ResponseBody    string    `json:"response_body"`
	ResponseHeaders string    `json:"response_headers"`
	DurationMs      int       `json:"duration_ms"`
	CreatedAt       time.Time `json:"created_at"`
}
```

- [ ] **Step 6: Create `database/database.go`**

```go
package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() error {
	appDataDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home dir: %w", err)
	}
	dbDir := filepath.Join(appDataDir, ".paw-api")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("failed to create db dir: %w", err)
	}
	dbPath := filepath.Join(dbDir, "paw-api.db")
	DB, err = sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)")
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}
	DB.SetMaxOpenConns(1)
	if err := runMigrations(DB); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
```

- [ ] **Step 7: Create `database/migrations.go`**

```go
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
```

- [ ] **Step 8: Add `modernc.org/sqlite` dependency**

Run: `cd D:\javap\paw-api && go get modernc.org/sqlite`

- [ ] **Step 9: Verify Go compiles**

Run: `cd D:\javap\paw-api && go build ./...`
Expected: No errors

---

### Task 2: Go — Create repositories

**Files:**
- Create: `repositories/project_repo.go`
- Create: `repositories/collection_repo.go`
- Create: `repositories/request_repo.go`
- Create: `repositories/environment_repo.go`
- Create: `repositories/history_repo.go`

- [ ] **Step 1: Create `repositories/project_repo.go`**

```go
package repositories

import (
	"database/sql"
	"paw-api/database"
	"paw-api/models"
)

type ProjectRepo struct{}

func (r *ProjectRepo) Create(project *models.Project) error {
	_, err := database.DB.Exec(
		`INSERT INTO projects (id, name, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		project.ID, project.Name, project.Description, project.CreatedAt, project.UpdatedAt,
	)
	return err
}

func (r *ProjectRepo) GetByID(id string) (*models.Project, error) {
	p := &models.Project{}
	err := database.DB.QueryRow(
		`SELECT id, name, description, created_at, updated_at FROM projects WHERE id = ?`, id,
	).Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return p, err
}

func (r *ProjectRepo) List() ([]models.Project, error) {
	rows, err := database.DB.Query(`SELECT id, name, description, created_at, updated_at FROM projects ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var projects []models.Project
	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

func (r *ProjectRepo) Update(project *models.Project) error {
	_, err := database.DB.Exec(
		`UPDATE projects SET name = ?, description = ?, updated_at = ? WHERE id = ?`,
		project.Name, project.Description, project.UpdatedAt, project.ID,
	)
	return err
}

func (r *ProjectRepo) Delete(id string) error {
	_, err := database.DB.Exec(`DELETE FROM projects WHERE id = ?`, id)
	return err
}
```

- [ ] **Step 2: Create `repositories/collection_repo.go`**

```go
package repositories

import (
	"database/sql"
	"paw-api/database"
	"paw-api/models"
)

type CollectionRepo struct{}

func (r *CollectionRepo) Create(c *models.Collection) error {
	_, err := database.DB.Exec(
		`INSERT INTO collections (id, project_id, parent_id, name, sort_order, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		c.ID, c.ProjectID, c.ParentID, c.Name, c.SortOrder, c.CreatedAt, c.UpdatedAt,
	)
	return err
}

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
	var collections []models.Collection
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

func (r *CollectionRepo) Delete(id string) error {
	_, err := database.DB.Exec(`DELETE FROM collections WHERE id = ?`, id)
	return err
}
```

- [ ] **Step 3: Create `repositories/request_repo.go`**

```go
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
```

- [ ] **Step 4: Create `repositories/environment_repo.go`**

```go
package repositories

import (
	"database/sql"
	"paw-api/database"
	"paw-api/models"
)

type EnvironmentRepo struct{}

func (r *EnvironmentRepo) Create(env *models.Environment) error {
	_, err := database.DB.Exec(
		`INSERT INTO environments (id, project_id, name, variables, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		env.ID, env.ProjectID, env.Name, env.Variables, env.IsActive, env.CreatedAt, env.UpdatedAt,
	)
	return err
}

func (r *EnvironmentRepo) GetByID(id string) (*models.Environment, error) {
	e := &models.Environment{}
	err := database.DB.QueryRow(
		`SELECT id, project_id, name, variables, is_active, created_at, updated_at FROM environments WHERE id = ?`, id,
	).Scan(&e.ID, &e.ProjectID, &e.Name, &e.Variables, &e.IsActive, &e.CreatedAt, &e.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return e, err
}

func (r *EnvironmentRepo) ListByProject(projectID string) ([]models.Environment, error) {
	rows, err := database.DB.Query(
		`SELECT id, project_id, name, variables, is_active, created_at, updated_at FROM environments WHERE project_id = ? ORDER BY name`, projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var envs []models.Environment
	for rows.Next() {
		var e models.Environment
		if err := rows.Scan(&e.ID, &e.ProjectID, &e.Name, &e.Variables, &e.IsActive, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		envs = append(envs, e)
	}
	return envs, rows.Err()
}

func (r *EnvironmentRepo) Update(env *models.Environment) error {
	_, err := database.DB.Exec(
		`UPDATE environments SET name = ?, variables = ?, is_active = ?, updated_at = ? WHERE id = ?`,
		env.Name, env.Variables, env.IsActive, env.UpdatedAt, env.ID,
	)
	return err
}

func (r *EnvironmentRepo) DeactivateAll(projectID string) error {
	_, err := database.DB.Exec(
		`UPDATE environments SET is_active = 0 WHERE project_id = ?`, projectID,
	)
	return err
}

func (r *EnvironmentRepo) Delete(id string) error {
	_, err := database.DB.Exec(`DELETE FROM environments WHERE id = ?`, id)
	return err
}
```

- [ ] **Step 5: Create `repositories/history_repo.go`**

```go
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
	var history []models.History
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
```

- [ ] **Step 6: Verify Go compiles**

Run: `cd D:\javap\paw-api && go build ./...`
Expected: No errors

---

### Task 3: Go — Create services layer

**Files:**
- Create: `services/project_service.go`
- Create: `services/collection_service.go`
- Create: `services/request_service.go`
- Create: `services/environment_service.go`
- Create: `services/history_service.go`

Each service wraps its repo and adds basic validation. Services return plain data or errors.

- [ ] **Step 1: Create `services/project_service.go`**

```go
package services

import (
	"errors"
	"time"
	"paw-api/models"
	"paw-api/repositories"

	"github.com/google/uuid"
)

type ProjectService struct {
	repo *repositories.ProjectRepo
}

func NewProjectService() *ProjectService {
	return &ProjectService{repo: &repositories.ProjectRepo{}}
}

func (s *ProjectService) Create(name, description string) (*models.Project, error) {
	if name == "" {
		return nil, errors.New("project name is required")
	}
	now := time.Now()
	p := &models.Project{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	return p, s.repo.Create(p)
}

func (s *ProjectService) GetByID(id string) (*models.Project, error) {
	return s.repo.GetByID(id)
}

func (s *ProjectService) List() ([]models.Project, error) {
	return s.repo.List()
}

func (s *ProjectService) Update(id, name, description string) (*models.Project, error) {
	if name == "" {
		return nil, errors.New("project name is required")
	}
	p, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errors.New("project not found")
	}
	p.Name = name
	p.Description = description
	p.UpdatedAt = time.Now()
	return p, s.repo.Update(p)
}

func (s *ProjectService) Delete(id string) error {
	return s.repo.Delete(id)
}
```

- [ ] **Step 2: Create `services/collection_service.go`**

```go
package services

import (
	"errors"
	"time"
	"paw-api/models"
	"paw-api/repositories"

	"github.com/google/uuid"
)

type CollectionService struct {
	repo *repositories.CollectionRepo
}

func NewCollectionService() *CollectionService {
	return &CollectionService{repo: &repositories.CollectionRepo{}}
}

func (s *CollectionService) Create(projectID, parentID, name string, sortOrder int) (*models.Collection, error) {
	if name == "" {
		return nil, errors.New("collection name is required")
	}
	now := time.Now()
	var pid *string
	if parentID != "" {
		pid = &parentID
	}
	c := &models.Collection{
		ID:        uuid.New().String(),
		ProjectID: projectID,
		ParentID:  pid,
		Name:      name,
		SortOrder: sortOrder,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return c, s.repo.Create(c)
}

func (s *CollectionService) GetByID(id string) (*models.Collection, error) {
	return s.repo.GetByID(id)
}

func (s *CollectionService) ListByProject(projectID string) ([]models.Collection, error) {
	return s.repo.ListByProject(projectID)
}

func (s *CollectionService) Update(id, name string, parentID *string, sortOrder int) (*models.Collection, error) {
	c, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, errors.New("collection not found")
	}
	c.Name = name
	c.ParentID = parentID
	c.SortOrder = sortOrder
	c.UpdatedAt = time.Now()
	return c, s.repo.Update(c)
}

func (s *CollectionService) Delete(id string) error {
	return s.repo.Delete(id)
}
```

- [ ] **Step 3: Create `services/request_service.go`**

```go
package services

import (
	"errors"
	"time"
	"paw-api/models"
	"paw-api/repositories"

	"github.com/google/uuid"
)

type RequestService struct {
	repo *repositories.RequestRepo
}

func NewRequestService() *RequestService {
	return &RequestService{repo: &repositories.RequestRepo{}}
}

func (s *RequestService) Create(collectionID, name, method, url, headers, params, body, auth, script string, sortOrder int) (*models.Request, error) {
	if name == "" {
		return nil, errors.New("request name is required")
	}
	now := time.Now()
	r := &models.Request{
		ID:           uuid.New().String(),
		CollectionID: collectionID,
		Name:         name,
		Method:       method,
		URL:          url,
		Headers:      headers,
		Params:       params,
		Body:         body,
		Auth:         auth,
		Script:       script,
		SortOrder:    sortOrder,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	return r, s.repo.Create(r)
}

func (s *RequestService) GetByID(id string) (*models.Request, error) {
	return s.repo.GetByID(id)
}

func (s *RequestService) ListByCollection(collectionID string) ([]models.Request, error) {
	return s.repo.ListByCollection(collectionID)
}

func (s *RequestService) Update(id, name, method, url, headers, params, body, auth, script string, sortOrder int) (*models.Request, error) {
	r, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, errors.New("request not found")
	}
	r.Name = name
	r.Method = method
	r.URL = url
	r.Headers = headers
	r.Params = params
	r.Body = body
	r.Auth = auth
	r.Script = script
	r.SortOrder = sortOrder
	r.UpdatedAt = time.Now()
	return r, s.repo.Update(r)
}

func (s *RequestService) Delete(id string) error {
	return s.repo.Delete(id)
}
```

- [ ] **Step 4: Create `services/environment_service.go`**

```go
package services

import (
	"errors"
	"time"
	"paw-api/models"
	"paw-api/repositories"

	"github.com/google/uuid"
)

type EnvironmentService struct {
	repo *repositories.EnvironmentRepo
}

func NewEnvironmentService() *EnvironmentService {
	return &EnvironmentService{repo: &repositories.EnvironmentRepo{}}
}

func (s *EnvironmentService) Create(projectID, name, variables string, isActive bool) (*models.Environment, error) {
	if name == "" {
		return nil, errors.New("environment name is required")
	}
	now := time.Now()
	e := &models.Environment{
		ID:        uuid.New().String(),
		ProjectID: projectID,
		Name:      name,
		Variables: variables,
		IsActive:  isActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return e, s.repo.Create(e)
}

func (s *EnvironmentService) GetByID(id string) (*models.Environment, error) {
	return s.repo.GetByID(id)
}

func (s *EnvironmentService) ListByProject(projectID string) ([]models.Environment, error) {
	return s.repo.ListByProject(projectID)
}

func (s *EnvironmentService) SetActive(id, projectID string) (*models.Environment, error) {
	if err := s.repo.DeactivateAll(projectID); err != nil {
		return nil, err
	}
	e, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, errors.New("environment not found")
	}
	e.IsActive = true
	e.UpdatedAt = time.Now()
	return e, s.repo.Update(e)
}
```

- [ ] **Step 5: Create `services/history_service.go`**

```go
package services

import (
	"time"
	"paw-api/models"
	"paw-api/repositories"

	"github.com/google/uuid"
)

type HistoryService struct {
	repo *repositories.HistoryRepo
}

func NewHistoryService() *HistoryService {
	return &HistoryService{repo: &repositories.HistoryRepo{}}
}

func (s *HistoryService) Record(projectID, method, url, headers, body string) (*models.History, error) {
	now := time.Now()
	h := &models.History{
		ID:        uuid.New().String(),
		ProjectID: projectID,
		Method:    method,
		URL:       url,
		Headers:   headers,
		Body:      body,
		CreatedAt: now,
	}
	return h, s.repo.Create(h)
}

func (s *HistoryService) ListByProject(projectID string, limit int) ([]models.History, error) {
	return s.repo.ListByProject(projectID, limit)
}

func (s *HistoryService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *HistoryService) ClearByProject(projectID string) error {
	return s.repo.ClearByProject(projectID)
}
```

- [ ] **Step 6: Add `google/uuid` dependency and verify compile**

Run: `cd D:\javap\paw-api && go get github.com/google/uuid && go build ./...`
Expected: No errors

---

### Task 4: Go — Create handlers and wire up App

**Files:**
- Create: `handlers/project_handler.go`
- Create: `handlers/collection_handler.go`
- Create: `handlers/request_handler.go`
- Create: `handlers/environment_handler.go`
- Create: `handlers/history_handler.go`
- Create: `pkg/httpclient/client.go`
- Modify: `app.go`
- Modify: `main.go`

Handlers are methods on `*App` that Wails binds directly. Each handler delegates to its service.

- [ ] **Step 1: Create `handlers/project_handler.go`**

```go
package handlers

import (
	"paw-api/models"
	"paw-api/services"
)

type ProjectHandler struct {
	service *services.ProjectService
}

func NewProjectHandler() *ProjectHandler {
	return &ProjectHandler{service: services.NewProjectService()}
}

func (h *ProjectHandler) CreateProject(name, description string) (*models.Project, error) {
	return h.service.Create(name, description)
}

func (h *ProjectHandler) GetProject(id string) (*models.Project, error) {
	return h.service.GetByID(id)
}

func (h *ProjectHandler) ListProjects() ([]models.Project, error) {
	return h.service.List()
}

func (h *ProjectHandler) UpdateProject(id, name, description string) (*models.Project, error) {
	return h.service.Update(id, name, description)
}

func (h *ProjectHandler) DeleteProject(id string) error {
	return h.service.Delete(id)
}
```

- [ ] **Step 2: Create `handlers/collection_handler.go`**

```go
package handlers

import (
	"paw-api/models"
	"paw-api/services"
)

type CollectionHandler struct {
	service *services.CollectionService
}

func NewCollectionHandler() *CollectionHandler {
	return &CollectionHandler{service: services.NewCollectionService()}
}

func (h *CollectionHandler) CreateCollection(projectID, parentID, name string, sortOrder int) (*models.Collection, error) {
	return h.service.Create(projectID, parentID, name, sortOrder)
}

func (h *CollectionHandler) GetCollection(id string) (*models.Collection, error) {
	return h.service.GetByID(id)
}

func (h *CollectionHandler) ListCollections(projectID string) ([]models.Collection, error) {
	return h.service.ListByProject(projectID)
}

func (h *CollectionHandler) UpdateCollection(id, name string, parentID *string, sortOrder int) (*models.Collection, error) {
	return h.service.Update(id, name, parentID, sortOrder)
}

func (h *CollectionHandler) DeleteCollection(id string) error {
	return h.service.Delete(id)
}
```

- [ ] **Step 3: Create `handlers/request_handler.go`**

```go
package handlers

import (
	"paw-api/models"
	"paw-api/services"
)

type RequestHandler struct {
	service *services.RequestService
}

func NewRequestHandler() *RequestHandler {
	return &RequestHandler{service: services.NewRequestService()}
}

func (h *RequestHandler) CreateRequest(collectionID, name, method, url, headers, params, body, auth, script string, sortOrder int) (*models.Request, error) {
	return h.service.Create(collectionID, name, method, url, headers, params, body, auth, script, sortOrder)
}

func (h *RequestHandler) GetRequest(id string) (*models.Request, error) {
	return h.service.GetByID(id)
}

func (h *RequestHandler) ListRequests(collectionID string) ([]models.Request, error) {
	return h.service.ListByCollection(collectionID)
}

func (h *RequestHandler) UpdateRequest(id, name, method, url, headers, params, body, auth, script string, sortOrder int) (*models.Request, error) {
	return h.service.Update(id, name, method, url, headers, params, body, auth, script, sortOrder)
}

func (h *RequestHandler) DeleteRequest(id string) error {
	return h.service.Delete(id)
}
```

- [ ] **Step 4: Create `handlers/environment_handler.go`**

```go
package handlers

import (
	"paw-api/models"
	"paw-api/services"
)

type EnvironmentHandler struct {
	service *services.EnvironmentService
}

func NewEnvironmentHandler() *EnvironmentHandler {
	return &EnvironmentHandler{service: services.NewEnvironmentService()}
}

func (h *EnvironmentHandler) CreateEnvironment(projectID, name, variables string, isActive bool) (*models.Environment, error) {
	return h.service.Create(projectID, name, variables, isActive)
}

func (h *EnvironmentHandler) GetEnvironment(id string) (*models.Environment, error) {
	return h.service.GetByID(id)
}

func (h *EnvironmentHandler) ListEnvironments(projectID string) ([]models.Environment, error) {
	return h.service.ListByProject(projectID)
}

func (h *EnvironmentHandler) SetActiveEnvironment(id, projectID string) (*models.Environment, error) {
	return h.service.SetActive(id, projectID)
}
```

- [ ] **Step 5: Create `handlers/history_handler.go`**

```go
package handlers

import (
	"paw-api/models"
	"paw-api/services"
)

type HistoryHandler struct {
	service *services.HistoryService
}

func NewHistoryHandler() *HistoryHandler {
	return &HistoryHandler{service: services.NewHistoryService()}
}

func (h *HistoryHandler) RecordHistory(projectID, method, url, headers, body string) (*models.History, error) {
	return h.service.Record(projectID, method, url, headers, body)
}

func (h *HistoryHandler) ListHistory(projectID string, limit int) ([]models.History, error) {
	return h.service.ListByProject(projectID, limit)
}

func (h *HistoryHandler) DeleteHistory(id string) error {
	return h.service.Delete(id)
}

func (h *HistoryHandler) ClearHistory(projectID string) error {
	return h.service.ClearByProject(projectID)
}
```

- [ ] **Step 6: Create `pkg/httpclient/client.go`**

```go
package httpclient

type Client struct{}

func NewClient() *Client {
	return &Client{}
}
```

- [ ] **Step 7: Rewrite `app.go`**

```go
package main

import (
	"context"
	"paw-api/database"
	"paw-api/handlers"
)

type App struct {
	ctx                context.Context
	ProjectHandler     *handlers.ProjectHandler
	CollectionHandler  *handlers.CollectionHandler
	RequestHandler     *handlers.RequestHandler
	EnvironmentHandler *handlers.EnvironmentHandler
	HistoryHandler     *handlers.HistoryHandler
}

func NewApp() *App {
	return &App{
		ProjectHandler:     handlers.NewProjectHandler(),
		CollectionHandler:  handlers.NewCollectionHandler(),
		RequestHandler:     handlers.NewRequestHandler(),
		EnvironmentHandler: handlers.NewEnvironmentHandler(),
		HistoryHandler:     handlers.NewHistoryHandler(),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	if err := database.InitDB(); err != nil {
		panic("failed to initialize database: " + err.Error())
	}
}

func (a *App) shutdown(ctx context.Context) {
	database.CloseDB()
}
```

- [ ] **Step 8: Update `main.go` to bind all handlers and register shutdown**

```go
package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:     "Paw API",
		Width:     1280,
		Height:    800,
		MinWidth:  960,
		MinHeight: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind: []interface{}{
			app,
			app.ProjectHandler,
			app.CollectionHandler,
			app.RequestHandler,
			app.EnvironmentHandler,
			app.HistoryHandler,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
```

- [ ] **Step 9: Verify Go compiles**

Run: `cd D:\javap\paw-api && go build ./...`
Expected: No errors

---

### Task 5: Frontend — Install dependencies

**Files:**
- Modify: `frontend/package.json`

- [ ] **Step 1: Install Naive UI, Pinia, vue-router, and icons**

Run: `cd D:\javap\paw-api\frontend && npm install naive-ui pinia vue-router@4 @vicons/ionicons5`

- [ ] **Step 2: Verify frontend install**

Run: `cd D:\javap\paw-api\frontend && npm ls naive-ui pinia vue-router @vicons/ionicons5`
Expected: All listed with version numbers, no errors

---

### Task 6: Frontend — Types, stores, and composables

**Files:**
- Create: `frontend/src/types/project.ts`
- Create: `frontend/src/types/request.ts`
- Create: `frontend/src/types/environment.ts`
- Create: `frontend/src/stores/project.ts`
- Create: `frontend/src/stores/request.ts`
- Create: `frontend/src/stores/environment.ts`
- Create: `frontend/src/stores/history.ts`
- Create: `frontend/src/stores/tabs.ts`
- Create: `frontend/src/composables/useTheme.ts`
- Create: `frontend/src/composables/useRequest.ts`
- Create: `frontend/src/composables/useCollection.ts`

- [ ] **Step 1: Create `frontend/src/types/project.ts`**

```typescript
export interface Project {
  id: string
  name: string
  description: string
  created_at: string
  updated_at: string
}

export interface Collection {
  id: string
  project_id: string
  parent_id: string | null
  name: string
  sort_order: number
  created_at: string
  updated_at: string
}
```

- [ ] **Step 2: Create `frontend/src/types/request.ts`**

```typescript
export interface Request {
  id: string
  collection_id: string
  name: string
  method: string
  url: string
  headers: string
  params: string
  body: string
  auth: string
  script: string
  sort_order: number
  created_at: string
  updated_at: string
}
```

- [ ] **Step 3: Create `frontend/src/types/environment.ts`**

```typescript
export interface Environment {
  id: string
  project_id: string
  name: string
  variables: string
  is_active: boolean
  created_at: string
  updated_at: string
}
```

- [ ] **Step 4: Create `frontend/src/stores/project.ts`**

```typescript
import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Project, Collection } from '../types/project'

export const useProjectStore = defineStore('project', () => {
  const projects = ref<Project[]>([])
  const currentProject = ref<Project | null>(null)
  const collections = ref<Collection[]>([])

  function setCurrentProject(p: Project) {
    currentProject.value = p
  }

  return { projects, currentProject, collections, setCurrentProject }
})
```

- [ ] **Step 5: Create `frontend/src/stores/request.ts`**

```typescript
import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Request } from '../types/request'

export const useRequestStore = defineStore('request', () => {
  const requests = ref<Request[]>([])
  const currentRequest = ref<Request | null>(null)

  function setCurrentRequest(r: Request) {
    currentRequest.value = r
  }

  return { requests, currentRequest, setCurrentRequest }
})
```

- [ ] **Step 6: Create `frontend/src/stores/environment.ts`**

```typescript
import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Environment } from '../types/environment'

export const useEnvironmentStore = defineStore('environment', () => {
  const environments = ref<Environment[]>([])
  const activeEnvironment = ref<Environment | null>(null)

  function setActiveEnvironment(env: Environment) {
    activeEnvironment.value = env
  }

  return { environments, activeEnvironment, setActiveEnvironment }
})
```

- [ ] **Step 7: Create `frontend/src/stores/history.ts`**

```typescript
import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface HistoryItem {
  id: string
  project_id: string
  request_id: string | null
  method: string
  url: string
  headers: string
  body: string
  response_status: number
  response_body: string
  response_headers: string
  duration_ms: number
  created_at: string
}

export const useHistoryStore = defineStore('history', () => {
  const history = ref<HistoryItem[]>([])

  function setHistory(items: HistoryItem[]) {
    history.value = items
  }

  return { history, setHistory }
})
```

- [ ] **Step 8: Create `frontend/src/stores/tabs.ts`**

```typescript
import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface Tab {
  id: string
  title: string
  requestId: string
  active: boolean
}

export const useTabsStore = defineStore('tabs', () => {
  const tabs = ref<Tab[]>([])
  const activeTabId = ref<string | null>(null)

  function addTab(tab: Tab) {
    tabs.value.push(tab)
    activeTabId.value = tab.id
  }

  function removeTab(tabId: string) {
    const idx = tabs.value.findIndex(t => t.id === tabId)
    if (idx !== -1) {
      tabs.value.splice(idx, 1)
      if (activeTabId.value === tabId) {
        activeTabId.value = tabs.value.length > 0 ? tabs.value[Math.min(idx, tabs.value.length - 1)].id : null
      }
    }
  }

  function setActiveTab(tabId: string) {
    activeTabId.value = tabId
  }

  return { tabs, activeTabId, addTab, removeTab, setActiveTab }
})
```

- [ ] **Step 9: Create `frontend/src/composables/useTheme.ts`**

```typescript
import { computed, ref, watch } from 'vue'
import { darkTheme } from 'naive-ui'
import type { GlobalTheme, GlobalThemeOverrides } from 'naive-ui'

type ColorMode = 'light' | 'dark'
type ThemeColor = 'green' | 'blue' | 'purple'

const themeOverrides: Record<ThemeColor, GlobalThemeOverrides> = {
  green: {
    common: { primaryColor: '#18a058', primaryColorHover: '#36ad6a', primaryColorPressed: '#0c7a43' }
  },
  blue: {
    common: { primaryColor: '#2080f0', primaryColorHover: '#4098f7', primaryColorPressed: '#1060c0' }
  },
  purple: {
    common: { primaryColor: '#8a63d2', primaryColorHover: '#b794f4', primaryColorPressed: '#6b46c0' }
  }
}

const colorMode = ref<ColorMode>('light')
const themeColor = ref<ThemeColor>('green')

export function useTheme() {
  const theme = computed<GlobalTheme | null>(() => colorMode.value === 'dark' ? darkTheme : null)

  const themeOverridesComputed = computed<GlobalThemeOverrides>(() => themeOverrides[themeColor.value])

  function toggleColorMode() {
    colorMode.value = colorMode.value === 'light' ? 'dark' : 'light'
  }

  function setThemeColor(color: ThemeColor) {
    themeColor.value = color
  }

  function applyBodyTheme(mode: ColorMode) {
    document.body.setAttribute('data-naive-ui-theme', mode)
  }

  watch(colorMode, (mode) => {
    localStorage.setItem('paw-color-mode', mode)
    applyBodyTheme(mode)
  })

  watch(themeColor, (color) => {
    localStorage.setItem('paw-theme-color', color)
  })

  applyBodyTheme(colorMode.value)

  return { theme, themeOverrides: themeOverridesComputed, colorMode, themeColor, toggleColorMode, setThemeColor }
}
```

- [ ] **Step 10: Create `frontend/src/composables/useRequest.ts`**

```typescript
export function useRequest() {
  function parseParams(url: string): Record<string, string> {
    const idx = url.indexOf('?')
    if (idx === -1) return {}
    const params: Record<string, string> = {}
    new URLSearchParams(url.slice(idx)).forEach((v, k) => { params[k] = v })
    return params
  }

  return { parseParams }
}
```

- [ ] **Step 11: Create `frontend/src/composables/useCollection.ts`**

```typescript
export function useCollection() {
  function flattenTree(items: Array<{ id: string; parent_id: string | null }>): string[] {
    const ids: string[] = []
    function walk(items: Array<{ id: string; parent_id: string | null }>, parentId: string | null) {
      for (const item of items) {
        if (item.parent_id === parentId) {
          ids.push(item.id)
          walk(items, item.id)
        }
      }
    }
    walk(items, null)
    return ids
  }

  return { flattenTree }
}
```

---

### Task 7: Frontend — Router, views, layout components, and App

**Files:**
- Create: `frontend/src/router/index.ts`
- Create: `frontend/src/views/Workspace.vue`
- Create: `frontend/src/views/History.vue`
- Create: `frontend/src/views/Docs.vue`
- Create: `frontend/src/views/TestRunner.vue`
- Create: `frontend/src/components/AppSidebar.vue`
- Create: `frontend/src/components/TabBar.vue`
- Create: `frontend/src/components/RequestEditor.vue`
- Create: `frontend/src/components/ResponseViewer.vue`
- Modify: `frontend/src/main.ts`
- Modify: `frontend/src/App.vue`
- Modify: `frontend/src/style.css`
- Delete: `frontend/src/components/HelloWorld.vue`
- Delete: `frontend/src/assets/images/logo-universal.png` (optional)

- [ ] **Step 1: Create `frontend/src/router/index.ts`**

```typescript
import { createRouter, createWebHashHistory } from 'vue-router'
import Workspace from '../views/Workspace.vue'

const routes = [
  { path: '/', redirect: '/workspace' },
  { path: '/workspace', name: 'Workspace', component: Workspace },
  { path: '/history', name: 'History', component: () => import('../views/History.vue') },
  { path: '/docs', name: 'Docs', component: () => import('../views/Docs.vue') },
  { path: '/test-runner', name: 'TestRunner', component: () => import('../views/TestRunner.vue') },
]

export const router = createRouter({
  history: createWebHashHistory(),
  routes,
})
```

- [ ] **Step 2: Create `frontend/src/views/Workspace.vue`**

```vue
<script lang="ts" setup>
</script>

<template>
  <div class="workspace">
    <TabBar />
    <div class="workspace-body">
      <RequestEditor />
      <ResponseViewer />
    </div>
  </div>
</template>

<style scoped>
.workspace {
  display: flex;
  flex-direction: column;
  height: 100%;
}
.workspace-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
</style>
```

- [ ] **Step 3: Create `frontend/src/views/History.vue`**

```vue
<script lang="ts" setup>
</script>

<template>
  <div class="history-view">
    <h2>History</h2>
    <p class="subtitle">Request history coming soon</p>
  </div>
</template>

<style scoped>
.history-view {
  padding: 24px;
}
.subtitle {
  color: #888;
  font-size: 14px;
}
</style>
```

- [ ] **Step 4: Create `frontend/src/views/Docs.vue`**

```vue
<script lang="ts" setup>
</script>

<template>
  <div class="docs-view">
    <h2>API Docs</h2>
    <p class="subtitle">Documentation coming soon</p>
  </div>
</template>

<style scoped>
.docs-view {
  padding: 24px;
}
.subtitle {
  color: #888;
  font-size: 14px;
}
</style>
```

- [ ] **Step 5: Create `frontend/src/views/TestRunner.vue`**

```vue
<script lang="ts" setup>
</script>

<template>
  <div class="test-runner-view">
    <h2>Test Runner</h2>
    <p class="subtitle">Automated testing coming soon</p>
  </div>
</template>

<style scoped>
.test-runner-view {
  padding: 24px;
}
.subtitle {
  color: #888;
  font-size: 14px;
}
</style>
```

- [ ] **Step 6: Create `frontend/src/components/AppSidebar.vue`**

```vue
<script lang="ts" setup>
import { NLayoutSider } from 'naive-ui'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

const menuItems = [
  { label: 'Workspace', key: '/workspace', icon: '▣' },
  { label: 'History', key: '/history', icon: '◷' },
  { label: 'Docs', key: '/docs', icon: '◰' },
  { label: 'Test Runner', key: '/test-runner', icon: '▶' },
]

function navigateTo(path: string) {
  router.push(path)
}
</script>

<template>
  <NLayoutSider
    bordered
    width="220"
    :native-scrollbar="false"
    class="app-sidebar"
  >
    <div class="sidebar-header">
      <span class="sidebar-title">Paw API</span>
    </div>
    <div class="sidebar-menu">
      <div
        v-for="item in menuItems"
        :key="item.key"
        class="sidebar-menu-item"
        :class="{ active: route.path === item.key }"
        @click="navigateTo(item.key)"
      >
        <span class="menu-icon">{{ item.icon }}</span>
        <span>{{ item.label }}</span>
      </div>
    </div>
  </NLayoutSider>
</template>

<style scoped>
.app-sidebar {
  height: 100%;
}
.sidebar-header {
  padding: 16px 16px 12px;
  border-bottom: 1px solid var(--border-color);
}
.sidebar-title {
  font-size: 16px;
  font-weight: 700;
}
.sidebar-menu {
  padding: 8px;
}
.sidebar-menu-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: background 0.15s;
}
.sidebar-menu-item:hover {
  background: var(--hover-color);
}
.sidebar-menu-item.active {
  background: var(--active-color);
  font-weight: 600;
}
.menu-icon {
  font-size: 16px;
  width: 20px;
  text-align: center;
}
</style>
```

- [ ] **Step 7: Create `frontend/src/components/TabBar.vue`**

```vue
<script lang="ts" setup>
import { NTabs } from 'naive-ui'
</script>

<template>
  <div class="tab-bar">
    <span class="tab-placeholder">Tabs coming in Phase 4</span>
  </div>
</template>

<style scoped>
.tab-bar {
  display: flex;
  align-items: center;
  padding: 4px 8px;
  border-bottom: 1px solid var(--border-color);
  background: var(--tab-bar-bg);
  min-height: 36px;
}
.tab-placeholder {
  font-size: 12px;
  color: #999;
  padding: 0 8px;
}
</style>
```

- [ ] **Step 8: Create `frontend/src/components/RequestEditor.vue`**

```vue
<script lang="ts" setup>
import { NInput, NSelect, NButton, NTabs, NTabPane } from 'naive-ui'

const httpMethods = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'HEAD', 'OPTIONS']
</script>

<template>
  <div class="request-editor">
    <div class="url-row">
      <NSelect
        :options="httpMethods.map(m => ({ label: m, value: m }))"
        value="GET"
        style="width: 110px"
        size="small"
      />
      <NInput
        placeholder="https://api.example.com/endpoint"
        size="small"
        class="url-input"
      />
      <NButton type="primary" size="small">
        Send
      </NButton>
    </div>
    <NTabs type="line" size="small" class="editor-tabs">
      <NTabPane name="params" tab="Params">
        <div class="tab-placeholder">Key-value editor coming soon</div>
      </NTabPane>
      <NTabPane name="headers" tab="Headers">
        <div class="tab-placeholder">Key-value editor coming soon</div>
      </NTabPane>
      <NTabPane name="body" tab="Body">
        <div class="tab-placeholder">Body editor coming soon</div>
      </NTabPane>
      <NTabPane name="auth" tab="Auth">
        <div class="tab-placeholder">Auth config coming soon</div>
      </NTabPane>
    </NTabs>
  </div>
</template>

<style scoped>
.request-editor {
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
}
.url-row {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
}
.url-input {
  flex: 1;
}
.editor-tabs {
  margin-top: 4px;
}
.tab-placeholder {
  padding: 16px 0;
  color: #999;
  font-size: 13px;
}
</style>
```

- [ ] **Step 9: Create `frontend/src/components/ResponseViewer.vue`**

```vue
<script lang="ts" setup>
</script>

<template>
  <div class="response-viewer">
    <div class="response-header">
      <span class="response-status">Waiting for request...</span>
    </div>
    <div class="response-body">
      <div class="response-placeholder">Send a request to see the response</div>
    </div>
  </div>
</template>

<style scoped>
.response-viewer {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.response-header {
  padding: 8px 16px;
  border-bottom: 1px solid var(--border-color);
  font-size: 13px;
}
.response-status {
  color: #999;
}
.response-body {
  flex: 1;
  padding: 16px;
  overflow: auto;
}
.response-placeholder {
  color: #999;
  font-size: 13px;
  text-align: center;
  margin-top: 48px;
}
</style>
```

- [ ] **Step 10: Rewrite `frontend/src/main.ts`**

```typescript
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import naive from 'naive-ui'
import App from './App.vue'
import { router } from './router'
import './style.css'

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(naive)
app.mount('#app')
```

- [ ] **Step 11: Rewrite `frontend/src/App.vue`**

```vue
<script lang="ts" setup>
import { NConfigProvider, NLayout, NLayoutSider, NLayoutContent, NMessageProvider, NDialogProvider } from 'naive-ui'
import { useTheme } from './composables/useTheme'
import AppSidebar from './components/AppSidebar.vue'

const { theme, themeOverrides } = useTheme()
</script>

<template>
  <NConfigProvider :theme="theme" :theme-overrides="themeOverrides">
    <NMessageProvider>
      <NDialogProvider>
        <NLayout class="app-layout" has-sider>
          <AppSidebar />
          <NLayoutContent class="main-content">
            <router-view />
          </NLayoutContent>
        </NLayout>
      </NDialogProvider>
    </NMessageProvider>
  </NConfigProvider>
</template>

<style scoped>
.app-layout {
  height: 100vh;
  width: 100vw;
}
.main-content {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
}
</style>
```

- [ ] **Step 12: Rewrite `frontend/src/style.css`**

```css
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body, #app {
  height: 100%;
  width: 100%;
  overflow: hidden;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue", sans-serif;
  -webkit-font-smoothing: antialiased;
}

:root {
  --border-color: #e0e0e0;
  --hover-color: #f5f5f5;
  --active-color: #e8f5e9;
  --tab-bar-bg: #fafafa;
}

body[data-naive-ui-theme="dark"] {
  --border-color: #333;
  --hover-color: #2a2a2a;
  --active-color: #1b3a2a;
  --tab-bar-bg: #1e1e1e;
}
```

- [ ] **Step 13: Clean up template files**

```bash
Remove-Item -LiteralPath "D:\javap\paw-api\frontend\src\components\HelloWorld.vue" -Force
Remove-Item -LiteralPath "D:\javap\paw-api\frontend\src\assets\images\logo-universal.png" -Force
```

---

### Task 8: Verify the build

- [ ] **Step 1: Verify frontend type-check**

Run: `cd D:\javap\paw-api\frontend && npm run build`
Expected: `vue-tsc --noEmit` passes, then `vite build` produces `frontend/dist/`

- [ ] **Step 2: Verify Go build**

Run: `cd D:\javap\paw-api && go build ./...`
Expected: No errors

- [ ] **Step 3: Verify Wails build**

Run: `cd D:\javap\paw-api && wails build`
Expected: Builds successfully, outputs to `build/bin/`
