<template>
  <n-config-provider :theme-overrides="themeOverrides">
    <n-dialog-provider>
      <n-message-provider>
        <div class="app-container">
          <AppHeader
            v-if="projectStore.currentId"
            :project-id="projectStore.currentId"
            @project-changed="onProjectChanged"
            @back-to-home="projectStore.currentId = null"
          />

          <ProjectHome
            v-if="!projectStore.currentId"
            @enter-project="onEnterProject"
          />

          <AppBody v-else :project-id="projectStore.currentId" />
        </div>
      </n-message-provider>
    </n-dialog-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { NConfigProvider, NMessageProvider, NDialogProvider } from 'naive-ui'
import AppHeader from './components/layout/AppHeader.vue'
import AppBody from './components/layout/AppBody.vue'
import ProjectHome from './components/layout/ProjectHome.vue'
import { useProjectStore } from './stores/project'
import { useCollectionStore } from './stores/collection'
import { useEnvStore } from './stores/env'

const projectStore = useProjectStore()
const collectionStore = useCollectionStore()
const envStore = useEnvStore()

async function onEnterProject(id: number) {
  await projectStore.switchProject(id)
  collectionStore.loadTree(id)
  envStore.loadEnvironments(id)
}

async function onProjectChanged(id: number) {
  await projectStore.switchProject(id)
  collectionStore.loadTree(id)
  envStore.loadEnvironments(id)
}

onMounted(async () => {
  await projectStore.loadProjects()
})

const themeOverrides = {
  common: {
    primaryColor: '#18a058',
    primaryColorHover: '#0c7a43',
  },
}
</script>

<style>
:root {
  --green: #18a058;
  --green-hover: #0c7a43;
  --green-soft: #e6f7ec;
  --red: #d03050;
  --red-soft: #fef0f0;
  --amber: #d97706;
  --amber-soft: #fffbeb;
  --blue: #2563eb;
  --blue-soft: #eff6ff;
  --purple: #7c3aed;
  --purple-soft: #f5f3ff;
  --gray-50: #f9fafb;
  --gray-100: #f3f4f6;
  --gray-200: #e5e7eb;
  --gray-300: #d1d5db;
  --gray-400: #9ca3af;
  --gray-500: #6b7280;
  --gray-600: #4b5563;
  --gray-700: #374151;
  --gray-800: #1f2937;
  --gray-900: #111827;
  --radius-sm: 5px;
  --radius: 7px;
  --radius-lg: 10px;
  --shadow-sm: 0 1px 2px rgba(0,0,0,0.04);
  --shadow: 0 1px 3px rgba(0,0,0,0.06), 0 1px 2px rgba(0,0,0,0.04);
  --shadow-md: 0 4px 6px rgba(0,0,0,0.04), 0 2px 4px rgba(0,0,0,0.04);
  --shadow-lg: 0 10px 25px rgba(0,0,0,0.08), 0 4px 10px rgba(0,0,0,0.04);
  --transition: 0.15s ease;
}

html, body, #app {
  margin: 0;
  padding: 0;
  height: 100%;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  font-size: 13px;
  color: var(--gray-700);
  background: var(--gray-50);
  -webkit-font-smoothing: antialiased;
}

.app-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
}
</style>
