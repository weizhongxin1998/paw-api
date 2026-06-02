# Paw API — Go 后端设计

> 版本 1.0 | 2026-06-02

---

## 1. 包结构

```
paw-api/
├── main.go
├── app.go                       # Wails App 结构体，绑定注册
├── go.mod / go.sum
├── wails.json
│
├── models/                      # 数据模型（纯结构体，无逻辑）
│   ├── project.go
│   ├── collection.go
│   ├── request.go
│   ├── environment.go
│   ├── history.go
│   └── settings.go
│
├── database/                    # 数据库初始化
│   ├── database.go              # 连接管理
│   └── migrations.go            # 建表 / 迁移
│
├── repositories/                # 数据访问（接口 + 实现）
│   ├── project_repo.go
│   ├── collection_repo.go
│   ├── request_repo.go
│   ├── environment_repo.go
│   ├── history_repo.go
│   └── settings_repo.go
│
├── services/                    # 业务逻辑（接口 + 实现）
│   ├── project_service.go
│   ├── collection_service.go
│   ├── request_service.go
│   ├── environment_service.go
│   ├── history_service.go
│   ├── settings_service.go
│   ├── import_service.go
│   ├── export_service.go
│   ├── docs_service.go
│   └── variable_service.go     # 变量替换
│
├── handlers/                    # Wails 绑定（入参校验 + 调用 service）
│   ├── project_handler.go
│   ├── collection_handler.go
│   ├── request_handler.go
│   ├── environment_handler.go
│   ├── history_handler.go
│   ├── settings_handler.go
│   ├── import_handler.go
│   ├── export_handler.go
│   ├── backup_handler.go
│   └── websocket_handler.go
│
└── pkg/
    ├── httpclient/
    │   ├── client.go            # HTTP 请求执行
    │   └── websocket.go         # WebSocket 客户端
    └── snowflake/
        └── snowflake.go         # ID 生成器
```

---

## 2. 分层约定

```
handlers  → 只做参数校验 + 类型转换，调用 service 返回结果
services  → 全部业务逻辑、事务协调
repos     → 纯 CRUD，不含业务判断
models    → 纯结构体，不含方法
```

规则：
- `handlers/` 不 import `repositories/`
- `repositories/` 不 import `services/` 或 `handlers/`
- `services/` 不 import `handlers/`
- `models/` 不 import 任何内部包
- 所有依赖注入在 `app.go` 的 `Startup` 中完成

---

## 3. 数据模型

### 3.1 models/project.go

```go
package models

type Project struct {
    ID          int64  `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    CreatedAt   string `json:"created_at"`
    UpdatedAt   string `json:"updated_at"`
}
```

### 3.2 models/collection.go

```go
package models

type Collection struct {
    ID        int64  `json:"id"`
    ProjectID int64  `json:"project_id"`
    ParentID  *int64 `json:"parent_id"`
    Name      string `json:"name"`
    SortOrder int    `json:"sort_order"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
}

// TreeItem 前端渲染用的统一树节点
type TreeItem struct {
    ID        int64      `json:"id"`
    Name      string     `json:"name"`
    Type      string     `json:"type"` // "folder" | "request"
    Method    string     `json:"method,omitempty"`
    URL       string     `json:"url,omitempty"`
    Children  []TreeItem `json:"children"`
    SortOrder int        `json:"sort_order"`
}
```

### 3.3 models/request.go

```go
package models

type Request struct {
    ID           int64  `json:"id"`
    CollectionID int64  `json:"collection_id"`
    Name         string `json:"name"`
    Description  string `json:"description"`
    Method       string `json:"method"`
    URL          string `json:"url"`
    Headers      string `json:"headers"`    // JSON
    Params       string `json:"params"`     // JSON
    BodyType     string `json:"body_type"`
    Body         string `json:"body"`       // JSON
    Auth         string `json:"auth"`       // JSON
    SortOrder    int    `json:"sort_order"`
    CreatedAt    string `json:"created_at"`
    UpdatedAt    string `json:"updated_at"`
}
```

### 3.4 models/environment.go

```go
package models

type Environment struct {
    ID        int64  `json:"id"`
    ProjectID int64  `json:"project_id"`
    Name      string `json:"name"`
    IsActive  bool   `json:"is_active"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
}

