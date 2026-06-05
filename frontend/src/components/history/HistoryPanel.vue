<template>
  <div class="history-panel">
    <!-- Toolbar with filter count indicator -->
    <div class="history-toolbar">
      <div class="search-wrap">
        <n-input
          v-model:value="searchKeyword"
          placeholder="搜索 URL..."
          size="small"
          clearable
          class="search-input"
        >
          <template #prefix>
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
          </template>
        </n-input>
        <span class="filter-count" v-if="items.length > 0">{{ filteredItems.length }}/{{ items.length }}</span>
      </div>
      <n-select
        v-model:value="methodFilter"
        :options="methodOptions"
        size="small"
        :consistent-menu-width="false"
        class="method-filter"
      />
    </div>

    <!-- History list with fade gradient -->
    <div class="history-list-wrap" v-if="filteredItems.length > 0">
      <div class="history-list" ref="listRef">
        <div
          v-for="(item, idx) in filteredItems"
          :key="item.id"
          class="history-item"
          :class="[
            { selected: selectedId === item.id },
            { focused: focusedIndex === idx },
            methodBorderClass(item.method),
          ]"
          @click="onSelect(item)"
          @dblclick="onReplay(item)"
          @contextmenu.prevent="onContextMenu($event, item)"
        >
          <div class="item-row">
            <span class="item-method" :class="item.method?.toLowerCase()">{{ item.method }}</span>
            <span class="item-status-pill" :class="statusClass(item.response_status)">{{ item.response_status }}</span>
            <span class="item-duration" :class="durationClass(item.duration_ms)">{{ formatDuration(item.duration_ms) }}</span>
            <span class="item-time">{{ formatTime(item.created_at) }}</span>
          </div>
          <div class="hist-url" :title="item.url">{{ item.url }}</div>
        </div>
      </div>
      <div class="list-fade"></div>
    </div>

    <!-- Empty state with clock icon -->
    <div v-else class="history-empty">
      <svg class="empty-clock-icon" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="12" cy="12" r="10"/>
        <polyline points="12 6 12 12 16 14"/>
      </svg>
      <span class="empty-text">{{ searchKeyword || methodFilter !== '全部' ? '没有匹配的历史记录' : '发送请求后，历史记录将在此显示' }}</span>
      <span class="empty-hint" v-if="!searchKeyword && methodFilter === '全部'">双击历史记录可重放请求</span>
    </div>

    <!-- Footer -->
    <div class="history-footer">
      <n-button size="tiny" quaternary @click="onClearAll" class="clear-btn">
        <template #icon>
          <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/></svg>
        </template>
        清空全部
      </n-button>
      <n-select
        v-model:value="retentionDays"
        :options="retentionOptions"
        size="tiny"
        :consistent-menu-width="false"
        class="retention-select"
      />
    </div>

    <!-- Right-click context menu -->
    <Teleport to="body">
      <Transition name="ctx-menu">
        <div
          v-if="ctxMenu.visible"
          class="hist-ctx-menu"
          :style="{ left: ctxMenu.x + 'px', top: ctxMenu.y + 'px' }"
        >
          <div class="ctx-item" @click="ctxReplay">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><polygon points="5 3 19 12 5 21 5 3"/></svg>
            <span>重放此请求</span>
          </div>
          <div class="ctx-item" @click="ctxCopyUrl">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"/><path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/></svg>
            <span>复制 URL</span>
          </div>
          <div class="ctx-divider"></div>
          <div class="ctx-item ctx-danger" @click="ctxDelete">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/></svg>
            <span>删除此记录</span>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { NInput, NSelect, NButton, useDialog, useMessage } from 'naive-ui'
import { ListHistory, ClearHistory, DeleteHistory } from '../../../wailsjs/go/main/App'
import { ClipboardSetText, EventsOn, EventsOff } from '../../../wailsjs/runtime/runtime'
import { useProjectStore } from '../../stores/project'

interface HistoryItem {
  id: number
  project_id: number
  request_id?: number
  method: string
  url: string
  request_headers: string
  request_body: string
  response_status: number
  response_headers: string
  response_body: string
  duration_ms: number
  created_at: string
}

const emit = defineEmits<{
  (e: 'open-tab', item: HistoryItem): void
  (e: 'select-detail', item: HistoryItem): void
}>()

const projectStore = useProjectStore()
const dialog = useDialog()
const message = useMessage()

const items = ref<HistoryItem[]>([])
const searchKeyword = ref('')
const methodFilter = ref('全部')
const selectedId = ref<number | null>(null)
const retentionDays = ref(30)
const focusedIndex = ref(-1)
const listRef = ref<HTMLElement | null>(null)

