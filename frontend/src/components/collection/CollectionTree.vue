<template>
  <div class="collection-tree">
    <!-- Empty state with helpful illustration -->
    <div v-if="filteredTree.length === 0" class="tree-empty">
      <template v-if="searchFilter && searchFilter.trim()">
        <!-- No results for search -->
        <div class="empty-illustr">
          <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.2" stroke-linecap="round" stroke-linejoin="round" class="empty-icon">
            <circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/>
            <line x1="8" y1="11" x2="14" y2="11"/>
          </svg>
        </div>
        <p class="empty-text">未找到匹配 "{{ searchFilter }}" 的请求</p>
        <p class="empty-hint">尝试使用不同的关键词</p>
      </template>
      <template v-else>
        <!-- No collections at all -->
        <div class="empty-illustr">
          <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round" class="empty-icon">
            <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
            <line x1="12" y1="11" x2="12" y2="17"/>
            <line x1="9" y1="14" x2="15" y2="14"/>
          </svg>
        </div>
        <p class="empty-text">创建第一个集合</p>
        <p class="empty-hint">点击右上角 + 按钮开始组织你的 API 请求</p>
      </template>
    </div>

    <TreeNode
      v-for="node in filteredTree"
      :key="node.id"
      :node="node"
      :depth="0"
      @click="onClick"
      @dbl-click="onDblClick"
      @ctx-menu="onCtxMenu"
      @action="onAction"
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import TreeNode from './TreeNode.vue'
import type { TreeItem } from '../../types/collection'

const props = defineProps<{
  tree: TreeItem[]
  /** Optional search filter string to filter requests by name */
  searchFilter?: string
}>()

const emit = defineEmits<{
  (e: 'open-request', node: TreeItem): void
  (e: 'dbl-click', node: TreeItem): void
  (e: 'ctx-menu', node: TreeItem, event: MouseEvent): void
  (e: 'action', action: string, node: TreeItem): void
}>()

/**
 * Recursively filter the tree to only include nodes that match the search query.
 * Folders are kept if they have matching children.
 * Root nodes are always kept (they act as collection headers).
 */
function filterTree(nodes: TreeItem[], query: string): TreeItem[] {
  const q = query.toLowerCase().trim()
  if (!q) return nodes

  return nodes.reduce<TreeItem[]>((acc, node) => {
    if (node.type === 'root') {
      // Always show root, but filter its children
      const filteredChildren = filterTree(node.children || [], query)
      acc.push({ ...node, children: filteredChildren })
      return acc
    }

    if (node.type === 'request') {
      // Match by name or method
      const nameMatch = node.name.toLowerCase().includes(q)
      const methodMatch = node.method?.toLowerCase().includes(q)
      if (nameMatch || methodMatch) {
        acc.push(node)
      }
      return acc
    }

    if (node.type === 'folder') {
      // Keep folder if its name matches or any children match
      const filteredChildren = filterTree(node.children || [], query)
      const nameMatch = node.name.toLowerCase().includes(q)
      if (nameMatch || filteredChildren.length > 0) {
        acc.push({ ...node, children: filteredChildren })
      }
      return acc
    }

    acc.push(node)
    return acc
  }, [])
}

const filteredTree = computed(() => {
  if (!props.searchFilter || !props.searchFilter.trim()) {
    return props.tree
  }
  return filterTree(props.tree, props.searchFilter)
})

function onClick(node: TreeItem) { if (node.type === 'request') emit('open-request', node) }
function onDblClick(node: TreeItem) { emit('dbl-click', node) }
function onCtxMenu(node: TreeItem, event: MouseEvent) { emit('ctx-menu', node, event) }
function onAction(action: string, node: TreeItem) { emit('action', action, node) }
</script>

<style scoped>
.collection-tree {
  flex: 1;
  overflow-y: auto;
  padding: 4px 0;
}

/* Empty state styling */
.tree-empty {
  padding: 32px 20px;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
}
.empty-illustr {
  margin-bottom: 8px;
  opacity: 0.3;
  transition: opacity var(--transition);
}
.empty-icon {
  color: var(--text-muted);
}
.tree-empty:hover .empty-illustr {
  opacity: 0.45;
}
.empty-text {
  margin: 0;
  font-size: var(--fs-sm);
  font-weight: 600;
  color: var(--text-secondary);
  font-family: var(--font-ui);
}
.empty-hint {
  margin: 0;
  font-size: var(--fs-xs);
  color: var(--text-muted);
  line-height: 1.5;
  max-width: 180px;
}
</style>
