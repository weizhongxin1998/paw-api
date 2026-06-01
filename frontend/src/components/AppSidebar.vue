<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue'
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
import { ImportPostman, ImportSwagger, ImportCurl } from '../../wailsjs/go/handlers/ImporterHandler'
import { ExportPostman, ExportSwagger } from '../../wailsjs/go/handlers/ExporterHandler'
import EnvSelector from './EnvSelector.vue'
import EnvManager from './EnvManager.vue'

const router = useRouter()
const route = useRoute()
const projectStore = useProjectStore()
const tabsStore = useTabsStore()
const { toggleColorMode, setThemeColor, themeColor, colorMode } = useTheme()
const { locale } = useI18n()
const message = useMessage()
const dialog = useDialog()

const activeSection = ref<'workspace' | 'project' | 'history' | 'docs' | 'tests' | 'settings'>('workspace')
const showAddModal = ref(false)
const newCollectionName = ref('')
const showEnvManager = ref(false)
const showProjectModal = ref(false)
const newProjectName = ref('')
const searchQuery = ref('')
const showImportModal = ref(false)
const importFileContent = ref('')
const importFormat = ref('postman')
const showCurlModal = ref(false)
const curlCommand = ref('')
const showExportModal = ref(false)
const exportFormat = ref('postman')
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
      await loadCollections(p.id)
    } else {
      const p = list[0]
      if (!p) return
      projectStore.setCurrentProject(p)
      await loadCollections(p.id)
    }
  } catch { console.error('Failed to load projects') }
}

async function loadCollections(projectId: string) {
  try {
    projectStore.setCollections(await ListCollections(projectId))
    await loadAllRequests()
  }
  catch { projectStore.setCollections([]); allRequests.value = [] }
}

async function switchProject(projectId: string) {
  const p = projectStore.projects.find(x => x.id === projectId)
  if (p) { projectStore.setCurrentProject(p); await loadCollections(p.id) }
}

function startAdd() { newCollectionName.value = ''; showAddModal.value = true }

async function confirmAdd() {
  if (!newCollectionName.value.trim()) return
  if (!projectStore.currentProject) { message.error(t('sidebar.noProject')); return }
  try {
    const col = await CreateCollection(projectStore.currentProject.id, '', newCollectionName.value.trim(), 0)
    projectStore.addCollection(col)
    showAddModal.value = false
    message.success(t('sidebar.created'))
    await loadCollections(projectStore.currentProject.id)
  } catch (e: any) { message.error(e.message || t('sidebar.failedCreate')) }
}

