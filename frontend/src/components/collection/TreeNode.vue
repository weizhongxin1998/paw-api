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

  <div v-if="ctxVisible" class="ctx-overlay" @click="ctxVisible = false">
    <div class="ctx-menu" :style="{ left: ctxX + 'px', top: ctxY + 'px' }" @click.stop>
      <template v-if="node.type === 'folder'">
        <div class="ctx-item" @click="onAction('new-request')">+ 新建请求</div>
        <div class="ctx-item" @click="onAction('new-folder')">新建文件夹</div>
        <div class="ctx-sep"></div>
        <div class="ctx-item" @click="onAction('rename')">重命名</div>
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
  (e: 'dbl-click', node: TreeItem): void
  (e: 'ctx-menu', node: TreeItem, event: MouseEvent): void
  (e: 'action', action: string, node: TreeItem): void
}>()

const expanded = ref(true)
const ctxVisible = ref(false)
const ctxX = ref(0)
const ctxY = ref(0)

function onDblClick() {
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

function onChildCtxMenu(node: TreeItem, event: MouseEvent) {
  emit('ctx-menu', node, event)
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
  padding: 4px 8px;
  cursor: pointer;
  white-space: nowrap;
  font-size: 12px;
  gap: 4px;
  user-select: none;
}
.tree-node:hover { background: #f0f0f0; }
.arrow { font-size: 9px; width: 12px; color: #888; cursor: pointer; }
.method-tag {
  display: inline-block; width: 36px; font-size: 9px; font-weight: 700;
  text-align: center; padding: 0 2px; border-radius: 2px;
}
.method-tag.get { color: #18a058; }
.method-tag.post { color: #f0a020; }
.method-tag.put { color: #2080f0; }
.method-tag.delete { color: #d03050; }
.method-tag.patch { color: #9c27b0; }
.node-name { overflow: hidden; text-overflow: ellipsis; }
.node-url { color: #aaa; font-size: 10px; margin-left: 4px; overflow: hidden; text-overflow: ellipsis; }
.ctx-overlay { position: fixed; inset: 0; z-index: 1000; }
.ctx-menu {
  position: fixed; background: #fff; border: 1px solid #e0e0e0; border-radius: 8px;
  box-shadow: 0 3px 12px rgba(0,0,0,0.12); padding: 4px 0; min-width: 160px; z-index: 1001;
}
.ctx-item { padding: 6px 14px; cursor: pointer; font-size: 12px; color: #333; }
.ctx-item:hover { background: #f0f0f0; }
.ctx-item.danger { color: #d03050; }
.ctx-sep { border-top: 1px solid #eee; margin: 4px 0; }
</style>
