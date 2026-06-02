<template>
  <div class="workspace">
    <n-empty v-if="!activeTab" description="点击左侧集合中的请求开始调试" />
    <div v-else class="workspace-editor">
      <div class="tabs-bar">TODO: TabBar</div>
      <UrlBar
        :model-method="'GET'"
        :model-url="''"
        @update:model-method="() => {}"
        @update:model-url="() => {}"
        @send="() => {}"
      />
      <div class="editor-body">
        <n-empty description="请求编辑区" size="small" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { NEmpty } from 'naive-ui'
import UrlBar from '../request/UrlBar.vue'

const activeTab = ref(null)
</script>

<style scoped>
.workspace {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fff;
  min-width: 0;
}
.workspace-editor {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}
.tabs-bar {
  padding: 6px 10px;
  border-bottom: 1px solid #eee;
  color: #aaa;
  font-size: 12px;
}
.editor-body {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
