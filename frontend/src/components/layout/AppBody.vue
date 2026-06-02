<template>
  <div class="body">
    <Sidebar />
    <Workspace />
  </div>
</template>

<script setup lang="ts">
import Sidebar from './Sidebar.vue'
import Workspace from './Workspace.vue'
</script>

<style scoped>
.body {
  display: flex;
  flex: 1;
  overflow: hidden;
}
</style>
