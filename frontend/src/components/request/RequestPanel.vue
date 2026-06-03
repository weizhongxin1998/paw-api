<template>
  <div class="request-panel">
    <div class="sub-tabs">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        :class="{ active: activeTab === tab.key }"
        @click="activeTab = tab.key"
      >
        {{ tab.label }}
        <span v-if="tab.count != null" class="cnt">{{ tab.count }}</span>
      </button>
      <div style="flex:1"></div>
      <button v-if="activeTab === 'headers'" class="bulk-btn" @click="isBulkEdit = !isBulkEdit">
        {{ isBulkEdit ? 'Table' : 'Bulk' }}
      </button>
    </div>

    <div class="sub-content">
      <div v-if="activeTab === 'params'" class="params-content">
        <div class="params-section">
          <div class="params-section-hdr">
            <label class="section-toggle">
              <input type="checkbox" v-model="queryParamsEnabled" />
              <span>Query Params</span>
            </label>
          </div>
          <KeyValueTable
            v-if="queryParamsEnabled"
            :items="paramsItems"
            @update:items="onParamsChange"
          />
        </div>
      </div>

      <div v-else-if="activeTab === 'path'" class="params-content">
        <div v-if="pathVariables.length > 0" class="params-section">
          <div class="params-section-hdr">
            <span class="section-label">Path Variables</span>
            <span class="section-hint">URL 中 :param 自动识别</span>
          </div>
          <table class="kvt">
            <thead>
              <tr>
                <th>Key</th>
                <th>Value</th>
                <th>Description</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="pv in pathVariables" :key="pv.key">
                <td><input class="kvt-input readonly" :value="pv.key" readonly /></td>
                <td><input class="kvt-input" v-model="pv.value" placeholder="值" @input="onPathVarChange" /></td>
                <td><input class="kvt-input" v-model="pv.description" placeholder="描述" @input="onPathVarChange" /></td>
              </tr>
            </tbody>
          </table>
        </div>
        <div v-else class="path-empty">
          <span class="hint-text">URL 中未检测到路径参数，使用 <code>:param</code> 或 <code>{param}</code> 语法</span>
        </div>
      </div>

      <KeyValueTable
        v-else-if="activeTab === 'headers'"
        :items="headersItems"
        :show-bulk-edit="false"
        :bulk-mode="isBulkEdit"
        @update:items="onHeadersChange"
      />
      <BodyEditor
        v-else-if="activeTab === 'body'"
        :body-type="bodyType"
        :body-data="bodyData"
        @update:body-type="v => { bodyType = v; sync() }"
        @update:body-data="v => { bodyData = v; sync() }"
      />
      <AuthEditor
        v-else-if="activeTab === 'auth'"
        :model-value="authData"
        @update:model-value="v => { authData = v; sync() }"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import KeyValueTable from '../shared/KeyValueTable.vue'
import BodyEditor from './BodyEditor.vue'
import AuthEditor from './AuthEditor.vue'
import type { KvItem } from '../../types/request'

interface PathVariable { key: string; value: string; description: string }

const props = defineProps<{
  headers: string; params: string; bodyType: string; bodyData: string
  authData: string; url: string; pathVars: string
}>()

const emit = defineEmits<{
  (e: 'update:headers', v: string): void
  (e: 'update:params', v: string): void
  (e: 'update:bodyType', v: string): void
  (e: 'update:bodyData', v: string): void
  (e: 'update:authData', v: string): void
  (e: 'update:pathVars', v: string): void
}>()

const tabs = [
  { key: 'params', label: 'Params', count: undefined as number | undefined },
  { key: 'path', label: 'Path', count: undefined as number | undefined },
  { key: 'headers', label: 'Headers', count: undefined as number | undefined },
  { key: 'body', label: 'Body' },
  { key: 'auth', label: 'Auth' },
]

