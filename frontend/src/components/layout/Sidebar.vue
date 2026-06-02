<template>
  <div class="sidebar">
    <div class="sidebar-tabs">
      <n-button
        :type="activePanel === 'collection' ? 'primary' : 'default'"
        size="small"
        @click="activePanel = 'collection'"
      >Collections</n-button>
      <n-button
        :type="activePanel === 'history' ? 'primary' : 'default'"
        size="small"
        @click="activePanel = 'history'"
      >History</n-button>
    </div>

    <CollectionTree
      v-if="activePanel === 'collection'"
      :tree="tree"
      @open-request="onOpenRequest"
      @ctx-menu="onContextMenu"
    />
    <HistoryPanel
      v-else
      @open-tab="onHistoryReplay"
      ref="historyPanelRef"
    />

    <div v-if="activePanel === 'collection'" class="sidebar-footer">
      <n-button text size="tiny" @click="onNewCollection">+ 新建集合</n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { NButton } from 'naive-ui'
import CollectionTree from '../collection/CollectionTree.vue'
import HistoryPanel from '../history/HistoryPanel.vue'
import type { TreeItem } from '../../types/collection'

const activePanel = ref<'collection' | 'history'>('collection')
const tree = ref<TreeItem[]>([])
const historyPanelRef = ref<InstanceType<typeof HistoryPanel> | null>(null)

const emit = defineEmits<{
  (e: 'open-request', node: TreeItem): void
  (e: 'history-replay', item: any): void
}>()

function onOpenRequest(node: TreeItem) {
  emit('open-request', node)
}

function onContextMenu(node: TreeItem, event: MouseEvent) {
  // TODO: show context menu
}

function onNewCollection() {
  // TODO: create collection
}

function onHistoryReplay(item: any) {
  emit('history-replay', item)
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
  padding: 8px;
  border-bottom: 1px solid #e8e8e8;
  display: flex;
  gap: 4px;
}
.sidebar-footer {
  padding: 6px 10px;
  border-top: 1px solid #e8e8e8;
}
</style>