type EnvVariable struct {
    ID            int64  `json:"id"`
    EnvironmentID int64  `json:"environment_id"`
    Key           string `json:"key"`
    Value         string `json:"value"`
    Enabled       bool   `json:"enabled"`
    SortOrder     int    `json:"sort_order"`
    CreatedAt     string `json:"created_at"`
}
```

### 3.5 models/history.go

```go
package models

type History struct {
    ID              int64  `json:"id"`
    ProjectID       int64  `json:"project_id"`
    RequestID       *int64 `json:"request_id"`
    Method          string `json:"method"`
    URL             string `json:"url"`
    RequestHeaders  string `json:"request_headers"`
    RequestBody     string `json:"request_body"`
    ResponseStatus  int    `json:"response_status"`
    ResponseHeaders string `json:"response_headers"`
    ResponseBody    string `json:"response_body"`
    DurationMs      int    `json:"duration_ms"`
    CreatedAt       string `json:"created_at"`
}
```

### 3.6 models/settings.go

```go
package models

type Setting struct {
    Key       string `json:"key"`
    Value     string `json:"value"`
    UpdatedAt string `json:"updated_at"`
}
```

### 3.7 models/response.go（运行时，不持久化）

```go
package models

type HTTPResponse struct {
    Status       int               `json:"status"`
    StatusText   string            `json:"status_text"`
    Time         int64             `json:"time"`          // ms
    Size         int64             `json:"size"`          // bytes
    Headers      map[string]string `json:"headers"`
    Cookies      []Cookie          `json:"cookies"`
    Body         string            `json:"body"`
    RawRequest   string            `json:"raw_request"`
    CurlCommand  string            `json:"curl_command"`
}

