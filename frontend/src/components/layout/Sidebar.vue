<template>
  <div class="sidebar">
    <div class="sidebar-tabs">
      <button :class="{ active: activePanel === 'collection' }" @click="switchPanel('collection')">集合</button>
      <button :class="{ active: activePanel === 'history' }" @click="switchPanel('history')">历史</button>
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
  position: relative;
}
.sidebar-tabs {
  display: flex;
  border-bottom: 1px solid var(--border-primary);
  background: var(--bg-base);
}
.sidebar-tabs button {
  flex: 1; padding: 8px 6px; text-align: center; font-size: 10px; cursor: pointer;
  background: transparent; border: none; color: var(--text-muted); outline: none;
  font-weight: 600; letter-spacing: 0.5px; text-transform: uppercase;
  border-bottom: 2px solid transparent; transition: all var(--transition);
  font-family: var(--font-mono);
}
.sidebar-tabs button.active {
  color: var(--accent); border-bottom-color: var(--accent);
}
.sidebar-tabs button:hover { color: var(--text-secondary); background: rgba(255,255,255,0.02); }
.sidebar-add {
  display: flex; align-items: center; justify-content: center;
  width: 26px; font-size: 15px; color: var(--accent); cursor: pointer;
  font-weight: 700; border-radius: var(--radius-sm); margin: 3px 3px 3px 0;
  transition: all var(--transition);
}
.sidebar-add:hover { background: var(--accent-soft); color: var(--accent-hover); }
.add-menu {
  position: absolute; top: 36px; left: 6px; z-index: 50;
  background: var(--bg-elevated); border: 1px solid var(--border-primary);
  border-radius: var(--radius); box-shadow: 0 8px 24px rgba(0,0,0,0.4);
  padding: 3px 0; min-width: 140px;
}
.add-menu-item {
  padding: 6px 14px; cursor: pointer; font-size: 11px; color: var(--text-secondary);
  transition: background var(--transition); font-family: var(--font-mono);
}
.add-menu-item:hover { background: var(--bg-hover); color: var(--text-primary); }

.modal-overlay {
  position: fixed; inset: 0; background: rgba(0,0,0,0.6);
  display: flex; align-items: center; justify-content: center; z-index: 200;
}
.modal-box {
  background: var(--bg-surface); border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg); padding: 22px; width: 340px;
  box-shadow: 0 12px 40px rgba(0,0,0,0.5);
}
.modal-box h3 { margin: 0 0 16px; font-size: 14px; font-weight: 600; color: var(--text-primary); }
.modal-box label { display: block; font-size: 10px; color: var(--text-muted); margin-bottom: 3px; letter-spacing: 0.5px; text-transform: uppercase; }
.modal-box input, .modal-box select {
  width: 100%; padding: 7px 10px; border: 1px solid var(--border-primary); border-radius: var(--radius);
  font-size: 12px; outline: none; margin-bottom: 10px; box-sizing: border-box;
  background: var(--bg-base); color: var(--text-primary); font-family: var(--font-mono);
  transition: border-color var(--transition);
}
.modal-box input:focus, .modal-box select:focus { border-color: var(--accent); }
.method-select { cursor: pointer; }
.modal-acts { display: flex; justify-content: flex-end; gap: 6px; margin-top: 16px; }
.modal-acts button {
  padding: 6px 16px; border: 1px solid var(--border-primary); border-radius: var(--radius);
  font-size: 11px; cursor: pointer; background: var(--bg-base); color: var(--text-secondary);
  transition: all var(--transition); font-family: var(--font-mono);
}
.modal-acts button:hover { border-color: var(--border-hover); color: var(--text-primary); }
.modal-acts .btn-primary { background: var(--accent); color: #000; border-color: var(--accent); font-weight: 600; }
.modal-acts .btn-primary:hover { background: var(--accent-hover); border-color: var(--accent-hover); }
.modal-acts .btn-primary:disabled { background: var(--bg-elevated); border-color: var(--border-primary); color: var(--text-muted); cursor: not-allowed; }
</style>
