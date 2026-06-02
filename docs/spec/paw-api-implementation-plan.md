# Paw API — 实施计划

> 版本 1.0 | 2026-06-02 | 基于 [paw-api-requirements.md](./paw-api-requirements.md)

---

## Phase 1 — 框架搭建

**目标**：可运行的 Wails 骨架，数据库就绪，前后端目录结构建成，基础布局渲染。

| # | 任务 | 产出 | 验证 |
|---|---|---|---|
| 1.1 | `wails init` 重建项目，保留 `go.mod` 依赖项 | 可编译运行的空窗口 | `wails dev` 启动成功 |
| 1.2 | Go 包结构：建 `models/`、`database/`、`repositories/`、`services/`、`handlers/`、`pkg/snowflake/`、`pkg/httpclient/` | 目录树 + 空文件 | 编译通过，无 import 错误 |
| 1.3 | `database/database.go` + `migrations.go`：SQLite 连接、`SetMaxOpenConns(1)`、9 张表的 DDL | 运行后生成 `paw.db` | 查询各表存在 |
| 1.4 | `pkg/snowflake/snowflake.go`：48 位 ID 生成器 | `snowflake.New().Next()` 返回唯一递增值 | 并发 10000 次无碰撞 |
| 1.5 | `models/` 全部结构体定义（project, collection, request, environment, history, settings, response） | 与数据库、JSON 序列化对齐 | 字段与 DB 列名一致 |
| 1.6 | Vue 项目初始化：`npm create vue`，装 NaiveUI + Pinia + Ionicons，建立 `frontend/src/` 目录结构 | `npm run dev` 可启动 | 页面渲染无报错 |
| 1.7 | 主布局骨架：`AppHeader` + `Sidebar` + `Workspace`（左右分栏 + 占位文本） | 三栏布局渲染 | 浏览器尺寸 1920×1080 比例正常 |

---

## Phase 2 — HTTP 请求调试（核心）

**目标**：能创建项目、集合、请求，能编辑请求参数，能发送 HTTP 请求并查看响应。

| # | 任务 | 产出 | 验证 |
|---|---|---|---|
| 2.1 | Repository 层全部实现（project, collection, request, environment, history, settings） | 完整 CRUD + 搜索查询 | 各方法返回正确数据 |
| 2.2 | Service 层：project, collection, request（含 Send/SendQuick） | 业务逻辑完整 | 创建项目→创建集合→创建请求→发送→返回响应 |
| 2.3 | Handler 层全部实现，Wails 绑定注册 | 前端可调用全部 API | `@/wailsjs/go/main/App` 方法可用 |
| 2.4 | `pkg/httpclient/client.go`：HTTP 请求执行，支持 method/url/headers/body/form-data/binary | 各种请求类型成功 | curl 等价请求返回一致 |
| 2.5 | `CollectionTree` + `TreeNode` 组件：递归渲染，右键菜单 | 侧栏树可操作 | 新建/重命名/删除集合和请求 |
| 2.6 | `RequestPanel` + `UrlBar` + 子Tab（Params/Headers/Body/Auth） | 请求编辑器可交互 | 填参数→点 Send→看响应 |
| 2.7 | `ResponsePanel` + `StatusBar` + `ResponseBody(Pretty/Raw/Preview)` | 响应完整展示 | JSON 格式化、HTML 预览正常 |
| 2.8 | 请求发送关联 `sessionID`，支持取消（`CancelRequest`） | 可中断请求 | 慢接口点取消后显示"已取消" |

---

## Phase 3 — 导入/导出

**目标**：支持 Postman Collection v2.x、OpenAPI 3.x / Swagger 2.x 导入，Postman Collection 导出。

| # | 任务 | 产出 | 验证 |
|---|---|---|---|
| 3.1 | `ImportService.ImportPostman`：解析 JSON → 创建集合+请求 | Postman 文件导入成功 | 用 Postman 导出的 JSON 测试 |
| 3.2 | `ImportService.ImportOpenAPI`：解析 YAML/JSON → 创建集合+请求 | OpenAPI/Swagger 导入成功 | 用 petstore.yaml 测试 |
| 3.3 | `ExportService.ExportPostman`：项目→Postman Collection JSON | Postman 格式导出 | 导出的文件可被 Postman 打开 |
| 3.4 | `ExportService.ExportPaw`：项目→Paw 自有 JSON 格式 | .paw 文件导出 | 导出→重新导入数据一致 |
| 3.5 | `ImportModal` + `ExportModal` 前端组件 | 可视化的导入导出流程 | 选择文件→确认→完成反馈 |

---

## Phase 4 — 环境变量

**目标**：多环境创建/切换，`{{var}}` 插值替换，URL/Header/Body/Auth 全覆盖。

