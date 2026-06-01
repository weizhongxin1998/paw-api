<script lang="ts" setup>
import { ref, computed, onMounted, watch, h } from 'vue'
import { NTree, NButton, NIcon, NInput, NModal, NForm, NFormItem, NSpace, NSelect, NTag, NDropdown, useMessage, useDialog } from 'naive-ui'
import { Add, CodeSlash, FolderOpen, Time, DocumentText, Bug, Settings as SettingsIcon, Search, Download, Trash } from '@vicons/ionicons5'
import { useRouter, useRoute } from 'vue-router'
import { useProjectStore } from '../stores/project'
import { useTabsStore } from '../stores/tabs'
import { useI18n } from 'vue-i18n'
import { useTheme } from '../composables/useTheme'
import { t, setLocale as setI18nLocale } from '../i18n'
import { CreateProject, ListProjects } from '../../wailsjs/go/handlers/ProjectHandler'
import { CreateCollection, ListCollections, UpdateCollection, DeleteCollection } from '../../wailsjs/go/handlers/CollectionHandler'
import { CreateRequest, DeleteRequest, UpdateRequest, ListRequests } from '../../wailsjs/go/handlers/RequestHandler'
import { ListEnvironments } from '../../wailsjs/go/handlers/EnvironmentHandler'
import { useEnvironmentStore } from '../stores/environment'
import EnvSelector from './EnvSelector.vue'
import EnvManager from './EnvManager.vue'

const router = useRouter()
const route = useRoute()
const projectStore = useProjectStore()
const tabsStore = useTabsStore()
const envStore = useEnvironmentStore()
const { toggleColorMode, setThemeColor, themeColor, colorMode } = useTheme()
const { locale } = useI18n()
const message = useMessage()
const dialog = useDialog()

const activeSection = ref<'workspace' | 'project' | 'history' | 'docs' | 'tests' | 'settings'>('workspace')
const showAddModal = ref(false)
const newCollectionName = ref('')
const selectedParentId = ref<string | null>(null)
const showEnvManager = ref(false)

const searchQuery = ref('')

const collectionMenuId = ref<string | null>(null)
const collectionMenuPos = ref({ x: 0, y: 0 })
const showRenameModal = ref(false)
const renameTarget = ref<{ id: string; name: string; type: 'collection' | 'request' } | null>(null)
const renameValue = ref('')
const showMoveModal = ref(false)
const moveTargetId = ref('')
const contextMenuColId = ref<string | null>(null)
const contextMenuReqId = ref<string | null>(null)
const showContextMenu = ref(false)
const contextMenuX = ref(0)
const contextMenuY = ref(0)

const sections = [
  { id: 'workspace', labelKey: 'sidebar.workspaceLabel', icon: CodeSlash, route: '/workspace' },
  { id: 'project', labelKey: 'sidebar.projectLabel', icon: FolderOpen, route: '/projects' },
  { id: 'history', labelKey: 'sidebar.historyLabel', icon: Time, route: '/history' },
  { id: 'docs', labelKey: 'sidebar.docsLabel', icon: DocumentText, route: '/docs' },
  { id: 'tests', labelKey: 'sidebar.testsLabel', icon: Bug, route: '/test-runner' },
  { id: 'settings', labelKey: 'sidebar.settingsLabel', icon: SettingsIcon, route: '' },
]

const themeColorOptions = [
  { label: 'Green', value: 'green' },
  { label: 'Blue', value: 'blue' },
  { label: 'Purple', value: 'purple' },
]

// Load all requests for search and display
const allRequests = ref<any[]>([])

async function loadAllRequests() {
  if (!projectStore.currentProject) return
  try {
    const cols = await ListCollections(projectStore.currentProject.id)
    const all: any[] = []
    for (const col of cols) {
      const reqs = await ListRequests(col.id)
      all.push(...reqs.map((r: any) => ({ ...r, _collectionName: col.name, _collectionId: col.id })))
    }
    allRequests.value = all
  } catch { allRequests.value = [] }
}

