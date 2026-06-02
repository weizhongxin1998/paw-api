import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { TreeItem } from '../types/collection'

export const useCollectionStore = defineStore('collection', () => {
  const tree = ref<TreeItem[]>([])
  const expandedKeys = ref<number[]>([])

  return { tree, expandedKeys }
})
