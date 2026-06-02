<template>
  <div class="response-panel">
    <div class="status-bar">
      <n-tag v-if="response" :type="statusType" size="small">
        {{ response.status }} {{ statusText }}
      </n-tag>
      <span v-if="response" class="meta">{{ response.time }} ms</span>
      <span v-if="response" class="meta">{{ formatSize(response.size) }}</span>
      <span v-if="!response" class="placeholder">点击 Send 发送请求</span>
    </div>

    <div class="sub-tabs">
      <n-button
        v-for="tab in bodyTabs"
        :key="tab"
        :type="activeTab === tab ? 'primary' : 'default'"
        size="small"
        :ghost="activeTab !== tab"
        @click="activeTab = tab"
      >
        {{ tab }}
      </n-button>
    </div>

    <div v-if="!response" class="body-empty">
      <span class="hint">响应体将在此显示</span>
    </div>

    <div v-else class="body-content">
      <div v-if="activeTab === 'Body'" class="body-modes">
        <n-radio-group v-model:value="bodyMode" size="small">
          <n-radio-button value="pretty">Pretty</n-radio-button>
          <n-radio-button value="raw">Raw</n-radio-button>
          <n-radio-button value="preview">Preview</n-radio-button>
        </n-radio-group>

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
        <n-button size="small" @click="copyCurl">复制 cURL</n-button>
        <pre class="code-block">{{ response.rawRequest || '(暂无)' }}</pre>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { NTag, NButton, NRadioGroup, NRadioButton } from 'naive-ui'
import type { HttpResponse } from '../../types/response'

const props = defineProps<{
  response: HttpResponse | null
}>()

const activeTab = ref('Body')
const bodyMode = ref('pretty')

const bodyTabs = ['Body', 'Headers', 'Cookies', 'Log']

const statusType = computed(() => {
  if (!props.response) return 'default'
  const s = props.response.status
  if (s < 300) return 'success'
  if (s < 400) return 'info'
  if (s < 500) return 'warning'
  return 'error'
})

const statusText = computed(() => {
  if (!props.response) return ''
  const s = props.response.status
  if (s === 200) return 'OK'
  if (s === 201) return 'Created'
  if (s === 204) return 'No Content'
  if (s === 400) return 'Bad Request'
  if (s === 401) return 'Unauthorized'
  if (s === 404) return 'Not Found'
  if (s === 500) return 'Server Error'
  return ''
})

function formatSize(bytes: number): string {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

function formatBody(raw: string): string {
  try {
    return JSON.stringify(JSON.parse(raw), null, 2)
  } catch {
    return raw
  }
}

function copyCurl() {
  if (props.response?.curlCommand) {
    navigator.clipboard.writeText(props.response.curlCommand)
  }
}
</script>

<style scoped>
.response-panel {
  border-top: 1px solid #e0e0e0;
  display: flex;
  flex-direction: column;
}
.status-bar {
  display: flex;
  padding: 6px 12px;
  background: #f8f9fa;
  border-bottom: 1px solid #eee;
  align-items: center;
  gap: 12px;
}
.meta {
  font-size: 11px;
  color: #999;
}
.placeholder {
  font-size: 11px;
  color: #aaa;
}
.sub-tabs {
  display: flex;
  padding: 4px 10px;
  gap: 4px;
  border-bottom: 1px solid #eee;
  background: #fafafa;
}
.body-content {
  flex: 1;
  overflow-y: auto;
  max-height: 200px;
  padding: 8px;
}
.body-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 30px;
}
.body-modes {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.code-block {
  font-family: 'SF Mono', 'Consolas', monospace;
  font-size: 12px;
  line-height: 1.6;
  margin: 0;
  white-space: pre-wrap;
  word-break: break-all;
  color: #333;
}
.code-block.raw {
  color: #555;
  font-size: 11px;
}
.kv-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 11px;
}
.kv-table th, .kv-table td {
  text-align: left;
  padding: 4px 8px;
  border-bottom: 1px solid #f0f0f0;
}
.kv-table th {
  color: #999;
  font-weight: 500;
}
.preview-hint {
  color: #aaa;
  font-size: 12px;
  padding: 20px;
}
.hint {
  color: #aaa;
  font-size: 12px;
}
</style>
