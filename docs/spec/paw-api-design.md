# Paw API — 本地 API 调试与管理工具 设计文档

## 1. 项目概述

基于 Wails + Go + Vue 3 的桌面端 API 调试与管理工具，类似 Apifox 的核心体验，完全本地化运行（无账号、无云端同步）。

## 2. 技术选型

| 层 | 技术 | 备注 |
|---|---|---|
| 桌面框架 | Wails v2 | |
| 后端语言 | Go | |
| 前端框架 | Vue 3 + TypeScript | |
| 状态管理 | Pinia | |
| 数据存储 | SQLite | `modernc.org/sqlite`（纯 Go 实现，无 CGO） |
| UI 组件 | Naive UI + Ionicons v5 | Vue 3 + TypeScript 原生支持，tree-shaking |
| 图标集 | `@vicons/ionicons5` | 线条极简风格 |

## 3. 核心功能

- **HTTP 请求调试** — GET/POST/PUT/DELETE/PATCH 等，支持 Params、Headers、Body、Auth
- **集合管理** — 树形结构分组，支持拖拽排序，导入/导出（Postman/Swagger）
- **多标签页** — 每个请求在一个独立 Tab 中编辑
- **环境变量** — 多环境（开发/测试/生产），变量引用 `{{variable}}`
- **历史记录** — 自动保存请求历史，支持搜索和筛选
- **WebSocket 调试** — WebSocket 连接测试与消息收发
- **自动化测试** — 编写测试脚本，批量运行断言
- **API 文档** — 自动生成可分享的 API 文档（导出为 Markdown/HTML）
- **Mock 服务** — 可选，后续扩展

## 4. 架构方案

采用 **分层架构（方案 B）**：

```
┌─────────────────────────────────────────┐
│  Vue 3 前端                              │
│  ┌─────────┐ ┌──────────┐ ┌──────────┐  │
│  │  Views   │ │  Stores  │ │Composable│  │
│  └────┬────┘ └────┬─────┘ └────┬─────┘  │
│       │           │            │         │
│  ┌────┴───────────┴────────────┴────┐    │
│  │        Wails Binding             │    │
│  └──────────────────────────────────┘    │
├─────────────────────────────────────────┤
│  Go 后端                                │
│  ┌──────────┐                           │
│  │ Handlers │ ← Wails Bind              │
│  └────┬─────┘                           │
│       │                                  │
│  ┌────▼─────┐                           │
│  │ Services │ ← 业务逻辑                 │
│  └────┬─────┘                           │
│       │                                  │
│  ┌────▼──────────┐                      │
│  │ Repositories  │ ← 数据库操作          │
│  └────┬──────────┘                      │
│       │                                  │
│  ┌────▼────┐                            │
│  │ SQLite  │                            │
│  └─────────┘                            │
└─────────────────────────────────────────┘
```

## 5. 数据流

前端（用户操作） → Wails Binding → Handlers → Services → Repositories → SQLite

所有数据库操作通过 Repository 层访问，业务逻辑集中在 Service 层，Handler 层只做入参校验和数据转换。

## 6. 目录结构

```
paw-api/
├── main.go                    # Wails 应用入口
├── app.go                     # Wails 应用结构体、Binding 注册
├── wails.json                 # Wails 配置
│
├── models/                    # 数据模型定义
│   ├── project.go
│   ├── collection.go
│   ├── request.go
│   ├── environment.go
│   ├── history.go
│   └── response.go
│
├── database/                  # 数据库初始化与迁移
│   ├── database.go            # 连接管理
│   └── migrations.go          # 建表/升级
│
├── repositories/              # 数据访问层
│   ├── project_repo.go
│   ├── collection_repo.go
│   ├── request_repo.go
│   ├── environment_repo.go
│   └── history_repo.go
│
├── services/                  # 业务逻辑层
│   ├── project_service.go
│   ├── collection_service.go
│   ├── request_service.go
│   ├── environment_service.go
│   ├── history_service.go
│   ├── test_service.go
│   └── docs_service.go
│
├── handlers/                  # Wails Binding 层
│   ├── project_handler.go
│   ├── collection_handler.go
│   ├── request_handler.go
│   ├── environment_handler.go
│   ├── history_handler.go
│   ├── test_handler.go
│   └── docs_handler.go
│
├── pkg/
│   └── httpclient/            # HTTP 客户端封装
│       ├── client.go          # 通用 HTTP 请求
│       └── websocket.go       # WebSocket 客户端
│
└── frontend/
    └── src/
        ├── main.ts
        ├── App.vue
        ├── components/        # 通用组件
        │   ├── Sidebar.vue
        │   ├── TabBar.vue
        │   ├── RequestEditor.vue
        │   ├── ResponseViewer.vue
        │   ├── KeyValueEditor.vue
        │   ├── EnvSelector.vue
        │   └── JsonEditor.vue
        │
        ├── views/             # 页面
        │   ├── Workspace.vue
        │   ├── History.vue
        │   ├── Docs.vue
        │   └── TestRunner.vue
        │
        ├── stores/            # Pinia store
        │   ├── project.ts
        │   ├── request.ts
        │   ├── environment.ts
        │   ├── history.ts
        │   └── tabs.ts
        │
        ├── composables/       # 可复用组合式函数
        │   ├── useRequest.ts
        │   ├── useCollection.ts
        │   └── useTheme.ts
        │
        ├── types/             # TypeScript 类型定义
        │   ├── project.ts
        │   ├── request.ts
        │   └── environment.ts
        │
        └── assets/            # 样式、图片
            └── styles/
                └── main.css
```

