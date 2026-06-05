<template>
  <div class="url-bar" :class="{ 'url-bar--focused': isBarFocused }" @focusin="isBarFocused = true" @focusout="onBarFocusOut">
    <div class="method-select" :style="{ '--method-accent': methodColors[method] || 'var(--accent)' }">
      <n-select
        v-model:value="method"
        :options="methodSelectOptions"
        :render-label="renderMethodLabel"
        size="small"
        :consistent-menu-width="false"
      />
    </div>
    <div class="url-input-wrapper" @click="startEdit">
      <span v-if="baseURLPrefix" class="url-prefix">{{ baseURLPrefix }}</span>
      <input
        v-if="isEditing"
        ref="inputRef"
        v-model="editingURL"
        class="url-input-editing"
        :class="{ hasPrefix: baseURLPrefix }"
        placeholder="/v1/users"
        @blur="stopEdit"
        @keydown.enter="onEnter"
        @focus="onInputFocus"
        @input="onInputChanged"
      />
      <div v-else class="url-display" :class="{ hasPrefix: baseURLPrefix }" @mouseover="onVarHover" @mouseleave="onVarLeave">
        <template v-for="(seg, i) in urlSegments" :key="i">
          <span v-if="seg.type === 'var'" class="var-highlight">{{ seg.text }}</span>
          <span v-else>{{ seg.text }}</span>
        </template>
      </div>
      <div v-if="hoveredVar && resolvedVarText" class="var-tooltip" v-text="varTooltipText"></div>
      <!-- URL History Dropdown -->
      <div v-if="showHistory && filteredHistory.length > 0" class="url-history-dropdown" @mousedown.prevent>
        <div
          v-for="(item, idx) in filteredHistory"
          :key="idx"
          class="url-history-item"
          @click="selectHistory(item)"
        >
          <span class="history-method" :style="{ color: methodColors[item.method] || 'inherit' }">{{ item.method }}</span>
          <span class="history-url">{{ item.url }}</span>
        </div>
      </div>
    </div>
    <button
      class="send-btn"
      :class="{ 'send-btn--loading': isLoading }"
      @click="onSendClick"
      :disabled="isLoading"
      title="发送请求 (Ctrl+Enter)"
    >
      <span v-if="isLoading" class="send-spinner"></span>
      <span class="send-label">发送</span>
      <span class="send-arrow">→</span>
    </button>
    <button
      class="save-btn"
      @click="$emit('save')"
      title="保存接口 (Ctrl+S)"
    >
      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M19 21H5a2 2 0 01-2-2V5a2 2 0 012-2h11l5 5v11z"/>
        <polyline points="17 21 17 13 7 13 7 21"/>
      </svg>
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick, h, onMounted, onBeforeUnmount } from 'vue'
import { NSelect } from 'naive-ui'
import type { SelectOption } from 'naive-ui'
import { ResolveVariable } from '../../../wailsjs/go/main/App'
import { useEnvStore } from '../../stores/env'

interface UrlSegment {
  type: 'text' | 'var'
  text: string
}

interface HistoryEntry {
  method: string
  url: string
}

const HISTORY_KEY = 'paw-api-url-history'
const MAX_HISTORY = 10

const props = defineProps<{
  modelMethod: string
  modelUrl: string
  loading?: boolean
}>()

const emit = defineEmits<{
  'update:modelMethod': [value: string]
  'update:modelUrl': [value: string]
  'send': []
  'save': []
}>()

const envStore = useEnvStore()
const isEditing = ref(false)
const inputRef = ref<HTMLInputElement | null>(null)
const hoveredVar = ref<string | null>(null)
const resolvedVarText = ref('')
let resolveTimer: ReturnType<typeof setTimeout> | null = null

const isBarFocused = ref(false)
const showHistory = ref(false)

const isLoading = computed(() => props.loading ?? false)

// ── URL History ──
const urlHistory = ref<HistoryEntry[]>(loadHistory())

function loadHistory(): HistoryEntry[] {
  try {
    const raw = localStorage.getItem(HISTORY_KEY)
    if (raw) return JSON.parse(raw)
  } catch { /* ignore */ }
  return []
}

function saveHistory() {
  try {
    localStorage.setItem(HISTORY_KEY, JSON.stringify(urlHistory.value))
  } catch { /* ignore */ }
}

function addToHistory(method: string, url: string) {
  const list = urlHistory.value.filter(h => !(h.method === method && h.url === url))
  list.unshift({ method, url })
  urlHistory.value = list.slice(0, MAX_HISTORY)
  saveHistory()
}

const filteredHistory = computed(() => {
  if (!editingURL.value) return urlHistory.value
  const q = editingURL.value.toLowerCase()
  return urlHistory.value.filter(h => h.url.toLowerCase().includes(q))
})

