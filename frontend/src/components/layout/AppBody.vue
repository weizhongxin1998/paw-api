<template>
  <div class="body">
    <Sidebar
      ref="sidebarRef"
      :project-id="projectId"
      @open-request="onPreviewFromCollection"
      @open-request-persist="onOpenFromCollection"
      @history-replay="onHistoryReplay"
      @history-select="onHistorySelect"
      @open-docs="docsModalShow = true"
      @tree-action="onTreeAction"
    />
    <Workspace ref="workspaceRef" :project-id="projectId" @request-saved="onRequestSaved" />
    <DocsPreviewModal
      v-model:show="docsModalShow"
      :project-id="currentProjectId"
    />

    <n-modal v-model:show="showRenameModal" preset="card" :title="t('appBody.rename')" :class="modalClass" style="width: 360px" :mask-closable="false">
      <n-form label-placement="top">
        <n-form-item :label="t('common.newName')">
          <n-input v-model:value="renameValue" :placeholder="t('common.newNamePlaceholder')" @keydown.enter="onRenameConfirm" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button @click="showRenameModal = false">{{ t('common.cancel') }}</n-button>
        <n-button type="primary" :disabled="!renameValue.trim()" @click="onRenameConfirm">{{ t('common.confirm') }}</n-button>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useDialog, useMessage, NModal, NForm, NFormItem, NInput, NButton } from 'naive-ui'
import Sidebar from './Sidebar.vue'
import Workspace from './Workspace.vue'
import DocsPreviewModal from '../modals/DocsPreviewModal.vue'
import { useCollectionStore } from '../../stores/collection'
import { GetRequest, UpdateRequest } from '../../../wailsjs/go/main/App'
import { models } from '../../../wailsjs/go/models'
import type { TreeItem } from '../../types/collection'

const { t } = useI18n()

const props = defineProps<{
  projectId: number | null
}>()

const workspaceRef = ref<InstanceType<typeof Workspace> | null>(null)
const sidebarRef = ref<InstanceType<typeof Sidebar> | null>(null)
const docsModalShow = ref(false)
const currentProjectId = computed(() => props.projectId)
const collectionStore = useCollectionStore()
const dialog = useDialog()
const message = useMessage()

const showRenameModal = ref(false)
const renameValue = ref('')
let renameNodeId = 0
let renameNodeType = ''

const isLightMode = ref(false)
const modalClass = computed(() => isLightMode.value ? 'rename-modal theme-light' : 'rename-modal')

onMounted(() => {
  const check = () => { isLightMode.value = !!document.querySelector('.theme-light') }
  check()
  const observer = new MutationObserver(check)
  observer.observe(document.documentElement, { attributes: true, attributeFilter: ['class'], subtree: true })
})

function onRequestSaved() {
  sidebarRef.value?.refreshTree()
}

function onHistorySelect(item: any) {
  if (!workspaceRef.value) return
  if (item) {
    workspaceRef.value.showHistoryDetail(item)
  } else {
    workspaceRef.value.clearHistoryDetail()
  }
}

async function onTreeAction(action: string, node: TreeItem) {
  if (!props.projectId) return
  switch (action) {
    case 'rename': {
      renameNodeId = node.id
      renameNodeType = node.type
      renameValue.value = node.name || ''
      showRenameModal.value = true
      break
    }
    case 'delete': {
      dialog.warning({
        title: t('appBody.confirmDeleteTitle'),
        content: t('appBody.confirmDeleteContent', { name: node.name }),
        positiveText: t('common.delete'),
        negativeText: t('common.cancel'),
        onPositiveClick: async () => {
          try {
            if (node.type === 'folder' || node.type === 'root')
              await collectionStore.removeCollection(node.id)
            else
              await collectionStore.removeRequest(node.id)
            await sidebarRef.value?.refreshTree()
            message.success(t('appBody.deleted', { name: node.name }))
          } catch (e: any) {
            message.error(t('appBody.deleteFailed', { error: e?.message || String(e) }))
          }
        },
      })
      break
    }
    case 'duplicate': {
      try {
        await collectionStore.duplicateRequest(node.id)
        await sidebarRef.value?.refreshTree()
        message.success(t('appBody.duplicated', { name: node.name }))
      } catch (e: any) {
        message.error(t('appBody.duplicateFailed', { error: e?.message || String(e) }))
      }
      break
    }
    case 'new-request': {
      const parentId = (node.type === 'folder' || node.type === 'root') ? node.id : 0
      if (parentId) {
        try {
          await collectionStore.createRequest(parentId, t('appBody.newRequest'), 'GET')
          await sidebarRef.value?.refreshTree()
          message.success(t('appBody.newRequestCreated'))
        } catch (e: any) {
          message.error(t('appBody.createFailed', { error: e?.message || String(e) }))
        }
      }
      break
    }
    case 'new-folder': {
      const parentId = (node.type === 'folder' || node.type === 'root') ? node.id : 0
      if (parentId) {
        try {
          await collectionStore.createCollection(props.projectId, parentId, t('appBody.newFolder'))
          await sidebarRef.value?.refreshTree()
          message.success(t('appBody.newFolderCreated'))
        } catch (e: any) {
          message.error(t('appBody.createFailed', { error: e?.message || String(e) }))
        }
      }
      break
    }
  }
}

async function onRenameConfirm() {
  const name = renameValue.value.trim()
  if (!name) return
  try {
    if (renameNodeType === 'request') {
      // For request nodes, fetch current data then update name via UpdateRequest
      const req = await GetRequest(renameNodeId)
      await UpdateRequest({
        ...req,
        name,
      } as any)
    } else {
      await collectionStore.renameCollection(renameNodeId, name)
    }
    await sidebarRef.value?.refreshTree()
    if (renameNodeType === 'request') {
      workspaceRef.value?.renameOpenTab(renameNodeId, name)
    }
    showRenameModal.value = false
    message.success(t('appBody.renamed'))
  } catch (e: any) {
    message.error(t('appBody.renameFailed', { error: e?.message || String(e) }))
  }
}

let tabCounter = 0
function newTabId(): string {
  return 'tab-' + (++tabCounter) + '-' + Date.now()
}

async function onPreviewFromCollection(node: TreeItem) {
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
    pathVars: '[]',
    isDirty: false,
    isPreview: true,
    headers: request?.headers || '[]',
    params: request?.params || '[]',
    paramsEnabled: true,
    bodyType: request?.body_type || 'none',
    bodyData: request?.body || '{}',
    authData: request?.auth || '{"type":"none"}',
    collectionId: request?.collection_id || 0,
  }
  workspaceRef.value.previewTab(tab)
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
    pathVars: '[]',
    isDirty: false,
    isPreview: false,
    headers: request?.headers || '[]',
    params: request?.params || '[]',
    paramsEnabled: true,
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
    name: t('appBody.replayPrefix', { name: request?.name || item.url || '' }),
    url: request?.url || item.url || '',
    pathVars: '[]',
    isDirty: false,
    isPreview: false,
    headers: request?.headers || item.request_headers || '[]',
    params: request?.params || '[]',
    paramsEnabled: true,
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
