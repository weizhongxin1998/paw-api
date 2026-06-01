<script lang="ts" setup>
import { computed } from 'vue'
import TabBar from '../components/TabBar.vue'
import RequestEditor from '../components/RequestEditor.vue'
import ResponseViewer from '../components/ResponseViewer.vue'
import WebSocketTester from './WebSocketTester.vue'
import { useTabsStore } from '../stores/tabs'

const tabsStore = useTabsStore()
const activeType = computed(() => tabsStore.activeTab?.type)
</script>

<template>
  <div class="workspace">
    <TabBar />
    <div class="workspace-body">
      <RequestEditor v-if="activeType === 'http'" />
      <WebSocketTester v-else-if="activeType === 'websocket'" />
      <div v-else class="workspace-empty">
        <p>Open a request or create a new one to get started</p>
      </div>
      <ResponseViewer v-if="activeType === 'http'" />
    </div>
  </div>
</template>

<style scoped>
.workspace { display: flex; flex-direction: column; height: 100%; }
.workspace-body { flex: 1; display: flex; flex-direction: column; overflow: hidden; }
.workspace-empty { flex: 1; display: flex; align-items: center; justify-content: center; color: #999; font-size: 14px; }
</style>
