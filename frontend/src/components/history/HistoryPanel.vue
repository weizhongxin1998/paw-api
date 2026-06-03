<template>
  <div class="history-panel">
    <div class="history-toolbar">
      <input
        v-model="searchKeyword"
        placeholder="搜索 URL..."
        class="search-input"
      />
      <select v-model="methodFilter" class="method-filter">
        <option>全部</option><option>GET</option><option>POST</option>
        <option>PUT</option><option>DELETE</option>
      </select>
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
      <span class="empty-text">暂无历史记录</span>
    </div>

    <div class="history-footer">
      <button class="footer-btn" @click="onClearAll">清空全部</button>
      <select v-model="retentionDays" class="retention-select">
        <option :value="30">保留 30 天</option>
        <option :value="60">保留 60 天</option>
        <option :value="0">永久保留</option>
      </select>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
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
  gap: 3px;
  padding: 6px 8px;
  border-bottom: 1px solid var(--border-primary);
  background: var(--bg-base);
}
.search-input {
  flex: 1;
  padding: 4px 7px;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-sm);
  font-size: var(--fs-xs);
  outline: none;
  background: var(--bg-surface);
  color: var(--text-primary);
  font-family: var(--font-mono);
  transition: border-color var(--transition);
}
.search-input:focus { border-color: var(--accent); }
.search-input::placeholder { color: var(--text-muted); }
.method-filter {
  width: 60px;
  font-size: var(--fs-xs);
  padding: 3px;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-sm);
  outline: none;
  cursor: pointer;
  background: var(--bg-surface);
  color: var(--text-secondary);
  font-family: var(--font-mono);
}
.history-list { flex: 1; overflow-y: auto; }
.history-item {
  padding: 7px 8px;
  border-bottom: 1px solid var(--border-primary);
  cursor: pointer;
  transition: background var(--transition), border-color var(--transition);
  animation: fadeIn 0.2s ease both;
}
.history-item:hover { background: var(--bg-hover); }
.history-item.selected { background: var(--accent-soft); border-left: 2px solid var(--accent); padding-left: 6px; }
.item-row { display: flex; gap: 5px; align-items: center; }
.hist-url {
  font-size: var(--fs-2xs);
  color: var(--text-muted);
  margin-top: 2px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-family: var(--font-mono);
}
.item-method {
  font-size: var(--fs-2xs);
  font-weight: 700;
  padding: 1px 4px;
  border-radius: 2px;
  letter-spacing: 0.3px;
}
.item-method.get { background: var(--accent-soft); color: var(--accent); }
.item-method.post { background: var(--amber-soft); color: var(--amber); }
.item-method.put { background: var(--blue-soft); color: var(--blue); }
.item-method.delete { background: var(--red-soft); color: var(--red); }
.item-status {
  font-size: var(--fs-sm);
  font-weight: 700;
  font-family: var(--font-mono);
}
.status-2xx { color: var(--accent); }
.status-3xx { color: var(--blue); }
.status-4xx { color: var(--amber); }
.status-5xx { color: var(--red); }
.item-time { font-size: var(--fs-2xs); color: var(--text-muted); margin-left: auto; font-family: var(--font-mono); }
.history-empty {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}
.empty-text { color: var(--text-muted); font-size: var(--fs-sm); font-family: var(--font-mono); }
.history-footer {
  padding: 6px 8px;
  border-top: 1px solid var(--border-primary);
  display: flex;
  gap: 6px;
  align-items: center;
}
.footer-btn {
  flex: 1;
  font-size: var(--fs-2xs);
  padding: 3px 6px;
  background: var(--bg-surface);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-sm);
  cursor: pointer;
  color: var(--text-muted);
  font-family: var(--font-mono);
  transition: all var(--transition);
}
.footer-btn:hover { border-color: var(--red); color: var(--red); }
.retention-select {
  font-size: var(--fs-2xs);
  padding: 3px;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-sm);
  outline: none;
  cursor: pointer;
  background: var(--bg-surface);
  color: var(--text-muted);
  font-family: var(--font-mono);
}
</style>
