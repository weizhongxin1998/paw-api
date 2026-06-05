<template>
  <div
    class="response-panel"
    ref="panelRef"
    tabindex="-1"
    @keydown="onPanelKeydown"
  >
    <!-- Status Bar -->
    <div class="status-bar" :class="'status-bar--' + statusClass">
      <div v-if="response" class="status-bar-inner">
        <span class="status-badge" :class="statusClass">
          {{ response.status }} {{ statusText }}
        </span>
        <span class="meta meta--time">
          <svg class="meta-icon" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="8" cy="8" r="6"/>
            <path d="M8 5v3l2 1.5"/>
          </svg>
          {{ response.time }}ms
        </span>
        <span class="meta meta--size">
          <svg class="meta-icon" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
            <path d="M2 4l6-2 6 2v8l-6 2-6-2V4z"/>
            <path d="M2 4l6 2 6-2"/>
            <path d="M8 6v8"/>
          </svg>
          {{ formatSize(response.size) }}
        </span>
        <span v-if="response.size > SIZE_WARN_THRESHOLD" class="size-warning" title="响应体较大，可能影响渲染">
          <svg class="meta-icon" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
            <path d="M8 1.5l6.5 12H1.5L8 1.5z"/>
            <path d="M8 6v3"/>
            <circle cx="8" cy="11.5" r="0.5" fill="currentColor"/>
          </svg>
          大响应
        </span>
        <span style="flex:1"></span>
        <span class="meta meta--timestamp">{{ new Date().toLocaleTimeString() }}</span>
      </div>
      <div v-if="!response" class="status-bar-inner">
        <span class="placeholder">等待响应...</span>
      </div>
    </div>

    <!-- Sub Tabs -->
    <div class="sub-tabs">
      <button :class="{ active: activeTab === 'Body' }" @click="activeTab = 'Body'">响应体</button>
      <button :class="{ active: activeTab === 'Headers' }" @click="activeTab = 'Headers'">响应头</button>
      <button :class="{ active: activeTab === 'Cookies' }" @click="activeTab = 'Cookies'">Cookies</button>
      <button :class="{ active: activeTab === 'Log' }" @click="activeTab = 'Log'">日志</button>
    </div>

    <!-- Empty State -->
    <div v-if="!response" class="body-empty">
      <div class="empty-terminal">
        <div class="terminal-titlebar">
          <span class="terminal-dot terminal-dot--red"></span>
          <span class="terminal-dot terminal-dot--yellow"></span>
          <span class="terminal-dot terminal-dot--green"></span>
          <span class="terminal-title">response</span>
        </div>
        <div class="terminal-body">
          <span class="terminal-prompt">$</span>
          <span class="terminal-cursor">等待响应...</span>
        </div>
      </div>
    </div>

    <!-- Body Content -->
    <div v-else class="body-content" :class="{ 'copy-flash': showCopyFlash }">
      <div v-if="activeTab === 'Body'">
        <!-- Body Mode Toggle (Segmented Control) + extras -->
        <div class="body-toolbar">
          <div class="segmented-control">
            <button :class="{ active: bodyMode === 'pretty' }" @click="bodyMode = 'pretty'">格式化</button>
            <button :class="{ active: bodyMode === 'raw' }" @click="bodyMode = 'raw'">原始</button>
            <button :class="{ active: bodyMode === 'preview' }" @click="bodyMode = 'preview'">预览</button>
          </div>
          <span style="flex:1"></span>
          <span v-if="bodyMode === 'pretty' && lineCount > 0" class="line-count">{{ lineCount }} 行</span>
          <button
            v-if="bodyMode === 'pretty' || bodyMode === 'raw'"
            class="toolbar-btn"
            :class="{ active: wordWrap }"
            @click="wordWrap = !wordWrap"
            :title="wordWrap ? '关闭自动换行' : '开启自动换行'"
          >
            <svg class="toolbar-icon" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round">
              <path d="M2 4h12"/>
              <path d="M2 8h8c1.1 0 2 .9 2 2s-.9 2-2 2H9"/>
              <path d="M11 10l-2 2 2 2"/>
              <path d="M2 14h5"/>
            </svg>
            换行
          </button>
          <button class="toolbar-btn copy-btn" @click="copyResponse" v-if="response.body">
            <svg v-if="!copiedResponse" class="toolbar-icon" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round">
              <rect x="5" y="5" width="9" height="9" rx="1"/>
              <path d="M11 5V3a1 1 0 00-1-1H3a1 1 0 00-1 1v7a1 1 0 001 1h2"/>
            </svg>
            <svg v-else class="toolbar-icon" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
              <path d="M3 8.5l3 3 7-7"/>
            </svg>
            {{ copiedResponse ? '已复制' : '复制' }}
          </button>
        </div>

        <!-- Pretty Code Block with line numbers -->
        <div v-if="bodyMode === 'pretty'" class="code-container" :class="{ 'code-wrap': wordWrap }">
          <div class="code-line-numbers">
            <div v-for="n in lineCount" :key="n" class="line-num">{{ n }}</div>
          </div>
          <pre class="code-block" :class="{ 'code-nowrap': !wordWrap }"><template
            v-for="(line, idx) in prettyLines" :key="idx"
          ><div class="code-line" :class="{ 'code-line--alt': idx % 2 === 1 }">{{ line }}</div></template></pre>
        </div>

        <!-- Raw Code Block -->
        <pre v-else-if="bodyMode === 'raw'" class="code-block raw" :class="{ 'code-nowrap': !wordWrap }">{{ response.body }}</pre>

        <!-- Preview -->
        <div v-else class="preview-hint">
          <iframe v-if="isHtml(response.body)" :srcdoc="response.body" class="preview-iframe" sandbox="allow-scripts"></iframe>
          <span v-else>非 HTML 响应，无法预览</span>
        </div>
      </div>

      <!-- Headers Tab -->
      <div v-else-if="activeTab === 'Headers'">
        <table class="kv-table">
          <thead><tr><th>键</th><th>值</th></tr></thead>
          <tbody>
            <tr v-for="(v, k) in response.headers" :key="k">
              <td class="key-cell">{{ k }}</td><td class="val-cell">{{ v }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Cookies Tab -->
      <div v-else-if="activeTab === 'Cookies'">
        <table v-if="response.cookies?.length" class="kv-table">
          <thead><tr><th>名称</th><th>值</th><th>域</th><th>路径</th></tr></thead>
          <tbody>
            <tr v-for="c in response.cookies" :key="c.name">
              <td>{{ c.name }}</td><td>{{ c.value }}</td><td>{{ c.domain }}</td><td>{{ c.path }}</td>
            </tr>
          </tbody>
        </table>
        <span v-else class="hint">无 Cookie</span>
      </div>

      <!-- Log Tab -->
      <div v-else-if="activeTab === 'Log'">
        <button class="copy-curl-btn" @click="copyCurl">{{ copiedCurl ? '已复制 ✓' : '复制 cURL' }}</button>
        <pre class="code-block">{{ response.rawRequest || '(暂无)' }}</pre>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { HttpResponse } from '../../types/response'

const SIZE_WARN_THRESHOLD = 1024 * 1024 // 1MB

const props = defineProps<{
  response: HttpResponse | null
}>()

const panelRef = ref<HTMLElement | null>(null)
const activeTab = ref('Body')
const bodyMode = ref('pretty')
const copiedResponse = ref(false)
const copiedCurl = ref(false)
const wordWrap = ref(true)
const showCopyFlash = ref(false)

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
  const map: Record<number, string> = {
    200:'OK', 201:'Created', 202:'Accepted', 204:'No Content',
    301:'Moved Permanently', 302:'Found', 304:'Not Modified',
    400:'Bad Request', 401:'Unauthorized', 403:'Forbidden', 404:'Not Found',
    405:'Method Not Allowed', 408:'Timeout', 409:'Conflict', 413:'Payload Too Large',
    422:'Unprocessable Entity', 429:'Too Many Requests',
    500:'Internal Server Error', 502:'Bad Gateway', 503:'Service Unavailable', 504:'Gateway Timeout',
  }
  return map[s] || ''
})

