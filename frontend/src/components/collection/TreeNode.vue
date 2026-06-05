<template>
  <div>
    <div
      class="tree-node"
      :class="[node.type, { 'is-active': isActive, 'is-drop-target': isDropTarget }]"
      :style="{ paddingLeft: node.type === 'request' ? (12 + depth * 14 + 10) + 'px' : (8 + depth * 14) + 'px' }"
      @click="onNodeClick"
      @dblclick="onNodeDblClick"
      @contextmenu.prevent="onCtxMenu($event)"
      @mouseenter="isHovered = true"
      @mouseleave="isHovered = false"
      @dragover.prevent="onDragOver"
      @dragleave="onDragLeave"
      @drop.prevent="onDrop"
      ref="nodeRef"
    >
      <!-- Folder arrow -->
      <span v-if="node.type === 'folder'" class="arrow" @click.stop="toggleExpand">
        <svg width="8" height="8" viewBox="0 0 10 10" :class="{ rotated: expanded }">
          <path d="M3 1 L7 5 L3 9" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
      </span>

      <!-- Root name -->
      <span v-if="node.type === 'root'" class="node-name root-name">{{ node.name }}</span>

      <!-- Folder: icon + name -->
      <template v-else-if="node.type === 'folder'">
        <svg class="node-icon folder-icon" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
          <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
        </svg>
        <span class="node-name">{{ node.name }}</span>
      </template>

      <!-- Request: method tag + icon + name -->
      <template v-else-if="node.type === 'request'">
        <span class="method-tag" :class="node.method?.toLowerCase()">{{ node.method }}</span>
        <span class="node-name">{{ node.name }}</span>
      </template>

      <!-- Drop indicator line -->
      <div v-if="isDropTarget" class="drop-indicator"></div>
    </div>

    <!-- Children with smooth animation -->
    <Transition name="tree-expand">
      <div v-if="node.type === 'folder' && expanded" class="tree-children">
        <TreeNode
          v-for="child in node.children || []"
          :key="child.id"
          :node="child"
          :depth="depth + 1"
          @click="onChildClick"
          @dbl-click="onChildDblClick"
          @ctx-menu="onChildCtxMenu"
          @action="onChildAction"
        />
      </div>
    </Transition>

    <n-dropdown
      placement="bottom-start"
      trigger="manual"
      :x="ctxX"
      :y="ctxY"
      :options="ctxMenuOptions"
      :show="ctxVisible"
      :on-clickoutside="() => { ctxVisible = false }"
      @select="onAction"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { NDropdown } from 'naive-ui'
import type { TreeItem } from '../../types/collection'

const props = defineProps<{
  node: TreeItem
  depth: number
  /** ID of the currently active/open request (for highlighting) */
  activeRequestId?: number
}>()

const emit = defineEmits<{
  (e: 'click', node: TreeItem): void
  (e: 'dbl-click', node: TreeItem): void
  (e: 'ctx-menu', node: TreeItem, event: MouseEvent): void
  (e: 'action', action: string, node: TreeItem): void
}>()

const expanded = ref(true)
const ctxVisible = ref(false)
const ctxX = ref(0)
const ctxY = ref(0)
const isHovered = ref(false)
const isDropTarget = ref(false)
const nodeRef = ref<HTMLElement | null>(null)

const isActive = computed(() => {
  return props.node.type === 'request' && props.activeRequestId === props.node.id
})

const ctxMenuOptions = computed(() => {
  if (props.node.type === 'root' || props.node.type === 'folder') {
    return [
      { label: '新建请求', key: 'new-request' },
      { label: '新建文件夹', key: 'new-folder' },
      { type: 'divider' as const, key: 'd1' },
      { label: '重命名              F2', key: 'rename' },
      ...(props.node.type === 'folder' ? [{ label: '导出此集合', key: 'export' }] : []),
      { type: 'divider' as const, key: 'd2' },
      { label: '删除', key: 'delete' },
    ]
  }
  return [
    { label: '在标签中打开', key: 'open-in-tab' },
    { type: 'divider' as const, key: 'd0' },
    { label: '重命名              F2', key: 'rename' },
    { label: '复制', key: 'duplicate' },
    { type: 'divider' as const, key: 'd1' },
    { label: '删除', key: 'delete' },
  ]
})

function toggleExpand() {
  expanded.value = !expanded.value
}

function onNodeClick() { emit('click', props.node) }
function onNodeDblClick() { emit('dbl-click', props.node) }
function onCtxMenu(event: MouseEvent) {
  emit('ctx-menu', props.node, event)
  ctxX.value = event.clientX
  ctxY.value = event.clientY
  ctxVisible.value = true
}
function onChildDblClick(node: TreeItem) { emit('dbl-click', node) }
function onChildClick(node: TreeItem) { emit('click', node) }
function onChildCtxMenu(node: TreeItem, event: MouseEvent) { emit('ctx-menu', node, event) }
function onChildAction(action: string, node: TreeItem) { emit('action', action, node) }
function onAction(key: string) {
  ctxVisible.value = false
  emit('action', key, props.node)
}

