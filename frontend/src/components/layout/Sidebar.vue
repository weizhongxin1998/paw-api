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
      <button class="footer-btn" @click="onNewRequest">+ 新建请求</button>
      <button class="footer-btn" @click="onNewCollection">+ 新建集合</button>
      <button class="footer-btn" @click="onOpenDocs">Docs</button>
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

async function onNewCollection() {
  if (!props.projectId) return
  const name = prompt('集合名称:')
  if (!name) return
  await collectionStore.createCollection(props.projectId, null, name)
  await refreshTree()
}

async function onNewRequest() {
  if (!props.projectId) return
  const name = prompt('请求名称:')
  if (!name) return
  const method = prompt('请求方法 (GET/POST/PUT/DELETE):', 'GET')
  if (!method) return
  if (tree.value.length > 0) {
    await collectionStore.createRequest(tree.value[0].id, name, method)
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
</style>
