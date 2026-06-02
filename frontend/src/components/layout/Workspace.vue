<template>
  <div class="workspace">
    <div v-if="!activeTab && !historyItem" class="empty-state">
      <div class="empty-logo">
        <svg width="48" height="48" viewBox="0 0 48 48" fill="none" stroke="currentColor" stroke-width="1.5" opacity="0.4">
          <rect x="4" y="8" width="40" height="32" rx="4" />
          <line x1="4" y1="18" x2="44" y2="18" />
          <circle cx="10" cy="13" r="1.5" fill="currentColor" stroke="none" />
          <circle cx="15" cy="13" r="1.5" fill="currentColor" stroke="none" />
          <circle cx="20" cy="13" r="1.5" fill="currentColor" stroke="none" />
        </svg>
      </div>
      <h2>Paw API</h2>
      <p>点击左侧集合中的请求<br/>或按 <kbd>Ctrl + N</kbd> 新建标签</p>
    </div>

    <div v-if="historyItem" class="workspace-editor">
      <div class="tabs-bar">
        <span class="tabs-msg">历史记录详情 · 双击可回放至新标签</span>
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
          <span class="tab-close" @click.stop="onCloseTab(tab)">&times;</span>
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
      >
        <div class="resize-line"></div>
      </div>

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
  if (activeTab.value.isPreview) activeTab.value.isPreview = false
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
    id, requestId: 0, method: 'GET', name: '新建请求', url: '',
    pathVars: '[]', isDirty: true, isPreview: false,
    headers: '[]', params: '[]', bodyType: 'none', bodyData: '{}',
    authData: '{"type":"none"}', collectionId: 0,
  }
  tabs.value.push(tab)
  selectTab(tab.id)
}

function showHistoryDetail(item: any) { historyItem.value = item }
function clearHistoryDetail() { historyItem.value = null }

function statusClass(code: number): string {
  if (code < 300) return 'status-2xx'
  if (code < 400) return 'status-3xx'
  if (code < 500) return 'status-4xx'
  return 'status-5xx'
}

function formatJson(raw: string): string {
  try { return JSON.stringify(JSON.parse(raw), null, 2) } catch { return raw }
}

async function onSend() {
  if (isSending.value) return
  isSending.value = true
  syncToActiveTab()

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
      status: 0, time: 0, size: 0, headers: {}, cookies: [],
      body: 'Error: ' + (e?.message || e), rawRequest: '', curlCommand: '',
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
    name: activeTab.value.name, description: '',
    method: currentMethod.value, url: currentUrl.value,
    headers: currentHeaders.value, params: currentParams.value,
    body_type: currentBodyType.value, body: currentBodyData.value,
    auth: currentAuthData.value, sort_order: 0,
  }
  try {
    await UpdateRequest(reqModel as any)
    if (activeTab.value) activeTab.value.isDirty = false
  } catch (e: any) {
    dialog.warning({ title: '保存失败', content: e?.message || String(e), positiveText: '确定' })
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
      positiveText: '确定', negativeText: '取消',
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
      activeTabId.value = null; activeTab.value = null
    }
  }
}

function closeOthers(tabId: string) {
  const dirtyOthers = tabs.value.some(t => t.id !== tabId && t.isDirty)
  const doClose = () => { tabs.value = tabs.value.filter(t => t.id === tabId); selectTab(tabId) }
  if (dirtyOthers) {
    dialog.warning({ title: '确认关闭', content: '未保存的修改将丢失，确定关闭？', positiveText: '确定', negativeText: '取消', onPositiveClick: doClose })
  } else { doClose() }
}

function closeRight(idx: number) {
  const dirtyRight = tabs.value.slice(idx + 1).some(t => t.isDirty)
  const doClose = () => {
    tabs.value = tabs.value.slice(0, idx + 1)
    if (activeTabId.value && !tabs.value.find(t => t.id === activeTabId.value))
      selectTab(tabs.value[tabs.value.length - 1].id)
  }
  if (dirtyRight) {
    dialog.warning({ title: '确认关闭', content: '未保存的修改将丢失，确定关闭？', positiveText: '确定', negativeText: '取消', onPositiveClick: doClose })
  } else { doClose() }
}

