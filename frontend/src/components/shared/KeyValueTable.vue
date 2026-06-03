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
            <n-checkbox v-model:checked="item.enabled" size="small" />
          </td>
          <td><n-input v-model:value="item.key" size="small" borderless placeholder="Key" /></td>
          <td><n-input v-model:value="item.value" size="small" borderless placeholder="Value" /></td>
          <td><n-input v-model:value="item.description" size="small" borderless placeholder="Desc" /></td>
          <td v-if="showType">
            <n-select v-model:value="(item as any).fieldType" size="small" :options="typeOptions" />
          </td>
          <td class="td-remove" @click="onRemove(index)">&times;</td>
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
      <n-button v-if="showBulkEdit" text size="tiny" @click="isBulkEdit = !isBulkEdit">
        {{ isBulkEdit ? 'Table' : 'Bulk' }}
      </n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
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

// Watch the last row: when it gains content, auto-append a new empty row below it.
// Uses a stable boolean so typing multiple chars in the same cell doesn't re-trigger.
watch(
  () => {
    const arr = props.items
    if (arr.length === 0) return 'empty' as const
    const last = arr[arr.length - 1]
    return !!(last.key || last.value || last.description)
  },
  (filled, prev) => {
    if (prev === undefined) {
      // immediate: ensure exactly one trailing empty row
      if (props.items.length === 0) {
        emit('update:items', [{ id: String(++idCounter), key: '', value: '', description: '', enabled: true }])
      } else if (filled) {
        emit('update:items', [...props.items, { id: String(++idCounter), key: '', value: '', description: '', enabled: true }])
      }
      return
    }
    // subsequent: last row just became non-empty → add a new empty row below
    if (filled && !prev) {
      emit('update:items', [...props.items, { id: String(++idCounter), key: '', value: '', description: '', enabled: true }])
    }
  },
  { immediate: true }
)

function onRemove(index: number) {
  const newItems = props.items.filter((_, i) => i !== index)
  if (newItems.length === 0 || (newItems[newItems.length - 1].key || newItems[newItems.length - 1].value || newItems[newItems.length - 1].description)) {
    newItems.push({ id: String(++idCounter), key: '', value: '', description: '', enabled: true })
  }
  emit('update:items', newItems)
}

const bulkText = computed(() =>
  props.items.map((i) => `${i.key}: ${i.value}`).join('\n')
)

function onBulkChange(text: string) {
  const lines = text.split('\n').filter(Boolean)
  const newItems: KvItem[] = lines.map((line) => {
    const idx = line.indexOf(':')
    const key = idx !== -1 ? line.substring(0, idx).trim() : line
    const value = idx !== -1 ? line.substring(idx + 1).trim() : ''
    return { id: String(++idCounter), key, value, description: '', enabled: true }
  })
  emit('update:items', newItems)
}
</script>

<style scoped>
.key-value-table {
  padding: 2px;
}
table {
  width: 100%;
  border-collapse: collapse;
}
th {
  text-align: left;
  padding: 5px 6px;
  font-size: var(--fs-xs);
  color: var(--text-muted);
  text-transform: uppercase;
  border-bottom: 1px solid var(--border-primary);
  font-weight: 500;
  background: var(--bg-surface);
  letter-spacing: 0.3px;
}
td {
  padding: 2px 6px;
}
.td-check { text-align: center; width: 26px; }
.td-remove {
  text-align: center;
  color: var(--text-muted);
  cursor: pointer;
  font-size: var(--fs-md);
  font-weight: 600;
  transition: color var(--transition);
}
.td-remove:hover { color: var(--red); }
.kv-footer {
  display: flex;
  gap: 10px;
  padding: 3px 6px;
}
</style>
