<script lang="ts" setup>
import { computed, ref } from 'vue'
import { t } from '../i18n'
import { NInput, NSelect, NButton, NTabs, NTabPane, NIcon, NCheckbox, NPopover, NSlider, NSpace, useMessage } from 'naive-ui'
import { Send, Settings as SettingsIcon, Add, Trash } from '@vicons/ionicons5'
import KeyValueEditor from './KeyValueEditor.vue'
import type { KeyValuePair, BodyFileItem } from '../stores/tabs'
import { useRequestStore } from '../stores/request'
import { useEnvironmentStore } from '../stores/environment'
import { useVariableResolver } from '../composables/useVariableResolver'
import { useTabsStore } from '../stores/tabs'
import { useProjectStore } from '../stores/project'
import { SendRequest, CreateRequest, UpdateRequest } from '../../wailsjs/go/handlers/RequestHandler'
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
const bodyFiles = computed({ get: () => data.value?.bodyFiles ?? [], set: (val) => tabsStore.updateTabData({ bodyFiles: val }) })
const authType = computed({ get: () => data.value?.authType ?? 'none', set: (val) => tabsStore.updateTabData({ authType: val }) })
const authData = computed({ get: () => data.value?.authData ?? {}, set: (val) => tabsStore.updateTabData({ authData: val }) })
const timeoutMs = computed({ get: () => data.value?.timeoutMs ?? 30000, set: (val) => tabsStore.updateTabData({ timeoutMs: val }) })
const followRedirect = computed({ get: () => data.value?.followRedirect ?? true, set: (val) => tabsStore.updateTabData({ followRedirect: val }) })
const activeTab = computed(() => tabsStore.activeTab)
const requestName = computed({
  get: () => activeTab.value?.title ?? '',
  set: (val) => { if (activeTab.value) tabsStore.updateTabTitle(activeTab.value.id, val) },
})
const sending = ref(false)
const saving = ref(false)
const showSettingsPopover = ref(false)

const authOptions = [
  { label: t('auth.typeNone'), value: 'none' },
  { label: t('auth.typeBasic'), value: 'basic' },
  { label: t('auth.typeBearer'), value: 'bearer' },
  { label: t('auth.typeDigest'), value: 'digest' },
  { label: t('auth.typeOAuth2'), value: 'oauth2' },
]

const bodyTypeOptions = [
  { label: t('request.bodyNone'), value: 'none' },
  { label: t('request.bodyJSON'), value: 'json' },
  { label: t('request.bodyText'), value: 'text' },
  { label: t('bodyType.formData'), value: 'form-data' },
  { label: t('bodyType.urlencoded'), value: 'urlencoded' },
  { label: t('bodyType.binary'), value: 'binary' },
]

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

function addBodyFile() {
  bodyFiles.value = [...bodyFiles.value, { key: '', value: '', file_path: '', enabled: true }]
}

function updateBodyFile(index: number, field: string, val: any) {
  const next = bodyFiles.value.map((f, i) => i === index ? { ...f, [field]: val } : f)
  bodyFiles.value = next
}

function removeBodyFile(index: number) {
  bodyFiles.value = bodyFiles.value.filter((_, i) => i !== index)
}

function updateAuth(key: string, val: string) {
  const next = { ...authData.value }
  next[key] = val
  tabsStore.updateTabData({ authData: next })
}

async function handleSave() {
  const tab = tabsStore.activeTab
  if (!tab || !projectStore.currentProject) return
  saving.value = true
  try {
    const colId = projectStore.collections[0]?.id || ''
    if (!colId) { message.warning('No collection to save to'); return }
    const headersJSON = JSON.stringify(headers.value.filter(h => h.enabled))
    const paramsJSON = JSON.stringify(params.value.filter(p => p.enabled))
    const bodyJSON = JSON.stringify({ body_type: bodyType.value, content: body.value })
    const authJSON = JSON.stringify({ type: authType.value, ...authData.value })
    if (tab.requestId) {
      await UpdateRequest(tab.requestId, requestName.value, method.value, url.value, headersJSON, paramsJSON, bodyJSON, authJSON, '', 0)
      message.success('Saved')
    } else {
      const r = await CreateRequest(colId, requestName.value, method.value, url.value, headersJSON, paramsJSON, bodyJSON, authJSON, '', 0)
      if (r) {
        tab.requestId = r.id
        message.success('Created')
      }
    }
  } catch (e: any) { message.error(e.message || 'Save failed')
  } finally { saving.value = false }
}

