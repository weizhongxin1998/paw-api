# Paw API — 数据库设计

> 版本 1.0 | 2026-06-02

---

## 1. 概述

- 数据库：SQLite（通过 `modernc.org/sqlite` 访问，无 CGO）
- 主键：雪花风格 ID（`INTEGER`，48 位有序整数，详见 [ID 生成方案](#5-id-生成方案)）
- 时间戳：ISO 8601 格式字符串（TEXT）
- 复杂结构字段：JSON 字符串（TEXT），Go 端序列化/反序列化
- 所有表使用 `IF NOT EXISTS` 创建，支持增量迁移

---

## 2. ER 图

```
┌──────────┐
│ projects │
└────┬─────┘
     │ 1
     │
     ├────────────┐
     │ *          │ *
┌────┴─────┐  ┌──┴────────────┐
│collections│  │ environments  │
└────┬─────┘  └──┬────────────┘
     │ 1         │ 1
     │           │
     │ *         │ *
┌────┴─────┐  ┌──┴──────────────┐
│ requests │  │ env_variables   │
└────┬─────┘  └─────────────────┘
     │ *
     │
┌────┴─────┐
│ history  │
└──────────┘
```

---

## 3. 表定义

### 3.1 projects

| 列 | 类型 | 约束 | 说明 |
|---|---|---|---|
| id | INTEGER | PK | 雪花 ID |
| name | TEXT | NOT NULL | 项目名称 |
| description | TEXT | DEFAULT '' | 项目描述 |
| created_at | TEXT | NOT NULL | 创建时间 |
| updated_at | TEXT | NOT NULL | 最后修改时间 |

```sql
CREATE TABLE IF NOT EXISTS projects (
    id          INTEGER PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at  TEXT NOT NULL,
    updated_at  TEXT NOT NULL
);
```

### 3.2 collections

| 列 | 类型 | 约束 | 说明 |
|---|---|---|---|
| id | INTEGER | PK | 雪花 ID |
| project_id | INTEGER | NOT NULL, FK → projects(id) | 所属项目 |
| parent_id | INTEGER | DEFAULT NULL | 父集合 ID，NULL 为根级 |
| name | TEXT | NOT NULL | 集合名称 |
| sort_order | INTEGER | NOT NULL DEFAULT 0 | 同级排序 |
| created_at | TEXT | NOT NULL | 创建时间 |
| updated_at | TEXT | NOT NULL | 最后修改时间 |

```sql
CREATE TABLE IF NOT EXISTS collections (
    id          INTEGER PRIMARY KEY,
    project_id  INTEGER NOT NULL,
    parent_id   INTEGER DEFAULT NULL,
    name        TEXT NOT NULL,
    sort_order  INTEGER NOT NULL DEFAULT 0,
    created_at  TEXT NOT NULL,
    updated_at  TEXT NOT NULL,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES collections(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_collections_project ON collections(project_id);
CREATE INDEX IF NOT EXISTS idx_collections_parent ON collections(parent_id);
```

### 3.3 requests

| 列 | 类型 | 约束 | 说明 |
|---|---|---|---|
| id | INTEGER | PK | 雪花 ID |
| collection_id | INTEGER | NOT NULL, FK → collections(id) | 所属集合 |
| name | TEXT | NOT NULL | 接口名称 |
| description | TEXT | DEFAULT '' | 接口描述 |
| method | TEXT | NOT NULL | 请求方法 |
| url | TEXT | NOT NULL DEFAULT '' | 请求 URL |
| headers | TEXT | DEFAULT '[]' | JSON: `[{key, value, description, enabled}]` |
| params | TEXT | DEFAULT '[]' | JSON: `[{key, value, description, enabled}]` |
| body_type | TEXT | NOT NULL DEFAULT 'none' | none / form-data / x-www-form-urlencoded / raw / binary |
| body | TEXT | DEFAULT '{}' | JSON: body_type 对应的内容结构 |
| auth | TEXT | DEFAULT '{"type":"none"}' | JSON: `{type, ...credentials}` |
| sort_order | INTEGER | NOT NULL DEFAULT 0 | 同级排序 |
| created_at | TEXT | NOT NULL | 创建时间 |
| updated_at | TEXT | NOT NULL | 最后修改时间 |

```sql
CREATE TABLE IF NOT EXISTS requests (
    id              INTEGER PRIMARY KEY,
    collection_id   INTEGER NOT NULL,
    name            TEXT NOT NULL,
    description     TEXT NOT NULL DEFAULT '',
    method          TEXT NOT NULL,
    url             TEXT NOT NULL DEFAULT '',
    headers         TEXT NOT NULL DEFAULT '[]',
    params          TEXT NOT NULL DEFAULT '[]',
    body_type       TEXT NOT NULL DEFAULT 'none',
    body            TEXT NOT NULL DEFAULT '{}',
    auth            TEXT NOT NULL DEFAULT '{"type":"none"}',
    sort_order      INTEGER NOT NULL DEFAULT 0,
    created_at      TEXT NOT NULL,
    updated_at      TEXT NOT NULL,
    FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_requests_collection ON requests(collection_id);
```

### 3.4 Body JSON 结构

**none**：
```json
{}
```

**form-data**：
```json
[
  { "key": "field1", "value": "val1", "type": "text", "description": "", "enabled": true },
  { "key": "file1", "value": "/path/to/file", "type": "file", "description": "", "enabled": true }
]
```

**x-www-form-urlencoded**：
```json
[
  { "key": "field1", "value": "val1", "description": "", "enabled": true }
]
```

**raw**：
```json
{
  "subType": "json",
  "content": "{\n  \"key\": \"value\"\n}"
}
```
> `subType` 可选值：json / xml / html / text / javascript

**binary**：
```json
{
  "fileName": "report.pdf",
  "filePath": "/local/path/to/report.pdf"
}
```

### 3.5 Auth JSON 结构

**No Auth**：
```json
{ "type": "none" }
```

**Bearer Token**：
```json
{ "type": "bearer", "token": "{{token_var}}" }
```

**Basic Auth**：
```json
{ "type": "basic", "username": "admin", "password": "{{pwd_var}}" }
```

**API Key**：
```json
{ "type": "apikey", "key": "X-API-Key", "value": "{{key_var}}", "addTo": "header" }
```
> `addTo` 可选值：header / query

### 3.6 environments

| 列 | 类型 | 约束 | 说明 |
|---|---|---|---|
| id | INTEGER | PK | 雪花 ID |
| project_id | INTEGER | NOT NULL, FK → projects(id) | 所属项目 |
| name | TEXT | NOT NULL | 环境名称 |
| is_active | INTEGER | NOT NULL DEFAULT 0 | 是否为当前激活环境 |
| created_at | TEXT | NOT NULL | 创建时间 |
| updated_at | TEXT | NOT NULL | 最后修改时间 |

```sql
CREATE TABLE IF NOT EXISTS environments (
    id          INTEGER PRIMARY KEY,
    project_id  INTEGER NOT NULL,
    name        TEXT NOT NULL,
    is_active   INTEGER NOT NULL DEFAULT 0,
    created_at  TEXT NOT NULL,
    updated_at  TEXT NOT NULL,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_environments_project ON environments(project_id);
```

### 3.7 env_variables

| 列 | 类型 | 约束 | 说明 |
|---|---|---|---|
| id | INTEGER | PK | 雪花 ID |
| environment_id | INTEGER | NOT NULL, FK → environments(id) | 所属环境 |
| key | TEXT | NOT NULL | 变量名 |
| value | TEXT | NOT NULL DEFAULT '' | 变量值 |
| enabled | INTEGER | NOT NULL DEFAULT 1 | 是否启用 |
| sort_order | INTEGER | NOT NULL DEFAULT 0 | 排序 |
| created_at | TEXT | NOT NULL | 创建时间 |

```sql
CREATE TABLE IF NOT EXISTS env_variables (
    id              INTEGER PRIMARY KEY,
    environment_id  INTEGER NOT NULL,
    key             TEXT NOT NULL,
    value           TEXT NOT NULL DEFAULT '',
    enabled         INTEGER NOT NULL DEFAULT 1,
    sort_order      INTEGER NOT NULL DEFAULT 0,
    created_at      TEXT NOT NULL,
    FOREIGN KEY (environment_id) REFERENCES environments(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_env_vars_env ON env_variables(environment_id);
```

### 3.8 history

| 列 | 类型 | 约束 | 说明 |
|---|---|---|---|
| id | INTEGER | PK | 雪花 ID |
| project_id | INTEGER | NOT NULL, FK → projects(id) | 所属项目 |
| request_id | INTEGER | DEFAULT NULL, FK → requests(id) | 来源请求（可为空） |
| method | TEXT | NOT NULL | 请求方法 |
| url | TEXT | NOT NULL | 请求 URL |
| request_headers | TEXT | DEFAULT '[]' | JSON：发送的 Headers |
| request_body | TEXT | DEFAULT '' | 发送的 Body |
| response_status | INTEGER | NOT NULL | HTTP 状态码 |
| response_headers | TEXT | DEFAULT '[]' | JSON：响应 Headers |
| response_body | TEXT | DEFAULT '' | 响应体 |
| duration_ms | INTEGER | NOT NULL | 耗时（毫秒） |
| created_at | TEXT | NOT NULL | 请求时间 |

```sql
CREATE TABLE IF NOT EXISTS history (
    id                INTEGER PRIMARY KEY,
    project_id        INTEGER NOT NULL,
    request_id        INTEGER DEFAULT NULL,
    method            TEXT NOT NULL,
    url               TEXT NOT NULL,
    request_headers   TEXT NOT NULL DEFAULT '[]',
    request_body      TEXT NOT NULL DEFAULT '',
    response_status   INTEGER NOT NULL,
    response_headers  TEXT NOT NULL DEFAULT '[]',
    response_body     TEXT NOT NULL DEFAULT '',
    duration_ms       INTEGER NOT NULL,
    created_at        TEXT NOT NULL,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (request_id) REFERENCES requests(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_history_project ON history(project_id);
CREATE INDEX IF NOT EXISTS idx_history_created ON history(created_at);
```

### 3.9 settings

应用全局设置，单行 key-value 结构。

| 列 | 类型 | 约束 | 说明 |
|---|---|---|---|
| key | TEXT | PK | 设置键 |
| value | TEXT | NOT NULL | 设置值（JSON） |
| updated_at | TEXT | NOT NULL | 最后修改时间 |

```sql
CREATE TABLE IF NOT EXISTS settings (
    key         TEXT PRIMARY KEY,
    value       TEXT NOT NULL,
    updated_at  TEXT NOT NULL
);
```

**预设设置项**：

| key | value | 说明 |
|---|---|---|
| `general.timeout` | `30` | 请求超时（秒） |
| `general.follow_redirects` | `true` | 是否跟随重定向 |
| `general.max_redirects` | `10` | 最大重定向次数 |
| `general.ssl_verify` | `true` | 是否验证 SSL 证书 |
| `proxy.mode` | `"none"` | none / system / custom |
| `proxy.http` | `""` | HTTP 代理地址 |
| `proxy.https` | `""` | HTTPS 代理地址 |
| `theme.mode` | `"light"` | light / dark |
| `theme.color` | `"green"` | 颜色方案 |
| `theme.font_size` | `14` | 编辑器字体大小 |
| `history.retention_days` | `30` | 历史保留天数（0 = 永久） |
| `app.last_project_id` | `0` | 上次打开的项目 ID |

---

## 4. 级联删除规则

| 父表 | 子表 | 规则 |
|---|---|---|
| projects | collections | CASCADE |
| projects | environments | CASCADE |
| projects | history | CASCADE |
| collections | collections (parent_id) | SET NULL |
| collections | requests | CASCADE |
| environments | env_variables | CASCADE |
| requests | history (request_id) | SET NULL |

---

## 5. ID 生成方案

### 5.1 结构

采用雪花风格 48 位整数，既保证全局唯一又天然按时间排序：

```
|-- 36 bits 时间戳 --|-- 12 bits 序列号 --|
       毫秒级                    0~4095
```

- **时间戳**（36 位）：毫秒，起点自定义（2024-01-01），可用约 68 年
- **序列号**（12 位）：每毫秒 0~4095
- 总占用 48 位，`INTEGER`（SQLite 64 位）安全容纳
- 十进制最长 15 位（< 16），JSON 中直接输出数字，JS `Number`（53 位精度）无精度丢失

### 5.2 Go 实现

包路径：`pkg/snowflake/snowflake.go`

```go
package snowflake

import (
    "sync"
    "time"
)

// Epoch 自定义起始时间 (2024-01-01 00:00:00 UTC)
const epoch int64 = 1704067200000

// Generator 雪花 ID 生成器
type Generator struct {
    mu       sync.Mutex
    lastTime int64
    sequence int64
}

// New 创建生成器
func New() *Generator {
    return &Generator{}
}

// Next 生成下一个 ID（48 位，十进制最长 15 位）
func (g *Generator) Next() int64 {
    g.mu.Lock()
    defer g.mu.Unlock()

    now := time.Now().UnixMilli()

    if now == g.lastTime {
        g.sequence = (g.sequence + 1) & 0xFFF // 0~4095
        if g.sequence == 0 {
            for now <= g.lastTime {
                now = time.Now().UnixMilli()
            }
        }
    } else {
        g.sequence = 0
    }

    g.lastTime = now
    return ((now - epoch) << 12) | g.sequence
}
```

### 5.3 使用

```go
var snow = snowflake.New()

project := &models.Project{
    ID:   snow.Next(),
    Name: "My Project",
    // ...
}
```

### 5.4 导入导出 ID 处理

导入时，导入源（Postman / OpenAPI）的外部 ID 保留在额外字段或直接丢弃，数据行始终使用新生成的雪花 ID。导出时同理，雪花 ID 不暴露给外部格式。

