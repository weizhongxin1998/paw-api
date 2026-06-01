import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface HistoryItem {
  id: string
  project_id: string
  request_id: string | null
  method: string
  url: string
  headers: string
  body: string
  response_status: number
  response_body: string
  response_headers: string
  duration_ms: number
  created_at: string
}

export const useHistoryStore = defineStore('history', () => {
  const history = ref<HistoryItem[]>([])

  function setHistory(items: HistoryItem[]) {
    history.value = items
  }

  return { history, setHistory }
})