type Cookie struct {
    Name   string `json:"name"`
    Value  string `json:"value"`
    Domain string `json:"domain"`
    Path   string `json:"path"`
}
```

---

## 4. Repository 接口

每个 repo 文件包含接口 + 实现。

### 4.1 ProjectRepository

```go
type ProjectRepository interface {
    List() ([]models.Project, error)
    GetByID(id int64) (*models.Project, error)
    Create(project *models.Project) error
    Update(project *models.Project) error
    Delete(id int64) error
}
```

### 4.2 CollectionRepository

```go
type CollectionRepository interface {
    ListByProject(projectID int64) ([]models.Collection, error)
    GetByID(id int64) (*models.Collection, error)
    Create(collection *models.Collection) error
    Update(collection *models.Collection) error
    Delete(id int64) error
    UpdateSortOrder(id int64, sortOrder int) error
    MoveToParent(id int64, parentID *int64) error
}
```

### 4.3 RequestRepository

```go
type RequestRepository interface {
    ListByCollection(collectionID int64) ([]models.Request, error)
    GetByID(id int64) (*models.Request, error)
    Create(request *models.Request) error
    Update(request *models.Request) error
    Delete(id int64) error
    UpdateSortOrder(id int64, sortOrder int) error
    MoveToCollection(id int64, collectionID int64) error
}
```

### 4.4 EnvironmentRepository

```go
type EnvironmentRepository interface {
    ListByProject(projectID int64) ([]models.Environment, error)
    GetByID(id int64) (*models.Environment, error)
    Create(env *models.Environment) error
    Update(env *models.Environment) error
    Delete(id int64) error
    SetActive(projectID int64, envID int64) error
    GetActive(projectID int64) (*models.Environment, error)
}
```

### 4.5 EnvVariableRepository

```go
type EnvVariableRepository interface {
    ListByEnvironment(envID int64) ([]models.EnvVariable, error)
    Create(variable *models.EnvVariable) error
    Update(variable *models.EnvVariable) error
    Delete(id int64) error
    BatchReplace(envID int64, variables []models.EnvVariable) error
}
```

### 4.6 HistoryRepository

```go
type HistoryRepository interface {
    ListByProject(projectID int64, limit, offset int) ([]models.History, error)
    Search(projectID int64, keyword string, method string, statusMin, statusMax int, limit, offset int) ([]models.History, int, error)
    Create(history *models.History) error
    Delete(id int64) error
    DeleteByProject(projectID int64) error
    DeleteOlderThan(projectID int64, days int) error
}
```

### 4.7 SettingsRepository

```go
type SettingsRepository interface {
    Get(key string) (string, error)
    Set(key, value string) error
    GetAll() (map[string]string, error)
    SetAll(settings map[string]string) error
}
```

---

## 5. Service 接口

### 5.1 ProjectService

```go
type ProjectService interface {
    List() ([]models.Project, error)
    Get(id int64) (*models.Project, error)
    // Create 创建项目并自动创建默认集合
    Create(name, description string) (*models.Project, error)
    Update(id int64, name, description string) (*models.Project, error)
    Delete(id int64) error
}
```

### 5.2 CollectionService

```go
type CollectionService interface {
    GetTree(projectID int64) ([]models.TreeItem, error)
    Get(id int64) (*models.Collection, error)
    Create(projectID int64, parentID *int64, name string) (*models.Collection, error)
    Rename(id int64, name string) error
    Move(id int64, parentID *int64, sortOrder int) error
    Delete(id int64) error
}
```

### 5.3 RequestService

```go
type RequestService interface {
    Get(id int64) (*models.Request, error)
    Create(collectionID int64, name, method string) (*models.Request, error)
    Update(req *models.Request) error
    Clone(id int64) (*models.Request, error)
    Move(id int64, collectionID int64, sortOrder int) error
    Delete(id int64) error
    // Send 执行 HTTP 请求（接受 context 支持取消）
    Send(ctx context.Context, req *models.Request, envID int64) (*models.HTTPResponse, error)
    // SendQuick 直接发送 URL（不保存、不关联请求，用于临时调试）
    SendQuick(ctx context.Context, method, url, headers, body string, envID int64) (*models.HTTPResponse, error)
}
```

### 5.4 EnvironmentService

```go
type EnvironmentService interface {
    List(projectID int64) ([]models.Environment, error)
    Create(projectID int64, name string, cloneFromID *int64) (*models.Environment, error)
    Rename(id int64, name string) error
    Delete(id int64) error
    Activate(id int64) error
    GetActive(projectID int64) (*models.Environment, error)

    ListVariables(envID int64) ([]models.EnvVariable, error)
    SaveVariables(envID int64, variables []models.EnvVariable) error
}
```

### 5.5 HistoryService

```go
type HistoryService interface {
    List(projectID int64, page, pageSize int) ([]models.History, int, error)
    Search(projectID int64, keyword, method string, statusMin, statusMax int, page, pageSize int) ([]models.History, int, error)
    Get(id int64) (*models.History, error)
    Clear(projectID int64) error
    Delete(id int64) error
    CleanOld(projectID int64, days int) (int, error)
}
```

### 5.6 SettingsService

```go
type SettingsService interface {
    Get(key string) (string, error)
    Set(key, value string) error
    GetAll() (map[string]string, error)
}
```

### 5.7 ImportService

```go
type ImportService interface {
    // ImportPostman 导入 Postman Collection v2.x
    ImportPostman(projectID int64, filePath string) (*ImportResult, error)
    // ImportOpenAPI 导入 OpenAPI 3.x / Swagger 2.x
    ImportOpenAPI(projectID int64, filePath string) (*ImportResult, error)
    // ImportPaw 导入 Paw 自有格式（.paw，JSON 格式，用于项目分享）
    ImportPaw(filePath string) (*ImportResult, error)
}

