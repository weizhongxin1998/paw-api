<script lang="ts" setup>
import { computed } from 'vue'
import { NTag, NInput, NTabs, NTabPane, NEmpty } from 'naive-ui'
import { useRequestStore } from '../stores/request'

const requestStore = useRequestStore()
const response = computed(() => requestStore.lastResponse)
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
function formatJSON(text: string): string {
  try { return JSON.stringify(JSON.parse(text), null, 2) } catch { return text }
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
</style>
