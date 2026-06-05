# Paw API 国际化 (i18n) 实施计划

## 项目概况

Paw API 是一个基于 Wails v2 (Go + Vue 3) 的桌面 API 调试工具。前端使用 Vue 3 + Naive UI + Pinia，目前没有任何国际化基础设施，所有用户可见文本均为硬编码（约 85% 中文、15% 英文）。

**扫描结果：** 25 个前端文件、约 345+ 条硬编码字符串需要国际化。Go 后端无用户可见字符串，无需改动。

---

## 技术选型

| 项目 | 选择 | 理由 |
|------|------|------|
| i18n 库 | **vue-i18n v9** | Vue 3 官方推荐，Composition API 友好，社区成熟 |
| Naive UI 国际化 | **NConfigProvider locale 属性** | Naive UI 内置 `zhCN` / `enUS` 语言包，只需动态绑定 |
| 语言存储 | **settings store 增加 `locale` 字段** | 复用现有 key-value 设置存储，无需改后端 schema |
| 翻译文件格式 | **JSON（扁平 + 嵌套混合）** | 简洁直观，便于维护 |

---

## 目录结构设计

```
frontend/src/
├── i18n/
│   ├── index.ts          # vue-i18n 初始化 + 语言切换逻辑
│   ├── locales/
│   │   ├── zh-CN.json    # 简体中文翻译
│   │   └── en-US.json    # 英文翻译
│   └── naive-locales.ts  # Naive UI locale 映射
```

---

## 翻译文件 Key 设计

采用按功能模块分组，嵌套不超过 3 层：

```json
{
  "common": {
    "confirm": "确定",
    "cancel": "取消",
    "delete": "删除",
    "save": "保存",
    "close": "关闭",
    "create": "创建",
    "edit": "编辑",
    "copy": "复制",
    "rename": "重命名",
    "search": "搜索",
    "comingSoon": "即将推出"
  },
  "header": { ... },
  "sidebar": { ... },
  "workspace": { ... },
  "request": { ... },
  "response": { ... },
  "collection": { ... },
  "environment": { ... },
  "history": { ... },
  "settings": { ... },
  "import": { ... },
  "export": { ... },
  "docs": { ... },
  "websocket": { ... },
  "message": {
    "created": "已创建{name}",
    "saved": "已保存",
    "deleted": "已删除{name}",
    ...
  },
  "dialog": {
    "confirmClose": { "title": "确认关闭", "content": "未保存的修改将丢失..." },
    "confirmDelete": { "title": "确认删除", "content": "..." },
    ...
  },
  "time": {
    "justNow": "刚刚",
    "minutesAgo": "{n}分钟前",
    ...
  }
}
```

---

## 实施步骤

### 第一步：基础设施搭建（预计 1 个文件新建 + 2 个文件修改）

1. **安装 vue-i18n**
   ```bash
   cd frontend && npm install vue-i18n@9
   ```

2. **创建 `src/i18n/index.ts`**
   - 初始化 vue-i18n 实例
   - 导出 `i18n` 实例和 `t()` 便捷函数
   - 导出 `setLocale(locale)` 函数（同步更新 i18n locale + Naive UI locale + settings store）

3. **创建 `src/i18n/naive-locales.ts`**
   - 映射 `zh-CN` → `zhCN` (from `naive-ui/es/locales/common/zhCN`)
   - 映射 `en-US` → `enUS` (from `naive-ui/es/locales/common/enUS`)

4. **修改 `main.ts`**
   - 导入并注册 `i18n` 插件

5. **修改 `App.vue`**
   - `<n-config-provider>` 增加动态 `:locale` 和 `:date-locale` 绑定
   - 根据 settings store 中的 locale 值自动切换

### 第二步：Settings Store 扩展

1. **修改 `stores/settings.ts`**
   - `Settings` 接口增加 `locale: 'zh-CN' | 'en-US'` 字段
   - 默认值 `'zh-CN'`
   - 语言切换时调用 `setLocale()` 联动更新

### 第三步：创建翻译文件

1. **创建 `zh-CN.json` 和 `en-US.json`**
   - 从所有组件中提取约 345 条字符串
   - 按上述 Key 设计组织

### 第四步：逐组件改造（按优先级分批）

**批次 A — 核心布局（4 个文件，影响最大）**
- `components/layout/AppHeader.vue` (~15 条)
- `components/layout/Sidebar.vue` (~21 条)
- `components/layout/Workspace.vue` (~47 条)
- `components/layout/AppBody.vue` (~19 条)

