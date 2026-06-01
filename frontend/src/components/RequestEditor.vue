<script lang="ts" setup>
import { computed, ref } from 'vue'
import { t } from '../i18n'
import { NInput, NSelect, NButton, NTabs, NTabPane, NIcon, useMessage } from 'naive-ui'
import { Send } from '@vicons/ionicons5'
import KeyValueEditor from './KeyValueEditor.vue'
import type { KeyValuePair } from '../stores/tabs'
import { useRequestStore } from '../stores/request'
import { useEnvironmentStore } from '../stores/environment'
import { useVariableResolver } from '../composables/useVariableResolver'
import { useTabsStore } from '../stores/tabs'
import { useProjectStore } from '../stores/project'
import { SendRequest } from '../../wailsjs/go/handlers/RequestHandler'
import { RecordHistory } from '../../wailsjs/go/handlers/HistoryHandler'

const httpMethods = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'HEAD', 'OPTIONS']
const requestStore = useRequestStore()
const envStore = useEnvironmentStore()
const tabsStore = useTabsStore()
const projectStore = useProjectStore()
const { resolve } = useVariableResolver()
const message = useMessage()

const data = computed(() => tabsStore.activeTabData)
const method = computed({ get: () => data.value?.method ?? 'GET', set: (val) => tabsStore.updateTabData({ method: val }) })
const url = computed({ get: () => data.value?.url ?? '', set: (val) => tabsStore.updateTabData({ url: val }) })
const params = computed({ get: () => data.value?.params ?? [], set: (val) => tabsStore.updateTabData({ params: val }) })
const headers = computed({ get: () => data.value?.headers ?? [{ key: 'Content-Type', value: 'application/json', enabled: true }], set: (val) => tabsStore.updateTabData({ headers: val }) })
const body = computed({ get: () => data.value?.body ?? '', set: (val) => tabsStore.updateTabData({ body: val }) })
const bodyType = computed({ get: () => data.value?.bodyType ?? 'none', set: (val) => tabsStore.updateTabData({ bodyType: val }) })
const auth = computed({ get: () => data.value?.auth ?? [], set: (val) => tabsStore.updateTabData({ auth: val }) })
const sending = ref(false)

function buildHeadersMap(): Record<string, string> {
  const map: Record<string, string> = {}
  for (const h of headers.value) {
    if (h.enabled && h.key) map[h.key] = h.value
  }
  return map
}

function appendParamsToURL(baseURL: string): string {
  const active = params.value.filter(p => p.enabled && p.key)
  if (active.length === 0) return baseURL
  const sep = baseURL.includes('?') ? '&' : '?'
  return baseURL + sep + active.map(p => `${encodeURIComponent(p.key)}=${encodeURIComponent(p.value)}`).join('&')
}

async function handleSend() {
  if (!url.value.trim()) { message.warning(t('request.urlRequired')); return }
  sending.value = true
  try {
    const resolvedURL = resolve(appendParamsToURL(url.value.trim()), envStore.activeVariables)
    const resolvedBody = bodyType.value === 'none' ? '' : resolve(body.value, envStore.activeVariables)
    const resolvedHeaders = buildHeadersMap()
    for (const k of Object.keys(resolvedHeaders)) resolvedHeaders[k] = resolve(resolvedHeaders[k], envStore.activeVariables)
    const resp = await SendRequest({ Method: method.value, URL: resolvedURL, Headers: resolvedHeaders, Body: resolvedBody })
    requestStore.setLastResponse(resp)
    message.success(`${resp.status} ${resp.status_text} �?${resp.duration_ms}ms`)
    if (resolvedURL.startsWith('http') && projectStore.currentProject) {
      RecordHistory({
        ProjectID: projectStore.currentProject.id,
        RequestID: '',
        Method: method.value,
        URL: resolvedURL,
        Headers: JSON.stringify(buildHeadersMap()),
        Body: resolvedBody,
        ResponseStatus: resp.status,
        ResponseBody: resp.body,
        ResponseHeaders: JSON.stringify(resp.headers),
        DurationMs: resp.duration_ms,
      }).catch(() => {})
    }
  } catch (e: any) { message.error(e.message || t('request.requestFailed'))
  } finally { sending.value = false }
}
</script>

<template>
  <div v-if="data" class="request-editor">
    <div class="url-row">
      <NSelect :options="httpMethods.map(m => ({ label: m, value: m }))" v-model:value="method" style="width: 110px" size="small" />
      <NInput v-model:value="url" :placeholder="$t('request.placeholder')" size="small" class="url-input" />
      <NButton type="primary" size="small" :loading="sending" @click="handleSend">
        <template #icon><NIcon><Send /></NIcon></template>
        {{ $t('request.send') }}
      </NButton>
    </div>
    <NTabs type="line" size="small" class="editor-tabs">
      <NTabPane name="params" :tab="$t('request.params')">
        <KeyValueEditor v-model="params" :key-placeholder="$t('request.paramName')" :value-placeholder="$t('request.value')" />
      </NTabPane>
      <NTabPane name="headers" :tab="$t('request.headers')">
        <KeyValueEditor v-model="headers" />
      </NTabPane>
      <NTabPane name="body" :tab="$t('request.body')">
        <div class="body-controls">
          <NSelect :options="[
            { label: $t('request.bodyNone'), value: 'none' },
            { label: $t('request.bodyJSON'), value: 'json' },
            { label: $t('request.bodyText'), value: 'text' },
            { label: $t('request.bodyForm'), value: 'form' },
          ]" v-model:value="bodyType" size="tiny" style="width: 120px; margin-bottom: 8px;" />
          <NInput v-if="bodyType !== 'none'" v-model:value="body" type="textarea" :rows="6" :placeholder="$t('request.bodyPlaceholder')" class="body-input" />
        </div>
      </NTabPane>
      <NTabPane name="auth" :tab="$t('request.auth')">
        <KeyValueEditor v-model="auth" :key-placeholder="$t('request.authKey')" :value-placeholder="$t('request.authValue')" />
      </NTabPane>
    </NTabs>
  </div>
  <div v-else class="request-editor empty">
    <p class="empty-text">{{ $t('request.openRequest') }}</p>
  </div>
</template>

<style scoped>
.request-editor { padding: 12px 16px; border-bottom: 1px solid var(--border-color); }
.request-editor.empty { display: flex; align-items: center; justify-content: center; min-height: 120px; }
.empty-text { color: #999; font-size: 14px; }
.url-row { display: flex; gap: 8px; margin-bottom: 8px; }
.url-input { flex: 1; }
.editor-tabs { margin-top: 4px; }
.body-controls { padding: 8px 0; }
.body-input { font-family: 'Cascadia Code', 'Fira Code', 'Consolas', monospace; font-size: 13px; }
</style>
