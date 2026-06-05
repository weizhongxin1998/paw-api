<template>
  <div class="header">
    <div class="header-left">
      <n-button text size="tiny" @click="goHome" title="返回项目列表" class="back-btn">
        <template #icon>
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="19" y1="12" x2="5" y2="12"/>
            <polyline points="12 19 5 12 12 5"/>
          </svg>
        </template>
      </n-button>

      <n-dropdown
        trigger="click"
        :options="projectDropdownOptions"
        @select="onProjectDropdownSelect"
        ref="projectDropdownRef"
      >
        <n-button size="tiny" class="project-btn" ref="projectBtnRef">
          <span class="project-btn-label">{{ projectStore.currentProject?.name || '未选择项目' }}</span>
          <svg class="project-btn-chevron" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="6 9 12 15 18 9"/>
          </svg>
        </n-button>
      </n-dropdown>

      <span class="header-divider"></span>

      <EnvSelector
        :project-id="projectStore.currentId"
        @update:active-env-id="onActiveEnvChange"
      />
    </div>

    <div class="header-right">
      <n-button text size="tiny" @click="emit('toggleTheme')" :title="themeLabel" class="icon-btn">
        <template #icon>
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
            <circle v-if="props.themeMode === 'dark'" cx="12" cy="12" r="4"/>
            <line v-if="props.themeMode === 'dark'" x1="12" y1="2" x2="12" y2="4"/>
            <line v-if="props.themeMode === 'dark'" x1="12" y1="20" x2="12" y2="22"/>
            <line v-if="props.themeMode === 'dark'" x1="4.93" y1="4.93" x2="6.34" y2="6.34"/>
            <line v-if="props.themeMode === 'dark'" x1="17.66" y1="17.66" x2="19.07" y2="19.07"/>
            <line v-if="props.themeMode === 'dark'" x1="2" y1="12" x2="4" y2="12"/>
            <line v-if="props.themeMode === 'dark'" x1="20" y1="12" x2="22" y2="12"/>
            <line v-if="props.themeMode === 'dark'" x1="4.93" y1="19.07" x2="6.34" y2="17.66"/>
            <line v-if="props.themeMode === 'dark'" x1="17.66" y1="6.34" x2="19.07" y2="4.93"/>
            <path v-else d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
          </svg>
        </template>
      </n-button>
      <n-button text size="tiny" @click="showSettings = true" class="settings-btn">
        <template #icon>
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/>
          </svg>
        </template>
      </n-button>
    </div>

    <SettingsModal v-model:show="showSettings" />

    <n-modal
      v-model:show="showCreateModal"
      preset="card"
      title="新建项目"
      style="width: 380px"
      :mask-closable="false"
      @after-enter="onCreateModalOpened"
    >
      <n-form label-placement="top">
        <n-form-item label="名称">
          <n-input
            ref="createNameInputRef"
            v-model:value="newProjectName"
            placeholder="项目名称"
            @keydown.enter="onCreateProject"
            @keydown="onCreateModalKeydown"
          />
        </n-form-item>
        <n-form-item label="描述">
          <n-input
            v-model:value="newProjectDesc"
            placeholder="项目描述（可选）"
            @keydown.enter="onCreateProject"
            @keydown="onCreateModalKeydown"
          />
        </n-form-item>
      </n-form>
      <template #footer>
        <div class="modal-footer">
          <n-button @click="showCreateModal = false">取消</n-button>
          <n-button type="primary" :disabled="!newProjectName.trim()" @click="onCreateProject">
            创建
            <span class="shortcut-hint">Ctrl+Enter</span>
          </n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { NButton, NDropdown, NModal, NForm, NFormItem, NInput, useMessage } from 'naive-ui'
import { useProjectStore } from '../../stores/project'
import { useEnvStore } from '../../stores/env'
import EnvSelector from '../environment/EnvSelector.vue'
import SettingsModal from '../modals/SettingsModal.vue'

const props = defineProps<{ themeMode: 'dark' | 'light' }>()
const projectStore = useProjectStore()
const envStore = useEnvStore()
const message = useMessage()

const emit = defineEmits<{
  'projectChanged': [id: number]
  'backToHome': []
  'toggleTheme': []
}>()

const showCreateModal = ref(false)
const showSettings = ref(false)
const newProjectName = ref('')
const newProjectDesc = ref('')

const projectBtnRef = ref<InstanceType<typeof NButton> | null>(null)
const projectDropdownRef = ref<any>(null)
const createNameInputRef = ref<InstanceType<typeof NInput> | null>(null)

const themeLabel = computed(() => props.themeMode === 'dark' ? '日间模式' : '夜间模式')

