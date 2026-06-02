<template>
  <div class="workspace">
    <n-empty v-if="!activeTab" description="点击左侧集合中的请求开始调试" class="empty-state" />
    <div v-else class="workspace-editor">
      <div class="tabs-bar">
        <div v-if="tabs.length === 0" class="no-tabs">无打开 Tab</div>
        <div v-for="tab in tabs" :key="tab.id" class="tab" :class="{ active: tab.id === activeTabId }" @click="activeTabId = tab.id">
          <span class="tab-method" :class="tab.method?.toLowerCase()">{{ tab.method }}</span>
          <span class="tab-name">{{ tab.name }}</span>
          <span v-if="tab.isDirty" class="tab-dirty"></span>
          <span class="tab-close" @click.stop="closeTab(tab.id)">x</span>
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
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { NEmpty } from 'naive-ui'
import UrlBar from '../request/UrlBar.vue'
import RequestPanel from '../request/RequestPanel.vue'
import ResponsePanel from '../response/ResponsePanel.vue'
import type { HttpResponse } from '../../types/response'

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
}

const activeTab = ref<Tab | null>(null)
const activeTabId = ref<string | null>(null)
const tabs = ref<Tab[]>([])
const response = ref<HttpResponse | null>(null)

const currentMethod = ref('GET')
const currentUrl = ref('')
const currentHeaders = ref('[]')
const currentParams = ref('[]')
const currentBodyType = ref('none')
const currentBodyData = ref('{}')
const currentAuthData = ref('{"type":"none"}')

let sessionCounter = 0

function onMethodChange(v: string) {
  currentMethod.value = v
  if (activeTab.value) { activeTab.value.method = v; activeTab.value.isDirty = true }
}

function onUrlChange(v: string) {
  currentUrl.value = v
  if (activeTab.value) { activeTab.value.url = v; activeTab.value.isDirty = true }
}

function onHeadersChange(v: string) { currentHeaders.value = v }
function onParamsChange(v: string) { currentParams.value = v }
function onBodyTypeChange(v: string) { currentBodyType.value = v }
function onBodyDataChange(v: string) { currentBodyData.value = v }
function onAuthDataChange(v: string) { currentAuthData.value = v }

function onSend() {
  response.value = {
    status: 200,
    time: 0,
    size: 0,
    headers: {},
    cookies: [],
    body: '',
    rawRequest: '',
    curlCommand: '',
  }
}

function closeTab(id: string) {
  tabs.value = tabs.value.filter(t => t.id !== id)
  if (activeTabId.value === id) {
    activeTabId.value = tabs.value.length > 0 ? tabs.value[0].id : null
    activeTab.value = tabs.value.length > 0 ? tabs.value[0] : null
  }
}
</script>

<style scoped>
.workspace {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fff;
  min-width: 0;
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
}
.tab.active {
  background: #fff;
  border-bottom: 2px solid #18a058;
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
