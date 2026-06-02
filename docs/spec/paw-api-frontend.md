# Paw API — 前端设计

> 版本 1.0 | 2026-06-02

---

## 1. 路由

不需要 Vue Router。应用为单窗口工作台，所有视图通过组件条件渲染切换。

---

## 2. 组件树

```
App.vue
├── AppHeader.vue
│   ├── 项目名称（点击切换/管理项目）
│   ├── EnvSelector.vue（环境下拉切换）
│   ├── ThemeToggle.vue
│   └── 设置按钮（打开 SettingsModal）
│
├── AppBody.vue（左右分栏）
│   │
│   ├── Sidebar.vue（左侧 240px）
│   │   ├── CollectionTree.vue
│   │   │   └── TreeNode.vue（递归组件）
│   │   └── HistoryPanel.vue（侧栏底部 Tab：集合 / 历史）
│   │
│   └── Workspace.vue（右侧，上下分栏）
│       │
│       ├── TabBar.vue
│       ├── RequestPanel.vue
│       │   ├── UrlBar.vue（方法下拉 + URL 输入 + Send 按钮）
│       │   └── RequestSubTabs.vue
│       │       ├── ParamsEditor.vue
│       │       │   └── KeyValueTable.vue（共享组件）
│       │       ├── HeadersEditor.vue
│       │       │   └── KeyValueTable.vue + BulkEdit 切换
│       │       ├── BodyEditor.vue
│       │       │   └── 按 body_type 渲染不同编辑器
│       │       │       ├── FormDataEditor.vue
│       │       │       ├── UrlEncodedEditor.vue
│       │       │       ├── RawEditor.vue
│       │       │       │   └── CodeEditor.vue（语法高亮）
│       │       │       └── BinaryEditor.vue
│       │       └── AuthEditor.vue
│       │           └── 按 auth.type 渲染不同表单
│       │
│       └── ResponsePanel.vue
│           ├── StatusBar.vue（状态码 + 耗时 + 大小）
│           └── ResponseSubTabs.vue
│               ├── ResponseBody.vue（Pretty / Raw / Preview）
│               ├── ResponseHeaders.vue（键值对表格，只读）
│               ├── CookiesPanel.vue
│               └── RequestLog.vue（原始报文 + 复制 cURL）
```

### 全局弹窗/模态框

```
├── ProjectModal.vue（新建/编辑项目）
├── EnvManagerModal.vue（环境管理）
├── ImportModal.vue（导入选择器）
├── ExportModal.vue（导出格式选择）
├── SettingsModal.vue（设置面板）
└── DocsPreviewModal.vue（API 文档预览）
```

---

## 3. WebSocket 面板

WebSocket 不使用 TabBar，直接从集合树右键打开，以独立面板替换 RequestPanel + ResponsePanel：

```
WebSocketPanel.vue
├── WsUrlBar.vue（URL 输入 + Connect/Disconnect）
├── WsHeadersEditor.vue（连接附带的 Headers）
├── MessagePane.vue（收发消息历史）
├── MessageInput.vue（输入框 + Send 按钮）
└── WsStatusBar.vue（连接状态指示）
```

### 面板切换规则

```
[HTTP 模式]                           [WebSocket 模式]
┌─────────────────┐                   ┌─────────────────┐
│ TabBar           │                   │ 隐藏 TabBar       │
│ RequestPanel     │  ── 打开 WS ──→   │ WebSocketPanel   │
│ ResponsePanel    │                   │                  │
└─────────────────┘  ←── 关闭/断开 ──  └─────────────────┘
```

- 打开 WebSocket 时隐藏 TabBar，保持已打开的 HTTP Tab 不关闭
- 关闭或断开 WebSocket 连接后，恢复 HTTP TabBar 和之前活跃的 Tab
- 同一时间只有一个 WebSocket 连接存在

---

## 4. 数据流

```
+--------------------------------------------+
|  Vue 组件                                   |
|    |  dispatch / commit                     |
|    v                                        |
|  Pinia Store                                |
|    |  call                                  |
|    v                                        |
|  Wails Binding (Go)                        |
|    +-> handlers -> services -> repo -> SQLite  |
|    <- 返回结果                                |
|  Store 更新 -> 组件响应                       |
+--------------------------------------------+
```

- 所有数据库操作走 Go 后端，前端不直接读写 SQLite
- 前端发起 HTTP/WebSocket 请求也走 Go 后端
- 前端只通过 Pinia Store 调用 `@/wailsjs/go/main/App` 暴露的方法

---

## 5. 组件详情

### 5.1 KeyValueTable.vue（共享组件）

被 Params、Headers、FormData、UrlEncoded 等复用。

Props：
- `items: KvItem[]`（key, value, description, enabled）
- `presets?: { label, key, value? }[]`（Headers 的预设下拉）
- `showBulkEdit?: boolean`（是否显示 Bulk Edit 切换）
- `showType?: boolean`（form-data 的 text/file 列）

Events：
- `update:items`
- `bulkEdit`（切换至纯文本模式时）

### 5.2 CodeEditor.vue

基于 CodeMirror 6 或 Monaco Editor 封装。

Props：
- `modelValue: string`
- `language: string`（json / xml / html / javascript / text）
- `readonly?: boolean`

Events：`update:modelValue`

功能：语法高亮、行号、括号匹配、Beautify 按钮（JSON 模式）

### 5.3 TreeNode.vue

递归组件，渲染集合树节点。

Props：
- `node: TreeItem`（id, name, type, children）
- `depth: number`

功能：展开/折叠、右键菜单（新建请求、新建集合、重命名、删除、复制、导出）、拖拽（vuedraggable 或原生 Drag API）

### 5.4 TabBar.vue

