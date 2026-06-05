<template>
  <div class="request-docs">
    <div class="docs-toolbar">
      <div class="docs-title">
        <span class="docs-icon">&#128196;</span>
        <span class="docs-label">{{ requestName || '接口文档' }}</span>
        <span class="method-badge" :class="requestMethod?.toLowerCase()">{{ requestMethod }}</span>
      </div>
      <div class="docs-actions">
        <n-button-group size="tiny">
          <n-button :type="viewMode === 'html' ? 'primary' : 'default'" @click="viewMode = 'html'">预览</n-button>
          <n-button :type="viewMode === 'markdown' ? 'primary' : 'default'" @click="viewMode = 'markdown'">源码</n-button>
        </n-button-group>
        <n-button size="tiny" quaternary @click="onRefresh" :loading="loading" title="刷新">
          <template #icon><span class="refresh-icon">&#8635;</span></template>
        </n-button>
        <n-button size="tiny" quaternary @click="onExportMD" :disabled="!markdownContent" title="导出 .md 文件">
          <template #icon><span class="export-icon">MD</span></template>
        </n-button>
        <n-button size="tiny" quaternary @click="onExportHTML" :disabled="!htmlContent" title="导出 .html 文件">
          <template #icon><span class="export-icon">&lt;/&gt;</span></template>
        </n-button>
      </div>
    </div>

    <div v-if="loading && !markdownContent" class="docs-loading">
      <span class="loading-spinner"></span>
      <span>正在生成文档...</span>
    </div>

    <div v-else-if="error" class="docs-error">
      <span class="error-icon">&#9888;</span>
      <span>{{ error }}</span>
    </div>

    <div v-else class="docs-content">
      <iframe v-if="viewMode === 'html'" :srcdoc="htmlContent" class="docs-iframe" sandbox="allow-scripts" />
      <pre v-else class="docs-markdown">{{ markdownContent }}</pre>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { NButton, NButtonGroup } from 'naive-ui'
import { GenerateRequestDocsMarkdown, GenerateRequestDocsHTML } from '../../../wailsjs/go/main/App'

const props = defineProps<{
  requestId: number
  requestName: string
  requestMethod: string
}>()

const viewMode = ref<'html' | 'markdown'>('html')
const loading = ref(false)
const error = ref('')
const markdownContent = ref('')
const htmlContent = ref('')

async function generateDocs() {
  if (!props.requestId) return
  loading.value = true
  error.value = ''
  try {
    const [md, html] = await Promise.all([
      GenerateRequestDocsMarkdown(props.requestId),
      GenerateRequestDocsHTML(props.requestId),
    ])
    markdownContent.value = md
    htmlContent.value = html
  } catch (err: any) {
    error.value = err?.message || err?.toString() || '未知错误'
  } finally {
    loading.value = false
  }
}

function onRefresh() {
  generateDocs()
}

watch(() => props.requestId, (newId) => {
  if (newId) generateDocs()
}, { immediate: true })

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
  const name = props.requestName || 'api-docs'
  downloadFile(markdownContent.value, name + '.md', 'text/markdown')
}

function onExportHTML() {
  const name = props.requestName || 'api-docs'
  downloadFile(htmlContent.value, name + '.html', 'text/html')
}
</script>

<style scoped>
.request-docs {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
  background: var(--bg-base);
}

.docs-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  border-bottom: 1px solid var(--border-primary);
  background: var(--bg-surface);
  flex-shrink: 0;
  gap: 8px;
}

.docs-title {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.docs-icon {
  font-size: var(--fs-md);
  opacity: 0.7;
}

.docs-label {
  font-size: var(--fs-sm);
  font-weight: 600;
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.method-badge {
  font-size: var(--fs-2xs);
  font-weight: 700;
  padding: 1px 5px;
  border-radius: var(--radius-xs);
  letter-spacing: 0.04em;
  text-transform: uppercase;
  flex-shrink: 0;
  font-family: var(--font-mono);
}
.method-badge.get    { background: var(--blue-soft); color: var(--method-get); }
.method-badge.post   { background: rgba(34,197,94,0.1); color: var(--method-post); }
.method-badge.put    { background: var(--amber-soft); color: var(--method-put); }
.method-badge.delete { background: var(--red-soft); color: var(--method-delete); }
.method-badge.patch  { background: var(--purple-soft); color: var(--method-patch); }
.method-badge.head,
.method-badge.options { background: rgba(113,113,122,0.1); color: var(--text-secondary); }

.docs-actions {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
}

.refresh-icon {
  font-size: var(--fs-md);
  display: inline-block;
}

.export-icon {
  font-size: var(--fs-2xs);
  font-weight: 700;
  font-family: var(--font-mono);
  letter-spacing: 0.02em;
  opacity: 0.8;
}

.docs-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 60px;
  color: var(--text-muted);
  font-size: var(--fs-sm);
}

.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid var(--border-primary);
  border-top-color: var(--accent);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.docs-error {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 20px;
  color: var(--red);
  font-size: var(--fs-sm);
}

.error-icon {
  font-size: var(--fs-lg);
}

.docs-content {
  flex: 1;
  min-height: 0;
  overflow: auto;
}

.docs-iframe {
  width: 100%;
  height: 100%;
  min-height: 300px;
  border: none;
}

.docs-markdown {
  padding: 16px 20px;
  margin: 0;
  font-size: var(--fs-sm);
  line-height: 1.7;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: var(--font-mono);
  background: var(--bg-surface);
  color: var(--text-primary);
  min-height: 100%;
  box-sizing: border-box;
}
</style>
