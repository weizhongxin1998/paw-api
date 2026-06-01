<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { t } from '../i18n'
import { NButton, NTag, NIcon, NCheckbox, NSpace, NSpin, NEmpty, NSelect, NInput, NPopover } from 'naive-ui'
import { Play, Add, Trash, Settings as SettingsIcon } from '@vicons/ionicons5'
import { useProjectStore } from '../stores/project'
import { ListCollections } from '../../wailsjs/go/handlers/CollectionHandler'
import { ListRequests } from '../../wailsjs/go/handlers/RequestHandler'
import { RunAsserts } from '../../wailsjs/go/handlers/RequestHandler'

const projectStore = useProjectStore()
const requests = ref<any[]>([])
const selected = ref<Set<string>>(new Set())
const results = ref<any[]>([])
const asserts = ref<Record<string, any[]>>({})
const running = ref(false)

const assertTypes = [
  { label: 'Status Code', value: 'status' },
  { label: 'Body Contains', value: 'body_contains' },
  { label: 'JSONPath', value: 'body_jsonpath' },
  { label: 'Header Equals', value: 'header_equals' },
  { label: 'Duration <', value: 'duration_lt' },
]

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

function getAsserts(reqId: string): any[] {
  if (!asserts.value[reqId]) asserts.value[reqId] = []
  return asserts.value[reqId]
}

function addAssert(reqId: string) {
  if (!asserts.value[reqId]) asserts.value[reqId] = []
  asserts.value[reqId].push({ type: 'status', target: '', value: '200' })
}

function removeAssert(reqId: string, index: number) {
  asserts.value[reqId].splice(index, 1)
}

function updateAssert(reqId: string, index: number, field: string, val: any) {
  if (asserts.value[reqId]?.[index]) {
    asserts.value[reqId][index][field] = val
  }
}

