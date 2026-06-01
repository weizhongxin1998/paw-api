<script lang="ts" setup>
import { NInput } from 'naive-ui'

const props = defineProps<{
  modelValue: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', val: string): void
}>()
</script>

<template>
  <NInput
    :value="props.modelValue"
    type="textarea"
    :rows="8"
    placeholder="{}"
    class="json-editor"
    @update:value="emit('update:modelValue', $event)"
  />
</template>

<style scoped>
.json-editor {
  font-family: 'Cascadia Code', 'Fira Code', 'Consolas', monospace;
  font-size: 13px;
}
</style>
