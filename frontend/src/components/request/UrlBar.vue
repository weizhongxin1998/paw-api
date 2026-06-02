<template>
  <div class="url-bar">
    <select v-model="method" class="method-select">
      <option v-for="opt in methodOptions" :key="opt" :value="opt">{{ opt }}</option>
    </select>
    <div class="url-input-wrapper" @click="startEdit">
      <input
        v-if="isEditing"
        ref="inputRef"
        v-model="url"
        class="url-input-editing"
        placeholder="https://api.example.com/v1/users"
        @blur="stopEdit"
        @keydown.enter="onEnter"
      />
      <div v-else class="url-display" @mouseover="onVarHover" @mouseleave="onVarLeave">
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
import { ref, computed, nextTick } from 'vue'
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

const url = computed({
  get: () => props.modelUrl,
  set: (v) => emit('update:modelUrl', v),
})

const urlSegments = computed<UrlSegment[]>(() => {
  const val = props.modelUrl || ''
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

async function startEdit() {
  isEditing.value = true
  await nextTick()
  inputRef.value?.focus()
}

function stopEdit() {
  isEditing.value = false
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
  border-bottom: 1px solid #eee;
}
.method-select {
  width: 80px;
  padding: 7px 8px;
  border: 1px solid #d0d0d0;
  border-right: none;
  border-radius: 6px 0 0 6px;
  font-family: 'SF Mono', 'Consolas', monospace;
  font-size: 12px;
  font-weight: 600;
  background: #fafafa;
  cursor: pointer;
  outline: none;
  appearance: none;
  -webkit-appearance: none;
  text-align: center;
  text-align-last: center;
}
.method-select:focus {
  border-color: #18a058;
}
.url-input-wrapper {
  flex: 1;
  position: relative;
}
.url-input-editing {
  width: 100%;
  padding: 7px 10px;
  border: 1px solid #d0d0d0;
  border-left: none;
  border-right: none;
  font-family: 'SF Mono', 'Consolas', monospace;
  font-size: 12px;
  outline: none;
  box-sizing: border-box;
}
.url-input-editing:focus {
  border-color: #18a058;
}
.url-display {
  width: 100%;
  height: 34px;
  display: flex;
  align-items: center;
  padding: 0 10px;
  border: 1px solid #d0d0d0;
  border-left: none;
  border-right: none;
  font-size: 13px;
  overflow: hidden;
  white-space: nowrap;
  cursor: text;
  background: #fff;
  box-sizing: border-box;
}
.var-highlight {
  background: #fff3cd;
  border-radius: 3px;
  padding: 0 2px;
  cursor: pointer;
  color: #856404;
}
.var-tooltip {
  position: absolute;
  bottom: 100%;
  left: 60px;
  margin-bottom: 4px;
  background: #333;
  color: #fff;
  padding: 3px 8px;
  border-radius: 4px;
  font-size: 10px;
  white-space: nowrap;
  z-index: 100;
}
.send-btn {
  padding: 7px 22px;
  background: #18a058;
  color: #fff;
  border: 1px solid #18a058;
  border-radius: 0 6px 6px 0;
  font-weight: 600;
  cursor: pointer;
  font-size: 12px;
  white-space: nowrap;
}
.send-btn:hover {
  background: #0c7a43;
  border-color: #0c7a43;
}
</style>
