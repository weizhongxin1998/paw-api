package main

import (
	"context"
	"database/sql"

	"paw-api/database"
	"paw-api/handlers"
	"paw-api/models"
	"paw-api/pkg/httpclient"
	"paw-api/pkg/snowflake"
	"paw-api/repositories"
	"paw-api/services"
)

type App struct {
	ctx       context.Context
	db        *sql.DB
	snowflake *snowflake.Generator
	httpClient *httpclient.Client

	projectH     *handlers.ProjectHandler
	collectionH  *handlers.CollectionHandler
	requestH     *handlers.RequestHandler
	environmentH *handlers.EnvironmentHandler
	historyH     *handlers.HistoryHandler
	settingsH    *handlers.SettingsHandler
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	var err error
	a.db, err = database.Init()
	if err != nil {
		panic(err)
	}

	a.snowflake = snowflake.New()
	a.httpClient = httpclient.New()

	projectRepo := repositories.NewProjectRepo(a.db, a.snowflake)
	collectionRepo := repositories.NewCollectionRepo(a.db, a.snowflake)
	requestRepo := repositories.NewRequestRepo(a.db, a.snowflake)
	envRepo := repositories.NewEnvironmentRepo(a.db, a.snowflake)
	varRepo := repositories.NewEnvVariableRepo(a.db, a.snowflake)
	historyRepo := repositories.NewHistoryRepo(a.db, a.snowflake)
	settingsRepo := repositories.NewSettingsRepo(a.db)

	settingsSvc := services.NewSettingsService(settingsRepo)

	historySvc := services.NewHistoryService(historyRepo)
	envSvc := services.NewEnvironmentService(envRepo, varRepo, a.snowflake)
	requestSvc := services.NewRequestService(requestRepo, historyRepo, a.snowflake, a.httpClient)
	collectionSvc := services.NewCollectionService(collectionRepo, requestRepo, a.snowflake)
	projectSvc := services.NewProjectService(projectRepo, a.snowflake)

	// Wire project service to also create default collection
	_ = projectSvc
	_ = collectionSvc

	a.settingsH = handlers.NewSettingsHandler(settingsSvc)
	a.historyH = handlers.NewHistoryHandler(historySvc)
	a.environmentH = handlers.NewEnvironmentHandler(envSvc)
	a.requestH = handlers.NewRequestHandler(requestSvc)
	a.collectionH = handlers.NewCollectionHandler(collectionSvc)
	a.projectH = handlers.NewProjectHandler(projectSvc)

	// Load settings into http client
	settings, err := settingsSvc.GetAll()
	if err == nil {
		a.httpClient.ApplySettings(settings)
	}
}

// ========== Project Handlers ==========

func (a *App) ListProjects() ([]models.Project, error) {
	return a.projectH.List()
}

func (a *App) GetProject(id int64) (*models.Project, error) {
	return a.projectH.Get(id)
}

func (a *App) CreateProject(name, description string) (*models.Project, error) {
	return a.projectH.Create(name, description)
}

func (a *App) UpdateProject(id int64, name, description string) (*models.Project, error) {
	return a.projectH.Update(id, name, description)
}

func (a *App) DeleteProject(id int64) error {
	return a.projectH.Delete(id)
}

// ========== Collection Handlers ==========

func (a *App) GetCollectionTree(projectID int64) ([]models.TreeItem, error) {
	return a.collectionH.GetTree(projectID)
}

func (a *App) CreateCollection(projectID int64, parentID *int64, name string) (*models.Collection, error) {
	return a.collectionH.Create(projectID, parentID, name)
}

func (a *App) RenameCollection(id int64, name string) error {
	return a.collectionH.Rename(id, name)
}

func (a *App) MoveCollection(id int64, parentID *int64, sortOrder int) error {
	return a.collectionH.Move(id, parentID, sortOrder)
}

func (a *App) DeleteCollection(id int64) error {
	return a.collectionH.Delete(id)
}

// ========== Request Handlers ==========

func (a *App) GetRequest(id int64) (*models.Request, error) {
	return a.requestH.Get(id)
}

func (a *App) CreateRequest(collectionID int64, name, method string) (*models.Request, error) {
	return a.requestH.Create(collectionID, name, method)
}

func (a *App) UpdateRequest(req *models.Request) error {
	return a.requestH.Update(req)
}

func (a *App) CloneRequest(id int64) (*models.Request, error) {
	return a.requestH.Clone(id)
}

func (a *App) DeleteRequest(id int64) error {
	return a.requestH.Delete(id)
}

func (a *App) SendRequest(sessionID int64, req *models.Request, envID int64) (*models.HTTPResponse, error) {
	return a.requestH.SendRequest(sessionID, req, envID)
}

func (a *App) SendQuickRequest(sessionID int64, method, url, headers, body string, envID int64) (*models.HTTPResponse, error) {
	return a.requestH.SendQuick(sessionID, method, url, headers, body, envID)
}

func (a *App) CancelRequest(sessionID int64) {
	a.requestH.CancelRequest(sessionID)
}

// ========== Environment Handlers ==========

func (a *App) ListEnvironments(projectID int64) ([]models.Environment, error) {
	return a.environmentH.List(projectID)
}

func (a *App) CreateEnvironment(projectID int64, name string, cloneFromID *int64) (*models.Environment, error) {
	return a.environmentH.Create(projectID, name, cloneFromID)
}

func (a *App) RenameEnvironment(id int64, name string) error {
	return a.environmentH.Rename(id, name)
}

func (a *App) DeleteEnvironment(id int64) error {
	return a.environmentH.Delete(id)
}

func (a *App) ActivateEnvironment(id int64) error {
	return a.environmentH.Activate(id)
}

func (a *App) ListEnvVariables(envID int64) ([]models.EnvVariable, error) {
	return a.environmentH.ListVariables(envID)
}

func (a *App) SaveEnvVariables(envID int64, variables []models.EnvVariable) error {
	return a.environmentH.SaveVariables(envID, variables)
}

// ========== History Handlers ==========

func (a *App) ListHistory(projectID int64, page, pageSize int) ([]models.History, error) {
	return a.historyH.List(projectID, page, pageSize)
}

func (a *App) SearchHistory(projectID int64, keyword, method string, statusMin, statusMax int, page, pageSize int) ([]models.History, int, error) {
	return a.historyH.Search(projectID, keyword, method, statusMin, statusMax, page, pageSize)
}

func (a *App) ClearHistory(projectID int64) error {
	return a.historyH.Clear(projectID)
}

func (a *App) DeleteHistory(id int64) error {
	return a.historyH.Delete(id)
}

// ========== Settings Handlers ==========

func (a *App) GetSetting(key string) (string, error) {
	return a.settingsH.Get(key)
}

func (a *App) SetSetting(key, value string) error {
	return a.settingsH.Set(key, value)
}

func (a *App) GetAllSettings() (map[string]string, error) {
	return a.settingsH.GetAll()
}
