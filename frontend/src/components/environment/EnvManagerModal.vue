<template>
  <n-modal
    :show="show"
    preset="card"
    title="环境管理"
    style="width: 720px"
    :mask-closable="false"
    @update:show="$emit('update:show', $event)"
  >
    <div class="manager-body">
      <div class="env-list-panel">
        <div class="panel-header">
          <span>环境列表</span>
          <n-button size="tiny" @click="addEnv">新建</n-button>
        </div>
        <div class="env-list">
          <div
            v-for="env in environments"
            :key="env.id"
            class="env-item"
            :class="{ active: editingEnvId === env.id }"
            @click="selectEnv(env)"
          >
            <n-input
              v-if="renamingId === env.id"
              v-model:value="renameText"
              size="small"
              @blur="confirmRename(env)"
              @keydown.enter="confirmRename(env)"
            />
            <span v-else>{{ env.name }}</span>
            <n-dropdown trigger="click" :options="envMenuOptions(env)" @select="(k: string) => onEnvMenu(k, env)">
              <n-button text size="tiny" style="margin-left:auto">&#8943;</n-button>
            </n-dropdown>
          </div>
        </div>
      </div>

      <div class="var-list-panel">
        <template v-if="editingEnvId">
          <div class="panel-header">
            <span>变量</span>
            <n-button size="tiny" @click="addVariable">添加</n-button>
          </div>
          <div class="base-url-section">
            <label>前置 URL</label>
            <n-input v-model:value="baseURL" size="small" placeholder="https://api.example.com" />
          </div>
          <div class="var-table">
            <div class="var-row var-header">
              <span class="col-check"></span>
              <span class="col-key">Key</span>
              <span class="col-value">Value</span>
              <span class="col-action"></span>
            </div>
            <div v-for="(variable, idx) in variables" :key="idx" class="var-row">
              <span class="col-check">
                <n-checkbox v-model:checked="variable.enabled" size="small" />
              </span>
              <span class="col-key">
                <n-input v-model:value="variable.key" size="small" placeholder="变量名" />
              </span>
              <span class="col-value">
                <n-input v-model:value="variable.value" size="small" placeholder="变量值" />
              </span>
              <span class="col-action">
                <n-button text size="tiny" @click="removeVariable(idx)" style="color:#ff4444">&#10005;</n-button>
              </span>
            </div>
          </div>
          <div class="panel-footer">
            <n-button type="primary" size="small" :loading="saving" @click="saveVariables">保存</n-button>
          </div>
        </template>
        <div v-else class="panel-empty">选择一个环境以编辑变量</div>
      </div>
    </div>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { NModal, NButton, NInput, NCheckbox, NDropdown, useMessage } from 'naive-ui'
import {
  ListEnvironments, CreateEnvironment, RenameEnvironment, DeleteEnvironment,
  ListEnvVariables, SaveEnvVariables, SaveEnvBaseURL,
} from '../../../wailsjs/go/main/App'
import { useEnvStore } from '../../stores/env'
import type { Environment, EnvVariable } from '../../types/environment'

interface Props { show: boolean; projectId: number | null }
const props = defineProps<Props>()
const emit = defineEmits<{ 'update:show': [value: boolean]; refresh: [] }>()
const environments = ref<Environment[]>([])
const editingEnvId = ref<number | null>(null)
const renamingId = ref<number | null>(null)
const renameText = ref('')
const variables = ref<EnvVariable[]>([])
const saving = ref(false)
const baseURL = ref('')
const message = useMessage()
const envStore = useEnvStore()

watch(() => props.show, async (v) => { if (v) await loadEnvs() })

async function loadEnvs() {
  if (!props.projectId) return
  try { environments.value = await ListEnvironments(props.projectId) || [] } catch { environments.value = [] }
}

async function selectEnv(env: Environment) {
  editingEnvId.value = env.id; baseURL.value = env.base_url || ''
  try {
    const vars = await ListEnvVariables(env.id)
    variables.value = vars.map(v => ({ id: v.id, key: v.key, value: v.value, enabled: v.enabled }))
  } catch { variables.value = [] }
}

