<template>
  <n-modal
    :show="show"
    preset="card"
    title="API Docs"
    style="width: 860px; max-height: 85vh"
    :mask-closable="false"
    @update:show="emit('update:show', $event)"
  >
    <template #header-extra>
      <n-button-group size="tiny">
        <n-button
          :type="viewMode === 'html' ? 'primary' : 'default'"
          @click="viewMode = 'html'"
        >
          HTML
        </n-button>
        <n-button
          :type="viewMode === 'markdown' ? 'primary' : 'default'"
          @click="viewMode = 'markdown'"
        >
          Markdown
        </n-button>
      </n-button-group>
    </template>

    <div v-if="loading" class="docs-loading">
      <n-spin size="medium" />
      <p>Generating docs...</p>
    </div>

    <div v-else-if="error" class="docs-error">
      <n-result status="error" title="Generation failed" :description="error" />
    </div>

    <div v-else class="docs-content">
      <iframe
        v-if="viewMode === 'html'"
        :srcdoc="htmlContent"
        class="docs-iframe"
        sandbox="allow-scripts"
      />
      <pre v-else class="docs-markdown">{{ markdownContent }}</pre>
    </div>

    <template #footer>
      <div class="modal-footer">
        <n-button @click="emit('update:show', false)">Close</n-button>
        <n-button type="primary" :disabled="!markdownContent" @click="onExportMD">
          Export Markdown
        </n-button>
        <n-button type="primary" :disabled="!htmlContent" @click="onExportHTML">
          Export HTML
        </n-button>
      </div>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import {
  NModal,
  NButton,
  NButtonGroup,
  NSpin,
  NResult,
} from 'naive-ui'
import {
  GenerateDocsMarkdown,
  GenerateDocsHTML,
} from '../../../wailsjs/go/main/App'

interface Props {
  show: boolean
  projectId: number | null
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:show': [value: boolean]
}>()

const viewMode = ref<'html' | 'markdown'>('html')
const loading = ref(false)
const error = ref('')
const markdownContent = ref('')
const htmlContent = ref('')

watch(
  () => props.show,
  (val) => {
    if (val && props.projectId) {
      generateDocs()
    }
  }
)

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
    error.value = err?.message || err?.toString() || 'Unknown error'
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
</script>

<style scoped>
.docs-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px;
  gap: 12px;
}

.docs-error {
  padding: 20px;
}

.docs-content {
  min-height: 300px;
  max-height: 55vh;
  overflow: auto;
  border: 1px solid #e8e8e8;
  border-radius: 4px;
}

.docs-iframe {
  width: 100%;
  height: 55vh;
  border: none;
}

.docs-markdown {
  padding: 16px;
  margin: 0;
  font-size: 13px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: Consolas, Monaco, 'Courier New', monospace;
  background: #fafafa;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}
</style>