function selectHistory(item: HistoryEntry) {
  emit('update:modelMethod', item.method)
  editingURL.value = item.url
  showHistory.value = false
}

function onInputFocus() {
  if (!editingURL.value || editingURL.value === '/') {
    showHistory.value = true
  }
}

function onInputChanged() {
  if (!editingURL.value) {
    showHistory.value = true
  } else {
    showHistory.value = false
  }
}

function onBarFocusOut(e: FocusEvent) {
  const related = e.relatedTarget as Node | null
  const root = (e.currentTarget as HTMLElement)
  if (!related || !root.contains(related)) {
    isBarFocused.value = false
    showHistory.value = false
  }
}

// ── Ctrl+L keyboard shortcut ──
function onGlobalKeydown(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.key === 'l') {
    e.preventDefault()
    startEdit()
    nextTick(() => {
      inputRef.value?.select()
    })
  }
}

onMounted(() => {
  window.addEventListener('keydown', onGlobalKeydown)
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', onGlobalKeydown)
})

const varTooltipText = computed(() => {
  if (!hoveredVar.value || !resolvedVarText.value) return ''
  return '{{' + hoveredVar.value + '}} = ' + resolvedVarText.value
})

const method = computed({
  get: () => props.modelMethod,
  set: (v) => emit('update:modelMethod', v),
})

const baseURLPrefix = computed(() => {
  if (!envStore.activeEnvId) return ''
  return envStore.getActiveEnvBaseURL() || ''
})

const relativeURL = computed(() => {
  const prefix = baseURLPrefix.value
  if (!prefix) return props.modelUrl || ''
  const u = props.modelUrl || ''
  if (u.startsWith(prefix)) return u.slice(prefix.length)
  if (u.startsWith('http://') || u.startsWith('https://')) return u
  return u
})

const editingURL = ref('')

const urlDisplay = computed(() => relativeURL.value)

const urlSegments = computed<UrlSegment[]>(() => {
  const val = urlDisplay.value || ''
  const segments: UrlSegment[] = []
  let lastIndex = 0
  const re = /\{\{(\w+)\}\}/g
  let match: RegExpExecArray | null
  while ((match = re.exec(val)) !== null) {
    if (match.index > lastIndex) {
      segments.push({ type: 'text', text: val.slice(lastIndex, match.index) })
    }
    segments.push({ type: 'var', text: match[0] })
    lastIndex = match.index + match[0].length
  }
  if (lastIndex < val.length) {
    segments.push({ type: 'text', text: val.slice(lastIndex) })
  }
  return segments
})

const methodColors: Record<string, string> = {
  GET: 'var(--method-get)',
  POST: 'var(--method-post)',
  PUT: 'var(--method-put)',
  DELETE: 'var(--method-delete)',
  PATCH: 'var(--method-patch)',
  HEAD: 'var(--method-head)',
  OPTIONS: 'var(--method-options)',
}

const methodSelectOptions: SelectOption[] = [
  'GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'HEAD', 'OPTIONS',
].map(m => ({
  label: m,
  value: m,
  class: 'urlbar-method-opt-' + m.toLowerCase(),
}))

function renderMethodLabel(option: SelectOption) {
  const m = option.value as string
  return h('span', { style: { color: methodColors[m] || 'inherit', fontWeight: '700', fontFamily: 'var(--font-mono)', fontSize: 'var(--fs-sm)', letterSpacing: '0.5px' } }, m)
}

function fullURL(relative: string): string {
  const prefix = baseURLPrefix.value
  if (!prefix) return relative
  if (relative.startsWith('http://') || relative.startsWith('https://')) return relative
  const p = prefix.replace(/\/$/, '')
  const r = relative.startsWith('/') ? relative : '/' + relative
  return p + r
}

async function startEdit() {
  editingURL.value = relativeURL.value
  isEditing.value = true
  await nextTick()
  inputRef.value?.focus()
}

function stopEdit() {
  isEditing.value = false
  showHistory.value = false
  emit('update:modelUrl', fullURL(editingURL.value || '/'))
}

function onEnter() {
  stopEdit()
  // Don't send if URL is empty
  const url = fullURL(editingURL.value || '')
  if (!url || url === '/' || url.trim() === '') return
  // Add to history
  addToHistory(method.value, url)
  emit('send')
}

function onSendClick() {
  if (isLoading.value) return
  // If editing, commit the URL first
  if (isEditing.value) {
    stopEdit()
  }
  const url = props.modelUrl || ''
  if (!url || url === '/' || url.trim() === '') return
  addToHistory(method.value, url)
  emit('send')
}