const treeData = computed(() => {
  if (activeSection.value !== 'workspace') return []
  const q = searchQuery.value.toLowerCase().trim()

  function buildTree(parentId: string | null): any[] {
    return (projectStore.collections || [])
      .filter(c => c.parent_id === parentId)
      .sort((a, b) => a.sort_order - b.sort_order)
      .map(c => {
        let children = buildTree(c.id)
        const reqs = allRequests.value
          .filter(r => r._collectionId === c.id)
          .filter(r => !q || r.name.toLowerCase().includes(q) || r.url.toLowerCase().includes(q))
          .map(r => ({ label: r.name, key: 'req-' + r.id, isLeaf: true, _reqId: r.id, _method: r.method }))
        if (children.length > 0 || reqs.length > 0) {
          children = [...children, ...reqs]
        }
        const label = q && !c.name.toLowerCase().includes(q) && children.length === 0 ? null : c.name
        if (label === null) return null
        return { label: c.name, key: c.id, isLeaf: false, children }
      })
      .filter(Boolean)
  }

  if (q) {
    const allNodes: any[] = []
    const cols = projectStore.collections || []
    for (const col of cols) {
      const reqs = allRequests.value
        .filter(r => r._collectionId === col.id)
        .filter(r => r.name.toLowerCase().includes(q) || r.url.toLowerCase().includes(q))
        .map(r => ({ label: r.name, key: 'req-' + r.id, isLeaf: true, _reqId: r.id, _method: r.method }))
      if (reqs.length > 0) {
        allNodes.push({ label: col.name, key: col.id, isLeaf: false, children: reqs })
      }
    }
    return allNodes
  }

  return [
    { label: t('sidebar.collections'), key: 'collections-header', isLeaf: true, disabled: true },
    ...buildTree(null),
  ]
})

function selectSection(id: string) {
  const s = sections.find(x => x.id === id)
  if (!s) return
  activeSection.value = id as any
  if (s.route && route.path !== s.route) router.push(s.route)
}

const parentOptions = computed(() => {
  const cols = projectStore.collections || []
  return cols.map(c => ({ label: c.name, value: c.id }))
})

const routeSectionMap: Record<string, string> = {
  '/workspace': 'workspace',
  '/projects': 'project',
  '/history': 'history',
  '/docs': 'docs',
  '/test-runner': 'tests',
}

watch(() => route.path, async (path) => {
  const section = routeSectionMap[path]
  if (section) activeSection.value = section as any
  if (path === '/workspace' && projectStore.currentProject) {
    const pid = projectStore.currentProject.id
    if (projectStore.collections.length === 0 || allRequests.value.length === 0) {
      await loadCollections(pid)
    }
  }
})

watch(() => projectStore.refreshKey, () => {
  if (projectStore.currentProject) {
    loadCollections(projectStore.currentProject.id)
  }
})

async function loadProjects() {
  try {
    const list = await ListProjects()
    if (!list) return
    projectStore.setProjects(list)
    if (list.length === 0) {
      const p = await CreateProject(t('project.defaultName'), '')
      if (!p) return
      projectStore.addProject(p)
      projectStore.setCurrentProject(p)
      localStorage.setItem('paw-current-project', p.id)
      await loadCollections(p.id)
    } else {
      const savedId = localStorage.getItem('paw-current-project')
      const saved = savedId ? list.find((p: any) => p.id === savedId) : null
      const target = saved || list[0]
      if (!target) return
      projectStore.setCurrentProject(target)
      await loadCollections(target.id)
    }
  } catch { console.error('Failed to load projects') }
}

async function loadCollections(projectId: string) {
  try {
    projectStore.setCollections(await ListCollections(projectId))
    await loadAllRequests()
  }
  catch { projectStore.setCollections([]); allRequests.value = [] }
  try {
    const envs = await ListEnvironments(projectId)
    envStore.setEnvironments(envs)
    const active = envs.find((e: any) => e.is_active)
    if (active) envStore.setActiveEnvironment(active)
    else envStore.setActiveEnvironment(null)
  } catch { envStore.setEnvironments([]); envStore.setActiveEnvironment(null) }
}

function startAdd(parentId?: string) {
  newCollectionName.value = ''
  selectedParentId.value = parentId || null
  showAddModal.value = true
}

