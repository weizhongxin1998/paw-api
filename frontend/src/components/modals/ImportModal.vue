<template>
  <n-modal
    :show="show"
    preset="card"
    title="导入"
    :class="modalClass"
    style="width: 540px"
    :mask-closable="false"
    @update:show="emit('update:show', $event)"
  >
    <!-- Result screen -->
    <template v-if="status === 'result'">
      <n-result :status="resultType" :title="resultTitle" :description="resultDesc">
        <template #footer>
          <n-button @click="onReset">继续导入</n-button>
          <n-button type="primary" @click="onClose">完成</n-button>
        </template>
      </n-result>
    </template>

    <!-- Import form -->
    <template v-else>
      <!-- Format selector with icon badges -->
      <div class="field-section">
        <label class="field-label">导入格式</label>
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
            <span class="format-name">OpenAPI 3.x</span>
            <span class="format-badge">即将推出</span>
          </div>
          <div class="format-option disabled">
            <span class="format-icon">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
            </span>
            <span class="format-name">Swagger 2.x</span>
            <span class="format-badge">即将推出</span>
          </div>
        </div>
      </div>

      <!-- Drag-and-drop file zone -->
      <div class="field-section">
        <label class="field-label">文件路径</label>
        <div
          class="drop-zone"
          :class="{ dragging: isDragging, 'has-file': !!filePath }"
          @dragover.prevent="onDragOver"
          @dragleave.prevent="onDragLeave"
          @drop.prevent="onDrop"
          @click="onBrowse"
        >
          <template v-if="filePath">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round">
              <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/>
              <polyline points="14 2 14 8 20 8"/>
            </svg>
            <div class="drop-file-info">
              <span class="drop-file-name">{{ fileName }}</span>
              <span class="drop-file-path" :title="filePath">{{ filePath }}</span>
            </div>
          </template>
          <template v-else>
            <svg class="drop-icon" width="36" height="36" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.2" stroke-linecap="round">
              <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/>
              <polyline points="17 8 12 3 7 8"/>
              <line x1="12" y1="3" x2="12" y2="15"/>
            </svg>
            <span class="drop-text">{{ isDragging ? '松开鼠标以选择文件' : '拖拽文件到此处，或点击浏览' }}</span>
            <span class="drop-hint">支持 .json (Postman Collection)</span>
          </template>
        </div>

        <!-- Manual path input (collapsible) -->
        <div class="manual-input-wrap">
          <n-input-group>
            <n-input
              v-model:value="filePath"
              placeholder="或手动输入文件路径"
              size="small"
              :disabled="status === 'importing'"
            >
              <template #prefix>
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
              </template>
            </n-input>
            <n-button size="small" :disabled="status === 'importing'" @click.stop="onBrowse" title="浏览文件夹">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M22 19a2 2 0 01-2 2H4a2 2 0 01-2-2V5a2 2 0 012-2h5l2 3h9a2 2 0 012 2z"/></svg>
            </n-button>
          </n-input-group>
        </div>
      </div>

      <!-- Progress indicator -->
      <div class="import-progress" v-if="status === 'importing'">
        <n-progress type="line" :percentage="importProgress" :show-indicator="false" processing />
        <span class="progress-text">正在导入...</span>
      </div>

      <!-- Error with suggestion -->
      <n-alert v-if="importError" type="error" :bordered="false" style="margin-top: 12px">
        <div class="error-content">
          <span class="error-msg">{{ importError }}</span>
          <span class="error-hint">请确认文件为有效的 Postman Collection JSON 格式</span>
        </div>
      </n-alert>

      <!-- Footer -->
      <div class="modal-footer">
        <n-button @click="onClose" :disabled="status === 'importing'">取消</n-button>
        <n-button
          type="primary"
          :loading="status === 'importing'"
          :disabled="!canImport"
          @click="onImport"
        >
          <template #icon>
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
          </template>
          导入
        </n-button>
      </div>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  NModal, NForm, NFormItem, NSelect, NInput, NInputGroup,
  NButton, NResult, NProgress, NAlert, useMessage,
} from 'naive-ui'
import { ImportPostman } from '../../../wailsjs/go/main/App'
import { OnFileDrop, OnFileDropOff } from '../../../wailsjs/runtime/runtime'

interface Props { show: boolean; projectId: number | null }
const props = defineProps<Props>()
const emit = defineEmits<{ 'update:show': [value: boolean]; imported: [] }>()
const message = useMessage()

type Status = 'idle' | 'importing' | 'result'
const status = ref<Status>('idle')
const format = ref('postman')
const filePath = ref('')
const resultType = ref<'success' | 'error'>('success')
const resultTitle = ref('')
const resultDesc = ref('')
const importProgress = ref(0)
const importError = ref('')
const isDragging = ref(false)

const isLightMode = ref(false)
const modalClass = computed(() => isLightMode.value ? 'import-modal theme-light' : 'import-modal')