const prettyBody = computed(() => {
  if (!props.response?.body) return ''
  return formatBody(props.response.body)
})

const prettyLines = computed(() => {
  if (!prettyBody.value) return []
  return prettyBody.value.split('\n')
})

const lineCount = computed(() => prettyLines.value.length)

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
    copiedCurl.value = true
    setTimeout(() => { copiedCurl.value = false }, 1500)
  }
}

function copyResponse() {
  if (props.response?.body) {
    navigator.clipboard.writeText(props.response.body)
    copiedResponse.value = true
    triggerCopyFlash()
    setTimeout(() => { copiedResponse.value = false }, 1500)
  }
}

function triggerCopyFlash() {
  showCopyFlash.value = true
  setTimeout(() => { showCopyFlash.value = false }, 500)
}

// ── Ctrl+C shortcut when panel is focused ──
function onPanelKeydown(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.key === 'c') {
    // Only intercept if no text is selected (let normal copy work for selections)
    const sel = window.getSelection()
    if (!sel || sel.isCollapsed) {
      e.preventDefault()
      copyResponse()
    }
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
  outline: none;
}

/* ── Status Bar ── */
.status-bar {
  display: flex;
  padding: 0;
  background: var(--bg-surface);
  border-bottom: 1px solid var(--border-primary);
  position: relative;
  overflow: hidden;
}
.status-bar::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 4px;
  background: var(--border-primary);
  transition: background 0.3s ease;
}
.status-bar--green::before { background: var(--method-post); }
.status-bar--blue::before { background: var(--blue); }
.status-bar--orange::before { background: var(--amber); }
.status-bar--red::before { background: var(--red); }

