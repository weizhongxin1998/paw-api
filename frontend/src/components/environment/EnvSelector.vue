<template>
  <div class="env-selector">
    <n-select
      v-model:value="selectedId"
      :options="envOptions"
      size="small"
      placeholder="选择环境"
      style="width: 140px"
      @update:value="onSelect"
    />
    <n-button text size="small" @click="showManager = true" title="管理环境">
      <template #icon>
        <span style="font-size:14px">&#9881;</span>
      </template>
    </n-button>

    <EnvManagerModal
      v-model:show="showManager"
      :project-id="projectId"
      @refresh="fetchEnvs"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { NSelect, NButton } from 'naive-ui'
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

const envOptions = computed(() => {
  return environments.value.map(e => ({
    label: e.name + (e.is_active ? ' ✓' : ''),
    value: e.id,
  }))
})

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

function onSelect(id: number) {
  ActivateEnvironment(id).catch(() => {})
  emit('update:activeEnvId', id)
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
</style>