async function handleSend() {
  if (!url.value.trim()) { message.warning(t('request.urlRequired')); return }
  sending.value = true
  try {
    const resolvedURL = resolve(appendParamsToURL(url.value.trim()), envStore.activeVariables)
    const resolvedBody = bodyType.value === 'none' ? '' : resolve(body.value, envStore.activeVariables)
    const resolvedHeaders = buildHeadersMap()
    for (const k of Object.keys(resolvedHeaders)) resolvedHeaders[k] = resolve(resolvedHeaders[k], envStore.activeVariables)

    const bodyFilesPayload = bodyFiles.value.filter(f => f.enabled).map(f => ({
      key: f.key,
      value: f.value,
      file_path: f.file_path,
      enabled: f.enabled,
    }))

    const resp = await SendRequest({
      Method: method.value,
      URL: resolvedURL,
      Headers: resolvedHeaders,
      Body: resolvedBody,
      BodyType: bodyType.value,
      BodyFiles: bodyFilesPayload,
      AuthType: authType.value,
      AuthData: authData.value,
      TimeoutMs: timeoutMs.value,
      FollowRedirect: followRedirect.value,
    } as any)
    requestStore.setLastResponse(resp)
    message.success(`${resp.status} ${resp.status_text} ${resp.duration_ms}ms`)
    if (resolvedURL.startsWith('http') && projectStore.currentProject) {
      ;(RecordHistory as any)({
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
    <div class="name-row">
      <NInput v-model:value="requestName" placeholder="Request name" size="small" class="name-input" />
      <NButton size="tiny" :loading="saving" @click="handleSave">
        {{ activeTab?.requestId ? 'Save' : 'Create' }}
      </NButton>
    </div>
    <div class="url-row">
      <NSelect :options="httpMethods.map(m => ({ label: m, value: m }))" v-model:value="method" style="width: 110px" size="small" />
      <NInput v-model:value="url" :placeholder="$t('request.placeholder')" size="small" class="url-input" />
      <NPopover v-model:show="showSettingsPopover" trigger="click" placement="bottom">
        <template #trigger>
          <NButton quaternary size="small" class="settings-btn">
            <template #icon><NIcon><SettingsIcon /></NIcon></template>
          </NButton>
        </template>
        <div class="settings-popover">
          <div class="settings-row">
            <span class="settings-label">{{ $t('request.timeout') }}</span>
            <NSlider v-model:value="timeoutMs" :min="1000" :max="120000" :step="1000" style="width: 120px" />
            <span class="settings-value">{{ (timeoutMs / 1000).toFixed(0) }}s</span>
          </div>
          <div class="settings-row">
            <NCheckbox v-model:checked="followRedirect">{{ $t('request.followRedirect') }}</NCheckbox>
          </div>
        </div>
      </NPopover>
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
        <KeyValueEditor v-model="headers" show-presets />
      </NTabPane>
      <NTabPane name="body" :tab="$t('request.body')">
        <div class="body-controls">
          <NSelect :options="bodyTypeOptions" v-model:value="bodyType" size="tiny" style="width: 160px; margin-bottom: 8px;" />
          <NInput v-if="bodyType !== 'none' && bodyType !== 'form-data' && bodyType !== 'binary'" v-model:value="body" type="textarea" :rows="6" :placeholder="$t('request.bodyPlaceholder')" class="body-input" />
          <div v-if="bodyType === 'form-data'" class="form-data-editor">
            <div v-for="(f, i) in bodyFiles" :key="i" class="form-data-row">
              <input type="checkbox" :checked="f.enabled" class="kv-checkbox" @change="updateBodyFile(i, 'enabled', ($event.target as HTMLInputElement).checked)" />
              <NInput :value="f.key" size="tiny" placeholder="Key" class="kv-input" @update:value="updateBodyFile(i, 'key', $event)" />
              <NInput :value="f.value" size="tiny" placeholder="Value" class="kv-input" @update:value="updateBodyFile(i, 'value', $event)" />
              <NInput :value="f.file_path" size="tiny" :placeholder="$t('bodyType.filePath')" class="kv-input" @update:value="updateBodyFile(i, 'file_path', $event)" />
              <NButton quaternary circle size="tiny" @click="removeBodyFile(i)">
                <template #icon><NIcon><Trash /></NIcon></template>
              </NButton>
            </div>
            <NButton size="tiny" quaternary @click="addBodyFile">
              <template #icon><NIcon><Add /></NIcon></template>
              Add
            </NButton>
          </div>
          <NInput v-if="bodyType === 'binary'" v-model:value="body" type="textarea" :rows="4" placeholder="Base64 encoded data..." class="body-input mono" />
        </div>
      </NTabPane>
      <NTabPane name="auth" :tab="$t('request.auth')">
        <div class="auth-controls">
          <NSelect :options="authOptions" v-model:value="authType" size="tiny" style="width: 160px; margin-bottom: 8px;" />
          <div v-if="authType === 'basic'" class="auth-fields">
            <NInput :value="authData['username'] || ''" :placeholder="$t('auth.username')" size="tiny" class="auth-input" @update:value="updateAuth('username', $event)" />
            <NInput :value="authData['password'] || ''" :placeholder="$t('auth.password')" type="password" size="tiny" class="auth-input" @update:value="updateAuth('password', $event)" />
          </div>
          <div v-if="authType === 'bearer'" class="auth-fields">
            <NInput :value="authData['token'] || ''" :placeholder="$t('auth.token')" size="tiny" class="auth-input" @update:value="updateAuth('token', $event)" />
          </div>
          <div v-if="authType === 'digest'" class="auth-fields">
            <NInput :value="authData['username'] || ''" :placeholder="$t('auth.username')" size="tiny" class="auth-input" @update:value="updateAuth('username', $event)" />
            <NInput :value="authData['password'] || ''" :placeholder="$t('auth.password')" type="password" size="tiny" class="auth-input" @update:value="updateAuth('password', $event)" />
          </div>
          <div v-if="authType === 'oauth2'" class="auth-fields">
            <NInput :value="authData['client_id'] || ''" :placeholder="$t('auth.clientId')" size="tiny" class="auth-input" @update:value="updateAuth('client_id', $event)" />
            <NInput :value="authData['client_secret'] || ''" :placeholder="$t('auth.clientSecret')" size="tiny" class="auth-input" @update:value="updateAuth('client_secret', $event)" />
            <NInput :value="authData['token_url'] || ''" :placeholder="$t('auth.tokenUrl')" size="tiny" class="auth-input" @update:value="updateAuth('token_url', $event)" />
          </div>
        </div>
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
.name-row { display: flex; gap: 8px; margin-bottom: 6px; }
.name-input { flex: 1; }
.url-row { display: flex; gap: 8px; margin-bottom: 8px; }
.url-input { flex: 1; }
.settings-btn { flex-shrink: 0; }
.settings-popover { padding: 8px; display: flex; flex-direction: column; gap: 8px; min-width: 220px; }
.settings-row { display: flex; align-items: center; gap: 8px; }
.settings-label { font-size: 12px; white-space: nowrap; }
.settings-value { font-size: 12px; min-width: 30px; text-align: right; }
.editor-tabs { margin-top: 4px; }
.body-controls { padding: 8px 0; }
.body-input { font-family: 'Cascadia Code', 'Fira Code', 'Consolas', monospace; font-size: 13px; }
.body-input.mono { font-family: 'Cascadia Code', 'Fira Code', 'Consolas', monospace; font-size: 13px; }
.auth-controls { padding: 8px 0; }
.auth-fields { display: flex; flex-direction: column; gap: 6px; padding: 4px 0; }
.auth-input { width: 100%; }
.form-data-editor { padding: 4px 0; }
.form-data-row { display: flex; align-items: center; gap: 4px; margin-bottom: 4px; }
.kv-checkbox { width: 14px; height: 14px; cursor: pointer; flex-shrink: 0; }
.kv-input { flex: 1; }
</style>