.status-bar-inner {
  display: flex;
  padding: 8px 14px 8px 16px;
  align-items: center;
  gap: 14px;
  width: 100%;
}

/* ── Status Badge ── */
.status-badge {
  font-weight: 700;
  font-size: var(--fs-sm);
  padding: 4px 14px;
  border-radius: var(--radius-sm);
  font-family: var(--font-mono);
  letter-spacing: 0.03em;
  position: relative;
}
.status-badge.green {
  background: rgba(34,197,94,0.12);
  color: var(--method-post);
  box-shadow: 0 0 8px rgba(34,197,94,0.15);
}
.status-badge.blue {
  background: var(--blue-soft);
  color: var(--blue);
  box-shadow: 0 0 8px color-mix(in srgb, var(--blue) 20%, transparent);
}
.status-badge.orange {
  background: var(--amber-soft);
  color: var(--amber);
  box-shadow: 0 0 8px color-mix(in srgb, var(--amber) 20%, transparent);
}
.status-badge.red {
  background: var(--red-soft);
  color: var(--red);
  box-shadow: 0 0 8px color-mix(in srgb, var(--red) 20%, transparent);
}

/* ── Meta Info ── */
.meta {
  font-size: var(--fs-xs);
  color: var(--text-muted);
  font-family: var(--font-mono);
  letter-spacing: 0.01em;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}
.meta-icon {
  width: 13px;
  height: 13px;
  flex-shrink: 0;
  opacity: 0.6;
}
.meta--timestamp {
  opacity: 0.45;
}
.size-warning {
  font-size: var(--fs-xs);
  color: var(--amber);
  font-family: var(--font-mono);
  font-weight: 600;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  background: var(--amber-soft);
  padding: 2px 8px;
  border-radius: var(--radius-sm);
}
.placeholder {
  font-size: var(--fs-xs);
  color: var(--text-muted);
  font-family: var(--font-mono);
  font-style: italic;
}

