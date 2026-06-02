<template>
  <div class="history-panel">
    <div class="history-toolbar">
      <n-input
        v-model:value="searchKeyword"
        placeholder="搜索 URL..."
        size="small"
        clearable
        class="search-input"
      />
      <n-select
        v-model:value="methodFilter"
        :options="methodFilterOptions"
        size="small"
        class="method-filter"
        :consistent-menu-width="false"
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
        <div class="item-left">
          <span class="item-method" :class="item.method?.toLowerCase()">{{ item.method }}</span>
          <span class="item-url" :title="item.url">{{ truncateUrl(item.url, 60) }}</span>
        </div>
        <div class="item-right">
          <span class="item-status" :class="statusClass(item.response_status)">{{ item.response_status }}</span>
          <span class="item-time">{{ formatTime(item.created_at) }}</span>
          <span class="item-duration">{{ item.duration_ms }}ms</span>
        </div>
      </div>
    </div>
    <div v-else class="history-empty">
      <n-empty description="暂无历史记录" size="small" />
    </div>

    <div class="history-footer">
      <n-button size="tiny" text type="error" @click="onClearAll">清除全部</n-button>
    </div>

    <n-modal v-model:show="detailShow" title="请求详情" preset="card" style="width: 600px;">
      <div v-if="detailItem" class="detail-content">
        <div class="detail-row"><span class="detail-label">方法</span><span class="method-badge" :class="detailItem.method?.toLowerCase()">{{ detailItem.method }}</span></div>
        <div class="detail-row"><span class="detail-label">URL</span><span class="detail-value">{{ detailItem.url }}</span></div>
        <div class="detail-row"><span class="detail-label">状态码</span><span class="detail-value" :class="statusClass(detailItem.response_status)">{{ detailItem.response_status }}</span></div>
        <div class="detail-row"><span class="detail-label">耗时</span><span class="detail-value">{{ detailItem.duration_ms }}ms</span></div>
        <div class="detail-row"><span class="detail-label">时间</span><span class="detail-value">{{ detailItem.created_at }}</span></div>
        <div v-if="detailItem.request_headers" class="detail-section">
          <h4>请求头</h4>
          <pre class="detail-pre">{{ formatJson(detailItem.request_headers) }}</pre>
        </div>
        <div v-if="detailItem.request_body" class="detail-section">
          <h4>请求体</h4>
          <pre class="detail-pre">{{ formatJson(detailItem.request_body) }}</pre>
        </div>
        <div v-if="detailItem.response_headers" class="detail-section">
          <h4>响应头</h4>
          <pre class="detail-pre">{{ formatJson(detailItem.response_headers) }}</pre>
        </div>
        <div v-if="detailItem.response_body" class="detail-section">
          <h4>响应体</h4>
          <pre class="detail-pre">{{ formatJson(detailItem.response_body) }}</pre>
        </div>
      </div>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { NInput, NSelect, NButton, NEmpty, NModal } from 'naive-ui'
import { ListHistory, DeleteHistory, ClearHistory, GetRequest } from '../../../wailsjs/go/main/App'
import { models } from '../../../wailsjs/go/models'
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
}>()

const projectStore = useProjectStore()

const items = ref<HistoryItem[]>([])
const searchKeyword = ref('')
const methodFilter = ref('全部')
const selectedId = ref<number | null>(null)

const detailShow = ref(false)
const detailItem = ref<HistoryItem | null>(null)

const methodFilterOptions = [
  { label: '全部', value: '全部' },
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'DELETE', value: 'DELETE' },
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
  } catch {
    items.value = []
  }
}

function onSelect(item: HistoryItem) {
  selectedId.value = item.id
  detailItem.value = item
  detailShow.value = true
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
  } catch { /* ignore */ }
}

function truncateUrl(url: string, max: number): string {
  if (!url) return ''
  return url.length > max ? url.substring(0, max) + '...' : url
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
    const pad = (n: number) => String(n).padStart(2, '0')
    return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
  } catch {
    return raw
  }
}

function formatJson(raw: string): string {
  try {
    return JSON.stringify(JSON.parse(raw), null, 2)
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
  height: 100%;
}
.history-toolbar {
  display: flex;
  gap: 6px;
  padding: 8px;
  border-bottom: 1px solid #e8e8e8;
}
.search-input {
  flex: 1;
}
.method-filter {
  width: 80px;
}
.history-list {
  flex: 1;
  overflow-y: auto;
}
.history-item {
  display: flex;
  flex-direction: column;
  padding: 8px 10px;
  border-bottom: 1px solid #f0f0f0;
  cursor: pointer;
  transition: background 0.1s;
}
.history-item:hover {
  background: #f5f5f5;
}
.history-item.selected {
  background: #e8f0fe;
}
.item-left {
  display: flex;
  align-items: center;
  gap: 6px;
}
.item-right {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 4px;
  font-size: 11px;
  color: #999;
}
.item-method {
  font-size: 9px;
  font-weight: 700;
  padding: 0 4px;
  border-radius: 2px;
  flex-shrink: 0;
}
.item-method.get { background: #d4edda; color: #155724; }
.item-method.post { background: #fff3cd; color: #856404; }
.item-method.put { background: #d0e8ff; color: #004085; }
.item-method.delete { background: #f8d7da; color: #721c24; }
.item-method.patch { background: #f3e5f5; color: #6a1b9a; }
.item-url {
  font-size: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: #333;
}
.item-status {
  font-size: 11px;
  font-weight: 600;
  padding: 0 4px;
  border-radius: 2px;
}
.status-2xx { color: #18a058; background: #d4edda; }
.status-3xx { color: #0288d1; background: #d0e8ff; }
.status-4xx { color: #f0a020; background: #fff3cd; }
.status-5xx { color: #d03050; background: #f8d7da; }
.item-time {
  white-space: nowrap;
}
.item-duration {
  white-space: nowrap;
}
.history-empty {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}
.history-footer {
  padding: 6px 10px;
  border-top: 1px solid #e8e8e8;
  display: flex;
  justify-content: flex-end;
}
.detail-content {
  font-size: 13px;
}
.detail-row {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px 0;
  border-bottom: 1px solid #f5f5f5;
}
.detail-label {
  color: #999;
  width: 80px;
  flex-shrink: 0;
}
.detail-value {
  word-break: break-all;
}
.method-badge {
  font-size: 10px;
  font-weight: 700;
  padding: 1px 6px;
  border-radius: 2px;
}
.method-badge.get { background: #d4edda; color: #155724; }
.method-badge.post { background: #fff3cd; color: #856404; }
.method-badge.put { background: #d0e8ff; color: #004085; }
.method-badge.delete { background: #f8d7da; color: #721c24; }
.method-badge.patch { background: #f3e5f5; color: #6a1b9a; }
.detail-section {
  margin-top: 12px;
}
.detail-section h4 {
  margin: 0 0 6px 0;
  font-size: 12px;
  color: #666;
  font-weight: 600;
}
.detail-pre {
  background: #f8f9fa;
  border: 1px solid #e8e8e8;
  border-radius: 4px;
  padding: 10px;
  font-size: 11px;
  font-family: 'SF Mono', Consolas, monospace;
  white-space: pre-wrap;
  word-break: break-all;
  max-height: 200px;
  overflow-y: auto;
  margin: 0;
}
</style>
