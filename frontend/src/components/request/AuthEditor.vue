<template>
  <div class="auth-editor">
    <n-select
      v-model:value="authType"
      :options="authTypeOptions"
      size="small"
      class="type-select"
    />

    <div v-if="authType === 'none'" class="auth-empty">
      <span class="hint">此请求不使用认证</span>
    </div>

    <div v-else-if="authType === 'bearer'" class="auth-form">
      <label>Token</label>
      <n-input v-model:value="token" size="small" placeholder="token 或 {{variable}}" />
      <span class="hint-suffix">前缀 Bearer 自动添加</span>
    </div>

    <div v-else-if="authType === 'basic'" class="auth-form">
      <label>Username</label>
      <n-input v-model:value="username" size="small" placeholder="用户名" />
      <label>Password</label>
      <n-input v-model:value="password" size="small" type="password" placeholder="密码" />
    </div>

    <div v-else-if="authType === 'apikey'" class="auth-form">
      <label>Key</label>
      <n-input v-model:value="apiKey" size="small" placeholder="X-API-Key" />
      <label>Value</label>
      <n-input v-model:value="apiValue" size="small" placeholder="key 值或 {{variable}}" />
      <label>Add to</label>
      <n-select v-model:value="apiAddTo" size="small" :options="addToOptions" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { NSelect, NInput } from 'naive-ui'

const props = defineProps<{
  modelValue: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', v: string): void
}>()

const authTypeOptions = [
  { label: 'No Auth', value: 'none' },
  { label: 'Bearer Token', value: 'bearer' },
  { label: 'Basic Auth', value: 'basic' },
  { label: 'API Key', value: 'apikey' },
]

const addToOptions = [
  { label: 'Header', value: 'header' },
  { label: 'Query Params', value: 'query' },
]

const authType = ref('none')
const token = ref('')
const username = ref('')
const password = ref('')
const apiKey = ref('')
const apiValue = ref('')
const apiAddTo = ref('header')

function parseAuth(raw: string) {
  try {
    const obj = JSON.parse(raw)
    authType.value = obj.type || 'none'
    token.value = obj.token || ''
    username.value = obj.username || ''
    password.value = obj.password || ''
    apiKey.value = obj.key || ''
    apiValue.value = obj.value || ''
    apiAddTo.value = obj.addTo || 'header'
  } catch {
    authType.value = 'none'
  }
}

function syncAuth() {
  let auth: Record<string, any> = { type: authType.value }
  if (authType.value === 'bearer') auth.token = token.value
  if (authType.value === 'basic') { auth.username = username.value; auth.password = password.value }
  if (authType.value === 'apikey') { auth.key = apiKey.value; auth.value = apiValue.value; auth.addTo = apiAddTo.value }
  emit('update:modelValue', JSON.stringify(auth))
}

watch(() => props.modelValue, parseAuth, { immediate: true })
watch([authType, token, username, password, apiKey, apiValue, apiAddTo], syncAuth)
</script>

<style scoped>
.auth-editor {
  padding: 12px;
}
.type-select {
  margin-bottom: 12px;
}
.auth-empty {
  padding: 20px;
}
.auth-form {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.auth-form label {
  font-size: 11px;
  color: #888;
}
.hint-suffix {
  font-size: 10px;
  color: #aaa;
}
.hint {
  color: #aaa;
  font-size: 12px;
}
</style>