async function confirmAdd() {
  if (!newCollectionName.value.trim()) return
  if (!projectStore.currentProject) { message.error(t('sidebar.noProject')); return }
  try {
    const col = await CreateCollection(projectStore.currentProject.id, selectedParentId.value || '', newCollectionName.value.trim(), 0)
    projectStore.addCollection(col)
    showAddModal.value = false
    message.success(t('sidebar.created'))
    await loadCollections(projectStore.currentProject.id)
  } catch (e: any) { message.error(e.message || t('sidebar.failedCreate')) }
}

function toggleLocale() {
  const next = locale.value === 'zh-CN' ? 'en' : 'zh-CN'
  locale.value = next
  localStorage.setItem('paw-locale', next)
}

// Drag & drop
async function handleDrop({ node, dragNode, dropPosition }: any) {
  if (!dragNode || !node || dragNode.key === node.key) return
  const draggedId = dragNode.key
  const targetId = node.key
  if (draggedId.startsWith('req-') || targetId.startsWith('req-') || targetId === 'collections-header') return
  try {
    await UpdateCollection(draggedId, '', targetId, 0)
    if (projectStore.currentProject) await loadCollections(projectStore.currentProject.id)
    message.success('Moved')
  } catch (e: any) { message.error(e.message) }
}

// Tree node click
function handleTreeSelect(keys: any) {
  if (!keys?.[0]) return
  const key: string = keys[0]
  if (key.startsWith('req-')) {
    const reqId = key.slice(4)
    const req = allRequests.value.find(r => r.id === reqId)
    if (req) {
      const existing = tabsStore.getTabByRequestId(reqId)
      if (existing) {
        tabsStore.setActiveTab(existing.id)
        return
      }
      const tabId = tabsStore.addHttpTab(req.id, req.name)
      tabsStore.updateTabData({
        method: req.method,
        url: req.url,
        headers: safeParse(req.headers, [{ key: 'Content-Type', value: 'application/json', enabled: true }]),
        params: safeParse(req.params, []),
        body: safeParseBody(req.body),
        bodyType: safeParseBodyType(req.body),
        authType: safeParseAuth(req.auth),
      })
    }
  }
}

function safeParse(str: string, fallback: any): any {
  if (!str) return fallback
  try { return JSON.parse(str) } catch { return fallback }
}

function safeParseBody(body: string): string {
  if (!body) return ''
  try {
    const obj = JSON.parse(body)
    return obj.content || obj.body || ''
  } catch { return body }
}

function safeParseBodyType(body: string): string {
  if (!body) return 'none'
  try {
    const obj = JSON.parse(body)
    return obj.body_type || obj.type || 'none'
  } catch { return 'none' }
}

function safeParseAuth(auth: string): string {
  if (!auth) return 'none'
  try {
    const obj = JSON.parse(auth)
    return obj.type || obj.auth_type || 'none'
  } catch { return 'none' }
}

// Right-click context menu
function handleContextMenu(e: MouseEvent, key: string) {
  e.preventDefault()
  contextMenuX.value = e.clientX
  contextMenuY.value = e.clientY
  showContextMenu.value = true
  if (key.startsWith('req-')) {
    contextMenuReqId.value = key.slice(4)
    contextMenuColId.value = null
  } else {
    contextMenuColId.value = key
    contextMenuReqId.value = null
  }
}

// Rename
function startRename(type: 'collection' | 'request', id: string, name: string) {
  showContextMenu.value = false
  renameTarget.value = { id, name, type }
  renameValue.value = name
  showRenameModal.value = true
}

async function confirmRename() {
  if (!renameTarget.value || !renameValue.value.trim()) return
  try {
    if (renameTarget.value.type === 'collection') {
      await UpdateCollection(renameTarget.value.id, renameValue.value.trim(), '', 0)
      if (projectStore.currentProject) await loadCollections(projectStore.currentProject.id)
    } else {
      const req = allRequests.value.find(r => r.id === renameTarget.value!.id)
      if (req) {
        await UpdateRequest(req.id, renameValue.value.trim(), req.method, req.url, req.headers, req.params, req.body, req.auth, req.script, req.sort_order)
        await loadAllRequests()
      }
    }
    showRenameModal.value = false
    message.success('Renamed')
  } catch (e: any) { message.error(e.message) }
}

