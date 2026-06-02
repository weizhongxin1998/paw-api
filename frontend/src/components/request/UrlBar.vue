<template>
  <div class="url-bar">
    <select v-model="method" class="method-select">
      <option v-for="opt in methodOptions" :key="opt" :value="opt">{{ opt }}</option>
    </select>
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
      <div v-if="hoveredVar && resolvedVarText" class="var-tooltip">{{ resolvedVarText }}</div>
    </div>
    <button class="send-btn" @click="$emit('send')">Send</button>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
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

const urlDisplay = computed(() => {
  return relativeURL.value
})

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

const methodOptions = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'HEAD', 'OPTIONS']

function fullURL(relative: string): string {
  const prefix = baseURLPrefix.value
  if (!prefix) return relative
  if (relative.startsWith('http://') || relative.startsWith('https://')) return relative
  const p = prefix.replace(/\/$/, '')
  const r = relative.startsWith('/') ? relative : '/' + relative
  return p + r
}

const url = computed({
  get: () => props.modelUrl,
  set: (v) => emit('update:modelUrl', v),
})

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
        resolvedVarText.value = '解析失败'
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
  border-bottom: 1px solid var(--gray-200);
  gap: 0;
}
.method-select {
  width: 84px;
  padding: 7px 8px;
  border: 1px solid var(--gray-300);
  border-right: none;
  border-radius: var(--radius) 0 0 var(--radius);
  font-family: 'SF Mono', 'Consolas', monospace;
  font-size: 12px;
  font-weight: 600;
  background: var(--gray-50);
  cursor: pointer;
  outline: none;
  color: var(--gray-700);
  appearance: none;
  -webkit-appearance: none;
  text-align: center;
  text-align-last: center;
  transition: border-color var(--transition);
}
.method-select:focus {
  border-color: var(--green);
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
  background: var(--green-soft);
  color: var(--green);
  font-family: 'SF Mono', 'Consolas', monospace;
  font-size: 12px;
  font-weight: 500;
  border: 1px solid var(--gray-300);
  border-left: none;
  border-right: none;
  white-space: nowrap;
  flex-shrink: 0;
}
.url-input-editing {
  flex: 1;
  padding: 7px 10px;
  border: 1px solid var(--gray-300);
  border-left: none;
  border-right: none;
  font-family: 'SF Mono', 'Consolas', monospace;
  font-size: 12px;
  outline: none;
  box-sizing: border-box;
  min-width: 0;
  color: var(--gray-700);
  transition: border-color var(--transition);
}
.url-input-editing.hasPrefix {
  border-left: none;
}
.url-input-editing:focus {
  border-color: var(--green);
}
.url-display {
  flex: 1;
  height: 32px;
  display: flex;
  align-items: center;
  padding: 0 10px;
  border: 1px solid var(--gray-300);
  border-left: none;
  border-right: none;
  font-size: 12px;
  overflow: hidden;
  white-space: nowrap;
  cursor: text;
  background: #fff;
  box-sizing: border-box;
  min-width: 0;
  color: var(--gray-700);
}
.url-display.hasPrefix {
  border-left: none;
}
.var-highlight {
  background: var(--amber-soft);
  border-radius: 3px;
  padding: 0 3px;
  cursor: pointer;
  color: var(--amber);
}
.var-tooltip {
  position: absolute;
  bottom: 100%;
  left: 60px;
  margin-bottom: 6px;
  background: var(--gray-800);
  color: #fff;
  padding: 4px 10px;
  border-radius: var(--radius-sm);
  font-size: 11px;
  white-space: nowrap;
  z-index: 100;
  box-shadow: var(--shadow-md);
}
.send-btn {
  padding: 7px 20px;
  background: var(--green);
  color: #fff;
  border: 1px solid var(--green);
  border-radius: 0 var(--radius) var(--radius) 0;
  font-weight: 600;
  cursor: pointer;
  font-size: 12px;
  white-space: nowrap;
  transition: all var(--transition);
}
.send-btn:hover {
  background: var(--green-hover);
  border-color: var(--green-hover);
}
.send-btn:active {
  transform: scale(0.98);
}
</style>