async function createProject() {
  if (!newProjectName.value.trim()) return
  try {
    const p = await CreateProject(newProjectName.value.trim(), '')
    projectStore.addProject(p)
    projectStore.setCurrentProject(p)
    projectStore.setCollections([])
    allRequests.value = []
    showProjectModal.value = false
    newProjectName.value = ''
    message.success(t('project.created'))
  } catch (e: any) { message.error(e.message || t('project.failedCreate')) }
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

// Import
async function handleImport() {
  try {
    let result
    if (importFormat.value === 'postman') {
      result = await ImportPostman(importFileContent.value)
    } else {
      result = await ImportSwagger(importFileContent.value)
    }
    if (result) {
      message.success(`Imported ${result.requests?.length || 0} requests`)
      if (projectStore.currentProject) await loadCollections(projectStore.currentProject.id)
    }
    showImportModal.value = false
  } catch (e: any) { message.error(e.message || 'Import failed') }
}

async function handleImportCurl() {
  if (!curlCommand.value.trim()) return
  try {
    const result = await ImportCurl(curlCommand.value.trim())
    if (result && projectStore.currentProject) {
      const tabId = tabsStore.addHttpTab(undefined, result.name || 'cURL Import')
      tabsStore.updateTabData({
        method: result.method || 'GET',
        url: result.url || '',
        headers: safeParse(result.headers, []),
        params: safeParse(result.params, []),
        body: result.body || '',
        bodyType: result.body ? 'json' : 'none',
      })
      message.success('cURL imported')
    }
    showCurlModal.value = false
    curlCommand.value = ''
  } catch (e: any) { message.error(e.message || 'Import failed') }
}

// Export
async function handleExport() {
  if (!projectStore.currentProject) return
  try {
    const cols = JSON.stringify(projectStore.collections || [])
    const reqs = JSON.stringify(allRequests.value || [])
    let result: string
    if (exportFormat.value === 'postman') {
      result = await ExportPostman(cols, reqs, projectStore.currentProject.name)
    } else {
      result = await ExportSwagger(cols, reqs, projectStore.currentProject.name)
    }
    const blob = new Blob([result], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${projectStore.currentProject.name}.${exportFormat.value === 'postman' ? 'postman_collection.json' : 'swagger.json'}`
    a.click()
    URL.revokeObjectURL(url)
    showExportModal.value = false
    message.success('Exported')
  } catch (e: any) { message.error(e.message || 'Export failed') }
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
          <NTree
            :data="treeData"
            :default-expand-all="true"
            block-line
            selectable
            draggable
            @update:selected-keys="handleTreeSelect"
            @drop="handleDrop"
            @contextmenu="(e: any, opt: any) => handleContextMenu(e, opt?.key)"
          />
        </div>
        <div class="panel-footer">
          <NSpace size="small">
            <NButton size="tiny" quaternary @click="startAdd"><template #icon><NIcon size="14"><Add /></NIcon></template>{{ $t('sidebar.newCollection') }}</NButton>
            <NButton size="tiny" quaternary @click="showImportModal = true"><template #icon><NIcon size="14"><Download /></NIcon></template>Import</NButton>
            <NButton size="tiny" quaternary @click="showExportModal = true"><template #icon><NIcon size="14"><Download /></NIcon></template>Export</NButton>
            <NButton size="tiny" quaternary @click="showCurlModal = true"><template #icon><NIcon size="14"><CodeSlash /></NIcon></template>cURL</NButton>
          </NSpace>
        </div>
      </div>

      <div v-show="activeSection === 'project'" class="panel-section">
        <div class="panel-header"><span class="panel-title">{{ $t('project.select') }}</span></div>
        <div class="panel-list">
          <div v-for="p in projectStore.projects" :key="p.id" class="panel-list-item" :class="{ active: projectStore.currentProject?.id === p.id }" @click="switchProject(p.id)">{{ p.name }}</div>
        </div>
        <div class="panel-footer">
          <NButton size="tiny" quaternary @click="showProjectModal = true"><template #icon><NIcon><Add /></NIcon></template>{{ $t('project.newProject') }}</NButton>
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
    <div v-if="showContextMenu" class="context-menu" :style="{ left: contextMenuX + 'px', top: contextMenuY + 'px' }">
      <div v-if="contextMenuColId" class="context-item" @click="startRename('collection', contextMenuColId!, '')">Rename Collection</div>
      <div v-if="contextMenuColId" class="context-item" @click="handleDelete('collection', contextMenuColId!)">Delete Collection</div>
      <div v-if="contextMenuReqId" class="context-item" @click="startRename('request', contextMenuReqId!, '')">Rename</div>
      <div v-if="contextMenuReqId" class="context-item" @click="copyRequest(contextMenuReqId!)">Copy Request</div>
      <div v-if="contextMenuReqId" class="context-item" @click="startMove(contextMenuReqId!)">Move to Collection</div>
      <div v-if="contextMenuReqId" class="context-item" @click="handleDelete('request', contextMenuReqId!)">Delete</div>
    </div>
    <div v-if="showContextMenu" class="context-overlay" @click="showContextMenu = false" @contextmenu.prevent="showContextMenu = false" />

    <!-- Modals -->
    <NModal v-model:show="showAddModal" :title="$t('sidebar.newCollection')" preset="card" style="width:360px">
      <NForm><NFormItem :label="$t('sidebar.collectionName')"><NInput v-model:value="newCollectionName" :placeholder="$t('sidebar.collectionName')" /></NFormItem></NForm>
      <template #footer><NSpace justify="end"><NButton @click="showAddModal = false">{{ $t('sidebar.cancel') }}</NButton><NButton type="primary" @click="confirmAdd">{{ $t('sidebar.create') }}</NButton></NSpace></template>
    </NModal>

    <NModal v-model:show="showProjectModal" :title="$t('project.newProject')" preset="card" style="width:360px">
      <NForm><NFormItem :label="$t('project.projectName')"><NInput v-model:value="newProjectName" :placeholder="$t('project.projectName')" /></NFormItem></NForm>
      <template #footer><NSpace justify="end"><NButton @click="showProjectModal = false">{{ $t('common.cancel') }}</NButton><NButton type="primary" @click="createProject">{{ $t('project.create') }}</NButton></NSpace></template>
    </NModal>

    <NModal v-model:show="showRenameModal" title="Rename" preset="card" style="width:360px">
      <NForm><NFormItem label="Name"><NInput v-model:value="renameValue" /></NFormItem></NForm>
      <template #footer><NSpace justify="end"><NButton @click="showRenameModal = false">Cancel</NButton><NButton type="primary" @click="confirmRename">Rename</NButton></NSpace></template>
    </NModal>

    <NModal v-model:show="showMoveModal" title="Move to Collection" preset="card" style="width:400px">
      <div class="move-list">
        <div v-for="col in projectStore.collections" :key="col.id" class="move-item" @click="confirmMove(col.id)">
          {{ col.name }}
        </div>
      </div>
      <template #footer><NSpace justify="end"><NButton @click="showMoveModal = false">Cancel</NButton></NSpace></template>
    </NModal>

    <NModal v-model:show="showImportModal" title="Import" preset="card" style="width:500px">
      <NSelect v-model:value="importFormat" :options="[{ label: 'Postman v2.1', value: 'postman' }, { label: 'OpenAPI/Swagger', value: 'swagger' }]" size="small" class="import-format" />
      <NInput v-model:value="importFileContent" type="textarea" :rows="10" placeholder="Paste JSON content here..." class="import-textarea" />
      <template #footer><NSpace justify="end"><NButton @click="showImportModal = false">Cancel</NButton><NButton type="primary" @click="handleImport">Import</NButton></NSpace></template>
    </NModal>

    <NModal v-model:show="showCurlModal" title="Import cURL" preset="card" style="width:500px">
      <NInput v-model:value="curlCommand" type="textarea" :rows="4" placeholder="Paste cURL command..." class="import-textarea" />
      <template #footer><NSpace justify="end"><NButton @click="showCurlModal = false">Cancel</NButton><NButton type="primary" @click="handleImportCurl">Import</NButton></NSpace></template>
    </NModal>

    <NModal v-model:show="showExportModal" title="Export" preset="card" style="width:360px">
      <NSelect v-model:value="exportFormat" :options="[{ label: 'Postman v2.1', value: 'postman' }, { label: 'OpenAPI/Swagger', value: 'swagger' }]" size="small" />
      <template #footer><NSpace justify="end"><NButton @click="showExportModal = false">Cancel</NButton><NButton type="primary" @click="handleExport">Export</NButton></NSpace></template>
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
.panel-title { font-size: 13px; font-weight: 600; color: #888; text-transform: uppercase; letter-spacing: 0.5px; }
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
