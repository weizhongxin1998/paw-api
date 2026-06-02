# Paw API — 需求规格文档

> 版本 1.0 | 2026-06-02

---

## 1. 项目概述

### 1.1 产品定位

本地桌面端 API 调试与管理工具，对标 Postman / Apifox 的核心调试体验。完全离线运行，不依赖任何在线服务。

### 1.2 技术栈

| 层 | 技术 |
|---|---|
| 桌面框架 | Wails v2 |
| 后端 | Go 1.23 |
| 前端 | Vue 3 + TypeScript |
| 状态管理 | Pinia |
| 数据存储 | SQLite（modernc.org/sqlite，纯 Go，无 CGO） |
| UI 组件库 | Naive UI |
| 图标 | @vicons/ionicons5 |
| 构建 | Vite |

### 1.3 明确包含

- HTTP 请求调试（全方法、全参数类型）
- 集合（Collection）管理：树形分组、拖拽排序
- 多标签页：多请求并行编辑
- 环境变量：多环境切换，`{{var}}` 插值
- 历史记录：自动保存、搜索、回放
- WebSocket 调试：连接、收发消息
- API 文档：生成与导出（Markdown / HTML）
- 导入/导出：Postman Collection v2.x、OpenAPI 3.x / Swagger 2.x

### 1.4 明确排除

- 账号体系、登录注册
- 云端同步、团队协作
- Mock 服务

---

## 2. HTTP 请求调试

### 2.1 请求方法

- 内置：`GET`、`POST`、`PUT`、`DELETE`、`PATCH`、`HEAD`、`OPTIONS`
- 下拉选择，同时支持手动输入自定义方法

### 2.2 Params（Query 参数）

- 键值对表格：Key、Value、Description、Enabled 复选框
- 支持逐行启用/禁用
- 自动拼接到请求 URL
- Value 支持变量插值 `{{var}}`

### 2.2a Path Variables（路径参数）

- URL 中 `:paramName` 或 `{paramName}` 语法自动识别为路径参数
- 在 Params 下方自动展开 Path Variables 区域
- 表格格式：Key（自动填充，不可编辑）、Value、Description
- 支持变量插值

### 2.3 Headers

- 键值对表格：Key、Value、Description、Enabled
- 常用 Header 预设下拉（Content-Type、Authorization、Accept 等），选中自动填充
- 支持 **Bulk Edit** 模式：切换为纯文本编辑，按 `Key: Value` 格式，每行一条
- Value 支持变量插值

### 2.4 Body

| 类型 | 说明 |
|---|---|
| none | 无请求体（GET / HEAD 等默认） |
| form-data | 文本字段 + 文件字段（可本地上传）；每行可设类型 text/file |
| x-www-form-urlencoded | 键值对 URL 编码 |
| raw | 自由文本，下拉切换子类型：JSON / XML / HTML / Text / JavaScript |
| binary | 选择本地文件直接作为请求体 |

- raw 模式下提供**代码编辑器**（语法高亮、行号、括号匹配）
- JSON 模式下提供 **Beautify** 一键格式化
- Body 内容支持变量插值

### 2.5 Auth

| 类型 | 参数 |
|---|---|
| No Auth | — |
| Bearer Token | Token 输入框，支持变量插值 |
| Basic Auth | Username + Password，自动编码 Base64 |
| API Key | Key + Value + Add to（Header / Query Params） |

### 2.6 发送请求

- URL 输入框 + **Send** 按钮（主色调醒目按钮）
- 快捷键：`Ctrl + Enter`
- **不自动保存**，发送的是当前编辑中的值（对标 Postman/Apifox）
- 请求前自动替换环境变量
- 请求中显示加载状态，**支持取消请求**
- 支持设置请求超时
- 未保存修改的 Tab 显示圆点，关闭时提示

### 2.7 响应查看

- **状态行**：HTTP 状态码（颜色标识：2xx 绿、3xx 蓝、4xx 橙、5xx 红）+ 耗时（ms）+ 响应体大小
- **Body**：
  - **Pretty** 模式：自动格式化 JSON/XML，语法高亮，支持节点折叠/展开
  - **Raw** 模式：原始响应体，等宽字体
  - **Preview** 模式：HTML 内嵌 iframe 渲染预览（仅 HTML 响应类型）
- **Headers**：响应头键值对表格
- **Cookies**：解析 Set-Cookie 响应头，以表格展示
- **请求日志**：展示实际发送的原始请求报文（请求行 + Headers + Body）；支持一键复制为 cURL 命令

---

## 3. 集合管理（Collections）

### 3.1 数据结构

