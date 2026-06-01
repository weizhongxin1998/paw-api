# Bug 修复 & 功能补全计划

审计来源：全量代码审计，发现 20 个问题。

---

## Batch 1: 关键 Bug（功能直接不可用）

### 1.1 修复导入功能 — 数据未持久化

**问题**: `ImporterService` 解析 Postman/Swagger/cURL 后只返回数据，`ProjectGallery.handleImport` 收到后直接丢弃。

**Go 后端:**
- `handlers/importer_handler.go` — `ImportPostman`/`ImportSwagger` 方法：解析成功后调用 `CollectionService.Create` + `RequestService.Create` 逐个写入数据库。保留原始返回格式（兼容）。

**前端:**
- `ProjectGallery.vue` — `handleImport` 调用成功后刷新项目列表。

**Work**: 中

---

### 1.2 修复导出功能 — 请求数组为空

**问题**: `handleExport` 中 `reqsJSON = JSON.stringify([])` 硬编码空数组。

**前端:**
- `ProjectGallery.vue` — 导出前先调 `ListCollections` + `ListRequests` 遍历所有集合收集实际请求数据，传入 `ExportPostman`/`ExportSwagger`。

**Work**: 小

---

### 1.3 修复拖拽排序 — UpdateCollection 空 name 校验失败

**问题**: `UpdateCollection(draggedId, '', targetId, 0)` 传空 name，`CollectionService.Update` 校验 `name == ""` 报错。

**Go 后端:**
- `services/collection_service.go` — `Update` 方法：修改为仅当新 name 非空时才校验和更新 name 字段，允许仅更新 parent_id/sort_order。

**前端:**
- `AppSidebar.vue` — drop handler 传集合的实际 name 而非空字符串。

**Work**: 小

---

### 1.4 修复保存接口 — 始终写第一个集合

**问题**: `handleSave` 中 `colId = projectStore.collections[0]?.id`，不管当前请求属于哪个集合。

**前端:**
- `RequestEditor.vue` — `handleSave` 新请求时需知道目标 collection。方案：在 tab 数据中新增 `collectionId` 字段：
  - `stores/tabs.ts` — `HttpTabData` 加 `collectionId?: string`
  - `AppSidebar.vue` — `handleTreeSelect` 和 `createRequestInCollection` 设置 `collectionId`
  - `RequestEditor.vue` — Save 时优先用 `tab.httpData?.collectionId`，其次才用 `collections[0]`

**Work**: 中

---

### 1.5 修复 Digest 认证 — probe 请求无 body

**问题**: `applyDigestAuth` 创建 probe 请求时 `http.NewRequest(method, rawURL, nil)`，对 POST/PUT 无效。

**Go 后端:**
- `services/request_service.go` — `applyDigestAuth` 改为接收 body 参数，probe 请求使用相同 body（或至少非 nil body）。

**Work**: 小

---

### 1.6 修复历史页面 — 首次无项目时空白

**问题**: 无当前项目时 `projectID = ''`, `ListHistory('', 100)` 无结果。

**Go 后端:**
- `handlers/history_handler.go` — `ListHistory`: projectID 为空时返回所有最近历史，而非空列表。

**前端:**
- `History.vue` — 无项目时显示所有历史。

**Work**: 小

---

## Batch 2: 缺失功能

### 2.1 项目删除 / 重命名 UI

**后端已就绪**（`DeleteProject`、`UpdateProject` handler 已存在），只需前端。

**前端:**
- `stores/project.ts` — 添加 `removeProject(id)`、`updateProject(p)` 方法
- `ProjectGallery.vue` — 项目卡片悬浮显示 `...` 菜单 → 右键或点击显示 Rename / Delete
- 删除前弹窗确认；删除当前项目后自动切换到第一个项目

**Work**: 小

---

### 2.2 历史记录删除 / 清空 UI

**后端已就绪**（`DeleteHistory`、`ClearHistory` handler 已存在）。

**前端:**
- `stores/history.ts` — 添加 `removeHistory(id)`、`clearHistory(projectId)` 方法
- `History.vue` — 每条历史记录右边加删除按钮；工具栏加"清空"按钮

