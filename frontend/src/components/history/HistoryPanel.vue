<template>
  <div class="history-panel">
    <div class="history-toolbar">
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
      <n-select
        v-model:value="methodFilter"
        :options="methodOptions"
        size="small"
        :consistent-menu-width="false"
        class="method-filter"
      />
    </div>

    <div class="history-list" v-if="filteredItems.length > 0">
      <div
        v-for="item in filteredItems"
        :key="item.id"
        class="history-item"
        :class="{ selected: selectedId === item.id }"
        @click="onSelect(item)"
        @dblclick="onReplay(item)"
      >
        <div class="item-row">
          <span class="item-method" :class="item.method?.toLowerCase()">{{ item.method }}</span>
          <span class="item-status" :class="statusClass(item.response_status)">{{ item.response_status }}</span>
          <span class="item-time">{{ formatTime(item.created_at) }}</span>
        </div>
        <div class="hist-url">{{ item.url }}</div>
      </div>
    </div>
    <div v-else class="history-empty">
      <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.2" opacity="0.3">
        <circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/>
      </svg>
      <span class="empty-text">暂无历史记录</span>
    </div>

    <div class="history-footer">
      <n-button size="tiny" quaternary @click="onClearAll" class="clear-btn">清空全部</n-button>
      <n-select
        v-model:value="retentionDays"
        :options="retentionOptions"
        size="tiny"
        :consistent-menu-width="false"
        class="retention-select"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { NInput, NSelect, NButton } from 'naive-ui'
import { ListHistory, ClearHistory } from '../../../wailsjs/go/main/App'
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

const items = ref<HistoryItem[]>([])
const searchKeyword = ref('')
const methodFilter = ref('全部')
const selectedId = ref<number | null>(null)
const retentionDays = ref(30)

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

async function loadHistory() {
  const pid = projectStore.currentId
  if (!pid) return
  try {
    const res = await ListHistory(pid, 200, 0)
    items.value = (res || []) as HistoryItem[]
  } catch { items.value = [] }
}

function onSelect(item: HistoryItem) {
  selectedId.value = item.id
  emit('select-detail', item)
}

async function onReplay(item: HistoryItem) { emit('open-tab', item) }

async function onClearAll() {
  const pid = projectStore.currentId
  if (!pid) return
  try {
    await ClearHistory(pid)
    items.value = []
    selectedId.value = null
    emit('select-detail', null as any)
  } catch { /* ignore */ }
}

function statusClass(code: number): string {
  if (code < 300) return 'status-2xx'
  if (code < 400) return 'status-3xx'
  if (code < 500) return 'status-4xx'
  return 'status-5xx'
}

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

watch(() => projectStore.currentId, () => { if (projectStore.currentId) loadHistory() })
onMounted(() => { if (projectStore.currentId) loadHistory() })
</script>

<style scoped>
.history-panel {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
}

.history-toolbar {
  display: flex;
  gap: 4px;
  padding: 8px 10px;
  border-bottom: 1px solid var(--border-subtle);
  background: var(--bg-base);
}
.search-input {
  flex: 1;
}
.search-input :deep(.n-input) {
  height: 28px;
  border-radius: var(--radius-sm);
}
.search-input :deep(.n-input__prefix svg) {
  color: var(--text-muted);
}
.method-filter {
  width: 76px;
  flex-shrink: 0;
}
.method-filter :deep(.n-base-selection) {
  height: 28px;
  border-radius: var(--radius-sm);
}

/* ── History List ── */
.history-list {
  flex: 1;
  overflow-y: auto;
  padding: 2px 0;
}
.history-item {
  padding: 8px 10px;
  margin: 1px 4px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all var(--transition-fast);
  animation: slideUp 0.2s var(--ease-out) both;
  border-left: 2px solid transparent;
}
.history-item:hover {
  background: var(--bg-hover);
}
.history-item.selected {
  background: var(--accent-soft);
  border-left-color: var(--accent);
}
.item-row {
  display: flex;
  gap: 6px;
  align-items: center;
}
.hist-url {
  font-size: var(--fs-2xs);
  color: var(--text-muted);
  margin-top: 3px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-family: var(--font-mono);
  padding-left: 2px;
}

.item-method {
  font-size: var(--fs-2xs);
  font-weight: 700;
  padding: 2px 5px;
  border-radius: var(--radius-xs);
  letter-spacing: 0.04em;
  text-transform: uppercase;
  line-height: 1.3;
}
.item-method.get    { background: var(--blue-soft); color: var(--method-get); }
.item-method.post   { background: rgba(34,197,94,0.1); color: var(--method-post); }
.item-method.put    { background: var(--amber-soft); color: var(--method-put); }
.item-method.delete { background: var(--red-soft); color: var(--method-delete); }

.item-status {
  font-size: var(--fs-xs);
  font-weight: 700;
  font-family: var(--font-mono);
}
.status-2xx { color: var(--method-post); }
.status-3xx { color: var(--blue); }
.status-4xx { color: var(--amber); }
.status-5xx { color: var(--red); }

.item-time {
  font-size: var(--fs-2xs);
  color: var(--text-muted);
  margin-left: auto;
  font-family: var(--font-mono);
}

/* ── Empty State ── */
.history-empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: var(--text-muted);
}
.empty-text {
  font-size: var(--fs-sm);
  font-family: var(--font-mono);
}

/* ── Footer ── */
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
.clear-btn:hover {
  color: var(--red) !important;
}
.retention-select {
  width: 80px;
  margin-left: auto;
}
.retention-select :deep(.n-base-selection) {
  height: 24px;
  font-size: var(--fs-2xs);
}
</style>