/* ── Sub Tabs ── */
.sub-tabs {
  display: flex;
  border-bottom: 1px solid var(--border-primary);
  background: var(--bg-surface);
  padding: 0 10px;
}
.sub-tabs button {
  padding: 7px 14px;
  font-size: var(--fs-sm);
  cursor: pointer;
  color: var(--text-muted);
  border: none;
  background: transparent;
  border-bottom: 2px solid transparent;
  outline: none;
  font-family: var(--font-mono);
  font-weight: 500;
  letter-spacing: 0.02em;
  transition: all var(--transition);
}
.sub-tabs button.active {
  color: var(--accent);
  border-bottom-color: var(--accent);
  font-weight: 600;
}
.sub-tabs button:hover:not(.active) {
  color: var(--text-secondary);
  background: var(--bg-hover);
  border-radius: var(--radius-sm) var(--radius-sm) 0 0;
}

/* ── Empty State ── */
.body-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px 32px;
}
.empty-terminal {
  background: var(--bg-surface);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  width: 320px;
  overflow: hidden;
  box-shadow: var(--shadow-lg, 0 4px 24px rgba(0,0,0,0.08));
}
.terminal-titlebar {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  background: var(--bg-elevated);
  border-bottom: 1px solid var(--border-primary);
}
.terminal-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}
.terminal-dot--red { background: #ff5f57; }
.terminal-dot--yellow { background: #febc29; }
.terminal-dot--green { background: #28c840; }
.terminal-title {
  margin-left: 8px;
  font-family: var(--font-mono);
  font-size: var(--fs-2xs, 11px);
  color: var(--text-muted);
  opacity: 0.6;
}
.terminal-body {
  padding: 16px 16px 20px;
  font-family: var(--font-mono);
  font-size: var(--fs-sm);
  display: flex;
  align-items: center;
  gap: 8px;
}
.terminal-prompt {
  color: var(--accent);
  font-weight: 700;
}
.terminal-cursor {
  color: var(--text-muted);
  position: relative;
}
.terminal-cursor::after {
  content: '';
  display: inline-block;
  width: 7px;
  height: 14px;
  background: var(--accent);
  margin-left: 2px;
  vertical-align: text-bottom;
  animation: blink 1s step-end infinite;
  opacity: 0.7;
}
@keyframes blink {
  50% { opacity: 0; }
}

/* ── Body Content ── */
.body-content {
  flex: 1;
  overflow-y: auto;
  padding: 12px 14px;
  background: var(--bg-base);
  position: relative;
}
.body-content.copy-flash::after {
  content: '';
  position: absolute;
  inset: 0;
  background: rgba(34,197,94,0.08);
  border: 2px solid rgba(34,197,94,0.25);
  border-radius: var(--radius);
  pointer-events: none;
  animation: flashFade 0.5s ease-out forwards;
}
@keyframes flashFade {
  0% { opacity: 1; }
  100% { opacity: 0; }
}

/* ── Body Toolbar ── */
.body-toolbar {
  display: flex;
  gap: 8px;
  margin-bottom: 10px;
  align-items: center;
}

/* ── Segmented Control ── */
.segmented-control {
  display: inline-flex;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-sm);
  overflow: hidden;
}
.segmented-control button {
  padding: 4px 14px;
  border: none;
  border-right: 1px solid var(--border-primary);
  background: var(--bg-surface);
  font-size: var(--fs-xs);
  cursor: pointer;
  color: var(--text-muted);
  outline: none;
  font-family: var(--font-mono);
  font-weight: 500;
  transition: all var(--transition);
  white-space: nowrap;
}
.segmented-control button:last-child {
  border-right: none;
}
.segmented-control button.active {
  color: var(--accent);
  font-weight: 600;
  background: var(--accent-soft);
}
.segmented-control button:hover:not(.active) {
  background: var(--bg-hover);
  color: var(--text-secondary);
}

/* ── Toolbar Buttons ── */
.toolbar-btn {
  padding: 4px 10px;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-sm);
  background: var(--bg-surface);
  font-size: var(--fs-xs);
  cursor: pointer;
  color: var(--text-muted);
  outline: none;
  font-family: var(--font-mono);
  transition: all var(--transition);
  display: inline-flex;
  align-items: center;
  gap: 4px;
  white-space: nowrap;
}
.toolbar-btn:hover {
  border-color: var(--border-hover);
  color: var(--text-secondary);
}
.toolbar-btn.active {
  border-color: var(--accent);
  color: var(--accent);
  background: var(--accent-soft);
}
.toolbar-icon {
  width: 14px;
  height: 14px;
  flex-shrink: 0;
}
.line-count {
  font-family: var(--font-mono);
  font-size: var(--fs-2xs, 11px);
  color: var(--text-muted);
  opacity: 0.7;
}
.copy-btn:hover {
  color: var(--accent);
  border-color: var(--accent);
  background: var(--accent-soft);
}
.copy-btn:active { transform: scale(0.95); }

/* ── Code Block with Line Numbers ── */
.code-container {
  display: flex;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  overflow: hidden;
  background: var(--bg-surface);
}
.code-line-numbers {
  flex-shrink: 0;
  padding: 12px 0;
  background: var(--bg-elevated);
  border-right: 1px solid var(--border-primary);
  user-select: none;
  min-width: 36px;
  text-align: right;
}
.line-num {
  font-family: var(--font-mono);
  font-size: var(--fs-xs);
  line-height: 1.7;
  padding: 0 8px 0 6px;
  color: var(--text-muted);
  opacity: 0.4;
}
.code-block {
  font-family: var(--font-mono);
  font-size: var(--fs-sm);
  line-height: 1.7;
  margin: 0;
  white-space: pre-wrap;
  word-break: break-all;
  color: var(--text-primary);
  background: transparent;
  padding: 12px 14px;
  flex: 1;
  min-width: 0;
  overflow-x: auto;
}
.code-block.code-nowrap {
  white-space: pre;
  word-break: normal;
}
.code-block.raw {
  color: var(--text-secondary);
  background: var(--bg-surface);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  padding: 12px 14px;
}

/* ── Alternating Line Backgrounds ── */
.code-line {
  min-height: 1.7em;
}
.code-line--alt {
  background: color-mix(in srgb, var(--bg-elevated) 50%, transparent);
}

/* ── Tables ── */
.kv-table {
  width: 100%;
  border-collapse: separate;
  border-spacing: 0;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  overflow: hidden;
  font-size: var(--fs-sm);
}
.kv-table th {
  text-align: left; padding: 6px 10px;
  border-bottom: 1px solid var(--border-primary);
  background: var(--bg-elevated);
}
.kv-table th {
  color: var(--text-muted);
  font-weight: 600;
  text-transform: uppercase;
  font-size: var(--fs-2xs);
  letter-spacing: 0.05em;
  font-family: var(--font-ui);
}
.kv-table td {
  text-align: left; padding: 5px 10px;
  border-bottom: 1px solid var(--border-subtle);
  color: var(--text-primary);
  font-family: var(--font-mono);
  font-size: var(--fs-sm);
}
.kv-table tr:last-child td { border-bottom: none; }
.kv-table tr:hover td { background: var(--bg-hover); }
.key-cell { color: var(--accent) !important; font-weight: 500; }

.preview-hint { color: var(--text-muted); font-size: var(--fs-sm); padding: 20px; }
.preview-iframe {
  width: 100%; height: 400px; border: 1px solid var(--border-primary);
  background: var(--bg-surface); border-radius: var(--radius);
}
.hint {
  color: var(--text-muted); font-size: var(--fs-sm);
  font-family: var(--font-mono);
  font-style: italic;
}
.copy-curl-btn {
  padding: 4px 12px;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-sm);
  background: var(--bg-surface);
  font-size: var(--fs-xs);
  cursor: pointer;
  color: var(--text-muted);
  margin-bottom: 8px;
  font-family: var(--font-mono);
  transition: all var(--transition);
}
.copy-curl-btn:hover {
  border-color: var(--accent);
  color: var(--accent);
  background: var(--accent-soft);
}
.copy-curl-btn:active { transform: scale(0.95); }
</style>
