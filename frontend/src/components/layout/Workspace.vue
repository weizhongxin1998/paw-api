<template>
  <div class="workspace">
    <n-empty description="点击左侧集合中的请求开始调试" />
  </div>
</template>

<script setup lang="ts">
import { NEmpty } from 'naive-ui'
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
</style>