```
Project
├── Default Collection    (新建项目自动创建)
├── Collection A          (文件夹/分组)
│   ├── GET /api/user     (请求)
│   └── POST /api/login
├── Collection B
│   ├── GET /api/orders
│   └── Sub-Collection    (嵌套子分组)
│       ├── PUT /api/item/1
│       └── DELETE /api/item/1
```

- 集合支持无限层级嵌套
- 集合和请求均可通过**拖拽**调整顺序或移动归属
- 新建项目时自动创建一个默认集合

### 3.2 集合操作

- **新建请求**：选中集合 → 右键或工具栏按钮 → 新建请求，选择方法并命名
- **新建集合**：右键菜单 / 顶部按钮 → 输入名称
- **重命名**：双击或右键 → Rename
- **删除**：右键 → Delete（确认对话框）
- **复制/剪切/粘贴**：支持跨集合复制请求或子集合
- **导出**：导出整个项目或选中集合为 Paw API 自有格式

### 3.3 导入/导出格式

| 格式 | 导入 | 导出 |
|---|---|---|
| Postman Collection v2.0 / v2.1 | ✅ | ✅ |
| OpenAPI 3.x | ✅ | — |
| Swagger 2.x | ✅ | — |
| Paw API 自有格式 (JSON) | ✅ | ✅ |

---

## 4. 多标签页（Tabs）

### 4.1 行为

- 双击集合树中的请求 → 开启新 Tab（如已打开则激活已有 Tab）
- Tab 显示：HTTP 方法 + 请求名称（过长则截断加省略号）
- 未保存修改时 Tab 名后显示圆点标识
- 支持 **关闭**（× 按钮 / 中键点击 / Ctrl+W）
- 支持 **右键关闭**：关闭当前、关闭其他、关闭右侧、关闭全部
- 切换 Tab 保留编辑状态
- 标签过多时水平滚动，不缩窄 Tab 宽度

---

## 5. 环境变量

### 5.1 环境管理

- 每个项目可创建多个环境（如：开发、测试、预发布、生产）
- 环境列表：名称 + 变量数量
- 当前激活环境用对勾标识，一键切换
- 支持新增、重命名、删除、复制环境

### 5.2 变量

- 每个环境包含一组变量，键值对格式：`variable_name` → `value`
- 支持启用/禁用单个变量
- Value 可引用其他变量（如 `{{base_url}}/api`）
- 引用语法：`{{variable_name}}`

### 5.3 插值范围

- 请求 URL
- Headers 值
- Params 值
- Body 内容（所有类型）
- Auth 凭证字段

### 5.4 环境变量高亮

- URL 输入框中 `{{变量}}` 使用特殊背景色高亮，方便识别
- 鼠标悬停显示解析后的实际值

---

## 6. 历史记录

### 6.1 自动记录

- 每次发送 HTTP 请求自动保存一条历史记录
- 记录内容：请求方法、URL、Headers、Body、响应状态码、响应体、响应头、耗时、时间戳

### 6.2 历史列表

- 独立视图 / 侧边栏面板
- 按时间倒序排列，显示：方法标签 + URL + 状态码 + 时间
- 支持搜索（URL 模糊匹配）
- 支持按 HTTP 方法、状态码范围筛选

### 6.3 回放

- 点击历史记录 → 在新 Tab 中打开，自动填充请求参数
- 支持一键重新发送

### 6.4 清理

- 支持清空全部历史
- 可设置保留天数（如 30 天后自动删除）

---

## 7. WebSocket 调试

### 7.1 连接

- URL 输入框：`ws://` 或 `wss://` 协议
- **Connect** / **Disconnect** 按钮
- 连接状态指示：绿色（已连接）、红色（已断开）、黄色（连接中）
- 显示连接时间

### 7.2 消息收发

- 文本框输入消息内容，支持 JSON 语法高亮
- **Send** 按钮 + `Ctrl+Enter`
- 消息历史列表保留发送和接收的消息
- 消息类型标识（发送/接收）和时间戳

### 7.3 Headers

- 连接时可附加自定义 Headers 键值对

---

## 8. API 文档生成

### 8.1 生成来源

- 基于项目中的集合和请求自动生成文档
- 文档使用请求的 name、description（需请求支持描述字段）
- URL 中的路径参数自动识别

### 8.2 文档内容

- 接口分组（对应集合结构）
- 每个接口：方法、路径、参数、Headers、Body 示例、响应示例

### 8.3 导出格式

- **Markdown**：适合嵌入项目 README 或 Wiki
- **HTML**：独立页面，可离线浏览

### 8.4 预览

- 应用内提供文档预览视图，实时反映最新请求配置

---

## 9. 通用能力

### 9.1 快捷键