async function runSelected() {
  running.value = true; results.value = []
  for (const req of requests.value) {
    if (!selected.value.has(req.id)) continue
    try {
      const reqAsserts = getAsserts(req.id)
      const headers = parseJSON(req.headers, {})
      const headerMap: Record<string, string> = {}
      if (typeof headers === 'object' && !Array.isArray(headers)) {
        Object.entries(headers).forEach(([k, v]: [string, any]) => { headerMap[k] = String(v) })
      }
      const result = await (RunAsserts as any)({
        Method: req.method,
        URL: req.url,
        Headers: headerMap,
        Body: '',
        BodyType: 'none',
        BodyFiles: [],
        AuthType: 'none',
        AuthData: {},
        TimeoutMs: 30000,
        FollowRedirect: true,
        Asserts: reqAsserts,
      })
      if (result && result.response) {
        const response = result.response
        const assertResults = result.asserts || []
        const allPassed = assertResults.length === 0 || assertResults.every((a: any) => a.passed)
        results.value.push({
          name: req.name, method: req.method, url: req.url,
          status: response.status, duration: response.duration_ms,
          passed: allPassed,
          asserts: assertResults,
        })
      }
    } catch (e: any) {
      results.value.push({ name: req.name, method: req.method, url: req.url, status: 0, duration: 0, passed: false, error: e.message || t('testRunner.err'), asserts: [] })
    }
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
      <h2 class="title">{{ $t('testRunner.title') }}</h2>
      <NSpace>
        <NButton size="small" @click="toggleAll">{{ selected.size === requests.length ? $t('testRunner.deselectAll') : $t('testRunner.selectAll') }}</NButton>
        <NButton type="primary" size="small" :loading="running" :disabled="selected.size === 0" @click="runSelected">
          <template #icon><NIcon><Play /></NIcon></template>
          {{ $t('testRunner.run') }} {{ selected.size }} {{ $t('testRunner.tests') }}
        </NButton>
      </NSpace>
    </div>
    <NSpin :show="running">
      <div v-if="requests.length === 0" class="empty"><NEmpty :description="$t('testRunner.empty')" /></div>
      <div v-else class="content">
        <div class="request-list">
          <div v-for="req in requests" :key="req.id" class="req-row">
            <div class="req-row-main" @click="toggleSelect(req.id)">
              <NCheckbox :checked="selected.has(req.id)" />
              <NTag size="tiny" class="method-tag">{{ req.method }}</NTag>
              <span class="req-name">{{ req.name }}</span>
              <span class="req-url">{{ req.url }}</span>
            </div>
            <NPopover trigger="click" placement="left">
              <template #trigger>
                <NButton quaternary circle size="tiny">
                  <template #icon><NIcon size="14"><SettingsIcon /></NIcon></template>
                </NButton>
              </template>
              <div class="assert-editor">
                <div class="assert-title">Assertions</div>
                <div v-for="(a, i) in getAsserts(req.id)" :key="i" class="assert-row">
                  <NSelect :options="assertTypes" :value="a.type" size="tiny" style="width:110px" @update:value="updateAssert(req.id, i, 'type', $event)" />
                  <NInput v-if="a.type === 'status'" :value="a.value" size="tiny" placeholder="200" style="width:60px" @update:value="updateAssert(req.id, i, 'value', $event)" />
                  <NInput v-if="a.type === 'body_contains'" :value="a.value" size="tiny" placeholder="text" style="width:80px" @update:value="updateAssert(req.id, i, 'value', $event)" />
                  <NInput v-if="a.type === 'body_jsonpath'" :value="a.target" size="tiny" placeholder="$.field" style="width:80px" @update:value="updateAssert(req.id, i, 'target', $event)" />
                  <NInput v-if="a.type === 'body_jsonpath'" :value="a.value" size="tiny" placeholder="value" style="width:60px" @update:value="updateAssert(req.id, i, 'value', $event)" />
                  <NInput v-if="a.type === 'header_equals'" :value="a.target" size="tiny" placeholder="Header" style="width:70px" @update:value="updateAssert(req.id, i, 'target', $event)" />
                  <NInput v-if="a.type === 'header_equals'" :value="a.value" size="tiny" placeholder="val" style="width:60px" @update:value="updateAssert(req.id, i, 'value', $event)" />
                  <NInput v-if="a.type === 'duration_lt'" :value="a.value" size="tiny" placeholder="1000" style="width:60px" @update:value="updateAssert(req.id, i, 'value', $event)" />
                  <NButton quaternary circle size="tiny" @click="removeAssert(req.id, i)">
                    <template #icon><NIcon size="12"><Trash /></NIcon></template>
                  </NButton>
                </div>
                <NButton size="tiny" quaternary @click="addAssert(req.id)">
                  <template #icon><NIcon size="12"><Add /></NIcon></template>
                  Add Assert
                </NButton>
              </div>
            </NPopover>
          </div>
        </div>
        <div v-if="results.length > 0" class="results-section">
          <h3>{{ $t('testRunner.results') }}</h3>
          <div v-for="(r, i) in results" :key="i" class="result-row" :class="{ passed: r.passed, failed: !r.passed }">
            <NTag size="tiny" :type="r.passed ? 'success' : 'error'">{{ r.passed ? $t('testRunner.pass') : $t('testRunner.fail') }}</NTag>
            <span class="result-name">{{ r.name }}</span>
            <span class="result-status">{{ r.status || $t('testRunner.err') }}</span>
            <span class="result-duration">{{ r.duration }}ms</span>
            <span v-if="r.error" class="result-error">{{ r.error }}</span>
            <div v-if="r.asserts && r.asserts.length > 0" class="assert-results">
              <div v-for="(a, ai) in r.asserts" :key="ai" class="assert-result-item" :class="{ 'assert-pass': a.passed, 'assert-fail': !a.passed }">
                <NTag size="tiny" :type="a.passed ? 'success' : 'error'" style="margin-right:4px">{{ a.passed ? 'PASS' : 'FAIL' }}</NTag>
                <span>{{ a.rule.type }}: {{ a.actual }}</span>
              </div>
            </div>
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
.req-row { display: flex; align-items: center; gap: 4px; padding: 4px 8px; border-bottom: 1px solid var(--border-color); font-size: 13px; }
.req-row:hover { background: var(--hover-color); }
.req-row-main { display: flex; align-items: center; gap: 8px; flex: 1; cursor: pointer; }
.method-tag { width: 50px; text-align: center; }
.req-name { font-weight: 500; min-width: 150px; }
.req-url { color: #666; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; flex: 1; }
.results-section { border-top: 2px solid var(--border-color); padding-top: 12px; }
.result-row { display: flex; flex-wrap: wrap; align-items: center; gap: 8px; padding: 6px 8px; font-size: 13px; }
.result-row.passed { background: #f0faf5; }
.result-row.failed { background: #fef0f0; }
.result-name { font-weight: 500; min-width: 120px; }
.result-status { width: 50px; text-align: right; }
.result-duration { width: 60px; text-align: right; color: #999; }
.result-error { color: #e74c3c; font-size: 12px; width: 100%; }
.empty { flex: 1; display: flex; align-items: center; justify-content: center; }
.assert-editor { padding: 8px; min-width: 340px; }
.assert-title { font-size: 12px; font-weight: 600; margin-bottom: 6px; }
.assert-row { display: flex; align-items: center; gap: 4px; margin-bottom: 4px; }
.assert-results { width: 100%; display: flex; flex-wrap: wrap; gap: 4px; margin-left: 40px; }
.assert-result-item { display: flex; align-items: center; gap: 4px; font-size: 11px; }
</style>
