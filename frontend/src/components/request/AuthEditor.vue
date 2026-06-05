<template>
  <div class="auth-editor">
    <!-- ── Auth type cards selector ── -->
    <div class="auth-cards">
      <button
        v-for="opt in authTypeOptions"
        :key="opt.value"
        :class="{ active: authType === opt.value }"
        class="auth-card"
        @click="authType = opt.value"
      >
        <div class="auth-card-icon" v-html="opt.icon"></div>
        <div class="auth-card-info">
          <span class="auth-card-label">{{ opt.label }}</span>
          <span class="auth-card-desc">{{ opt.desc }}</span>
        </div>
        <div v-if="authType === opt.value" class="auth-card-check">
          <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="20 6 9 17 4 12" />
          </svg>
        </div>
      </button>
    </div>

    <!-- ── Auth forms ── -->
    <div class="auth-form-wrap">
      <!-- None -->
      <div v-if="authType === 'none'" class="auth-empty">
        <svg class="auth-empty-icon" viewBox="0 0 24 24" width="32" height="32" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
          <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
          <path d="M7 11V7a5 5 0 0110 0v4" />
          <line x1="4" y1="4" x2="20" y2="20" />
        </svg>
        <span class="auth-empty-text">{{ $t('auth.noneText') }}</span>
        <span class="auth-empty-hint">{{ $t('auth.noneHint') }}</span>
      </div>

      <!-- Bearer Token -->
      <div v-else-if="authType === 'bearer'" class="auth-form">
        <div class="auth-form-hdr">
          <svg class="auth-form-icon bearer" viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
          </svg>
          <span class="auth-form-title">{{ $t('auth.bearer.title') }}</span>
        </div>
        <div class="field-group">
          <label>{{ $t('auth.bearer.tokenLabel') }}</label>
          <n-input v-model:value="token" size="small" :placeholder="$t('auth.bearer.tokenPlaceholder')" />
          <div class="field-hint">
            <svg viewBox="0 0 24 24" width="11" height="11" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="10" />
              <line x1="12" y1="16" x2="12" y2="12" />
              <line x1="12" y1="8" x2="12.01" y2="8" />
            </svg>
            <span>{{ $t('auth.bearer.hint') }}</span>
          </div>
        </div>
      </div>

      <!-- Basic Auth -->
      <div v-else-if="authType === 'basic'" class="auth-form">
        <div class="auth-form-hdr">
          <svg class="auth-form-icon basic" viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2" />
            <circle cx="12" cy="7" r="4" />
          </svg>
          <span class="auth-form-title">{{ $t('auth.basic.title') }}</span>
        </div>
        <div class="field-group">
          <label>{{ $t('auth.basic.usernameLabel') }}</label>
          <n-input v-model:value="username" size="small" :placeholder="$t('auth.basic.usernamePlaceholder')" />
        </div>
        <div class="field-group">
          <label>{{ $t('auth.basic.passwordLabel') }}</label>
          <div class="pwd-wrap">
            <n-input
              v-model:value="password"
              size="small"
              :type="showPwd ? 'text' : 'password'"
              :placeholder="$t('auth.basic.passwordPlaceholder')"
            />
            <button class="pwd-toggle" @click="showPwd = !showPwd" :title="showPwd ? $t('auth.basic.hidePassword') : $t('auth.basic.showPassword')">
              <!-- Eye open -->
              <svg v-if="!showPwd" viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
                <circle cx="12" cy="12" r="3" />
              </svg>
              <!-- Eye off -->
              <svg v-else viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94M9.9 4.24A9.12 9.12 0 0112 4c7 0 11 8 11 8a18.5 18.5 0 01-2.16 3.19m-6.72-1.07a3 3 0 11-4.24-4.24" />
                <line x1="1" y1="1" x2="23" y2="23" />
              </svg>
            </button>
          </div>
          <div class="field-hint">
            <svg viewBox="0 0 24 24" width="11" height="11" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="10" />
              <line x1="12" y1="16" x2="12" y2="12" />
              <line x1="12" y1="8" x2="12.01" y2="8" />
            </svg>
            <span>{{ $t('auth.basic.hint') }}</span>
          </div>
        </div>
      </div>

      <!-- API Key -->
      <div v-else-if="authType === 'apikey'" class="auth-form">
        <div class="auth-form-hdr">
          <svg class="auth-form-icon apikey" viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 11-7.778 7.778 5.5 5.5 0 017.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4" />
          </svg>
          <span class="auth-form-title">{{ $t('auth.apikey.title') }}</span>
        </div>
        <div class="field-group">
          <label>{{ $t('auth.apikey.keyLabel') }}</label>
          <n-input v-model:value="apiKey" size="small" :placeholder="$t('auth.apikey.keyPlaceholder')" />
        </div>
        <div class="field-group">
          <label>{{ $t('auth.apikey.valueLabel') }}</label>
          <n-input v-model:value="apiValue" size="small" :placeholder="$t('auth.apikey.valuePlaceholder')" />
        </div>
        <div class="field-group">
          <label>{{ $t('auth.apikey.addToLabel') }}</label>
          <div class="addto-bar">
            <button :class="{ active: apiAddTo === 'header' }" class="addto-btn" @click="apiAddTo = 'header'">
              <svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" />
                <polyline points="14 2 14 8 20 8" />
              </svg>
              <span>{{ $t('auth.apikey.addToHeader') }}</span>
            </button>
            <button :class="{ active: apiAddTo === 'query' }" class="addto-btn" @click="apiAddTo = 'query'">
              <svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="11" cy="11" r="8" />
                <line x1="21" y1="21" x2="16.65" y2="16.65" />
              </svg>
              <span>{{ $t('auth.apikey.addToQuery') }}</span>
            </button>
          </div>
          <div class="field-hint">
            <svg viewBox="0 0 24 24" width="11" height="11" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="10" />
              <line x1="12" y1="16" x2="12" y2="12" />
              <line x1="12" y1="8" x2="12.01" y2="8" />
            </svg>
            <span v-if="apiAddTo === 'header'">
              {{ $t('auth.apikey.hintHeader', { keyName: apiKey || $t('auth.apikey.fallbackKeyName') }) }}
            </span>
            <span v-else>
              {{ $t('auth.apikey.hintQuery', { keyName: apiKey || $t('auth.apikey.fallbackKeyName') }) }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { NInput } from 'naive-ui'

const { t } = useI18n()

const props = defineProps<{
  modelValue: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', v: string): void
}>()

const authTypeOptions = computed(() => [
  {
    label: t('auth.type.none'), value: 'none',
    desc: t('auth.type.noneDesc'),
    icon: '<svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="4.93" y1="4.93" x2="19.07" y2="19.07"/></svg>',
  },
  {
    label: t('auth.type.bearer'), value: 'bearer',
    desc: t('auth.type.bearerDesc'),
    icon: '<svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>',
  },
  {
    label: t('auth.type.basic'), value: 'basic',
    desc: t('auth.type.basicDesc'),
    icon: '<svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>',
  },
  {
    label: t('auth.type.apikey'), value: 'apikey',
    desc: t('auth.type.apikeyDesc'),
    icon: '<svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 11-7.778 7.778 5.5 5.5 0 017.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/></svg>',
  },
])

const authType = ref('none')
const token = ref('')
const username = ref('')
const password = ref('')
const apiKey = ref('')
const apiValue = ref('')
const apiAddTo = ref<'header' | 'query'>('header')
const showPwd = ref(false)

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
  } catch { authType.value = 'none' }
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
  padding: 10px;
}

