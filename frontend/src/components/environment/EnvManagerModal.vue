<template>
  <n-modal
    :show="show"
    preset="card"
    title="环境管理"
    :class="modalClass"
    style="width: 760px"
    :mask-closable="false"
    @update:show="$emit('update:show', $event)"
  >
    <div class="manager-body">
      <!-- Environment list panel with card-style items -->
      <div class="env-list-panel">
        <div class="panel-header">
          <span>环境列表</span>
          <n-button size="tiny" @click="addEnv">
            <template #icon>
              <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
            </template>
            新建
          </n-button>
        </div>
        <div class="env-list">
          <div
            v-for="env in environments"
            :key="env.id"
            class="env-card"
            :class="{ active: editingEnvId === env.id }"
            @click="selectEnv(env)"
          >
            <!-- Active indicator bar -->
            <div class="env-card-active-bar" v-if="env.is_active"></div>

            <div class="env-card-body">
              <div class="env-card-top">
                <!-- Rename inline -->
                <n-input
                  v-if="renamingId === env.id"
                  v-model:value="renameText"
                  size="small"
                  @blur="confirmRename(env)"
                  @keydown.enter="confirmRename(env)"
                  class="rename-input"
                />
                <template v-else>
                  <span class="env-name">{{ env.name }}</span>
                  <!-- Environment type badge -->
                  <span class="env-type-badge" :class="detectEnvType(env)">{{ envTypeLabel(env) }}</span>
                </template>
              </div>
              <div class="env-card-url" v-if="env.base_url" :title="env.base_url">
                {{ env.base_url }}
              </div>
              <!-- Active label -->
              <span class="env-active-label" v-if="env.is_active">当前活跃</span>
            </div>

            <!-- Actions dropdown -->
            <n-dropdown trigger="click" :options="envMenuOptions(env)" @select="(k: string) => onEnvMenu(k, env)">
              <n-button text size="tiny" class="env-card-menu" @click.stop>&#8943;</n-button>
            </n-dropdown>
          </div>

          <!-- Empty state -->
          <div v-if="environments.length === 0" class="env-list-empty">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.2" opacity="0.3">
              <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0118 0z"/>
              <circle cx="12" cy="10" r="3"/>
            </svg>
            <span>暂无环境</span>
          </div>
        </div>
      </div>

      <!-- Variable editing panel -->
      <div class="var-list-panel">
        <template v-if="editingEnvId">
          <div class="panel-header">
            <span>变量配置</span>
            <div class="panel-header-actions">
              <!-- Test connection button -->
              <n-button size="tiny" :loading="testing" @click="onTestConnection" :disabled="!baseURL" secondary>
                <template #icon>
                  <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M22 11.08V12a10 10 0 11-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
                </template>
                测试连接
              </n-button>
              <n-button size="tiny" @click="addVariable">
                <template #icon>
                  <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
                </template>
                添加
              </n-button>
            </div>
          </div>

          <!-- Base URL section with validation -->
          <div class="base-url-section">
            <label>前置 URL</label>
            <n-input
              v-model:value="baseURL"
              size="small"
              placeholder="https://api.example.com"
              :status="baseUrlError ? 'error' : undefined"
            >
              <template #prefix>
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><circle cx="12" cy="12" r="10"/><path d="M2 12h20"/><path d="M12 2a15.3 15.3 0 014 10 15.3 15.3 0 01-4 10 15.3 15.3 0 01-4-10 15.3 15.3 0 014-10z"/></svg>
              </template>
            </n-input>
            <span class="field-error" v-if="baseUrlError">{{ baseUrlError }}</span>
            <!-- Connection test result -->
            <div class="test-result" v-if="testResult" :class="testResult.success ? 'success' : 'error'">
              <svg v-if="testResult.success" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="20 6 9 17 4 12"/></svg>
              <svg v-else width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>
              <span>{{ testResult.message }}</span>
            </div>
          </div>

          <!-- Variable table with better layout -->
          <div class="var-table">
            <div class="var-row var-header">
              <span class="col-check">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 11 12 14 22 4"/><path d="M21 12v7a2 2 0 01-2 2H5a2 2 0 01-2-2V5a2 2 0 012-2h11"/></svg>
              </span>
              <span class="col-key">键</span>
              <span class="col-value">值</span>
              <span class="col-action"></span>
            </div>
            <div v-for="(variable, idx) in variables" :key="idx" class="var-row" :class="{ 'var-disabled': !variable.enabled }">
              <span class="col-check">
                <n-checkbox v-model:checked="variable.enabled" size="small" />
              </span>
              <span class="col-key">
                <n-input
                  v-model:value="variable.key"
                  size="small"
                  placeholder="变量名"
                  :status="validateVarKey(variable.key, idx) ? undefined : 'warning'"
                />
              </span>
              <span class="col-value">
                <n-input
                  v-model:value="variable.value"
                  size="small"
                  placeholder="变量值"
                  :type="isSecretValue(variable.key) ? 'password' : 'text'"
                  show-password-on="click"
                />
              </span>
              <span class="col-action">
                <n-button text size="tiny" @click="removeVariable(idx)" class="var-remove-btn">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                </n-button>
              </span>
            </div>
            <!-- Empty variables hint -->
            <div v-if="variables.length === 0" class="var-empty">
              暂无变量，点击上方"添加"按钮
            </div>
          </div>

          <!-- Validation warnings -->
          <div class="var-warnings" v-if="varWarnings.length > 0">
            <div v-for="(w, i) in varWarnings" :key="i" class="var-warning">
              <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
              {{ w }}
            </div>
          </div>

          <div class="panel-footer">
            <n-button type="primary" size="small" :loading="saving" @click="saveVariables">
              <template #icon>
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M19 21H5a2 2 0 01-2-2V5a2 2 0 012-2h11l5 5v11a2 2 0 01-2 2z"/><polyline points="17 21 17 13 7 13 7 21"/><polyline points="7 3 7 8 15 8"/></svg>
              </template>
              保存
            </n-button>
          </div>
        </template>
        <div v-else class="panel-empty">
          <svg width="36" height="36" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" opacity="0.2">
            <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0118 0z"/>
            <circle cx="12" cy="10" r="3"/>
          </svg>
          <span>选择一个环境以编辑变量</span>
        </div>
      </div>
    </div>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { NModal, NButton, NInput, NCheckbox, NDropdown, useMessage } from 'naive-ui'
