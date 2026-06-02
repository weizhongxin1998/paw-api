<template>
  <div class="header">
    <div class="header-left">
      <n-button text size="tiny" @click="goHome" title="返回项目列表" class="back-btn">
        <template #icon>
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/></svg>
        </template>
      </n-button>

      <n-dropdown trigger="click" :options="projectDropdownOptions" @select="onProjectDropdownSelect">
        <n-button size="tiny" class="project-btn">
          {{ projectStore.currentProject?.name || '未选择项目' }}
        </n-button>
      </n-dropdown>

      <EnvSelector
        :project-id="projectStore.currentId"
        @update:active-env-id="onActiveEnvChange"
      />
    </div>

    <div class="header-right">
      <n-button text size="tiny" @click="emit('toggleTheme')" :title="themeLabel">
        <template #icon>
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
            <circle v-if="props.themeMode === 'dark'" cx="12" cy="12" r="5"/>
            <path v-else d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
          </svg>
        </template>
      </n-button>
      <n-button text size="tiny" @click="showSettings = true">设置</n-button>
    </div>

    <SettingsModal v-model:show="showSettings" />

    <n-modal v-model:show="showCreateModal" preset="card" title="新建项目" style="width: 360px" :mask-closable="false">
      <n-form label-placement="top">
        <n-form-item label="名称">
          <n-input v-model:value="newProjectName" placeholder="项目名称" @keydown.enter="onCreateProject" />
        </n-form-item>
        <n-form-item label="描述">
          <n-input v-model:value="newProjectDesc" placeholder="项目描述（可选）" @keydown.enter="onCreateProject" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button @click="showCreateModal = false">取消</n-button>
        <n-button type="primary" :disabled="!newProjectName.trim()" @click="onCreateProject">创建</n-button>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { NButton, NDropdown, NModal, NForm, NFormItem, NInput } from 'naive-ui'
import { useProjectStore } from '../../stores/project'
import { useEnvStore } from '../../stores/env'
import EnvSelector from '../environment/EnvSelector.vue'
import SettingsModal from '../modals/SettingsModal.vue'

const props = defineProps<{ themeMode: 'dark' | 'light' }>()
const projectStore = useProjectStore()
const envStore = useEnvStore()

const emit = defineEmits<{
  'projectChanged': [id: number]
  'backToHome': []
  'toggleTheme': []
}>()

const showCreateModal = ref(false)
const showSettings = ref(false)
const newProjectName = ref('')
const newProjectDesc = ref('')

const themeLabel = computed(() => props.themeMode === 'dark' ? '日间模式' : '夜间模式')

const projectDropdownOptions = computed(() => {
  const items: Array<{ label?: string; key: string; disabled?: boolean; type?: string }> = projectStore.projects.map(p => ({
    label: p.name + (p.description ? ' - ' + p.description : ''),
    key: String(p.id),
    disabled: p.id === projectStore.currentId,
  }))
  items.push(
    { type: 'divider', key: 'div' },
    { label: '+ 新建项目', key: '__create__' },
    { label: '返回项目列表', key: '__home__' },
  )
  return items as any
})

function onProjectDropdownSelect(key: string) {
  if (key === '__create__') { showCreateModal.value = true; return }
  if (key === '__home__') { emit('backToHome'); return }
  emit('projectChanged', Number(key))
}

function goHome() { emit('backToHome') }

async function onCreateProject() {
  const name = newProjectName.value.trim()
  if (!name) return
  const desc = newProjectDesc.value.trim()
  const p = await projectStore.createProject(name, desc)
  showCreateModal.value = false
  newProjectName.value = ''
  newProjectDesc.value = ''
  emit('projectChanged', p.id)
}

function onActiveEnvChange(id: number | null) {
  envStore.activeEnvId = id
}
</script>

<style scoped>
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 10px;
  height: 38px;
  background: var(--bg-surface);
  border-bottom: 1px solid var(--border-primary);
  flex-shrink: 0;
}
.header-left { display: flex; align-items: center; gap: 10px; }
.header-left > :first-child { color: var(--text-muted); }
.header-right { display: flex; align-items: center; gap: 4px; }
.project-btn { font-weight: 600; max-width: 180px; }
</style>
