# Apifox 对标 — 功能补全计划

## 不做

- 账号体系
- 云端同步
- Mock 服务
- 前置/后置脚本引擎

---

## Task 1: Auth 认证模块

**目标**：支持 None / Basic Auth / Bearer Token / Digest / OAuth 2.0 五种认证方式。

**Go 后端：**
- `services/request_service.go` `Send` 方法：根据 auth 类型构造对应认证头：
  - `none` — 无操作
  - `basic` — `Authorization: Basic base64(user:pass)`
  - `bearer` — `Authorization: Bearer <token>`
  - `digest` — 先发请求获取 401 的 `WWW-Authenticate` challenge，再计算 hash 重发
  - `oauth2` — 客户端凭证模式：POST token_url 获取 access_token
- `handlers/request_handler.go` — `SendRequestInput` 新增字段：
  ```go
  AuthType string // none | basic | bearer | digest | oauth2
  AuthData map[string]string // 认证参数
  ```

**前端：**
- `RequestEditor.vue` Auth 标签页改为认证类型下拉 + 对应表单：
  - Basic → 用户名 + 密码输入
  - Bearer → Token 输入
  - Digest → 用户名 + 密码输入（Go 侧自动处理 challenge）
  - OAuth2 → Client ID / Secret / Token URL

**Work**：中型

---

## Task 2: Body 类型扩展

**目标**：支持 `form-data`、`x-www-form-urlencoded`、`binary`、现有 `none`/`json`/`text`。

**Go 后端：**
- `httpclient/client.go` `Do` 方法：根据 `BodyType` 设置 Content-Type 并编码：
  - `form-data` → `multipart/form-data`（KV + 可选的 file_path）
  - `urlencoded` → `application/x-www-form-urlencoded`
  - `binary` → `application/octet-stream`，body 为 base64 解码后的原始字节
  - `json` / `text` / `none` — 已有

**前端：**
- `RequestEditor.vue` Body 标签页新增选项：`form-data` / `x-www-form-urlencoded` / `binary`
- `form-data` 用 `KeyValueEditor.vue` 加一行"文件路径"输入

**Work**：中型

---

## Task 3: 断言系统

**目标**：运行测试时自动执行断言，显示 PASS/FAIL。

**Go 后端（新文件）：**
- `services/assert_service.go` — 执行断言列表，返回结果：
  ```go
  type AssertRule struct {
      Type   string // status | body_contains | body_jsonpath | header_equals | duration_lt
      Target string
      Value  string
  }
  type AssertResult struct {
      Rule   AssertRule
      Passed bool
      Actual string
      Error  string
  }
  ```
- 断言类型：
  - `status` — 状态码等于/不等于期望值
  - `body_contains` — 响应体包含/不包含指定字符串
  - `body_jsonpath` — JSONPath 提取值等于/存在
  - `header_equals` — 响应头等于/包含
  - `duration_lt` — 耗时小于 N ms
- `handlers/request_handler.go` — `RunAsserts` 方法：发送请求 + 执行断言论 → 返回结果

**前端：**
- `TestRunner.vue` 每个请求可展开断言论列表
- 断言编辑器：类型下拉 + 目标/值输入
- 运行后在请求结果旁显示 PASS/FAIL + 实际值

**Work**：大型

---

## Task 4: 导入导出

**目标**：支持 Postman、Swagger/OpenAPI 3.0 的导入和导出，以及 cURL 导入。

**Go 后端（新文件）：**
- `services/importer.go` — 解析文件内容 → `[]models.Collection` + `[]models.Request`
  - `ImportPostman(json string)` — Postman v2.1 collection JSON
  - `ImportSwagger(json string)` — OpenAPI 3.0 JSON/YAML
  - `ImportCurl(curl string)` — 解析 cURL 命令行为 Request
- `services/exporter.go` — 将集合树序列化：
  - `ExportPostman(projectID string) string` → Postman v2.1 JSON
  - `ExportSwagger(projectID string) string` → OpenAPI 3.0 JSON
- `handlers/importer_handler.go` — Wails bind
- `handlers/exporter_handler.go` — Wails bind

