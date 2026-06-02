<template>
  <div class="sidebar">
    <div class="sidebar-tabs">
      <n-button type="primary" size="small" block>Collections</n-button>
    </div>
    <CollectionTree
      :tree="tree"
      @open-request="onOpenRequest"
      @ctx-menu="onContextMenu"
    />
    <div class="sidebar-footer">
      <n-button text size="tiny" @click="onNewCollection">+ 新建集合</n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { NButton } from 'naive-ui'
import CollectionTree from '../collection/CollectionTree.vue'
import type { TreeItem } from '../../types/collection'

const tree = ref<TreeItem[]>([])

function onOpenRequest(node: TreeItem) {
  // TODO: open tab
}

function onContextMenu(node: TreeItem, event: MouseEvent) {
  // TODO: show context menu
}

function onNewCollection() {
  // TODO: create collection
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
}
.sidebar-footer {
  padding: 6px 10px;
  border-top: 1px solid #e8e8e8;
}
</style>
