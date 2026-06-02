import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Environment } from '../types/environment'

export const useEnvStore = defineStore('env', () => {
  const environments = ref<Environment[]>([])
  const activeEnvId = ref<number | null>(null)

  return { environments, activeEnvId }
})
