<template>
  <div class="workspace">
    <!-- ════════ Enhanced Empty State ════════ -->
    <div v-if="!activeTab && !historyItem" class="empty-state">
      <div class="empty-logo">
        <svg width="64" height="64" viewBox="0 0 64 64" fill="none">
          <!-- Paw print: main pad -->
          <ellipse cx="32" cy="42" rx="10" ry="8.5" fill="currentColor" stroke="currentColor" stroke-width="1.2" opacity="0.25"/>
          <!-- Toe beans -->
          <ellipse cx="18" cy="26" rx="5" ry="6" fill="currentColor" stroke="currentColor" stroke-width="1" opacity="0.2"/>
          <ellipse cx="28" cy="18" rx="4.5" ry="5.5" fill="currentColor" stroke="currentColor" stroke-width="1" opacity="0.2"/>
          <ellipse cx="36" cy="18" rx="4.5" ry="5.5" fill="currentColor" stroke="currentColor" stroke-width="1" opacity="0.2"/>
          <ellipse cx="46" cy="26" rx="5" ry="6" fill="currentColor" stroke="currentColor" stroke-width="1" opacity="0.2"/>
        </svg>
      </div>
      <h2>Paw API</h2>
      <p class="empty-subtitle">轻量级 API 调试工具 · 为效率而生</p>

      <div class="feature-cards">
        <div class="feature-card" @click="addNewTab">
          <div class="feature-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
              <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
            </svg>
          </div>
          <span class="feature-title">新建请求</span>
          <span class="feature-desc">创建 API 调试请求</span>
        </div>
        <div class="feature-card" @click="reopenLastClosed">
          <div class="feature-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
              <polyline points="1 4 1 10 7 10"/><path d="M3.51 15a9 9 0 1 0 2.13-9.36L1 10"/>
            </svg>
          </div>
          <span class="feature-title">恢复标签</span>
          <span class="feature-desc">重新打开已关闭标签</span>
        </div>
        <div class="feature-card">
          <div class="feature-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
              <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/>
            </svg>
          </div>
          <span class="feature-title">快速发送</span>
          <span class="feature-desc">Ctrl+Enter 立即调试</span>
        </div>
      </div>

      <div class="quick-start">
        <h3>快速开始</h3>
        <div class="qs-items">
          <div class="qs-item" @click="addNewTab">
            <span class="qs-step">1</span>
            <span>点击左侧集合中的请求，或创建新请求</span>
          </div>
          <div class="qs-item">
            <span class="qs-step">2</span>
            <span>输入 URL，配置请求头和参数</span>
          </div>
          <div class="qs-item">
            <span class="qs-step">3</span>
            <span>按 <kbd>Ctrl</kbd>+<kbd>Enter</kbd> 发送请求查看响应</span>
          </div>
        </div>
      </div>

      <div class="shortcuts-hint">
        <span class="shortcut"><kbd>Ctrl</kbd>+<kbd>N</kbd> 新建标签</span>
        <span class="shortcut"><kbd>Ctrl</kbd>+<kbd>S</kbd> 保存</span>
        <span class="shortcut"><kbd>Ctrl</kbd>+<kbd>Enter</kbd> 发送</span>
        <span class="shortcut"><kbd>Ctrl</kbd>+<kbd>W</kbd> 关闭</span>
        <span class="shortcut"><kbd>Ctrl</kbd>+<kbd>Shift</kbd>+<kbd>T</kbd> 恢复</span>
        <span class="shortcut"><kbd>Ctrl</kbd>+<kbd>1-9</kbd> 切换</span>
      </div>
    </div>

    <!-- ════════ History Detail View ════════ -->
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

    <!-- ════════ Active Tab Editor ════════ -->
    <div v-else-if="activeTab" class="workspace-editor">
      <!-- ── Tab Bar with Scroll Arrows ── -->
      <div class="tabs-bar">
        <button
          v-if="tabsScrollable"
          class="tab-scroll-btn tab-scroll-left"
          @click="scrollTabs(-1)"
        >
          <svg width="10" height="10" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><polyline points="8 2 4 6 8 10"/></svg>
        </button>

        <div class="tabs-scroll-container" ref="tabsScrollRef">
          <div
            v-for="(tab, idx) in tabs"
            :key="tab.id"
            class="tab"
            :class="{
              active: tab.id === activeTabId,
              dirty: tab.isDirty,
              preview: tab.isPreview,
              'drag-over': dragOverId === tab.id,
              renaming: renamingTabId === tab.id,
            }"
            draggable="true"
            @click="selectTab(tab.id)"
            @auxclick.middle="onMiddleClick(tab)"
            @dblclick="startRename(tab)"
            @contextmenu.prevent="onTabContextMenu($event, tab, idx)"
            @dragstart="onDragStart($event, idx)"
            @dragover.prevent="onDragOver($event, tab.id)"
            @dragleave="onDragLeave(tab.id)"
            @drop="onDrop($event, idx)"
            @dragend="onDragEnd"
          >
            <span class="tab-method" :class="tab.method?.toLowerCase()">{{ tab.method }}</span>

            <!-- Inline rename input -->
            <input
              v-if="renamingTabId === tab.id"
              ref="renameInputRef"
              v-model="renameValue"
              class="tab-rename-input"
              @click.stop
              @blur="finishRename"
              @keydown.enter="($event.target as HTMLInputElement).blur()"
              @keydown.escape="cancelRename"
              @mousedown.stop
            />
            <span v-else class="tab-name">{{ tab.name }}</span>

            <span v-if="tab.isDirty" class="tab-dirty" title="未保存的更改"></span>
            <span class="tab-close" @click.stop="onCloseTab(tab)" title="关闭 (Ctrl+W)">
              <svg width="10" height="10" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><line x1="3" y1="3" x2="9" y2="9"/><line x1="9" y1="3" x2="3" y2="9"/></svg>
            </span>
          </div>
        </div>

        <button
          v-if="tabsScrollable"
          class="tab-scroll-btn tab-scroll-right"
          @click="scrollTabs(1)"
        >
          <svg width="10" height="10" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><polyline points="4 2 8 6 4 10"/></svg>
        </button>

        <span class="tab-plus" @click="addNewTab" title="新建标签 (Ctrl+N)">+</span>

        <span v-if="tabs.length > 1" class="tab-count">{{ activeTabIndex + 1 }}/{{ tabs.length }}</span>
      </div>

      <!-- ── Pill-style View Switcher ── -->
      <div class="view-switcher">
        <div class="view-toggle">
          <button
            class="view-toggle-btn"
            :class="{ active: viewMode === 'request' }"
            @click="viewMode = 'request'"
          >
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
              <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
            </svg>
            请求
          </button>
          <button
            class="view-toggle-btn"
            :class="{ active: viewMode === 'docs' }"
            @click="onSwitchToDocs"
          >
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
              <polyline points="14 2 14 8 20 8"/>
              <line x1="16" y1="13" x2="8" y2="13"/>
              <line x1="16" y1="17" x2="8" y2="17"/>
            </svg>
            文档
          </button>
        </div>

        <!-- Name editor with pencil icon -->
        <div class="name-editor">
          <div class="name-field" :class="{ focused: nameFocused }">
            <svg class="pencil-icon" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
              <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
            </svg>
            <input
              v-model="currentName"
              class="name-input"
              placeholder="接口名称"
              @focus="nameFocused = true"
              @blur="onNameBlur"
              @keydown.enter="($event.target as HTMLInputElement).blur()"
            />
          </div>
        </div>
      </div>

      <template v-if="viewMode === 'request'">
        <UrlBar
          :model-method="currentMethod"
          :model-url="currentUrl"
          :loading="isSending"
          @update:model-method="onMethodChange"
          @update:model-url="onUrlChange"
          @send="onSend"
          @save="onSave"
        />

        <RequestPanel
          :headers="currentHeaders"
          :params="currentParams"
          :params-enabled="currentParamsEnabled"
          :body-type="currentBodyType"
          :body-data="currentBodyData"
          :auth-data="currentAuthData"
          :url="currentUrl"
          :path-vars="currentPathVars"
          @update:headers="onHeadersChange"
          @update:params="onParamsChange"
          @update:params-enabled="onParamsEnabledChange"
          @update:body-type="onBodyTypeChange"
          @update:body-data="onBodyDataChange"
          @update:auth-data="onAuthDataChange"
          @update:path-vars="onPathVarsChange"
        />

        <!-- Improved resize handle with dots pattern -->
        <div
          v-if="response"
          class="resize-handle"
          @mousedown="onResizeStart"
        >
          <div class="resize-dots">
            <span class="resize-dot"></span>
            <span class="resize-dot"></span>
            <span class="resize-dot"></span>
            <span class="resize-dot"></span>
            <span class="resize-dot"></span>
          </div>
        </div>

        <Transition name="response-slide">
          <ResponsePanel
            v-if="response"
            :response="response"
            :style="{ height: responseHeight + 'px' }"
          />
        </Transition>
      </template>

      <RequestDocsView
        v-else
        :request-id="activeTab?.requestId || 0"
        :request-name="currentName"
        :request-method="currentMethod"
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
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { NDropdown, useDialog, useMessage } from 'naive-ui'
import UrlBar from '../request/UrlBar.vue'
import RequestPanel from '../request/RequestPanel.vue'
import ResponsePanel from '../response/ResponsePanel.vue'
import RequestDocsView from '../request/RequestDocsView.vue'
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
  paramsEnabled: boolean
  bodyType: string
  bodyData: string
  authData: string
  collectionId: number
}