状态来源：`tabsStore`

功能：渲染已打开的请求 Tab，支持关闭（x / 中键 / Ctrl+W）、右键菜单（关闭当前/其他/右侧/全部）、拖拽排序。

未保存修改显示圆点标识。

### 5.5 UrlBar.vue

功能：方法下拉 + URL 输入 + 发送按钮。URL 中 `{{变量}}` 高亮渲染，悬停显示解析值。Ctrl+Enter 触发发送。

### 5.6 ResponseBody.vue

三种模式切换：
- Pretty：CodeEditor（readonly），JSON/XML 自动格式化
- Raw：等宽字体 `<pre>` 展示原始文本
- Preview：iframe 渲染 HTML（仅 Content-Type: text/html）

### 5.7 RequestLog.vue

展示实际发出的请求报文（拼接变量后的最终值），格式：

```
GET /api/user?id=123 HTTP/1.1
Host: example.com
Authorization: Bearer xxx
```

一键复制为 cURL 命令。

---

## 6. Pinia Store

| Store | 职责 | 核心字段 |
|---|---|---|
| `projectStore` | 项目管理 | `projects[]`, `currentId` |
| `collectionStore` | 集合树 CRUD | `tree[]`, `expandedKeys` |
| `requestStore` | 请求编辑与发送 | `currentRequest`, `isSending`, `response` |
| `tabsStore` | Tab 状态 | `tabs[]`, `activeTabId` |
| `envStore` | 环境与变量 | `environments[]`, `activeEnvId`, `resolvedVars` |
| `historyStore` | 历史列表 | `items[]`, `searchKeyword`, `filter` |
| `settingsStore` | 全局设置 | `settings: Record<string, any>` |
| `themeStore` | 主题 | `mode`, `color` |

### Store 数据加载时机

- `projectStore`：应用启动时加载列表，选中上次使用的项目
- `collectionStore`：选中项目后加载
- `envStore`：选中项目后加载
- `historyStore`：按需加载（切换到历史面板时）
- `settingsStore` / `themeStore`：应用启动时加载

---

## 7. 视图状态

### 7.1 Empty State

- 无项目：主区域显示"创建项目"引导
- 无请求打开：RequestPanel 显示"点击集合中的请求开始调试"
- 未发送过请求：ResponsePanel 显示占位提示
- 历史为空：HistoryPanel 显示"暂无历史记录"

### 7.2 Loading State

- 发送请求时：Send 按钮变为 loading 状态，显示"取消"按钮
- 切换项目时：集合树显示骨架屏

### 7.3 Error State

- 请求失败（网络错误、超时）：ResponsePanel 显示错误信息
- 导入格式不支持：ImportModal 显示友好错误

---

## 8. 目录结构

```
frontend/src/
├── main.ts
├── App.vue
│
├── components/
│   ├── layout/
│   │   ├── AppHeader.vue
│   │   ├── AppBody.vue
│   │   ├── Sidebar.vue
│   │   └── Workspace.vue
│   │
│   ├── request/
│   │   ├── UrlBar.vue
│   │   ├── RequestPanel.vue
│   │   ├── RequestSubTabs.vue
│   │   ├── ParamsEditor.vue
│   │   ├── HeadersEditor.vue
│   │   ├── BodyEditor.vue
│   │   ├── AuthEditor.vue
│   │   ├── FormDataEditor.vue
│   │   ├── UrlEncodedEditor.vue
│   │   ├── RawEditor.vue
│   │   └── BinaryEditor.vue
│   │
│   ├── response/
│   │   ├── ResponsePanel.vue
│   │   ├── StatusBar.vue
│   │   ├── ResponseSubTabs.vue
│   │   ├── ResponseBody.vue
│   │   ├── ResponseHeaders.vue
│   │   ├── CookiesPanel.vue
│   │   └── RequestLog.vue
│   │
│   ├── collection/
│   │   ├── CollectionTree.vue
│   │   └── TreeNode.vue
│   │
│   ├── tabs/
│   │   └── TabBar.vue
│   │
│   ├── history/
│   │   └── HistoryPanel.vue
│   │
│   ├── websocket/
│   │   ├── WebSocketPanel.vue
│   │   ├── WsUrlBar.vue
│   │   ├── WsMessagePane.vue
│   │   └── WsMessageInput.vue
│   │
│   ├── environment/
│   │   ├── EnvSelector.vue
│   │   └── EnvManagerModal.vue
│   │
│   ├── modals/
│   │   ├── ProjectModal.vue
│   │   ├── ImportModal.vue
│   │   ├── ExportModal.vue
│   │   ├── SettingsModal.vue
│   │   └── DocsPreviewModal.vue
│   │
│   └── shared/
│       ├── KeyValueTable.vue
│       ├── CodeEditor.vue
│       └── EmptyState.vue
│
├── stores/
│   ├── project.ts
│   ├── collection.ts
│   ├── request.ts
│   ├── tabs.ts
│   ├── env.ts
│   ├── history.ts
│   ├── settings.ts
│   └── theme.ts
│
├── composables/
│   ├── useRequest.ts
│   ├── useVariables.ts
│   ├── useDragDrop.ts
│   ├── useImportExport.ts
│   └── useTheme.ts
│
├── types/
│   ├── project.ts
│   ├── collection.ts
│   ├── request.ts
│   ├── environment.ts
│   ├── history.ts
│   └── settings.ts
│
└── assets/
    └── styles/
        └── main.css
```

---

## 9. TypeScript 核心类型

```typescript
// collection.ts
interface TreeItem {
  id: number
  name: string
  type: 'folder' | 'request'
  method?: string       // 仅 type='request'
  url?: string          // 仅 type='request'
  children: TreeItem[]
  sort_order: number
}
```