async function addEnv() {
  if (!props.projectId) return
  const name = `环境 ${environments.value.length + 1}`
  try { await CreateEnvironment(props.projectId, name, '', null); await loadEnvs(); emit('refresh') } catch {}
}

function envMenuOptions(env: Environment) {
  return [{ label: '重命名', key: 'rename' }, { label: '复制', key: 'copy' }, { label: '删除', key: 'delete' }]
}

async function onEnvMenu(key: string, env: Environment) {
  switch (key) {
    case 'rename': renamingId.value = env.id; renameText.value = env.name; break
    case 'copy':
      if (!props.projectId) return
      try { await CreateEnvironment(props.projectId, env.name + ' (副本)', env.base_url || '', env.id as any); await loadEnvs(); emit('refresh') } catch {}
      break
    case 'delete':
      try { await DeleteEnvironment(env.id); if (editingEnvId.value === env.id) { editingEnvId.value = null; variables.value = [] }; await loadEnvs(); emit('refresh') } catch {}
      break
  }
}

async function confirmRename(env: Environment) {
  if (renameText.value.trim() && renameText.value !== env.name) {
    try { await RenameEnvironment(env.id, renameText.value.trim()); await loadEnvs(); emit('refresh') } catch {}
  }
  renamingId.value = null
}

function addVariable() { variables.value.push({ id: 0, key: '', value: '', enabled: true }) }
function removeVariable(idx: number) { variables.value.splice(idx, 1) }

async function saveVariables() {
  if (!editingEnvId.value) return
  saving.value = true
  try {
    await SaveEnvVariables(editingEnvId.value, variables.value.map(v => ({ id: v.id, environment_id: editingEnvId.value!, key: v.key, value: v.value, enabled: v.enabled, sort_order: 0, created_at: '' } as any)))
    await SaveEnvBaseURL(editingEnvId.value, baseURL.value)
    if (props.projectId) await envStore.loadEnvironments(props.projectId)
    message.success('环境变量已保存')
    await selectEnv({ id: editingEnvId.value, project_id: 0, name: '', base_url: baseURL.value, is_active: false, created_at: '', updated_at: '' })
  } catch (e: any) { message.error('保存失败: ' + (e?.message || String(e))) }
  saving.value = false
}
</script>

<style scoped>
.manager-body { display: flex; gap: 16px; height: 420px; }
.env-list-panel { width: 200px; border-right: 1px solid var(--border-primary); padding-right: 12px; display: flex; flex-direction: column; }
.var-list-panel { flex: 1; display: flex; flex-direction: column; overflow: hidden; }
.panel-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; font-weight: 600; font-size: var(--fs-base); color: var(--text-primary); }
.base-url-section { margin-bottom: 10px; }
.base-url-section label { display: block; font-size: var(--fs-xs); color: var(--text-muted); margin-bottom: 3px; text-transform: uppercase; letter-spacing: 0.3px; }
.env-list { flex: 1; overflow-y: auto; }
.env-item { display: flex; align-items: center; padding: 5px 8px; border-radius: 4px; cursor: pointer; font-size: var(--fs-sm); color: var(--text-secondary); transition: background var(--transition); }
.env-item:hover { background: var(--bg-hover); }
.env-item.active { background: var(--accent-soft); color: var(--accent); }
.var-table { flex: 1; overflow-y: auto; }
.var-row { display: flex; align-items: center; padding: 3px 0; gap: 6px; }
.var-row.var-header { font-size: var(--fs-sm); font-weight: 600; color: var(--text-muted); padding-bottom: 6px; border-bottom: 1px solid var(--border-primary); margin-bottom: 4px; }
.col-check { width: 24px; text-align: center; flex-shrink: 0; }
.col-key { width: 150px; flex-shrink: 0; }
.col-value { flex: 1; }
.col-action { width: 28px; text-align: center; flex-shrink: 0; }
.panel-footer { margin-top: 10px; text-align: right; }
.panel-empty { flex: 1; display: flex; align-items: center; justify-content: center; color: var(--text-muted); font-size: var(--fs-sm); }
</style>