type ImportResult struct {
    Collections int `json:"collections"`
    Requests    int `json:"requests"`
}
```

### 5.8 ExportService

> 仅处理单个项目/集合的导出（用于分享）。全局备份/恢复走 `database.Backup()` / `database.Restore()` 直接操作 SQLite 文件。

```go
type ExportService interface {
    // ExportPostman 导出为 Postman Collection
    ExportPostman(projectID int64, collectionID *int64) (string, error) // 返回 JSON 字符串
    // ExportPaw 导出为 Paw 自有格式（.paw，JSON 格式，分享用）
    ExportPaw(projectID int64, collectionID *int64, filePath string) error
}
```

### 5.9 DocsService

```go
type DocsService interface {
    // GenerateMarkdown 生成 Markdown 文档
    GenerateMarkdown(projectID int64, collectionID *int64) (string, error)
    // GenerateHTML 生成 HTML 文档
    GenerateHTML(projectID int64, collectionID *int64) (string, error)
}
```

### 5.10 VariableService

```go
type VariableService interface {
    // Resolve 将文本中的 {{var}} 替换为实际值
    Resolve(input string, envID int64) (string, error)
    // ResolveMap 批量替换 map 中的变量
    ResolveMap(input map[string]string, envID int64) (map[string]string, error)
}
```

---

## 6. Handler 注册

所有 handler 方法绑定到 `app.go` 的 `App` 结构体上，Wails 自动生成前端调用代码。

```go
// app.go

type App struct {
    ctx         context.Context
    db          *sql.DB
    snowflake   *snowflake.Generator
    httpClient  *httpclient.Client

    // 运行时状态
    cancels sync.Map    // map[int64]context.CancelFunc — 请求取消（key: sessionID）
    wsConns sync.Map    // map[string]*wsclient.WSClient — WebSocket 连接池

    projectSvc     services.ProjectService
    collectionSvc  services.CollectionService
    requestSvc     services.RequestService
    envSvc         services.EnvironmentService
    historySvc     services.HistoryService
    settingsSvc    services.SettingsService
    importSvc      services.ImportService
    exportSvc      services.ExportService
    docsSvc        services.DocsService
    variableSvc    services.VariableService
}

func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
    a.db = database.Init()
    a.snowflake = snowflake.New()
    a.httpClient = httpclient.New()

    // 依赖注入
    projectRepo := repositories.NewProjectRepo(a.db)
    collectionRepo := repositories.NewCollectionRepo(a.db)
    // ... 其余 repo

    a.projectSvc = services.NewProjectService(projectRepo, a.snowflake)
    a.collectionSvc = services.NewCollectionService(collectionRepo, requestRepo, a.snowflake)
    // ... 其余 service
}
```

### 6.1 Handler 示例

```go
// handlers/project_handler.go

// ListProjects 返回所有项目
func (a *App) ListProjects() ([]models.Project, error) {
    return a.projectSvc.List()
}

// CreateProject 创建项目
func (a *App) CreateProject(name, description string) (*models.Project, error) {
    if strings.TrimSpace(name) == "" {
        return nil, errors.New("project name is required")
    }
    return a.projectSvc.Create(name, description)
}

// DeleteProject 删除项目
func (a *App) DeleteProject(id int64) error {
    return a.projectSvc.Delete(id)
}
```

### 6.2 WebSocket Handler（事件驱动）

Wails 是请求-响应模式，WebSocket 是长连接推流，需要 Wails Events 桥接：

```go
// handlers/websocket_handler.go

// Connect 建立 WebSocket 连接
func (a *App) WSConnect(url string, headersJSON string) error {
    var headers map[string]string
    json.Unmarshal([]byte(headersJSON), &headers)

    conn := wsclient.NewWSClient()
    err := conn.Connect(url, headers)
    if err != nil {
        return err
    }

    a.wsConns.Store(url, conn)

    // 启动读取协程，通过 Wails Events 推消息到前端
    go func() {
        for {
            msg, err := conn.Read()
            if err != nil {
                runtime.EventsEmit(a.ctx, "ws:error", url, err.Error())
                break
            }
            runtime.EventsEmit(a.ctx, "ws:message", url, string(msg))
        }
        runtime.EventsEmit(a.ctx, "ws:closed", url)
    }()

    runtime.EventsEmit(a.ctx, "ws:connected", url)
    return nil
}

// WSSend 发送消息
func (a *App) WSSend(url string, message string) error

// WSDisconnect 断开连接
func (a *App) WSDisconnect(url string) error
```

前端订阅 Wails Events：

```typescript
import { EventsOn } from '@wailsjs/runtime/runtime'

