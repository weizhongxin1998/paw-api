package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

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
	importH      *handlers.ImportHandler
	exportH      *handlers.ExportHandler
	variableH    *handlers.VariableHandler
	wsH          *handlers.WebSocketHandler
	docsH        *handlers.DocsHandler
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
	requestSvc := services.NewRequestService(requestRepo, historyRepo, envRepo, a.snowflake, a.httpClient)
	collectionSvc := services.NewCollectionService(collectionRepo, requestRepo, a.snowflake)
	projectSvc := services.NewProjectService(projectRepo, a.snowflake)

	// Wire project service to also create default collection
	_ = projectSvc
	_ = collectionSvc

	importSvc := services.NewImportService(projectRepo, collectionRepo, requestRepo, a.snowflake)
	exportSvc := services.NewExportService(projectRepo, collectionRepo, requestRepo, a.snowflake)

	variableSvc := services.NewVariableService(varRepo)

	a.settingsH = handlers.NewSettingsHandler(settingsSvc)
	a.historyH = handlers.NewHistoryHandler(historySvc)
	a.environmentH = handlers.NewEnvironmentHandler(envSvc)
	a.requestH = handlers.NewRequestHandler(requestSvc)
	a.collectionH = handlers.NewCollectionHandler(collectionSvc)
	a.projectH = handlers.NewProjectHandler(projectSvc)
	a.importH = handlers.NewImportHandler(importSvc)
	a.exportH = handlers.NewExportHandler(exportSvc)
	a.variableH = handlers.NewVariableHandler(variableSvc)

	a.wsH = handlers.NewWebSocketHandler()

	docsSvc := services.NewDocsService(collectionRepo, requestRepo)
	a.docsH = handlers.NewDocsHandler(docsSvc)

	// Set context for WebSocket handler events
	a.wsH.SetContext(ctx)

	// Load settings into http client
	settings, err := settingsSvc.GetAll()
	if err == nil {
		a.httpClient.ApplySettings(settings)
	}

	// Ensure demo data on first run
	a.seedDemoData()
}

func (a *App) seedDemoData() {
	// Dev mode (wails.json exists in cwd): always ensure demo project exists
	// Production mode: only seed on first run (no projects)
	_, isDev := os.Stat("wails.json")
	if isDev == nil {
		a.ensureDemoProject()
		return
	}
	// Production: only on first launch
	projects, err := a.projectH.List()
	if err != nil || len(projects) > 0 {
		return
	}
	fmt.Println("seeding demo data (first launch)...")
	a.ensureDemoProject()
}

func (a *App) ensureDemoProject() {
	projects, err := a.projectH.List()
	if err != nil {
		return
	}
	for _, p := range projects {
		if p.Name == "演示项目" {
			return
		}
	}
	// Reuse the same seeding logic by temporarily setting up
	// We call CreateProject which also creates "内置环境"
	project, err := a.CreateProject("演示项目", "内置演示项目，包含用户和订单模块")
	if err != nil {
		fmt.Printf("seed: failed to create project: %v\n", err)
		return
	}
	a.seedRequests(project)
}

