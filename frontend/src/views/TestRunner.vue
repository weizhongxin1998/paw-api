<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { NButton, NTag, NIcon, NCheckbox, NSpace, NSpin, NEmpty } from 'naive-ui'
import { Play } from '@vicons/ionicons5'
import { useProjectStore } from '../stores/project'
import { ListCollections } from '../../wailsjs/go/handlers/CollectionHandler'
import { ListRequests } from '../../wailsjs/go/handlers/RequestHandler'
import { SendRequest } from '../../wailsjs/go/handlers/RequestHandler'

const { t } = useI18n()
const projectStore = useProjectStore()
const requests = ref<any[]>([])
const selected = ref<Set<string>>(new Set())
const results = ref<any[]>([])
const running = ref(false)

async function loadRequests() {
  if (!projectStore.currentProject) return
  const cols = await ListCollections(projectStore.currentProject.id)
  projectStore.setCollections(cols)
  const all: any[] = []
  for (const col of cols) {
    const reqs = await ListRequests(col.id)
    all.push(...reqs.map((r: any) => ({ ...r, collectionName: col.name })))
  }
  requests.value = all
}

async function runSelected() {
  running.value = true; results.value = []
  for (const req of requests.value) {
    if (!selected.value.has(req.id)) continue
    try {
      const headers = parseJSON(req.headers, {})
      const resp = await SendRequest({ Method: req.method, URL: req.url, Headers: typeof headers === 'object' && !Array.isArray(headers) ? headers : {}, Body: '' })
      results.value.push({ name: req.name, method: req.method, url: req.url, status: resp.status, duration: resp.duration_ms, passed: resp.status >= 200 && resp.status < 400 })
    } catch (e: any) { results.value.push({ name: req.name, method: req.method, url: req.url, status: 0, duration: 0, passed: false, error: e.message || t('testRunner.err') }) }
  }
  running.value = false
}

function parseJSON(str: string, fallback: any): any { try { return JSON.parse(str) } catch { return fallback } }
function toggleSelect(id: string) { if (selected.value.has(id)) selected.value.delete(id); else selected.value.add(id) }
function toggleAll() { if (selected.value.size === requests.value.length) selected.value.clear(); else selected.value = new Set(requests.value.map(r => r.id)) }

onMounted(loadRequests)
</script>

<template>
  <div class="test-runner">
    <div class="toolbar">
      <h2 class="title">{{ t('testRunner.title') }}</h2>
      <NSpace>
        <NButton size="small" @click="toggleAll">{{ selected.size === requests.length ? t('testRunner.deselectAll') : t('testRunner.selectAll') }}</NButton>
        <NButton type="primary" size="small" :loading="running" :disabled="selected.size === 0" @click="runSelected">
          <template #icon><NIcon><Play /></NIcon></template>
          {{ t('testRunner.run') }} {{ selected.size }} {{ t('testRunner.tests') }}
        </NButton>
      </NSpace>
    </div>
    <NSpin :show="running">
      <div v-if="requests.length === 0" class="empty"><NEmpty :description="t('testRunner.empty')" /></div>
      <div v-else class="content">
        <div class="request-list">
          <div v-for="req in requests" :key="req.id" class="req-row" @click="toggleSelect(req.id)">
            <NCheckbox :checked="selected.has(req.id)" />
            <NTag size="tiny" class="method-tag">{{ req.method }}</NTag>
            <span class="req-name">{{ req.name }}</span>
            <span class="req-url">{{ req.url }}</span>
          </div>
        </div>
        <div v-if="results.length > 0" class="results-section">
          <h3>{{ t('testRunner.results') }}</h3>
          <div v-for="(r, i) in results" :key="i" class="result-row" :class="{ passed: r.passed, failed: !r.passed }">
            <NTag size="tiny" :type="r.passed ? 'success' : 'error'">{{ r.passed ? t('testRunner.pass') : t('testRunner.fail') }}</NTag>
            <span class="result-name">{{ r.name }}</span>
            <span class="result-status">{{ r.status || t('testRunner.err') }}</span>
            <span class="result-duration">{{ r.duration }}ms</span>
            <span v-if="r.error" class="result-error">{{ r.error }}</span>
          </div>
        </div>
      </div>
    </NSpin>
  </div>
</template>

<style scoped>
.test-runner { padding: 16px 20px; height: 100%; display: flex; flex-direction: column; }
.toolbar { display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px; }
.title { font-size: 18px; font-weight: 600; }
.content { flex: 1; overflow-y: auto; }
.request-list { margin-bottom: 16px; }
.req-row { display: flex; align-items: center; gap: 8px; padding: 6px 8px; border-bottom: 1px solid var(--border-color); cursor: pointer; font-size: 13px; }
.req-row:hover { background: var(--hover-color); }
.method-tag { width: 50px; text-align: center; }
.req-name { font-weight: 500; min-width: 150px; }
.req-url { color: #666; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.results-section { border-top: 2px solid var(--border-color); padding-top: 12px; }
.result-row { display: flex; align-items: center; gap: 8px; padding: 6px 8px; font-size: 13px; }
.result-row.passed { background: #f0faf5; }
.result-row.failed { background: #fef0f0; }
.result-name { flex: 1; font-weight: 500; }
.result-status { width: 50px; text-align: right; }
.result-duration { width: 60px; text-align: right; color: #999; }
.result-error { color: #e74c3c; font-size: 12px; }
.empty { flex: 1; display: flex; align-items: center; justify-content: center; }
</style>
