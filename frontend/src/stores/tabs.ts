import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface KeyValuePair {
  key: string
  value: string
  enabled: boolean
}

export interface HttpTabData {
  method: string
  url: string
  params: KeyValuePair[]
  headers: KeyValuePair[]
  body: string
  bodyType: string
  bodyFiles: BodyFileItem[]
  authType: string
  authData: Record<string, string>
  timeoutMs: number
  followRedirect: boolean
}

export interface BodyFileItem {
  key: string
  value: string
  file_path: string
  enabled: boolean
}

export interface WsMessage {
  type: 'sent' | 'received' | 'system'
  content: string
  time: number
}

export interface WsTabData {
  url: string
  messages: WsMessage[]
  connected: boolean
}

export type TabType = 'http' | 'websocket'

export interface Tab {
  id: string
  title: string
  type: TabType
  requestId: string | null
  httpData?: HttpTabData
  wsData?: WsTabData
}

export const useTabsStore = defineStore('tabs', () => {
  const tabs = ref<Tab[]>([])
  const activeTabId = ref<string | null>(null)

  const activeTab = computed<Tab | undefined>(() => tabs.value.find(t => t.id === activeTabId.value))
  const activeTabData = computed<HttpTabData | undefined>(() => activeTab.value?.httpData)

  function createHttpData(): HttpTabData {
    return {
      method: 'GET', url: '', params: [], body: '', bodyType: 'none', bodyFiles: [],
      authType: 'none', authData: {},
      timeoutMs: 30000, followRedirect: true,
      headers: [{ key: 'Content-Type', value: 'application/json', enabled: true }],
    }
  }

  function createWsData(): WsTabData {
    return { url: '', messages: [], connected: false }
  }

  function addHttpTab(requestId?: string, title?: string): string {
    const id = crypto.randomUUID()
    tabs.value.push({ id, title: title || 'New Request', type: 'http', requestId: requestId || null, httpData: createHttpData() })
    activeTabId.value = id
    return id
  }

  function addWsTab(): string {
    const id = crypto.randomUUID()
    tabs.value.push({ id, title: 'WebSocket', type: 'websocket', requestId: null, wsData: createWsData() })
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

  function updateHttpData(data: Partial<HttpTabData>) {
    const tab = tabs.value.find(t => t.id === activeTabId.value)
    if (tab && tab.httpData) Object.assign(tab.httpData, data)
  }

  function updateTabData(data: Partial<HttpTabData>) {
    updateHttpData(data)
  }

  function updateWsData(data: Partial<WsTabData>) {
    const tab = tabs.value.find(t => t.id === activeTabId.value)
    if (tab && tab.wsData) Object.assign(tab.wsData, data)
  }

  function updateTabTitle(tabId: string, title: string) {
    const tab = tabs.value.find(t => t.id === tabId)
    if (tab) tab.title = title
  }

  return {
    tabs, activeTabId, activeTab, activeTabData,
    addHttpTab, addWsTab, removeTab, setActiveTab,
    updateHttpData, updateTabData, updateWsData, updateTabTitle,
  }
})
