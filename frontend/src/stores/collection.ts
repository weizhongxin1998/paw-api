import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { TreeItem } from '../types/collection'
import { GetCollectionTree, CreateCollection, CreateRequest } from '../../wailsjs/go/main/App'

export const useCollectionStore = defineStore('collection', () => {
  const tree = ref<TreeItem[]>([])
  const expandedKeys = ref<number[]>([])

  async function loadTree(projectId: number | null) {
    if (!projectId) {
      tree.value = []
      return
    }
    try {
      const raw = await GetCollectionTree(projectId)
      tree.value = (raw || []) as any as TreeItem[]
    } catch {
      tree.value = []
    }
  }

  async function createCollection(projectId: number, parentId: number | null, name: string) {
    await CreateCollection(projectId, parentId as any, name)
    await loadTree(projectId)
  }

  async function createRequestInCollection(projectId: number, collectionId: number, name: string, method: string) {
    await CreateRequest(collectionId, name, method)
    await loadTree(projectId)
  }

  return { tree, expandedKeys, loadTree, createCollection, createRequestInCollection }
})
