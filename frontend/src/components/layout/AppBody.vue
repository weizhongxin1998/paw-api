<template>
  <div class="body">
    <Sidebar
      ref="sidebarRef"
      :project-id="projectId"
      @open-request="onOpenFromCollection"
      @history-replay="onHistoryReplay"
      @open-docs="docsModalShow = true"
      @tree-action="onTreeAction"
    />
    <Workspace ref="workspaceRef" :project-id="projectId" />
    <DocsPreviewModal
      v-model:show="docsModalShow"
      :project-id="currentProjectId"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import Sidebar from './Sidebar.vue'
import Workspace from './Workspace.vue'
import DocsPreviewModal from '../modals/DocsPreviewModal.vue'
import { useCollectionStore } from '../../stores/collection'
import { GetRequest } from '../../../wailsjs/go/main/App'
import { models } from '../../../wailsjs/go/models'
import type { TreeItem } from '../../types/collection'

const props = defineProps<{
  projectId: number | null
}>()

const workspaceRef = ref<InstanceType<typeof Workspace> | null>(null)
const sidebarRef = ref<InstanceType<typeof Sidebar> | null>(null)
const docsModalShow = ref(false)
const currentProjectId = computed(() => props.projectId)
const collectionStore = useCollectionStore()

function onPanelChange(_: string) {}

async function onTreeAction(action: string, node: TreeItem) {
  if (!props.projectId) return
  switch (action) {
    case 'rename': {
      const name = prompt('新名称:', node.name)
      if (name && name !== node.name) {
        await collectionStore.renameCollection(node.id, name)
        await sidebarRef.value?.refreshTree()
      }
      break
    }
    case 'delete': {
      if (confirm(`确定删除 "${node.name}"?`)) {
        if (node.type === 'folder')
          await collectionStore.removeCollection(node.id)
        else
          await collectionStore.removeRequest(node.id)
        await sidebarRef.value?.refreshTree()
      }
      break
    }
    case 'duplicate': {
      await collectionStore.duplicateRequest(node.id)
      await sidebarRef.value?.refreshTree()
      break
    }
    case 'new-folder': {
      // delegate to sidebar
      break
    }
    case 'new-request': {
      // delegate to sidebar
      break
    }
  }
}

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
    collectionId: request?.collection_id || 0,
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
    collectionId: request?.collection_id || 0,
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