// Drag & drop feedback
function onDragOver(e: DragEvent) {
  if (props.node.type === 'folder' || props.node.type === 'root') {
    isDropTarget.value = true
  }
}
function onDragLeave() {
  isDropTarget.value = false
}
function onDrop(e: DragEvent) {
  isDropTarget.value = false
  // Delegate drop handling via action emit
  emit('action', 'drop', props.node)
}

// F2 to rename when hovered
function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    ctxVisible.value = false
    return
  }
  if (e.key === 'F2' && isHovered.value) {
    e.preventDefault()
    emit('action', 'rename', props.node)
  }
}

onMounted(() => document.addEventListener('keydown', onKeydown))
onUnmounted(() => document.removeEventListener('keydown', onKeydown))
</script>

<style scoped>
.tree-node {
  display: flex;
  align-items: center;
  padding: 5px 10px;
  cursor: pointer;
  white-space: nowrap;
  font-size: var(--fs-sm);
  gap: 6px;
  user-select: none;
  color: var(--text-primary);
  transition: background var(--transition-fast), border-left-color 0.15s var(--ease-out);
  border-left: 2px solid transparent;
  position: relative;
}
.tree-node:hover {
  background: var(--bg-hover);
}
.tree-node.request:hover {
  border-left-color: var(--accent);
}
.tree-node.folder:hover {
  border-left-color: var(--amber);
}

/* Active/selected state for open requests */
.tree-node.request.is-active {
  background: var(--accent-soft);
  border-left-color: var(--accent);
}
.tree-node.request.is-active .node-name {
  color: var(--accent);
  font-weight: 600;
}

/* Drop target indicator */
.tree-node.is-drop-target {
  background: var(--accent-soft);
  border-left-color: var(--accent);
}
.drop-indicator {
  position: absolute;
  bottom: 0;
  left: 8px;
  right: 8px;
  height: 2px;
  background: var(--accent);
  border-radius: 1px;
  box-shadow: 0 0 4px var(--accent-glow);
}

.tree-node.root {
  font-weight: 600;
  color: var(--text-secondary);
  font-size: var(--fs-xs);
  text-transform: uppercase;
  letter-spacing: 0.04em;
  padding: 6px 10px;
  border-left: none;
}

/* Folder arrow */
.arrow {
  display: flex; align-items: center; width: 12px;
  color: var(--text-muted); flex-shrink: 0;
  cursor: pointer;
}
.arrow svg {
  transition: transform 0.22s var(--ease-out);
}
.arrow svg.rotated { transform: rotate(90deg); }

/* Node icons */
.node-icon {
  flex-shrink: 0;
  opacity: 0.6;
  transition: opacity var(--transition-fast);
}
.tree-node:hover .node-icon {
  opacity: 0.85;
}
.folder-icon {
  color: var(--amber);
}

/* Smooth expand/collapse animation */
.tree-children {
  overflow: hidden;
}
.tree-expand-enter-active {
  transition: all 0.25s var(--ease-out);
  overflow: hidden;
}
.tree-expand-leave-active {
  transition: all 0.18s var(--ease-out);
  overflow: hidden;
}
.tree-expand-enter-from {
  opacity: 0;
  max-height: 0;
  transform: translateY(-4px);
}
.tree-expand-enter-to {
  opacity: 1;
  max-height: 2000px;
}
.tree-expand-leave-from {
  opacity: 1;
  max-height: 2000px;
}
.tree-expand-leave-to {
  opacity: 0;
  max-height: 0;
  transform: translateY(-4px);
}

/* Method tag: slightly larger, more distinct */
.method-tag {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 42px;
  font-size: var(--fs-2xs);
  font-weight: 700;
  text-align: center;
  padding: 2px 6px;
  border-radius: var(--radius-xs);
  flex-shrink: 0;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  line-height: 1.5;
  font-family: var(--font-mono);
}
.method-tag.get    { background: var(--blue-soft); color: var(--method-get); }
.method-tag.post   { background: rgba(34,197,94,0.12); color: var(--method-post); }
.method-tag.put    { background: var(--amber-soft); color: var(--method-put); }
.method-tag.delete { background: var(--red-soft); color: var(--method-delete); }
.method-tag.patch  { background: var(--purple-soft); color: var(--method-patch); }
.method-tag.head,
.method-tag.options { background: rgba(113,113,122,0.1); color: var(--text-secondary); }

.node-name {
  overflow: hidden;
  text-overflow: ellipsis;
  color: var(--text-primary);
  font-weight: 500;
  transition: color var(--transition-fast);
}
.root-name {
  font-family: var(--font-ui);
}
</style>
