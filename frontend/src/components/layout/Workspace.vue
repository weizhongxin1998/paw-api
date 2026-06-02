<template>
  <div class="workspace" @keydown="onGlobalKeydown" tabindex="0">
    <n-empty v-if="!activeTab" description="点击左侧集合中的请求开始调试" class="empty-state" />
    <div v-else class="workspace-editor">
      <div class="tabs-bar">
        <div v-if="tabs.length === 0" class="no-tabs">无打开 Tab</div>
        <div
          v-for="(tab, idx) in tabs"
          :key="tab.id"
          class="tab"
          :class="{ active: tab.id === activeTabId, 'drag-over': dragOverId === tab.id, dirty: tab.isDirty }"
          draggable="true"
          @click="selectTab(tab.id)"
          @contextmenu.prevent="onTabContextMenu($event, tab, idx)"
          @dragstart="onDragStart($event, idx)"
          @dragover.prevent="onDragOver($event, tab.id)"
          @dragleave="onDragLeave(tab.id)"
          @drop="onDrop($event, idx)"
          @dragend="onDragEnd"
        >
          <span class="tab-method" :class="tab.method?.toLowerCase()">{{ tab.method }}</span>
          <span class="tab-name">{{ tab.name }}</span>
          <span v-if="tab.isDirty" class="tab-dirty"></span>
          <span class="tab-close" @click.stop="onCloseTab(tab)">x</span>
        </div>
      </div>

      <UrlBar
        :model-method="currentMethod"
        :model-url="currentUrl"
        @update:model-method="onMethodChange"
        @update:model-url="onUrlChange"
        @send="onSend"
      />

      <RequestPanel
        :headers="currentHeaders"
        :params="currentParams"
        :body-type="currentBodyType"
        :body-data="currentBodyData"
        :auth-data="currentAuthData"
        @update:headers="onHeadersChange"
        @update:params="onParamsChange"
        @update:body-type="onBodyTypeChange"
        @update:body-data="onBodyDataChange"
        @update:auth-data="onAuthDataChange"
      />

      <ResponsePanel :response="response" />
    </div>

    <n-dropdown
      placement="bottom-start"
      trigger="manual"
      :x="ctxMenuX"
      :y="ctxMenuY"
      :options="ctxMenuOptions"
      :show="ctxMenuShow"
      :on-clickoutside="onCtxMenuClose"
      @select="onCtxMenuSelect"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { NEmpty, NDropdown, NDialog, useMessage } from 'naive-ui'
import UrlBar from '../request/UrlBar.vue'
import RequestPanel from '../request/RequestPanel.vue'
import ResponsePanel from '../response/ResponsePanel.vue'
import type { HttpResponse } from '../../types/response'
import type { DropdownOption } from 'naive-ui'
import { SendRequest, UpdateRequest } from '../../../wailsjs/go/main/App'
import { models } from '../../../wailsjs/go/models'
import { useEnvStore } from '../../stores/env'

interface Tab {
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

const props = defineProps<{
  projectId: number | null
}>()

const envStore = useEnvStore()
const message = useMessage()

const activeTab = ref<Tab | null>(null)
const activeTabId = ref<string | null>(null)
const tabs = ref<Tab[]>([])
const response = ref<HttpResponse | null>(null)
const isSending = ref(false)

const currentMethod = ref('GET')
const currentUrl = ref('')
const currentHeaders = ref('[]')
const currentParams = ref('[]')
const currentBodyType = ref('none')
const currentBodyData = ref('{}')
const currentAuthData = ref('{"type":"none"}')

let sessionCounter = 0

function markDirty() {
  if (activeTab.value) {
    activeTab.value.isDirty = true
  }
}

function onMethodChange(v: string) {
  currentMethod.value = v
  if (activeTab.value) { activeTab.value.method = v; markDirty() }
}

function onUrlChange(v: string) {
  currentUrl.value = v
  if (activeTab.value) { activeTab.value.url = v; markDirty() }
}

function onHeadersChange(v: string) {
  currentHeaders.value = v
  if (activeTab.value) { activeTab.value.headers = v; markDirty() }
}

function onParamsChange(v: string) {
  currentParams.value = v
  if (activeTab.value) { activeTab.value.params = v; markDirty() }
}

function onBodyTypeChange(v: string) {
  currentBodyType.value = v
  if (activeTab.value) { activeTab.value.bodyType = v; markDirty() }
}

function onBodyDataChange(v: string) {
  currentBodyData.value = v
  if (activeTab.value) { activeTab.value.bodyData = v; markDirty() }
}

function onAuthDataChange(v: string) {
  currentAuthData.value = v
  if (activeTab.value) { activeTab.value.authData = v; markDirty() }
}

function buildRequest(): models.Request {
  const t = activeTab.value
  return models.Request.createFrom({
    id: t?.requestId || 0,
    collection_id: t?.collectionId || 0,
    name: t?.name || '',
    description: '',
    method: currentMethod.value,
    url: currentUrl.value,
    headers: currentHeaders.value,
    params: currentParams.value,
    body_type: currentBodyType.value,
    body: currentBodyData.value,
    auth: currentAuthData.value,
    sort_order: 0,
  })
}

async function onSend() {
  if (!activeTab.value) return
  isSending.value = true
  const sessionId = ++sessionCounter
  const req = buildRequest()
  const envId = envStore.activeEnvId ?? 0
  try {
    const resp = await SendRequest(sessionId, req, envId)
    response.value = {
      status: resp.status,
      time: resp.time,
      size: resp.size,
      headers: resp.headers || {},
      cookies: (resp.cookies || []).map((c: any) => ({
        name: c.name || '',
        value: c.value || '',
        domain: c.domain || '',
        path: c.path || '',
      })),
      body: resp.body || '',
      rawRequest: resp.raw_request || '',
      curlCommand: resp.curl_command || '',
    }
  } catch (e: any) {
    response.value = {
      status: 0,
      time: 0,
      size: 0,
      headers: {},
      cookies: [],
      body: e?.message || String(e),
      rawRequest: '',
      curlCommand: '',
    }
  } finally {
    isSending.value = false
  }
}

function onGlobalKeydown(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.key === 's') {
    e.preventDefault()
    onSaveRequest()
  }
}

