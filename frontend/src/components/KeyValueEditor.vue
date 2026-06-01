<script lang="ts" setup>
import { NButton, NInput, NIcon, NDropdown } from 'naive-ui'
import { Add, Trash } from '@vicons/ionicons5'
import type { KeyValuePair } from '../stores/tabs'

export type { KeyValuePair }

const props = withDefaults(defineProps<{
  modelValue: KeyValuePair[]
  keyPlaceholder?: string
  valuePlaceholder?: string
  showPresets?: boolean
}>(), {
  keyPlaceholder: 'Key',
  valuePlaceholder: 'Value',
  showPresets: false,
})

const emit = defineEmits<{
  (e: 'update:modelValue', val: KeyValuePair[]): void
}>()

const headerPresets = [
  { label: 'Authorization: Bearer <token>', key: 'header-preset-bearer' },
  { label: 'Content-Type: application/json', key: 'header-preset-json' },
  { label: 'Content-Type: application/x-www-form-urlencoded', key: 'header-preset-form' },
  { label: 'Accept: application/json', key: 'header-preset-accept' },
  { label: 'User-Agent: PawAPI/1.0', key: 'header-preset-ua' },
  { label: 'Cache-Control: no-cache', key: 'header-preset-cache' },
  { label: 'Referer:', key: 'header-preset-referer' },
]

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

function addPreset(key: string) {
  const map: Record<string, { key: string; value: string }> = {
    'header-preset-bearer': { key: 'Authorization', value: 'Bearer <token>' },
    'header-preset-json': { key: 'Content-Type', value: 'application/json' },
    'header-preset-form': { key: 'Content-Type', value: 'application/x-www-form-urlencoded' },
    'header-preset-accept': { key: 'Accept', value: 'application/json' },
    'header-preset-ua': { key: 'User-Agent', value: 'PawAPI/1.0' },
    'header-preset-cache': { key: 'Cache-Control', value: 'no-cache' },
    'header-preset-referer': { key: 'Referer', value: '' },
  }
  const preset = map[key]
  if (!preset) return
  const exists = props.modelValue.some(p => p.key === preset.key)
  if (exists) return
  emit('update:modelValue', [...props.modelValue, { key: preset.key, value: preset.value, enabled: true }])
}
</script>

<template>
  <div class="kv-editor">
    <div class="kv-toolbar">
      <NButton size="tiny" quaternary @click="addRow" class="kv-add-btn">
        <template #icon><NIcon><Add /></NIcon></template>
        Add
      </NButton>
      <NDropdown v-if="showPresets" :options="headerPresets" @select="addPreset">
        <NButton size="tiny" quaternary>{{ $t('request.presets') }}</NButton>
      </NDropdown>
    </div>
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
