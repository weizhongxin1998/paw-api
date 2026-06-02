<template>
  <div class="workspace">
    <div v-if="!activeTab && !historyItem" class="empty-state">
      <div class="empty-icon">📦</div>
      <h2>Paw API</h2>
      <p>点击左侧集合中的请求开始调试</p>
    </div>

    <div v-if="historyItem" class="workspace-editor">
      <div class="tabs-bar">
        <span class="tabs-msg">选择一条历史记录查看详情，或双击在新 Tab 中回放</span>
      </div>
      <div class="history-detail">
        <div class="hist-detail-row"><span class="hist-detail-label">方法</span><span class="method-badge" :class="historyItem.method?.toLowerCase()">{{ historyItem.method }}</span></div>
        <div class="hist-detail-row"><span class="hist-detail-label">URL</span><span class="hist-detail-value">{{ historyItem.url }}</span></div>
        <div class="hist-detail-row"><span class="hist-detail-label">状态码</span><span class="hist-detail-value" :class="statusClass(historyItem.response_status)">{{ historyItem.response_status }}</span></div>
        <div class="hist-detail-row"><span class="hist-detail-label">耗时</span><span class="hist-detail-value">{{ historyItem.duration_ms }}ms</span></div>
        <div v-if="historyItem.request_headers" class="hist-detail-section">
          <h4>请求头</h4>
          <pre class="hist-detail-pre">{{ formatJson(historyItem.request_headers) }}</pre>
        </div>
        <div v-if="historyItem.request_body" class="hist-detail-section">
          <h4>请求体</h4>
          <pre class="hist-detail-pre">{{ formatJson(historyItem.request_body) }}</pre>
        </div>
        <div v-if="historyItem.response_headers" class="hist-detail-section">
          <h4>响应头</h4>
          <pre class="hist-detail-pre">{{ formatJson(historyItem.response_headers) }}</pre>
        </div>
        <div v-if="historyItem.response_body" class="hist-detail-section">
          <h4>响应体</h4>
          <pre class="hist-detail-pre">{{ formatJson(historyItem.response_body) }}</pre>
        </div>
      </div>
    </div>

    <div v-else-if="activeTab" class="workspace-editor">
      <div class="tabs-bar">
        <div
          v-for="(tab, idx) in tabs"
          :key="tab.id"
          class="tab"
          :class="{ active: tab.id === activeTabId, dirty: tab.isDirty, preview: tab.isPreview }"
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
        <span class="tab-plus" @click="addNewTab">+</span>
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
        :url="currentUrl"
        :path-vars="currentPathVars"
        @update:headers="onHeadersChange"
        @update:params="onParamsChange"
        @update:body-type="onBodyTypeChange"
        @update:body-data="onBodyDataChange"
        @update:auth-data="onAuthDataChange"
        @update:path-vars="onPathVarsChange"
      />

      <div
        v-if="response"
        class="resize-handle"
        @mousedown="onResizeStart"
      ></div>

      <ResponsePanel
        v-if="response"
        :response="response"
        :style="{ height: responseHeight + 'px' }"
      />
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
  isPreview: boolean
  pathVars: string
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
const historyItem = ref<any>(null)

const currentMethod = ref('GET')
const currentUrl = ref('')
const currentPathVars = ref('[]')
const currentHeaders = ref('[]')
const currentParams = ref('[]')
const currentBodyType = ref('none')
const currentBodyData = ref('{}')
const currentAuthData = ref('{"type":"none"}')

let sessionCounter = 0
let tabIdCounter = 0

function markDirty() {
  if (!activeTab.value) return
  activeTab.value.isDirty = true
  if (activeTab.value.isPreview) {
    activeTab.value.isPreview = false
  }
}