const props = defineProps<{
  projectId: number | null
}>()

const emit = defineEmits<{
  (e: 'request-saved'): void
}>()

const envStore = useEnvStore()
const dialog = useDialog()
const message = useMessage()

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
const currentParamsEnabled = ref(true)
const currentBodyType = ref('none')
const currentBodyData = ref('{}')
const currentAuthData = ref('{"type":"none"}')
const currentName = ref('')
const viewMode = ref<'request' | 'docs'>('request')

let sessionCounter = 0
let tabIdCounter = 0

// ── Closed tabs stack (max 10) ──
const closedTabs: Tab[] = []
const MAX_CLOSED_TABS = 10

// ── Inline rename state ──
const renamingTabId = ref<string | null>(null)
const renameValue = ref('')
const renameInputRef = ref<HTMLInputElement[] | null>(null)

// ── Tab scroll state ──
const tabsScrollRef = ref<HTMLElement | null>(null)
const tabsScrollable = ref(false)

// ── Name input focus state ──
const nameFocused = ref(false)

// ── Active tab index (for count display) ──
const activeTabIndex = computed(() => {
  if (!activeTabId.value) return 0
  const idx = tabs.value.findIndex(t => t.id === activeTabId.value)
  return idx >= 0 ? idx : 0
})

// ── Detect tab overflow ──
function detectOverflow() {
  const el = tabsScrollRef.value
  if (el) {
    tabsScrollable.value = el.scrollWidth > el.clientWidth + 4
  } else {
    tabsScrollable.value = false
  }
}

