<template>
  <div class="env-selector">
    <div class="env-select-wrap">
      <!-- Active environment indicator dot -->
      <span class="env-active-dot" v-if="selectedId" :class="activeEnvType"></span>
      <n-select
        v-model:value="selectedId"
        :options="envOptions"
        placeholder="选择环境"
        size="tiny"
        class="env-select"
        @update:value="onSelect"
        clearable
      />
    </div>
    <n-button text size="tiny" @click="showManager = true" title="管理环境 (Ctrl+E 切换)" class="env-manage-btn">
      <template #icon>
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
          <circle cx="12" cy="12" r="3"/>
          <path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 11-2.83 2.83l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 11-4 0v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 11-2.83-2.83l.06-.06A1.65 1.65 0 004.68 15a1.65 1.65 0 00-1.51-1H3a2 2 0 110-4h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 112.83-2.83l.06.06A1.65 1.65 0 009 4.68a1.65 1.65 0 001-1.51V3a2 2 0 114 0v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 112.83 2.83l-.06.06A1.65 1.65 0 0019.32 9a1.65 1.65 0 001.51 1H21a2 2 0 110 4h-.09a1.65 1.65 0 00-1.51 1z"/>
        </svg>
      </template>
      环境
    </n-button>
    <!-- Visual feedback when switching -->
    <Transition name="env-switch">
      <div v-if="switching" class="env-switch-toast">
        <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="20 6 9 17 4 12"/></svg>
        已切换环境
      </div>
    </Transition>
    <EnvManagerModal
      v-model:show="showManager"
      :project-id="projectId"
      @refresh="fetchEnvs"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted, computed } from 'vue'
import { NSelect, NButton, useMessage } from 'naive-ui'
import { ListEnvironments, ActivateEnvironment } from '../../../wailsjs/go/main/App'
import { useEnvStore } from '../../stores/env'
import type { Environment } from '../../types/environment'
import EnvManagerModal from './EnvManagerModal.vue'

const props = defineProps<{ projectId: number | null }>()
const emit = defineEmits<{ 'update:activeEnvId': [value: number | null] }>()
const envStore = useEnvStore()
const message = useMessage()
const selectedId = ref<number | null>(null)
const environments = ref<Environment[]>([])
const showManager = ref(false)
const switching = ref(false)

const envOptions = ref<{ label: string; value: number }[]>([])

// Detect environment type from name/url for icon coloring
function detectEnvType(env: Environment): string {
  const name = (env.name || '').toLowerCase()
  const url = (env.base_url || '').toLowerCase()
  const combined = name + ' ' + url
  if (combined.includes('dev') || combined.includes('local') || combined.includes('localhost')) return 'dev'
  if (combined.includes('staging') || combined.includes('stg') || combined.includes('test') || combined.includes('uat')) return 'staging'
  if (combined.includes('prod') || combined.includes('production') || combined.includes('live')) return 'prod'
  return 'default'
}

// Current active environment type for the dot indicator
const activeEnvType = computed(() => {
  const env = environments.value.find(e => e.id === selectedId.value)
  return env ? detectEnvType(env) : 'default'
})

async function fetchEnvs() {
  if (!props.projectId) {
    environments.value = []
    envOptions.value = []
    selectedId.value = null
    return
  }
  try {
    const envs = await ListEnvironments(props.projectId) || []
    environments.value = envs
    envOptions.value = envs.map(e => {
      const type = detectEnvType(e)
      const prefix = type === 'dev' ? '[DEV] ' : type === 'staging' ? '[STG] ' : type === 'prod' ? '[PRD] ' : ''
      return {
        label: prefix + e.name,
        value: e.id,
      }
    })
    const active = envs.find(e => e.is_active)
    if (active) {
      selectedId.value = active.id
      emit('update:activeEnvId', active.id)
    }
  } catch {}
}

async function onSelect(val: number | null) {
  if (val == null) {
    envStore.activeEnvId = null
    emit('update:activeEnvId', null)
    return
  }
  try {
    await ActivateEnvironment(val)
    envStore.activeEnvId = val
    await fetchEnvs()
    await envStore.loadEnvironments(props.projectId!)
    showSwitchToast()
  } catch { await fetchEnvs() }
  emit('update:activeEnvId', val)
}

function showSwitchToast() {
  switching.value = true
  setTimeout(() => { switching.value = false }, 1500)
}

// Ctrl+E to cycle through environments
function handleKeydown(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.key === 'e') {
    e.preventDefault()
    cycleEnvironment()
  }
}

async function cycleEnvironment() {
  const envs = environments.value
  if (envs.length <= 1) return

  const currentIdx = envs.findIndex(e => e.id === selectedId.value)
  const nextIdx = (currentIdx + 1) % envs.length
  const nextEnv = envs[nextIdx]

  if (nextEnv) {
    selectedId.value = nextEnv.id
    await onSelect(nextEnv.id)
    message.success(`已切换到: ${nextEnv.name}`)
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
})
onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})

watch(() => props.projectId, () => { fetchEnvs() }, { immediate: true })
</script>

<style scoped>
.env-selector {
  display: flex;
  align-items: center;
  gap: 6px;
  position: relative;
}

/* Active environment indicator dot */
.env-select-wrap {
  position: relative;
  display: flex;
  align-items: center;
}
.env-active-dot {
  position: absolute;
  left: 6px;
  top: 50%;
  transform: translateY(-50%);
  width: 6px;
  height: 6px;
  border-radius: 50%;
  z-index: 2;
  pointer-events: none;
}
.env-active-dot.dev     { background: #3b82f6; box-shadow: 0 0 4px rgba(59,130,246,0.4); }
.env-active-dot.staging { background: #f59e0b; box-shadow: 0 0 4px rgba(245,158,11,0.4); }
.env-active-dot.prod    { background: #22c55e; box-shadow: 0 0 4px rgba(34,197,94,0.4); }
.env-active-dot.default { background: var(--accent, #22c55e); box-shadow: 0 0 4px rgba(34,197,94,0.3); }

.env-select {
  width: 160px;
}
.env-select :deep(.n-base-selection) {
  height: 28px;
  --n-height: 28px !important;
  border-radius: var(--radius-sm);
}
.env-select :deep(.n-base-selection-input) {
  text-align: left;
  padding-left: 20px;
}
.env-select :deep(.n-base-selection-label),
.env-select :deep(.n-base-selection-input) {
  font-size: var(--fs-xs) !important;
}
.env-select :deep(.n-base-selection-input__prefix svg) {
  opacity: 0.4;
}

/* Manage button with settings icon */
.env-manage-btn {
  color: var(--text-muted) !important;
  font-size: var(--fs-xs) !important;
  border-radius: var(--radius-sm) !important;
  height: 28px;
  transition: all 0.15s ease;
}
.env-manage-btn:hover {
  color: var(--accent) !important;
  background: var(--accent-soft) !important;
}

/* Switch toast */
.env-switch-toast {
  position: absolute;
  top: 100%;
  left: 0;
  margin-top: 4px;
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 3px 10px;
  font-size: 9px;
  color: var(--accent);
  background: var(--accent-soft, rgba(0,224,90,0.08));
  border-radius: 6px;
  white-space: nowrap;
  pointer-events: none;
}
.env-switch-enter-active { transition: opacity 0.2s, transform 0.2s; }
.env-switch-leave-active { transition: opacity 0.3s, transform 0.3s; }
.env-switch-enter-from { opacity: 0; transform: translateY(-4px); }
.env-switch-leave-to { opacity: 0; transform: translateY(4px); }
</style>