const activeTab = ref('params')
const isBulkEdit = ref(false)
const queryParamsEnabled = ref(true)
const bodyType = ref(props.bodyType)
const bodyData = ref(props.bodyData)
const authData = ref(props.authData)
const paramsItems = ref<KvItem[]>([])
const headersItems = ref<KvItem[]>([])
const pathVariables = ref<PathVariable[]>([])

let iid = Date.now()

const pathVarPatterns = computed(() => {
  const url = props.url || ''
  const patterns: string[] = []
  const re = /:([a-zA-Z_]\w*)/g
  let match: RegExpExecArray | null
  while ((match = re.exec(url)) !== null) {
    if (!patterns.includes(match[1])) patterns.push(match[1])
  }
  const re2 = /\{([a-zA-Z_]\w*)\}/g
  while ((match = re2.exec(url)) !== null) {
    if (!patterns.includes(match[1])) patterns.push(match[1])
  }
  return patterns
})

watch(pathVarPatterns, (patterns) => {
  let stored: PathVariable[] = []
  try { stored = JSON.parse(props.pathVars || '[]') } catch {}
  const storedMap = new Map(stored.map(p => [p.key, p]))
  const existing = new Map(pathVariables.value.map(p => [p.key, p]))
  pathVariables.value = patterns.map(key => {
    const s = storedMap.get(key); const e = existing.get(key)
    return { key, value: s?.value || e?.value || '', description: s?.description || e?.description || '' }
  })
  tabs[1].count = patterns.length > 0 ? patterns.length : undefined
  syncPathVars()
}, { immediate: true })

function onPathVarChange() { syncPathVars() }
function syncPathVars() { emit('update:pathVars', JSON.stringify(pathVariables.value)) }

function parseKv(raw: string): KvItem[] {
  try {
    const items = JSON.parse(raw).map((i: any) => ({ ...i, id: i.id || String(++iid) }))
    // Always ensure a trailing empty row so users can add new entries
    // (KeyValueTable skips its auto-append on initial mount to avoid markDirty)
    const last = items[items.length - 1]
    if (!last || last.key || last.value || last.description) {
      items.push({ id: String(++iid), key: '', value: '', description: '', enabled: true })
    }
    return items
  } catch {
    return [{ id: String(++iid), key: '', value: '', description: '', enabled: true }]
  }
}

watch(() => props.params, (v) => { paramsItems.value = parseKv(v) }, { immediate: true })
watch(() => props.headers, (v) => { headersItems.value = parseKv(v) }, { immediate: true })
watch(() => props.bodyType, (v) => { bodyType.value = v })
watch(() => props.bodyData, (v) => { bodyData.value = v })
watch(() => props.authData, (v) => { authData.value = v })

function sync() {
  emit('update:headers', JSON.stringify(headersItems.value))
  emit('update:params', JSON.stringify(paramsItems.value))
  emit('update:bodyType', bodyType.value)
  emit('update:bodyData', bodyData.value)
  emit('update:authData', authData.value)
}

function onHeadersChange(items: KvItem[]) {
  headersItems.value = items
  tabs[2].count = items.filter(i => i.enabled && i.key).length
  sync()
}

function onParamsChange(items: KvItem[]) {
  paramsItems.value = items
  tabs[0].count = items.filter(i => i.enabled && i.key).length
  sync()
}
</script>

<style scoped>
.request-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: var(--bg-base);
  min-height: 0;
}

