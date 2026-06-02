<template>
  <div class="header">
    <div class="header-left">
      <n-button text class="project-btn">
        {{ projectName || '未选择项目' }}
      </n-button>
      <n-select
        v-model:value="activeEnv"
        :options="envOptions"
        size="small"
        placeholder="环境"
        style="width: 140px"
      />
    </div>
    <div class="header-right">
      <n-button text size="small">日间</n-button>
      <n-button text size="small">设置</n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { NButton, NSelect } from 'naive-ui'

const projectName = ref('')
const activeEnv = ref<string | null>(null)
const envOptions = ref<{ label: string; value: string }[]>([])
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
  font-size: 13px;
  font-weight: 600;
}
</style>
