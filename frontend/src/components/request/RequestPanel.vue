<template>
  <div
    ref="panelRef"
    class="request-panel"
    tabindex="-1"
    @keydown="onPanelKeydown"
  >
    <!-- ── Sub Tabs Bar ── -->
    <div class="sub-tabs">
      <div class="sub-tabs-group">
        <button
          v-for="(tab, i) in tabs"
          :key="tab.key"
          :class="{ active: activeTab === tab.key }"
          @click="switchTab(tab.key)"
        >
          <span class="tab-label">{{ tab.label }}</span>
          <span v-if="tab.count != null && tab.count > 0" class="cnt" :class="{ 'cnt-accent': tab.count > 0 }">{{ tab.count }}</span>
          <span class="tab-shortcut">{{ i + 1 }}</span>
        </button>
      </div>
      <div style="flex:1"></div>
      <!-- Clear all button (params / headers) -->
      <button
        v-if="(activeTab === 'params' || activeTab === 'headers') && hasActiveItems"
        class="toolbar-btn clear-btn"
        title="清空全部"
        @click="onClearAll"
      >
        <svg viewBox="0 0 24 24" width="13" height="13" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <polyline points="3 6 5 6 21 6" />
          <path d="M19 6l-1 14a2 2 0 01-2 2H8a2 2 0 01-2-2L5 6" />
          <path d="M10 11v6" /><path d="M14 11v6" />
        </svg>
        <span>清空</span>
      </button>
      <!-- Bulk edit toggle (headers) -->
      <button
        v-if="activeTab === 'headers'"
        class="toolbar-btn"
        :class="{ 'toolbar-btn-active': isBulkEdit }"
        @click="isBulkEdit = !isBulkEdit"
      >
        <svg v-if="!isBulkEdit" viewBox="0 0 24 24" width="13" height="13" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" />
          <polyline points="14 2 14 8 20 8" />
          <line x1="16" y1="13" x2="8" y2="13" />
          <line x1="16" y1="17" x2="8" y2="17" />
        </svg>
        <svg v-else viewBox="0 0 24 24" width="13" height="13" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <rect x="3" y="3" width="18" height="18" rx="2" />
          <line x1="3" y1="9" x2="21" y2="9" />
          <line x1="9" y1="21" x2="9" y2="9" />
        </svg>
        <span>{{ isBulkEdit ? '表格' : '批量' }}</span>
      </button>
    </div>

    <!-- ── Content area ── -->
    <div class="sub-content" ref="contentRef">
      <!-- Params -->
      <div v-if="activeTab === 'params'" class="params-content">
        <KeyValueTable
          :items="paramsItems"
          :header-check="paramsEnabled"
          @update:items="onParamsChange"
          @update:header-check="(v: boolean) => { paramsEnabled = v; emit('update:paramsEnabled', v) }"
        />
      </div>

      <!-- Path Variables -->
      <div v-else-if="activeTab === 'path'" class="params-content">
        <div v-if="pathVariables.length > 0" class="path-section">
          <div class="path-section-hdr">
            <div class="path-section-title">
              <svg class="path-section-icon" viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M10 13a5 5 0 007.54.54l3-3a5 5 0 00-7.07-7.07l-1.72 1.71" />
                <path d="M14 11a5 5 0 00-7.54-.54l-3 3a5 5 0 007.07 7.07l1.71-1.71" />
              </svg>
              <span>路径变量</span>
            </div>
            <span class="path-section-hint">URL 中 <code>:param</code> 或 <code>{param}</code> 自动识别</span>
          </div>
          <div class="path-table-wrap">
            <table class="kvt">
              <thead>
                <tr>
                  <th class="kvt-th-key">变量名</th>
                  <th class="kvt-th-val">值</th>
                  <th class="kvt-th-desc">描述</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="pv in pathVariables" :key="pv.key" :class="{ 'kvt-row-warn': !pv.value }">
                  <td>
                    <div class="kvt-key-cell">
                      <span class="kvt-key-badge">{{ pv.key }}</span>
                    </div>
                  </td>
                  <td><input class="kvt-input" v-model="pv.value" placeholder="输入值..." @input="onPathVarChange" /></td>
                  <td><input class="kvt-input" v-model="pv.description" placeholder="描述 (可选)" @input="onPathVarChange" /></td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
        <div v-else class="path-empty">
          <svg class="path-empty-icon" viewBox="0 0 24 24" width="32" height="32" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
            <path d="M10 13a5 5 0 007.54.54l3-3a5 5 0 00-7.07-7.07l-1.72 1.71" />
            <path d="M14 11a5 5 0 00-7.54-.54l-3 3a5 5 0 007.07 7.07l1.71-1.71" />
          </svg>
          <span class="path-empty-title">未检测到路径变量</span>
          <span class="path-empty-hint">在 URL 中使用以下语法定义路径变量：</span>
          <div class="path-empty-examples">
            <code>GET /users/:userId/posts/:postId</code>
            <code>GET /items/{itemId}/detail</code>
          </div>
        </div>
      </div>

      <!-- Headers -->
      <KeyValueTable
        v-else-if="activeTab === 'headers'"
        :items="headersItems"
        :show-bulk-edit="false"
        :bulk-mode="isBulkEdit"
        @update:items="onHeadersChange"
      />

      <!-- Body -->
      <BodyEditor
        v-else-if="activeTab === 'body'"
        :body-type="bodyType"
        :body-data="bodyData"
        @update:body-type="v => { bodyType = v; sync() }"
        @update:body-data="v => { bodyData = v; sync() }"
      />

      <!-- Auth -->
      <AuthEditor
        v-else-if="activeTab === 'auth'"
        :model-value="authData"
        @update:model-value="v => { authData = v; sync() }"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import KeyValueTable from '../shared/KeyValueTable.vue'
