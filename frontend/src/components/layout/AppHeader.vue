<template>
  <div class="header">
    <div class="header-left">
      <button class="project-btn">
        {{ projectName || '未选择项目' }}
        <span class="proj-arrow">&#x25BC;</span>
      </button>
      <EnvSelector
        :project-id="projectId"
        @update:active-env-id="onActiveEnvChange"
      />
    </div>
    <div class="header-right">
      <button class="header-btn">日间</button>
      <button class="header-btn">设置</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useEnvStore } from '../../stores/env'
import EnvSelector from '../environment/EnvSelector.vue'

const projectName = ref('')
const projectId = ref<number | null>(null)
const envStore = useEnvStore()

function onActiveEnvChange(id: number | null) {
  envStore.activeEnvId = id
}
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
.project-btn {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 4px 10px;
  border: 1px solid #ddd;
  border-radius: 6px;
  background: #fff;
  cursor: pointer;
  font-size: 12px;
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
.header-btn {
  font-size: 11px;
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
</style>
