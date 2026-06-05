<template>
  <div class="sidebar" :class="{ collapsed: sidebarCollapsed }">
    <!-- Expanded state -->
    <template v-if="!sidebarCollapsed">
      <div class="sidebar-tabs">
        <div class="tab-pills">
          <button
            class="tab-pill"
            :class="{ active: activePanel === 'collection' }"
            @click="activePanel = 'collection'"
          >{{ t('sidebar.collections') }}</button>
          <button
            class="tab-pill"
            :class="{ active: activePanel === 'history' }"
            @click="activePanel = 'history'"
          >{{ t('sidebar.history') }}</button>
        </div>
        <n-dropdown trigger="click" :options="addMenuOptions" @select="onAddMenuSelect">
          <n-button size="tiny" class="add-btn" circle>
            <template #icon>
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
            </template>
          </n-button>
        </n-dropdown>
      </div>

      <!-- Search filter for collections -->
      <div v-if="activePanel === 'collection'" class="sidebar-search">
        <n-input
          ref="searchInputRef"
          v-model:value="searchQuery"
          :placeholder="t('sidebar.searchCollections')"
          size="tiny"
          clearable
          class="search-input"
        >
          <template #prefix>
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="opacity:0.65">
              <circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/>
            </svg>
          </template>
        </n-input>
      </div>

      <div class="sidebar-body">
        <CollectionTree
          v-if="activePanel === 'collection'"
          :tree="tree"
          :search-filter="searchQuery"
          @open-request="onPreviewRequest"
          @dbl-click="onPersistRequest"
          @ctx-menu="onContextMenu"
          @action="onTreeAction"
          ref="collectionTreeRef"
        />
        <HistoryPanel
          v-else
          @open-tab="onHistoryReplay"
          @select-detail="onHistorySelect"
          ref="historyPanelRef"
        />
      </div>
    </template>

    <!-- Collapsed state: icon-only strip -->
    <template v-else>
      <div class="collapsed-rail">
        <button
          class="rail-btn"
          :class="{ active: activePanel === 'collection' }"
          @click="activePanel = 'collection'; sidebarCollapsed = false"
          :title="t('sidebar.collectionsExpand')"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
            <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
          </svg>
        </button>
        <button
          class="rail-btn"
          :class="{ active: activePanel === 'history' }"
          @click="activePanel = 'history'; sidebarCollapsed = false"
          :title="t('sidebar.historyTooltip')"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/>
          </svg>
        </button>
        <div class="rail-spacer"></div>
        <n-dropdown trigger="click" :options="addMenuOptions" @select="onAddMenuSelect">
          <button class="rail-btn rail-add-btn" :title="t('sidebar.createNew')">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
          </button>
        </n-dropdown>
      </div>
    </template>

    <n-modal v-model:show="showCollectionModal" preset="card" :title="t('sidebar.createCollection')" :class="modalClass" style="width: 360px" :mask-closable="false">
      <n-form label-placement="top">
        <n-form-item :label="t('common.name')">
          <n-input v-model:value="collectionName" :placeholder="t('sidebar.collectionNamePlaceholder')" @keydown.enter="onCreateCollection" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button @click="showCollectionModal = false">{{ t('common.cancel') }}</n-button>
        <n-button type="primary" :disabled="!collectionName.trim()" @click="onCreateCollection">{{ t('common.create') }}</n-button>
      </template>
    </n-modal>

    <n-modal v-model:show="showRequestModal" preset="card" :title="t('sidebar.createRequest')" :class="modalClass" style="width: 360px" :mask-closable="false">
      <n-form label-placement="top">
        <n-form-item :label="t('common.name')">
          <n-input v-model:value="requestName" :placeholder="t('sidebar.requestNamePlaceholder')" @keydown.enter="onCreateRequest" />
        </n-form-item>
        <n-form-item :label="t('sidebar.method')">
          <n-select v-model:value="requestMethod" :options="methodOptions" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button @click="showRequestModal = false">{{ t('common.cancel') }}</n-button>
        <n-button type="primary" :disabled="!requestName.trim()" @click="onCreateRequest">{{ t('common.create') }}</n-button>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { NButton, NDropdown, NModal, NForm, NFormItem, NInput, NSelect, useMessage } from 'naive-ui'
import type { DropdownOption } from 'naive-ui'
import { useCollectionStore } from '../../stores/collection'
import CollectionTree from '../collection/CollectionTree.vue'
import HistoryPanel from '../history/HistoryPanel.vue'
import type { TreeItem } from '../../types/collection'