import BodyEditor from './BodyEditor.vue'
import AuthEditor from './AuthEditor.vue'
import type { KvItem } from '../../types/request'

interface PathVariable { key: string; value: string; description: string }

const props = defineProps<{
  headers: string; params: string; bodyType: string; bodyData: string
  authData: string; url: string; pathVars: string; paramsEnabled: boolean
}>()

const emit = defineEmits<{
  (e: 'update:headers', v: string): void
  (e: 'update:params', v: string): void
  (e: 'update:bodyType', v: string): void
  (e: 'update:bodyData', v: string): void
  (e: 'update:authData', v: string): void
  (e: 'update:pathVars', v: string): void
  (e: 'update:paramsEnabled', v: boolean): void
}>()

const tabs = [
  { key: 'params', label: '参数', count: undefined as number | undefined },
  { key: 'path', label: '路径', count: undefined as number | undefined },
  { key: 'headers', label: '请求头', count: undefined as number | undefined },
  { key: 'body', label: '请求体' },
  { key: 'auth', label: '认证' },
]

const activeTab = ref('params')
const isBulkEdit = ref(false)
const paramsEnabled = ref(props.paramsEnabled)
const bodyType = ref(props.bodyType)
const bodyData = ref(props.bodyData)
const authData = ref(props.authData)
const paramsItems = ref<KvItem[]>([])
const headersItems = ref<KvItem[]>([])
const pathVariables = ref<PathVariable[]>([])
const panelRef = ref<HTMLElement | null>(null)
const contentRef = ref<HTMLElement | null>(null)

// ── Scroll position preservation ──
const scrollPositions = new Map<string, number>()

function switchTab(key: string) {
  if (key === activeTab.value) return
  // Save current scroll
  if (contentRef.value) {
    scrollPositions.set(activeTab.value, contentRef.value.scrollTop)
  }
  activeTab.value = key
  // Restore scroll for new tab
  nextTick(() => {
    if (contentRef.value) {
      contentRef.value.scrollTop = scrollPositions.get(key) ?? 0
    }
  })
}