/* ══════════════════════════════════════════
   Auth type cards
   ══════════════════════════════════════════ */
.auth-cards {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 6px;
  margin-bottom: 12px;
}
.auth-card {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 10px;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  background: var(--bg-surface);
  cursor: pointer;
  transition: all var(--transition-fast);
  text-align: left;
  position: relative;
}
.auth-card:hover:not(.active) {
  border-color: var(--border-hover);
  background: var(--bg-hover);
}
.auth-card.active {
  border-color: var(--accent);
  background: var(--accent-soft);
}
.auth-card-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: var(--radius-sm);
  background: var(--bg-elevated);
  color: var(--text-muted);
  flex-shrink: 0;
  transition: all var(--transition-fast);
}
.auth-card.active .auth-card-icon {
  background: var(--accent-soft);
  color: var(--accent);
}
.auth-card-info {
  display: flex;
  flex-direction: column;
  gap: 1px;
  min-width: 0;
  flex: 1;
}
.auth-card-label {
  font-size: var(--fs-xs);
  font-weight: 600;
  font-family: var(--font-ui);
  color: var(--text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.auth-card-desc {
  font-size: var(--fs-2xs);
  font-family: var(--font-ui);
  color: var(--text-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.auth-card-check {
  position: absolute;
  top: 6px;
  right: 6px;
  color: var(--accent);
  display: flex;
  align-items: center;
}

/* ══════════════════════════════════════════
   Auth form wrap
   ══════════════════════════════════════════ */
.auth-form-wrap {
  min-height: 80px;
}

/* ── Empty state ── */
.auth-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 28px 16px;
  text-align: center;
}
.auth-empty-icon {
  color: var(--text-muted);
  opacity: 0.25;
  margin-bottom: 4px;
}
.auth-empty-text {
  font-size: var(--fs-sm);
  font-weight: 600;
  color: var(--text-secondary);
  font-family: var(--font-ui);
}
.auth-empty-hint {
  font-size: var(--fs-xs);
  color: var(--text-secondary);
  font-family: var(--font-ui);
}

/* ── Auth form ── */
.auth-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
  background: var(--bg-surface);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  padding: 14px;
}
.auth-form-hdr {
  display: flex;
  align-items: center;
  gap: 8px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--border-subtle);
}
.auth-form-icon {
  color: var(--text-muted);
}
.auth-form-icon.bearer  { color: var(--accent); }
.auth-form-icon.basic   { color: var(--blue); }
.auth-form-icon.apikey  { color: var(--amber); }
.auth-form-title {
  font-size: var(--fs-sm);
  font-weight: 700;
  color: var(--text-primary);
  font-family: var(--font-ui);
}

/* ── Field groups ── */
.field-group {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.field-group label {
  font-size: var(--fs-xs);
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  font-family: var(--font-ui);
  font-weight: 600;
}
.field-hint {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: var(--fs-2xs);
  color: var(--text-muted);
  font-family: var(--font-ui);
  padding: 2px 0;
}
.field-hint svg {
  flex-shrink: 0;
  opacity: 0.6;
}
.field-hint code {
  background: var(--bg-elevated);
  padding: 1px 4px;
  border-radius: var(--radius-xs);
  font-family: var(--font-mono);
  font-size: var(--fs-2xs);
  border: 1px solid var(--border-primary);
  color: var(--text-secondary);
  white-space: nowrap;
}

/* ── Password toggle ── */
.pwd-wrap {
  position: relative;
  display: flex;
}
.pwd-wrap :deep(.n-input) {
  padding-right: 32px;
}
.pwd-toggle {
  position: absolute;
  right: 4px;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  border: none;
  background: transparent;
  color: var(--text-muted);
  cursor: pointer;
  border-radius: var(--radius-xs);
  transition: all var(--transition-fast);
  z-index: 1;
}
.pwd-toggle:hover {
  color: var(--text-secondary);
  background: var(--bg-hover);
}

/* ── Add-to segmented buttons ── */
.addto-bar {
  display: flex;
  gap: 0;
  background: var(--bg-elevated);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-sm);
  padding: 2px;
  width: fit-content;
}
.addto-btn {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 4px 10px;
  font-size: var(--fs-xs);
  font-family: var(--font-ui);
  font-weight: 500;
  color: var(--text-muted);
  background: transparent;
  border: none;
  border-radius: var(--radius-xs);
  cursor: pointer;
  transition: all var(--transition-fast);
  white-space: nowrap;
}
.addto-btn:hover:not(.active) {
  color: var(--text-secondary);
  background: var(--bg-hover);
}
.addto-btn.active {
  color: var(--text-primary);
  background: var(--bg-surface);
  box-shadow: 0 1px 3px rgba(0,0,0,0.12), 0 0 0 1px var(--border-primary);
  font-weight: 600;
}
</style>
