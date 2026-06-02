<template>
  <div class="sidebar">
    <div class="sidebar-tabs">
      <n-button
        :type="activeTab === 'collections' ? 'primary' : 'default'"
        size="small"
        @click="activeTab = 'collections'"
      >
        Collections
      </n-button>
      <n-button
        :type="activeTab === 'history' ? 'primary' : 'default'"
        size="small"
        @click="activeTab = 'history'"
      >
        History
      </n-button>
    </div>
    <div class="sidebar-content">
      <n-empty v-if="activeTab === 'collections'" description="暂无集合" />
      <n-empty v-else description="暂无历史记录" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { NButton, NEmpty } from 'naive-ui'

const activeTab = ref('collections')
</script>

<style scoped>
.sidebar {
  width: 240px;
  background: #fafafa;
  border-right: 1px solid #e0e0e0;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}
.sidebar-tabs {
  display: flex;
  padding: 8px;
  gap: 4px;
  border-bottom: 1px solid #e8e8e8;
}
.sidebar-tabs .n-button {
  flex: 1;
}
.sidebar-content {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
}
</style>