// ── Keyboard shortcuts (Ctrl+1..5) ──
function onPanelKeydown(e: KeyboardEvent) {
  if (!e.ctrlKey || e.shiftKey || e.altKey) return
  const num = parseInt(e.key)
  if (num >= 1 && num <= 5) {
    e.preventDefault()
    const tab = tabs[num - 1]
    if (tab) switchTab(tab.key)
  }
}

// ── Has active items (for clear button) ──
const hasActiveItems = computed(() => {
  if (activeTab.value === 'params') {
    return paramsItems.value.some(i => i.key || i.value)
  }
  if (activeTab.value === 'headers') {
    return headersItems.value.some(i => i.key || i.value)
  }
  return false
})

function onClearAll() {
  const empty = [{ id: String(++iid), key: '', value: '', description: '', enabled: true }]
  if (activeTab.value === 'params') {
    paramsItems.value = empty
    tabs[0].count = undefined
    sync()
  } else if (activeTab.value === 'headers') {
    headersItems.value = empty
    tabs[2].count = 0
    sync()
  }
}

let iid = Date.now()

const pathVarPatterns = computed(() => {
  const url = props.url || ''
  const patterns: string[] = []
  const re = /:([a-zA-Z_]\w*)/g
  let match: RegExpExecArray | null
  while ((match = re.exec(url)) !== null) {
    if (!patterns.includes(match[1])) patterns.push(match[1])
  }
  const re2 = /\{([a-zA-Z_]\w*)\}/g
  while ((match = re2.exec(url)) !== null) {
    if (!patterns.includes(match[1])) patterns.push(match[1])
  }
  return patterns
})

watch(pathVarPatterns, (patterns) => {
  let stored: PathVariable[] = []
  try { stored = JSON.parse(props.pathVars || '[]') } catch {}
  const storedMap = new Map(stored.map(p => [p.key, p]))
  const existing = new Map(pathVariables.value.map(p => [p.key, p]))
  pathVariables.value = patterns.map(key => {
    const s = storedMap.get(key); const e = existing.get(key)
    return { key, value: s?.value || e?.value || '', description: s?.description || e?.description || '' }
  })
  tabs[1].count = patterns.length > 0 ? patterns.length : undefined
  syncPathVars()
}, { immediate: true })

function onPathVarChange() { syncPathVars() }
function syncPathVars() { emit('update:pathVars', JSON.stringify(pathVariables.value)) }

function parseKv(raw: string): KvItem[] {
  try {
    const items = JSON.parse(raw).map((i: any) => ({ ...i, id: i.id || String(++iid) }))
    const last = items[items.length - 1]
    if (!last || last.key || last.value || last.description) {
      items.push({ id: String(++iid), key: '', value: '', description: '', enabled: true })
    }
    return items
  } catch {
    return [{ id: String(++iid), key: '', value: '', description: '', enabled: true }]
  }
}

watch(() => props.params, (v) => { paramsItems.value = parseKv(v) }, { immediate: true })
watch(() => props.headers, (v) => { headersItems.value = parseKv(v) }, { immediate: true })
watch(() => props.bodyType, (v) => { bodyType.value = v })
watch(() => props.bodyData, (v) => { bodyData.value = v })
watch(() => props.authData, (v) => { authData.value = v })
watch(() => props.paramsEnabled, (v) => { paramsEnabled.value = v })

function sync() {
  emit('update:headers', JSON.stringify(headersItems.value))
  emit('update:params', JSON.stringify(paramsItems.value))
  emit('update:bodyType', bodyType.value)
  emit('update:bodyData', bodyData.value)
  emit('update:authData', authData.value)
}

function onHeadersChange(items: KvItem[]) {
  headersItems.value = items
  tabs[2].count = items.filter(i => i.enabled && i.key).length
  sync()
}

