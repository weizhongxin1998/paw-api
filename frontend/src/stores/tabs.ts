import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface KeyValuePair {
  key: string
  value: string
  enabled: boolean
}

export interface TabData {
  method: string
  url: string
  params: KeyValuePair[]
  headers: KeyValuePair[]
  body: string
  bodyType: string
  auth: KeyValuePair[]
}

export interface Tab {
  id: string
  title: string
  requestId: string | null
  data: TabData
}

export const useTabsStore = defineStore('tabs', () => {
  const tabs = ref<Tab[]>([])
  const activeTabId = ref<string | null>(null)

  const activeTab = computed<Tab | undefined>(() => tabs.value.find(t => t.id === activeTabId.value))
  const activeTabData = computed<TabData | undefined>(() => activeTab.value?.data)

  function createEmptyData(): TabData {
    return {
      method: 'GET',
      url: '',
      params: [],
      headers: [{ key: 'Content-Type', value: 'application/json', enabled: true }],
      body: '',
      bodyType: 'none',
      auth: [],
    }
  }

  function addTab(requestId?: string, title?: string): string {
    const id = crypto.randomUUID()
    const tab: Tab = {
      id,
      title: title || 'New Request',
      requestId: requestId || null,
      data: createEmptyData(),
    }
    tabs.value.push(tab)
    activeTabId.value = id
    return id
  }

  function removeTab(tabId: string) {
    const idx = tabs.value.findIndex(t => t.id === tabId)
    if (idx === -1) return
    tabs.value.splice(idx, 1)
    if (activeTabId.value === tabId) {
      activeTabId.value = tabs.value.length > 0 ? tabs.value[Math.min(idx, tabs.value.length - 1)].id : null
    }
  }

  function setActiveTab(tabId: string) {
    activeTabId.value = tabId
  }

  function updateTabData(data: Partial<TabData>) {
    const tab = tabs.value.find(t => t.id === activeTabId.value)
    if (tab) {
      Object.assign(tab.data, data)
    }
  }

  function updateTabTitle(tabId: string, title: string) {
    const tab = tabs.value.find(t => t.id === tabId)
    if (tab) tab.title = title
  }

  return {
    tabs, activeTabId, activeTab, activeTabData,
    addTab, removeTab, setActiveTab, updateTabData, updateTabTitle,
  }
})