function onVarHover(e: MouseEvent) {
  const target = e.target as HTMLElement
  if (target.classList.contains('var-highlight')) {
    const varName = target.textContent || ''
    const match = varName.match(/^\{\{(\w+)\}\}$/)
    if (!match) return
    hoveredVar.value = match[1]
    if (resolveTimer) clearTimeout(resolveTimer)
    resolveTimer = setTimeout(async () => {
      const envId = envStore.activeEnvId ?? 0
      try {
        resolvedVarText.value = await ResolveVariable(varName, envId)
      } catch {
        resolvedVarText.value = '(解析失败)'
      }
    }, 200)
  }
}

function onVarLeave(e: MouseEvent) {
  const target = e.target as HTMLElement
  if (target.classList.contains('var-highlight')) {
    hoveredVar.value = null
    resolvedVarText.value = ''
    if (resolveTimer) clearTimeout(resolveTimer)
  }
}
</script>

<style scoped>
.url-bar {
  display: flex;
  padding: 10px 12px;
  gap: 0;
  background: var(--bg-base);
  border: 2px solid transparent;
  border-radius: var(--radius);
  transition: border-color 0.25s ease, box-shadow 0.25s ease;
}
.url-bar--focused {
  border-color: color-mix(in srgb, var(--accent) 40%, transparent);
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--accent) 12%, transparent),
              0 0 12px -2px color-mix(in srgb, var(--accent) 20%, transparent);
}

/* ── Method Select ── */
.method-select {
  width: 100px;
  flex-shrink: 0;
  position: relative;
}
.method-select::before {
  content: '';
  position: absolute;
  left: 0;
  top: 2px;
  bottom: 2px;
  width: 3px;
  border-radius: 3px 0 0 3px;
  background: var(--method-accent, var(--accent));
  z-index: 2;
  transition: background 0.2s ease;
}
.method-select :deep(.n-base-selection) {
  border-top-right-radius: 0;
  border-bottom-right-radius: 0;
  border-right: none;
  height: 34px;
  padding-left: 8px;
}
.method-select :deep(.n-base-selection-label) {
  font-family: var(--font-mono);
  font-weight: 700;
  letter-spacing: 0.04em;
}

/* ── URL Input Wrapper ── */
.url-input-wrapper {
  flex: 1;
  position: relative;
  display: flex;
  align-items: stretch;
}

/* ── URL Prefix Chip ── */
.url-prefix {
  display: flex;
  align-items: center;
  padding: 0 10px 0 14px;
  background: var(--accent-soft);
  color: var(--accent);
  font-family: var(--font-mono);
  font-size: var(--fs-sm);
  font-weight: 500;
  border: 1px solid var(--border-primary);
  border-left: none;
  border-right: none;
  white-space: nowrap;
  flex-shrink: 0;
  border-radius: 0;
  position: relative;
}
.url-prefix::before {
  content: '';
  position: absolute;
  left: 0;
  top: 20%;
  bottom: 20%;
  width: 2px;
  background: var(--accent);
  border-radius: 2px;
  opacity: 0.5;
}

.url-input-editing {
  flex: 1;
  padding: 7px 10px;
  border: 1px solid var(--border-primary);
  border-left: none;
  border-right: none;
  font-family: var(--font-mono);
  font-size: var(--fs-sm);
  outline: none;
  box-sizing: border-box;
  min-width: 0;
  color: var(--text-primary);
  background: var(--bg-surface);
  transition: border-color var(--transition), background var(--transition);
}
.url-input-editing:focus {
  border-color: var(--accent);
}

/* ── URL Display ── */
.url-display {
  flex: 1;
  display: flex;
  align-items: center;
  padding: 0 10px;
  border: 1px solid var(--border-primary);
  border-left: none;
  border-right: none;
  font-size: var(--fs-sm);
  overflow: hidden;
  white-space: nowrap;
  cursor: text;
  background: var(--bg-surface);
  box-sizing: border-box;
  min-width: 0;
  color: var(--text-primary);
  font-family: var(--font-mono);
  transition: background 0.3s ease;
}
.url-display:hover {
  background: linear-gradient(90deg, var(--bg-surface) 0%, color-mix(in srgb, var(--accent) 4%, var(--bg-surface)) 50%, var(--bg-surface) 100%);
}

.var-highlight {
  background: var(--amber-soft);
  border: 1px solid rgba(245,158,11,0.15);
  border-radius: var(--radius-xs);
  padding: 0 4px;
  cursor: pointer;
  color: var(--amber);
  font-weight: 600;
  transition: all var(--transition-fast);
}
.var-highlight:hover {
  background: rgba(245,158,11,0.15);
  border-color: var(--amber);
}
.var-tooltip {
  position: absolute;
  bottom: 100%;
  left: 60px;
  margin-bottom: 6px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-hover);
  color: var(--text-primary);
  padding: 5px 10px;
  border-radius: var(--radius);
  font-size: var(--fs-xs);
  white-space: nowrap;
  z-index: 100;
  font-family: var(--font-mono);
  box-shadow: var(--shadow-lg);
  animation: fadeInScale 0.15s var(--ease-out) both;
}
.tooltip-key { color: var(--amber); }

