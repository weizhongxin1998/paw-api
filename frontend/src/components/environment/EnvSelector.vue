<template>
  <div class="env-selector">
    <select v-model="selectedId" class="env-select" @change="onSelect($event)">
      <option :value="null" disabled>选择环境</option>
      <option v-for="e in environments" :key="e.id" :value="e.id">
        {{ e.name }}{{ e.is_active ? ' ✓' : '' }}
      </option>
    </select>
    <button class="gear-btn" @click="showManager = true" title="管理环境">&#9881;</button>

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
import type { Environment } from '../../types/environment'
import EnvManagerModal from './EnvManagerModal.vue'

const props = defineProps<{
  projectId: number | null
}>()

const emit = defineEmits<{
  'update:activeEnvId': [value: number | null]
}>()

const selectedId = ref<number | null>(null)
const environments = ref<Environment[]>([])
const showManager = ref(false)

async function fetchEnvs() {
  if (!props.projectId) {
    environments.value = []
    selectedId.value = null
    return
  }
  try {
    const envs = await ListEnvironments(props.projectId)
    environments.value = envs
    const active = envs.find(e => e.is_active)
    if (active) {
      selectedId.value = active.id
      emit('update:activeEnvId', active.id)
    }
  } catch {}
}

function onSelect(_e: Event) {
  if (selectedId.value != null) {
    ActivateEnvironment(selectedId.value).catch(() => {})
    emit('update:activeEnvId', selectedId.value)
  }
}

watch(() => props.projectId, () => {
  fetchEnvs()
}, { immediate: true })
</script>

<style scoped>
.env-selector {
  display: flex;
  align-items: center;
  gap: 2px;
}
.env-select {
  font-size: 11px;
  padding: 4px 10px;
  background: #fff;
  border: 1px solid #ccc;
  border-radius: 6px;
  color: #555;
  cursor: pointer;
  outline: none;
  width: 140px;
  appearance: none;
  -webkit-appearance: none;
}
.env-select:focus {
  border-color: #18a058;
}
.gear-btn {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 14px;
  color: #888;
  padding: 2px 4px;
}
.gear-btn:hover {
  color: #333;
}
</style>
