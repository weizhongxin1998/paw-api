<script lang="ts" setup>
import TabBar from '../components/TabBar.vue'
import RequestEditor from '../components/RequestEditor.vue'
import ResponseViewer from '../components/ResponseViewer.vue'
</script>

<template>
  <div class="workspace">
    <TabBar />
    <div class="workspace-body">
      <RequestEditor />
      <ResponseViewer />
    </div>
  </div>
</template>

<style scoped>
.workspace {
  display: flex;
  flex-direction: column;
  height: 100%;
}
.workspace-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
</style>