import {
  ListEnvironments, CreateEnvironment, RenameEnvironment, DeleteEnvironment,
  ListEnvVariables, SaveEnvVariables, SaveEnvBaseURL, SendQuickRequest,
} from '../../../wailsjs/go/main/App'
import { useEnvStore } from '../../stores/env'
import type { Environment, EnvVariable } from '../../types/environment'

interface Props { show: boolean; projectId: number | null }
const props = defineProps<Props>()
const emit = defineEmits<{ 'update:show': [value: boolean]; refresh: [] }>()

// Detect light mode so teleported modal gets correct CSS variables
const isLightMode = ref(false)
onMounted(() => {
  const check = () => { isLightMode.value = !!document.querySelector('.theme-light') }
  check()
  const observer = new MutationObserver(check)
  observer.observe(document.documentElement, { attributes: true, attributeFilter: ['class'], subtree: true })
})
const modalClass = computed(() => isLightMode.value ? 'env-manager-modal theme-light' : 'env-manager-modal')
const environments = ref<Environment[]>([])
const editingEnvId = ref<number | null>(null)
const renamingId = ref<number | null>(null)
const renameText = ref('')
const variables = ref<EnvVariable[]>([])
const saving = ref(false)
const testing = ref(false)
const baseURL = ref('')
const baseUrlError = ref('')
const testResult = ref<{ success: boolean; message: string } | null>(null)
const message = useMessage()
const envStore = useEnvStore()