function onMethodChange(v: string) { currentMethod.value = v; markDirty() }
function onUrlChange(v: string) { currentUrl.value = v; markDirty() }
function onPathVarsChange(v: string) { currentPathVars.value = v }
function onHeadersChange(v: string) { currentHeaders.value = v; markDirty() }
function onParamsChange(v: string) { currentParams.value = v; markDirty() }
function onBodyTypeChange(v: string) { currentBodyType.value = v; markDirty() }
function onBodyDataChange(v: string) { currentBodyData.value = v; markDirty() }
function onAuthDataChange(v: string) { currentAuthData.value = v; markDirty() }

function syncToActiveTab() {
  if (!activeTab.value) return
  activeTab.value.method = currentMethod.value
  activeTab.value.url = currentUrl.value
  activeTab.value.pathVars = currentPathVars.value
  activeTab.value.headers = currentHeaders.value
  activeTab.value.params = currentParams.value
  activeTab.value.bodyType = currentBodyType.value
  activeTab.value.bodyData = currentBodyData.value
  activeTab.value.authData = currentAuthData.value
}

function addNewTab() {
  const id = 'tab-' + (++tabIdCounter) + '-' + Date.now()
  const tab: Tab = {
    id,
    requestId: 0,
    method: 'GET',
    name: '新建请求',
    url: '',
    pathVars: '[]',
    isDirty: true,
    isPreview: false,
    headers: '[]',
    params: '[]',
    bodyType: 'none',
    bodyData: '{}',
    authData: '{"type":"none"}',
    collectionId: 0,
  }
  tabs.value.push(tab)
  selectTab(tab.id)
}

function showHistoryDetail(item: any) {
  historyItem.value = item
}

function clearHistoryDetail() {
  historyItem.value = null
}

function statusClass(code: number): string {
  if (code < 300) return 'status-2xx'
  if (code < 400) return 'status-3xx'
  if (code < 500) return 'status-4xx'
  return 'status-5xx'
}

function formatJson(raw: string): string {
  try {
    return JSON.stringify(JSON.parse(raw), null, 2)
  } catch {
    return raw
  }
}