onMounted(() => {
  const check = () => { isLightMode.value = !!document.querySelector('.theme-light') }
  check()
  const observer = new MutationObserver(check)
  observer.observe(document.documentElement, { attributes: true, attributeFilter: ['class'], subtree: true })
})

const canImport = computed(() => props.projectId !== null && filePath.value.trim() !== '')

// Extract file name from path
const fileName = computed(() => {
  if (!filePath.value) return ''
  const parts = filePath.value.replace(/\\/g, '/').split('/')
  return parts[parts.length - 1] || filePath.value
})

// Drag-and-drop handlers
function onDragOver(e: DragEvent) {
  isDragging.value = true
}

function onDragLeave(e: DragEvent) {
  isDragging.value = false
}

function onDrop(e: DragEvent) {
  isDragging.value = false
  const files = e.dataTransfer?.files
  if (files && files.length > 0) {
    // Try to get file path from the dropped file
    const file = files[0]
    // In Wails, the file's path property may be available
    const path = (file as any).path || file.name || ''
    if (path) {
      filePath.value = path
    } else {
      message.warning('无法获取文件路径，请使用手动输入路径')
    }
  }
}

function onBrowse() {
  // Wails runtime does not expose OpenFile in this build.
  // Show a hint for the user to input the path manually.
  message.info('请手动输入文件路径，或将文件拖拽到上方区域')
}

function onReset() {
  status.value = 'idle'
  filePath.value = ''
  resultTitle.value = ''
  resultDesc.value = ''
  importError.value = ''
  importProgress.value = 0
}

function onClose() {
  emit('update:show', false)
  status.value = 'idle'
  filePath.value = ''
  importError.value = ''
  importProgress.value = 0
}

async function onImport() {
  if (!props.projectId || !filePath.value.trim()) return
  status.value = 'importing'
  importError.value = ''
  importProgress.value = 10

  // Simulate progress
  const progressInterval = setInterval(() => {
    if (importProgress.value < 90) {
      importProgress.value += Math.random() * 12
    }
  }, 400)

  try {
    const result = await ImportPostman(props.projectId, filePath.value.trim())
    clearInterval(progressInterval)
    importProgress.value = 100
    resultType.value = 'success'
    resultTitle.value = '导入成功'
    resultDesc.value = `成功导入 ${result.collections} 个集合, ${result.requests} 个请求`
    status.value = 'result'
    emit('imported')
  } catch (err: any) {
    clearInterval(progressInterval)
    importProgress.value = 0
    const msg = err?.message || err?.toString() || '未知错误'
    importError.value = msg
    resultType.value = 'error'
    resultTitle.value = '导入失败'
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

/* -- Drag-and-Drop Zone -- */
.drop-zone {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 28px 16px;
  border: 2px dashed var(--border-primary);
  border-radius: var(--radius-lg, 10px);
  cursor: pointer;
  transition: all 0.2s ease;
  background: var(--bg-surface, transparent);
  min-height: 120px;
}
.drop-zone:hover {
  border-color: var(--text-muted);
  background: var(--bg-hover);
}
.drop-zone.dragging {
  border-color: var(--accent);
  background: var(--accent-soft, rgba(0,224,90,0.06));
}
.drop-zone.has-file {
  flex-direction: row;
  padding: 14px 16px;
  border-style: solid;
  border-color: var(--accent);
  background: var(--accent-soft, rgba(0,224,90,0.04));
  min-height: auto;
}
.drop-zone.has-file svg { opacity: 0.5; flex-shrink: 0; }
.drop-icon { opacity: 0.25; }
.drop-text {
  font-size: var(--fs-sm);
  color: var(--text-secondary);
}
.drop-hint {
  font-size: var(--fs-xs);
  color: var(--text-muted);
  opacity: 0.6;
}
.drop-file-info {
  display: flex;
  flex-direction: column;
}
.drop-file-name {
  font-size: var(--fs-sm);
  font-weight: 600;
  color: var(--text-primary);
}
.drop-file-path {
  font-size: var(--fs-xs);
  font-family: var(--font-mono);
  color: var(--text-muted);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 380px;
}

.manual-input-wrap { margin-top: 8px; }

/* -- Progress -- */
.import-progress {
  margin-top: 12px;
  display: flex;
  align-items: center;
  gap: 10px;
}
.import-progress :deep(.n-progress) { flex: 1; }
.progress-text { font-size: var(--fs-xs); color: var(--text-muted); white-space: nowrap; }

/* -- Error -- */
.error-content { display: flex; flex-direction: column; gap: 2px; }
.error-msg { font-size: var(--fs-sm); }
.error-hint { font-size: var(--fs-xs); opacity: 0.7; }

/* -- Footer -- */
.modal-footer { display: flex; justify-content: flex-end; gap: 10px; margin-top: 20px; }
</style>