**批次 B — 请求/响应面板（5 个文件）**
- `components/request/RequestPanel.vue` (~19 条)
- `components/request/UrlBar.vue` (~6 条)
- `components/request/AuthEditor.vue` (~30 条)
- `components/request/BodyEditor.vue` (~16 条)
- `components/response/ResponsePanel.vue` (~26 条)

**批次 C — 弹窗与模态框（4 个文件）**
- `components/modals/SettingsModal.vue` (~40 条)
- `components/modals/ImportModal.vue` (~22 条)
- `components/modals/ExportModal.vue` (~22 条)
- `components/modals/DocsPreviewModal.vue` (~16 条)

**批次 D — 集合/环境/历史（6 个文件）**
- `components/collection/CollectionTree.vue` (~5 条)
- `components/collection/TreeNode.vue` (~8 条)
- `components/environment/EnvSelector.vue` (~8 条)
- `components/environment/EnvManagerModal.vue` (~25 条)
- `components/history/HistoryPanel.vue` (~15 条)
- `components/websocket/WebSocketPanel.vue` (~8 条)

**批次 E — 共享组件 + 其他（3 个文件）**
- `components/shared/KeyValueTable.vue` (~12 条)
- `components/request/RequestDocsView.vue` (~8 条)
- `components/layout/ProjectHome.vue` (~32 条)

### 第五步：特殊处理

1. **Naive UI `message` / `dialog` API**
   - 14 个文件使用了 `message.success/error/info/warning`
   - 4 个文件使用了 `dialog.warning`
   - 所有调用点的文本参数替换为 `t('message.xxx')`

2. **相对时间函数**
   - `ProjectHome.vue` 中的 `relativeTime()` 函数（6 条中文时间字符串）
   - `HistoryPanel.vue` 中的 `formatTime()` 函数（4 条中文时间字符串）
   - 替换为基于 locale 的实现

3. **模板中的插值表达式**
   - 部分字符串使用了模板字面量（如 `` `已创建项目 "${name}"` ``）
   - 改为 vue-i18n 的命名参数：`t('message.projectCreated', { name })`

4. **Naive UI 内置组件文本**
   - `<n-result>`、`<n-empty>`、`<n-progress>` 等组件自带文本
   - 通过 `NConfigProvider` 的 locale 属性统一处理

### 第六步：Settings UI 添加语言切换入口

- 在 `SettingsModal.vue` 的"通用"分组中添加语言选择下拉框
- 选项：`简体中文` / `English`
- 切换后即时生效（无需重启）

---

## 改造模式示例

**改造前：**
```vue
<n-button @click="handleDelete">删除</n-button>
```

**改造后：**
```vue
<n-button @click="handleDelete">{{ $t('common.delete') }}</n-button>
```

**带参数的消息改造前：**
```ts
message.success(`已创建项目 "${name}"`)
```

**改造后：**
```ts
message.success(t('message.projectCreated', { name }))
```

**Naive UI locale 联动：**
```vue
<!-- App.vue -->
<n-config-provider :locale="naiveLocale" :date-locale="naiveDateLocale" ...>
```

---

## 风险与注意事项

1. **Wails 绑定方法中的文本**：部分 Go 方法返回的 error 信息可能在 JS 层展示给用户，需确认是否有此类情况。目前扫描未发现。

2. **动态拼接字符串**：如右键菜单 label 在 `computed` 或函数中构建，需确保 `t()` 调用位于响应式上下文中。

3. **Naive UI message/dialog 的创建时机**：这些 API 在 setup 中通过 `useMessage()` / `useDialog()` 获取实例，`t()` 函数可在同一上下文中直接使用。

4. **首次加载语言**：应在应用初始化时从 settings store 读取上次保存的 locale，如无记录则默认 `zh-CN`。

5. **翻译一致性**：同一概念（如"集合"）在所有文件中应使用同一个 key，避免翻译不一致。

---

## 工作量估算

| 阶段 | 预计时间 |
|------|----------|
| 基础设施搭建 | ~30 分钟 |
| 翻译文件编写 | ~2 小时 |
| 组件改造（25 个文件） | ~4-5 小时 |
| 测试与校对 | ~1 小时 |
| **合计** | **~8 小时** |

---

## 建议的实施顺序

建议从 **基础设施搭建 → 翻译文件 → Settings 语言切换 → 逐组件改造** 的顺序推进。每完成一个批次即可运行验证，确保不影响现有功能。
