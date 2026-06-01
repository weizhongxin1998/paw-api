<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { NButton, NTag, NIcon, NEmpty } from 'naive-ui'
import { Download } from '@vicons/ionicons5'
import { useProjectStore } from '../stores/project'
import { ListCollections } from '../../wailsjs/go/handlers/CollectionHandler'
import { ListRequests } from '../../wailsjs/go/handlers/RequestHandler'

const { t } = useI18n()
const projectStore = useProjectStore()
const docs = ref<Array<{ collection: any; requests: any[] }>>([])

async function loadDocs() {
  if (!projectStore.currentProject) return
  const cols = await ListCollections(projectStore.currentProject.id)
  const result: Array<{ collection: any; requests: any[] }> = []
  for (const col of cols) {
    const reqs = await ListRequests(col.id)
    result.push({ collection: col, requests: reqs })
  }
  docs.value = result
}

function formatJSON(text: string): string { try { return JSON.stringify(JSON.parse(text), null, 2) } catch { return text } }

function methodType(method: string): 'success' | 'info' | 'warning' | 'error' | 'default' {
  const map: Record<string, any> = { GET: 'success', POST: 'info', PUT: 'warning', DELETE: 'error', PATCH: 'info' }
  return map[method] ?? 'default'
}

function exportHTML() {
  let html = `<!DOCTYPE html><html><head><meta charset="utf-8"><title>API Docs</title>
<style>body{font-family:system-ui,sans-serif;max-width:800px;margin:0 auto;padding:2rem}
h1{color:#333} h2{color:#555;border-bottom:1px solid #ddd;padding-bottom:4px}
.endpoint{margin:1rem 0;padding:.5rem 1rem;background:#f5f5f5;border-radius:4px}
.method{display:inline-block;padding:2px 8px;border-radius:3px;color:#fff;font-size:12px;font-weight:600}
.method-GET{background:#18a058} .method-POST{background:#2080f0} .method-PUT{background:#f0a020} .method-DELETE{background:#e74c3c}
.url{font-family:monospace;margin-left:8px}
pre{background:#1e1e1e;color:#d4d4d4;padding:8px;border-radius:4px;overflow-x:auto}
</style></head><body><h1>API Documentation</h1>`
  for (const { collection, requests } of docs.value) {
    html += `<h2>${collection.name}</h2>`
    for (const req of requests) {
      html += `<div class="endpoint"><span class="method method-${req.method}">${req.method}</span><span class="url">${req.url}</span>
<p><strong>${req.name}</strong></p>
${req.params && req.params !== '{}' ? `<pre>Params:\n${formatJSON(req.params)}</pre>` : ''}
${req.headers && req.headers !== '{}' ? `<pre>Headers:\n${formatJSON(req.headers)}</pre>` : ''}
${req.body && req.body !== '{}' ? `<pre>Body:\n${formatJSON(req.body)}</pre>` : ''}</div>`
    }
  }
  html += '</body></html>'
  const blob = new Blob([html], { type: 'text/html' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url; a.download = 'api-docs.html'; a.click()
  URL.revokeObjectURL(url)
}

onMounted(loadDocs)
</script>

<template>
  <div class="docs-view">
    <div class="toolbar">
      <h2 class="title">{{ t('docs.title') }}</h2>
      <NButton size="small" @click="exportHTML" :disabled="docs.length === 0">
        <template #icon><NIcon><Download /></NIcon></template>
        {{ t('docs.export') }}
      </NButton>
    </div>
    <div v-if="docs.length === 0" class="empty"><NEmpty :description="t('docs.empty')" /></div>
    <div v-else class="doc-content">
      <div v-for="{ collection, requests } in docs" :key="collection.id" class="doc-section">
        <h3 class="collection-title">{{ collection.name }}</h3>
        <div v-for="req in requests" :key="req.id" class="doc-endpoint">
          <div class="endpoint-header">
            <NTag size="small" :type="methodType(req.method)" class="method-tag">{{ req.method }}</NTag>
            <span class="endpoint-url">{{ req.url }}</span>
          </div>
          <p class="endpoint-name">{{ req.name }}</p>
          <div v-if="req.params && req.params !== '{}'" class="doc-block"><span class="block-label">{{ t('docs.params') }}</span><pre>{{ formatJSON(req.params) }}</pre></div>
          <div v-if="req.headers && req.headers !== '{}'" class="doc-block"><span class="block-label">{{ t('docs.headers') }}</span><pre>{{ formatJSON(req.headers) }}</pre></div>
          <div v-if="req.body && req.body !== '{}'" class="doc-block"><span class="block-label">{{ t('docs.body') }}</span><pre>{{ formatJSON(req.body) }}</pre></div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.docs-view { padding: 16px 20px; height: 100%; display: flex; flex-direction: column; }
.toolbar { display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px; }
.title { font-size: 18px; font-weight: 600; }
.doc-content { flex: 1; overflow-y: auto; }
.doc-section { margin-bottom: 24px; }
.collection-title { font-size: 16px; font-weight: 600; border-bottom: 1px solid var(--border-color); padding-bottom: 4px; margin-bottom: 12px; }
.doc-endpoint { padding: 8px 12px; margin-bottom: 8px; background: var(--tab-bar-bg); border-radius: 4px; border: 1px solid var(--border-color); }
.endpoint-header { display: flex; align-items: center; gap: 8px; margin-bottom: 4px; }
.method-tag { width: 60px; text-align: center; }
.endpoint-url { font-family: monospace; font-size: 13px; color: #333; }
.endpoint-name { font-size: 13px; color: #666; margin-bottom: 8px; }
.doc-block { margin-top: 4px; }
.block-label { font-size: 11px; color: #999; text-transform: uppercase; letter-spacing: 0.5px; }
pre { background: #1e1e1e; color: #d4d4d4; padding: 8px; border-radius: 4px; font-size: 12px; overflow-x: auto; margin-top: 2px; }
.empty { flex: 1; display: flex; align-items: center; justify-content: center; }
</style>
