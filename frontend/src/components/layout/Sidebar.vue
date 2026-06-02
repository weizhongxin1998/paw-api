<template>
  <div class="sidebar">
    <div class="sidebar-tabs">
      <n-tabs
        v-model:value="activePanel"
        type="bar"
        size="small"
        animated
      >
        <n-tab-pane name="collection" tab="集合" />
        <n-tab-pane name="history" tab="历史" />
      </n-tabs>
      <n-dropdown trigger="click" :options="addMenuOptions" @select="onAddMenuSelect">
        <n-button text size="tiny" class="add-btn">+</n-button>
      </n-dropdown>
    </div>

    <div class="sidebar-body">
      <CollectionTree
        v-if="activePanel === 'collection'"
        :tree="tree"
        @open-request="onPreviewRequest"
        @dbl-click="onPersistRequest"
        @ctx-menu="onContextMenu"
        @action="onTreeAction"
      />
      <HistoryPanel
        v-else
        @open-tab="onHistoryReplay"
        @select-detail="onHistorySelect"
        ref="historyPanelRef"
      />
    </div>

    <n-modal v-model:show="showCollectionModal" preset="card" title="新建集合" style="width: 360px" :mask-closable="false">
      <n-form label-placement="top">
        <n-form-item label="名称">
          <n-input v-model:value="collectionName" placeholder="集合名称" @keydown.enter="onCreateCollection" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button @click="showCollectionModal = false">取消</n-button>
        <n-button type="primary" :disabled="!collectionName.trim()" @click="onCreateCollection">创建</n-button>
      </template>
    </n-modal>

    <n-modal v-model:show="showRequestModal" preset="card" title="新建请求" style="width: 360px" :mask-closable="false">
      <n-form label-placement="top">
        <n-form-item label="名称">
          <n-input v-model:value="requestName" placeholder="请求名称" @keydown.enter="onCreateRequest" />
        </n-form-item>
        <n-form-item label="方法">
          <n-select v-model:value="requestMethod" :options="methodOptions" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button @click="showRequestModal = false">取消</n-button>
        <n-button type="primary" :disabled="!requestName.trim()" @click="onCreateRequest">创建</n-button>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { NTabs, NTabPane, NButton, NDropdown, NModal, NForm, NFormItem, NInput, NSelect } from 'naive-ui'
import type { DropdownOption } from 'naive-ui'
import { useCollectionStore } from '../../stores/collection'
import CollectionTree from '../collection/CollectionTree.vue'
import HistoryPanel from '../history/HistoryPanel.vue'
import type { TreeItem } from '../../types/collection'

const props = defineProps<{
  projectId: number | null
}>()

const activePanel = ref('collection')
const historyPanelRef = ref<InstanceType<typeof HistoryPanel> | null>(null)
const collectionStore = useCollectionStore()

const emit = defineEmits<{
  (e: 'open-request', node: TreeItem): void
  (e: 'open-request-persist', node: TreeItem): void
  (e: 'history-replay', item: any): void
  (e: 'history-select', item: any | null): void
  (e: 'open-docs'): void
  (e: 'tree-action', action: string, node: TreeItem): void
}>()

const tree = ref<TreeItem[]>([])
const showCollectionModal = ref(false)
const showRequestModal = ref(false)
const collectionName = ref('')
const requestName = ref('')
const requestMethod = ref('GET')

const methodOptions = [
  { label: 'GET', value: 'GET' }, { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' }, { label: 'DELETE', value: 'DELETE' },
  { label: 'PATCH', value: 'PATCH' }, { label: 'HEAD', value: 'HEAD' },
  { label: 'OPTIONS', value: 'OPTIONS' },
]

const addMenuOptions: DropdownOption[] = [
  { label: '新建请求', key: 'request' },
  { label: '新建集合', key: 'collection' },
  { label: 'API 文档', key: 'docs' },
]

function onAddMenuSelect(key: string) {
  switch (key) {
    case 'request': showRequestModal.value = true; break
    case 'collection': showCollectionModal.value = true; break
    case 'docs': emit('open-docs'); break
  }
}

watch(activePanel, (val) => {
  if (val === 'collection') emit('history-select', null)
})

watch(() => props.projectId, async (id) => {
  if (id) tree.value = await collectionStore.loadTree(id)
  else tree.value = []
}, { immediate: true })

async function refreshTree() {
  if (props.projectId) tree.value = await collectionStore.loadTree(props.projectId)
}

function onPreviewRequest(node: TreeItem) { emit('open-request', node) }
function onPersistRequest(node: TreeItem) {
  if (node.type === 'request') emit('open-request-persist', node)
}
function onContextMenu(_: TreeItem, _ev: MouseEvent) {}
function onTreeAction(action: string, node: TreeItem) { emit('tree-action', action, node) }

async function onCreateCollection() {
  const name = collectionName.value.trim()
  if (!name || !props.projectId) return
  await collectionStore.createCollection(props.projectId, null, name)
  showCollectionModal.value = false
  collectionName.value = ''
  await refreshTree()
}

async function onCreateRequest() {
  const name = requestName.value.trim()
  if (!name || !props.projectId) return
  if (tree.value.length > 0) {
    await collectionStore.createRequest(tree.value[0].id, name, requestMethod.value)
    showRequestModal.value = false
    requestName.value = ''
    requestMethod.value = 'GET'
    await refreshTree()
  }
}

function onHistoryReplay(item: any) { emit('history-replay', item) }
function onHistorySelect(item: any) { emit('history-select', item) }

defineExpose({ refreshTree })
</script>

<style scoped>
.sidebar {
  width: 220px;
  background: var(--bg-surface);
  border-right: 1px solid var(--border-primary);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}
.sidebar-tabs {
  display: flex;
  align-items: center;
}
.sidebar-tabs :deep(.n-tabs) {
  flex: 1;
}
.sidebar-tabs :deep(.n-tabs .n-tabs-tab__label) {
  font-size: var(--fs-sm) !important;
}
.sidebar-tabs :deep(.n-tabs .n-tabs-tab) {
  padding: 6px 10px;
}
.add-btn {
  font-size: 16px;
  font-weight: 700;
  padding: 1px 5px;
  margin-right: 2px;
  color: var(--accent);
  line-height: 1;
}
.sidebar-body {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}
</style>