**Work**: 小

---

### 2.3 Cookie 管理 UI 面板接入

**后端已就绪**（`CookieHandler` 已创建、已绑定），但 Cookie jar 与 HTTP client 未连接。

**Go 后端:**
- `services/request_service.go` — `Send()` 中复用同一个 `httpclient.Client` 实例（而非每次 new），让 cookie jar 跨请求共享
- `pkg/httpclient/client.go` — `NewClient()` 接受可选 `http.CookieJar` 参数

**前端:**
- `AppSidebar.vue` — Settings 面板新增 Cookies 区域
- 表格显示域名、名称、值、过期时间；支持删除单个和清空全部
- 调 `ListCookies` / `ClearCookies`

**Work**: 中

---

### 2.4 断言规则持久化

**问题**: 断言配置只在组件 ref 中，刷新丢失。

**方案**: 将断言规则存入 `request.script` 字段（JSON 数组）。

**前端:**
- `TestRunner.vue` — 加载请求时从 `req.script` 解析已有断言规则
- 修改断言时写回 `request.script`
- 新增 `UpdateRequest` 调用在修改断言时

**Work**: 小

---

### 2.5 Request `collection_id` 支持更新

**Go 后端:**
- `repositories/request_repo.go` — `Update` SQL 增加 `collection_id = ?`
- `handlers/request_handler.go` — `UpdateRequest` 增加 `collectionID` 参数
- `services/request_service.go` — `Update` 增加 `collectionID` 参数

**前端:**
- `AppSidebar.vue` — `confirmMove` 改为直接调 `UpdateRequest` 改 `collection_id`，而非删+重建

**Work**: 小

---

## Batch 3: 代码质量 / 体验优化

### 3.1 form-data 保存/恢复修复

**前端:**
- `AppSidebar.vue` — `safeParseBody` 和 `safeParseBodyType` 支持恢复 form-data 文件字段
- `RequestEditor.vue` — `handleSave` bodyJSON 包含 `body_files` 数据

**Work**: 小

---

### 3.2 代码生成完善 body 类型支持

**Go 后端:**
- `services/codegen_service.go` — curl 生成器加 `Content-Type` 头；JS/Go 生成器支持 urlencoded/form-data

**Work**: 小

---

### 3.3 右键菜单 i18n 补充

**前端:**
- `i18n/en.ts`、`zh-CN.ts` — 添加 context menu 相关 key
- `AppSidebar.vue` — 硬编码字符串替换为 `$t()` 调用

**Work**: 小

---

### 3.4 请求加载优化 — N+1 查询

**Go 后端:**
- `repositories/request_repo.go` — 新增 `ListByProjectID(projectID string)` 方法，一次查询返回项目下所有请求
- `handlers/request_handler.go` — 暴露 `ListRequestsByProject(projectID)` 绑定

**前端:**
- `AppSidebar.vue` — `loadAllRequests` 改为单次调用

**Work**: 小

---

### 3.5 WebSocket goroutine 泄漏修复

**Go 后端:**
- `pkg/httpclient/websocket.go` — `Disconnect` 等待 `done` channel 确认 goroutine 退出后再返回

**Work**: 小

---

### 3.6 HTTP Client 线程安全

**Go 后端:**
- `pkg/httpclient/client.go` — `CheckRedirect` 作为 `Request` 字段而非 Client 字段，避免每次 `Do()` 修改共享状态。或使用 `sync.Mutex` 保护。

**Work**: 小

---

### 3.7 新建集合/请求 sort_order 自动递增

**Go 后端:**
- `services/collection_service.go` — `Create` 计算当前 parent 下 max sort_order + 1
- `services/request_service.go` — `Create` 计算当前 collection 下 max sort_order + 1

**Work**: 小

---

## 工期估算

| Batch | Tasks | Work | 预计 |
|-------|-------|------|------|
| Batch 1 | 6 个关键 Bug | 中×3 + 小×3 | **优先** |
| Batch 2 | 5 个缺失功能 | 中×1 + 小×4 | |
| Batch 3 | 7 个体验优化 | 小×7 | **低优先** |

总计 18 项修复，全部为小或中工作量。