async function onSend() {
  if (isSending.value) return
  isSending.value = true
  syncToActiveTab()

  // Resolve path variables in URL
  let resolvedUrl = currentUrl.value
  try {
    const pv = JSON.parse(currentPathVars.value || '[]') as { key: string; value: string }[]
    for (const v of pv) {
      if (v.value) {
        resolvedUrl = resolvedUrl.replace(':' + v.key, v.value)
        resolvedUrl = resolvedUrl.replace('{' + v.key + '}', v.value)
      }
    }
  } catch {}

  const envId = envStore.activeEnvId ?? 0
  const reqModel = {
    id: activeTab.value?.requestId || 0,
    collection_id: activeTab.value?.collectionId || 0,
    name: activeTab.value?.name || '',
    description: '',
    method: currentMethod.value,
    url: resolvedUrl,
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
    if (responseHeight.value < 150) responseHeight.value = 200
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
    currentPathVars.value = activeTab.value.pathVars || '[]'
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

function openTab(tab: Tab) {
  historyItem.value = null
  const existing = tabs.value.find(t => t.requestId > 0 && t.requestId === tab.requestId)
  if (existing) {
    if (existing.isPreview) existing.isPreview = false
    selectTab(existing.id)
    return
  }
  tab.isPreview = false
  tabs.value.push(tab)
  selectTab(tab.id)
}

function previewTab(tab: Tab) {
  historyItem.value = null
  const existing = tabs.value.find(t => t.requestId > 0 && t.requestId === tab.requestId)
  if (existing) { selectTab(existing.id); return }
  const oldPreview = tabs.value.findIndex(t => t.isPreview)
  if (oldPreview !== -1) {
    tabs.value.splice(oldPreview, 1)
  }
  tab.isPreview = true
  tabs.value.push(tab)
  selectTab(tab.id)
}

const responseHeight = ref(200)
let resizeStartY = 0
let resizeStartH = 0

function onResizeStart(e: MouseEvent) {
  e.preventDefault()
  resizeStartY = e.clientY
  resizeStartH = responseHeight.value
  document.addEventListener('mousemove', onResizeMove)
  document.addEventListener('mouseup', onResizeEnd)
}

function onResizeMove(e: MouseEvent) {
  const delta = resizeStartY - e.clientY
  responseHeight.value = Math.max(80, Math.min(800, resizeStartH + delta))
}

function onResizeEnd() {
  document.removeEventListener('mousemove', onResizeMove)
  document.removeEventListener('mouseup', onResizeEnd)
}

function onKeydown(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.key === 's') {
    e.preventDefault()
    onSave()
  }
}

onMounted(() => document.addEventListener('keydown', onKeydown))
onUnmounted(() => document.removeEventListener('keydown', onKeydown))

defineExpose({ openTab, previewTab, tabs, activeTabId, selectTab, showHistoryDetail, clearHistoryDetail })
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
.tabs-msg {
  font-size: 12px;
  color: #aaa;
  padding: 5px 12px;
}
.tab {
  padding: 5px 12px;
  font-size: 13px;
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
.tab.preview { font-style: italic; color: #888; }
.tab:hover { background: #d9d9d9; }
.tab.active:hover { background: #fff; }
.tab-method {
  font-size: 10px; font-weight: 700; padding: 0 4px; border-radius: 2px;
}
.tab-method.get { background: #d4edda; color: #155724; }
.tab-method.post { background: #fff3cd; color: #856404; }
.tab-method.put { background: #d0e8ff; color: #004085; }
.tab-method.delete { background: #f8d7da; color: #721c24; }
.tab-method.patch { background: #f3e5f5; color: #6a1b9a; }
.tab-name { font-size: 13px; }
.tab-dirty { width: 6px; height: 6px; background: #bbb; border-radius: 50%; }
.tab-close { color: #aaa; font-size: 13px; margin-left: 4px; }
.tab-close:hover { color: #d03050; }
.tab-plus {
  padding: 5px 8px;
  font-size: 16px;
  cursor: pointer;
  color: #18a058;
  user-select: none;
}
.resize-handle {
  height: 5px;
  background: #e0e0e0;
  cursor: ns-resize;
  flex-shrink: 0;
  border-top: 1px solid #d0d0d0;
  border-bottom: 1px solid #d0d0d0;
}
.resize-handle:hover {
  background: #18a058;
}

.history-detail {
  flex: 1;
  overflow-y: auto;
  padding: 12px 16px;
  font-size: 13px;
}
.hist-detail-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 6px 0;
  border-bottom: 1px solid #f0f0f0;
}
.hist-detail-label {
  color: #999;
  width: 80px;
  flex-shrink: 0;
  font-size: 12px;
}
.hist-detail-value {
  word-break: break-all;
  color: #333;
}
.hist-detail-value.status-2xx { color: #18a058; font-weight: 600; }
.hist-detail-value.status-3xx { color: #0288d1; font-weight: 600; }
.hist-detail-value.status-4xx { color: #f0a020; font-weight: 600; }
.hist-detail-value.status-5xx { color: #d03050; font-weight: 600; }
.method-badge {
  font-size: 10px;
  font-weight: 700;
  padding: 1px 6px;
  border-radius: 2px;
}
.method-badge.get { background: #d4edda; color: #155724; }
.method-badge.post { background: #fff3cd; color: #856404; }
.method-badge.put { background: #d0e8ff; color: #004085; }
.method-badge.delete { background: #f8d7da; color: #721c24; }
.hist-detail-section {
  margin-top: 12px;
}
.hist-detail-section h4 {
  margin: 0 0 6px;
  font-size: 12px;
  color: #666;
  font-weight: 600;
}
.hist-detail-pre {
  background: #f8f9fa;
  border: 1px solid #e8e8e8;
  border-radius: 4px;
  padding: 10px;
  font-size: 11px;
  font-family: 'SF Mono', Consolas, monospace;
  white-space: pre-wrap;
  word-break: break-all;
  max-height: 200px;
  overflow-y: auto;
  margin: 0;
}
</style>