const { t } = useI18n()

const props = defineProps<{
  projectId: number | null
}>()

const activePanel = ref('collection')
const historyPanelRef = ref<InstanceType<typeof HistoryPanel> | null>(null)
const collectionTreeRef = ref<InstanceType<typeof CollectionTree> | null>(null)
const searchInputRef = ref<InstanceType<typeof NInput> | null>(null)
const collectionStore = useCollectionStore()
const message = useMessage()

const emit = defineEmits<{
  (e: 'open-request', node: TreeItem): void
  (e: 'open-request-persist', node: TreeItem): void
  (e: 'history-replay', item: any): void
  (e: 'history-select', item: any | null): void
  (e: 'open-docs'): void
  (e: 'tree-action', action: string, node: TreeItem): void
}>()

const tree = ref<TreeItem[]>([])
const showCollectionModal = ref(false)
const showRequestModal = ref(false)
const collectionName = ref('')
const requestName = ref('')
const requestMethod = ref('GET')
const searchQuery = ref('')
const sidebarCollapsed = ref(false)

const isLightMode = ref(false)
const modalClass = computed(() => isLightMode.value ? 'sidebar-modal theme-light' : 'sidebar-modal')

const methodOptions = [
  { label: 'GET', value: 'GET' }, { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' }, { label: 'DELETE', value: 'DELETE' },
  { label: 'PATCH', value: 'PATCH' }, { label: 'HEAD', value: 'HEAD' },
  { label: 'OPTIONS', value: 'OPTIONS' },
]

const addMenuOptions: DropdownOption[] = [
  { label: t('sidebar.createRequest'), key: 'request' },
  { label: t('sidebar.createCollection'), key: 'collection' },
  { type: 'divider', key: 'div' },
  { label: t('sidebar.apiDocs'), key: 'docs' },
]

function onAddMenuSelect(key: string) {
  switch (key) {
    case 'request': showRequestModal.value = true; break
    case 'collection': showCollectionModal.value = true; break
    case 'docs': emit('open-docs'); break
  }
}

watch(activePanel, (val) => {
  if (val === 'collection') emit('history-select', null)
})

watch(() => props.projectId, async (id) => {
  if (id) tree.value = await collectionStore.loadTree(id)
  else tree.value = []
}, { immediate: true })

async function refreshTree() {
  if (props.projectId) tree.value = await collectionStore.loadTree(props.projectId)
}

function onPreviewRequest(node: TreeItem) { emit('open-request', node) }
function onPersistRequest(node: TreeItem) {
  if (node.type === 'request') emit('open-request-persist', node)
}
function onContextMenu(_: TreeItem, _ev: MouseEvent) {}
function onTreeAction(action: string, node: TreeItem) { emit('tree-action', action, node) }

async function onCreateCollection() {
  const name = collectionName.value.trim()
  if (!name || !props.projectId) return
  try {
    await collectionStore.createCollection(props.projectId, null, name)
    showCollectionModal.value = false
    collectionName.value = ''
    await refreshTree()
    message.success(t('sidebar.collectionCreated', { name }))
    // Auto-scroll to newly created item
    nextTick(() => {
      scrollToNewestItem()
    })
  } catch (e: any) {
    message.error(t('sidebar.createFailed', { error: e?.message || String(e) }))
  }
}

async function onCreateRequest() {
  const name = requestName.value.trim()
  if (!name || !props.projectId) return
  if (tree.value.length > 0) {
    try {
      await collectionStore.createRequest(tree.value[0].id, name, requestMethod.value)
      showRequestModal.value = false
      requestName.value = ''
      requestMethod.value = 'GET'
      await refreshTree()
      message.success(t('sidebar.requestCreated', { name }))
      nextTick(() => {
        scrollToNewestItem()
      })
    } catch (e: any) {
      message.error(t('sidebar.createFailed', { error: e?.message || String(e) }))
    }
  }
}

function scrollToNewestItem() {
  const treeEl = collectionTreeRef.value?.$el as HTMLElement | undefined
  if (!treeEl) return
  const items = treeEl.querySelectorAll('.tree-node')
  if (items.length > 0) {
    const lastItem = items[items.length - 1] as HTMLElement
    lastItem.scrollIntoView({ behavior: 'smooth', block: 'nearest' })
  }
}

function onHistoryReplay(item: any) { emit('history-replay', item) }
function onHistorySelect(item: any) { emit('history-select', item) }

function onGlobalKeydown(e: KeyboardEvent) {
  // Ctrl+B: toggle sidebar
  if (e.ctrlKey && !e.shiftKey && e.key === 'b') {
    e.preventDefault()
    sidebarCollapsed.value = !sidebarCollapsed.value
  }
  // Ctrl+Shift+F: focus search input
  if (e.ctrlKey && e.shiftKey && e.key === 'F') {
    e.preventDefault()
    if (sidebarCollapsed.value) sidebarCollapsed.value = false
    activePanel.value = 'collection'
    nextTick(() => {
      const inputEl = searchInputRef.value?.$el?.querySelector('input') as HTMLInputElement | null
      if (inputEl) {
        inputEl.focus()
        inputEl.select()
      }
    })
  }
}

onMounted(() => {
  document.addEventListener('keydown', onGlobalKeydown)
  const check = () => { isLightMode.value = !!document.querySelector('.theme-light') }
  check()
  const observer = new MutationObserver(check)
  observer.observe(document.documentElement, { attributes: true, attributeFilter: ['class'], subtree: true })
})
onUnmounted(() => {
  document.removeEventListener('keydown', onGlobalKeydown)
})

defineExpose({ refreshTree })
</script>

<style scoped>
.sidebar {
  width: 240px;
  background: var(--bg-surface);
  border-right: 1px solid var(--border-primary);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  position: relative;
  /* Subtle top shadow to separate from header */
  box-shadow: inset 0 1px 3px rgba(0,0,0,0.08);
  transition: width 0.2s var(--ease-out);
}
.sidebar.collapsed {
  width: 42px;
}

/* ── Pill / Segmented Tabs ── */
.sidebar-tabs {
  display: flex;
  align-items: center;
  padding: 0 8px 0 8px;
  height: 38px;
  border-bottom: 1px solid var(--border-subtle);
  gap: 6px;
}

.tab-pills {
  display: flex;
  flex: 1;
  background: var(--bg-inset);
  border-radius: var(--radius);
  padding: 2px;
  gap: 2px;
}
.tab-pill {
  flex: 1;
  border: none;
  background: transparent;
  color: var(--text-secondary);
  font-size: var(--fs-sm);
  font-weight: 500;
  font-family: var(--font-ui);
  padding: 5px 0;
  border-radius: calc(var(--radius) - 2px);
  cursor: pointer;
  transition: all var(--transition);
  letter-spacing: 0.01em;
  outline: none;
}
.tab-pill:hover:not(.active) {
  color: var(--text-secondary);
  background: var(--bg-hover);
}
.tab-pill.active {
  color: var(--text-primary);
  background: var(--bg-elevated);
  box-shadow: var(--shadow-sm);
}

/* ── Add Button ── */
.add-btn {
  width: 26px !important;
  height: 26px !important;
  min-width: 26px !important;
  margin-left: 2px;
  background: var(--accent-soft) !important;
  border: 1px solid transparent !important;
  color: var(--accent) !important;
  transition: all var(--transition) !important;
}
.add-btn:hover {
  background: linear-gradient(135deg, var(--accent-glow), var(--accent-soft)) !important;
  border-color: var(--accent) !important;
  transform: scale(1.1);
  box-shadow: 0 0 8px var(--accent-glow);
}
.add-btn:active {
  transform: scale(0.95);
}

/* ── Search Input ── */
.sidebar-search {
  padding: 6px 10px;
  border-bottom: 1px solid var(--border-subtle);
}
.search-input :deep(.n-input) {
  --n-height: 28px !important;
  border-radius: var(--radius-sm) !important;
}
.search-input :deep(.n-input__border) {
  border-color: var(--border-subtle) !important;
}
.search-input :deep(.n-input__border:hover) {
  border-color: var(--border-hover) !important;
}
.search-input :deep(input) {
  font-size: var(--fs-xs) !important;
}

/* ── Sidebar Body ── */
.sidebar-body {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

/* ── Collapsed Rail ── */
.collapsed-rail {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 8px 0;
  gap: 4px;
  height: 100%;
}
.rail-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: transparent;
  color: var(--text-muted);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all var(--transition);
  outline: none;
}
.rail-btn:hover {
  color: var(--text-primary);
  background: var(--bg-hover);
}
.rail-btn.active {
  color: var(--accent);
  background: var(--accent-soft);
}
.rail-add-btn {
  color: var(--accent) !important;
}
.rail-add-btn:hover {
  background: var(--accent-glow) !important;
  box-shadow: 0 0 8px var(--accent-glow);
}
.rail-spacer {
  flex: 1;
}
</style>
