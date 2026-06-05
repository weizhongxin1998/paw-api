<template>
  <div class="key-value-table">
    <!-- ── Table mode ── -->
    <div v-if="!isBulkEdit" class="kv-table-wrap">
      <table>
        <thead>
          <tr>
            <th class="th-check">
              <n-checkbox
                v-if="headerCheck !== undefined"
                :checked="headerCheck"
                size="small"
                @update:checked="(v: boolean) => emit('update:headerCheck', v)"
              />
            </th>
            <th style="width:25%">键</th>
            <th style="width:35%">值</th>
            <th style="width:25%">描述</th>
            <th v-if="showType" style="width:40px">类型</th>
            <th class="th-actions"></th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="(item, index) in items"
            :key="item.id"
            :class="{ 'row-alt': index % 2 === 1, 'row-filled': item.key || item.value }"
          >
            <td class="td-check">
              <n-checkbox v-model:checked="item.enabled" size="small" />
            </td>
            <td>
              <n-input
                :ref="(el: any) => setKeyRef(el, index)"
                v-model:value="item.key"
                size="small"
                borderless
                placeholder="键"
                @keydown="onCellKeydown($event, index, 'key')"
              />
            </td>
            <td>
              <n-input
                :ref="(el: any) => setValRef(el, index)"
                v-model:value="item.value"
                size="small"
                borderless
                placeholder="值"
                @keydown="onCellKeydown($event, index, 'value')"
              />
            </td>
            <td>
              <n-input
                :ref="(el: any) => setDescRef(el, index)"
                v-model:value="item.description"
                size="small"
                borderless
                placeholder="描述"
                @keydown="onCellKeydown($event, index, 'desc')"
              />
            </td>
            <td v-if="showType">
              <n-select v-model:value="(item as any).fieldType" size="small" :options="typeOptions" />
            </td>
            <td class="td-actions">
              <!-- Duplicate button -->
              <button class="act-btn dup-btn" title="复制行 (Ctrl+D)" @click="onDuplicate(index)">
                <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <rect x="9" y="9" width="13" height="13" rx="2" />
                  <path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1" />
                </svg>
              </button>
              <!-- Remove button -->
              <button class="act-btn rm-btn" title="删除行" @click="onRemove(index)">
                <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <polyline points="3 6 5 6 21 6" />
                  <path d="M19 6l-1 14a2 2 0 01-2 2H8a2 2 0 01-2-2L5 6" />
                  <path d="M10 11v6" />
                  <path d="M14 11v6" />
                  <path d="M9 6V4a1 1 0 011-1h4a1 1 0 011 1v2" />
                </svg>
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- ── Bulk edit mode ── -->
    <div v-else class="kv-bulk-wrap">
      <div class="kv-bulk-editor">
        <div class="kv-bulk-lines">{{ bulkLineNumbers }}</div>
        <n-input
          type="textarea"
          :value="bulkText"
          :rows="8"
          placeholder="键: 值 (每行一个)"
          class="kv-bulk-input"
          @update:value="onBulkChange"
        />
      </div>
      <div class="kv-bulk-hint">每行格式: <code>键: 值</code>，支持多行批量编辑</div>
    </div>

    <!-- ── Footer ── -->
    <div class="kv-footer">
      <div style="flex:1"></div>
      <n-button v-if="showBulkEdit" text size="tiny" class="kv-bulk-toggle" @click="isBulkEdit = !isBulkEdit">
        <template #icon>
          <svg v-if="!isBulkEdit" viewBox="0 0 24 24" width="13" height="13" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" />
            <polyline points="14 2 14 8 20 8" />
            <line x1="16" y1="13" x2="8" y2="13" />
            <line x1="16" y1="17" x2="8" y2="17" />
            <polyline points="10 9 9 9 8 9" />
          </svg>
          <svg v-else viewBox="0 0 24 24" width="13" height="13" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <rect x="3" y="3" width="18" height="18" rx="2" />
            <line x1="3" y1="9" x2="21" y2="9" />
            <line x1="9" y1="21" x2="9" y2="9" />
          </svg>
        </template>
        {{ isBulkEdit ? '表格模式' : '批量编辑' }}
      </n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import { NInput, NCheckbox, NButton, NSelect } from 'naive-ui'
