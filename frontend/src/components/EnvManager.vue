<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { NModal, NButton, NInput, NIcon, NForm, NFormItem, NSpace, NTabPane, NTabs, useMessage } from 'naive-ui'
import { Add, Trash } from '@vicons/ionicons5'
import { useEnvironmentStore, type EnvVariable } from '../stores/environment'
import {
  ListEnvironments,
  CreateEnvironment,
  UpdateEnvironment,
  DeleteEnvironment,
  SetActiveEnvironment,
} from '../../wailsjs/go/handlers/EnvironmentHandler'
import { useProjectStore } from '../stores/project'

const { t } = useI18n()
const props = defineProps<{ show: boolean }>()
const emit = defineEmits<{ (e: 'update:show', val: boolean): void }>()

const envStore = useEnvironmentStore()
const projectStore = useProjectStore()
const message = useMessage()

const editingEnvId = ref<string | null>(null)
const envName = ref('')
const variables = ref<EnvVariable[]>([])

function resetForm() {
  editingEnvId.value = null
  envName.value = ''
  variables.value = []
}

function editEnv(env: typeof envStore.environments[0]) {
  editingEnvId.value = env.id
  envName.value = env.name
  try { variables.value = JSON.parse(env.variables) } catch { variables.value = [] }
}

function addVariable() { variables.value.push({ key: '', value: '', enabled: true }) }
function removeVariable(index: number) { variables.value.splice(index, 1) }

async function save() {
  if (!envName.value.trim()) { message.warning(t('env.nameRequired')); return }
  if (!projectStore.currentProject) { message.error(t('common.error')); return }
  const varsJSON = JSON.stringify(variables.value)
  try {
    if (editingEnvId.value) {
      const updated = await UpdateEnvironment(editingEnvId.value, envName.value.trim(), varsJSON)
      envStore.updateEnvironment(updated)
      message.success(t('env.updated'))
    } else {
      const created = await CreateEnvironment(projectStore.currentProject.id, envName.value.trim(), varsJSON, envStore.environments.length === 0)
      envStore.addEnvironment(created)
      if (created.is_active) envStore.setActiveEnvironment(created)
      message.success(t('env.created'))
    }
    resetForm()
  } catch (e: any) { message.error(e.message || t('env.failedSave')) }
}

async function removeEnv(id: string) {
  try {
    await DeleteEnvironment(id)
    envStore.removeEnvironment(id)
    message.success(t('env.deleted'))
  } catch (e: any) { message.error(e.message || t('env.failedDelete')) }
}

async function setActive(id: string) {
  if (!projectStore.currentProject) return
  try {
    const env = await SetActiveEnvironment(id, projectStore.currentProject.id)
    envStore.setActiveEnvironment(env)
    message.success(`${env.name} ${t('env.activeNow')}`)
  } catch (e: any) { message.error(e.message || t('env.failedActivate')) }
}

onMounted(async () => {
  if (projectStore.currentProject) {
    try {
      const envs = await ListEnvironments(projectStore.currentProject.id)
      envStore.setEnvironments(envs)
      const active = envs.find(e => e.is_active)
      if (active) envStore.setActiveEnvironment(active)
    } catch { console.error('Failed to load environments') }
  }
})
</script>

<template>
  <NModal :show="props.show" @update:show="emit('update:show', $event)" :title="t('env.manager')" preset="card" style="width: 520px">
    <NTabs type="line">
      <NTabPane name="list" :tab="t('env.environments')">
        <div v-if="envStore.environments.length === 0" class="empty-hint">{{ t('env.noEnvs') }}</div>
        <div v-for="env in envStore.environments" :key="env.id" class="env-row">
          <div class="env-info" @click="editEnv(env)">
            <span class="env-name">{{ env.name }}</span>
            <span v-if="envStore.activeEnvironment?.id === env.id" class="active-badge">{{ t('env.active') }}</span>
          </div>
          <NSpace>
            <NButton size="tiny" quaternary @click.stop="setActive(env.id)">{{ t('env.activate') }}</NButton>
            <NButton size="tiny" quaternary @click.stop="editEnv(env)">{{ t('env.edit') }}</NButton>
            <NButton size="tiny" quaternary @click.stop="removeEnv(env.id)">
              <template #icon><NIcon><Trash /></NIcon></template>
            </NButton>
          </NSpace>
        </div>
      </NTabPane>
      <NTabPane name="edit" :tab="t('common.edit')">
        <NForm>
          <NFormItem :label="t('env.name')">
            <NInput v-model:value="envName" :placeholder="t('env.namePlaceholder')" />
          </NFormItem>
          <NFormItem :label="t('env.variables')">
            <div class="var-list">
              <div v-for="(v, i) in variables" :key="i" class="var-row">
                <input type="checkbox" v-model="v.enabled" class="var-checkbox" />
                <NInput v-model:value="v.key" size="tiny" :placeholder="t('env.varKey')" class="var-input" />
                <NInput v-model:value="v.value" size="tiny" :placeholder="t('env.varValue')" class="var-input" />
                <NButton quaternary circle size="tiny" @click="removeVariable(i)">
                  <template #icon><NIcon><Trash /></NIcon></template>
                </NButton>
              </div>
              <NButton size="tiny" quaternary @click="addVariable">
                <template #icon><NIcon><Add /></NIcon></template>
                {{ t('env.addVariable') }}
              </NButton>
            </div>
          </NFormItem>
        </NForm>
        <NSpace justify="end" style="margin-top: 12px;">
          <NButton @click="resetForm">{{ t('env.reset') }}</NButton>
          <NButton type="primary" @click="save">{{ t('env.save') }}</NButton>
        </NSpace>
      </NTabPane>
    </NTabs>
  </NModal>
</template>

<style scoped>
.env-row { display: flex; align-items: center; justify-content: space-between; padding: 8px 4px; border-bottom: 1px solid var(--border-color); }
.env-info { cursor: pointer; flex: 1; }
.env-name { font-weight: 500; }
.active-badge { margin-left: 8px; font-size: 10px; background: #18a058; color: white; padding: 1px 6px; border-radius: 3px; }
.empty-hint { color: #999; font-size: 13px; padding: 24px 0; text-align: center; }
.var-list { width: 100%; }
.var-row { display: flex; align-items: center; gap: 4px; margin-bottom: 4px; }
.var-checkbox { width: 14px; height: 14px; flex-shrink: 0; }
.var-input { flex: 1; }
</style>
