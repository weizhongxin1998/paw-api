<template>
  <div>
    <div
      class="tree-node"
      :class="node.type"
      :style="{ paddingLeft: 8 + depth * 16 + 'px' }"
      @contextmenu.prevent="onCtxMenu($event)"
    >
      <span v-if="node.type === 'folder'" class="arrow" @click.stop="expanded = !expanded">
        {{ expanded ? '\u25BC' : '\u25BA' }}
      </span>
      <span v-if="node.type === 'request'" class="method-tag" :class="node.method?.toLowerCase()">
        {{ node.method }}
      </span>
      <span class="node-name" @dblclick="onDblClick">{{ node.name }}</span>
      <span v-if="node.type === 'request'" class="node-url">{{ node.url }}</span>
    </div>
    <div v-if="node.type === 'folder' && expanded">
      <TreeNode
        v-for="child in node.children || []"
        :key="child.id"
        :node="child"
        :depth="depth + 1"
        @dbl-click="onChildDblClick"
        @ctx-menu="onChildCtxMenu"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { TreeItem } from '../../types/collection'

const props = defineProps<{
  node: TreeItem
  depth: number
}>()

const emit = defineEmits<{
  (e: 'dbl-click', node: TreeItem): void
  (e: 'ctx-menu', node: TreeItem, event: MouseEvent): void
}>()

const expanded = ref(true)

function onDblClick() {
  emit('dbl-click', props.node)
}

function onCtxMenu(event: MouseEvent) {
  emit('ctx-menu', props.node, event)
}

function onChildDblClick(node: TreeItem) {
  emit('dbl-click', node)
}

function onChildCtxMenu(node: TreeItem, event: MouseEvent) {
  emit('ctx-menu', node, event)
}
</script>

<style scoped>
.tree-node {
  display: flex;
  align-items: center;
  padding: 4px 8px;
  cursor: pointer;
  white-space: nowrap;
  font-size: 12px;
  gap: 4px;
  user-select: none;
}
.tree-node:hover {
  background: #f0f0f0;
}
.arrow {
  font-size: 9px;
  width: 12px;
  color: #888;
  cursor: pointer;
}
.method-tag {
  display: inline-block;
  width: 36px;
  font-size: 9px;
  font-weight: 700;
  text-align: center;
  padding: 0 2px;
  border-radius: 2px;
}
.method-tag.get { color: #18a058; }
.method-tag.post { color: #f0a020; }
.method-tag.put { color: #2080f0; }
.method-tag.delete { color: #d03050; }
.method-tag.patch { color: #9c27b0; }
.node-name {
  overflow: hidden;
  text-overflow: ellipsis;
}
.node-url {
  color: #aaa;
  font-size: 10px;
  margin-left: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
