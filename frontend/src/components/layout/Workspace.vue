<template>
  <div class="workspace">
    <div v-if="!activeTab" class="empty-state">
      <div class="empty-icon">📦</div>
      <h2>Paw API</h2>
      <p>点击左侧集合中的请求开始调试</p>
    </div>

    <div v-else class="workspace-editor">
      <div class="tabs-bar">
        <div
          v-for="(tab, idx) in tabs"
          :key="tab.id"
          class="tab"
          :class="{ active: tab.id === activeTabId, dirty: tab.isDirty }"
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
      :on-clickoutside="() => { ctxMenuShow = false }"
      @select="onCtxMenuSelect"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { NDropdown, useDialog } from 'naive-ui'
import UrlBar from '../request/UrlBar.vue'
import RequestPanel from '../request/RequestPanel.vue'
import ResponsePanel from '../response/ResponsePanel.vue'
import type { HttpResponse } from '../../types/response'
import { SendRequest, UpdateRequest } from '../../../wailsjs/go/main/App'
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
const dialog = useDialog()

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
  if (activeTab.value) activeTab.value.isDirty = true
}

function onMethodChange(v: string) { currentMethod.value = v; markDirty() }
function onUrlChange(v: string) { currentUrl.value = v; markDirty() }
function onHeadersChange(v: string) { currentHeaders.value = v; markDirty() }
function onParamsChange(v: string) { currentParams.value = v; markDirty() }
function onBodyTypeChange(v: string) { currentBodyType.value = v; markDirty() }
function onBodyDataChange(v: string) { currentBodyData.value = v; markDirty() }
function onAuthDataChange(v: string) { currentAuthData.value = v; markDirty() }

function syncToActiveTab() {
  if (!activeTab.value) return
  activeTab.value.method = currentMethod.value
  activeTab.value.url = currentUrl.value
  activeTab.value.headers = currentHeaders.value
  activeTab.value.params = currentParams.value
  activeTab.value.bodyType = currentBodyType.value
  activeTab.value.bodyData = currentBodyData.value
  activeTab.value.authData = currentAuthData.value
}

async function onSend() {
  if (isSending.value) return
  isSending.value = true
  syncToActiveTab()

  const envId = envStore.activeEnvId ?? 0
  const reqModel = {
    id: activeTab.value?.requestId || 0,
    collection_id: activeTab.value?.collectionId || 0,
    name: activeTab.value?.name || '',
    description: '',
    method: currentMethod.value,
    url: currentUrl.value,
    headers: currentHeaders.value,
    params: currentParams.value,
    body_type: currentBodyType.value,
    body: currentBodyData.value,
    auth: currentAuthData.value,
    sort_order: 0,
  }

  try {
    const sid = ++sessionCounter
    const resp = await SendRequest(sid, reqModel as any, envId)
    response.value = resp as any
  } catch (e: any) {
    response.value = {
      status: 0,
      time: 0,
      size: 0,
      headers: {},
      cookies: [],
      body: 'Error: ' + (e?.message || e),
      rawRequest: '',
      curlCommand: '',
    }
  } finally {
    isSending.value = false
  }
}

async function onSave() {
  if (!activeTab.value) return
  syncToActiveTab()
  const reqModel = {
    id: activeTab.value.requestId || 0,
    collection_id: activeTab.value.collectionId || 0,
    name: activeTab.value.name,
    description: '',
    method: currentMethod.value,
    url: currentUrl.value,
    headers: currentHeaders.value,
    params: currentParams.value,
    body_type: currentBodyType.value,
    body: currentBodyData.value,
    auth: currentAuthData.value,
    sort_order: 0,
  }
  try {
    await UpdateRequest(reqModel as any)
    if (activeTab.value) activeTab.value.isDirty = false
  } catch (e: any) {
    dialog.warning({
      title: '保存失败',
      content: e?.message || String(e),
      positiveText: '确定',
    })
  }
}

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
  }
}

function onCloseTab(tab: Tab) {
  if (tab.isDirty) {
    dialog.warning({
      title: '确认关闭',
      content: '未保存的修改将丢失，确定关闭？',
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: () => closeTab(tab.id),
    })
  } else {
    closeTab(tab.id)
  }
}

function closeTab(id: string) {
  const idx = tabs.value.findIndex(t => t.id === id)
  tabs.value = tabs.value.filter(t => t.id !== id)
  if (activeTabId.value === id) {
    if (tabs.value.length > 0) {
      selectTab(tabs.value[Math.min(idx, tabs.value.length - 1)].id)
    } else {
      activeTabId.value = null
      activeTab.value = null
    }
  }
}

