<template>
  <n-config-provider :theme-overrides="themeOverrides">
    <n-message-provider>
      <div class="app-container">
        <AppHeader :project-id="projectStore.currentId" @project-changed="onProjectChanged" />

        <div v-if="!projectStore.currentId && projectStore.projects.length === 0" class="welcome-screen">
          <div class="welcome-icon">📦</div>
          <h2>欢迎使用 Paw API</h2>
          <p>创建你的第一个项目，开始调试 API</p>
          <button class="welcome-btn" @click="showCreateModal = true">+ 新建项目</button>
        </div>

        <div v-else-if="!projectStore.currentId && projectStore.projects.length > 0" class="welcome-screen">
          <p>请从标题栏选择一个项目</p>
        </div>

        <AppBody v-else :project-id="projectStore.currentId" />

        <div v-if="showCreateModal" class="modal-overlay" @click.self="showCreateModal = false">
          <div class="modal-box">
            <h3>新建项目</h3>
            <label>名称</label>
            <input v-model="newProjectName" placeholder="项目名称" />
            <label>描述</label>
            <input v-model="newProjectDesc" placeholder="描述（可选）" />
            <div class="modal-acts">
              <button @click="showCreateModal = false">取消</button>
              <button class="btn-primary" @click="onCreateProject" :disabled="!newProjectName.trim()">创建</button>
            </div>
          </div>
        </div>
      </div>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NConfigProvider, NMessageProvider } from 'naive-ui'
import AppHeader from './components/layout/AppHeader.vue'
import AppBody from './components/layout/AppBody.vue'
import { useProjectStore } from './stores/project'
import { useCollectionStore } from './stores/collection'
import { useEnvStore } from './stores/env'

const projectStore = useProjectStore()
const collectionStore = useCollectionStore()
const envStore = useEnvStore()

const showCreateModal = ref(false)
const newProjectName = ref('')
const newProjectDesc = ref('')

async function onProjectChanged(id: number) {
  await projectStore.switchProject(id)
  collectionStore.loadTree(id)
  envStore.activeEnvId = null
}

async function onCreateProject() {
  const name = newProjectName.value.trim()
  if (!name) return
  const p = await projectStore.createProject(name, newProjectDesc.value.trim())
  showCreateModal.value = false
  newProjectName.value = ''
  newProjectDesc.value = ''
  await onProjectChanged(p.id)
}

onMounted(async () => {
  await projectStore.loadProjects()
  await projectStore.loadLastProject()
  if (projectStore.currentId) {
    collectionStore.loadTree(projectStore.currentId)
  }
})

const themeOverrides = {
  common: {
    primaryColor: '#18a058',
    primaryColorHover: '#0c7a43',
  },
}
</script>

<style>
html, body, #app {
  margin: 0;
  padding: 0;
  height: 100%;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}
.app-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
}
.welcome-screen {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: #fff;
}
.welcome-icon { font-size: 48px; margin-bottom: 12px; opacity: 0.3; }
.welcome-screen h2 { font-size: 20px; color: #555; margin: 0 0 6px; font-weight: 500; }
.welcome-screen p { color: #999; font-size: 13px; margin: 0 0 20px; }
.welcome-btn {
  padding: 10px 28px;
  background: #18a058;
  color: #fff;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
}
.welcome-btn:hover { background: #0c7a43; }

.modal-overlay {
  position: fixed; inset: 0;
  background: rgba(0,0,0,0.3);
  display: flex; align-items: center; justify-content: center;
  z-index: 200;
}
.modal-box {
  background: #fff; border-radius: 10px;
  padding: 24px; width: 380px;
  box-shadow: 0 8px 30px rgba(0,0,0,0.15);
}
.modal-box h3 { margin: 0 0 16px; font-size: 16px; }
.modal-box label { display: block; font-size: 12px; color: #888; margin-bottom: 4px; }
.modal-box input {
  width: 100%; padding: 7px 10px; border: 1px solid #ddd; border-radius: 6px;
  font-size: 13px; outline: none; margin-bottom: 10px; box-sizing: border-box;
}
.modal-box input:focus { border-color: #18a058; }
.modal-acts { display: flex; justify-content: flex-end; gap: 8px; margin-top: 14px; }
.modal-acts button {
  padding: 6px 18px; border: 1px solid #ddd; border-radius: 6px;
  font-size: 12px; cursor: pointer; background: #fff;
}
.modal-acts .btn-primary {
  background: #18a058; color: #fff; border-color: #18a058;
}
.modal-acts .btn-primary:hover { background: #0c7a43; }
.modal-acts .btn-primary:disabled { background: #aaa; border-color: #aaa; cursor: not-allowed; }
</style>
