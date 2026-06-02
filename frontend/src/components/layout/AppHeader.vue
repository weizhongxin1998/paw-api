<template>
  <div class="header">
    <div class="header-left">
      <button class="back-btn" @click="goHome" title="返回项目列表">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/></svg>
      </button>
      <div class="project-btn-wrapper" ref="projectBtnRef">
        <button class="project-btn" @click="dropdownShow = !dropdownShow">
          <span class="project-name">{{ projectStore.currentProject?.name || '未选择项目' }}</span>
          <svg class="proj-arrow" width="8" height="8" viewBox="0 0 10 10"><path d="M2 3 L5 6 L8 3" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/></svg>
        </button>
        <div v-if="dropdownShow" class="project-dropdown" ref="dropdownRef">
          <div
            v-for="p in projectStore.projects"
            :key="p.id"
            class="project-item"
            :class="{ active: p.id === projectStore.currentId }"
            @click="selectProject(p.id)"
          >
            <span class="project-item-name">{{ p.name }}</span>
            <span class="project-item-desc">{{ p.description }}</span>
          </div>
          <div class="dropdown-divider"></div>
          <div class="dropdown-create-btn" @click="showCreateModal = true; dropdownShow = false">+ 新建项目</div>
          <div class="dropdown-home-btn" @click="goHome">返回项目列表</div>
        </div>
      </div>
      <EnvSelector
        :project-id="projectStore.currentId"
        @update:active-env-id="onActiveEnvChange"
      />
    </div>
    <div class="header-right">
      <button class="header-btn" @click="emit('toggleTheme')" :title="themeLabel">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <circle v-if="props.themeMode === 'dark'" cx="12" cy="12" r="5"/><line v-if="props.themeMode === 'dark'" x1="12" y1="1" x2="12" y2="3"/><line v-if="props.themeMode === 'dark'" x1="12" y1="21" x2="12" y2="23"/><line v-if="props.themeMode === 'dark'" x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line v-if="props.themeMode === 'dark'" x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line v-if="props.themeMode === 'dark'" x1="1" y1="12" x2="3" y2="12"/><line v-if="props.themeMode === 'dark'" x1="21" y1="12" x2="23" y2="12"/><line v-if="props.themeMode === 'dark'" x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line v-if="props.themeMode === 'dark'" x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>
          <path v-else d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
        </svg>
      </button>
      <button class="header-btn" @click="showSettings = true">设置</button>
    </div>
    <SettingsModal v-model:show="showSettings" />

    <div v-if="showCreateModal" class="modal-overlay" @click.self="showCreateModal = false">
      <div class="modal-box">
        <h3>新建项目</h3>
        <label>名称</label>
        <input v-model="newProjectName" placeholder="项目名称" @keydown.enter="onCreateProject" />
        <label>描述</label>
        <input v-model="newProjectDesc" placeholder="项目描述（可选）" @keydown.enter="onCreateProject" />
        <div class="modal-actions">
          <button @click="showCreateModal = false">取消</button>
          <button class="btn-primary" @click="onCreateProject" :disabled="!newProjectName.trim()">创建</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
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

const dropdownShow = ref(false)
const showCreateModal = ref(false)
const showSettings = ref(false)
const newProjectName = ref('')
const newProjectDesc = ref('')
const projectBtnRef = ref<HTMLElement | null>(null)
const dropdownRef = ref<HTMLElement | null>(null)

function selectProject(id: number) {
  dropdownShow.value = false
  emit('projectChanged', id)
}

function goHome() {
  dropdownShow.value = false
  emit('backToHome')
}

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

function onClickOutside(e: MouseEvent) {
  if (dropdownShow.value && dropdownRef.value && !dropdownRef.value.contains(e.target as Node)) {
    dropdownShow.value = false
  }
}

const themeLabel = computed(() => props.themeMode === 'dark' ? '日间模式' : '夜间模式')

onMounted(() => {
  document.addEventListener('click', onClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', onClickOutside)
})
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
  z-index: 10;
}
.header-left { display: flex; align-items: center; gap: 6px; }
.header-right { display: flex; align-items: center; gap: 2px; }

.back-btn {
  display: flex; align-items: center; justify-content: center;
  width: 28px; height: 28px;
  background: transparent; border: 1px solid transparent; color: var(--text-muted);
  cursor: pointer; border-radius: var(--radius-sm);
  transition: all var(--transition);
}
.back-btn:hover { background: var(--bg-hover); color: var(--accent); border-color: var(--border-hover); }

