import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { TreeItem } from '../types/collection'
import {
  GetCollectionTree, CreateCollection, CreateRequest,
  DeleteCollection, RenameCollection, DeleteRequest, CloneRequest,
} from '../../wailsjs/go/main/App'

export const useCollectionStore = defineStore('collection', () => {
  const tree = ref<TreeItem[]>([])
  const expandedKeys = ref<number[]>([])

  async function loadTree(projectId: number | null) {
    if (!projectId) {
      tree.value = []
      return tree.value
    }
    try {
      tree.value = (await GetCollectionTree(projectId) || []) as any as TreeItem[]
    } catch {
      tree.value = []
    }
    return tree.value
  }

  async function createCollection(projectId: number, parentId: number | null, name: string) {
    await CreateCollection(projectId, parentId as any, name)
  }

  async function createRequest(collectionId: number, name: string, method: string) {
    await CreateRequest(collectionId, name, method)
  }

  async function renameCollection(id: number, name: string) {
    await RenameCollection(id, name)
  }

  async function removeCollection(id: number) {
    await DeleteCollection(id)
  }

  async function removeRequest(id: number) {
    await DeleteRequest(id)
  }

  async function duplicateRequest(id: number) {
    await CloneRequest(id)
  }

  return {
    tree, expandedKeys,
    loadTree, createCollection, createRequest,
    renameCollection, removeCollection, removeRequest, duplicateRequest,
  }
})
