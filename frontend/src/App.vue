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
import { ref, onMounted } from 'vue'
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
</style>