function scrollTabs(direction: number) {
  const el = tabsScrollRef.value
  if (el) {
    el.scrollBy({ left: direction * 180, behavior: 'smooth' })
    // Re-check overflow after scroll settles
    setTimeout(detectOverflow, 300)
  }
}

let resizeObserver: ResizeObserver | null = null

watch(tabs, () => {
  nextTick(detectOverflow)
}, { deep: true })

watch(activeTabId, () => {
  nextTick(detectOverflow)
  // Scroll active tab into view
  nextTick(() => {
    const container = tabsScrollRef.value
    if (!container || !activeTabId.value) return
    const activeEl = container.querySelector('.tab.active') as HTMLElement
    if (activeEl) {
      activeEl.scrollIntoView({ block: 'nearest', inline: 'nearest', behavior: 'smooth' })
    }
  })
})

// ── Inline rename ──
function startRename(tab: Tab) {
  renamingTabId.value = tab.id
  renameValue.value = tab.name
  nextTick(() => {
    const inputs = renameInputRef.value
    if (inputs && inputs.length > 0) {
      inputs[0].focus()
      inputs[0].select()
    }
  })
}

function finishRename() {
  if (!renamingTabId.value) return
  const tab = tabs.value.find(t => t.id === renamingTabId.value)
  if (tab && renameValue.value.trim()) {
    tab.name = renameValue.value.trim()
    tab.isDirty = true
    if (tab.id === activeTabId.value) {
      currentName.value = tab.name
    }
  }
  renamingTabId.value = null
}

function cancelRename() {
  renamingTabId.value = null
}

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
function onParamsEnabledChange(v: boolean) { currentParamsEnabled.value = v; markDirty() }
function onBodyTypeChange(v: string) { currentBodyType.value = v; markDirty() }
function onBodyDataChange(v: string) { currentBodyData.value = v; markDirty() }
function onAuthDataChange(v: string) { currentAuthData.value = v; markDirty() }

async function onNameBlur() {
  nameFocused.value = false
  if (!activeTab.value) return
  if (activeTab.value.name !== currentName.value) {
    activeTab.value.name = currentName.value
    if (activeTab.value.isPreview) activeTab.value.isPreview = false
    // Auto-save just the name
    if (activeTab.value.requestId > 0) {
      try {
        await UpdateRequest({
          id: activeTab.value.requestId,
          collection_id: activeTab.value.collectionId || 0,
          name: currentName.value,
          description: '',
          method: currentMethod.value,
          url: currentUrl.value,
          headers: currentHeaders.value,
          params: currentParams.value,
          body_type: currentBodyType.value,
          body: currentBodyData.value,
          auth: currentAuthData.value,
          sort_order: 0,
        } as any)
      } catch { /* silent */ }
      emit('request-saved')
    }
  }
}

function onSwitchToDocs() {
  viewMode.value = 'docs'
}

function syncToActiveTab() {
  if (!activeTab.value) return
  activeTab.value.name = currentName.value
  activeTab.value.method = currentMethod.value
  activeTab.value.url = currentUrl.value
  activeTab.value.pathVars = currentPathVars.value
  activeTab.value.headers = currentHeaders.value
  activeTab.value.params = currentParams.value
  activeTab.value.paramsEnabled = currentParamsEnabled.value
  activeTab.value.bodyType = currentBodyType.value
  activeTab.value.bodyData = currentBodyData.value
  activeTab.value.authData = currentAuthData.value
}

function addNewTab() {
  const id = 'tab-' + (++tabIdCounter) + '-' + Date.now()
  const tab: Tab = {
    id, requestId: 0, method: 'GET', name: '新建请求', url: '',
    pathVars: '[]', isDirty: true, isPreview: false,
    headers: '[]', params: '[]', paramsEnabled: true, bodyType: 'none', bodyData: '{}',
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
    params: currentParamsEnabled.value ? currentParams.value : '[]',
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
      body: '请求失败: ' + (e?.message || e), rawRequest: '', curlCommand: '',
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
    message.success('已保存')
    emit('request-saved')
  } catch (e: any) {
    message.error('保存失败: ' + (e?.message || String(e)))
  }
}

