<template>
  <div class="response-panel">
    <div class="status-bar">
      <span v-if="response" class="status-badge" :class="statusClass">
        {{ response.status }} {{ statusText }}
      </span>
      <span v-if="response" class="meta">{{ response.time }} ms</span>
      <span v-if="response" class="meta">{{ formatSize(response.size) }}</span>
      <span v-if="!response" class="placeholder">点击 Send 发送请求</span>
    </div>

    <div class="sub-tabs">
      <button :class="{ active: activeTab === 'Body' }" @click="activeTab = 'Body'">Body</button>
      <button :class="{ active: activeTab === 'Headers' }" @click="activeTab = 'Headers'">Headers</button>
      <button :class="{ active: activeTab === 'Cookies' }" @click="activeTab = 'Cookies'">Cookies</button>
      <button :class="{ active: activeTab === 'Log' }" @click="activeTab = 'Log'">Log</button>
    </div>

    <div v-if="!response" class="body-empty">
      <span class="hint">响应体将在此显示</span>
    </div>

    <div v-else class="body-content">
      <div v-if="activeTab === 'Body'">
        <div class="body-modes">
          <button :class="{ active: bodyMode === 'pretty' }" @click="bodyMode = 'pretty'">Pretty</button>
          <button :class="{ active: bodyMode === 'raw' }" @click="bodyMode = 'raw'">Raw</button>
          <button :class="{ active: bodyMode === 'preview' }" @click="bodyMode = 'preview'">Preview</button>
        </div>
        <pre v-if="bodyMode === 'pretty'" class="code-block">{{ formatBody(response.body) }}</pre>
        <pre v-else-if="bodyMode === 'raw'" class="code-block raw">{{ response.body }}</pre>
        <div v-else class="preview-hint">Preview 模式需要 HTML 类型响应</div>
      </div>

      <div v-else-if="activeTab === 'Headers'">
        <table class="kv-table">
          <thead><tr><th>Key</th><th>Value</th></tr></thead>
          <tbody>
            <tr v-for="(v, k) in response.headers" :key="k">
              <td>{{ k }}</td><td>{{ v }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-else-if="activeTab === 'Cookies'">
        <table v-if="response.cookies?.length" class="kv-table">
          <thead><tr><th>Name</th><th>Value</th><th>Domain</th><th>Path</th></tr></thead>
          <tbody>
            <tr v-for="c in response.cookies" :key="c.name">
              <td>{{ c.name }}</td><td>{{ c.value }}</td><td>{{ c.domain }}</td><td>{{ c.path }}</td>
            </tr>
          </tbody>
        </table>
        <span v-else class="hint">无 Cookie</span>
      </div>

      <div v-else-if="activeTab === 'Log'">
        <button class="copy-curl-btn" @click="copyCurl">复制 cURL</button>
        <pre class="code-block">{{ response.rawRequest || '(暂无)' }}</pre>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { HttpResponse } from '../../types/response'

const props = defineProps<{
  response: HttpResponse | null
}>()

const activeTab = ref('Body')
const bodyMode = ref('pretty')

const statusClass = computed(() => {
  if (!props.response) return ''
  const s = props.response.status
  if (s < 300) return 'green'
  if (s < 400) return 'blue'
  if (s < 500) return 'orange'
  return 'red'
})

const statusText = computed(() => {
  if (!props.response) return ''
  const s = props.response.status
  const map: Record<number, string> = { 200: 'OK', 201: 'Created', 204: 'No Content', 400: 'Bad Request', 401: 'Unauthorized', 404: 'Not Found', 500: 'Server Error' }
  return map[s] || ''
})

function formatSize(bytes: number): string {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

function formatBody(raw: string): string {
  try { return JSON.stringify(JSON.parse(raw), null, 2) } catch { return raw }
}

function copyCurl() {
  if (props.response?.curlCommand) {
    navigator.clipboard.writeText(props.response.curlCommand)
  }
}
</script>

<style scoped>
.response-panel {
  display: flex;
  flex-direction: column;
  min-height: 100px;
  background: #fff;
}
.status-bar {
  display: flex;
  padding: 8px 14px;
  background: var(--gray-50);
  border-bottom: 1px solid var(--gray-200);
  align-items: center;
  gap: 14px;
}
.status-badge {
  font-weight: 700;
  font-size: 12px;
  padding: 3px 10px;
  border-radius: var(--radius-sm);
}
.status-badge.green { background: var(--green-soft); color: var(--green); }
.status-badge.blue { background: var(--blue-soft); color: var(--blue); }
.status-badge.orange { background: var(--amber-soft); color: var(--amber); }
.status-badge.red { background: var(--red-soft); color: var(--red); }
.meta { font-size: 11px; color: var(--gray-500); }
.placeholder { font-size: 11px; color: var(--gray-400); }
.sub-tabs {
  display: flex;
  border-bottom: 1px solid var(--gray-200);
  background: var(--gray-50);
  padding: 0 12px;
}
.sub-tabs button {
  padding: 7px 14px;
  font-size: 12px;
  cursor: pointer;
  color: var(--gray-500);
  border: none;
  background: transparent;
  border-bottom: 2px solid transparent;
  outline: none;
  transition: all var(--transition);
}
.sub-tabs button.active {
  color: var(--green);
  border-bottom-color: var(--green);
  font-weight: 600;
}
.sub-tabs button:hover:not(.active) { color: var(--gray-700); }
.body-content {
  flex: 1;
  overflow-y: auto;
  padding: 12px 14px;
  background: #fff;
}
.body-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 30px;
}
.body-modes {
  display: flex;
  gap: 4px;
  margin-bottom: 10px;
}
.body-modes button {
  padding: 3px 12px;
  border: 1px solid var(--gray-200);
  border-radius: var(--radius-sm);
  background: #fff;
  font-size: 11px;
  cursor: pointer;
  color: var(--gray-500);
  outline: none;
  transition: all var(--transition);
}
.body-modes button.active {
  border-color: var(--green);
  color: var(--green);
  font-weight: 600;
  background: var(--green-soft);
}
.code-block {
  font-family: 'SF Mono', 'Consolas', monospace;
  font-size: 12px;
  line-height: 1.7;
  margin: 0;
  white-space: pre-wrap;
  word-break: break-all;
  color: var(--gray-700);
  background: var(--gray-50);
  border: 1px solid var(--gray-200);
  border-radius: var(--radius);
  padding: 12px;
}
.code-block.raw {
  color: var(--gray-600);
  font-size: 11px;
}
.kv-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 12px;
}
.kv-table th, .kv-table td {
  text-align: left;
  padding: 5px 10px;
  border-bottom: 1px solid var(--gray-100);
}
.kv-table th { color: var(--gray-500); font-weight: 500; text-transform: uppercase; font-size: 10px; letter-spacing: 0.3px; }
.kv-table td { color: var(--gray-700); }
.preview-hint { color: var(--gray-400); font-size: 12px; padding: 20px; }
.hint { color: var(--gray-400); font-size: 12px; }
.copy-curl-btn {
  padding: 4px 12px;
  border: 1px solid var(--gray-200);
  border-radius: var(--radius-sm);
  background: #fff;
  font-size: 11px;
  cursor: pointer;
  color: var(--gray-500);
  margin-bottom: 8px;
  transition: all var(--transition);
}
.copy-curl-btn:hover { border-color: var(--gray-300); color: var(--gray-700); }
</style>