// Context menu state
const ctxMenu = ref({ visible: false, x: 0, y: 0, item: null as HistoryItem | null })

const methodOptions = [
  { label: '全部', value: '全部' },
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'DELETE', value: 'DELETE' },
]

const retentionOptions = [
  { label: '30 天', value: 30 },
  { label: '60 天', value: 60 },
  { label: '永久', value: 0 },
]

const filteredItems = computed(() => {
  let result = items.value
  if (searchKeyword.value) {
    const kw = searchKeyword.value.toLowerCase()
    result = result.filter(i => i.url.toLowerCase().includes(kw))
  }
  if (methodFilter.value !== '全部') {
    result = result.filter(i => i.method === methodFilter.value)
  }
  return result
})

// -- Method border color class --
function methodBorderClass(method: string): string {
  return 'border-' + (method || 'get').toLowerCase()
}

// -- Status pill class --
function statusClass(code: number): string {
  if (!code) return 'status-0'
  if (code < 300) return 'status-2xx'
  if (code < 400) return 'status-3xx'
  if (code < 500) return 'status-4xx'
  return 'status-5xx'
}

// -- Duration formatting and color coding --
function formatDuration(ms: number): string {
  if (ms == null) return ''
  if (ms < 1000) return ms + 'ms'
  return (ms / 1000).toFixed(1) + 's'
}

function durationClass(ms: number): string {
  if (ms == null) return ''
  if (ms < 200) return 'dur-fast'   // green
  if (ms < 1000) return 'dur-mid'   // amber
  return 'dur-slow'                  // red
}

// -- Relative time formatting --
function formatTime(raw: string): string {
  if (!raw) return ''
  try {
    const d = new Date(raw)
    if (isNaN(d.getTime())) return raw
    const now = Date.now()
    const diffMs = now - d.getTime()
    const diffMin = Math.floor(diffMs / 60000)
    if (diffMin < 1) return '刚刚'
    if (diffMin < 60) return `${diffMin}m`
    const diffHour = Math.floor(diffMin / 60)
    if (diffHour < 24) return `${diffHour}h`
    return `${Math.floor(diffHour / 24)}d`
  } catch { return raw }
}

// -- Load history from backend --
async function loadHistory() {
  const pid = projectStore.currentId
  if (!pid) return
  try {
    const res = await ListHistory(pid, 200, 0)
    items.value = (res || []) as HistoryItem[]
  } catch { items.value = [] }
}

// Expose refresh for parent components
function refresh() { loadHistory() }
defineExpose({ refresh })

// -- Selection --
function onSelect(item: HistoryItem) {
  selectedId.value = item.id
  focusedIndex.value = filteredItems.value.findIndex(i => i.id === item.id)
  emit('select-detail', item)
}

async function onReplay(item: HistoryItem) { emit('open-tab', item) }

// -- Context menu handlers --
function onContextMenu(e: MouseEvent, item: HistoryItem) {
  ctxMenu.value = { visible: true, x: e.clientX, y: e.clientY, item }
}

function closeCtxMenu() {
  ctxMenu.value.visible = false
  ctxMenu.value.item = null
}

function ctxReplay() {
  if (ctxMenu.value.item) onReplay(ctxMenu.value.item)
  closeCtxMenu()
}

async function ctxCopyUrl() {
  if (ctxMenu.value.item) {
    try {
      await ClipboardSetText(ctxMenu.value.item.url)
      message.success('URL 已复制到剪贴板')
    } catch {
      message.error('复制失败')
    }
  }
  closeCtxMenu()
}

async function ctxDelete() {
  const item = ctxMenu.value.item
  if (!item) return
  closeCtxMenu()
  try {
    await DeleteHistory(item.id)
    items.value = items.value.filter(i => i.id !== item.id)
    if (selectedId.value === item.id) {
      selectedId.value = null
      emit('select-detail', null as any)
    }
    message.success('已删除')
  } catch { message.error('删除失败') }
}

// -- Clear all history --
async function onClearAll() {
  const pid = projectStore.currentId
  if (!pid) return
  dialog.warning({
    title: '确认清空',
    content: '确定要清空所有历史记录吗？此操作不可撤销。',
    positiveText: '清空',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await ClearHistory(pid)
        items.value = []
        selectedId.value = null
        focusedIndex.value = -1
        emit('select-detail', null as any)
        message.success('已清空历史记录')
      } catch { /* ignore */ }
    },
  })
}

