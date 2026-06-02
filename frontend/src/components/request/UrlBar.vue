<template>
  <div class="url-bar">
    <n-select
      v-model:value="method"
      :options="methodOptions"
      :consistent-menu-width="false"
      size="medium"
      class="method-select"
    />
    <div class="url-input-wrapper" @click="startEdit">
      <n-input
        v-if="isEditing"
        ref="inputRef"
        v-model:value="url"
        placeholder="https://api.example.com/v1/users"
        size="medium"
        class="url-input-editing"
        @blur="stopEdit"
        @keydown.enter="onEnter"
      />
      <div
        v-else
        class="url-display"
        @mouseover="onVarHover"
        @mouseleave="onVarLeave"
      >
        <template v-for="(seg, i) in urlSegments" :key="i">
          <span v-if="seg.type === 'var'" class="var-highlight">{{ seg.text }}</span>
          <span v-else>{{ seg.text }}</span>
        </template>
      </div>
      <n-popover
        v-if="hoveredVar"
        trigger="manual"
        :show="true"
        :x="tooltipX"
        :y="tooltipY"
        placement="top"
        :to="false"
      >
        <template #trigger>
          <span style="position:fixed;display:none"></span>
        </template>
        <span>{{ resolvedVarText || '加载中...' }}</span>
      </n-popover>
    </div>
    <n-button type="primary" class="send-btn" @click="$emit('send')">
      Send
    </n-button>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick } from 'vue'
import { NSelect, NInput, NButton, NPopover } from 'naive-ui'
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
const inputRef = ref<InstanceType<typeof NInput> | null>(null)
const hoveredVar = ref<string | null>(null)
const resolvedVarText = ref('')
const tooltipX = ref(0)
const tooltipY = ref(0)
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

const methodOptions = [
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'DELETE', value: 'DELETE' },
  { label: 'PATCH', value: 'PATCH' },
  { label: 'HEAD', value: 'HEAD' },
  { label: 'OPTIONS', value: 'OPTIONS' },
]

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
    tooltipX.value = e.clientX
    tooltipY.value = e.clientY
    resolvedVarText.value = '加载中...'
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
  gap: 0;
  border-bottom: 1px solid #eee;
}
.method-select {
  width: 90px;
}
.method-select :deep(.n-base-selection) {
  border-radius: 6px 0 0 6px !important;
  border-right: none !important;
}
.url-input-wrapper {
  flex: 1;
  position: relative;
}
.url-input-editing {
  width: 100%;
}
.url-input-editing :deep(.n-input__border) {
  border-radius: 0 !important;
  border-left: none !important;
  border-right: none !important;
}
.url-display {
  width: 100%;
  height: 34px;
  display: flex;
  align-items: center;
  padding: 0 12px;
  border: 1px solid #ddd;
  border-left: none;
  border-right: none;
  font-size: 14px;
  font-family: inherit;
  overflow: hidden;
  white-space: nowrap;
  cursor: text;
  background: #fff;
  color: #333;
}
.var-highlight {
  background: #fff3cd;
  border-radius: 3px;
  padding: 0 2px;
  cursor: pointer;
}
.send-btn {
  border-radius: 0 6px 6px 0 !important;
  font-weight: 600;
  padding: 0 22px;
}
</style>