**前端：**
- 侧栏 Workspace 面板添加"导入"按钮 → 文件选择 → 解析 → 合并到当前项目集合树
- 侧栏添加"导出"按钮 → 选择格式（Postman / Swagger）→ 触发下载
- cURL 导入用独立弹窗粘贴命令

**Work**：大型

---

## Task 5: 请求头快捷预设

**目标**：一键添加常用请求头。

**前端：**
- `KeyValueEditor.vue` 新增"Presets"下拉按钮
- 预设列表（可配置）：
  - `Authorization: Bearer <token>`
  - `Content-Type: application/json`
  - `Content-Type: application/x-www-form-urlencoded`
  - `Accept: application/json`
  - `User-Agent: PawAPI/1.0`
  - `Cache-Control: no-cache`
  - `Referer:`
- 点击预设自动添加一行（已存在的 key 跳过）

**Work**：小型

---

## Task 6: 代码生成

**目标**：根据当前请求生成多语言代码片段。

**Go 后端（新文件）：**
- `services/codegen_service.go` — 模板生成：
  - cURL
  - JavaScript (fetch)
  - Python (requests)
  - Go (net/http)
- `handlers/request_handler.go` — `GenerateCode` 方法

**前端：**
- 响应区添加"Code"标签页
- 语言下拉选择 → 渲染代码 → 复制按钮

**Work**：中型

---

## Task 7: 拖拽排序 + 复制/移动请求

**目标**：集合树支持拖拽排序和右键操作。

**前端：**
- `AppSidebar.vue` 集合树 → `NTree` 设置 `draggable` + `on-drop`
- 拖拽结束 → 更新 `parent_id` + `sort_order` → 调后端 `UpdateCollection`
- 右键菜单（用 `NDropdown` + `NPopover` 或 contextmenu 事件）：
  - 重命名
  - 复制请求到当前集合
  - 移动请求到其他集合
  - 删除
  - 新建子集合

**Go 后端：**
- 已有 `UpdateCollection`（sort_order 更新）
- 可能需要 `UpdateRequest` 传 `collection_id` 变更（移动请求）

**Work**：中型

---

## Task 8: 全局搜索接口

**目标**：快速搜索并打开任意接口。

**前端：**
- TabBar 旁边添加快捷搜索框（`Ctrl+K` 唤起）
- 输入关键词 → 遍历当前项目所有 `collections` + `requests`，按名称和 URL 模糊匹配
- 下拉列表展示匹配结果，带集合名和请求方法标签
- 点击结果 → 在工作区打开对应标签页
- 前端本地过滤即可，无需后端

**Work**：小型

---

## Task 9: Cookie 管理面板

**目标**：查看和管理 HTTP 客户端维护的 cookies。

**Go 后端：**
- `httpclient/client.go` — HTTP Client 启用 cookie jar
- `services/cookie_service.go` — 列出当前 cookies 和删除
- `handlers/cookie_handler.go` — Wails bind

**前端：**
- Settings 面板或独立管理页，表格列出 cookies（域名、名称、值、过期时间）
- 支持手动添加和删除

**Work**：中型

---

## Task 10: 请求配置（超时 / 跟随重定向）

**目标**：每个请求可独立设置超时和重定向策略。

**Go 后端：**
- `httpclient/client.go` — `Client` 新增 `TimeoutMs`、`FollowRedirects` 字段
- `sendRequestInput` 新增：
  ```go
  TimeoutMs       int  // 默认 30000
  FollowRedirects bool // 默认 true
  ```

**前端：**
- `RequestEditor.vue` URL 行右侧添齿轮按钮 → 弹出小面板
- 超时滑块（1s - 120s）
- "跟随重定向" 复选框

**Work**：小型

---

## 工期估算

| 优先级 | Task | Work |
|---|---|---|
| P0 | 1 — Auth 认证 | 中 |
| P0 | 2 — Body 类型扩展 | 中 |
| P0 | 3 — 断言系统 | 大 |
| P1 | 4 — 导入导出 | 大 |
| P1 | 5 — 请求头预设 | 小 |
| P1 | 6 — 代码生成 | 中 |
| P2 | 7 — 拖拽排序 + 复制请求 | 中 |
| P2 | 8 — 全局搜索 | 小 |
| P2 | 9 — Cookie 管理 | 中 |
| P2 | 10 — 超时/重定向 | 小 |
