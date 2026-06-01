<script lang="ts" setup>
import { computed, ref } from 'vue'
import { NTag, NInput, NTabs, NTabPane, NEmpty, NSelect, NButton, useMessage } from 'naive-ui'
import { useRequestStore } from '../stores/request'
import { useTabsStore } from '../stores/tabs'
import { GenerateCode } from '../../wailsjs/go/handlers/RequestHandler'
import { t } from '../i18n'

const requestStore = useRequestStore()
const tabsStore = useTabsStore()
const message = useMessage()
const response = computed(() => requestStore.lastResponse)
const codeLang = ref('curl')
const codeOutput = ref('')
const codeLoading = ref(false)

const statusType = computed(() => {
  if (!response.value) return 'default'
  if (response.value.status >= 200 && response.value.status < 300) return 'success'
  if (response.value.status >= 300 && response.value.status < 400) return 'warning'
  if (response.value.status >= 400) return 'error'
  return 'default'
})
const headersText = computed(() => {
  if (!response.value) return ''
  return Object.entries(response.value.headers).map(([k, v]) => `${k}: ${v.join(', ')}`).join('\n')
})

const codeLangs = [
  { label: 'cURL', value: 'curl' },
  { label: 'JavaScript (fetch)', value: 'javascript' },
  { label: 'Python (requests)', value: 'python' },
  { label: 'Go (net/http)', value: 'go' },
]

function formatJSON(text: string): string {
  try { return JSON.stringify(JSON.parse(text), null, 2) } catch { return text }
}

async function generateCode() {
  const data = tabsStore.activeTabData
  if (!data) return
  codeLoading.value = true
  try {
    codeOutput.value = await GenerateCode({
      Method: data.method,
      URL: data.url,
      Headers: Object.fromEntries(data.headers.filter(h => h.enabled && h.key).map(h => [h.key, h.value])),
      Body: data.body,
      BodyType: data.bodyType,
    }, codeLang.value)
  } catch (e: any) {
    codeOutput.value = 'Error: ' + e.message
  } finally {
    codeLoading.value = false
  }
}

async function copyCode() {
  try {
    await navigator.clipboard.writeText(codeOutput.value)
    message.success(t('codegen.copied'))
  } catch { message.error('Copy failed') }
}
</script>

<template>
  <div class="response-viewer">
    <div v-if="response" class="response-header">
      <NTag :type="statusType" size="small">{{ response.status }} {{ response.status_text }}</NTag>
      <span class="duration">{{ response.duration_ms }}ms</span>
    </div>
    <div v-if="!response" class="response-empty">
      <NEmpty :description="$t('response.empty')" />
    </div>
    <div v-else class="response-body">
      <NTabs type="line" size="small">
        <NTabPane name="body" :tab="$t('response.body')">
          <NInput :value="formatJSON(response.body)" type="textarea" :rows="12" readonly class="resp-input" />
        </NTabPane>
        <NTabPane name="headers" :tab="$t('response.headers')">
          <NInput :value="headersText" type="textarea" :rows="8" readonly class="resp-input" />
        </NTabPane>
        <NTabPane name="code" :tab="$t('codegen.title')">
          <div class="codegen-controls">
            <div class="codegen-top">
              <NSelect :options="codeLangs" v-model:value="codeLang" size="tiny" style="width: 160px" />
              <NButton size="tiny" @click="generateCode">{{ $t('codegen.title') }}</NButton>
              <NButton v-if="codeOutput" size="tiny" @click="copyCode">{{ $t('codegen.copy') }}</NButton>
            </div>
            <NInput v-if="codeOutput" :value="codeOutput" type="textarea" :rows="10" readonly class="resp-input code-output" />
          </div>
        </NTabPane>
      </NTabs>
    </div>
  </div>
</template>

<style scoped>
.response-viewer { flex: 1; display: flex; flex-direction: column; overflow: hidden; }
.response-header { display: flex; align-items: center; gap: 8px; padding: 8px 16px; border-bottom: 1px solid var(--border-color); font-size: 13px; }
.duration { color: #999; font-size: 12px; }
.response-empty { flex: 1; display: flex; align-items: center; justify-content: center; }
.response-body { flex: 1; padding: 8px 16px; overflow: auto; }
.resp-input { font-family: 'Cascadia Code', 'Fira Code', 'Consolas', monospace; font-size: 13px; }
.codegen-controls { padding: 8px 0; }
.codegen-top { display: flex; gap: 8px; align-items: center; margin-bottom: 8px; }
.code-output { font-family: 'Cascadia Code', 'Fira Code', 'Consolas', monospace; font-size: 13px; }
</style>