// Delete
async function handleDelete(type: 'collection' | 'request', id: string) {
  showContextMenu.value = false
  dialog.warning({
    title: 'Delete',
    content: 'Are you sure?',
    positiveText: 'Delete',
    negativeText: 'Cancel',
    onPositiveClick: async () => {
      try {
        if (type === 'collection') {
          await DeleteCollection(id)
        } else {
          await DeleteRequest(id)
        }
        if (projectStore.currentProject) await loadCollections(projectStore.currentProject.id)
        message.success('Deleted')
      } catch (e: any) { message.error(e.message) }
    },
  })
}

// Copy request
async function copyRequest(reqId: string) {
  showContextMenu.value = false
  const req = allRequests.value.find(r => r.id === reqId)
  if (!req) return
  try {
    await CreateRequest(req._collectionId, req.name + ' (Copy)', req.method, req.url, req.headers, req.params, req.body, req.auth, req.script, 0)
    await loadAllRequests()
    message.success('Copied')
  } catch (e: any) { message.error(e.message) }
}

// Move request
function startMove(reqId: string) {
  showContextMenu.value = false
  moveTargetId.value = reqId
  showMoveModal.value = true
}

async function confirmMove(targetColId: string) {
  if (!moveTargetId.value) return
  const req = allRequests.value.find(r => r.id === moveTargetId.value)
  if (!req) return
  try {
    await UpdateRequest(req.id, req.name, req.method, req.url, req.headers, req.params, req.body, req.auth, req.script, 0)
    // Note: The actual collection_id move would need a separate API,
    // but we use UpdateRequest which doesn't support collection_id change.
    // For now, recreate in target collection.
    await CreateRequest(targetColId, req.name, req.method, req.url, req.headers, req.params, req.body, req.auth, req.script, 0)
    await DeleteRequest(req.id)
    await loadAllRequests()
    if (projectStore.currentProject) await loadCollections(projectStore.currentProject.id)
    showMoveModal.value = false
    message.success('Moved')
  } catch (e: any) { message.error(e.message) }
}

function renderLabel(info: { option: any }) {
  const opt = info.option
  if (opt.isLeaf || opt.disabled || opt.key === 'collections-header') {
    return opt.label
  }
  return h('div', { class: 'tree-node-label' }, [
    h('span', { class: 'tree-node-name' }, opt.label),
    h('span', { class: 'tree-node-actions' }, [
      h('button', {
        class: 'tree-node-dots',
        onClick: (e: MouseEvent) => {
          e.stopPropagation()
          collectionMenuId.value = opt.key
          collectionMenuPos.value = { x: e.clientX, y: e.clientY }
        },
      }, '···'),
    ]),
  ])
}

const colMenuOptions = [
  { label: 'New Request', key: 'new-request' },
]

function handleColMenuSelect(key: string) {
  const colId = collectionMenuId.value
  collectionMenuId.value = null
  if (key === 'new-request' && colId) {
    createRequestInCollection(colId)
  }
}

async function createRequestInCollection(colId: string) {
  try {
    const r = await CreateRequest(colId, 'New Request', 'GET', '', '[]', '[]', '{}', '{}', '', 0)
    if (r) {
      const tabId = tabsStore.addHttpTab(r.id, r.name)
      message.success('Request created')
      if (projectStore.currentProject) await loadCollections(projectStore.currentProject.id)
    }
  } catch (e: any) { message.error(e.message || 'Failed to create request') }
}

onMounted(loadProjects)
</script>

