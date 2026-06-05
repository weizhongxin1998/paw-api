<template>
  <div class="body-editor" @keydown="onEditorKeydown">
    <!-- ── Segmented body type selector ── -->
    <div class="body-type-bar">
      <button
        v-for="opt in bodyTypeOptions"
        :key="opt.value"
        :class="{ active: bodyType === opt.value }"
        class="seg-btn"
        @click="switchBodyType(opt.value)"
      >
        {{ opt.label }}
      </button>
    </div>

    <!-- ── None ── -->
    <div v-if="bodyType === 'none'" class="body-empty">
      <svg class="body-empty-icon" viewBox="0 0 24 24" width="28" height="28" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="12" cy="12" r="10" />
        <line x1="4.93" y1="4.93" x2="19.07" y2="19.07" />
      </svg>
      <span class="hint">{{ $t('body.noBody') }}</span>
    </div>

    <!-- ── form-data ── -->
    <KeyValueTable
      v-else-if="bodyType === 'form-data'"
      :items="formDataItems"
      :show-type="true"
      @update:items="onFormDataChange"
    />

    <!-- ── x-www-form-urlencoded ── -->
    <KeyValueTable
      v-else-if="bodyType === 'x-www-form-urlencoded'"
      :items="urlEncodedItems"
      @update:items="onUrlEncodedChange"
    />

    <!-- ── Raw editor ── -->
    <div v-else-if="bodyType === 'raw'" class="raw-editor">
      <div class="raw-header">
        <!-- Language indicator badge -->
        <span class="lang-badge" :class="`lang-${rawSubType}`">{{ rawSubTypeLabel }}</span>
        <!-- Sub-type selector -->
        <div class="raw-sub-bar">
          <button
            v-for="opt in rawSubOptions"
            :key="opt.value"
            :class="{ active: rawSubType === opt.value }"
            class="sub-btn"
            @click="rawSubType = opt.value"
          >
            {{ opt.label }}
          </button>
        </div>
        <div style="flex:1"></div>
        <!-- Beautify button (JSON only) -->
        <button
          v-if="rawSubType === 'json'"
          class="fmt-btn"
          :title="$t('body.raw.formatTitle')"
          @click="beautify"
        >
          <svg viewBox="0 0 24 24" width="13" height="13" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" />
            <polyline points="14 2 14 8 20 8" />
            <line x1="16" y1="13" x2="8" y2="13" />
            <line x1="16" y1="17" x2="8" y2="17" />
          </svg>
          <span>{{ $t('body.raw.format') }}</span>
        </button>
      </div>
      <n-input
        type="textarea"
        :value="rawContent"
        :rows="12"
        :placeholder="rawPlaceholder"
        class="raw-input"
        @update:value="onRawChange"
      />
      <div class="raw-footer-hint">
        <span v-if="rawSubType === 'json'" class="raw-hint">{{ $t('body.raw.formatHint') }}</span>
      </div>
    </div>

    <!-- ── Binary ── -->
    <div v-else-if="bodyType === 'binary'" class="binary-section">
      <div
        class="drop-zone"
        :class="{ 'drop-zone-active': isDragging }"
        @dragover.prevent="isDragging = true"
        @dragleave.prevent="isDragging = false"
        @drop.prevent="onDrop"
        @click="onSelectFile"
      >
        <svg class="drop-icon" viewBox="0 0 24 24" width="32" height="32" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
          <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" />
          <polyline points="17 8 12 3 7 8" />
          <line x1="12" y1="3" x2="12" y2="15" />
        </svg>
        <span v-if="binaryFileName" class="file-name">{{ binaryFileName }}</span>
        <span v-else class="drop-text">{{ $t('body.binary.dropText') }}</span>
        <span class="drop-hint">{{ $t('body.binary.dropHint') }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { NInput } from 'naive-ui'
import KeyValueTable from '../shared/KeyValueTable.vue'
import type { KvItem } from '../../types/request'

const { t } = useI18n()

const props = defineProps<{
  bodyType: string
  bodyData: string
}>()

const emit = defineEmits<{
  (e: 'update:bodyType', v: string): void
  (e: 'update:bodyData', v: string): void
}>()

const bodyTypeOptions = [
  { label: 'none', value: 'none' },
  { label: 'form-data', value: 'form-data' },
  { label: 'urlencoded', value: 'x-www-form-urlencoded' },
  { label: 'raw', value: 'raw' },
  { label: 'binary', value: 'binary' },
]

const rawSubOptions = [
  { label: 'JSON', value: 'json' },
  { label: 'XML', value: 'xml' },
  { label: 'HTML', value: 'html' },
  { label: 'Text', value: 'text' },
  { label: 'JS', value: 'javascript' },
]

const bodyType = ref(props.bodyType)
const formDataItems = ref<KvItem[]>([])
const urlEncodedItems = ref<KvItem[]>([])
const rawSubType = ref('json')
const rawContent = ref('')
const binaryFileName = ref('')
const isDragging = ref(false)

// ── Content preservation when switching body types ──
const bodyTypeSnapshots = new Map<string, { rawSubType: string; rawContent: string; binaryFileName: string }>()

function switchBodyType(newType: string) {
  // Save current state
  bodyTypeSnapshots.set(bodyType.value, {
    rawSubType: rawSubType.value,
    rawContent: rawContent.value,
    binaryFileName: binaryFileName.value,
  })
  bodyType.value = newType
  // Restore previous state if available
  const snap = bodyTypeSnapshots.get(newType)
  if (snap) {
    rawSubType.value = snap.rawSubType
    rawContent.value = snap.rawContent
    binaryFileName.value = snap.binaryFileName
  }
  emit('update:bodyType', newType)
}

const rawSubTypeLabel = computed(() => {
  const map: Record<string, string> = {
    json: 'JSON', xml: 'XML', html: 'HTML', text: 'TEXT', javascript: 'JS',
  }
  return map[rawSubType.value] || rawSubType.value.toUpperCase()
})

const rawPlaceholder = computed(() => {
  const hints: Record<string, string> = {
    json: '{\n  "key": "value"\n}',
    xml: '<root>\n  <item>value</item>\n</root>',
    html: '<html>\n  <body>...</body>\n</html>',
    text: t('body.raw.placeholderText'),
    javascript: '// JavaScript code\nconst data = {}',
  }
  return hints[rawSubType.value] || t('body.raw.placeholderDefault')
})

watch(() => props.bodyType, (v) => { bodyType.value = v })
watch(bodyType, (v) => emit('update:bodyType', v))

function onFormDataChange(items: KvItem[]) {
  formDataItems.value = items
  syncBody('form-data', JSON.stringify(items))
}

function onUrlEncodedChange(items: KvItem[]) {
  urlEncodedItems.value = items
  syncBody('x-www-form-urlencoded', JSON.stringify(items))
}

function onRawChange(v: string) {
  rawContent.value = v
  syncBody('raw', JSON.stringify({ subType: rawSubType.value, content: v }))
}

function syncBody(type: string, data: string) {
  emit('update:bodyData', data)
}

function beautify() {
  try {
    const obj = JSON.parse(rawContent.value)
    rawContent.value = JSON.stringify(obj, null, 2)
    syncBody('raw', JSON.stringify({ subType: rawSubType.value, content: rawContent.value }))
  } catch { /* ignore parse error */ }
}

// ── Ctrl+Shift+F to format JSON ──
function onEditorKeydown(e: KeyboardEvent) {
  if (e.ctrlKey && e.shiftKey && e.key === 'F') {
    if (bodyType.value === 'raw' && rawSubType.value === 'json') {
      e.preventDefault()
      beautify()
    }
  }
}

function onSelectFile() {
  binaryFileName.value = t('body.binary.notImplemented')
}

function onDrop(e: DragEvent) {
  isDragging.value = false
  const files = e.dataTransfer?.files
  if (files && files.length > 0) {
    binaryFileName.value = files[0].name
  }
}
</script>

<style scoped>
.body-editor {
  padding: 8px;
  outline: none;
}

/* ══════════════════════════════════════════
   Segmented body type selector
   ══════════════════════════════════════════ */
.body-type-bar {
  display: flex;
  gap: 0;
  background: var(--bg-elevated);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-sm);
  padding: 2px;
  margin-bottom: 10px;
  width: fit-content;
}
.seg-btn {
  padding: 4px 12px;
  font-size: var(--fs-xs);
  font-family: var(--font-mono);
  font-weight: 500;
  color: var(--text-secondary);
  background: transparent;
  border: none;
  border-radius: var(--radius-xs);
  cursor: pointer;
  transition: all var(--transition-fast);
  white-space: nowrap;
}
.seg-btn:hover:not(.active) {
  color: var(--text-secondary);
  background: var(--bg-hover);
}
.seg-btn.active {
  color: var(--text-primary);
  background: var(--bg-surface);
  box-shadow: 0 1px 3px rgba(0,0,0,0.15), 0 0 0 1px var(--border-primary);
  font-weight: 600;
}

