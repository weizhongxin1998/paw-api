<template>
  <div class="key-value-table">
    <table v-if="!isBulkEdit">
      <thead>
        <tr>
          <th style="width:24px"></th>
          <th style="width:25%">Key</th>
          <th style="width:35%">Value</th>
          <th style="width:25%">Description</th>
          <th v-if="showType" style="width:40px">Type</th>
          <th style="width:20px"></th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(item, index) in items" :key="item.id">
          <td class="td-check">
            <n-checkbox v-model:checked="item.enabled" size="medium" />
          </td>
          <td><n-input v-model:value="item.key" size="small" borderless placeholder="Key" /></td>
          <td><n-input v-model:value="item.value" size="small" borderless placeholder="Value" /></td>
          <td><n-input v-model:value="item.description" size="small" borderless placeholder="Desc" /></td>
          <td v-if="showType">
            <n-select v-model:value="(item as any).fieldType" size="small" :options="typeOptions" />
          </td>
          <td class="td-remove" @click="onRemove(index)">x</td>
        </tr>
      </tbody>
    </table>
    <n-input
      v-else
      type="textarea"
      :value="bulkText"
      :rows="6"
      placeholder="Key: Value"
      @update:value="onBulkChange"
    />
    <div class="kv-footer">
      <n-button text size="tiny" @click="onAdd">+ 添加</n-button>
      <n-button v-if="showBulkEdit" text size="tiny" @click="isBulkEdit = !isBulkEdit">
        {{ isBulkEdit ? 'Table' : 'Bulk Edit' }}
      </n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { NInput, NCheckbox, NButton, NSelect } from 'naive-ui'
import type { KvItem } from '../../types/request'

const props = defineProps<{
  items: KvItem[]
  showType?: boolean
  showBulkEdit?: boolean
  bulkMode?: boolean
}>()

const emit = defineEmits<{
  'update:items': [items: KvItem[]]
}>()

const isBulkEditInternal = ref(false)
const isBulkEdit = computed({
  get: () => props.bulkMode ?? isBulkEditInternal.value,
  set: (v) => { isBulkEditInternal.value = v },
})
const typeOptions = [
  { label: 'text', value: 'text' },
  { label: 'file', value: 'file' },
]

let idCounter = Date.now()

function onAdd() {
  const item: KvItem = {
    id: String(++idCounter),
    key: '',
    value: '',
    description: '',
    enabled: true,
  }
  const newItems = [...props.items, item]
  emit('update:items', newItems)
}

function onRemove(index: number) {
  const newItems = props.items.filter((_, i) => i !== index)
  emit('update:items', newItems)
}

const bulkText = computed(() =>
  props.items.map((i) => `${i.key}: ${i.value}`).join('\n')
)

function onBulkChange(text: string) {
  const lines = text.split('\n').filter(Boolean)
  const newItems: KvItem[] = lines.map((line, i) => {
    const idx = line.indexOf(':')
    const key = idx !== -1 ? line.substring(0, idx).trim() : line
    const value = idx !== -1 ? line.substring(idx + 1).trim() : ''
    return {
      id: String(++idCounter),
      key,
      value,
      description: '',
      enabled: true,
    }
  })
  emit('update:items', newItems)
}
</script>

<style scoped>
.key-value-table {
  padding: 4px;
}
table {
  width: 100%;
  border-collapse: collapse;
}
th {
  text-align: left;
  padding: 6px 8px;
  font-size: 12px;
  color: #999;
  text-transform: uppercase;
  border-bottom: 1px solid #eee;
  font-weight: 500;
  background: #fafafa;
}
td {
  padding: 3px 8px;
}
.td-check {
  text-align: center;
  width: 28px;
}
.td-remove {
  text-align: center;
  color: #ccc;
  cursor: pointer;
  font-size: 15px;
}
.td-remove:hover {
  color: #d03050;
}
.kv-footer {
  display: flex;
  gap: 12px;
  padding: 4px 8px;
}
</style>
