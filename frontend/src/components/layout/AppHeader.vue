<template>
  <div class="header">
    <div class="header-left">
      <button class="back-btn" @click="goHome" title="返回项目列表">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/></svg>
      </button>
      <div class="project-btn-wrapper" ref="projectBtnRef">
        <button class="project-btn" @click="dropdownShow = !dropdownShow">
          {{ projectStore.currentProject?.name || '未选择项目' }}
          <span class="proj-arrow">&#x25BC;</span>
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
          <div class="dropdown-create-btn" style="color:#888" @click="goHome">
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
  padding: 0 14px;
  height: 40px;
  background: #f0f0f0;
  border-bottom: 1px solid #ddd;
  flex-shrink: 0;
  z-index: 10;
}
.header-left {
  display: flex;
  align-items: center;
  gap: 14px;
}
.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}
.project-btn-wrapper {
  position: relative;
}
.project-btn {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 4px 10px;
  border: 1px solid #ddd;
  border-radius: 6px;
  background: #fff;
  cursor: pointer;
  font-size: 13px;
  font-weight: 600;
  color: #333;
  outline: none;
}
.project-btn:hover {
  background: #f8f8f8;
}
.proj-arrow {
  font-size: 9px;
  color: #888;
}
.project-dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  margin-top: 4px;
  width: 260px;
  background: #fff;
  border: 1px solid #ddd;
  border-radius: 8px;
  box-shadow: 0 4px 16px rgba(0,0,0,0.1);
  z-index: 100;
  overflow: hidden;
}
.project-item {
  padding: 8px 12px;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.project-item:hover {
  background: #f5f5f5;
}
.project-item.active {
  background: #e8f5e9;
}
.project-item-name {
  font-size: 12px;
  font-weight: 600;
  color: #333;
}
.project-item-desc {
  font-size: 10px;
  color: #999;
}
.dropdown-divider {
  height: 1px;
  background: #eee;
}
.dropdown-create-btn {
  padding: 8px 12px;
  font-size: 12px;
  color: #18a058;
  cursor: pointer;
  text-align: center;
  font-weight: 600;
}
.dropdown-create-btn:hover {
  background: #f0faf3;
}
.header-btn {
  font-size: 12px;
  padding: 3px 10px;
  background: transparent;
  border: 1px solid transparent;
  border-radius: 4px;
  color: #666;
  cursor: pointer;
  outline: none;
}
.header-btn:hover {
  background: #e8e8e8;
}
.back-btn {
  background: transparent;
  border: none;
  font-size: 18px;
  cursor: pointer;
  color: #888;
  padding: 0 6px;
  line-height: 1;
}
.back-btn:hover { color: #18a058; }

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 200;
}
.modal-box {
  background: #fff;
  border-radius: 10px;
  padding: 24px;
  width: 380px;
  box-shadow: 0 8px 30px rgba(0,0,0,0.15);
}
.modal-box h3 {
  margin: 0 0 16px;
  font-size: 16px;
  color: #333;
}
.modal-form {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.modal-form label {
  font-size: 12px;
  color: #888;
}
.modal-form input {
  padding: 8px 10px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 13px;
  outline: none;
}
.modal-form input:focus {
  border-color: #18a058;
}
.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 20px;
}
.modal-actions button {
  padding: 6px 18px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 12px;
  cursor: pointer;
  background: #fff;
}
.modal-actions .btn-primary {
  background: #18a058;
  color: #fff;
  border-color: #18a058;
}
.modal-actions .btn-primary:hover {
  background: #0c7a43;
}
.modal-actions .btn-primary:disabled {
  background: #aaa;
  border-color: #aaa;
  cursor: not-allowed;
}
</style>
