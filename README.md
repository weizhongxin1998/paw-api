# Paw API

> A local API debugging & management tool — like Postman/Apifox but runs entirely on your desktop. No accounts, no cloud.
> 本地 API 调试与管理工具，类似 Postman/Apifox，纯桌面运行，无账号无云端。

---

## Features / 功能特性

| English | 中文 |
|---------|------|
| HTTP request editor (GET/POST/PUT/DELETE/PATCH/HEAD/OPTIONS) | HTTP 请求编辑器 |
| Query params, headers, body (JSON / Text / Form-Data / URL-Encoded / Binary) | 多种 Body 类型 |
| Auth: None / Basic / Bearer / Digest / OAuth 2.0 | 五种认证方式 |
| Collection tree management with drag & drop | 集合树管理，支持拖拽排序 |
| Multi-tab workspace | 多标签页工作区 |
| Environment variables with `{{variable}}` interpolation | 环境变量及插值解析 |
| Response viewer with JSON formatting | 响应查看器（JSON 格式化） |
| WebSocket client | WebSocket 客户端 |
| Test runner with assertions (status / body / JSONPath / header / duration) | 测试运行器 + 断言系统 |
| Code generation (cURL / JS fetch / Python requests / Go net/http) | 代码生成 |
| Postman v2.1 & OpenAPI 3.0 import/export | 导入导出（Postman / Swagger / cURL） |
| Global search (Ctrl+K) | 全局搜索（Ctrl+K） |
| History with replay | 请求历史与回放 |
| Header presets | 请求头快捷预设 |
| Request timeout & redirect configuration | 超时/重定向配置 |
| Request collection as API docs + HTML export | API 文档导出 |
| Dark mode & 3 theme colors | 深色模式 + 3 种主题色 |
| i18n: English / 中文 | 中英文国际化 |

---

## Screenshots / 截图

*(Coming soon — run `wails dev` to see it live)*

---

## Architecture / 架构

```
Frontend (Vue 3 + TypeScript)
    ↕ Wails Binding
Handlers (Go) → Services → Repositories
    ↕
SQLite (modernc.org/sqlite, no CGO)
```

### Project Structure / 目录结构

```
paw-api/
├── app.go                  # Wails app entry, handler registration
├── main.go                 # Binary entry
├── wails.json              # Wails config
├── go.mod / go.sum         # Go dependencies
│
├── database/               # SQLite init & migrations
├── models/                 # Go data models
├── repositories/           # Database access layer
├── services/               # Business logic
│   ├── request_service.go  #   HTTP send + auth handling
│   ├── assert_service.go   #   Assertion engine
│   ├── codegen_service.go  #   Code generation
│   ├── importer.go         #   Postman/Swagger/cURL import
│   ├── exporter.go         #   Postman/Swagger export
│   └── cookie_service.go   #   Cookie management
├── handlers/               # Wails-bound handlers
├── pkg/httpclient/         # HTTP & WebSocket client
│
├── frontend/
│   ├── src/
│   │   ├── components/     # Reusable Vue components
│   │   ├── views/          # Route pages
│   │   ├── stores/         # Pinia stores
│   │   ├── composables/    # Vue composables
│   │   ├── router/         # Vue Router config
│   │   ├── i18n/           # en.ts / zh-CN.ts
│   │   └── types/          # TypeScript interfaces
│   ├── wailsjs/go/         # Auto-generated Wails bindings
│   └── package.json
│
└── docs/
    ├── spec/               # Design specs
    └── plans/              # Implementation plans
```

---

## Getting Started / 快速开始

### Prerequisites / 前置要求

- [Go](https://go.dev/) 1.23+
- [Node.js](https://nodejs.org/) 18+
- [Wails v2](https://wails.io/) CLI

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### Development / 开发

```bash
# Clone & enter
git clone <repo-url> && cd paw-api

# Install frontend deps
cd frontend && npm install && cd ..

# Live dev (Go backend + Vite HMR frontend)
wails dev
```

### Build / 构建

```bash
wails build
# Output: build/bin/paw-api (or paw-api.exe on Windows)
```

### Frontend dev only / 仅前端开发

```bash
cd frontend
npm run dev      # Vite dev server
npm run build    # Type-check + production build
```

---

## Tech Stack / 技术栈

**Backend:** Go 1.23, Wails v2, SQLite (modernc.org/sqlite, no CGO), Gorilla WebSocket

**Frontend:** Vue 3 (Composition API + `<script setup>`), TypeScript, Pinia, Vue Router, Naive UI, vue-i18n, Vite

---

## i18n / 国际化

The app supports English and 中文 (Chinese). Switch via sidebar Settings panel or click the language toggle.

应用支持英文和中文，通过侧栏设置面板或语言切换按钮切换。

---

## License / 许可证

MIT
