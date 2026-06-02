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
        <div style="display:flex;gap:6px;align-items:center">
          <span class="item-method" :class="item.method?.toLowerCase()">{{ item.method }}</span>
          <span class="item-status" :class="statusClass(item.response_status)">{{ item.response_status }}</span>
          <span style="color:#999;font-size:10px;margin-left:auto">{{ formatTime(item.created_at) }}</span>
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
  } catch {
    items.value = []
  }
}

function onSelect(item: HistoryItem) {
  selectedId.value = item.id
  emit('select-detail', item)
}

async function onReplay(item: HistoryItem) {
  emit('open-tab', item)
}

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
    if (diffMin < 60) return `${diffMin} min ago`
    const diffHour = Math.floor(diffMin / 60)
    if (diffHour < 24) return `${diffHour} hour ago`
    return `${Math.floor(diffHour / 24)} day ago`
  } catch {
    return raw
  }
}

watch(() => projectStore.currentId, () => {
  if (projectStore.currentId) loadHistory()
})

onMounted(() => {
  if (projectStore.currentId) loadHistory()
})
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
  border-bottom: 1px solid #e8e8e8;
}
.search-input {
  flex: 1;
  padding: 5px 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 12px;
  outline: none;
}
.search-input:focus { border-color: #18a058; }
.method-filter {
  width: 70px;
  font-size: 12px;
  padding: 4px;
  border: 1px solid #ddd;
  border-radius: 4px;
  outline: none;
  cursor: pointer;
}
.history-list {
  flex: 1;
  overflow-y: auto;
}
.history-item {
  padding: 8px 10px;
  border-bottom: 1px solid #f0f0f0;
  cursor: pointer;
}
.history-item:hover {
  background: #f8f8f8;
}
.history-item.selected {
  background: #e8f0fe;
}
.hist-url {
  font-size: 10px;
  color: #999;
  margin-top: 2px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.item-method {
  font-size: 10px;
  font-weight: 700;
  padding: 1px 5px;
  border-radius: 3px;
}
.item-method.get { background: #d4edda; color: #155724; }
.item-method.post { background: #fff3cd; color: #856404; }
.item-method.put { background: #d0e8ff; color: #004085; }
.item-method.delete { background: #f8d7da; color: #721c24; }
.item-method.patch { background: #f3e5f5; color: #6a1b9a; }
.item-status {
  font-size: 12px;
  font-weight: 600;
  padding: 0 4px;
  border-radius: 2px;
}
.status-2xx { color: #18a058; }
.status-3xx { color: #0288d1; }
.status-4xx { color: #f0a020; }
.status-5xx { color: #d03050; }
.history-empty {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}
.empty-text {
  color: #aaa;
  font-size: 12px;
}
.history-footer {
  padding: 8px 10px;
  border-top: 1px solid #eee;
  display: flex;
  gap: 8px;
  align-items: center;
}
.footer-btn {
  flex: 1;
  font-size: 10px;
  padding: 4px 8px;
  background: #fff;
  border: 1px solid #ddd;
  border-radius: 4px;
  cursor: pointer;
  color: #555;
}
.footer-btn:hover { background: #f8f8f8; border-color: #ccc; }
.retention-select {
  font-size: 10px;
  padding: 4px;
  border: 1px solid #ddd;
  border-radius: 4px;
  outline: none;
  cursor: pointer;
}
</style>