watch(() => props.show, async (v) => { if (v) await loadEnvs() })

// -- Environment type detection --
function detectEnvType(env: Environment): string {
  const name = (env.name || '').toLowerCase()
  const url = (env.base_url || '').toLowerCase()
  const combined = name + ' ' + url
  if (combined.includes('dev') || combined.includes('local') || combined.includes('localhost')) return 'dev'
  if (combined.includes('staging') || combined.includes('stg') || combined.includes('test') || combined.includes('uat')) return 'staging'
  if (combined.includes('prod') || combined.includes('production') || combined.includes('live')) return 'prod'
  return 'default'
}

function envTypeLabel(env: Environment): string {
  const type = detectEnvType(env)
  switch (type) {
    case 'dev': return 'DEV'
    case 'staging': return 'STG'
    case 'prod': return 'PRD'
    default: return 'ENV'
  }
}

// -- Variable validation --
function validateVarKey(key: string, idx: number): boolean {
  if (!key || key.trim() === '') return true // empty is allowed (just not filled yet)
  // Check for duplicate keys
  const dupes = variables.value.filter((v, i) => i !== idx && v.key === key && v.enabled)
  return dupes.length === 0
}

function isSecretValue(key: string): boolean {
  if (!key) return false
  const lk = key.toLowerCase()
  return lk.includes('password') || lk.includes('secret') || lk.includes('token') || lk.includes('api_key') || lk.includes('apikey')
}

const varWarnings = computed(() => {
  const warnings: string[] = []
  const keys = variables.value.filter(v => v.enabled).map(v => v.key)
  const seen = new Set<string>()
  for (const key of keys) {
    if (key && seen.has(key)) {
      warnings.push(`变量 "${key}" 存在重复键名`)
    }
    if (key) seen.add(key)
  }
  // Check for empty enabled keys
  const emptyKeys = variables.value.filter(v => v.enabled && (!v.key || v.key.trim() === ''))
  if (emptyKeys.length > 0) {
    warnings.push(`${emptyKeys.length} 个已启用的变量缺少键名`)
  }
  return warnings
})

// -- Data loading --
async function loadEnvs() {
  if (!props.projectId) return
  try { environments.value = await ListEnvironments(props.projectId) || [] } catch { environments.value = [] }
}

async function selectEnv(env: Environment) {
  editingEnvId.value = env.id
  baseURL.value = env.base_url || ''
  baseUrlError.value = ''
  testResult.value = null
  try {
    const vars = await ListEnvVariables(env.id)
    variables.value = vars.map(v => ({ id: v.id, key: v.key, value: v.value, enabled: v.enabled }))
  } catch { variables.value = [] }
}

async function addEnv() {
  if (!props.projectId) return
  const name = `环境 ${environments.value.length + 1}`
  try {
    await CreateEnvironment(props.projectId, name, '', null)
    await loadEnvs()
    emit('refresh')
  } catch {}
}

function envMenuOptions(env: Environment) {
  return [
    { label: '重命名', key: 'rename' },
    { label: '复制环境', key: 'duplicate' },
    { label: '复制', key: 'copy' },
    { type: 'divider', key: 'd1' },
    { label: '删除', key: 'delete', props: { style: 'color: var(--red, #ef4444)' } },
  ]
}

async function onEnvMenu(key: string, env: Environment) {
  switch (key) {
    case 'rename':
      renamingId.value = env.id
      renameText.value = env.name
      break
    case 'duplicate':
    case 'copy':
      if (!props.projectId) return
      try {
        // Duplicate: create new env with same name + suffix, then copy variables
        const newName = env.name + ' (副本)'
        const newEnv = await CreateEnvironment(props.projectId, newName, env.base_url || '', env.id as any)
        await loadEnvs()
        emit('refresh')
        message.success('环境已复制')
      } catch {}
      break
    case 'delete':
      try {
        await DeleteEnvironment(env.id)
        if (editingEnvId.value === env.id) {
          editingEnvId.value = null
          variables.value = []
        }
        await loadEnvs()
        emit('refresh')
        message.success('环境已删除')
      } catch {}
      break
  }
}