function onParamsChange(items: KvItem[]) {
  paramsItems.value = items
  tabs[0].count = items.filter(i => i.enabled && i.key).length
  sync()
}
</script>

<style scoped>
.request-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: var(--bg-base);
  min-height: 0;
  outline: none;
}

/* ══════════════════════════════════════════
   Sub Tabs
   ══════════════════════════════════════════ */
.sub-tabs {
  display: flex;
  align-items: center;
  border-bottom: 1px solid var(--border-primary);
  background: var(--bg-surface);
  padding: 0 10px;
  gap: 0;
  min-height: 38px;
}
.sub-tabs-group {
  display: flex;
  gap: 0;
}
.sub-tabs button {
  padding: 0 14px;
  height: 38px;
  font-size: var(--fs-sm);
  cursor: pointer;
  color: var(--text-muted);
  border: none;
  background: transparent;
  outline: none;
  font-family: var(--font-mono);
  font-weight: 500;
  letter-spacing: 0.02em;
  transition: color var(--transition), background var(--transition);
  position: relative;
  display: inline-flex;
  align-items: center;
  gap: 5px;
}
/* Gradient underline for active tab */
.sub-tabs button::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%; right: 50%;
  height: 2px;
  background: linear-gradient(90deg, var(--accent), var(--accent-hover));
  border-radius: 2px 2px 0 0;
  transition: left var(--transition), right var(--transition), opacity var(--transition);
  opacity: 0;
}
.sub-tabs button.active::after {
  left: 8px;
  right: 8px;
  opacity: 1;
}
/* Glow effect under active tab */
.sub-tabs button.active::before {
  content: '';
  position: absolute;
  bottom: -1px;
  left: 20%; right: 20%;
  height: 6px;
  background: var(--accent);
  filter: blur(8px);
  opacity: 0.25;
  pointer-events: none;
}
.sub-tabs button.active {
  color: var(--accent);
  font-weight: 600;
}
.sub-tabs button:hover:not(.active) {
  color: var(--text-secondary);
  background: var(--bg-hover);
  border-radius: var(--radius-sm) var(--radius-sm) 0 0;
}
.tab-label {
  pointer-events: none;
}
/* Shortcut hint */
.tab-shortcut {
  font-size: var(--fs-xs);
  color: var(--text-muted);
  font-family: var(--font-ui);
  font-weight: 600;
  opacity: 0;
  transition: opacity var(--transition-fast);
  margin-left: 2px;
  pointer-events: none;
}
.sub-tabs button:hover .tab-shortcut,
.sub-tabs button.active .tab-shortcut {
  opacity: 0.7;
}

/* ── Count badges ── */
.cnt {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: var(--fs-2xs);
  background: var(--bg-active);
  color: var(--text-muted);
  padding: 0 5px;
  border-radius: 10px;
  font-weight: 600;
  min-width: 16px;
  height: 16px;
  line-height: 1;
  transition: all var(--transition-fast);
  pointer-events: none;
}
.cnt-accent {
  background: var(--accent-soft);
  color: var(--accent);
}

/* ── Toolbar buttons ── */
.toolbar-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 9px;
  font-size: var(--fs-xs);
  font-family: var(--font-ui);
  font-weight: 500;
  color: var(--text-muted);
  background: transparent;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all var(--transition-fast);
  margin-left: 4px;
  white-space: nowrap;
}
.toolbar-btn:hover {
  color: var(--text-secondary);
  border-color: var(--border-hover);
  background: var(--bg-hover);
}
.toolbar-btn-active {
  color: var(--accent) !important;
  border-color: var(--accent) !important;
  background: var(--accent-soft) !important;
}
.clear-btn:hover {
  color: var(--red);
  border-color: var(--red);
  background: var(--red-soft);
}

/* ══════════════════════════════════════════
   Content
   ══════════════════════════════════════════ */
.sub-content {
  flex: 1;
  overflow-y: auto;
}
.params-content {
  padding: 12px;
}

