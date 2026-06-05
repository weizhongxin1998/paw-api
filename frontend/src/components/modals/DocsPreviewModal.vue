<template>
  <n-modal
    :show="show"
    preset="card"
    :title="$t('docs.title')"
    :class="modalClass"
    style="width: 880px; max-height: 85vh"
    :mask-closable="false"
    @update:show="emit('update:show', $event)"
  >
    <!-- Segmented control mode toggle -->
    <template #header-extra>
      <n-radio-group v-model:value="viewMode" size="small">
        <n-radio-button value="html">
          <span class="mode-label">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><rect x="2" y="3" width="20" height="14" rx="2" ry="2"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="12" y1="17" x2="12" y2="21"/></svg>
          {{ $t('docs.preview') }}
          </span>
        </n-radio-button>
        <n-radio-button value="markdown">
          <span class="mode-label">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/></svg>
          {{ $t('docs.source') }}
          </span>
        </n-radio-button>
      </n-radio-group>
    </template>

    <!-- Loading state -->
    <div v-if="loading" class="docs-loading">
      <n-spin size="medium" />
      <p class="loading-text">{{ $t('docs.loading') }}</p>
    </div>

    <!-- Error state -->
    <div v-else-if="error" class="docs-error">
      <n-result status="error" :title="$t('docs.generateFailed')" :description="error" />
    </div>

    <!-- Empty state -->
    <div v-else-if="!htmlContent && !markdownContent" class="docs-empty">
      <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" opacity="0.2">
        <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/>
        <polyline points="14 2 14 8 20 8"/>
        <line x1="16" y1="13" x2="8" y2="13"/>
        <line x1="16" y1="17" x2="8" y2="17"/>
      </svg>
      <span class="empty-text">{{ $t('docs.empty') }}</span>
      <span class="empty-hint">{{ $t('docs.emptyHint') }}</span>
    </div>

    <!-- HTML preview with browser-like frame -->
    <div v-else-if="viewMode === 'html'" class="docs-content">
      <div class="browser-frame">
        <div class="browser-chrome">
          <div class="browser-dots">
            <span class="dot red"></span>
            <span class="dot yellow"></span>
            <span class="dot green"></span>
          </div>
          <div class="browser-url-bar">
            <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="M2 12h20"/><path d="M12 2a15.3 15.3 0 014 10 15.3 15.3 0 01-4 10 15.3 15.3 0 01-4-10 15.3 15.3 0 014-10z"/></svg>
            <span class="browser-url-text">api-docs://preview</span>
          </div>
        </div>
        <iframe :srcdoc="htmlContent" class="docs-iframe" sandbox="allow-scripts" />
      </div>
    </div>

    <!-- Markdown source view with copy button -->
    <div v-else class="docs-content">
      <div class="markdown-toolbar">
        <span class="md-lang-badge">Markdown</span>
        <n-button size="tiny" quaternary @click="onCopyMarkdown" :disabled="!markdownContent">
          <template #icon>
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"/><path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/></svg>
          </template>
          {{ $t('docs.copyAll') }}
        </n-button>
      </div>
      <pre class="docs-markdown">{{ markdownContent }}</pre>
    </div>

    <template #footer>
      <div class="modal-footer">
        <span class="footer-hint" v-if="viewMode === 'html'">{{ $t('docs.printHint') }}</span>
        <n-button @click="emit('update:show', false)">{{ $t('common.close') }}</n-button>
        <n-button type="primary" :disabled="!markdownContent" @click="onExportMD" secondary>
          <template #icon>
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
          </template>
          {{ $t('docs.exportMD') }}
        </n-button>
        <n-button type="primary" :disabled="!htmlContent" @click="onExportHTML">
          <template #icon>
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
          </template>
          {{ $t('docs.exportHTML') }}
        </n-button>
      </div>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted, computed } from 'vue'
import { NModal, NButton, NRadioGroup, NRadioButton, NSpin, NResult, useMessage } from 'naive-ui'
import { GenerateDocsMarkdown, GenerateDocsHTML } from '../../../wailsjs/go/main/App'
import { ClipboardSetText } from '../../../wailsjs/runtime/runtime'
import { useI18n } from 'vue-i18n'

interface Props { show: boolean; projectId: number | null }
const props = defineProps<Props>()
const emit = defineEmits<{ 'update:show': [value: boolean] }>()
const message = useMessage()
const { t } = useI18n()

const viewMode = ref<'html' | 'markdown'>('html')
const loading = ref(false)
const error = ref('')
const markdownContent = ref('')
const htmlContent = ref('')