## 7. 数据库设计

### 7.1 projects

| 列 | 类型 | 说明 |
|---|---|---|
| id | TEXT (UUID) | 主键 |
| name | TEXT | 项目名称 |
| description | TEXT | 描述 |
| created_at | DATETIME | |
| updated_at | DATETIME | |

### 7.2 collections

| 列 | 类型 | 说明 |
|---|---|---|
| id | TEXT (UUID) | 主键 |
| project_id | TEXT | 外键 → projects |
| parent_id | TEXT | 自引用，支持树形 |
| name | TEXT | 集合名称 |
| sort_order | INTEGER | 排序 |
| created_at | DATETIME | |
| updated_at | DATETIME | |

### 7.3 requests

| 列 | 类型 | 说明 |
|---|---|---|
| id | TEXT (UUID) | 主键 |
| collection_id | TEXT | 外键 → collections |
| name | TEXT | 接口名称 |
| method | TEXT | GET/POST/PUT/DELETE... |
| url | TEXT | |
| headers | TEXT (JSON) | |
| params | TEXT (JSON) | |
| body | TEXT (JSON) | body_type, content |
| auth | TEXT (JSON) | auth_type, credentials |
| script | TEXT | 测试脚本 |
| sort_order | INTEGER | |
| created_at | DATETIME | |
| updated_at | DATETIME | |

### 7.4 environments

| 列 | 类型 | 说明 |
|---|---|---|
| id | TEXT (UUID) | 主键 |
| project_id | TEXT | 外键 → projects |
| name | TEXT | 环境名称（开发/测试/生产） |
| variables | TEXT (JSON) | [{key, value, enabled}] |
| is_active | BOOLEAN | 当前激活的环境 |
| created_at | DATETIME | |
| updated_at | DATETIME | |

### 7.5 history

| 列 | 类型 | 说明 |
|---|---|---|
| id | TEXT (UUID) | 主键 |
| project_id | TEXT | 外键 → projects |
| request_id | TEXT | 外键 → requests (可为空) |
| method | TEXT | |
| url | TEXT | |
| headers | TEXT (JSON) | |
| body | TEXT (JSON) | |
| response_status | INTEGER | |
| response_body | TEXT | |
| response_headers | TEXT (JSON) | |
| duration_ms | INTEGER | |
| created_at | DATETIME | |

## 8. UI 布局

采用 **三栏布局 + 多标签页**：

```
┌──────────────────────────────────────────────────┐
│ 标题栏          [环境选择器]            [设置]    │
├─────────┬────────────────────────────────────────┤
│         │  Tab1  Tab2  Tab3  +                   │
│ 集合树   ├────────────────────────────────────────┤
│ (左侧栏) │  Method  URL              [发送]      │
│ 240px    │━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━│
│          │  Params │ Headers │ Body │ Auth       │
│          │                                       │
│          │  [请求编辑区]                          │
│          │                                       │
│          ├━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━│
│          │  ● 200 OK    124ms                    │
│          │  [响应查看区]                          │
│          │                                       │
└─────────┴────────────────────────────────────────┘
```

- **左侧栏**：项目/集合树形列表，支持拖拽、搜索
- **中间标签栏**：每个请求一个 Tab，可关闭、可拖拽排序
- **请求编辑**：URL 行 + 选项卡（Params/Headers/Body/Auth）
- **底部响应**：状态码 + 耗时 + 响应体（格式化 JSON）+ 响应头

## 9. 实施顺序建议

1. **Phase 1 — 框架搭建**：Wails 项目初始化、SQLite 初始化、基础目录结构、前端路由和布局骨架
2. **Phase 2 — HTTP 请求调试**：集合管理、请求编辑、发送请求、查看响应
3. **Phase 3 — 环境变量**：多环境管理，变量插值
4. **Phase 4 — 多标签页**：Tab 管理、窗口状态
5. **Phase 5 — 历史记录**：历史列表、搜索、回放
6. **Phase 6 — WebSocket**：WebSocket 连接和消息
7. **Phase 7 — 自动化测试**：测试脚本引擎、批量运行
8. **Phase 8 — API 文档**：文档生成和导出
