import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { HistoryItem } from '../types/history'

export const useHistoryStore = defineStore('history', () => {
  const items = ref<HistoryItem[]>([])
  const searchKeyword = ref('')

  return { items, searchKeyword }
})
