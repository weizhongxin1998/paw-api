<script lang="ts" setup>
import { NButton, NInput, NIcon } from 'naive-ui'
import { Add, Trash } from '@vicons/ionicons5'
import type { KeyValuePair } from '../stores/tabs'

export type { KeyValuePair }

const props = withDefaults(defineProps<{
  modelValue: KeyValuePair[]
  keyPlaceholder?: string
  valuePlaceholder?: string
}>(), {
  keyPlaceholder: 'Key',
  valuePlaceholder: 'Value',
})

const emit = defineEmits<{
  (e: 'update:modelValue', val: KeyValuePair[]): void
}>()

function addRow() {
  emit('update:modelValue', [...props.modelValue, { key: '', value: '', enabled: true }])
}

function removeRow(index: number) {
  const next = props.modelValue.filter((_, i) => i !== index)
  emit('update:modelValue', next)
}

function updateRow(index: number, field: 'key' | 'value' | 'enabled', val: string | boolean) {
  const next = props.modelValue.map((pair, i) => {
    if (i !== index) return pair
    return { ...pair, [field]: val }
  })
  emit('update:modelValue', next)
}
</script>

<template>
  <div class="kv-editor">
    <div class="kv-rows">
      <div v-for="(pair, index) in modelValue" :key="index" class="kv-row">
        <input
          type="checkbox"
          :checked="pair.enabled"
          class="kv-checkbox"
          @change="updateRow(index, 'enabled', ($event.target as HTMLInputElement).checked)"
        />
        <NInput
          :value="pair.key"
          size="tiny"
          :placeholder="keyPlaceholder"
          class="kv-input"
          @update:value="updateRow(index, 'key', $event)"
        />
        <NInput
          :value="pair.value"
          size="tiny"
          :placeholder="valuePlaceholder"
          class="kv-input"
          @update:value="updateRow(index, 'value', $event)"
        />
        <NButton quaternary circle size="tiny" @click="removeRow(index)">
          <template #icon><NIcon><Trash /></NIcon></template>
        </NButton>
      </div>
    </div>
    <NButton size="tiny" quaternary @click="addRow" class="kv-add-btn">
      <template #icon><NIcon><Add /></NIcon></template>
      Add
    </NButton>
  </div>
</template>

<style scoped>
.kv-rows {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.kv-row {
  display: flex;
  align-items: center;
  gap: 4px;
}
.kv-checkbox {
  width: 14px;
  height: 14px;
  cursor: pointer;
  flex-shrink: 0;
}
.kv-input {
  flex: 1;
}
.kv-add-btn {
  margin-top: 4px;
}
</style>
