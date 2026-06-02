<template>
  <div class="header">
    <div class="header-left">
      <button class="back-btn" @click="goHome" title="返回项目列表">
        <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/></svg>
      </button>
      <div class="project-btn-wrapper" ref="projectBtnRef">
        <button class="project-btn" @click="dropdownShow = !dropdownShow">
          <span class="project-name">{{ projectStore.currentProject?.name || '未选择项目' }}</span>
          <svg class="proj-arrow" width="10" height="10" viewBox="0 0 10 10"><path d="M2 3 L5 6 L8 3" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
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
          <div class="dropdown-create-btn" @click="showCreateModal = true; dropdownShow = false">
            + 新建项目
          </div>
          <div class="dropdown-home-btn" @click="goHome">
            返回项目列表
          </div>
        </div>
      </div>
      <EnvSelector
        :project-id="projectStore.currentId"
        @update:active-env-id="onActiveEnvChange"
      />
    </div>
    <div class="header-right">
      <button class="header-btn">日间</button>
      <button class="header-btn" @click="showSettings = true">设置</button>
    </div>
    <SettingsModal v-model:show="showSettings" />

    <div v-if="showCreateModal" class="modal-overlay" @click.self="showCreateModal = false">
      <div class="modal-box">
        <h3>新建项目</h3>
        <div class="modal-form">
          <label>名称</label>
          <input v-model="newProjectName" placeholder="项目名称" @keydown.enter="onCreateProject" />
          <label>描述</label>
          <input v-model="newProjectDesc" placeholder="项目描述（可选）" @keydown.enter="onCreateProject" />
        </div>
        <div class="modal-actions">
          <button @click="showCreateModal = false">取消</button>
          <button class="btn-primary" @click="onCreateProject" :disabled="!newProjectName.trim()">创建</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useProjectStore } from '../../stores/project'
import { useEnvStore } from '../../stores/env'
import EnvSelector from '../environment/EnvSelector.vue'
import SettingsModal from '../modals/SettingsModal.vue'

const projectStore = useProjectStore()
const envStore = useEnvStore()

const emit = defineEmits<{
  'projectChanged': [id: number]
  'backToHome': []
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
  padding: 0 12px;
  height: 42px;
  background: #fff;
  border-bottom: 1px solid var(--gray-200);
  flex-shrink: 0;
  z-index: 10;
}
.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}
.header-right {
  display: flex;
  align-items: center;
  gap: 4px;
}

.back-btn {
  display: flex; align-items: center; justify-content: center;
  width: 30px; height: 30px;
  background: transparent; border: none; color: var(--gray-400);
  cursor: pointer; border-radius: var(--radius-sm);
  transition: all var(--transition);
}
.back-btn:hover { background: var(--gray-100); color: var(--green); }

.project-btn-wrapper { position: relative; }
.project-btn {
  display: flex; align-items: center; gap: 5px;
  padding: 5px 10px;
  border: 1px solid var(--gray-200); border-radius: var(--radius);
  background: #fff; cursor: pointer; font-size: 12px; font-weight: 600;
  color: var(--gray-700); outline: none;
  transition: all var(--transition);
}
.project-btn:hover { background: var(--gray-50); border-color: var(--gray-300); }
.project-name { max-width: 140px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.proj-arrow { color: var(--gray-400); flex-shrink: 0; transition: transform var(--transition); }
.project-btn-wrapper:hover .proj-arrow { color: var(--gray-500); }

.project-dropdown {
  position: absolute; top: 100%; left: 0; margin-top: 4px;
  width: 260px; background: #fff; border: 1px solid var(--gray-200);
  border-radius: var(--radius-lg); box-shadow: var(--shadow-lg);
  z-index: 100; overflow: hidden;
  animation: ddFadeIn 0.12s ease;
}
@keyframes ddFadeIn { from { opacity: 0; transform: translateY(-4px); } to { opacity: 1; transform: translateY(0); } }
.project-item {
  padding: 8px 14px; cursor: pointer; display: flex; flex-direction: column; gap: 2px;
  transition: background var(--transition);
}
.project-item:hover { background: var(--gray-50); }
.project-item.active { background: var(--green-soft); }
.project-item-name { font-size: 12px; font-weight: 600; color: var(--gray-700); }
.project-item-desc { font-size: 10px; color: var(--gray-400); }
.dropdown-divider { height: 1px; background: var(--gray-100); }
.dropdown-create-btn {
  padding: 8px 14px; font-size: 12px; color: var(--green); cursor: pointer;
  text-align: left; font-weight: 500; transition: background var(--transition);
}
.dropdown-create-btn:hover { background: var(--green-soft); }
.dropdown-home-btn {
  padding: 8px 14px; font-size: 12px; color: var(--gray-500); cursor: pointer;
  text-align: left; transition: background var(--transition);
}
.dropdown-home-btn:hover { background: var(--gray-50); }

.header-btn {
  font-size: 11px; padding: 5px 10px; background: transparent;
  border: 1px solid transparent; border-radius: var(--radius-sm); color: var(--gray-500);
  cursor: pointer; outline: none; transition: all var(--transition);
}
.header-btn:hover { background: var(--gray-100); color: var(--gray-700); }

.modal-overlay {
  position: fixed; inset: 0; background: rgba(0,0,0,0.3);
  display: flex; align-items: center; justify-content: center; z-index: 200;
}
.modal-box {
  background: #fff; border-radius: var(--radius-lg); padding: 26px; width: 380px;
  box-shadow: var(--shadow-lg);
}
.modal-box h3 { margin: 0 0 18px; font-size: 16px; font-weight: 600; color: var(--gray-800); }
.modal-form { display: flex; flex-direction: column; gap: 8px; }
.modal-form label { font-size: 11px; color: var(--gray-500); text-transform: uppercase; letter-spacing: 0.3px; }
.modal-form input {
  padding: 8px 10px; border: 1px solid var(--gray-200); border-radius: var(--radius);
  font-size: 13px; outline: none; transition: border-color var(--transition);
}
.modal-form input:focus { border-color: var(--green); box-shadow: 0 0 0 3px rgba(24,160,88,0.1); }
.modal-actions {
  display: flex; justify-content: flex-end; gap: 8px; margin-top: 22px;
}
.modal-actions button {
  padding: 7px 18px; border: 1px solid var(--gray-200); border-radius: var(--radius);
  font-size: 12px; cursor: pointer; background: #fff; color: var(--gray-600);
  transition: all var(--transition);
}
.modal-actions button:hover { border-color: var(--gray-300); }
.modal-actions .btn-primary { background: var(--green); color: #fff; border-color: var(--green); }
.modal-actions .btn-primary:hover { background: var(--green-hover); border-color: var(--green-hover); }
.modal-actions .btn-primary:disabled { background: var(--gray-300); border-color: var(--gray-300); cursor: not-allowed; }
</style>
