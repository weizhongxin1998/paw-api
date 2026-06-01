import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Project, Collection } from '../types/project'

export const useProjectStore = defineStore('project', () => {
  const projects = ref<Project[]>([])
  const currentProject = ref<Project | null>(null)
  const collections = ref<Collection[]>([])
  const selectedCollectionId = ref<string | null>(null)
  const refreshKey = ref(0)

  function setProjects(list: Project[]) {
    projects.value = list
  }

  function addProject(p: Project) {
    projects.value.push(p)
  }

  function setCurrentProject(p: Project | null) {
    currentProject.value = p
  }

  function setCollections(list: Collection[]) {
    collections.value = list
  }

  function addCollection(c: Collection) {
    collections.value.push(c)
  }

  function removeCollection(id: string) {
    collections.value = collections.value.filter(c => c.id !== id)
    if (selectedCollectionId.value === id) selectedCollectionId.value = null
  }

  function updateCollection(c: Collection) {
    const idx = collections.value.findIndex(x => x.id === c.id)
    if (idx !== -1) collections.value[idx] = c
  }

  function selectCollection(id: string | null) {
    selectedCollectionId.value = id
  }

  function removeProject(id: string) {
    projects.value = projects.value.filter(p => p.id !== id)
    if (currentProject.value?.id === id) {
      currentProject.value = projects.value[0] || null
    }
  }

  function updateProject(p: Project) {
    const idx = projects.value.findIndex(x => x.id === p.id)
    if (idx !== -1) projects.value[idx] = p
    if (currentProject.value?.id === p.id) currentProject.value = p
  }

  function triggerRefresh() {
    refreshKey.value++
  }

  return {
    projects, currentProject, collections, selectedCollectionId, refreshKey,
    setProjects, addProject, setCurrentProject, setCollections,
    addCollection, removeCollection, updateCollection, selectCollection,
    removeProject, updateProject, triggerRefresh,
  }
})