function closeAll() {
  const dirty = tabs.value.some(t => t.isDirty)
  const doClose = () => { tabs.value = []; activeTabId.value = null; activeTab.value = null }
  if (dirty) {
    dialog.warning({ title: '确认关闭', content: '未保存的修改将丢失，确定关闭？', positiveText: '确定', negativeText: '取消', onPositiveClick: doClose })
  } else { doClose() }
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

function onDragStart(e: DragEvent, idx: number) { dragIdx.value = idx; e.dataTransfer!.effectAllowed = 'move' }
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
  if (oldPreview !== -1) tabs.value.splice(oldPreview, 1)
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
  background: var(--bg-base);
  min-width: 0;
  overflow: hidden;
}
.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--text-muted);
  gap: 4px;
}
.empty-logo { margin-bottom: 8px; }
.empty-state h2 {
  font-size: var(--fs-lg); color: var(--text-secondary); margin: 0;
  font-weight: 600; letter-spacing: 1px; text-transform: uppercase;
}
.empty-state p {
  font-size: var(--fs-sm); color: var(--text-muted); margin: 0;
  text-align: center; line-height: 1.7;
}
.empty-state kbd {
  background: var(--bg-elevated); border: 1px solid var(--border-primary);
  padding: 1px 5px; border-radius: 2px; font-size: var(--fs-xs);
  font-family: var(--font-mono); color: var(--text-secondary);
}
.workspace-editor {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}
.tabs-bar {
  display: flex;
  background: var(--bg-surface);
  border-bottom: 1px solid var(--border-primary);
  height: 32px;
  align-items: flex-end;
  padding: 0 4px;
  gap: 2px;
  overflow-x: auto;
  flex-shrink: 0;
}
.tabs-msg {
  font-size: var(--fs-sm);
  color: var(--text-muted);
  padding: 5px 12px;
  font-family: var(--font-mono);
}
.tab {
  padding: 4px 10px;
  font-size: var(--fs-sm);
  background: var(--bg-elevated);
  border: 1px solid var(--border-primary);
  border-bottom: none;
  border-radius: var(--radius-sm) var(--radius-sm) 0 0;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 4px;
  white-space: nowrap;
  user-select: none;
  color: var(--text-muted);
  font-family: var(--font-mono);
  transition: all var(--transition);
  max-width: 200px;
}
.tab.active { background: var(--bg-base); border-color: var(--border-primary); border-bottom: 2px solid var(--accent); color: var(--text-primary); }
.tab.preview { font-style: italic; }
.tab:hover:not(.active) { background: var(--bg-hover); color: var(--text-secondary); }
.tab-method {
  font-size: var(--fs-2xs); font-weight: 700; padding: 1px 3px; border-radius: 2px; letter-spacing: 0.3px; flex-shrink: 0;
}
.tab-method.get { background: var(--accent-soft); color: var(--accent); }
.tab-method.post { background: var(--amber-soft); color: var(--amber); }
.tab-method.put { background: var(--blue-soft); color: var(--blue); }
.tab-method.delete { background: var(--red-soft); color: var(--red); }
.tab-method.patch { background: var(--purple-soft); color: var(--purple); }
.tab-method.head, .tab-method.options { background: var(--bg-hover); color: var(--text-secondary); }
.tab-name { font-size: var(--fs-sm); overflow: hidden; text-overflow: ellipsis; }
.tab-dirty { width: 5px; height: 5px; background: var(--accent); border-radius: 50%; flex-shrink: 0; }
.tab-close { color: var(--text-muted); font-size: var(--fs-md); margin-left: 1px; padding: 0 2px; border-radius: 2px; transition: all var(--transition); }
.tab-close:hover { color: var(--red); background: var(--red-soft); }
.tab-plus {
  padding: 2px 7px; font-size: var(--fs-md); cursor: pointer; color: var(--accent);
  user-select: none; border-radius: var(--radius-sm); transition: all var(--transition); flex-shrink: 0;
}
.tab-plus:hover { background: var(--accent-soft); }
.resize-handle {
  height: 6px; cursor: ns-resize; flex-shrink: 0;
  display: flex; align-items: center; justify-content: center;
  background: var(--bg-surface);
  border-top: 1px solid var(--border-primary);
  transition: background var(--transition);
}
.resize-handle:hover { background: var(--bg-elevated); }
.resize-line {
  width: 30px; height: 2px; border-radius: 1px;
  background: var(--border-hover); transition: all var(--transition);
}
.resize-handle:hover .resize-line { background: var(--accent); width: 50px; }

.history-detail {
  flex: 1;
  overflow-y: auto;
  padding: 16px 18px;
  font-size: var(--fs-sm);
  background: var(--bg-base);
}
.hist-detail-row {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px 0;
  border-bottom: 1px solid var(--border-primary);
}
.hist-detail-label {
  color: var(--text-muted);
  width: 60px;
  flex-shrink: 0;
  font-size: var(--fs-xs);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  font-weight: 500;
}
.hist-detail-value {
  word-break: break-all;
  color: var(--text-primary);
  font-size: var(--fs-sm);
  font-family: var(--font-mono);
}
.hist-detail-value.status-2xx { color: var(--accent); font-weight: 600; }
.hist-detail-value.status-3xx { color: var(--blue); font-weight: 600; }
.hist-detail-value.status-4xx { color: var(--amber); font-weight: 600; }
.hist-detail-value.status-5xx { color: var(--red); font-weight: 600; }
.method-badge {
  font-size: var(--fs-2xs); font-weight: 700; padding: 2px 5px; border-radius: 2px; letter-spacing: 0.3px;
}
.method-badge.get { background: var(--accent-soft); color: var(--accent); }
.method-badge.post { background: var(--amber-soft); color: var(--amber); }
.method-badge.put { background: var(--blue-soft); color: var(--blue); }
.method-badge.delete { background: var(--red-soft); color: var(--red); }
.method-badge.patch { background: var(--purple-soft); color: var(--purple); }
.hist-detail-section { margin-top: 12px; }
.hist-detail-section h4 {
  margin: 0 0 6px; font-size: var(--fs-xs); color: var(--text-muted);
  font-weight: 600; text-transform: uppercase; letter-spacing: 0.5px;
}
.hist-detail-pre {
  background: var(--bg-surface); border: 1px solid var(--border-primary); border-radius: var(--radius);
  padding: 10px; font-size: var(--fs-sm); font-family: var(--font-mono);
  white-space: pre-wrap; word-break: break-all; max-height: 180px;
  overflow-y: auto; margin: 0; color: var(--text-primary); line-height: 1.6;
}
</style>