async function onSaveRequest() {
  if (!activeTab.value || !activeTab.value.requestId) {
    message.warning('无法保存：请求数据不完整')
    return
  }
  const req = buildRequest()
  try {
    await UpdateRequest(req)
    if (activeTab.value) {
      activeTab.value.isDirty = false
    }
    message.success('已保存')
  } catch (e: any) {
    message.error('保存失败: ' + (e?.message || String(e)))
  }
}

// ---- Tab management ----

function selectTab(id: string) {
  activeTabId.value = id
  activeTab.value = tabs.value.find(t => t.id === id) || null
  if (activeTab.value) {
    currentMethod.value = activeTab.value.method
    currentUrl.value = activeTab.value.url
    currentHeaders.value = activeTab.value.headers
    currentParams.value = activeTab.value.params
    currentBodyType.value = activeTab.value.bodyType
    currentBodyData.value = activeTab.value.bodyData
    currentAuthData.value = activeTab.value.authData
    response.value = null
  }
}

function onCloseTab(tab: Tab) {
  if (tab.isDirty) {
    showCloseConfirm(() => closeTab(tab.id))
  } else {
    closeTab(tab.id)
  }
}

function closeTab(id: string) {
  const idx = tabs.value.findIndex(t => t.id === id)
  tabs.value = tabs.value.filter(t => t.id !== id)
  if (activeTabId.value === id) {
    if (tabs.value.length > 0) {
      const nextIdx = Math.min(idx, tabs.value.length - 1)
      selectTab(tabs.value[nextIdx].id)
    } else {
      activeTabId.value = null
      activeTab.value = null
    }
  }
}

function closeOthers(tabId: string) {
  const dirtyOthers = tabs.value.some(t => t.id !== tabId && t.isDirty)
  if (dirtyOthers) {
    showCloseConfirm(() => {
      tabs.value = tabs.value.filter(t => t.id === tabId)
      selectTab(tabId)
    })
  } else {
    tabs.value = tabs.value.filter(t => t.id === tabId)
    selectTab(tabId)
  }
}

function closeRight(idx: number) {
  const dirtyRight = tabs.value.slice(idx + 1).some(t => t.isDirty)
  if (dirtyRight) {
    showCloseConfirm(() => {
      tabs.value = tabs.value.slice(0, idx + 1)
      if (activeTabId.value && !tabs.value.find(t => t.id === activeTabId.value)) {
        selectTab(tabs.value[tabs.value.length - 1].id)
      }
    })
  } else {
    tabs.value = tabs.value.slice(0, idx + 1)
    if (activeTabId.value && !tabs.value.find(t => t.id === activeTabId.value)) {
      selectTab(tabs.value[tabs.value.length - 1].id)
    }
  }
}

function closeAll() {
  const dirtyAny = tabs.value.some(t => t.isDirty)
  if (dirtyAny) {
    showCloseConfirm(() => {
      tabs.value = []
      activeTabId.value = null
      activeTab.value = null
    })
  } else {
    tabs.value = []
    activeTabId.value = null
    activeTab.value = null
  }
}

// ---- Close confirmation ----

function showCloseConfirm(onOk: () => void) {
  const dialog = (NDialog as any).warning ?? (NDialog as any).create
  if (typeof dialog !== 'function') {
    onOk()
    return
  }
  dialog({
    title: '确认关闭',
    content: '未保存的修改将丢失，确定关闭？',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: () => { onOk() },
  })
}

