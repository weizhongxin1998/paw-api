<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { NCard, NIcon, NButton, NModal, NForm, NFormItem, NInput, NSpace, NEmpty, NSpin, NTag, NSelect, useMessage } from 'naive-ui'
import { FolderOpen, Add, Download, CodeSlash } from '@vicons/ionicons5'
import { useRouter } from 'vue-router'
import { useProjectStore } from '../stores/project'
import { useTabsStore } from '../stores/tabs'
import { ListProjects, CreateProject } from '../../wailsjs/go/handlers/ProjectHandler'
import { ListCollections } from '../../wailsjs/go/handlers/CollectionHandler'
import { ListRequests } from '../../wailsjs/go/handlers/RequestHandler'
import { ListEnvironments } from '../../wailsjs/go/handlers/EnvironmentHandler'
import { useEnvironmentStore } from '../stores/environment'
import { ImportPostman, ImportSwagger, ImportCurl } from '../../wailsjs/go/handlers/ImporterHandler'
import { ExportPostman, ExportSwagger } from '../../wailsjs/go/handlers/ExporterHandler'

const router = useRouter()
const projectStore = useProjectStore()
const tabsStore = useTabsStore()
const envStore = useEnvironmentStore()
const message = useMessage()

const showImportModal = ref(false)
const importFileContent = ref('')
const importFormat = ref('postman')
const showCurlModal = ref(false)
const curlCommand = ref('')
const showExportModal = ref(false)
const exportFormat = ref('postman')

const loading = ref(true)
const showNewModal = ref(false)
const newName = ref('')
const newDesc = ref('')
const projectStats = ref<Record<string, { collections: number; requests: number }>>({})

async function loadAll() {
  loading.value = true
  try {
    const list = await ListProjects()
    projectStore.setProjects(list || [])
    const stats: Record<string, { collections: number; requests: number }> = {}
    for (const p of list || []) {
      const cols = await ListCollections(p.id)
      let totalReqs = 0
      for (const c of cols || []) {
        const reqs = await ListRequests(c.id)
        totalReqs += reqs?.length || 0
      }
      stats[p.id] = { collections: cols?.length || 0, requests: totalReqs }
    }
    projectStats.value = stats
  } catch {
    message.error('Failed to load projects')
  } finally {
    loading.value = false
  }
}

async function selectProject(p: any) {
  projectStore.setCurrentProject(p)
  try {
    const cols = await ListCollections(p.id)
    projectStore.setCollections(cols || [])
  } catch {
    projectStore.setCollections([])
  }
  try {
    const envs = await ListEnvironments(p.id)
    envStore.setEnvironments(envs)
    const active = envs.find((e: any) => e.is_active)
    if (active) envStore.setActiveEnvironment(active)
    else envStore.setActiveEnvironment(null)
  } catch {
    envStore.setEnvironments([])
    envStore.setActiveEnvironment(null)
  }
  router.push('/workspace')
}

async function createProject() {
  if (!newName.value.trim()) return
  try {
    const p = await CreateProject(newName.value.trim(), newDesc.value.trim())
    if (p) {
      projectStore.addProject(p)
      projectStats.value[p.id] = { collections: 0, requests: 0 }
      message.success('Project created')
      showNewModal.value = false
      newName.value = ''
      newDesc.value = ''
    }
  } catch (e: any) {
    message.error(e.message || 'Failed to create project')
  }
}

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
    }
    showImportModal.value = false
    importFileContent.value = ''
    await loadAll()
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
      router.push('/workspace')
    }
    showCurlModal.value = false
    curlCommand.value = ''
  } catch (e: any) { message.error(e.message || 'Import failed') }
}

function safeParse(str: string, fallback: any): any {
  if (!str) return fallback
  try { return JSON.parse(str) } catch { return fallback }
}

