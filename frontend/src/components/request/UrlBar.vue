<template>
  <div class="url-bar">
    <n-select
      v-model:value="method"
      :options="methodOptions"
      :consistent-menu-width="false"
      size="medium"
      class="method-select"
    />
    <n-input
      v-model:value="url"
      placeholder="https://api.example.com/v1/users"
      size="medium"
      class="url-input"
      @keydown.enter="$emit('send')"
    />
    <n-button type="primary" class="send-btn" @click="$emit('send')">
      Send
    </n-button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { NSelect, NInput, NButton } from 'naive-ui'

const props = defineProps<{
  modelMethod: string
  modelUrl: string
}>()

const emit = defineEmits<{
  'update:modelMethod': [value: string]
  'update:modelUrl': [value: string]
  'send': []
}>()

const method = computed({
  get: () => props.modelMethod,
  set: (v) => emit('update:modelMethod', v),
})
const url = computed({
  get: () => props.modelUrl,
  set: (v) => emit('update:modelUrl', v),
})

const methodOptions = [
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'DELETE', value: 'DELETE' },
  { label: 'PATCH', value: 'PATCH' },
  { label: 'HEAD', value: 'HEAD' },
  { label: 'OPTIONS', value: 'OPTIONS' },
]
</script>

<style scoped>
.url-bar {
  display: flex;
  padding: 8px 10px;
  gap: 0;
  border-bottom: 1px solid #eee;
}
.method-select {
  width: 90px;
}
.method-select :deep(.n-base-selection) {
  border-radius: 6px 0 0 6px !important;
  border-right: none !important;
}
.url-input {
  flex: 1;
}
.url-input :deep(.n-input__border) {
  border-radius: 0 !important;
  border-left: none !important;
  border-right: none !important;
}
.send-btn {
  border-radius: 0 6px 6px 0 !important;
  font-weight: 600;
  padding: 0 22px;
}
</style>