| 快捷键 | 功能 |
|---|---|
| `Ctrl + Enter` | 发送请求 |
| `Ctrl + S` | 保存当前请求 |
| `Ctrl + W` | 关闭当前 Tab |
| `Ctrl + Tab` | 切换下一个 Tab |
| `Ctrl + Shift + Tab` | 切换上一个 Tab |

### 9.2 主题

- 日间模式（默认）+ 夜间模式
- 日间模式支持多套颜色主题预设（绿、蓝、紫等）
- 通过标题栏入口切换
- 偏好本地持久化存储

### 9.3 项目

- 多项目支持：创建、切换、删除项目
- 启动时默认打开上次使用的项目
- 切换项目时关闭所有已打开 Tab（Tab 属于项目）
- 项目数据统一存储（单个 SQLite 文件多项目）

---

## 10. 设置（Settings）

### 10.1 通用（General）

| 配置项 | 说明 | 默认值 |
|---|---|---|
| 请求超时 | 超时时间（秒），0 表示不限制 | 30s |
| 跟随重定向 | 是否自动跟随 3xx 重定向 | 开 |
| 最大重定向次数 | 跟随跳转的上限 | 10 |
| SSL 证书验证 | 是否验证 HTTPS 证书；关闭后可访问自签名站点 | 开 |

### 10.2 代理（Proxy）

- **不使用代理**（默认）
- **系统代理**：使用操作系统代理设置
- **自定义代理**：手动填写 HTTP 代理地址 + 端口；支持为 HTTPS 单独配置

### 10.3 证书（Certificates）

- **CA 证书**：可添加 PEM 格式的 CA 根证书文件，用于信任自签名或企业 CA
- **客户端证书**：支持配置 PFX/PEM 客户端证书（含私钥），用于双向 TLS 认证

### 10.4 外观（Appearance）

| 配置项 | 说明 |
|---|---|
| 主题模式 | 日间 / 夜间 |
| 颜色主题 | 日间模式下的预设颜色方案（绿 / 蓝 / 紫等） |
| 字体大小 | 编辑器字体大小（12-24px） |

### 10.5 数据（Data）

| 操作 | 说明 |
|---|---|
| 备份全部数据 | 将所有项目与设置归档为备份文件（.pawb，SQLite 文件副本） |
| 从备份恢复 | 从 .pawb 文件恢复（覆盖现有数据） |
| 清空历史 | 删除所有历史记录（保留项目和环境配置） |
| 重置应用 | 清除所有本地数据，恢复初始状态 |

---

## 11. 架构约束

### 11.1 架构分层

```
┌──────────────────────────────────┐
│  Vue 3 前端                       │
│  Views → Stores → Composables    │
├──────────────────────────────────┤
│  Wails Binding                   │
├──────────────────────────────────┤
│  Go 后端                          │
│  Handlers → Services → Repos → DB│
└──────────────────────────────────┘
```

- **Handlers** 层：Wails 绑定入口，入参校验 + 数据转换
- **Services** 层：全部业务逻辑
- **Repositories** 层：数据访问，仅做 CRUD
- Handler 不做业务逻辑，Repository 不做业务逻辑

### 11.2 HTTP 客户端

- Go 后端执行 HTTP 请求（不从前端发起，避免跨域/证书问题）
- 支持超时设置、重定向跟随策略
- 支持自签名证书（配置开关）

### 11.3 变量替换

- 请求发送前，由 Service 层统一执行变量替换
- 替换顺序：嵌套变量需多次扫描直到稳定

---

## 12. 数据库设计

> 详见 [paw-api-database.md](./paw-api-database.md)

---

## 13. 实施阶段

> 前端设计详见 [paw-api-frontend.md](./paw-api-frontend.md)
> 后端设计详见 [paw-api-backend.md](./paw-api-backend.md)

| Phase | 内容 | 依赖 |
|---|---|---|
| 1 | 框架搭建：Wails 项目初始化、SQLite 建表、Go 包结构、前端布局骨架 | — |
| 2 | HTTP 请求调试（核心）：请求编辑、发送、响应查看、集合 CRUD | Phase 1 |
| 3 | 导入/导出：Postman Collection、OpenAPI/Swagger | Phase 2 |
| 4 | 环境变量：多环境管理、`{{var}}` 插值、环境切换 | Phase 1 |
| 5 | 多标签页：Tab 状态管理、拖拽排序 | Phase 2 |
| 6 | 历史记录：自动记录、搜索、回放、清理 | Phase 2 |
| 7 | WebSocket 调试：连接管理、消息收发 | Phase 2 |
| 8 | API 文档：自动生成、Preview、导出 Markdown/HTML | Phase 2 |
