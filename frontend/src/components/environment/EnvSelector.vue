<template>
  <div class="env-selector">
    <select v-model="selectedId" class="env-select" @change="onSelect">
      <option :value="null" disabled selected>选择环境</option>
      <option v-for="e in environments" :key="e.id" :value="e.id">{{ e.name }}</option>
    </select>
    <button class="manage-btn" @click="showManager = true" title="管理环境">环境</button>
    <EnvManagerModal
      v-model:show="showManager"
      :project-id="projectId"
      @refresh="fetchEnvs"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
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

async function fetchEnvs() {
  if (!props.projectId) { environments.value = []; selectedId.value = null; return }
  try {
    const envs = await ListEnvironments(props.projectId) || []
    environments.value = envs
    const active = envs.find(e => e.is_active)
    if (active) { selectedId.value = active.id; emit('update:activeEnvId', active.id) }
  } catch {}
}

async function onSelect() {
  if (selectedId.value == null) return
  try {
    await ActivateEnvironment(selectedId.value)
    envStore.activeEnvId = selectedId.value
    await fetchEnvs()
    await envStore.loadEnvironments(props.projectId!)
  } catch { await fetchEnvs() }
  emit('update:activeEnvId', selectedId.value)
}

watch(() => props.projectId, () => { fetchEnvs() }, { immediate: true })
</script>

<style scoped>
.env-selector {
  display: flex;
  align-items: center;
  gap: 2px;
}
.env-select {
  font-size: 10px;
  padding: 4px 8px;
  background: var(--bg-base);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  color: var(--text-secondary);
  cursor: pointer;
  outline: none;
  width: 130px;
  appearance: none;
  -webkit-appearance: none;
  font-family: var(--font-mono);
  transition: border-color var(--transition);
}
.env-select:focus { border-color: var(--accent); }
.env-select:hover { border-color: var(--border-hover); }
.manage-btn {
  background: var(--bg-base);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  cursor: pointer;
  font-size: 10px;
  color: var(--text-muted);
  padding: 4px 7px;
  white-space: nowrap;
  font-family: var(--font-mono);
  transition: all var(--transition);
}
.manage-btn:hover { border-color: var(--accent); color: var(--accent); background: var(--accent-soft); }
</style>
