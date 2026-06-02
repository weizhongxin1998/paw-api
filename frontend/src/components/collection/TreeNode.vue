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
        <svg width="10" height="10" viewBox="0 0 10 10" :class="{ rotated: expanded }">
          <path d="M3 1 L7 5 L3 9" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
      </span>
      <span v-if="node.type === 'root'" class="node-name root-name">{{ node.name }}</span>
      <template v-else-if="node.type === 'folder'">
        <span class="node-name folder-name">{{ node.name }}</span>
      </template>
      <template v-else-if="node.type === 'request'">
        <span class="method-tag" :class="node.method?.toLowerCase()">{{ node.method }}</span>
        <span class="node-name">{{ node.name }}</span>
        <span class="node-url">{{ node.url }}</span>
      </template>
    </div>
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
  </div>

  <div v-if="ctxVisible" class="ctx-overlay" @click="ctxVisible = false">
    <div class="ctx-menu" :style="{ left: ctxX + 'px', top: ctxY + 'px' }" @click.stop>
      <template v-if="node.type === 'root' || node.type === 'folder'">
        <div class="ctx-item" @click="onAction('new-request')">+ 新建请求</div>
        <div class="ctx-item" @click="onAction('new-folder')">新建文件夹</div>
        <div class="ctx-sep"></div>
        <div class="ctx-item" @click="onAction('rename')">重命名</div>
        <div v-if="node.type === 'folder'" class="ctx-item" @click="onAction('export')">导出此集合</div>
        <div class="ctx-sep"></div>
        <div class="ctx-item danger" @click="onAction('delete')">删除</div>
      </template>
      <template v-else>
        <div class="ctx-item" @click="onAction('rename')">重命名</div>
        <div class="ctx-item" @click="onAction('duplicate')">复制</div>
        <div class="ctx-sep"></div>
        <div class="ctx-item danger" @click="onAction('delete')">删除</div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
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

function onNodeClick() {
  emit('click', props.node)
}

function onNodeDblClick() {
  emit('dbl-click', props.node)
}

function onCtxMenu(event: MouseEvent) {
  emit('ctx-menu', props.node, event)
  ctxX.value = event.clientX
  ctxY.value = event.clientY
  ctxVisible.value = true
}

function onChildDblClick(node: TreeItem) {
  emit('dbl-click', node)
}

function onChildClick(node: TreeItem) {
  emit('click', node)
}

function onChildCtxMenu(node: TreeItem, event: MouseEvent) {
  emit('ctx-menu', node, event)
}

function onChildAction(action: string, node: TreeItem) {
  emit('action', action, node)
}

function onAction(action: string) {
  ctxVisible.value = false
  emit('action', action, props.node)
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') ctxVisible.value = false
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
  font-size: 12px;
  gap: 5px;
  user-select: none;
  transition: background var(--transition);
  border-radius: 0 4px 4px 0;
  margin-right: 2px;
}
.tree-node:hover { background: var(--gray-100); }
.tree-node.root { font-weight: 600; color: var(--green); }
.arrow { display: flex; align-items: center; width: 12px; color: var(--gray-400); transition: transform var(--transition); flex-shrink: 0; }
.arrow svg { transition: transform 0.2s ease; }
.arrow svg.rotated { transform: rotate(90deg); }
.method-tag {
  display: inline-block;
  width: 34px;
  font-size: 9px;
  font-weight: 700;
  text-align: center;
  padding: 2px 3px;
  border-radius: 3px;
  flex-shrink: 0;
  letter-spacing: 0.3px;
}
.method-tag.get { background: var(--green-soft); color: var(--green); }
.method-tag.post { background: var(--amber-soft); color: var(--amber); }
.method-tag.put { background: var(--blue-soft); color: var(--blue); }
.method-tag.delete { background: var(--red-soft); color: var(--red); }
.method-tag.patch { background: var(--purple-soft); color: var(--purple); }
.root-name { font-weight: 600; color: var(--green); }
.folder-name { color: var(--gray-700); font-weight: 500; }
.node-name { overflow: hidden; text-overflow: ellipsis; }
.node-url { color: var(--gray-400); margin-left: 3px; font-size: 10px; overflow: hidden; text-overflow: ellipsis; }
.ctx-overlay { position: fixed; inset: 0; z-index: 1000; }
.ctx-menu {
  position: fixed;
  background: #fff;
  border: 1px solid var(--gray-200);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg);
  padding: 4px 0;
  min-width: 160px;
  z-index: 1001;
  animation: ctxFadeIn 0.12s ease;
}
@keyframes ctxFadeIn { from { opacity: 0; transform: translateY(-4px); } to { opacity: 1; transform: translateY(0); } }
.ctx-item { padding: 6px 14px; cursor: pointer; font-size: 12px; color: var(--gray-600); transition: background var(--transition); }
.ctx-item:hover { background: var(--gray-100); }
.ctx-item.danger { color: var(--red); }
.ctx-sep { border-top: 1px solid var(--gray-100); margin: 4px 0; }
</style>