function selectTab(id: string) {
  activeTabId.value = id
  activeTab.value = tabs.value.find(t => t.id === id) || null
  if (activeTab.value) {
    currentName.value = activeTab.value.name
    currentMethod.value = activeTab.value.method
    currentUrl.value = activeTab.value.url
    currentPathVars.value = activeTab.value.pathVars || '[]'
    currentHeaders.value = activeTab.value.headers
    currentParams.value = activeTab.value.params
    currentParamsEnabled.value = activeTab.value.paramsEnabled ?? true
    currentBodyType.value = activeTab.value.bodyType
    currentBodyData.value = activeTab.value.bodyData
    currentAuthData.value = activeTab.value.authData
  }
  viewMode.value = 'request'
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
  if (idx === -1) return
  const removed = tabs.value[idx]

  // Push to closed tabs stack
  closedTabs.push({ ...removed })
  if (closedTabs.length > MAX_CLOSED_TABS) {
    closedTabs.shift()
  }

  tabs.value = tabs.value.filter(t => t.id !== id)
  if (activeTabId.value === id) {
    if (tabs.value.length > 0) {
      selectTab(tabs.value[Math.min(idx, tabs.value.length - 1)].id)
    } else {
      activeTabId.value = null; activeTab.value = null
    }
  }
}

function reopenLastClosed() {
  if (closedTabs.length === 0) return
  const tab = closedTabs.pop()!
  // Avoid duplicate if already reopened
  const existing = tabs.value.find(t => t.id === tab.id)
  if (existing) {
    selectTab(existing.id)
    return
  }
  tabs.value.push(tab)
  selectTab(tab.id)
}

function closeOthers(tabId: string) {
  const dirtyOthers = tabs.value.some(t => t.id !== tabId && t.isDirty)
  const doClose = () => {
    // Save closed tabs to stack
    const others = tabs.value.filter(t => t.id !== tabId)
    for (const t of others) {
      closedTabs.push({ ...t })
      if (closedTabs.length > MAX_CLOSED_TABS) closedTabs.shift()
    }
    tabs.value = tabs.value.filter(t => t.id === tabId)
    selectTab(tabId)
  }
  if (dirtyOthers) {
    dialog.warning({ title: '确认关闭', content: '未保存的修改将丢失，确定关闭？', positiveText: '确定', negativeText: '取消', onPositiveClick: doClose })
  } else { doClose() }
}

function closeRight(idx: number) {
  const dirtyRight = tabs.value.slice(idx + 1).some(t => t.isDirty)
  const doClose = () => {
    const rightTabs = tabs.value.slice(idx + 1)
    for (const t of rightTabs) {
      closedTabs.push({ ...t })
      if (closedTabs.length > MAX_CLOSED_TABS) closedTabs.shift()
    }
    tabs.value = tabs.value.slice(0, idx + 1)
    if (activeTabId.value && !tabs.value.find(t => t.id === activeTabId.value))
      selectTab(tabs.value[tabs.value.length - 1].id)
  }
  if (dirtyRight) {
    dialog.warning({ title: '确认关闭', content: '未保存的修改将丢失，确定关闭？', positiveText: '确定', negativeText: '取消', onPositiveClick: doClose })
  } else { doClose() }
}

function closeLeft(idx: number) {
  const dirtyLeft = tabs.value.slice(0, idx).some(t => t.isDirty)
  const doClose = () => {
    const leftTabs = tabs.value.slice(0, idx)
    for (const t of leftTabs) {
      closedTabs.push({ ...t })
      if (closedTabs.length > MAX_CLOSED_TABS) closedTabs.shift()
    }
    tabs.value = tabs.value.slice(idx)
    if (activeTabId.value && !tabs.value.find(t => t.id === activeTabId.value))
      selectTab(tabs.value[0].id)
  }
  if (dirtyLeft) {
    dialog.warning({ title: '确认关闭', content: '未保存的修改将丢失，确定关闭？', positiveText: '确定', negativeText: '取消', onPositiveClick: doClose })
  } else { doClose() }
}

function closeAll() {
  const dirty = tabs.value.some(t => t.isDirty)
  const doClose = () => {
    for (const t of tabs.value) {
      closedTabs.push({ ...t })
      if (closedTabs.length > MAX_CLOSED_TABS) closedTabs.shift()
    }
    tabs.value = []; activeTabId.value = null; activeTab.value = null
  }
  if (dirty) {
    dialog.warning({ title: '确认关闭', content: '未保存的修改将丢失，确定关闭？', positiveText: '确定', negativeText: '取消', onPositiveClick: doClose })
  } else { doClose() }
}

// ── Middle-click to close ──
function onMiddleClick(tab: Tab) {
  onCloseTab(tab)
}

// ── Context menu ──
const ctxMenuShow = ref(false)
const ctxMenuX = ref(0)
const ctxMenuY = ref(0)
let ctxTarget: { tab: Tab; idx: number } | null = null

const ctxMenuOptions = computed(() => {
  const opts = [
    { label: '关闭', key: 'close' },
    { label: '关闭其他', key: 'close-others' },
  ]
  if (ctxTarget && ctxTarget.idx > 0) {
    opts.push({ label: '关闭左侧', key: 'close-left' })
  }
  if (ctxTarget && ctxTarget.idx < tabs.value.length - 1) {
    opts.push({ label: '关闭右侧', key: 'close-right' })
  }
  opts.push({ label: '关闭全部', key: 'close-all' })
  if (closedTabs.length > 0) {
    opts.push({ label: `恢复已关闭 (${closedTabs.length})`, key: 'reopen' })
  }
  return opts
})

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
    case 'close-left': closeLeft(ctxTarget.idx); break
    case 'close-right': closeRight(ctxTarget.idx); break
    case 'close-all': closeAll(); break
    case 'reopen': reopenLastClosed(); break
  }
}

// ── Drag and drop ──
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

// ── Response panel resize ──
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