<template>
  <div class="sidebar-container">
    <div class="icon-bar">
      <div class="icon-bar-top">
        <div v-for="s in sections" :key="s.id" class="icon-item" :class="{ active: activeSection === s.id }" @click="selectSection(s.id)">
          <NIcon size="22"><component :is="s.icon" /></NIcon>
          <span class="icon-label">{{ $t(s.labelKey) }}</span>
        </div>
      </div>
    </div>

    <div class="content-panel">
      <div v-show="activeSection === 'workspace'" class="panel-section">
        <div class="panel-header">
          <div class="panel-project-name">{{ projectStore.currentProject?.name || $t('project.noProject') }}</div>
          <span class="panel-title">{{ $t('sidebar.collections') }}</span>
        </div>
        <div class="panel-search">
          <NInput v-model:value="searchQuery" :placeholder="$t('history.search')" size="tiny" clearable>
            <template #prefix>
              <NIcon size="14"><Search /></NIcon>
            </template>
          </NInput>
        </div>
        <div class="panel-env"><EnvSelector @manage="showEnvManager = true" /></div>
        <div class="panel-tree-area">
          <div class="tree-top-bar">
            <span></span>
            <NButton quaternary circle size="tiny" @click="() => startAdd()" class="tree-add-btn">
              <template #icon><NIcon size="14"><Add /></NIcon></template>
            </NButton>
          </div>
          <NTree
            :data="treeData"
            :default-expand-all="true"
            block-line
            selectable
            draggable
            :render-label="renderLabel"
            @update:selected-keys="handleTreeSelect"
            @drop="handleDrop"
            @contextmenu="(e: any, opt: any) => handleContextMenu(e, opt?.key)"
          />
        </div>
      </div>

      <div v-show="activeSection === 'settings'" class="panel-section">
        <div class="panel-header"><span class="panel-title">Settings</span></div>
        <div class="panel-settings">
          <div class="setting-row"><span>Dark Mode</span><NButton size="tiny" @click="toggleColorMode">{{ colorMode === 'light' ? 'Dark' : 'Light' }}</NButton></div>
          <div class="setting-row"><span>Theme Color</span><NSelect :options="themeColorOptions" :value="themeColor" size="tiny" style="width:100px" @update:value="setThemeColor" /></div>
          <div class="setting-row"><span>Language</span><NButton size="tiny" @click="toggleLocale">{{ locale === 'zh-CN' ? 'English' : '中文' }}</NButton></div>
        </div>
      </div>
    </div>

    <!-- Context Menu -->
    <div v-if="collectionMenuId" class="collection-menu" :style="{ left: collectionMenuPos.x + 'px', top: collectionMenuPos.y + 'px' }">
      <div v-for="opt in colMenuOptions" :key="opt.key" class="collection-menu-item" @click="handleColMenuSelect(opt.key)">{{ opt.label }}</div>
    </div>
    <div v-if="collectionMenuId" class="context-overlay" @click="collectionMenuId = null" @contextmenu.prevent="collectionMenuId = null" />

    <div v-if="showContextMenu" class="context-menu" :style="{ left: contextMenuX + 'px', top: contextMenuY + 'px' }">
      <div v-if="contextMenuColId" class="context-item" @click="startRename('collection', contextMenuColId!, '')">Rename Collection</div>
      <div v-if="contextMenuColId" class="context-item" @click="() => startAdd(contextMenuColId ?? undefined)">New Sub-collection</div>
      <div v-if="contextMenuColId" class="context-item" @click="handleDelete('collection', contextMenuColId!)">Delete Collection</div>
      <div v-if="contextMenuReqId" class="context-item" @click="startRename('request', contextMenuReqId!, '')">Rename</div>
      <div v-if="contextMenuReqId" class="context-item" @click="copyRequest(contextMenuReqId!)">Copy Request</div>
      <div v-if="contextMenuReqId" class="context-item" @click="startMove(contextMenuReqId!)">Move to Collection</div>
      <div v-if="contextMenuReqId" class="context-item" @click="handleDelete('request', contextMenuReqId!)">Delete</div>
    </div>
    <div v-if="showContextMenu" class="context-overlay" @click="showContextMenu = false" @contextmenu.prevent="showContextMenu = false" />

    <!-- Modals -->
    <NModal v-model:show="showAddModal" :title="$t('sidebar.newCollection')" preset="card" style="width:360px">
      <NForm>
        <NFormItem :label="$t('sidebar.collectionName')">
          <NInput v-model:value="newCollectionName" :placeholder="$t('sidebar.collectionName')" />
        </NFormItem>
        <NFormItem label="Parent">
          <NSelect
            v-model:value="selectedParentId"
            :options="parentOptions"
            size="small"
            clearable
            :placeholder="'None (root level)'"
          />
        </NFormItem>
      </NForm>
      <template #footer><NSpace justify="end"><NButton @click="showAddModal = false">{{ $t('sidebar.cancel') }}</NButton><NButton type="primary" @click="confirmAdd">{{ $t('sidebar.create') }}</NButton></NSpace></template>
    </NModal>

    <EnvManager v-model:show="showEnvManager" />
  </div>
