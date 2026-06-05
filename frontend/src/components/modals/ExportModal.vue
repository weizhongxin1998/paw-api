<template>
  <n-modal
    :show="show"
    preset="card"
    :title="$t('export.title')"
    :class="modalClass"
    style="width: 540px"
    :mask-closable="false"
    @update:show="emit('update:show', $event)"
  >
    <!-- Result screen -->
    <template v-if="status === 'result'">
      <n-result :status="resultType" :title="resultTitle" :description="resultDesc">
        <template #footer>
          <n-button @click="onReset">{{ $t('export.continueExport') }}</n-button>
          <n-button type="primary" @click="onClose">{{ $t('common.done') }}</n-button>
        </template>
      </n-result>
    </template>

    <!-- Export form -->
    <template v-else>
      <!-- Scope selector: visual cards -->
      <div class="field-section">
        <label class="field-label">{{ $t('export.scopeLabel') }}</label>
        <div class="scope-cards">
          <div
            class="scope-card"
            :class="{ selected: scope === 'project' }"
            @click="scope = 'project'"
          >
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M22 19a2 2 0 01-2 2H4a2 2 0 01-2-2V5a2 2 0 012-2h5l2 3h9a2 2 0 012 2z"/></svg>
            <div class="scope-info">
              <span class="scope-name">{{ $t('export.scopeProject') }}</span>
              <span class="scope-desc">{{ $t('export.scopeProjectDesc') }}</span>
            </div>
          </div>
          <div
            class="scope-card disabled"
            :class="{ selected: scope === 'collection' }"
          >
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
            <div class="scope-info">
              <span class="scope-name">{{ $t('export.scopeCollection') }}</span>
              <span class="scope-desc">{{ $t('common.comingSoon') }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Format selector with icon badge -->
      <div class="field-section">
        <label class="field-label">{{ $t('export.formatLabel') }}</label>
        <div class="format-options">
          <div
            class="format-option"
            :class="{ selected: format === 'postman' }"
            @click="format = 'postman'"
          >
            <span class="format-icon postman-icon">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><circle cx="12" cy="12" r="10"/><path d="M8 12h8"/><path d="M12 8v8"/></svg>
            </span>
            <span class="format-name">Postman Collection</span>
            <span class="format-ext">.json</span>
          </div>
          <div class="format-option disabled">
            <span class="format-icon">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
            </span>
            <span class="format-name">{{ $t('export.formatPaw') }}</span>
            <span class="format-badge">{{ $t('common.comingSoon') }}</span>
          </div>
        </div>
      </div>

      <!-- Summary preview -->
      <div class="export-summary" v-if="stats">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/></svg>
        <span v-html="summaryHtml"></span>
      </div>

      <!-- File path with folder icon browse button -->
      <div class="field-section">
        <label class="field-label">{{ $t('export.savePathLabel') }}</label>
        <n-input-group>
          <n-input
            v-model:value="filePath"
            :placeholder="$t('export.savePathPlaceholder')"
            :disabled="status === 'exporting'"
          >
            <template #prefix>
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M22 19a2 2 0 01-2 2H4a2 2 0 01-2-2V5a2 2 0 012-2h5l2 3h9a2 2 0 012 2z"/></svg>
            </template>
          </n-input>
          <n-button :disabled="status === 'exporting'" @click="onBrowse" :title="$t('common.browseFolder')">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M22 19a2 2 0 01-2 2H4a2 2 0 01-2-2V5a2 2 0 012-2h5l2 3h9a2 2 0 012 2z"/></svg>
          </n-button>
        </n-input-group>
      </div>

      <!-- Progress indicator -->
      <div class="export-progress" v-if="status === 'exporting'">
        <n-progress type="line" :percentage="exportProgress" :show-indicator="false" processing />
        <span class="progress-text">{{ $t('export.exporting') }}</span>
      </div>

      <!-- Error message with suggestion -->
      <n-alert v-if="exportError" type="error" :bordered="false" style="margin-top: 12px">
        <div class="error-content">
          <span class="error-msg">{{ exportError }}</span>
          <span class="error-hint">{{ $t('export.errorHint') }}</span>
        </div>
      </n-alert>

      <!-- Footer -->
      <div class="modal-footer">
        <n-button @click="onClose" :disabled="status === 'exporting'">{{ $t('common.cancel') }}</n-button>
        <n-button
          type="primary"
          :loading="status === 'exporting'"
          :disabled="!canExport"
          @click="onExport"
        >
          <template #icon>
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
          </template>
          {{ $t('export.exportBtn') }}
        </n-button>
      </div>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import {
  NModal, NForm, NFormItem, NSelect, NInput, NInputGroup,
  NButton, NResult, NProgress, NAlert, useMessage,
} from 'naive-ui'
import { ExportPostman, GetProjectStats } from '../../../wailsjs/go/main/App'
import { useI18n } from 'vue-i18n'
import { useThemeClass } from '../../composables/useThemeClass'

interface Props { show: boolean; projectId: number | null }
const props = defineProps<Props>()
const emit = defineEmits<{ 'update:show': [value: boolean]; exported: [] }>()
const message = useMessage()
const { t } = useI18n()

type Status = 'idle' | 'exporting' | 'result'
const status = ref<Status>('idle')
const scope = ref('project')
const format = ref('postman')
const filePath = ref('')
const resultType = ref<'success' | 'error'>('success')
const resultTitle = ref('')
const resultDesc = ref('')
const exportProgress = ref(0)
const exportError = ref('')
const stats = ref<{ request_count: number; collection_count: number } | null>(null)

const { modalClass } = useThemeClass('export-modal')

const canExport = computed(() => props.projectId !== null)

const summaryHtml = computed(() => {
  if (!stats.value) return ''
  return t('export.summary', { requests: `<strong>${stats.value.request_count}</strong>`, collections: `<strong>${stats.value.collection_count}</strong>` })
})

// Fetch project stats when modal opens
watch(() => props.show, async (visible) => {
  if (visible && props.projectId) {
    try {
      stats.value = await GetProjectStats(props.projectId)
    } catch { stats.value = null }
  } else if (!visible) {
    stats.value = null
  }
})

function onBrowse() {
  // Wails does not expose a native SaveFile dialog in this runtime.
  // The parent component can listen for browse events or we show a hint.
  message.info(t('export.browseHint'))
}

function onReset() {
  status.value = 'idle'
  filePath.value = ''
  resultTitle.value = ''
  resultDesc.value = ''
  exportError.value = ''
  exportProgress.value = 0
}

function onClose() {
  emit('update:show', false)
  status.value = 'idle'
  filePath.value = ''
  exportError.value = ''
  exportProgress.value = 0
}

async function onExport() {
  if (!props.projectId) return
  status.value = 'exporting'
  exportError.value = ''
  exportProgress.value = 10

  // Simulate progress during export
  const progressInterval = setInterval(() => {
    if (exportProgress.value < 90) {
      exportProgress.value += Math.random() * 15
    }
  }, 300)

  try {
    const result = await ExportPostman(props.projectId)
    clearInterval(progressInterval)
    exportProgress.value = 100
    resultType.value = 'success'
    resultTitle.value = t('export.successTitle')
    resultDesc.value = result ? t('export.successDescSaved') + result : t('export.successDescDone')
    status.value = 'result'
    emit('exported')
  } catch (err: any) {
    clearInterval(progressInterval)
    exportProgress.value = 0
    const msg = err?.message || err?.toString() || t('common.unknownError')
    exportError.value = msg
    resultType.value = 'error'
    resultTitle.value = t('export.failTitle')
    resultDesc.value = msg
    status.value = 'result'
  }
}
</script>

<style scoped>
.field-section { margin-bottom: 16px; }
.field-label {
  display: block;
  font-size: var(--fs-sm);
  font-weight: 600;
  color: var(--text-secondary);
  margin-bottom: 8px;
}

/* -- Scope Cards -- */
.scope-cards { display: flex; gap: 10px; }
.scope-card {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  border: 2px solid var(--border-primary);
  border-radius: var(--radius-lg, 10px);
  cursor: pointer;
  transition: all 0.15s ease;
  background: var(--bg-surface, transparent);
}
.scope-card:hover:not(.disabled) { border-color: var(--text-muted); }
.scope-card.selected {
  border-color: var(--accent);
  background: var(--accent-soft, rgba(0,224,90,0.04));
}
.scope-card.disabled {
  opacity: 0.45;
  cursor: not-allowed;
}
.scope-card svg {
  flex-shrink: 0;
  opacity: 0.5;
}
.scope-card.selected svg { opacity: 0.8; color: var(--accent); }
.scope-info { display: flex; flex-direction: column; }
.scope-name { font-size: var(--fs-sm); font-weight: 600; color: var(--text-primary); }
.scope-desc { font-size: var(--fs-xs); color: var(--text-muted); margin-top: 1px; }

/* -- Format Options -- */
.format-options { display: flex; flex-direction: column; gap: 6px; }
.format-option {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  cursor: pointer;
  transition: all 0.15s ease;
}
.format-option:hover:not(.disabled) { background: var(--bg-hover); }
.format-option.selected {
  border-color: var(--accent);
  background: var(--accent-soft, rgba(0,224,90,0.04));
}
.format-option.disabled { opacity: 0.4; cursor: not-allowed; }
.format-icon {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  background: var(--bg-hover, rgba(255,255,255,0.04));
  flex-shrink: 0;
}
.postman-icon { color: #f97316; }
.format-name { font-size: var(--fs-sm); flex: 1; }
.format-ext {
  font-size: var(--fs-xs);
  font-family: var(--font-mono);
  color: var(--text-muted);
  background: var(--bg-hover);
  padding: 1px 6px;
  border-radius: 4px;
}
.format-badge {
  font-size: var(--fs-2xs);
  color: var(--text-muted);
  background: var(--bg-hover);
  padding: 1px 6px;
  border-radius: 4px;
}

/* -- Export Summary -- */
.export-summary {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  background: var(--accent-soft, rgba(0,224,90,0.04));
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  font-size: var(--fs-sm);
  color: var(--text-secondary);
  margin-bottom: 16px;
}
.export-summary svg { opacity: 0.5; flex-shrink: 0; }
.export-summary strong { color: var(--accent); font-weight: 700; }

/* -- Progress -- */
.export-progress {
  margin-top: 12px;
  display: flex;
  align-items: center;
  gap: 10px;
}
.export-progress :deep(.n-progress) { flex: 1; }
.progress-text { font-size: var(--fs-xs); color: var(--text-muted); white-space: nowrap; }

/* -- Error -- */
.error-content { display: flex; flex-direction: column; gap: 2px; }
.error-msg { font-size: var(--fs-sm); }
.error-hint { font-size: var(--fs-xs); opacity: 0.7; }

/* -- Footer -- */
.modal-footer { display: flex; justify-content: flex-end; gap: 10px; margin-top: 20px; }
</style>
