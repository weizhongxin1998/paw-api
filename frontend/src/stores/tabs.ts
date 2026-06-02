import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface Tab {
  id: string
  requestId: number
  method: string
  name: string
  url: string
  isDirty: boolean
  headers: string
  params: string
  bodyType: string
  bodyData: string
  authData: string
  collectionId: number
}

export const useTabsStore = defineStore('tabs', () => {
  const tabs = ref<Tab[]>([])
  const activeTabId = ref<string | null>(null)

  return { tabs, activeTabId }
})
