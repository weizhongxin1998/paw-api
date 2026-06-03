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
      <span v-if="isSendingExternal" class="sending-dot"></span>
      {{ isSendingExternal ? 'SEND' : 'SEND' }}
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
  GET: 'var(--accent)',
  POST: 'var(--amber)',
  PUT: 'var(--blue)',
  DELETE: 'var(--red)',
  PATCH: 'var(--purple)',
  HEAD: 'var(--text-secondary)',
  OPTIONS: 'var(--text-secondary)',
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
  padding: 6px 8px;
  gap: 0;
  background: var(--bg-base);
}
.method-select {
  width: 96px;
  flex-shrink: 0;
}
.method-select :deep(.n-base-selection) {
  border-top-right-radius: 0;
  border-bottom-right-radius: 0;
  border-right: none;
}
.method-select :deep(.n-base-selection-label) {
  font-family: var(--font-mono);
  font-weight: 700;
  letter-spacing: 0.5px;
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
  padding: 0 9px;
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
  padding: 6px 9px;
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
.url-input-editing:focus { border-color: var(--accent); }
.url-display {
  flex: 1;
  display: flex;
  align-items: center;
  padding: 0 9px;
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
  border-radius: 2px;
  padding: 0 3px;
  cursor: pointer;
  color: var(--amber);
  font-weight: 600;
}
.var-tooltip {
  position: absolute;
  bottom: 100%;
  left: 60px;
  margin-bottom: 5px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-primary);
  color: var(--text-primary);
  padding: 4px 9px;
  border-radius: var(--radius-sm);
  font-size: var(--fs-xs);
  white-space: nowrap;
  z-index: 100;
  font-family: var(--font-mono);
  box-shadow: 0 4px 12px rgba(0,0,0,0.3);
}
.tooltip-key { color: var(--amber); }
.send-btn {
  padding: 6px 20px;
  background: var(--accent);
  color: #000;
  border: 1px solid var(--accent);
  border-radius: 0 var(--radius) var(--radius) 0;
  font-weight: 700;
  cursor: pointer;
  font-size: var(--fs-sm);
  white-space: nowrap;
  font-family: var(--font-mono);
  letter-spacing: 1px;
  transition: all var(--transition);
  position: relative;
  overflow: hidden;
}
.send-btn::before {
  content: ''; position: absolute; inset: 0;
  background: linear-gradient(90deg, transparent, rgba(255,255,255,0.15), transparent);
  transform: translateX(-100%);
  transition: transform 0.4s ease;
}
.send-btn:hover::before { transform: translateX(100%); }
.send-btn:hover { background: var(--accent-hover); border-color: var(--accent-hover); box-shadow: 0 0 12px var(--accent-glow); }
.send-btn:active { transform: scale(0.97); }
.sending-dot {
  display: inline-block;
  width: 6px; height: 6px;
  background: #000;
  border-radius: 50%;
  animation: pulse 0.6s ease-in-out infinite;
  margin-right: 3px;
}
@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.3; }
}
</style>

<style>
.urlbar-method-opt-get { color: var(--accent) !important; font-weight: 700 !important; }
.urlbar-method-opt-post { color: var(--amber) !important; font-weight: 700 !important; }
.urlbar-method-opt-put { color: var(--blue) !important; font-weight: 700 !important; }
.urlbar-method-opt-delete { color: var(--red) !important; font-weight: 700 !important; }
.urlbar-method-opt-patch { color: var(--purple) !important; font-weight: 700 !important; }
.urlbar-method-opt-head { color: var(--text-secondary) !important; font-weight: 700 !important; }
.urlbar-method-opt-options { color: var(--text-secondary) !important; font-weight: 700 !important; }
</style>