// ── Keyboard shortcuts ──
function onKeydown(e: KeyboardEvent) {
  const ctrl = e.ctrlKey || e.metaKey

  // Ctrl+S — save
  if (ctrl && e.key === 's') {
    e.preventDefault()
    onSave()
  }

  // Ctrl+N — new tab
  if (ctrl && e.key === 'n') {
    e.preventDefault()
    addNewTab()
  }

  // Ctrl+W — close tab (or clear to empty if only one tab)
  if (ctrl && e.key === 'w') {
    e.preventDefault()
    if (activeTab.value) {
      if (tabs.value.length <= 1) {
        // Don't close the last tab — just clear to empty state
        if (activeTab.value.isDirty) {
          dialog.warning({
            title: '确认关闭',
            content: '未保存的修改将丢失，确定关闭？',
            positiveText: '确定', negativeText: '取消',
            onPositiveClick: () => {
              const removed = tabs.value[0]
              closedTabs.push({ ...removed })
              if (closedTabs.length > MAX_CLOSED_TABS) closedTabs.shift()
              tabs.value = []; activeTabId.value = null; activeTab.value = null
            },
          })
        } else {
          const removed = tabs.value[0]
          closedTabs.push({ ...removed })
          if (closedTabs.length > MAX_CLOSED_TABS) closedTabs.shift()
          tabs.value = []; activeTabId.value = null; activeTab.value = null
        }
      } else {
        onCloseTab(activeTab.value)
      }
    }
  }

  // Ctrl+Enter — send request
  if (ctrl && e.key === 'Enter') {
    e.preventDefault()
    if (activeTab.value && viewMode.value === 'request') onSend()
  }

  // Ctrl+Shift+T — reopen last closed tab
  if (ctrl && e.shiftKey && (e.key === 'T' || e.key === 't')) {
    e.preventDefault()
    reopenLastClosed()
  }

  // Ctrl+1 through Ctrl+9 — switch to tab N
  if (ctrl && !e.shiftKey && !e.altKey && e.key >= '1' && e.key <= '9') {
    e.preventDefault()
    const idx = parseInt(e.key) - 1
    if (idx < tabs.value.length) {
      selectTab(tabs.value[idx].id)
    }
  }
}

// Watch the scroll container ref (conditionally rendered) and observe when available
watch(tabsScrollRef, (el, oldEl) => {
  if (resizeObserver && oldEl) {
    resizeObserver.unobserve(oldEl)
  }
  if (resizeObserver && el) {
    resizeObserver.observe(el)
    detectOverflow()
  }
})

onMounted(() => {
  document.addEventListener('keydown', onKeydown)
  resizeObserver = new ResizeObserver(() => {
    detectOverflow()
  })
  // If the scroll container is already rendered (e.g., tab opened before mount)
  nextTick(() => {
    if (tabsScrollRef.value) {
      resizeObserver!.observe(tabsScrollRef.value)
    }
  })
})

onUnmounted(() => {
  document.removeEventListener('keydown', onKeydown)
  if (resizeObserver) {
    resizeObserver.disconnect()
    resizeObserver = null
  }
})

function renameOpenTab(requestId: number, newName: string) {
  const tab = tabs.value.find(t => t.requestId === requestId)
  if (tab) {
    tab.name = newName
    if (tab.id === activeTabId.value) {
      currentName.value = newName
    }
  }
}

defineExpose({ openTab, previewTab, tabs, activeTabId, selectTab, showHistoryDetail, clearHistoryDetail, renameOpenTab })
</script>

<style scoped>
.workspace {
  flex: 1;
  display: flex;
  background: var(--bg-base);
  min-width: 0;
  overflow: hidden;
}

/* ════════════════════════════════════════════
   Empty State — Paw Print + Feature Cards
   ════════════════════════════════════════════ */
.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--text-muted);
  gap: 10px;
  animation: slideUp 0.4s var(--ease-out) both;
  padding: 40px 28px;
  overflow-y: auto;
}

.empty-logo {
  margin-bottom: 10px;
  color: var(--accent);
  opacity: 0.75;
  animation: fadeIn 0.6s var(--ease-out) 0.1s both;
}
.empty-logo svg { width: 64px; height: 64px; }

.empty-state h2 {
  font-size: var(--fs-xl); color: var(--text-secondary); margin: 0;
  font-weight: 700; letter-spacing: 3px; text-transform: uppercase;
  font-family: var(--font-ui);
}
.empty-subtitle {
  font-size: var(--fs-sm); color: var(--text-secondary); margin: 0;
  text-align: center; line-height: 1.6;
}

/* ── Feature Cards ── */
.feature-cards {
  display: flex;
  gap: 12px;
  margin-top: 24px;
  flex-wrap: wrap;
  justify-content: center;
}
.feature-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 7px;
  padding: 16px 20px;
  background: var(--bg-surface);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all 0.22s var(--ease-out);
  min-width: 130px;
  animation: fadeIn 0.5s var(--ease-out) 0.2s both;
}
.feature-card:hover {
  border-color: var(--accent);
  background: var(--accent-soft);
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0,0,0,0.08);
}
.feature-icon {
  width: 36px; height: 36px;
  display: flex; align-items: center; justify-content: center;
  background: var(--accent-soft);
  border-radius: var(--radius);
  color: var(--accent);
  transition: all 0.22s var(--ease-out);
}
.feature-card:hover .feature-icon {
  background: var(--accent-glow);
  transform: scale(1.08);
}
.feature-title {
  font-size: var(--fs-sm);
  font-weight: 600;
  color: var(--text-primary);
  font-family: var(--font-ui);
}
.feature-desc {
  font-size: var(--fs-xs);
  color: var(--text-secondary);
  font-family: var(--font-ui);
}

