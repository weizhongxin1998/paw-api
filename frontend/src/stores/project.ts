import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Project } from '../types/project'
import { ListProjects, CreateProject, GetSetting, SetSetting } from '../../wailsjs/go/main/App'

export const useProjectStore = defineStore('project', () => {
  const projects = ref<Project[]>([])
  const currentId = ref<number | null>(null)

  async function loadProjects() {
    try {
      const list = await ListProjects()
      projects.value = list || []
    } catch {
      projects.value = []
    }
  }

  async function loadLastProject() {
    try {
      const lastId = await GetSetting('app.last_project_id')
      if (lastId && projects.value.some(p => p.id === Number(lastId))) {
        currentId.value = Number(lastId)
      } else if (projects.value.length > 0) {
        currentId.value = projects.value[0].id
      }
    } catch {
      if (projects.value.length > 0) {
        currentId.value = projects.value[0].id
      }
    }
  }

  async function switchProject(id: number) {
    currentId.value = id
    try {
      await SetSetting('app.last_project_id', String(id))
    } catch { /* ignore */ }
  }

  async function createProject(name: string, description: string): Promise<Project> {
    const p = await CreateProject(name, description)
    projects.value.push(p)
    return p
  }

  const currentProject = computed(() =>
    projects.value.find(p => p.id === currentId.value) || null
  )

  return { projects, currentId, currentProject, loadProjects, loadLastProject, switchProject, createProject }
})
