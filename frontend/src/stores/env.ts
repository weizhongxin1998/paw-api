import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Environment } from '../types/environment'
import { ListEnvironments } from '../../wailsjs/go/main/App'

export const useEnvStore = defineStore('env', () => {
  const environments = ref<Environment[]>([])
  const activeEnvId = ref<number | null>(null)

  const activeEnv = computed(() =>
    environments.value.find(e => e.id === activeEnvId.value) || null
  )

  function getActiveEnvBaseURL(): string {
    return activeEnv.value?.base_url || ''
  }

  async function loadEnvironments(projectId: number) {
    try {
      environments.value = await ListEnvironments(projectId) || []
      const active = environments.value.find(e => e.is_active)
      if (active) activeEnvId.value = active.id
    } catch {
      environments.value = []
    }
  }

  return { environments, activeEnvId, activeEnv, getActiveEnvBaseURL, loadEnvironments }
})