/* ── Quick Start ── */
.quick-start {
  margin-top: 28px;
  max-width: 420px;
  width: 100%;
  animation: fadeIn 0.5s var(--ease-out) 0.35s both;
}
.quick-start h3 {
  font-size: var(--fs-sm);
  color: var(--text-secondary);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  margin: 0 0 12px 0;
  font-family: var(--font-ui);
  text-align: center;
}
.qs-items {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.qs-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 14px;
  background: var(--bg-surface);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  font-size: var(--fs-sm);
  color: var(--text-secondary);
  font-family: var(--font-ui);
  cursor: default;
  transition: all var(--transition);
}
.qs-item:hover {
  border-color: var(--border-hover);
  background: var(--bg-hover);
}
.qs-step {
  width: 20px; height: 20px;
  display: flex; align-items: center; justify-content: center;
  background: var(--accent-soft);
  color: var(--accent);
  font-size: var(--fs-2xs);
  font-weight: 700;
  border-radius: 50%;
  flex-shrink: 0;
  font-family: var(--font-mono);
}
.qs-item kbd {
  display: inline-block;
  background: var(--bg-elevated); border: 1px solid var(--border-primary);
  padding: 1px 5px; border-radius: var(--radius-xs); font-size: var(--fs-2xs);
  font-family: var(--font-mono); color: var(--text-secondary);
  box-shadow: 0 1px 0 var(--border-primary);
  margin: 0 1px;
}

/* ── Shortcuts Hint ── */
.shortcuts-hint {
  display: flex;
  gap: 14px;
  margin-top: 24px;
  flex-wrap: wrap;
  justify-content: center;
  animation: fadeIn 0.5s var(--ease-out) 0.45s both;
}
.shortcut {
  font-size: var(--fs-xs);
  color: var(--text-muted);
  font-family: var(--font-ui);
  display: flex;
  align-items: center;
  gap: 4px;
}
.empty-state kbd {
  display: inline-block;
  background: var(--bg-elevated); border: 1px solid var(--border-primary);
  padding: 2px 6px; border-radius: var(--radius-xs); font-size: var(--fs-2xs);
  font-family: var(--font-mono); color: var(--text-secondary);
  box-shadow: 0 1px 0 var(--border-primary);
  margin: 0 1px;
}

/* ════════════════════════════════════════════
   Editor Layout
   ════════════════════════════════════════════ */
.workspace-editor {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

/* ════════════════════════════════════════════
   Tab Bar — Shadow + Scroll Arrows + Count
   ════════════════════════════════════════════ */
.tabs-bar {
  display: flex;
  background: var(--bg-surface);
  border-bottom: 1px solid var(--border-primary);
  height: 38px;
  align-items: flex-end;
  padding: 0 8px;
  gap: 3px;
  flex-shrink: 0;
  position: relative;
  box-shadow: 0 1px 4px rgba(0,0,0,0.06);
  z-index: 5;
}

.tabs-msg {
  font-size: var(--fs-sm);
  color: var(--text-secondary);
  padding: 8px 14px;
  font-family: var(--font-mono);
  letter-spacing: 0.01em;
}

/* ── Scroll container for tabs ── */
.tabs-scroll-container {
  display: flex;
  gap: 2px;
  overflow-x: auto;
  flex: 1;
  min-width: 0;
  align-items: flex-end;
  scrollbar-width: none;
  scroll-behavior: smooth;
}
.tabs-scroll-container::-webkit-scrollbar {
  display: none;
}

/* ── Scroll arrow buttons ── */
.tab-scroll-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 24px;
  background: var(--bg-surface);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-xs);
  cursor: pointer;
  color: var(--text-muted);
  flex-shrink: 0;
  align-self: center;
  transition: all var(--transition);
  margin-bottom: 2px;
  padding: 0;
}
.tab-scroll-btn:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
  border-color: var(--border-hover);
}
.tab-scroll-btn:active {
  background: var(--bg-elevated);
}

/* ── Tab count indicator ── */
.tab-count {
  font-size: var(--fs-xs);
  color: var(--text-muted);
  font-family: var(--font-mono);
  padding: 0 8px 6px 6px;
  white-space: nowrap;
  flex-shrink: 0;
  user-select: none;
  letter-spacing: 0.02em;
}

/* ════════════════════════════════════════════
   Tab Item — Enhanced
   ════════════════════════════════════════════ */
