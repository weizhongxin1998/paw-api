<template>
  <div class="url-bar">
    <n-select
      v-model:value="method"
      :options="methodSelectOptions"
      :render-label="renderMethodLabel"
      size="small"
      :consistent-menu-width="false"
      class="method-select"
    />
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
      />
      <div v-else class="url-display" :class="{ hasPrefix: baseURLPrefix }" @mouseover="onVarHover" @mouseleave="onVarLeave">
        <template v-for="(seg, i) in urlSegments" :key="i">
          <span v-if="seg.type === 'var'" class="var-highlight">{{ seg.text }}</span>
          <span v-else>{{ seg.text }}</span>
        </template>
      </div>
      <div v-if="hoveredVar && resolvedVarText" class="var-tooltip" v-text="varTooltipText"></div>
    </div>
    <button class="send-btn" @click="$emit('send')" :disabled="isSendingExternal">
      <span v-if="isSendingExternal" class="sending-indicator">
        <span class="sending-dot"></span>
        <span class="sending-dot" style="animation-delay: 0.15s"></span>
        <span class="sending-dot" style="animation-delay: 0.3s"></span>
      </span>
      <span class="send-label">SEND</span>
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, h } from 'vue'
import { NSelect } from 'naive-ui'
import type { SelectOption } from 'naive-ui'
import { ResolveVariable } from '../../../wailsjs/go/main/App'
import { useEnvStore } from '../../stores/env'

interface UrlSegment {
  type: 'text' | 'var'
  text: string
}

const props = defineProps<{
  modelMethod: string
  modelUrl: string
}>()

const emit = defineEmits<{
  'update:modelMethod': [value: string]
  'update:modelUrl': [value: string]
  'send': []
}>()

const envStore = useEnvStore()
const isEditing = ref(false)
const inputRef = ref<HTMLInputElement | null>(null)
const hoveredVar = ref<string | null>(null)
const resolvedVarText = ref('')
let resolveTimer: ReturnType<typeof setTimeout> | null = null

const isSendingExternal = ref(false)

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
  emit('update:modelUrl', fullURL(editingURL.value || '/'))
}

function onEnter() {
  stopEdit()
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
  padding: 8px 10px;
  gap: 0;
  background: var(--bg-base);
}
.method-select {
  width: 100px;
  flex-shrink: 0;
}
.method-select :deep(.n-base-selection) {
  border-top-right-radius: 0;
  border-bottom-right-radius: 0;
  border-right: none;
  height: 34px;
}
.method-select :deep(.n-base-selection-label) {
  font-family: var(--font-mono);
  font-weight: 700;
  letter-spacing: 0.04em;
}

.url-input-wrapper {
  flex: 1;
  position: relative;
  display: flex;
  align-items: stretch;
}
.url-prefix {
  display: flex;
  align-items: center;
  padding: 0 10px;
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
  transition: border-color var(--transition);
}
.url-input-editing:focus {
  border-color: var(--accent);
}
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

/* ── Send Button ── */
.send-btn {
  padding: 0 22px;
  background: var(--accent);
  color: var(--text-inverse);
  border: 1px solid var(--accent);
  border-radius: 0 var(--radius) var(--radius) 0;
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
.send-label { line-height: 1; }

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
