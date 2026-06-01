<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue'
import { NTree, NButton, NIcon, NInput, NModal, NForm, NFormItem, NSpace, NSelect, useMessage } from 'naive-ui'
import { Add, CodeSlash, FolderOpen, Time, DocumentText, Bug, Settings as SettingsIcon } from '@vicons/ionicons5'
import { useRouter, useRoute } from 'vue-router'
import { useProjectStore } from '../stores/project'
import { useI18n } from 'vue-i18n'
import { useTheme } from '../composables/useTheme'
import { t, setLocale as setI18nLocale } from '../i18n'
import { CreateProject, ListProjects } from '../../wailsjs/go/handlers/ProjectHandler'
import { CreateCollection, ListCollections } from '../../wailsjs/go/handlers/CollectionHandler'
import EnvSelector from './EnvSelector.vue'
import EnvManager from './EnvManager.vue'

const router = useRouter()
const route = useRoute()
const projectStore = useProjectStore()
const { toggleColorMode, setThemeColor, themeColor, colorMode } = useTheme()
const { locale } = useI18n()
const message = useMessage()

const activeSection = ref<'workspace' | 'project' | 'history' | 'docs' | 'tests' | 'settings'>('workspace')
const showAddModal = ref(false)
const newCollectionName = ref('')
const showEnvManager = ref(false)
const showProjectModal = ref(false)
const newProjectName = ref('')

const sections = [
  { id: 'workspace', labelKey: 'sidebar.workspaceLabel', icon: CodeSlash, route: '/workspace' },
  { id: 'project', labelKey: 'sidebar.projectLabel', icon: FolderOpen, route: '' },
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

const treeData = computed(() => {
  if (activeSection.value !== 'workspace') return []
  function buildTree(parentId: string | null): any[] {
    return (projectStore.collections || [])
      .filter(c => c.parent_id === parentId)
      .sort((a, b) => a.sort_order - b.sort_order)
      .map(c => ({ label: c.name, key: c.id, isLeaf: false, children: buildTree(c.id) }))
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
  try { projectStore.setCollections(await ListCollections(projectId)) }
  catch { projectStore.setCollections([]) }
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
  } catch (e: any) { message.error(e.message || t('sidebar.failedCreate')) }
}

async function createProject() {
  if (!newProjectName.value.trim()) return
  try {
    const p = await CreateProject(newProjectName.value.trim(), '')
    projectStore.addProject(p)
    projectStore.setCurrentProject(p)
    projectStore.setCollections([])
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
        <div class="panel-header"><span class="panel-title">{{ $t('sidebar.collections') }}</span></div>
        <div class="panel-env"><EnvSelector @manage="showEnvManager = true" /></div>
        <div class="panel-tree-area"><NTree :data="treeData" :default-expand-all="true" block-line selectable @update:selected-keys="(keys: any) => keys?.[0] && projectStore.selectCollection(keys[0])" /></div>
        <div class="panel-footer">
          <NButton size="tiny" quaternary @click="startAdd"><template #icon><NIcon><Add /></NIcon></template>{{ $t('sidebar.newCollection') }}</NButton>
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

    <NModal v-model:show="showAddModal" :title="$t('sidebar.newCollection')" preset="card" style="width:360px">
      <NForm><NFormItem :label="$t('sidebar.collectionName')"><NInput v-model:value="newCollectionName" :placeholder="$t('sidebar.collectionName')" /></NFormItem></NForm>
      <template #footer><NSpace justify="end"><NButton @click="showAddModal = false">{{ $t('sidebar.cancel') }}</NButton><NButton type="primary" @click="confirmAdd">{{ $t('sidebar.create') }}</NButton></NSpace></template>
    </NModal>

    <NModal v-model:show="showProjectModal" :title="$t('project.newProject')" preset="card" style="width:360px">
      <NForm><NFormItem :label="$t('project.projectName')"><NInput v-model:value="newProjectName" :placeholder="$t('project.projectName')" /></NFormItem></NForm>
      <template #footer><NSpace justify="end"><NButton @click="showProjectModal = false">{{ $t('common.cancel') }}</NButton><NButton type="primary" @click="createProject">{{ $t('project.create') }}</NButton></NSpace></template>
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
.panel-env { padding: 6px 8px; border-bottom: 1px solid var(--border-color); }
.panel-footer { padding: 6px 8px; border-top: 1px solid var(--border-color); }
.panel-tree-area { flex: 1; overflow-y: auto; }
.panel-list { flex: 1; overflow-y: auto; padding: 4px 0; }
.panel-list-item { padding: 6px 12px; cursor: pointer; font-size: 13px; }
.panel-list-item:hover { background: var(--hover-color); }
.panel-list-item.active { background: var(--active-color); font-weight: 600; }
.panel-settings { padding: 12px; display: flex; flex-direction: column; gap: 12px; }
.setting-row { display: flex; align-items: center; justify-content: space-between; font-size: 13px; }
</style>