func (a *App) seedRequests(project *models.Project) {
	userColl, err := a.collectionH.Create(project.ID, nil, "用户模块")
	if err != nil {
		return
	}
	orderColl, err := a.collectionH.Create(project.ID, nil, "订单模块")
	if err != nil {
		return
	}
	healthColl, err := a.collectionH.Create(project.ID, nil, "通用")
	if err != nil {
		return
	}

	type reqData struct {
		name     string
		method   string
		url      string
		params   string
		headers  string
		bodyType string
		body     string
	}

	createReq := func(collectionID int64, d reqData) {
		r, err := a.requestH.Create(collectionID, d.name, d.method)
		if err != nil {
			return
		}
		headers := d.headers
		if headers == "" {
			headers = `[{"key":"Accept","value":"application/json","enabled":true}]`
		}
		params := d.params
		if params == "" {
			params = `[]`
		}
		r.URL = d.url
		r.Headers = headers
		r.Params = params
		r.BodyType = d.bodyType
		r.Body = d.body
		a.requestH.Update(r)
	}

	jsonHdr := `[{"key":"Content-Type","value":"application/json","enabled":true}]`

	for _, d := range []reqData{
		{"用户列表", "GET", "/api/users", `[{"key":"page","value":"1","enabled":true,"description":"页码"},{"key":"size","value":"20","enabled":true,"description":"每页条数"}]`, "", "", ""},
		{"用户详情", "GET", "/api/users/{id}", "", "", "", ""},
		{"创建用户", "POST", "/api/users", "", jsonHdr, "raw", `{"name":"张三","email":"zhangsan@example.com","age":28,"role":"user"}`},
		{"更新用户", "PUT", "/api/users/{id}", "", jsonHdr, "raw", `{"name":"张三","email":"zhangsan@example.com","age":29,"role":"admin"}`},
		{"部分更新用户", "PATCH", "/api/users/{id}", "", jsonHdr, "raw", `{"name":"张三(已修改)"}`},
		{"删除用户", "DELETE", "/api/users/{id}", "", "", "", ""},
	} {
		createReq(userColl.ID, d)
	}

	for _, d := range []reqData{
		{"订单列表", "GET", "/api/orders", "", "", "", ""},
		{"订单详情", "GET", "/api/orders/{id}", "", "", "", ""},
		{"创建订单", "POST", "/api/orders", "", jsonHdr, "raw", `{"user_id":1,"items":[{"name":"Widget A","price":19.99,"quantity":2},{"name":"Widget B","price":9.99,"quantity":1}]}`},
		{"取消订单", "DELETE", "/api/orders/{id}", "", "", "", ""},
	} {
		createReq(orderColl.ID, d)
	}

	for _, d := range []reqData{
		{"健康检查", "GET", "/api/health", "", "", "", ""},
		{"健康检查(HEAD)", "HEAD", "/api/health", "", "", "", ""},
		{"查看支持方法", "OPTIONS", "/api/users", "", "", "", ""},
	} {
		createReq(healthColl.ID, d)
	}

	fmt.Println("demo project ensured successfully")
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

func (a *App) GetProjectStats(id int64) (models.ProjectStats, error) {
	return a.projectH.GetStats(id)
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

func (a *App) CreateEnvironment(projectID int64, name string, baseURL string, cloneFromID *int64) (*models.Environment, error) {
	return a.environmentH.Create(projectID, name, baseURL, cloneFromID)
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

func (a *App) SaveEnvBaseURL(envID int64, baseURL string) error {
	return a.environmentH.SaveBaseURL(envID, baseURL)
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

// ========== Import Handlers ==========

func (a *App) ImportPostman(projectID int64, filePath string) (*services.ImportResult, error) {
	return a.importH.ImportPostman(projectID, filePath)
}

// ========== Export Handlers ==========

func (a *App) ExportPostman(projectID int64) (string, error) {
	return a.exportH.ExportPostman(projectID)
}

// ========== Variable Handlers ==========

func (a *App) ResolveVariable(text string, envID int64) (string, error) {
	return a.variableH.Resolve(text, envID)
}

func (a *App) ResolveVariableMap(m map[string]string, envID int64) (map[string]string, error) {
	return a.variableH.ResolveMap(m, envID)
}

// ========== WebSocket Handlers ==========

func (a *App) WSConnect(url string, headersJSON string) error {
	return a.wsH.Connect(url, headersJSON)
}

func (a *App) WSSend(url string, message string) error {
	return a.wsH.Send(url, message)
}

func (a *App) WSDisconnect(url string) error {
	return a.wsH.Disconnect(url)
}

// ========== Docs Handlers ==========

func (a *App) GenerateDocsMarkdown(projectID int64) (string, error) {
	return a.docsH.GenerateMarkdown(projectID)
}

func (a *App) GenerateDocsHTML(projectID int64) (string, error) {
	return a.docsH.GenerateHTML(projectID)
}

func (a *App) GenerateRequestDocsMarkdown(requestID int64) (string, error) {
	return a.docsH.GenerateRequestMarkdown(requestID)
}

func (a *App) GenerateRequestDocsHTML(requestID int64) (string, error) {
	return a.docsH.GenerateRequestHTML(requestID)
}
