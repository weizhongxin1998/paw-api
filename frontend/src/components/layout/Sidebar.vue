<template>
  <div class="sidebar">
    <div class="sidebar-tabs">
      <button :class="{ active: activePanel === 'collection' }" @click="switchPanel('collection')">Collections</button>
      <button :class="{ active: activePanel === 'history' }" @click="switchPanel('history')">History</button>
      <span class="sidebar-add" @click="showAddMenu = !showAddMenu">+</span>
    </div>
    <div v-if="showAddMenu" class="add-menu" @click.self="showAddMenu = false">
      <div class="add-menu-item" @click="showRequestModal = true; showAddMenu = false">新建请求</div>
      <div class="add-menu-item" @click="showCollectionModal = true; showAddMenu = false">新建集合</div>
      <div class="add-menu-item" @click="$emit('open-docs'); showAddMenu = false">API 文档</div>
    </div>

    <CollectionTree
      v-if="activePanel === 'collection'"
      :tree="tree"
      @open-request="onPreviewRequest"
      @dbl-click="onPersistRequest"
      @ctx-menu="onContextMenu"
      @action="onTreeAction"
    />
    <HistoryPanel
      v-else-if="activePanel === 'history'"
      @open-tab="onHistoryReplay"
      @select-detail="onHistorySelect"
      ref="historyPanelRef"
    />

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
  (e: 'open-request-persist', node: TreeItem): void
  (e: 'history-replay', item: any): void
  (e: 'history-select', item: any | null): void
  (e: 'open-docs'): void
  (e: 'tree-action', action: string, node: TreeItem): void
}>()

const tree = ref<TreeItem[]>([])
const showCollectionModal = ref(false)
const showRequestModal = ref(false)
const showAddMenu = ref(false)
const collectionName = ref('')
const requestName = ref('')
const requestMethod = ref('GET')

function switchPanel(panel: 'collection' | 'history') {
  activePanel.value = panel
  if (panel === 'collection') {
    emit('history-select', null)
  }
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

function onPreviewRequest(node: TreeItem) {
  emit('open-request', node)
}

function onPersistRequest(node: TreeItem) {
  if (node.type === 'request') {
    emit('open-request-persist', node)
  }
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

function onHistorySelect(item: any) {
  emit('history-select', item)
}

defineExpose({ refreshTree })
</script>

<style scoped>
.sidebar {
  width: 220px;
  background: #fff;
  border-right: 1px solid var(--gray-200);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  position: relative;
}
.sidebar-tabs {
  display: flex;
  border-bottom: 1px solid var(--gray-200);
  background: var(--gray-50);
}
.sidebar-tabs button {
  flex: 1; padding: 9px 8px; text-align: center; font-size: 11px; cursor: pointer;
  background: transparent; border: none; color: var(--gray-400); outline: none;
  font-weight: 500; letter-spacing: 0.4px; text-transform: uppercase;
  border-bottom: 2px solid transparent; transition: all var(--transition);
}
.sidebar-tabs button.active {
  color: var(--green); border-bottom-color: var(--green); font-weight: 600;
}
.sidebar-tabs button:hover {
  color: var(--gray-600); background: rgba(0,0,0,0.02);
}
.sidebar-add {
  display: flex; align-items: center; justify-content: center;
  width: 28px; font-size: 17px; color: var(--green); cursor: pointer;
  font-weight: 700; border-radius: var(--radius-sm); margin: 4px 4px 4px 0;
  transition: all var(--transition);
}
.sidebar-add:hover { background: var(--green-soft); color: var(--green-hover); }
.add-menu {
  position: absolute; top: 38px; left: 6px; z-index: 50;
  background: #fff; border: 1px solid var(--gray-200);
  border-radius: var(--radius-lg); box-shadow: var(--shadow-lg);
  padding: 4px 0; min-width: 140px;
}
.add-menu-item {
  padding: 7px 14px; cursor: pointer; font-size: 12px; color: var(--gray-600);
  transition: background var(--transition);
}
.add-menu-item:hover { background: var(--gray-100); }

.modal-overlay {
  position: fixed; inset: 0; background: rgba(0,0,0,0.3);
  display: flex; align-items: center; justify-content: center; z-index: 200;
}
.modal-box {
  background: #fff; border-radius: var(--radius-lg); padding: 26px; width: 360px;
  box-shadow: var(--shadow-lg);
}
.modal-box h3 { margin: 0 0 18px; font-size: 15px; font-weight: 600; }
.modal-box label { display: block; font-size: 11px; color: var(--gray-500); margin-bottom: 4px; text-transform: uppercase; letter-spacing: 0.3px; }
.modal-box input, .modal-box select {
  width: 100%; padding: 8px 10px; border: 1px solid var(--gray-200); border-radius: var(--radius);
  font-size: 13px; outline: none; margin-bottom: 12px; box-sizing: border-box;
  transition: border-color var(--transition);
}
.modal-box input:focus, .modal-box select:focus { border-color: var(--green); box-shadow: 0 0 0 3px rgba(24,160,88,0.1); }
.method-select { background: #fff; cursor: pointer; }
.modal-acts { display: flex; justify-content: flex-end; gap: 8px; margin-top: 16px; }
.modal-acts button {
  padding: 7px 18px; border: 1px solid var(--gray-200); border-radius: var(--radius);
  font-size: 12px; cursor: pointer; background: #fff; color: var(--gray-600);
  transition: all var(--transition);
}
.modal-acts button:hover { border-color: var(--gray-300); }
.modal-acts .btn-primary { background: var(--green); color: #fff; border-color: var(--green); }
.modal-acts .btn-primary:hover { background: var(--green-hover); border-color: var(--green-hover); }
.modal-acts .btn-primary:disabled { background: var(--gray-300); border-color: var(--gray-300); cursor: not-allowed; }
</style>