import type { KvItem } from '../../types/request'

const props = defineProps<{
  items: KvItem[]
  showType?: boolean
  showBulkEdit?: boolean
  bulkMode?: boolean
  headerCheck?: boolean
}>()

const emit = defineEmits<{
  'update:items': [items: KvItem[]]
  'update:headerCheck': [value: boolean]
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

// ── Input ref maps for keyboard navigation ──
const keyRefs = new Map<number, any>()
const valRefs = new Map<number, any>()
const descRefs = new Map<number, any>()

function setKeyRef(el: any, idx: number) { if (el) keyRefs.set(idx, el) }
function setValRef(el: any, idx: number) { if (el) valRefs.set(idx, el) }
function setDescRef(el: any, idx: number) { if (el) descRefs.set(idx, el) }

function focusInput(refMap: Map<number, any>, idx: number) {
  nextTick(() => {
    const el = refMap.get(idx)
    if (el) {
      const inner = el.$el?.querySelector('input') || el.inputElRef
      if (inner && inner.focus) inner.focus()
      else if (el.focus) el.focus()
    }
  })
}

function onCellKeydown(e: KeyboardEvent, index: number, field: 'key' | 'value' | 'desc') {
  // Ctrl+D → duplicate row
  if (e.ctrlKey && e.key === 'd') {
    e.preventDefault()
    onDuplicate(index)
    return
  }

  if (e.key === 'Enter' && !e.shiftKey && !e.ctrlKey) {
    e.preventDefault()
    // Move to the same column in the next row
    const nextIdx = index + 1
    if (nextIdx < props.items.length) {
      if (field === 'key') focusInput(keyRefs, nextIdx)
      else if (field === 'value') focusInput(valRefs, nextIdx)
      else focusInput(descRefs, nextIdx)
    } else {
      // Auto-add a new row and focus it
      onAddRow()
      nextTick(() => {
        const ni = props.items.length // after add, the new index
        if (field === 'key') focusInput(keyRefs, ni - 1)
        else if (field === 'value') focusInput(valRefs, ni - 1)
        else focusInput(descRefs, ni - 1)
      })
    }
    return
  }

  if (e.key === 'Tab' && !e.shiftKey) {
    // Tab moves to next cell: key→value→desc→(next row key)
    e.preventDefault()
    if (field === 'key') {
      focusInput(valRefs, index)
    } else if (field === 'value') {
      focusInput(descRefs, index)
    } else {
      // desc → next row key
      const nextIdx = index + 1
      if (nextIdx < props.items.length) {
        focusInput(keyRefs, nextIdx)
      } else {
        onAddRow()
        nextTick(() => focusInput(keyRefs, props.items.length - 1))
      }
    }
    return
  }
}

// ── Auto-append watcher ──
let _kvInit = true
watch(
  () => {
    const arr = props.items
    if (arr.length === 0) return 'empty' as const
    const last = arr[arr.length - 1]
    return !!(last.key || last.value || last.description)
  },
  (filled, prev) => {
    if (_kvInit) { _kvInit = false; return }
    if (prev === undefined) {
      if (props.items.length === 0) {
        emit('update:items', [{ id: String(++idCounter), key: '', value: '', description: '', enabled: true }])
      } else if (filled) {
        emit('update:items', [...props.items, { id: String(++idCounter), key: '', value: '', description: '', enabled: true }])
      }
      return
    }
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

function onDuplicate(index: number) {
  const src = props.items[index]
  if (!src) return
  const clone: KvItem = {
    id: String(++idCounter),
    key: src.key,
    value: src.value,
    description: src.description,
    enabled: src.enabled,
  }
  const newItems = [...props.items]
  newItems.splice(index + 1, 0, clone)
  // Ensure trailing empty row
  const last = newItems[newItems.length - 1]
  if (last.key || last.value || last.description) {
    newItems.push({ id: String(++idCounter), key: '', value: '', description: '', enabled: true })
  }
  emit('update:items', newItems)
  // Focus the cloned row's key field
  nextTick(() => focusInput(keyRefs, index + 1))
}

function onAddRow() {
  const newItems = [...props.items]
  // If the last row is empty, focus it instead of adding another
  const last = newItems[newItems.length - 1]
  if (last && !last.key && !last.value && !last.description) {
    focusInput(keyRefs, newItems.length - 1)
    return
  }
  newItems.push({ id: String(++idCounter), key: '', value: '', description: '', enabled: true })
  emit('update:items', newItems)
  nextTick(() => focusInput(keyRefs, newItems.length - 1))
}

const bulkText = computed(() =>
  props.items.map((i) => `${i.key}: ${i.value}`).join('\n')
)

const bulkLineNumbers = computed(() => {
  const lines = props.items.filter(i => i.key || i.value).length
  const count = Math.max(lines, 1)
  return Array.from({ length: count }, (_, i) => i + 1).join('\n')
})

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
  padding: 4px;
}

/* ── Table ── */
table {
  width: 100%;
  border-collapse: separate;
  border-spacing: 0;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  overflow: hidden;
}
th {
  text-align: left;
  padding: 7px 10px;
  font-size: var(--fs-xs);
  color: var(--text-secondary);
  text-transform: uppercase;
  border-bottom: 1px solid var(--border-primary);
  font-weight: 600;
  background: var(--bg-elevated);
  letter-spacing: 0.04em;
  font-family: var(--font-ui);
}
.th-check {
  width: 30px;
  text-align: center;
}
.th-actions {
  width: 48px;
}
td {
  padding: 3px 4px;
  transition: background var(--transition-fast);
}
tr:not(:last-child) td {
  border-bottom: 1px solid var(--border-subtle);
}

/* ── Alternating row backgrounds ── */
tr.row-alt td {
  background: var(--bg-inset);
}
tr.row-alt:hover td {
  background: var(--bg-hover) !important;
}
tr:not(.row-alt):hover td {
  background: var(--bg-hover);
}

/* ── Checkbox column ── */
.td-check {
  text-align: center;
  width: 30px;
  padding: 3px 2px;
}

/* ── Action buttons column ── */
.td-actions {
  text-align: center;
  width: 48px;
  padding: 2px 2px;
  white-space: nowrap;
}
.act-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 22px;
  height: 22px;
  border: none;
  background: transparent;
  color: var(--text-muted);
  cursor: pointer;
  border-radius: var(--radius-xs);
  transition: all var(--transition-fast);
  opacity: 0;
  padding: 0;
}
tr:hover .act-btn {
  opacity: 0.7;
}
.act-btn:hover {
  opacity: 1 !important;
}
.dup-btn:hover {
  color: var(--blue);
  background: var(--blue-soft);
}
.rm-btn:hover {
  color: var(--red);
  background: var(--red-soft);
}

/* ── Bulk edit ── */
.kv-bulk-wrap {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.kv-bulk-editor {
  display: flex;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  overflow: hidden;
  background: var(--bg-surface);
}
.kv-bulk-lines {
  padding: 8px 8px 8px 10px;
  font-family: var(--font-mono);
  font-size: var(--fs-xs);
  line-height: 1.7;
  color: var(--text-placeholder);
  background: var(--bg-elevated);
  border-right: 1px solid var(--border-primary);
  text-align: right;
  user-select: none;
  white-space: pre;
  min-width: 30px;
}
.kv-bulk-input {
  flex: 1;
  font-family: var(--font-mono);
  font-size: var(--fs-sm);
}
.kv-bulk-input :deep(textarea) {
  border: none !important;
  box-shadow: none !important;
  padding: 8px 10px;
  line-height: 1.7;
}
.kv-bulk-hint {
  font-size: var(--fs-xs);
  color: var(--text-muted);
  font-family: var(--font-ui);
  padding: 0 4px;
}
.kv-bulk-hint code {
  background: var(--bg-elevated);
  padding: 1px 4px;
  border-radius: var(--radius-xs);
  font-family: var(--font-mono);
  font-size: var(--fs-2xs);
  border: 1px solid var(--border-primary);
  color: var(--accent);
}

/* ── Footer ── */
.kv-footer {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 4px;
  justify-content: flex-end;
}
.kv-bulk-toggle {
  font-size: var(--fs-xs) !important;
}
</style>