.tab {
  padding: 6px 14px;
  font-size: var(--fs-sm);
  background: var(--bg-elevated);
  border: 1px solid transparent;
  border-bottom: none;
  border-top-left-radius: var(--radius);
  border-top-right-radius: var(--radius);
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 7px;
  white-space: nowrap;
  user-select: none;
  color: var(--text-secondary);
  font-family: var(--font-mono);
  transition: all var(--transition);
  max-width: 200px;
  min-width: 80px;
  position: relative;
  flex-shrink: 0;
  animation: tabSlideIn 0.25s var(--ease-spring) both;
}
@keyframes tabSlideIn {
  from { opacity: 0; transform: translateY(6px) scale(0.97); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}

.tab.drag-over {
  border-top: 2px solid var(--accent);
  box-shadow: 0 -2px 8px var(--accent-glow);
}

/* Active tab with gradient accent line */
.tab.active {
  background: var(--bg-base);
  border-color: var(--border-primary);
  border-bottom: none;
  color: var(--text-primary);
  border-top-color: transparent;
}
.tab.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, var(--accent), var(--accent-hover));
  border-radius: 2px 2px 0 0;
  box-shadow: 0 0 8px var(--accent-glow);
}

/* Dirty tab: tinted background */
.tab.dirty:not(.active) {
  background: var(--accent-soft);
  border-color: rgba(0,224,90,0.06);
}
.tab.dirty:not(.active):hover {
  background: var(--accent-glow);
}

.tab.preview {
  font-style: italic;
  opacity: 0.75;
}
.tab:hover:not(.active) {
  background: var(--bg-hover);
  color: var(--text-secondary);
}

.tab-method {
  font-size: var(--fs-2xs); font-weight: 700; padding: 1px 5px;
  border-radius: var(--radius-xs); letter-spacing: 0.04em;
  flex-shrink: 0; line-height: 1.5; text-transform: uppercase;
}
.tab-method.get    { background: var(--blue-soft); color: var(--method-get); }
.tab-method.post   { background: rgba(34,197,94,0.1); color: var(--method-post); }
.tab-method.put    { background: var(--amber-soft); color: var(--method-put); }
.tab-method.delete { background: var(--red-soft); color: var(--method-delete); }
.tab-method.patch  { background: var(--purple-soft); color: var(--method-patch); }
.tab-method.head,
.tab-method.options { background: rgba(113,113,122,0.1); color: var(--text-secondary); }

.tab-name {
  font-size: var(--fs-sm);
  overflow: hidden;
  text-overflow: ellipsis;
}

/* ── Tab dirty indicator: pulsing dot ── */
.tab-dirty {
  width: 6px; height: 6px;
  background: var(--accent);
  border-radius: 50%;
  flex-shrink: 0;
  animation: pulseGlow 2s infinite, tabDirtyPulse 2s ease-in-out infinite;
}
@keyframes tabDirtyPulse {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.6; transform: scale(1.3); }
}

/* ── Tab rename input ── */
.tab-rename-input {
  font-size: var(--fs-sm);
  font-family: var(--font-mono);
  color: var(--text-primary);
  background: var(--bg-base);
  border: 1px solid var(--accent);
  border-radius: var(--radius-xs);
  padding: 0 4px;
  outline: none;
  width: 100%;
  min-width: 60px;
  max-width: 140px;
  height: 18px;
  box-shadow: 0 0 0 2px var(--accent-glow);
}

/* ── Tab close button ── */
.tab-close {
  color: var(--text-muted);
  margin-left: auto;
  padding: 3px;
  border-radius: var(--radius-xs);
  transition: all var(--transition-fast);
  opacity: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
  flex-shrink: 0;
}
.tab-close svg { width: 10px; height: 10px; }
.tab:hover .tab-close { opacity: 0.7; }
.tab.active .tab-close { opacity: 0.5; }
.tab.active .tab-close:hover,
.tab-close:hover {
  opacity: 1; color: var(--red); background: var(--red-soft);
}

/* Show close button when dirty too */
.tab.dirty .tab-close { opacity: 0.5; }
.tab.dirty .tab-close:hover { opacity: 1; }

.tab-plus {
  padding: 3px 8px; font-size: var(--fs-md); cursor: pointer;
  color: var(--text-muted); user-select: none;
  border-radius: var(--radius-sm);
  transition: all var(--transition);
  flex-shrink: 0; margin-bottom: 2px;
  display: flex; align-items: center; justify-content: center;
}
.tab-plus:hover { background: var(--accent-soft); color: var(--accent); }

/* ════════════════════════════════════════════
   View Switcher — Pill-Style Toggle
   ════════════════════════════════════════════ */
.view-switcher {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 12px;
  background: var(--bg-base);
  border-bottom: 1px solid var(--border-primary);
  flex-shrink: 0;
  height: 38px;
  gap: 14px;
}

/* ── Pill toggle container ── */
.view-toggle {
  display: flex;
  background: var(--bg-surface);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  padding: 2px;
  gap: 2px;
  flex-shrink: 0;
}

.view-toggle-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 5px 16px;
  font-size: var(--fs-sm);
  font-weight: 500;
  color: var(--text-secondary);
  background: transparent;
  border: none;
  border-radius: calc(var(--radius) - 2px);
  cursor: pointer;
  transition: all 0.2s var(--ease-out);
  font-family: var(--font-ui);
  white-space: nowrap;
}
.view-toggle-btn:hover {
  color: var(--text-secondary);
}
.view-toggle-btn.active {
  color: var(--accent);
  background: var(--accent-soft);
  font-weight: 600;
  box-shadow: 0 1px 3px rgba(0,0,0,0.06);
}
.view-toggle-btn svg {
  flex-shrink: 0;
  opacity: 0.7;
}
.view-toggle-btn.active svg {
  opacity: 1;
}