.project-btn-wrapper { position: relative; }
.project-btn {
  display: flex; align-items: center; gap: 4px;
  padding: 4px 8px;
  border: 1px solid var(--border-primary); border-radius: var(--radius);
  background: var(--bg-base); cursor: pointer; font-size: 11px; font-weight: 600;
  color: var(--text-primary); outline: none; font-family: var(--font-mono);
  transition: all var(--transition);
}
.project-btn:hover { background: var(--bg-elevated); border-color: var(--border-hover); }
.project-name { max-width: 150px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.proj-arrow { color: var(--text-muted); flex-shrink: 0; transition: transform var(--transition); }

.project-dropdown {
  position: absolute; top: 100%; left: 0; margin-top: 2px;
  width: 260px; background: var(--bg-elevated); border: 1px solid var(--border-primary);
  border-radius: var(--radius); box-shadow: 0 8px 24px rgba(0,0,0,0.4);
  z-index: 100; overflow: hidden;
  animation: ddFadeIn 0.1s ease;
}
@keyframes ddFadeIn { from { opacity: 0; transform: translateY(-2px); } to { opacity: 1; transform: translateY(0); } }
.project-item {
  padding: 7px 12px; cursor: pointer; display: flex; flex-direction: column; gap: 1px;
  transition: background var(--transition);
}
.project-item:hover { background: var(--bg-hover); }
.project-item.active { background: var(--accent-soft); }
.project-item-name { font-size: 11px; font-weight: 600; color: var(--text-primary); }
.project-item-desc { font-size: 10px; color: var(--text-muted); }
.dropdown-divider { height: 1px; background: var(--border-primary); }
.dropdown-create-btn {
  padding: 7px 12px; font-size: 11px; color: var(--accent); cursor: pointer;
  font-weight: 500; transition: background var(--transition);
}
.dropdown-create-btn:hover { background: var(--accent-soft); }
.dropdown-home-btn {
  padding: 7px 12px; font-size: 11px; color: var(--text-secondary); cursor: pointer;
  transition: background var(--transition);
}
.dropdown-home-btn:hover { background: var(--bg-hover); }

.header-btn {
  font-size: 11px; padding: 4px 8px; background: transparent;
  border: 1px solid transparent; border-radius: var(--radius-sm); color: var(--text-muted);
  cursor: pointer; outline: none; display: flex; align-items: center; gap: 3px;
  transition: all var(--transition); font-family: var(--font-mono);
}
.header-btn:hover { background: var(--bg-hover); color: var(--text-primary); border-color: var(--border-hover); }

.modal-overlay {
  position: fixed; inset: 0; background: rgba(0,0,0,0.6);
  display: flex; align-items: center; justify-content: center; z-index: 200;
}
.modal-box {
  background: var(--bg-surface); border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg); padding: 22px; width: 360px;
  box-shadow: 0 12px 40px rgba(0,0,0,0.5);
}
.modal-box h3 { margin: 0 0 16px; font-size: 14px; font-weight: 600; color: var(--text-primary); }
.modal-box label { display: block; font-size: 10px; color: var(--text-muted); margin-bottom: 3px; letter-spacing: 0.5px; text-transform: uppercase; }
.modal-box input {
  width: 100%; padding: 7px 10px; border: 1px solid var(--border-primary); border-radius: var(--radius);
  font-size: 12px; outline: none; background: var(--bg-base); color: var(--text-primary);
  margin-bottom: 10px; box-sizing: border-box; font-family: var(--font-mono);
  transition: border-color var(--transition);
}
.modal-box input:focus { border-color: var(--accent); }
.modal-actions { display: flex; justify-content: flex-end; gap: 6px; margin-top: 16px; }
.modal-actions button {
  padding: 6px 16px; border: 1px solid var(--border-primary); border-radius: var(--radius);
  font-size: 11px; cursor: pointer; background: var(--bg-base); color: var(--text-secondary);
  transition: all var(--transition); font-family: var(--font-mono);
}
.modal-actions button:hover { border-color: var(--border-hover); color: var(--text-primary); }
.modal-actions .btn-primary { background: var(--accent); color: #000; border-color: var(--accent); font-weight: 600; }
.modal-actions .btn-primary:hover { background: var(--accent-hover); border-color: var(--accent-hover); }
.modal-actions .btn-primary:disabled { background: var(--bg-elevated); border-color: var(--border-primary); color: var(--text-muted); cursor: not-allowed; }
</style>