// -- Keyboard navigation --
function handleKeydown(e: KeyboardEvent) {
  const list = filteredItems.value
  if (!list.length) return

  const tag = (e.target as HTMLElement)?.tagName
  const isInputFocused = tag === 'INPUT' || tag === 'TEXTAREA' || (e.target as HTMLElement)?.isContentEditable

  if (e.key === 'ArrowDown' && !isInputFocused) {
    e.preventDefault()
    focusedIndex.value = Math.min(focusedIndex.value + 1, list.length - 1)
    onSelect(list[focusedIndex.value])
    scrollToFocused()
  } else if (e.key === 'ArrowUp' && !isInputFocused) {
    e.preventDefault()
    focusedIndex.value = Math.max(focusedIndex.value - 1, 0)
    onSelect(list[focusedIndex.value])
    scrollToFocused()
  } else if (e.key === 'Enter' && focusedIndex.value >= 0 && !isInputFocused) {
    e.preventDefault()
    onReplay(list[focusedIndex.value])
  } else if (e.key === 'Delete' && focusedIndex.value >= 0 && !isInputFocused) {
    e.preventDefault()
    const item = list[focusedIndex.value]
    if (item) {
      ctxMenu.value.item = item
      ctxDelete()
    }
  }
}

function scrollToFocused() {
  nextTick(() => {
    const el = listRef.value?.querySelector('.history-item.focused')
    el?.scrollIntoView({ block: 'nearest', behavior: 'smooth' })
  })
}

// Close context menu on outside click
function onDocumentClick() { closeCtxMenu() }

// -- Lifecycle --
onMounted(() => {
  if (projectStore.currentId) loadHistory()
  document.addEventListener('click', onDocumentClick)
  document.addEventListener('keydown', handleKeydown)

  // Auto-refresh when backend emits history-updated event
  try {
    EventsOn('history-updated', () => loadHistory())
  } catch { /* event may not be registered */ }
})

onUnmounted(() => {
  document.removeEventListener('click', onDocumentClick)
  document.removeEventListener('keydown', handleKeydown)
  try { EventsOff('history-updated') } catch { /* ignore */ }
})

watch(() => projectStore.currentId, () => { if (projectStore.currentId) loadHistory() })
watch([searchKeyword, methodFilter], () => { focusedIndex.value = -1 })
</script>

<style scoped>
.history-panel {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
}

/* -- Toolbar -- */
.history-toolbar {
  display: flex;
  gap: 4px;
  padding: 8px 10px;
  border-bottom: 1px solid var(--border-subtle);
  background: var(--bg-base);
}
.search-wrap {
  flex: 1;
  position: relative;
  display: flex;
  align-items: center;
}
.search-input { flex: 1; }
.search-input :deep(.n-input) {
  height: 28px;
  border-radius: var(--radius-sm);
}
.search-input :deep(.n-input__prefix svg) {
  color: var(--text-muted);
}
.filter-count {
  position: absolute;
  right: 6px;
  font-size: var(--fs-2xs);
  font-family: var(--font-mono);
  color: var(--text-muted);
  background: var(--bg-elevated, var(--bg-base));
  padding: 1px 5px;
  border-radius: 6px;
  pointer-events: none;
  opacity: 0.85;
  line-height: 1.4;
}
.method-filter { width: 76px; flex-shrink: 0; }
.method-filter :deep(.n-base-selection) {
  height: 28px;
  border-radius: var(--radius-sm);
}

/* -- History list with fade gradient -- */
.history-list-wrap {
  flex: 1;
  position: relative;
  overflow: hidden;
  min-height: 0;
}
.history-list {
  height: 100%;
  overflow-y: auto;
  padding: 4px 0;
}
.list-fade {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 32px;
  background: linear-gradient(to bottom, transparent, var(--bg-base));
  pointer-events: none;
}

/* -- History item cards with method-colored left border -- */
.history-item {
  padding: 8px 10px;
  margin: 2px 6px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.15s ease;
  animation: slideUp 0.2s ease-out both;
  border-left: 3px solid transparent;
  position: relative;
}
.history-item:hover { background: var(--bg-hover); }
.history-item.selected { background: var(--accent-soft); }
.history-item.focused {
  outline: 1.5px solid var(--accent);
  outline-offset: -1.5px;
}