async function confirmRename(env: Environment) {
  if (renameText.value.trim() && renameText.value !== env.name) {
    try {
      await RenameEnvironment(env.id, renameText.value.trim())
      await loadEnvs()
      emit('refresh')
    } catch {}
  }
  renamingId.value = null
}

function addVariable() {
  variables.value.push({ id: 0, key: '', value: '', enabled: true })
}

function removeVariable(idx: number) {
  variables.value.splice(idx, 1)
}

// -- Save variables --
async function saveVariables() {
  if (!editingEnvId.value) return

  // Validate base URL format
  if (baseURL.value && !baseURL.value.match(/^https?:\/\//)) {
    baseUrlError.value = 'URL 应以 http:// 或 https:// 开头'
    return
  }
  baseUrlError.value = ''

  saving.value = true
  try {
    await SaveEnvVariables(
      editingEnvId.value,
      variables.value.map(v => ({
        id: v.id,
        environment_id: editingEnvId.value!,
        key: v.key,
        value: v.value,
        enabled: v.enabled,
        sort_order: 0,
        created_at: '',
      } as any))
    )
    await SaveEnvBaseURL(editingEnvId.value, baseURL.value)
    if (props.projectId) await envStore.loadEnvironments(props.projectId)
    message.success('环境变量已保存')
    // Reload current env variables
    await selectEnv({
      id: editingEnvId.value,
      project_id: 0,
      name: '',
      base_url: baseURL.value,
      is_active: false,
      created_at: '',
      updated_at: '',
    })
  } catch (e: any) {
    message.error('保存失败: ' + (e?.message || String(e)))
  }
  saving.value = false
}

// -- Test connection --
async function onTestConnection() {
  if (!baseURL.value || !editingEnvId.value) return
  testing.value = true
  testResult.value = null
  try {
    const resp = await SendQuickRequest(editingEnvId.value, 'GET', baseURL.value, '', '', 0)
    if (resp && resp.status >= 200 && resp.status < 400) {
      testResult.value = { success: true, message: `连接成功 (${resp.status} ${resp.status_text || ''}, ${resp.time}ms)` }
    } else if (resp) {
      testResult.value = { success: false, message: `连接返回 ${resp.status} ${resp.status_text || ''}` }
    } else {
      testResult.value = { success: false, message: '无响应' }
    }
  } catch (e: any) {
    testResult.value = { success: false, message: '连接失败: ' + (e?.message || String(e)) }
  }
  testing.value = false
}
</script>

<style scoped>
.manager-body { display: flex; gap: 16px; height: 460px; }

/* -- Environment List Panel -- */
.env-list-panel {
  width: 220px;
  border-right: 1px solid var(--border-primary);
  padding-right: 12px;
  display: flex;
  flex-direction: column;
}
.var-list-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  font-weight: 600;
  font-size: var(--fs-base);
  color: var(--text-primary);
}
.panel-header-actions {
  display: flex;
  gap: 6px;
}

