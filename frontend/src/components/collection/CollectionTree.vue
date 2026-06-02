<template>
  <div class="collection-tree">
    <div v-if="tree.length === 0" class="tree-empty">
      <n-empty description="暂无集合" size="small" />
    </div>
    <TreeNode
      v-for="node in tree"
      :key="node.id"
      :node="node"
      :depth="0"
      @dbl-click="onDblClick"
      @ctx-menu="onCtxMenu"
    />
  </div>
</template>

<script setup lang="ts">
import { NEmpty } from 'naive-ui'
import TreeNode from './TreeNode.vue'
import type { TreeItem } from '../../types/collection'

defineProps<{
  tree: TreeItem[]
}>()

const emit = defineEmits<{
  (e: 'open-request', node: TreeItem): void
  (e: 'dbl-click', node: TreeItem): void
  (e: 'ctx-menu', node: TreeItem, event: MouseEvent): void
}>()

function onDblClick(node: TreeItem) {
  if (node.type === 'request') {
    emit('open-request', node)
  } else {
    emit('dbl-click', node)
  }
}

function onCtxMenu(node: TreeItem, event: MouseEvent) {
  emit('ctx-menu', node, event)
}
</script>

<style scoped>
.collection-tree {
  flex: 1;
  overflow-y: auto;
  padding: 4px 0;
}
.tree-empty {
  padding: 24px;
  text-align: center;
}
</style>