| # | 任务 | 产出 | 验证 |
|---|---|---|---|
| 4.1 | `EnvironmentService` + `EnvVariableRepository` 全部实现 | 环境 CRUD + 变量 CRUD | 创建环境→添加变量→查询正确 |
| 4.2 | `VariableService.Resolve`：文本 `{{var}}` 替换 | 变量引擎 | 嵌套变量、缺失变量有明确行为 |
| 4.3 | `EnvSelector` 组件：环境下拉切换 | 标题栏可切换环境 | 切换后激活标识更新 |
| 4.4 | `EnvManagerModal`：变量增删改、启用/禁用 | 环境管理弹窗 | 添加/删除变量→保存生效 |
| 4.5 | URL 输入中 `{{var}}` 高亮 + 悬停解析值 | 前端视觉反馈 | 黄色背景高亮 + tooltip |

---

## Phase 5 — 多标签页

**目标**：多请求并行编辑，Tab 状态独立，关闭/拖拽/右键操作。

| # | 任务 | 产出 | 验证 |
|---|---|---|---|
| 5.1 | `tabsStore`：Tab 状态管理，active切换，dirty 标记 | 状态逻辑 | 打开 5 个 Tab→各自动编辑独立 |
| 5.2 | `TabBar` 右键菜单：关闭当前/其他/右侧/全部 | 右键操作 | 各操作正确 |
| 5.3 | Tab 拖拽排序（vuedraggable 或原生 Drag API） | 拖拽功能 | 拖拽后顺序持久 |
| 5.4 | 关闭未保存 Tab 弹窗确认 | 保存提示 | 修改后关闭→弹窗提示 |
| 5.5 | 切换项目时清空 Tab | 项目隔离 | 切项目→Tab 全部关闭 |

---

## Phase 6 — 历史记录

**目标**：每次请求自动记录，可搜索、筛选、回放。

| # | 任务 | 产出 | 验证 |
|---|---|---|---|
| 6.1 | `HistoryService`：自动写入、列表查询、搜索 | 历史 CRUD | 发送请求→历史列表出现 |
| 6.2 | `HistoryPanel`：列表 + 搜索 + 方法/状态码筛选 | 历史界面 | 搜索"users"过滤出匹配项 |
| 6.3 | 点击历史项→新 Tab 回放（自动填充请求参数） | 回放功能 | 回放后参数与原请求一致 |
| 6.4 | 清空全部历史、按保留天数自动清理 | 清理功能 | 设置 1 天→旧记录自动删除 |

---

## Phase 7 — WebSocket 调试

**目标**：WebSocket 连接/收发消息，Wails Events 事件驱动。

| # | 任务 | 产出 | 验证 |
|---|---|---|---|
| 7.1 | `pkg/httpclient/websocket.go`：gorilla/websocket 封装 | WS 客户端 | Connect→Send→收到回复→Disconnect |
| 7.2 | `websocket_handler.go`：Wails Events 桥接（ws:message/connected/closed/error） | 事件推送 | 前端 `EventsOn` 收到推送 |
| 7.3 | `WebSocketPanel` + `WsUrlBar` + `MessagePane` + `MessageInput` | WS 界面 | URL 连接→发消息→消息列表实时更新 |
| 7.4 | HTTP↔WS 面板切换：隐藏 TabBar，断开恢复 | 面板切换 | 打开 WS 面板→HTTP Tab 隐藏→关闭 WS→恢复 |

---

## Phase 8 — API 文档生成

**目标**：基于项目集合/请求自动生成文档，导出 Markdown/HTML。

| # | 任务 | 产出 | 验证 |
|---|---|---|---|
| 8.1 | `DocsService.GenerateMarkdown`：按集合分组的接口文档 | Markdown 输出 | 含方法、路径、参数、示例 |
| 8.2 | `DocsService.GenerateHTML`：HTML 独立页面 | HTML 输出 | 浏览器可打开，样式可读 |
| 8.3 | `DocsPreviewModal`：应用内文档预览 | 预览界面 | 实时反映集合最新配置 |

---

## 并行关系

```
Phase 1 ──┬── Phase 2 ──┬── Phase 3
          │             ├── Phase 5
          │             ├── Phase 6
          │             └── Phase 7
          │
          └── Phase 4 ── (注入 Phase 2)
                          
Phase 8 依赖 Phase 2（需集合和请求数据）
```

- Phase 4 可与 Phase 2 并行（环境变量独立模块）
- Phase 3/5/6/7 依赖 Phase 2 完成（需请求发送能力）
- Phase 8 最后执行

---

## 未纳入的功能（后续版本）

- 自动化测试（脚本引擎）
- GraphQL 支持
- 请求脚本（Pre-request Script）
- SOCKS 代理
- 请求链（变量提取+传递）