</template>

<style scoped>
.sidebar-container { display: flex; height: 100%; width: 100%; overflow: hidden; }
.icon-bar { width: 64px; min-width: 64px; display: flex; flex-direction: column; align-items: center; padding: 8px 0; border-right: 1px solid var(--border-color); background: var(--tab-bar-bg); }
.icon-bar-top { display: flex; flex-direction: column; align-items: center; gap: 2px; width: 100%; padding: 0 4px; }
.icon-item { width: 56px; padding: 6px 0 4px; display: flex; flex-direction: column; align-items: center; gap: 2px; border-radius: 8px; cursor: pointer; color: #888; transition: all 0.15s; }
.icon-item:hover { background: var(--hover-color); color: #333; }
.icon-item.active { background: var(--active-color); color: #18a058; }
.icon-label { font-size: 10px; line-height: 1; white-space: nowrap; }
.content-panel { flex: 1; display: flex; flex-direction: column; overflow: hidden; min-width: 0; }
.panel-section { display: flex; flex-direction: column; height: 100%; }
.panel-header { padding: 12px 12px 8px; border-bottom: 1px solid var(--border-color); }
.panel-project-name { font-size: 14px; font-weight: 700; margin-bottom: 2px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.panel-title { font-size: 11px; font-weight: 600; color: #888; text-transform: uppercase; letter-spacing: 0.5px; }
.tree-top-bar { display: flex; align-items: center; justify-content: flex-end; padding: 2px 4px 0 0; }
.tree-add-btn { flex-shrink: 0; }
.panel-search { padding: 6px 8px; border-bottom: 1px solid var(--border-color); }
.panel-env { padding: 6px 8px; border-bottom: 1px solid var(--border-color); }
.panel-footer { padding: 6px 8px; border-top: 1px solid var(--border-color); }
.panel-tree-area { flex: 1; overflow-y: auto; }
.panel-list { flex: 1; overflow-y: auto; padding: 4px 0; }
.panel-list-item { padding: 6px 12px; cursor: pointer; font-size: 13px; }
.panel-list-item:hover { background: var(--hover-color); }
.panel-list-item.active { background: var(--active-color); font-weight: 600; }
.panel-settings { padding: 12px; display: flex; flex-direction: column; gap: 12px; }
.setting-row { display: flex; align-items: center; justify-content: space-between; font-size: 13px; }
.tree-node-label { display: flex; align-items: center; width: 100%; padding-right: 8px; }
.tree-node-name { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; flex: 1; min-width: 0; }
.tree-node-actions { display: none; flex-shrink: 0; margin-left: auto; }
.tree-node-label:hover .tree-node-actions { display: flex; }
.tree-node-dots { background: none; border: none; cursor: pointer; font-size: 16px; line-height: 1; padding: 0 4px; color: #888; border-radius: 4px; letter-spacing: 1px; }
.tree-node-dots:hover { background: #e8e8e8; color: #333; }
.collection-menu { position: fixed; z-index: 9999; background: #fff; border: 1px solid #ddd; border-radius: 6px; box-shadow: 0 4px 12px rgba(0,0,0,0.15); padding: 4px 0; min-width: 140px; }
.collection-menu-item { padding: 6px 12px; font-size: 13px; cursor: pointer; }
.collection-menu-item:hover { background: #f0f0f0; }
.context-menu { position: fixed; z-index: 9999; background: #fff; border: 1px solid #ddd; border-radius: 6px; box-shadow: 0 4px 12px rgba(0,0,0,0.15); padding: 4px 0; min-width: 160px; }
.context-item { padding: 6px 12px; font-size: 13px; cursor: pointer; }
.context-item:hover { background: #f0f0f0; }
.context-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; z-index: 9998; }
.import-format { margin-bottom: 8px; }
.import-textarea { font-family: 'Cascadia Code', 'Fira Code', 'Consolas', monospace; font-size: 13px; }
.move-list { max-height: 300px; overflow-y: auto; }
.move-item { padding: 8px 12px; cursor: pointer; font-size: 13px; border-bottom: 1px solid var(--border-color); }
.move-item:hover { background: var(--hover-color); }
</style>
