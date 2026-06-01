import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Environment } from '../types/environment'

export interface EnvVariable {
  key: string
  value: string
  enabled: boolean
}

export const useEnvironmentStore = defineStore('environment', () => {
  const environments = ref<Environment[]>([])
  const activeEnvironment = ref<Environment | null>(null)

  const activeVariables = computed<EnvVariable[]>(() => {
    if (!activeEnvironment.value) return []
    try {
      return JSON.parse(activeEnvironment.value.variables)
    } catch {
      return []
    }
  })

  function setEnvironments(list: Environment[]) {
    environments.value = list
  }

  function addEnvironment(env: Environment) {
    environments.value.push(env)
  }

  function removeEnvironment(id: string) {
    environments.value = environments.value.filter(e => e.id !== id)
    if (activeEnvironment.value?.id === id) {
      activeEnvironment.value = null
    }
  }

  function updateEnvironment(env: Environment) {
    const idx = environments.value.findIndex(e => e.id === env.id)
    if (idx !== -1) {
      environments.value[idx] = env
    }
    if (activeEnvironment.value?.id === env.id) {
      activeEnvironment.value = env
    }
  }

  function setActiveEnvironment(env: Environment | null) {
    activeEnvironment.value = env
  }

  return {
    environments, activeEnvironment, activeVariables,
    setEnvironments, addEnvironment, removeEnvironment,
    updateEnvironment, setActiveEnvironment,
  }
})