/* ── URL History Dropdown ── */
.url-history-dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: var(--bg-elevated);
  border: 1px solid var(--border-hover);
  border-top: none;
  border-radius: 0 0 var(--radius) var(--radius);
  box-shadow: var(--shadow-lg);
  z-index: 200;
  max-height: 260px;
  overflow-y: auto;
  animation: slideDown 0.15s var(--ease-out) both;
}
@keyframes slideDown {
  from { opacity: 0; transform: translateY(-4px); }
  to   { opacity: 1; transform: translateY(0); }
}
.url-history-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px 12px;
  cursor: pointer;
  font-family: var(--font-mono);
  font-size: var(--fs-xs);
  transition: background 0.12s ease;
}
.url-history-item:hover {
  background: var(--bg-hover);
}
.history-method {
  font-weight: 700;
  font-size: var(--fs-2xs);
  letter-spacing: 0.04em;
  min-width: 52px;
  flex-shrink: 0;
}
.history-url {
  color: var(--text-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* ── Send Button ── */
.send-btn {
  padding: 0 20px;
  background: var(--accent);
  color: var(--text-inverse);
  border: 1px solid var(--accent);
  border-right: none;
  border-radius: 0;
  font-weight: 700;
  cursor: pointer;
  font-size: var(--fs-sm);
  white-space: nowrap;
  font-family: var(--font-mono);
  letter-spacing: 0.08em;
  transition: all var(--transition);
  position: relative;
  overflow: hidden;
  height: 34px;
  display: flex;
  align-items: center;
  gap: 6px;
}
.send-btn::before {
  content: '';
  position: absolute; inset: 0;
  background: linear-gradient(90deg, transparent, rgba(255,255,255,0.12), transparent);
  transform: translateX(-100%);
  transition: transform 0.5s ease;
}
.send-btn:hover::before { transform: translateX(100%); }
.send-btn:hover {
  background: var(--accent-hover);
  border-color: var(--accent-hover);
  box-shadow: var(--shadow-glow);
}
.send-btn:active { transform: scale(0.97); }
.send-btn:disabled { opacity: 0.6; cursor: not-allowed; }
.send-btn--loading { pointer-events: none; }
.send-label { line-height: 1; }
.send-arrow {
  font-size: var(--fs-md);
  line-height: 1;
  opacity: 0.85;
  transition: transform 0.2s ease;
}
.send-btn:hover .send-arrow {
  transform: translateX(2px);
}

/* ── Loading Spinner ── */
.send-spinner {
  display: inline-block;
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: currentColor;
  border-radius: 50%;
  animation: spin 0.65s linear infinite;
  flex-shrink: 0;
}
@keyframes spin {
  to { transform: rotate(360deg); }
}

.sending-indicator {
  display: flex; gap: 3px; align-items: center;
}
.sending-dot {
  display: inline-block;
  width: 5px; height: 5px;
  background: currentColor;
  border-radius: 50%;
  animation: sendingPulse 0.8s ease-in-out infinite;
}
@keyframes sendingPulse {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.3; transform: scale(0.7); }
}

/* ── Save Button ── */
.save-btn {
  padding: 0 12px;
  height: 34px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-surface);
  color: var(--text-muted);
  border: 1px solid var(--border-primary);
  border-left: none;
  border-radius: 0 var(--radius) var(--radius) 0;
  cursor: pointer;
  transition: all var(--transition);
  flex-shrink: 0;
}
.save-btn:hover {
  color: var(--accent);
  border-color: var(--accent);
  background: var(--accent-soft);
}
.save-btn:active {
  transform: scale(0.95);
}
.save-btn svg {
  width: 14px;
  height: 14px;
}
</style>

<style>
.urlbar-method-opt-get { color: var(--method-get) !important; font-weight: 700 !important; }
.urlbar-method-opt-post { color: var(--method-post) !important; font-weight: 700 !important; }
.urlbar-method-opt-put { color: var(--method-put) !important; font-weight: 700 !important; }
.urlbar-method-opt-delete { color: var(--method-delete) !important; font-weight: 700 !important; }
.urlbar-method-opt-patch { color: var(--method-patch) !important; font-weight: 700 !important; }
.urlbar-method-opt-head { color: var(--method-head) !important; font-weight: 700 !important; }
.urlbar-method-opt-options { color: var(--method-options) !important; font-weight: 700 !important; }
</style>
