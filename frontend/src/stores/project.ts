import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Project } from '../types/project'

export const useProjectStore = defineStore('project', () => {
  const projects = ref<Project[]>([])
  const currentId = ref<number | null>(null)

  return { projects, currentId }
})