/* ── Sub Tabs ── */
.sub-tabs {
  display: flex;
  border-bottom: 1px solid var(--border-primary);
  background: var(--bg-surface);
  padding: 0 10px;
  gap: 0;
}
.sub-tabs button {
  padding: 8px 14px;
  font-size: var(--fs-sm);
  cursor: pointer;
  color: var(--text-muted);
  border: none;
  background: transparent;
  border-bottom: 2px solid transparent;
  outline: none;
  font-family: var(--font-mono);
  font-weight: 500;
  letter-spacing: 0.02em;
  transition: all var(--transition);
  position: relative;
}
.sub-tabs button.active {
  color: var(--accent);
  border-bottom-color: var(--accent);
  font-weight: 600;
}
.sub-tabs button.active::after {
  content: '';
  position: absolute;
  bottom: -1px; left: 20%; right: 20%;
  height: 1px;
  background: var(--accent);
  filter: blur(4px);
  opacity: 0.5;
}
.sub-tabs button:hover:not(.active) {
  color: var(--text-secondary);
  background: var(--bg-hover);
  border-radius: var(--radius-sm) var(--radius-sm) 0 0;
}

.cnt {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: var(--fs-2xs);
  background: var(--bg-active);
  color: var(--text-secondary);
  padding: 0 5px;
  border-radius: 10px;
  margin-left: 4px;
  font-weight: 600;
  min-width: 16px;
  height: 16px;
  line-height: 1;
}

.bulk-btn {
  font-size: var(--fs-xs) !important;
  color: var(--accent) !important;
  padding: 8px 10px !important;
  border-radius: var(--radius-sm) !important;
  transition: all var(--transition) !important;
}
.bulk-btn:hover {
  background: var(--accent-soft) !important;
}

/* ── Content ── */
.sub-content {
  flex: 1;
  overflow-y: auto;
}
.params-content { padding: 10px; }
.params-section { margin-bottom: 12px; }
.params-section-hdr {
  display: flex; align-items: center; gap: 8px; margin-bottom: 4px;
}
.section-toggle {
  display: flex; align-items: center; gap: 6px;
  font-size: var(--fs-sm); color: var(--text-secondary);
  font-weight: 600; cursor: pointer; font-family: var(--font-mono);
}
.section-toggle input[type="checkbox"] {
  accent-color: var(--accent);
  width: 14px; height: 14px;
}
.section-label {
  font-size: var(--fs-sm); color: var(--text-secondary);
  font-weight: 600; font-family: var(--font-mono);
}
.section-hint {
  font-size: var(--fs-xs); color: var(--text-muted);
  font-family: var(--font-ui);
}

/* ── Path Variables Table ── */
.kvt {
  width: 100%;
  border-collapse: separate;
  border-spacing: 0;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  overflow: hidden;
}
.kvt th {
  text-align: left; padding: 6px 10px; font-size: var(--fs-xs); color: var(--text-muted);
  text-transform: uppercase; border-bottom: 1px solid var(--border-primary);
  font-weight: 600; background: var(--bg-elevated); letter-spacing: 0.04em;
  font-family: var(--font-ui);
}
.kvt td { padding: 3px 6px; }
.kvt tr:not(:last-child) td { border-bottom: 1px solid var(--border-subtle); }
.kvt-input {
  width: 100%; padding: 6px 8px; border: 1px solid transparent;
  font-family: var(--font-mono); font-size: var(--fs-sm); background: transparent;
  border-radius: var(--radius-xs); outline: none; color: var(--text-primary);
  transition: border-color var(--transition), background var(--transition);
}
.kvt-input:hover:not(.readonly) {
  border-color: var(--border-hover); background: var(--bg-surface);
}
.kvt-input:focus {
  border-color: var(--accent); background: var(--bg-surface);
}
.kvt-input.readonly {
  color: var(--text-muted); background: transparent; cursor: default;
}
.kvt-input.readonly:hover { border-color: transparent; background: transparent; }

.path-empty {
  padding: 28px; text-align: center;
}
.hint-text {
  font-size: var(--fs-sm); color: var(--text-muted); font-family: var(--font-mono);
}
.hint-text code {
  background: var(--bg-elevated); padding: 2px 6px;
  border-radius: var(--radius-xs); font-family: var(--font-mono); font-size: var(--fs-xs);
  border: 1px solid var(--border-primary);
  color: var(--accent);
}
</style>
