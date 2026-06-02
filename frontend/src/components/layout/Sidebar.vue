<template>
  <div class="sidebar">
    <div class="sidebar-tabs">
      <button :class="{ active: activePanel === 'collection' }" @click="switchPanel('collection')">Collections</button>
      <button :class="{ active: activePanel === 'history' }" @click="switchPanel('history')">History</button>
    </div>

    <CollectionTree
      v-if="activePanel === 'collection'"
      :tree="tree"
      @open-request="onOpenRequest"
      @ctx-menu="onContextMenu"
      @action="onTreeAction"
    />
    <HistoryPanel
      v-else-if="activePanel === 'history'"
      @open-tab="onHistoryReplay"
      ref="historyPanelRef"
    />

    <div v-if="activePanel === 'collection'" class="sidebar-footer">
      <button class="footer-btn" @click="showRequestModal = true">+ 新建请求</button>
      <button class="footer-btn" @click="showCollectionModal = true">+ 新建集合</button>
      <button class="footer-btn" @click="onOpenDocs">Docs</button>
    </div>

    <div v-if="showCollectionModal" class="modal-overlay" @click.self="showCollectionModal = false">
      <div class="modal-box">
        <h3>新建集合</h3>
        <label>名称</label>
        <input v-model="collectionName" placeholder="集合名称" @keydown.enter="onCreateCollection" />
        <div class="modal-acts">
          <button @click="showCollectionModal = false">取消</button>
          <button class="btn-primary" @click="onCreateCollection" :disabled="!collectionName.trim()">创建</button>
        </div>
      </div>
    </div>

    <div v-if="showRequestModal" class="modal-overlay" @click.self="showRequestModal = false">
      <div class="modal-box">
        <h3>新建请求</h3>
        <label>名称</label>
        <input v-model="requestName" placeholder="请求名称" @keydown.enter="onCreateRequest" />
        <label>方法</label>
        <select v-model="requestMethod" class="method-select">
          <option>GET</option><option>POST</option><option>PUT</option><option>DELETE</option>
          <option>PATCH</option><option>HEAD</option><option>OPTIONS</option>
        </select>
        <div class="modal-acts">
          <button @click="showRequestModal = false">取消</button>
          <button class="btn-primary" @click="onCreateRequest" :disabled="!requestName.trim()">创建</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useCollectionStore } from '../../stores/collection'
import CollectionTree from '../collection/CollectionTree.vue'
import HistoryPanel from '../history/HistoryPanel.vue'
import type { TreeItem } from '../../types/collection'

const props = defineProps<{
  projectId: number | null
}>()

const activePanel = ref<'collection' | 'history'>('collection')
const historyPanelRef = ref<InstanceType<typeof HistoryPanel> | null>(null)
const collectionStore = useCollectionStore()

const emit = defineEmits<{
  (e: 'open-request', node: TreeItem): void
  (e: 'history-replay', item: any): void
  (e: 'open-docs'): void
  (e: 'tree-action', action: string, node: TreeItem): void
}>()

const tree = ref<TreeItem[]>([])
const showCollectionModal = ref(false)
const showRequestModal = ref(false)
const collectionName = ref('')
const requestName = ref('')
const requestMethod = ref('GET')

function switchPanel(panel: 'collection' | 'history') {
  activePanel.value = panel
}

watch(() => props.projectId, async (id) => {
  if (id) {
    tree.value = await collectionStore.loadTree(id)
  } else {
    tree.value = []
  }
}, { immediate: true })

async function refreshTree() {
  if (props.projectId) {
    tree.value = await collectionStore.loadTree(props.projectId)
  }
}

function onOpenRequest(node: TreeItem) {
  emit('open-request', node)
}

function onContextMenu(_: TreeItem, _ev: MouseEvent) {}

function onTreeAction(action: string, node: TreeItem) {
  emit('tree-action', action, node)
}

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

function onHistoryReplay(item: any) {
  emit('history-replay', item)
}

function onOpenDocs() {
  emit('open-docs')
}

defineExpose({ refreshTree })
</script>

<style scoped>
.sidebar {
  width: 240px;
  background: #fafafa;
  border-right: 1px solid #e0e0e0;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}
.sidebar-tabs {
  display: flex;
  border-bottom: 1px solid #e8e8e8;
}
.sidebar-tabs button {
  flex: 1; padding: 8px; text-align: center; font-size: 11px; cursor: pointer;
  background: #fafafa; border: none; color: #888; outline: none;
  border-bottom: 2px solid transparent;
}
.sidebar-tabs button.active {
  background: #fff; color: #18a058; border-bottom-color: #18a058; font-weight: 600;
}
.sidebar-footer {
  padding: 6px 10px; border-top: 1px solid #e8e8e8;
  display: flex; justify-content: flex-start; gap: 12px;
}
.footer-btn {
  background: none; border: none; color: #18a058;
  font-size: 11px; cursor: pointer; padding: 0;
}
.footer-btn:hover { text-decoration: underline; }

.modal-overlay {
  position: fixed; inset: 0; background: rgba(0,0,0,0.3);
  display: flex; align-items: center; justify-content: center; z-index: 200;
}
.modal-box {
  background: #fff; border-radius: 10px; padding: 24px; width: 360px;
  box-shadow: 0 8px 30px rgba(0,0,0,0.15);
}
.modal-box h3 { margin: 0 0 16px; font-size: 15px; }
.modal-box label { display: block; font-size: 12px; color: #888; margin-bottom: 4px; }
.modal-box input, .modal-box select {
  width: 100%; padding: 7px 10px; border: 1px solid #ddd; border-radius: 6px;
  font-size: 13px; outline: none; margin-bottom: 10px; box-sizing: border-box;
}
.modal-box input:focus, .modal-box select:focus { border-color: #18a058; }
.method-select { background: #fff; cursor: pointer; }
.modal-acts { display: flex; justify-content: flex-end; gap: 8px; margin-top: 14px; }
.modal-acts button {
  padding: 6px 18px; border: 1px solid #ddd; border-radius: 6px;
  font-size: 12px; cursor: pointer; background: #fff;
}
.modal-acts .btn-primary { background: #18a058; color: #fff; border-color: #18a058; }
.modal-acts .btn-primary:hover { background: #0c7a43; }
.modal-acts .btn-primary:disabled { background: #aaa; border-color: #aaa; cursor: not-allowed; }
</style>