EventsOn('ws:message', (url: string, message: string) => {
  // 追加到消息历史
})
EventsOn('ws:connected', (url: string) => { ... })
EventsOn('ws:closed', (url: string) => { ... })
EventsOn('ws:error', (url: string, err: string) => { ... })
```

### 6.3 Backup Handler

```go
// handlers/backup_handler.go

// BackupAll 备份全部数据到指定文件
func (a *App) BackupAll(filePath string) error

// RestoreAll 从备份文件恢复（覆盖现有数据）
func (a *App) RestoreAll(filePath string) error
```

实现：直接复制 SQLite 文件 + WAL/SHM，或逐表导出为 JSON。优先采用文件级复制，简单可靠。

---

## 7. 请求取消

Handler 层维护 `sync.Map` 存储 `sessionID → context.CancelFunc`：

```go
// handlers/request_handler.go

// SendRequest 发送请求，前端传入 sessionID（前端自生成，如简单自增计数器）
func (a *App) SendRequest(sessionID int64, req *models.Request, envID int64) (*models.HTTPResponse, error) {
    ctx, cancel := context.WithCancel(a.ctx)
    a.cancels.Store(sessionID, cancel)
    defer a.cancels.Delete(sessionID)
    return a.requestSvc.Send(ctx, req, envID)
}

// CancelRequest 取消指定 sessionID 的请求
func (a *App) CancelRequest(sessionID int64) {
    if cancel, ok := a.cancels.Load(sessionID); ok {
        cancel()
        a.cancels.Delete(sessionID)
    }
}
```

前端：

```typescript
let sessionCounter = 0
function send() {
    const sid = ++sessionCounter
    SendRequest(sid, currentRequest, activeEnvId)
}
function cancel() {
    CancelRequest(sessionCounter)
}
```

---

## 8. HTTP 客户端

```go
package httpclient

type Client struct {
    timeout        time.Duration
    followRedirect bool
    maxRedirects   int
    verifySSL      bool
    proxyURL       *url.URL
    rootCAs        *x509.CertPool
    certs          []tls.Certificate
}

func New() *Client

// ApplySettings 应用 settings 中的配置
func (c *Client) ApplySettings(settings map[string]string)

// Execute 执行 HTTP 请求
func (c *Client) Execute(ctx context.Context, req *Request) (*Response, error)

type Request struct {
    Method  string
    URL     string
    Headers map[string]string
    Body    []byte
}

type Response struct {
    Status     int
    StatusText string
    Headers    map[string]string
    Body       []byte
    TimeMs     int64
    Size       int64
}
```

### 8.2 pkg/httpclient/websocket.go

```go
type WSClient struct {
    conn     *websocket.Conn
    isConnected bool
}

func NewWSClient() *WSClient

// Connect 建立 WebSocket 连接
func (w *WSClient) Connect(url string, headers map[string]string) error

// Disconnect 断开连接
func (w *WSClient) Disconnect() error

// Send 发送消息
func (w *WSClient) Send(message []byte) error

// Read 读取消息（阻塞）
func (w *WSClient) Read() ([]byte, error)

// IsConnected 连接状态
func (w *WSClient) IsConnected() bool
```

---

## 9. 启动流程

```
main.go
 └→ wails.Run()
     └→ App.startup()
         ├→ database.Init()          打开/创建 SQLite
         │   └→ migrations.Run()     执行建表
         ├→ snowflake.New()          ID 生成器
         ├→ httpclient.New()         HTTP 客户端
         ├→ 创建 Repository 实例
         ├→ 创建 Service 实例（注入 Repo）
         ├→ settingsSvc.GetAll()     加载全局设置
         │   └→ httpClient.ApplySettings()
         └→ projectSvc.List()        加载项目列表，选中上次项目
```

---

## 10. 数据库初始化

```go
// database/database.go
func Init() *sql.DB {
    db, err := sql.Open("sqlite", "./paw.db")
    // ...
    db.SetMaxOpenConns(1)  // SQLite 单写
    migrations.Run(db)
    return db
}

// database/migrations.go
func Run(db *sql.DB) {
    for _, ddl := range []string{
        createProjects,
        createCollections,
        createRequests,
        createEnvironments,
        createEnvVariables,
        createHistory,
        createSettings,
    } {
        db.Exec(ddl)
    }
}
```