function closeOthers(tabId: string) {
  const dirtyOthers = tabs.value.some(t => t.id !== tabId && t.isDirty)
  const doClose = () => {
    tabs.value = tabs.value.filter(t => t.id === tabId)
    selectTab(tabId)
  }
  if (dirtyOthers) {
    dialog.warning({
      title: '确认关闭',
      content: '未保存的修改将丢失，确定关闭？',
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: doClose,
    })
  } else {
    doClose()
  }
}

function closeRight(idx: number) {
  const dirtyRight = tabs.value.slice(idx + 1).some(t => t.isDirty)
  const doClose = () => {
    tabs.value = tabs.value.slice(0, idx + 1)
    if (activeTabId.value && !tabs.value.find(t => t.id === activeTabId.value))
      selectTab(tabs.value[tabs.value.length - 1].id)
  }
  if (dirtyRight) {
    dialog.warning({
      title: '确认关闭',
      content: '未保存的修改将丢失，确定关闭？',
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: doClose,
    })
  } else {
    doClose()
  }
}

function closeAll() {
  const dirty = tabs.value.some(t => t.isDirty)
  const doClose = () => {
    tabs.value = []
    activeTabId.value = null
    activeTab.value = null
  }
  if (dirty) {
    dialog.warning({
      title: '确认关闭',
      content: '未保存的修改将丢失，确定关闭？',
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: doClose,
    })
  } else {
    doClose()
  }
}

// Context menu
const ctxMenuShow = ref(false)
const ctxMenuX = ref(0)
const ctxMenuY = ref(0)
let ctxTarget: { tab: Tab; idx: number } | null = null

const ctxMenuOptions = computed(() => [
  { label: '关闭', key: 'close' },
  { label: '关闭其他', key: 'close-others' },
  { label: '关闭右侧', key: 'close-right' },
  { label: '关闭全部', key: 'close-all' },
])

function onTabContextMenu(e: MouseEvent, tab: Tab, idx: number) {
  ctxTarget = { tab, idx }
  ctxMenuX.value = e.clientX
  ctxMenuY.value = e.clientY
  ctxMenuShow.value = true
}

function onCtxMenuSelect(key: string) {
  ctxMenuShow.value = false
  if (!ctxTarget) return
  switch (key) {
    case 'close': onCloseTab(ctxTarget.tab); break
    case 'close-others': closeOthers(ctxTarget.tab.id); break
    case 'close-right': closeRight(ctxTarget.idx); break
    case 'close-all': closeAll(); break
  }
}

// Drag
const dragIdx = ref<number | null>(null)
const dragOverId = ref<string | null>(null)

function onDragStart(e: DragEvent, idx: number) {
  dragIdx.value = idx
  e.dataTransfer!.effectAllowed = 'move'
}
function onDragOver(_e: DragEvent, tabId: string) { dragOverId.value = tabId }
function onDragLeave(id: string) { if (dragOverId.value === id) dragOverId.value = null }
function onDrop(_e: DragEvent, targetIdx: number) {
  dragOverId.value = null
  if (dragIdx.value == null || dragIdx.value === targetIdx) return
  const arr = [...tabs.value]
  const [moved] = arr.splice(dragIdx.value, 1)
  arr.splice(targetIdx, 0, moved)
  tabs.value = arr
}
function onDragEnd() { dragIdx.value = null; dragOverId.value = null }

// Exposed
function openTab(tab: Tab) {
  const existing = tabs.value.find(t => t.id === tab.id)
  if (existing) { selectTab(existing.id); return }
  tabs.value.push(tab)
  selectTab(tab.id)
}

// Keyboard shortcuts
function onKeydown(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.key === 's') {
    e.preventDefault()
    onSave()
  }
}

onMounted(() => document.addEventListener('keydown', onKeydown))
onUnmounted(() => document.removeEventListener('keydown', onKeydown))

defineExpose({ openTab, tabs, activeTabId, selectTab })
</script>

<style scoped>
.workspace {
  flex: 1;
  display: flex;
  background: #fff;
  min-width: 0;
  overflow: hidden;
}
.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #aaa;
}
.empty-icon { font-size: 48px; margin-bottom: 12px; opacity: 0.3; }
.empty-state h2 { font-size: 18px; color: #555; margin: 0 0 6px; font-weight: 500; }
.empty-state p { font-size: 13px; color: #999; margin: 0; }
.workspace-editor {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
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
.tab.active { background: #fff; border-bottom: 2px solid #18a058; }
.tab:hover { background: #d9d9d9; }
.tab.active:hover { background: #fff; }
.tab-method {
  font-size: 9px; font-weight: 700; padding: 0 4px; border-radius: 2px;
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
