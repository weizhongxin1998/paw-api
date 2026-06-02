<template>
  <div class="body">
    <Sidebar
      @open-request="onOpenFromCollection"
      @history-replay="onHistoryReplay"
    />
    <Workspace ref="workspaceRef" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import Sidebar from './Sidebar.vue'
import Workspace from './Workspace.vue'
import { GetRequest } from '../../../wailsjs/go/main/App'
import { models } from '../../../wailsjs/go/models'
import type { TreeItem } from '../../types/collection'

const workspaceRef = ref<InstanceType<typeof Workspace> | null>(null)

let tabCounter = 0
function newTabId(): string {
  return 'tab-' + (++tabCounter) + '-' + Date.now()
}

async function onOpenFromCollection(node: TreeItem) {
  if (!workspaceRef.value) return
  let request: models.Request | null = null
  try {
    request = await GetRequest(node.id)
  } catch {
    // ignore
  }
  const tab = {
    id: newTabId(),
    requestId: node.id,
    method: request?.method || node.method || 'GET',
    name: node.name || '',
    url: request?.url || node.url || '',
    isDirty: false,
    headers: request?.headers || '[]',
    params: request?.params || '[]',
    bodyType: request?.body_type || 'none',
    bodyData: request?.body || '{}',
    authData: request?.auth || '{"type":"none"}',
  }
  workspaceRef.value.openTab(tab)
}

async function onHistoryReplay(item: any) {
  if (!workspaceRef.value) return
  let request: models.Request | null = null
  if (item.request_id) {
    try {
      request = await GetRequest(item.request_id)
    } catch {
      // ignore
    }
  }
  const tab = {
    id: newTabId(),
    requestId: request?.id || item.request_id || 0,
    method: request?.method || item.method || 'GET',
    name: 'Replay: ' + (request?.name || item.url || ''),
    url: request?.url || item.url || '',
    isDirty: false,
    headers: request?.headers || item.request_headers || '[]',
    params: request?.params || '[]',
    bodyType: request?.body_type || 'none',
    bodyData: request?.body || item.request_body || '{}',
    authData: request?.auth || '{"type":"none"}',
  }
  workspaceRef.value.openTab(tab)
}
</script>

<style scoped>
.body {
  display: flex;
  flex: 1;
  overflow: hidden;
}
</style>
