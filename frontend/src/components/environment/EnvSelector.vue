<template>
  <div class="env-selector">
    <n-select
      v-model:value="selectedId"
      :options="envOptions"
      placeholder="选择环境"
      size="tiny"
      class="env-select"
      @update:value="onSelect"
      clearable
    />
    <n-button text size="tiny" @click="showManager = true" title="管理环境" class="env-manage-btn">
      <template #icon>
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M12 22c5.523 0 10-4.477 10-10S17.523 2 12 2 2 6.477 2 12s4.477 10 10 10z"/><path d="M12 16v-4"/><path d="M12 8h.01"/></svg>
      </template>
      环境
    </n-button>
    <EnvManagerModal
      v-model:show="showManager"
      :project-id="projectId"
      @refresh="fetchEnvs"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { NSelect, NButton } from 'naive-ui'
import { ListEnvironments, ActivateEnvironment } from '../../../wailsjs/go/main/App'
import { useEnvStore } from '../../stores/env'
import type { Environment } from '../../types/environment'
import EnvManagerModal from './EnvManagerModal.vue'

const props = defineProps<{ projectId: number | null }>()
const emit = defineEmits<{ 'update:activeEnvId': [value: number | null] }>()
const envStore = useEnvStore()
const selectedId = ref<number | null>(null)
const environments = ref<Environment[]>([])
const showManager = ref(false)

const envOptions = ref<{ label: string; value: number }[]>([])

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
    envOptions.value = envs.map(e => ({ label: e.name, value: e.id }))
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
  } catch { await fetchEnvs() }
  emit('update:activeEnvId', val)
}

watch(() => props.projectId, () => { fetchEnvs() }, { immediate: true })
</script>

<style scoped>
.env-selector {
  display: flex;
  align-items: center;
  gap: 6px;
}
.env-select {
  width: 150px;
}
.env-select :deep(.n-base-selection) {
  height: 28px;
  --n-height: 28px !important;
  border-radius: var(--radius-sm);
}
/* 居中显示选中的文字 */
.env-select :deep(.n-base-selection-input) {
  text-align: center;
}
/* 统一触发区和下拉菜单的字体大小 */
.env-select :deep(.n-base-selection-label),
.env-select :deep(.n-base-selection-input) {
  font-size: var(--fs-xs) !important;
}
.env-manage-btn {
  color: var(--text-muted) !important;
  font-size: var(--fs-xs) !important;
  border-radius: var(--radius-sm) !important;
  height: 28px;
  transition: all var(--transition);
}
.env-manage-btn:hover {
  color: var(--accent) !important;
  background: var(--accent-soft) !important;
}
</style>