const projectDropdownOptions = computed(() => {
  const items: Array<{ label?: string; key: string; disabled?: boolean; type?: string }> = projectStore.projects.map(p => ({
    label: p.name + (p.description ? ' - ' + p.description : ''),
    key: String(p.id),
    disabled: p.id === projectStore.currentId,
  }))
  items.push(
    { type: 'divider', key: 'div' },
    { label: '+ 新建项目          Ctrl+N', key: '__create__' },
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

function openProjectSwitcher() {
  // Simulate click on the project button to open the dropdown
  const btnEl = projectBtnRef.value?.$el as HTMLElement | undefined
  if (btnEl) {
    btnEl.click()
  }
}

function onCreateModalOpened() {
  nextTick(() => {
    const inputEl = createNameInputRef.value?.$el?.querySelector('input') as HTMLInputElement | null
    if (inputEl) {
      inputEl.focus()
      inputEl.select()
    }
  })
}

function onCreateModalKeydown(e: KeyboardEvent) {
  if (e.ctrlKey && e.key === 'Enter') {
    e.preventDefault()
    onCreateProject()
  }
}

async function onCreateProject() {
  const name = newProjectName.value.trim()
  if (!name) return
  const desc = newProjectDesc.value.trim()
  try {
    const p = await projectStore.createProject(name, desc)
    showCreateModal.value = false
    newProjectName.value = ''
    newProjectDesc.value = ''
    message.success(`已创建项目 "${name}"`)
    emit('projectChanged', p.id)
  } catch (e: any) {
    message.error('创建失败: ' + (e?.message || String(e)))
  }
}

function onActiveEnvChange(id: number | null) {
  envStore.activeEnvId = id
}

function onGlobalKeydown(e: KeyboardEvent) {
  // Ctrl+P: open project switcher
  if (e.ctrlKey && !e.shiftKey && e.key === 'p') {
    e.preventDefault()
    openProjectSwitcher()
  }
}

onMounted(() => {
  document.addEventListener('keydown', onGlobalKeydown)
})
onUnmounted(() => {
  document.removeEventListener('keydown', onGlobalKeydown)
})
</script>

<style scoped>
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 12px;
  height: 40px;
  background: linear-gradient(180deg, var(--bg-elevated) 0%, var(--bg-surface) 100%);
  border-bottom: 1px solid var(--border-primary);
  flex-shrink: 0;
  position: relative;
  z-index: 10;
}
/* Always-visible accent gradient line at bottom */
.header::after {
  content: '';
  position: absolute;
  bottom: -1px; left: 0; right: 0;
  height: 1.5px;
  background: linear-gradient(90deg, transparent 5%, var(--accent-glow) 30%, var(--accent) 50%, var(--accent-glow) 70%, transparent 95%);
  opacity: 0.5;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

/* Back button: circular hover area with arrow-left icon */
.back-btn {
  color: var(--text-muted) !important;
  width: 28px; height: 28px;
  display: flex; align-items: center; justify-content: center;
  border-radius: 50% !important;
  transition: all var(--transition);
}
.back-btn:hover {
  color: var(--text-primary) !important;
  background: var(--bg-hover) !important;
}
.back-btn:active {
  background: var(--bg-active) !important;
  transform: scale(0.92);
}
.back-btn svg { width: 15px; height: 15px; }

/* Project button with dropdown chevron */
.project-btn {
  font-weight: 600;
  max-width: 280px;
  height: 28px;
  border-radius: var(--radius-sm) !important;
  transition: all var(--transition);
  letter-spacing: -0.01em;
  display: inline-flex;
  align-items: center;
  gap: 6px;
}
.project-btn-label {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 230px;
}
.project-btn-chevron {
  flex-shrink: 0;
  opacity: 0.5;
  transition: opacity var(--transition);
}
.project-btn:hover .project-btn-chevron {
  opacity: 0.9;
}

/* Divider between project btn and env selector */
.header-divider {
  display: inline-block;
  width: 1px;
  height: 16px;
  background: var(--border-primary);
  flex-shrink: 0;
  opacity: 0.7;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 2px;
}
.header-right .n-button {
  color: var(--text-muted) !important;
  border-radius: var(--radius-sm) !important;
  height: 28px;
  transition: all var(--transition);
}
.header-right .n-button:hover {
  color: var(--text-primary) !important;
  background: var(--bg-hover) !important;
}
.header-right .n-button svg { width: 15px; height: 15px; }

/* Modal footer with shortcut hint */
.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  align-items: center;
}
.shortcut-hint {
  display: inline-block;
  margin-left: 8px;
  font-size: var(--fs-2xs);
  opacity: 0.45;
  font-weight: 400;
  letter-spacing: 0.02em;
  font-family: var(--font-mono);
}
</style>
