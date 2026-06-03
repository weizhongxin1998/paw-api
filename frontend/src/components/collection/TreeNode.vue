<template>
  <div>
    <div
      class="tree-node"
      :class="[node.type]"
      :style="{ paddingLeft: node.type === 'request' ? (12 + depth * 14 + 10) + 'px' : (8 + depth * 14) + 'px' }"
      @click="onNodeClick"
      @dblclick="onNodeDblClick"
      @contextmenu.prevent="onCtxMenu($event)"
    >
      <span v-if="node.type === 'folder'" class="arrow" @click.stop="expanded = !expanded">
        <svg width="8" height="8" viewBox="0 0 10 10" :class="{ rotated: expanded }">
          <path d="M3 1 L7 5 L3 9" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
      </span>
      <span v-if="node.type === 'root'" class="node-name root-name">{{ node.name }}</span>
      <template v-else-if="node.type === 'folder'">
        <span class="node-name">{{ node.name }}</span>
      </template>
      <template v-else-if="node.type === 'request'">
        <span class="method-tag" :class="node.method?.toLowerCase()">{{ node.method }}</span>
        <span class="node-name">{{ node.name }}</span>
      </template>
    </div>
    <Transition name="tree-expand">
      <div v-if="node.type === 'folder' && expanded">
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

const ctxMenuOptions = computed(() => {
  if (props.node.type === 'root' || props.node.type === 'folder') {
    return [
      { label: '新建请求', key: 'new-request' },
      { label: '新建文件夹', key: 'new-folder' },
      { type: 'divider' as const, key: 'd1' },
      { label: '重命名', key: 'rename' },
      ...(props.node.type === 'folder' ? [{ label: '导出此集合', key: 'export' }] : []),
      { type: 'divider' as const, key: 'd2' },
      { label: '删除', key: 'delete' },
    ]
  }
  return [
    { label: '重命名', key: 'rename' },
    { label: '复制', key: 'duplicate' },
    { type: 'divider' as const, key: 'd1' },
    { label: '删除', key: 'delete' },
  ]
})

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
function onKeydown(e: KeyboardEvent) { if (e.key === 'Escape') ctxVisible.value = false }
onMounted(() => document.addEventListener('keydown', onKeydown))
onUnmounted(() => document.removeEventListener('keydown', onKeydown))
</script>

<style scoped>
.tree-node {
  display: flex;
  align-items: center;
  padding: 4px 8px;
  cursor: pointer;
  white-space: nowrap;
  font-size: var(--fs-sm);
  gap: 4px;
  user-select: none;
  color: var(--text-primary);
  transition: background var(--transition);
  border-left: 2px solid transparent;
}
.tree-node:hover { background: var(--bg-hover); }
.tree-node.request:hover { border-left-color: var(--accent); }
.arrow {
  display: flex; align-items: center; width: 10px;
  color: var(--text-muted); flex-shrink: 0;
}
.arrow svg { transition: transform 0.2s cubic-bezier(0.4, 0, 0.2, 1); }
.arrow svg.rotated { transform: rotate(90deg); }

.tree-expand-enter-active { transition: all 0.15s ease; overflow: hidden; }
.tree-expand-leave-active { transition: all 0.1s ease; overflow: hidden; }
.tree-expand-enter-from { opacity: 0; max-height: 0; }
.tree-expand-leave-to { opacity: 0; max-height: 0; }
.method-tag {
  display: inline-block;
  min-width: 34px;
  font-size: var(--fs-2xs);
  font-weight: 700;
  text-align: center;
  padding: 1px 3px;
  border-radius: 2px;
  flex-shrink: 0;
  letter-spacing: 0.3px;
}
.method-tag.get { background: var(--accent-soft); color: var(--accent); }
.method-tag.post { background: var(--amber-soft); color: var(--amber); }
.method-tag.put { background: var(--blue-soft); color: var(--blue); }
.method-tag.delete { background: var(--red-soft); color: var(--red); }
.method-tag.patch { background: var(--purple-soft); color: var(--purple); }
.method-tag.head, .method-tag.options { background: var(--bg-hover); color: var(--text-secondary); }
.node-name { overflow: hidden; text-overflow: ellipsis; color: var(--text-primary); font-weight: 500; }
</style>