/* Method-colored left borders */
.history-item.border-get    { border-left-color: var(--method-get, #3b82f6); }
.history-item.border-post   { border-left-color: var(--method-post, #22c55e); }
.history-item.border-put    { border-left-color: var(--method-put, #f59e0b); }
.history-item.border-delete { border-left-color: var(--method-delete, #ef4444); }
.history-item.border-patch  { border-left-color: #a855f7; }
.history-item.border-head   { border-left-color: #06b6d4; }

.item-row {
  display: flex;
  gap: 6px;
  align-items: center;
}
.hist-url {
  font-size: var(--fs-xs);
  color: var(--text-secondary);
  margin-top: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-family: var(--font-mono);
  padding-left: 2px;
}

/* Method badge */
.item-method {
  font-size: var(--fs-2xs);
  font-weight: 700;
  padding: 2px 6px;
  border-radius: var(--radius-xs);
  letter-spacing: 0.04em;
  text-transform: uppercase;
  line-height: 1.3;
  flex-shrink: 0;
}
.item-method.get    { background: var(--blue-soft, rgba(59,130,246,0.12)); color: var(--method-get); }
.item-method.post   { background: rgba(34,197,94,0.1); color: var(--method-post); }
.item-method.put    { background: var(--amber-soft, rgba(245,158,11,0.12)); color: var(--method-put); }
.item-method.delete { background: var(--red-soft, rgba(239,68,68,0.1)); color: var(--method-delete); }

/* Status code pill badges */
.item-status-pill {
  font-size: var(--fs-2xs);
  font-weight: 700;
  font-family: var(--font-mono);
  padding: 1px 6px;
  border-radius: 10px;
  line-height: 1.5;
  flex-shrink: 0;
}
.status-2xx { background: rgba(34,197,94,0.12); color: var(--method-post); }
.status-3xx { background: rgba(59,130,246,0.12); color: var(--blue, #3b82f6); }
.status-4xx { background: rgba(245,158,11,0.12); color: var(--amber, #f59e0b); }
.status-5xx { background: rgba(239,68,68,0.12); color: var(--red, #ef4444); }
.status-0   { background: rgba(113,113,122,0.12); color: var(--text-muted); }

/* Duration color coding: green <200ms, amber <1s, red >1s */
.item-duration {
  font-size: var(--fs-2xs);
  font-family: var(--font-mono);
  flex-shrink: 0;
}
.dur-fast { color: var(--method-post); }
.dur-mid  { color: var(--amber, #f59e0b); }
.dur-slow { color: var(--red, #ef4444); }

.item-time {
  font-size: var(--fs-2xs);
  color: var(--text-muted);
  margin-left: auto;
  font-family: var(--font-mono);
  flex-shrink: 0;
}

/* -- Empty State -- */
.history-empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;
  color: var(--text-muted);
  padding: 20px;
}
.empty-clock-icon {
  opacity: 0.18;
  stroke: var(--text-muted);
}
.empty-text {
  font-size: var(--fs-sm);
  font-family: var(--font-family, var(--font-mono));
  text-align: center;
  line-height: 1.5;
}
.empty-hint {
  font-size: var(--fs-xs);
  opacity: 0.6;
}

/* -- Footer -- */
.history-footer {
  padding: 8px 10px;
  border-top: 1px solid var(--border-subtle);
  display: flex;
  gap: 8px;
  align-items: center;
}
.clear-btn {
  color: var(--text-muted) !important;
  font-size: var(--fs-2xs) !important;
}
.clear-btn:hover { color: var(--red) !important; }
.retention-select { width: 80px; margin-left: auto; }
.retention-select :deep(.n-base-selection) {
  height: 24px;
  font-size: var(--fs-2xs);
}

/* -- Context Menu (teleported to body) -- */
.hist-ctx-menu {
  position: fixed;
  z-index: 99999;
  min-width: 160px;
  padding: 4px 0;
  background: var(--bg-elevated, #1e1e22);
  border: 1px solid var(--border-primary, #27272a);
  border-radius: var(--radius, 8px);
  box-shadow: 0 8px 24px rgba(0,0,0,0.35);
  font-size: var(--fs-sm);
  font-family: var(--font-family);
}
.ctx-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 7px 14px;
  cursor: pointer;
  color: var(--text-secondary, #a1a1aa);
  transition: background 0.12s, color 0.12s;
}
.ctx-item:hover {
  background: var(--bg-hover, rgba(255,255,255,0.04));
  color: var(--text-primary, #e4e4e7);
}
.ctx-item svg { opacity: 0.55; flex-shrink: 0; }
.ctx-item.ctx-danger { color: var(--red, #ef4444); }
.ctx-item.ctx-danger:hover { background: rgba(239,68,68,0.08); }
.ctx-divider {
  height: 1px;
  background: var(--border-primary, #27272a);
  margin: 4px 8px;
}

/* Context menu transition */
.ctx-menu-enter-active { transition: opacity 0.12s, transform 0.12s; }
.ctx-menu-leave-active { transition: opacity 0.1s, transform 0.1s; }
.ctx-menu-enter-from { opacity: 0; transform: scale(0.95) translateY(-4px); }
.ctx-menu-leave-to   { opacity: 0; transform: scale(0.95); }

@keyframes slideUp {
  from { opacity: 0; transform: translateY(6px); }
  to   { opacity: 1; transform: translateY(0); }
}
</style>