/* ══════════════════════════════════════════
   Empty state
   ══════════════════════════════════════════ */
.body-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 32px 16px;
}
.body-empty-icon {
  color: var(--text-muted);
  opacity: 0.25;
}
.hint {
  color: var(--text-secondary);
  font-size: var(--fs-sm);
  font-family: var(--font-mono);
}

/* ══════════════════════════════════════════
   Raw editor
   ══════════════════════════════════════════ */
.raw-editor {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.raw-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

/* Language badge */
.lang-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
  font-size: var(--fs-2xs);
  font-family: var(--font-mono);
  font-weight: 700;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  border-radius: var(--radius-xs);
  border: 1px solid var(--border-primary);
  background: var(--bg-elevated);
  color: var(--text-muted);
  white-space: nowrap;
}
.lang-badge.lang-json   { color: #f59e0b; border-color: rgba(245,158,11,0.3); background: rgba(245,158,11,0.08); }
.lang-badge.lang-xml    { color: #a855f7; border-color: rgba(168,85,247,0.3); background: rgba(168,85,247,0.08); }
.lang-badge.lang-html   { color: #ef4444; border-color: rgba(239,68,68,0.3);  background: rgba(239,68,68,0.08); }
.lang-badge.lang-text   { color: var(--text-secondary); }
.lang-badge.lang-javascript { color: #f59e0b; border-color: rgba(245,158,11,0.3); background: rgba(245,158,11,0.08); }

/* Sub-type selector bar */
.raw-sub-bar {
  display: flex;
  gap: 0;
  background: var(--bg-elevated);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-xs);
  padding: 1px;
}
.sub-btn {
  padding: 2px 7px;
  font-size: var(--fs-2xs);
  font-family: var(--font-mono);
  font-weight: 500;
  color: var(--text-secondary);
  background: transparent;
  border: none;
  border-radius: 2px;
  cursor: pointer;
  transition: all var(--transition-fast);
}
.sub-btn:hover:not(.active) {
  color: var(--text-secondary);
}
.sub-btn.active {
  color: var(--text-primary);
  background: var(--bg-surface);
  font-weight: 600;
}

/* Format button */
.fmt-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px 9px;
  font-size: var(--fs-xs);
  font-family: var(--font-ui);
  font-weight: 500;
  color: var(--text-muted);
  background: transparent;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all var(--transition-fast);
  white-space: nowrap;
}
.fmt-btn:hover {
  color: var(--accent);
  border-color: var(--accent);
  background: var(--accent-soft);
}

/* Raw input */
.raw-input {
  font-family: var(--font-mono);
  font-size: var(--fs-sm);
}
.raw-input :deep(textarea) {
  line-height: 1.65;
}

/* Footer hint */
.raw-footer-hint {
  min-height: 16px;
}
.raw-hint {
  font-size: var(--fs-2xs);
  color: var(--text-muted);
  font-family: var(--font-ui);
  padding: 0 2px;
}

/* ══════════════════════════════════════════
   Binary / drop zone
   ══════════════════════════════════════════ */
.binary-section {
  padding: 4px;
}
.drop-zone {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 40px 24px;
  border: 2px dashed var(--border-primary);
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--transition);
  background: var(--bg-surface);
  text-align: center;
}
.drop-zone:hover {
  border-color: var(--border-hover);
  background: var(--bg-hover);
}
.drop-zone-active {
  border-color: var(--accent) !important;
  background: var(--accent-soft) !important;
}
.drop-zone-active .drop-icon {
  color: var(--accent);
}
.drop-icon {
  color: var(--text-muted);
  opacity: 0.35;
  transition: all var(--transition);
}
.drop-zone:hover .drop-icon {
  opacity: 0.6;
  color: var(--text-secondary);
}
.drop-text {
  font-size: var(--fs-sm);
  color: var(--text-secondary);
  font-family: var(--font-ui);
  font-weight: 500;
}
.drop-hint {
  font-size: var(--fs-xs);
  color: var(--text-muted);
  font-family: var(--font-ui);
}
.file-name {
  color: var(--blue);
  font-size: var(--fs-sm);
  font-family: var(--font-mono);
  font-weight: 600;
}
</style>