/* -- Environment Cards -- */
.env-list { flex: 1; overflow-y: auto; }
.env-card {
  display: flex;
  align-items: center;
  padding: 8px 10px;
  margin-bottom: 4px;
  border-radius: var(--radius, 8px);
  cursor: pointer;
  font-size: var(--fs-sm);
  color: var(--text-secondary);
  transition: all 0.15s ease;
  border: 1px solid transparent;
  position: relative;
  overflow: hidden;
}
.env-card:hover { background: var(--bg-hover); }
.env-card.active {
  background: var(--accent-soft);
  border-color: var(--accent);
  color: var(--text-primary);
}
.env-card-active-bar {
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 3px;
  background: var(--accent);
  border-radius: 0 2px 2px 0;
}
.env-card-body {
  flex: 1;
  min-width: 0;
}
.env-card-top {
  display: flex;
  align-items: center;
  gap: 6px;
}
.env-name {
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.rename-input { width: 100%; }

/* Environment type badges */
.env-type-badge {
  font-size: var(--fs-2xs);
  font-weight: 700;
  letter-spacing: 0.5px;
  padding: 1px 5px;
  border-radius: 4px;
  flex-shrink: 0;
  line-height: 1.4;
}
.env-type-badge.dev     { background: rgba(59,130,246,0.15); color: var(--method-get); }
.env-type-badge.staging { background: rgba(245,158,11,0.15); color: var(--method-put); }
.env-type-badge.prod    { background: rgba(34,197,94,0.15); color: var(--method-post); }
.env-type-badge.default { background: rgba(113,113,122,0.12); color: var(--text-secondary); }

.env-card-url {
  font-size: var(--fs-2xs, 9px);
  color: var(--text-secondary);
  font-family: var(--font-mono);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-top: 2px;
}
.env-active-label {
  font-size: var(--fs-2xs);
  color: var(--accent-text);
  font-weight: 600;
  letter-spacing: 0.3px;
  margin-top: 2px;
  display: inline-block;
}
.env-card-menu {
  margin-left: auto;
  flex-shrink: 0;
  opacity: 0.4;
  transition: opacity 0.15s;
}
.env-card:hover .env-card-menu { opacity: 1; }

/* Empty list */
.env-list-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 30px;
  color: var(--text-secondary);
  font-size: var(--fs-sm);
}

/* -- Base URL Section -- */
.base-url-section {
  margin-bottom: 10px;
}
.base-url-section label {
  display: block;
  font-size: var(--fs-xs);
  color: var(--text-secondary);
  margin-bottom: 4px;
  text-transform: uppercase;
  letter-spacing: 0.3px;
  font-weight: 600;
}
.field-error {
  font-size: var(--fs-xs);
  color: var(--red, #ef4444);
  margin-top: 3px;
  display: block;
}
.test-result {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 6px;
  padding: 5px 10px;
  border-radius: var(--radius-sm);
  font-size: var(--fs-xs);
}
.test-result.success {
  background: rgba(34,197,94,0.08);
  color: var(--method-post);
}
.test-result.error {
  background: rgba(239,68,68,0.08);
  color: var(--red, #ef4444);
}

/* -- Variable Table -- */
.var-table { flex: 1; overflow-y: auto; }
.var-row {
  display: flex;
  align-items: center;
  padding: 4px 0;
  gap: 6px;
}
.var-row.var-header {
  font-size: var(--fs-xs);
  font-weight: 600;
  color: var(--text-secondary);
  padding-bottom: 6px;
  border-bottom: 1px solid var(--border-primary);
  margin-bottom: 4px;
}
.var-row.var-disabled { opacity: 0.5; }
.col-check { width: 24px; text-align: center; flex-shrink: 0; display: flex; align-items: center; justify-content: center; }
.col-key { width: 140px; flex-shrink: 0; }
.col-value { flex: 1; }
.col-action { width: 28px; text-align: center; flex-shrink: 0; }
.var-remove-btn {
  color: var(--text-muted) !important;
  transition: color 0.15s !important;
}
.var-remove-btn:hover { color: var(--red, #ef4444) !important; }

.var-empty {
  padding: 24px;
  text-align: center;
  color: var(--text-secondary);
  font-size: var(--fs-sm);
}

/* -- Validation Warnings -- */
.var-warnings {
  margin-top: 6px;
}
.var-warning {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: var(--fs-xs);
  color: var(--amber, #f59e0b);
  padding: 2px 0;
}
.var-warning svg { flex-shrink: 0; opacity: 0.7; }

/* -- Footer -- */
.panel-footer {
  margin-top: 10px;
  display: flex;
  justify-content: flex-end;
}

/* -- Panel empty state -- */
.panel-empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: var(--text-secondary);
  font-size: var(--fs-sm);
}
</style>