/* ══════════════════════════════════════════
   Path Variables Section
   ══════════════════════════════════════════ */
.path-section {
  margin-bottom: 12px;
}
.path-section-hdr {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  padding: 0 2px;
}
.path-section-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: var(--fs-sm);
  color: var(--text-secondary);
  font-weight: 600;
  font-family: var(--font-ui);
}
.path-section-icon {
  color: var(--accent);
  opacity: 0.7;
}
.path-section-hint {
  font-size: var(--fs-sm);
  color: var(--text-muted);
  font-family: var(--font-ui);
}
.path-section-hint code {
  background: var(--bg-elevated);
  padding: 1px 5px;
  border-radius: var(--radius-xs);
  font-family: var(--font-mono);
  font-size: var(--fs-xs);
  border: 1px solid var(--border-primary);
  color: var(--accent);
}

/* Path variables table */
.path-table-wrap {
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  overflow: hidden;
}
.kvt {
  width: 100%;
  border-collapse: separate;
  border-spacing: 0;
}
.kvt th {
  text-align: left;
  padding: 7px 10px;
  font-size: var(--fs-xs);
  color: var(--text-muted);
  text-transform: uppercase;
  border-bottom: 1px solid var(--border-primary);
  font-weight: 600;
  background: var(--bg-elevated);
  letter-spacing: 0.04em;
  font-family: var(--font-ui);
}
.kvt-th-key { width: 25%; }
.kvt-th-val { width: 35%; }
.kvt-th-desc { width: 40%; }
.kvt td {
  padding: 3px 6px;
}
.kvt tr:not(:last-child) td {
  border-bottom: 1px solid var(--border-subtle);
}
.kvt tr:hover td {
  background: var(--bg-hover);
}
.kvt-row-warn td {
  background: var(--amber-soft) !important;
}
.kvt-row-warn:hover td {
  background: var(--amber-soft) !important;
}
.kvt-key-cell {
  display: flex;
  align-items: center;
  padding: 2px 0;
}
.kvt-key-badge {
  display: inline-flex;
  align-items: center;
  font-family: var(--font-mono);
  font-size: var(--fs-sm);
  font-weight: 600;
  color: var(--accent);
  background: var(--accent-soft);
  padding: 2px 8px;
  border-radius: var(--radius-xs);
  letter-spacing: 0.01em;
}
.kvt-input {
  width: 100%;
  padding: 6px 8px;
  border: 1px solid transparent;
  font-family: var(--font-mono);
  font-size: var(--fs-sm);
  background: transparent;
  border-radius: var(--radius-xs);
  outline: none;
  color: var(--text-primary);
  transition: border-color var(--transition), background var(--transition);
}
.kvt-input::placeholder {
  color: var(--text-placeholder);
}
.kvt-input:hover {
  border-color: var(--border-hover);
  background: var(--bg-surface);
}
.kvt-input:focus {
  border-color: var(--accent);
  background: var(--bg-surface);
}

/* ── Path empty state ── */
.path-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 36px 24px;
  text-align: center;
}
.path-empty-icon {
  color: var(--text-muted);
  opacity: 0.25;
  margin-bottom: 4px;
}
.path-empty-title {
  font-size: var(--fs-sm);
  font-weight: 600;
  color: var(--text-secondary);
  font-family: var(--font-ui);
}
.path-empty-hint {
  font-size: var(--fs-xs);
  color: var(--text-muted);
  font-family: var(--font-ui);
}
.path-empty-examples {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-top: 6px;
}
.path-empty-examples code {
  display: block;
  background: var(--bg-elevated);
  padding: 5px 12px;
  border-radius: var(--radius-sm);
  font-family: var(--font-mono);
  font-size: var(--fs-xs);
  border: 1px solid var(--border-primary);
  color: var(--text-secondary);
  letter-spacing: 0.01em;
}
</style>
