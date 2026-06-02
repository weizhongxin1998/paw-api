<template>
  <div class="response-panel">
    <div class="status-bar">
      <span v-if="response" class="status-badge" :class="statusClass">
        {{ response.status }} {{ statusText }}
      </span>
      <span v-if="response" class="meta">{{ response.time }}ms</span>
      <span v-if="response" class="meta">{{ formatSize(response.size) }}</span>
      <span v-if="!response" class="placeholder">点击 Send 发送请求</span>
      <span style="flex:1"></span>
      <span v-if="response" class="meta time-badge">{{ new Date().toLocaleTimeString() }}</span>
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
          <span style="flex:1"></span>
          <button class="copy-btn" @click="copyResponse" v-if="response.body">复制</button>
        </div>
        <pre v-if="bodyMode === 'pretty'" class="code-block">{{ formatBody(response.body) }}</pre>
        <pre v-else-if="bodyMode === 'raw'" class="code-block raw">{{ response.body }}</pre>
        <div v-else class="preview-hint">
          <iframe v-if="isHtml(response.body)" :srcdoc="response.body" class="preview-iframe" sandbox="allow-scripts"></iframe>
          <span v-else>非 HTML 响应，无法预览</span>
        </div>
      </div>

      <div v-else-if="activeTab === 'Headers'">
        <table class="kv-table">
          <thead><tr><th>Key</th><th>Value</th></tr></thead>
          <tbody>
            <tr v-for="(v, k) in response.headers" :key="k">
              <td class="key-cell">{{ k }}</td><td class="val-cell">{{ v }}</td>
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
  const map: Record<number, string> = { 200:'OK', 201:'Created', 204:'No Content', 301:'Moved', 400:'Bad Request', 401:'Unauthorized', 403:'Forbidden', 404:'Not Found', 500:'Server Error', 502:'Bad Gateway', 503:'Unavailable' }
  return map[s] || ''
})

function formatSize(bytes: number): string {
  if (bytes < 1024) return bytes + 'B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + 'KB'
  return (bytes / (1024 * 1024)).toFixed(1) + 'MB'
}

function formatBody(raw: string): string {
  try { return JSON.stringify(JSON.parse(raw), null, 2) } catch { return raw }
}

function isHtml(raw: string): boolean {
  return /^\s*</.test(raw) && /<\/?\w+/.test(raw)
}

function copyCurl() {
  if (props.response?.curlCommand) {
    navigator.clipboard.writeText(props.response.curlCommand)
  }
}

function copyResponse() {
  if (props.response?.body) {
    navigator.clipboard.writeText(props.response.body)
  }
}
</script>

<style scoped>
.response-panel {
  display: flex;
  flex-direction: column;
  min-height: 100px;
  background: var(--bg-base);
  border-top: 1px solid var(--border-primary);
}
.status-bar {
  display: flex;
  padding: 6px 12px;
  background: var(--bg-surface);
  border-bottom: 1px solid var(--border-primary);
  align-items: center;
  gap: 10px;
}
.status-badge {
  font-weight: 700;
  font-size: var(--fs-sm);
  padding: 2px 8px;
  border-radius: var(--radius-sm);
  font-family: var(--font-mono);
  letter-spacing: 0.3px;
}
.status-badge.green { background: var(--accent-soft); color: var(--accent); }
.status-badge.blue { background: var(--blue-soft); color: var(--blue); }
.status-badge.orange { background: var(--amber-soft); color: var(--amber); }
.status-badge.red { background: var(--red-soft); color: var(--red); }
.meta { font-size: var(--fs-xs); color: var(--text-muted); font-family: var(--font-mono); }
.time-badge { opacity: 0.5; }
.placeholder { font-size: var(--fs-xs); color: var(--text-muted); font-family: var(--font-mono); }
.sub-tabs {
  display: flex;
  border-bottom: 1px solid var(--border-primary);
  background: var(--bg-surface);
  padding: 0 8px;
}
.sub-tabs button {
  padding: 6px 12px;
  font-size: var(--fs-sm);
  cursor: pointer;
  color: var(--text-muted);
  border: none;
  background: transparent;
  border-bottom: 2px solid transparent;
  outline: none;
  font-family: var(--font-mono);
  font-weight: 500;
  transition: all var(--transition);
}
.sub-tabs button.active {
  color: var(--accent);
  border-bottom-color: var(--accent);
  font-weight: 600;
}
.sub-tabs button:hover:not(.active) { color: var(--text-secondary); }
.body-content {
  flex: 1;
  overflow-y: auto;
  padding: 10px 12px;
  background: var(--bg-base);
}
.body-empty {
  display: flex; align-items: center; justify-content: center; padding: 24px;
}
.body-modes {
  display: flex; gap: 3px; margin-bottom: 8px; align-items: center;
}
.body-modes button {
  padding: 2px 10px; border: 1px solid var(--border-primary); border-radius: var(--radius-sm);
  background: var(--bg-surface); font-size: var(--fs-xs); cursor: pointer;
  color: var(--text-muted); outline: none; font-family: var(--font-mono);
  transition: all var(--transition);
}
.body-modes button.active {
  border-color: var(--accent); color: var(--accent); font-weight: 600;
  background: var(--accent-soft);
}
.copy-btn {
  padding: 2px 8px; border: 1px solid var(--border-primary); border-radius: var(--radius-sm);
  background: var(--bg-surface); font-size: var(--fs-xs); cursor: pointer;
  color: var(--text-muted); font-family: var(--font-mono);
  transition: all var(--transition);
}
.copy-btn:hover { color: var(--accent); border-color: var(--accent); }
.code-block {
  font-family: var(--font-mono);
  font-size: var(--fs-sm); line-height: 1.7; margin: 0;
  white-space: pre-wrap; word-break: break-all;
  color: var(--text-primary);
  background: var(--bg-surface);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  padding: 10px;
}
.code-block.raw { color: var(--text-secondary); }
.kv-table { width: 100%; border-collapse: collapse; font-size: var(--fs-sm); }
.kv-table th, .kv-table td { text-align: left; padding: 4px 8px; border-bottom: 1px solid var(--border-primary); }
.kv-table th { color: var(--text-muted); font-weight: 500; text-transform: uppercase; font-size: var(--fs-2xs); letter-spacing: 0.5px; }
.kv-table td { color: var(--text-primary); font-family: var(--font-mono); font-size: var(--fs-sm); }
.key-cell { color: var(--accent) !important; }
.preview-hint { color: var(--text-muted); font-size: var(--fs-sm); padding: 16px; }
.preview-iframe { width: 100%; height: 400px; border: none; background: #fff; border-radius: var(--radius); }
.hint { color: var(--text-muted); font-size: var(--fs-sm); }
.copy-curl-btn {
  padding: 3px 10px; border: 1px solid var(--border-primary); border-radius: var(--radius-sm);
  background: var(--bg-surface); font-size: var(--fs-xs); cursor: pointer;
  color: var(--text-muted); margin-bottom: 6px; font-family: var(--font-mono);
  transition: all var(--transition);
}
.copy-curl-btn:hover { border-color: var(--accent); color: var(--accent); }
</style>
