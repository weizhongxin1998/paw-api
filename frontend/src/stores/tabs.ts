import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useTabsStore = defineStore('tabs', () => {
  const tabs = ref<any[]>([])
  const activeTabId = ref<string | null>(null)

  return { tabs, activeTabId }
})