async function handleExport() {
  if (!projectStore.currentProject) return
  try {
    const cols = JSON.stringify(projectStore.collections || [])
    const reqsJSON = JSON.stringify([])
    let result: string
    if (exportFormat.value === 'postman') {
      result = await ExportPostman(cols, reqsJSON, projectStore.currentProject.name)
    } else {
      result = await ExportSwagger(cols, reqsJSON, projectStore.currentProject.name)
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

onMounted(loadAll)
</script>

<template>
  <div class="gallery">
    <div class="gallery-header">
      <h2 class="gallery-title">{{ $t('project.select') }}</h2>
      <NSpace size="small">
        <NButton size="tiny" quaternary @click="showImportModal = true">
          <template #icon><NIcon size="14"><Download /></NIcon></template>Import
        </NButton>
        <NButton size="tiny" quaternary @click="showExportModal = true">
          <template #icon><NIcon size="14"><Download /></NIcon></template>Export
        </NButton>
        <NButton size="tiny" quaternary @click="showCurlModal = true">
          <template #icon><NIcon size="14"><CodeSlash /></NIcon></template>cURL
        </NButton>
        <NButton size="small" type="primary" @click="showNewModal = true">
          <template #icon><NIcon><Add /></NIcon></template>
          {{ $t('project.newProject') }}
        </NButton>
      </NSpace>
    </div>

    <NSpin :show="loading">
      <div v-if="projectStore.projects.length === 0 && !loading" class="gallery-empty">
        <NEmpty :description="$t('project.noProject')">
          <template #extra>
            <NButton size="small" @click="showNewModal = true">{{ $t('project.newProject') }}</NButton>
          </template>
        </NEmpty>
      </div>
      <div v-else class="gallery-grid">
        <NCard
          v-for="p in projectStore.projects"
          :key="p.id"
          class="project-card"
          hoverable
          @click="selectProject(p)"
          :class="{ active: projectStore.currentProject?.id === p.id }"
        >
          <div class="card-body">
            <div class="card-icon">
              <NIcon size="48"><FolderOpen /></NIcon>
            </div>
            <div class="card-name">{{ p.name }}</div>
            <div v-if="p.description" class="card-desc">{{ p.description }}</div>
            <div class="card-stats">
              <NTag size="tiny">{{ projectStats[p.id]?.collections ?? 0 }} collections</NTag>
              <NTag size="tiny">{{ projectStats[p.id]?.requests ?? 0 }} requests</NTag>
            </div>
          </div>
        </NCard>
      </div>
    </NSpin>

    <NModal v-model:show="showNewModal" :title="$t('project.newProject')" preset="card" style="width:400px">
      <NForm>
        <NFormItem :label="$t('project.projectName')">
          <NInput v-model:value="newName" :placeholder="$t('project.projectName')" />
        </NFormItem>
        <NFormItem label="Description">
          <NInput v-model:value="newDesc" placeholder="Optional description" type="textarea" :rows="2" />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace justify="end">
          <NButton @click="showNewModal = false">{{ $t('common.cancel') }}</NButton>
          <NButton type="primary" @click="createProject">{{ $t('project.create') }}</NButton>
        </NSpace>
      </template>
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
  </div>
</template>

<style scoped>
.gallery {
  padding: 32px 40px;
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}
.gallery-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}
.gallery-title {
  font-size: 22px;
  font-weight: 700;
  margin: 0;
}
.gallery-empty {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}
.gallery-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
}
.project-card {
  cursor: pointer;
  transition: transform 0.15s, box-shadow 0.15s;
}
.project-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0,0,0,0.12);
}
.project-card.active {
  border-color: #18a058;
  box-shadow: 0 0 0 2px rgba(24,160,88,0.2);
}
.card-body {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  padding: 12px 0;
}
.card-icon {
  margin-bottom: 12px;
  color: #18a058;
}
.card-name {
  font-size: 15px;
  font-weight: 600;
  margin-bottom: 4px;
  word-break: break-all;
}
.card-desc {
  font-size: 12px;
  color: #999;
  margin-bottom: 8px;
  line-height: 1.4;
}
.card-stats {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
  justify-content: center;
}
.import-format { margin-bottom: 8px; }
.import-textarea { font-family: 'Cascadia Code', 'Fira Code', 'Consolas', monospace; font-size: 13px; }
</style>