/* ── Name Editor with Pencil Icon ── */
.name-editor {
  display: flex;
  align-items: center;
  gap: 6px;
  flex: 1;
  min-width: 0;
  justify-content: flex-end;
}

.name-field {
  display: flex;
  align-items: center;
  gap: 6px;
  background: var(--bg-surface);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-sm);
  padding: 4px 10px;
  transition: all var(--transition);
  max-width: 320px;
}
.name-field:hover {
  border-color: var(--border-hover);
}
.name-field.focused {
  border-color: var(--accent);
  box-shadow: 0 0 0 2px var(--accent-glow);
}

.pencil-icon {
  flex-shrink: 0;
  color: var(--text-placeholder);
  transition: color var(--transition);
  opacity: 0.7;
}
.name-field.focused .pencil-icon,
.name-field:hover .pencil-icon {
  color: var(--accent);
  opacity: 1;
}

.name-input {
  width: 180px;
  max-width: 240px;
  padding: 0;
  font-size: var(--fs-sm);
  font-family: var(--font-ui);
  color: var(--text-primary);
  background: transparent;
  border: none;
  outline: none;
  box-sizing: border-box;
}
.name-input::placeholder {
  color: var(--text-placeholder);
}

/* ════════════════════════════════════════════
   Resize Handle — Dots Pattern
   ════════════════════════════════════════════ */
.resize-handle {
  height: 10px;
  cursor: ns-resize;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-surface);
  border-top: 1px solid var(--border-primary);
  transition: all 0.25s var(--ease-out);
  position: relative;
  z-index: 10;
}
.resize-handle:hover {
  background: var(--bg-elevated);
  border-top-color: var(--accent);
}

.resize-dots {
  display: flex;
  align-items: center;
  gap: 4px;
  transition: all 0.25s var(--ease-out);
}
.resize-dot {
  width: 3px;
  height: 3px;
  border-radius: 50%;
  background: var(--border-hover);
  transition: all 0.25s var(--ease-out);
}
.resize-handle:hover .resize-dot {
  background: var(--accent);
  box-shadow: 0 0 4px var(--accent-glow);
}
.resize-handle:hover .resize-dot:nth-child(1),
.resize-handle:hover .resize-dot:nth-child(5) {
  transform: scale(0.7);
  opacity: 0.5;
}
.resize-handle:hover .resize-dot:nth-child(3) {
  transform: scale(1.4);
}

/* ════════════════════════════════════════════
   Response Transition
   ════════════════════════════════════════════ */
.response-slide-enter-active { transition: all 0.3s var(--ease-out); }
.response-slide-leave-active { transition: all 0.2s var(--ease-out); }
.response-slide-enter-from { opacity: 0; transform: translateY(10px); max-height: 0; }
.response-slide-leave-to { opacity: 0; transform: translateY(4px); }

/* ════════════════════════════════════════════
   History Detail
   ════════════════════════════════════════════ */
.history-detail {
  flex: 1;
  overflow-y: auto;
  padding: 20px 22px;
  font-size: var(--fs-sm);
  background: var(--bg-base);
  animation: fadeIn 0.25s var(--ease-out) both;
}
.hist-detail-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 12px;
  border-radius: var(--radius-sm);
  background: var(--bg-surface);
  margin-bottom: 4px;
}
.hist-detail-label {
  color: var(--text-secondary);
  width: 60px;
  flex-shrink: 0;
  font-size: var(--fs-sm);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  font-weight: 600;
}
.hist-detail-value {
  word-break: break-all;
  color: var(--text-primary);
  font-size: var(--fs-sm);
  font-family: var(--font-mono);
}
.hist-detail-value.status-2xx { color: var(--method-post); font-weight: 600; }
.hist-detail-value.status-3xx { color: var(--blue); font-weight: 600; }
.hist-detail-value.status-4xx { color: var(--amber); font-weight: 600; }
.hist-detail-value.status-5xx { color: var(--red); font-weight: 600; }

.method-badge {
  font-size: var(--fs-2xs); font-weight: 700; padding: 2px 7px;
  border-radius: var(--radius-xs); letter-spacing: 0.04em;
  text-transform: uppercase;
}
.method-badge.get    { background: var(--blue-soft); color: var(--method-get); }
.method-badge.post   { background: rgba(34,197,94,0.1); color: var(--method-post); }
.method-badge.put    { background: var(--amber-soft); color: var(--method-put); }
.method-badge.delete { background: var(--red-soft); color: var(--method-delete); }
.method-badge.patch  { background: var(--purple-soft); color: var(--method-patch); }

.hist-detail-section {
  margin-top: 16px;
  animation: slideUp 0.3s var(--ease-out) both;
}
.hist-detail-section h4 {
  margin: 0 0 8px; font-size: var(--fs-xs); color: var(--text-secondary);
  font-weight: 600; text-transform: uppercase; letter-spacing: 0.06em;
  font-family: var(--font-ui);
}
.hist-detail-pre {
  background: var(--bg-surface);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  padding: 12px 14px;
  font-size: var(--fs-sm);
  font-family: var(--font-mono);
  white-space: pre-wrap;
  word-break: break-all;
  max-height: 200px;
  overflow-y: auto;
  margin: 0;
  color: var(--text-primary);
  line-height: 1.7;
}
</style>