const isLightMode = ref(false)
const modalClass = computed(() => isLightMode.value ? 'docs-preview-modal theme-light' : 'docs-preview-modal')

watch(() => props.show, (val) => { if (val && props.projectId) generateDocs() })

async function generateDocs() {
  if (!props.projectId) return
  loading.value = true
  error.value = ''
  try {
    const [md, html] = await Promise.all([
      GenerateDocsMarkdown(props.projectId),
      GenerateDocsHTML(props.projectId),
    ])
    markdownContent.value = md
    htmlContent.value = html
  } catch (err: any) {
    error.value = err?.message || err?.toString() || t('common.unknownError')
  } finally {
    loading.value = false
  }
}

function downloadFile(content: string, filename: string, mime: string) {
  const blob = new Blob([content], { type: mime })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  a.click()
  URL.revokeObjectURL(url)
}

function onExportMD() {
  downloadFile(markdownContent.value, 'api-docs.md', 'text/markdown')
}

function onExportHTML() {
  downloadFile(htmlContent.value, 'api-docs.html', 'text/html')
}

async function onCopyMarkdown() {
  try {
    await ClipboardSetText(markdownContent.value)
    message.success(t('docs.copySuccess'))
  } catch {
    message.error(t('docs.copyFailed'))
  }
}

// Ctrl+P to print/export HTML preview
function handleKeydown(e: KeyboardEvent) {
  if (!props.show) return
  if ((e.ctrlKey || e.metaKey) && e.key === 'p') {
    e.preventDefault()
    if (viewMode.value === 'html' && htmlContent.value) {
      printHtmlContent()
    }
  }
}

function printHtmlContent() {
  const printWindow = window.open('', '_blank', 'width=800,height=600')
  if (printWindow) {
    printWindow.document.write(htmlContent.value)
    printWindow.document.close()
    printWindow.focus()
    setTimeout(() => { printWindow.print() }, 400)
  } else {
    message.warning(t('docs.printBlocked'))
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
  const check = () => { isLightMode.value = !!document.querySelector('.theme-light') }
  check()
  const observer = new MutationObserver(check)
  observer.observe(document.documentElement, { attributes: true, attributeFilter: ['class'], subtree: true })
})
onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})
</script>

<style scoped>
.mode-label {
  display: flex;
  align-items: center;
  gap: 4px;
}

/* -- Loading -- */
.docs-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px;
  gap: 12px;
  color: var(--text-muted);
}
.loading-text { font-size: var(--fs-sm); }

/* -- Error -- */
.docs-error { padding: 20px; }

/* -- Empty State -- */
.docs-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 60px;
  color: var(--text-muted);
}
.empty-text { font-size: var(--fs-base); font-weight: 500; }
.empty-hint { font-size: var(--fs-xs); opacity: 0.6; }

/* -- Content Area -- */
.docs-content {
  min-height: 300px;
  max-height: 58vh;
  overflow: auto;
}

/* -- Browser-like Frame -- */
.browser-frame {
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg, 10px);
  overflow: hidden;
}
.browser-chrome {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px 12px;
  background: var(--bg-hover, rgba(255,255,255,0.02));
  border-bottom: 1px solid var(--border-primary);
}
.browser-dots {
  display: flex;
  gap: 5px;
}
.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}
.dot.red { background: #ef4444; }
.dot.yellow { background: #f59e0b; }
.dot.green { background: #22c55e; }
.browser-url-bar {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 3px 10px;
  background: var(--bg-base, #0a0a0b);
  border-radius: 6px;
  font-size: var(--fs-xs);
  color: var(--text-muted);
}
.browser-url-bar svg { opacity: 0.4; flex-shrink: 0; }
.browser-url-text { font-family: var(--font-mono); }
.docs-iframe {
  width: 100%;
  height: 54vh;
  border: none;
  display: block;
}

/* -- Markdown Source View -- */
.markdown-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 4px 8px;
  border-bottom: 1px solid var(--border-primary);
}
.md-lang-badge {
  font-size: var(--fs-2xs);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: var(--text-muted);
  background: var(--bg-hover);
  padding: 2px 8px;
  border-radius: 4px;
}
.docs-markdown {
  padding: 16px;
  margin: 0;
  font-size: var(--fs-sm);
  line-height: 1.7;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: var(--font-mono);
  background: var(--bg-surface, transparent);
  color: var(--text-primary);
  max-height: 53vh;
  overflow-y: auto;
}

/* -- Footer -- */
.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  align-items: center;
}
.footer-hint {
  font-size: var(--fs-xs);
  color: var(--text-muted);
  opacity: 0.6;
  margin-right: auto;
}
</style>