// ---- Context menu ----

const ctxMenuShow = ref(false)
const ctxMenuX = ref(0)
const ctxMenuY = ref(0)
const ctxMenuTarget = ref<{ tab: Tab; idx: number } | null>(null)

const ctxMenuOptions = computed<DropdownOption[]>(() => [
  { label: '关闭', key: 'close' },
  { label: '关闭其他', key: 'close-others' },
  { label: '关闭右侧', key: 'close-right' },
  { label: '关闭全部', key: 'close-all' },
])

function onTabContextMenu(e: MouseEvent, tab: Tab, idx: number) {
  ctxMenuTarget.value = { tab, idx }
  ctxMenuX.value = e.clientX
  ctxMenuY.value = e.clientY
  ctxMenuShow.value = true
}

function onCtxMenuClose() {
  ctxMenuShow.value = false
  ctxMenuTarget.value = null
}

function onCtxMenuSelect(key: string) {
  const t = ctxMenuTarget.value
  if (!t) return
  switch (key) {
    case 'close': onCloseTab(t.tab); break
    case 'close-others': closeOthers(t.tab.id); break
    case 'close-right': closeRight(t.idx); break
    case 'close-all': closeAll(); break
  }
  ctxMenuShow.value = false
  ctxMenuTarget.value = null
}

// ---- Drag and drop ----

const dragIdx = ref<number | null>(null)
const dragOverId = ref<string | null>(null)

function onDragStart(e: DragEvent, idx: number) {
  dragIdx.value = idx
  if (e.dataTransfer) {
    e.dataTransfer.effectAllowed = 'move'
    e.dataTransfer.setData('text/plain', String(idx))
  }
}

function onDragOver(_e: DragEvent, tabId: string) {
  dragOverId.value = tabId
}

function onDragLeave(tabId: string) {
  if (dragOverId.value === tabId) {
    dragOverId.value = null
  }
}

function onDrop(_e: DragEvent, targetIdx: number) {
  dragOverId.value = null
  if (dragIdx.value == null || dragIdx.value === targetIdx) {
    dragIdx.value = null
    return
  }
  const arr = [...tabs.value]
  const [moved] = arr.splice(dragIdx.value, 1)
  arr.splice(targetIdx, 0, moved)
  tabs.value = arr
  dragIdx.value = null
}

function onDragEnd() {
  dragIdx.value = null
  dragOverId.value = null
}

// ---- Expose for other components ----

function openTab(tab: Tab) {
  const existing = tabs.value.find(t => t.id === tab.id)
  if (existing) {
    selectTab(existing.id)
    return
  }
  tabs.value.push(tab)
  selectTab(tab.id)
}

defineExpose({ openTab, tabs, activeTabId, selectTab })
</script>

<style scoped>
.workspace {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fff;
  min-width: 0;
  outline: none;
}
.empty-state {
  flex: 1;
}
.workspace-editor {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
  height: 100%;
}
.tabs-bar {
  display: flex;
  background: #f2f2f2;
  border-bottom: 1px solid #ddd;
  height: 32px;
  align-items: flex-end;
  padding: 0 4px;
  gap: 2px;
  overflow-x: auto;
  flex-shrink: 0;
}
.no-tabs {
  padding: 6px 10px;
  color: #aaa;
  font-size: 11px;
}
.tab {
  padding: 5px 12px;
  font-size: 11px;
  background: #e2e2e2;
  border: 1px solid #ddd;
  border-bottom: none;
  border-radius: 5px 5px 0 0;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 5px;
  white-space: nowrap;
  user-select: none;
}
.tab.active {
  background: #fff;
  border-bottom: 2px solid #18a058;
}
.tab.dirty {
  font-style: italic;
}
.tab.drag-over {
  background: #d0e8ff;
  border-color: #18a058;
}
.tab:hover {
  background: #d9d9d9;
}
.tab.active:hover {
  background: #fff;
}
.tab-method {
  font-size: 9px;
  font-weight: 700;
  padding: 0 4px;
  border-radius: 2px;
}
.tab-method.get { background: #d4edda; color: #155724; }
.tab-method.post { background: #fff3cd; color: #856404; }
.tab-method.put { background: #d0e8ff; color: #004085; }
.tab-method.delete { background: #f8d7da; color: #721c24; }
.tab-method.patch { background: #f3e5f5; color: #6a1b9a; }
.tab-name { font-size: 11px; }
.tab-dirty { width: 6px; height: 6px; background: #bbb; border-radius: 50%; }
.tab-close { color: #aaa; font-size: 13px; margin-left: 4px; }
.tab-close:hover { color: #d03050; }
</style>
