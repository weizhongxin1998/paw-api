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
    />
    <HistoryPanel
      v-else-if="activePanel === 'history'"
      @open-tab="onHistoryReplay"
      ref="historyPanelRef"
    />

    <div v-if="activePanel === 'collection'" class="sidebar-footer">
      <button class="footer-btn" @click="onNewCollection">+ 新建集合</button>
      <button class="footer-btn" @click="onOpenDocs">Docs</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import CollectionTree from '../collection/CollectionTree.vue'
import HistoryPanel from '../history/HistoryPanel.vue'
import type { TreeItem } from '../../types/collection'

const activePanel = ref<'collection' | 'history'>('collection')
const tree = ref<TreeItem[]>([])
const historyPanelRef = ref<InstanceType<typeof HistoryPanel> | null>(null)

const emit = defineEmits<{
  (e: 'open-request', node: TreeItem): void
  (e: 'history-replay', item: any): void
  (e: 'panel-change', panel: 'collection' | 'history' | 'websocket'): void
  (e: 'open-docs'): void
}>()

function switchPanel(panel: 'collection' | 'history') {
  activePanel.value = panel
  emit('panel-change', panel as any)
}

function onOpenRequest(node: TreeItem) {
  emit('open-request', node)
}

function onContextMenu(_: TreeItem, _ev: MouseEvent) {}

function onNewCollection() {}

function onHistoryReplay(item: any) {
  emit('history-replay', item)
}

function onOpenDocs() {
  emit('open-docs')
}
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
  flex: 1;
  padding: 8px;
  text-align: center;
  font-size: 11px;
  cursor: pointer;
  background: #fafafa;
  border: none;
  color: #888;
  outline: none;
  border-bottom: 2px solid transparent;
}
.sidebar-tabs button.active {
  background: #fff;
  color: #18a058;
  border-bottom-color: #18a058;
  font-weight: 600;
}
.sidebar-tabs button:hover {
  color: #555;
}
.sidebar-footer {
  padding: 6px 10px;
  border-top: 1px solid #e8e8e8;
  display: flex;
  justify-content: space-between;
}
.footer-btn {
  background: none;
  border: none;
  color: #18a058;
  font-size: 11px;
  cursor: pointer;
  padding: 0;
}
.footer-btn:hover {
  text-decoration: underline;
}
</style>
