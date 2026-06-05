<template>
  <n-modal
    :show="show"
    preset="card"
    title="设置"
    style="width: 720px"
    :mask-closable="false"
    :segmented="{ footer: true }"
    @update:show="onClose"
  >
    <div class="settings-body">
      <!-- Sidebar navigation with icons -->
      <div class="settings-nav">
        <div
          v-for="tab in navTabs"
          :key="tab.key"
          class="nav-item"
          :class="{ active: activeNav === tab.key }"
          @click="activeNav = tab.key"
        >
          <span class="nav-icon" v-html="tab.icon"></span>
          <span class="nav-label">{{ tab.label }}</span>
        </div>
        <div class="nav-spacer"></div>
        <div class="nav-item nav-reset" @click="onRestoreDefaults">
          <span class="nav-icon">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><polyline points="1 4 1 10 7 10"/><path d="M3.51 15a9 9 0 102.13-9.36L1 10"/></svg>
          </span>
          <span class="nav-label">恢复默认</span>
        </div>
      </div>

      <!-- Content area -->
      <div class="settings-content">
        <!-- General -->
        <div v-if="activeNav === 'general'" class="section-panel">
          <div class="section-hdr">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-2 2 2 2 0 01-2-2v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83 0 2 2 0 010-2.83l.06-.06A1.65 1.65 0 004.68 15a1.65 1.65 0 00-1.51-1H3a2 2 0 01-2-2 2 2 0 012-2h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 010-2.83 2 2 0 012.83 0l.06.06A1.65 1.65 0 009 4.68a1.65 1.65 0 001-1.51V3a2 2 0 012-2 2 2 0 012 2v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 0 2 2 0 010 2.83l-.06.06A1.65 1.65 0 0019.32 9a1.65 1.65 0 001.51 1H21a2 2 0 012 2 2 2 0 01-2 2h-.09a1.65 1.65 0 00-1.51 1z"/></svg>
            <span>通用设置</span>
          </div>
          <n-form label-placement="left" label-width="140px">
            <n-form-item label="请求超时 (s)">
              <n-input-number v-model:value="timeout" :min="1" :max="300" style="width: 100px" />
            </n-form-item>
            <n-form-item label="跟随重定向">
              <n-switch v-model:value="followRedirects" />
            </n-form-item>
            <n-form-item label="最大重定向次数">
              <n-input-number v-model:value="maxRedirects" :min="0" :max="20" style="width: 100px" :disabled="!followRedirects" />
            </n-form-item>
            <n-form-item label="SSL 证书验证">
              <n-switch v-model:value="sslVerify" />
              <template #feedback>
                <span class="form-hint">关闭后可能影响安全性，请谨慎使用</span>
              </template>
            </n-form-item>
          </n-form>
        </div>

        <!-- Proxy (placeholder) -->
        <div v-else-if="activeNav === 'proxy'" class="settings-muted">
          <span class="muted-icon">
            <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.2"><circle cx="12" cy="12" r="10"/><path d="M2 12h20"/><path d="M12 2a15.3 15.3 0 014 10 15.3 15.3 0 01-4 10 15.3 15.3 0 01-4-10 15.3 15.3 0 014-10z"/></svg>
          </span>
          <span>尚未实现代理功能</span>
          <span class="muted-sub">将在后续版本中支持 HTTP/SOCKS 代理</span>
        </div>

        <!-- Certificates (placeholder) -->
        <div v-else-if="activeNav === 'cert'" class="settings-muted">
          <span class="muted-icon">
            <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.2"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0110 0v4"/></svg>
          </span>
          <span>尚未配置证书管理</span>
          <span class="muted-sub">将在后续版本中支持自定义 SSL 证书</span>
        </div>

        <!-- Appearance -->
        <div v-else-if="activeNav === 'appearance'" class="section-panel">
          <div class="section-hdr">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/></svg>
            <span>外观设置</span>
          </div>

          <!-- Theme visual cards -->
          <div class="theme-section">
            <label class="setting-label">主题</label>
            <div class="theme-cards">
              <div
                class="theme-card"
                :class="{ selected: theme === 'dark' }"
                @click="theme = 'dark'"
              >
                <div class="theme-preview theme-dark-preview">
                  <div class="tp-titlebar"></div>
                  <div class="tp-body">
                    <div class="tp-sidebar"></div>
                    <div class="tp-main">
                      <div class="tp-line"></div>
                      <div class="tp-line short"></div>
                    </div>
                  </div>
                </div>
                <span class="theme-name">夜间</span>
              </div>
              <div
                class="theme-card"
                :class="{ selected: theme === 'light' }"
                @click="theme = 'light'"
              >
                <div class="theme-preview theme-light-preview">
                  <div class="tp-titlebar"></div>
                  <div class="tp-body">
                    <div class="tp-sidebar"></div>
                    <div class="tp-main">
                      <div class="tp-line"></div>
                      <div class="tp-line short"></div>
                    </div>
                  </div>
                </div>
                <span class="theme-name">日间</span>
              </div>
            </div>
          </div>

          <!-- Accent color -->
          <n-form label-placement="left" label-width="140px" class="accent-form">
            <n-form-item label="主色">
              <div class="accent-swatches">
                <div
                  v-for="opt in accentOptions"
                  :key="opt.value"
                  class="accent-swatch"
                  :class="[opt.value, { selected: accentColor === opt.value }]"
                  @click="accentColor = opt.value"
                  :title="opt.label"
                >
                  <svg v-if="accentColor === opt.value" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"><polyline points="20 6 9 17 4 12"/></svg>
                </div>
              </div>
            </n-form-item>
          </n-form>

          <!-- Font section with live preview -->
          <div class="setting-group">
            <div class="setting-group-hdr">字体</div>
            <n-form label-placement="left" label-width="140px">
              <n-form-item label="字体大小">
                <div class="font-slider-wrap">
                  <n-slider
                    v-model:value="fontSize"
                    :min="12"
                    :max="16"
                    :step="1"
                    :format-tooltip="(v: number) => v + 'px'"
                    style="flex: 1"
                    @update:value="onFontSizeLive"
                  />
                  <span class="fs-val">{{ fontSize }}px</span>
                </div>
              </n-form-item>
              <n-form-item label="字体族">
                <n-select v-model:value="fontFamily" :options="fontOptions" style="width: 220px" @update:value="onFontFamilyChange" />
              </n-form-item>
            </n-form>
            <div class="font-sample" :style="{ fontFamily: sampleFontStack, fontSize: fontSize + 'px' }">
              <div class="font-sample-line">ABC abc 0123 &lt;tag&gt; {json}</div>
              <div class="font-sample-line secondary">GET /api/users?page=1&amp;limit=20 HTTP/1.1</div>
            </div>
          </div>
        </div>

        <!-- Data -->
        <div v-else-if="activeNav === 'data'" class="section-panel">
          <div class="section-hdr">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"/><path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"/></svg>
            <span>数据管理</span>
          </div>
          <div class="data-actions">
            <n-button block @click="onBackup" class="data-btn">
              <template #icon>
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
              </template>
              <div class="data-btn-content">
                <span class="data-btn-title">备份数据</span>
                <span class="data-btn-desc">将所有项目、请求和环境导出为备份文件</span>
              </div>
            </n-button>
            <n-button block @click="onRestore" class="data-btn">
              <template #icon>
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="1 4 1 10 7 10"/><path d="M3.51 15a9 9 0 102.13-9.36L1 10"/></svg>
              </template>
              <div class="data-btn-content">
                <span class="data-btn-title">恢复数据</span>
                <span class="data-btn-desc">从备份文件恢复所有数据</span>
              </div>
            </n-button>
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="settings-footer">
        <span class="autosave-hint" v-if="hasChanges">
          <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="20 6 9 17 4 12"/></svg>
          已自动保存
        </span>
        <span class="autosave-hint" v-else>修改将自动保存</span>
        <n-button @click="onClose">关闭</n-button>
      </div>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import {
  NModal, NButton, NForm, NFormItem, NInputNumber,
  NSelect, NSwitch, NSlider, useMessage,
} from 'naive-ui'
import { useSettingsStore, applySettingsToDOM } from '../../stores/settings'

const DEFAULTS = {
  timeout: 30,
  followRedirects: true,
  maxRedirects: 10,
  sslVerify: true,
  theme: 'dark',
  accentColor: 'green',
  fontSize: 13,
  fontFamily: 'JetBrains Mono',
}

defineProps<{ show: boolean }>()
const emit = defineEmits<{ 'update:show': [value: boolean] }>()
const settingsStore = useSettingsStore()
const message = useMessage()

const navTabs = [
  { key: 'general', label: '通用', icon: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 11-2.83 2.83l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 11-4 0v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 11-2.83-2.83l.06-.06A1.65 1.65 0 004.68 15a1.65 1.65 0 00-1.51-1H3a2 2 0 110-4h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 112.83-2.83l.06.06A1.65 1.65 0 009 4.68a1.65 1.65 0 001-1.51V3a2 2 0 114 0v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 112.83 2.83l-.06.06A1.65 1.65 0 0019.32 9a1.65 1.65 0 001.51 1H21a2 2 0 110 4h-.09a1.65 1.65 0 00-1.51 1z"/></svg>' },
  { key: 'proxy', label: '代理', icon: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><circle cx="12" cy="12" r="10"/><path d="M2 12h20"/><path d="M12 2a15.3 15.3 0 014 10 15.3 15.3 0 01-4 10 15.3 15.3 0 01-4-10 15.3 15.3 0 014-10z"/></svg>' },
  { key: 'cert', label: '证书', icon: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0110 0v4"/></svg>' },
  { key: 'appearance', label: '外观', icon: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/></svg>' },
  { key: 'data', label: '数据', icon: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"/><path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"/></svg>' },
]

const accentOptions = [
  { label: '绿色', value: 'green' },
  { label: '蓝色', value: 'blue' },
  { label: '紫色', value: 'purple' },
]

const activeNav = ref('general')
const hasChanges = ref(false)

// Local reactive state
const timeout = ref(settingsStore.settings.timeout ?? 30)
const followRedirects = ref(settingsStore.settings.followRedirects ?? true)
const maxRedirects = ref(settingsStore.settings.maxRedirects ?? 10)
const sslVerify = ref(settingsStore.settings.sslVerify ?? true)
const theme = ref(settingsStore.settings.theme ?? 'dark')
const accentColor = ref(settingsStore.settings.accentColor ?? 'green')
const fontSize = ref(settingsStore.settings.fontSize ?? 13)
const fontFamily = ref(settingsStore.settings.fontFamily ?? 'JetBrains Mono')

const fontOptions = settingsStore.FONT_FAMILIES

const sampleFontStack = computed(() =>
  `'${fontFamily.value}', 'Cascadia Code', 'Fira Code', 'SF Mono', 'Consolas', monospace`
)

// -- Auto-save: watch all settings and apply to store --
let saveTimer: ReturnType<typeof setTimeout> | null = null

function applyToStore() {
  const s = settingsStore.settings
  s.timeout = timeout.value
  s.followRedirects = followRedirects.value
  s.maxRedirects = maxRedirects.value
  s.sslVerify = sslVerify.value
  s.theme = theme.value
  s.accentColor = accentColor.value
  s.fontSize = fontSize.value
  s.fontFamily = fontFamily.value
  hasChanges.value = true
  // Auto-hide the "saved" indicator after 2s
  if (saveTimer) clearTimeout(saveTimer)
  saveTimer = setTimeout(() => { hasChanges.value = false }, 2000)
}

// Watch all local refs and auto-save
watch([timeout, followRedirects, maxRedirects, sslVerify, theme, accentColor, fontSize, fontFamily], () => {
  applyToStore()
})

// -- Live preview for font size slider --
function onFontSizeLive(val: number) {
  document.documentElement.style.fontSize = val + 'px'
}

function onFontFamilyChange() {
  const fam = `'${fontFamily.value}', 'Cascadia Code', 'Fira Code', 'SF Mono', 'Consolas', monospace`
  document.documentElement.style.fontFamily = fam
}

// -- Restore defaults --
function onRestoreDefaults() {
  timeout.value = DEFAULTS.timeout
  followRedirects.value = DEFAULTS.followRedirects
  maxRedirects.value = DEFAULTS.maxRedirects
  sslVerify.value = DEFAULTS.sslVerify
  theme.value = DEFAULTS.theme
  accentColor.value = DEFAULTS.accentColor
  fontSize.value = DEFAULTS.fontSize
  fontFamily.value = DEFAULTS.fontFamily
  applyToStore()
  onFontFamilyChange()
  message.success('已恢复默认设置')
}

// -- Navigation --
function onClose() { emit('update:show', false) }

function onBackup() { message.info('备份功能即将推出') }
function onRestore() { message.info('恢复功能即将推出') }
</script>

<style scoped>
.settings-body { display: flex; min-height: 360px; }

/* -- Sidebar Navigation with icons -- */
.settings-nav {
  width: 100px;
  border-right: 1px solid var(--border-primary);
  padding: 4px 0;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
}
.nav-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 9px 14px;
  font-size: var(--fs-sm);
  cursor: pointer;
  color: var(--text-secondary);
  border-right: 2px solid transparent;
  transition: all 0.15s ease;
  font-family: var(--font-family);
}
.nav-item:hover {
  color: var(--text-secondary);
  background: var(--bg-hover);
}
.nav-item.active {
  color: var(--accent);
  font-weight: 600;
  border-right-color: var(--accent);
  background: var(--accent-soft);
}
.nav-icon {
  display: flex;
  align-items: center;
  opacity: 0.7;
  flex-shrink: 0;
}
.nav-item.active .nav-icon { opacity: 1; }
.nav-label { white-space: nowrap; }
.nav-spacer { flex: 1; }
.nav-reset {
  border-top: 1px solid var(--border-primary);
  margin-top: 4px;
  color: var(--text-muted);
}
.nav-reset:hover { color: var(--amber, #f59e0b); }

/* -- Content Area -- */
.settings-content { flex: 1; padding: 14px 20px; overflow-y: auto; }

.section-hdr {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: var(--fs-md);
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--border-subtle);
}
.section-hdr svg { opacity: 0.6; }

/* -- Placeholder sections -- */
.settings-muted {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  height: 100%;
  color: var(--text-muted);
  font-size: var(--fs-sm);
  font-family: var(--font-family);
}
.muted-icon { opacity: 0.2; }
.muted-sub { font-size: var(--fs-xs); opacity: 0.6; }

/* -- Theme Visual Cards -- */
.theme-section { margin-bottom: 14px; }
.setting-label {
  display: block;
  font-size: var(--fs-sm);
  color: var(--text-secondary);
  margin-bottom: 8px;
}
.theme-cards {
  display: flex;
  gap: 12px;
}
.theme-card {
  cursor: pointer;
  border: 2px solid var(--border-primary);
  border-radius: var(--radius-lg, 10px);
  padding: 6px;
  transition: all 0.2s ease;
  width: 120px;
}
.theme-card:hover { border-color: var(--text-muted); }
.theme-card.selected {
  border-color: var(--accent);
  box-shadow: 0 0 0 1px var(--accent);
}
.theme-preview {
  border-radius: 6px;
  overflow: hidden;
  height: 64px;
}
/* Dark theme preview */
.theme-dark-preview {
  background: #111113;
}
.theme-dark-preview .tp-titlebar {
  height: 10px;
  background: #18181b;
  border-bottom: 1px solid #27272a;
}
.theme-dark-preview .tp-body { display: flex; height: calc(100% - 10px); }
.theme-dark-preview .tp-sidebar {
  width: 22px;
  background: #0a0a0b;
  border-right: 1px solid #27272a;
}
.theme-dark-preview .tp-main { flex: 1; padding: 5px; }
.theme-dark-preview .tp-line {
  height: 3px;
  background: #27272a;
  border-radius: 2px;
  margin-bottom: 3px;
}
.theme-dark-preview .tp-line.short { width: 60%; }

/* Light theme preview */
.theme-light-preview {
  background: #ffffff;
}
.theme-light-preview .tp-titlebar {
  height: 10px;
  background: #f4f4f3;
  border-bottom: 1px solid #d4d4d1;
}
.theme-light-preview .tp-body { display: flex; height: calc(100% - 10px); }
.theme-light-preview .tp-sidebar {
  width: 22px;
  background: #f8f8f7;
  border-right: 1px solid #d4d4d1;
}
.theme-light-preview .tp-main { flex: 1; padding: 5px; }
.theme-light-preview .tp-line {
  height: 3px;
  background: #e4e4e2;
  border-radius: 2px;
  margin-bottom: 3px;
}
.theme-light-preview .tp-line.short { width: 60%; }

.theme-name {
  display: block;
  text-align: center;
  font-size: var(--fs-xs);
  color: var(--text-secondary);
  margin-top: 4px;
}
.theme-card.selected .theme-name { color: var(--accent); font-weight: 600; }

/* -- Accent Color Swatches -- */
.accent-form { margin-bottom: 4px; }
.accent-swatches {
  display: flex;
  gap: 8px;
}
.accent-swatch {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  cursor: pointer;
  border: 2px solid transparent;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;
}
.accent-swatch.green  { background: #22c55e; }
.accent-swatch.blue   { background: #3b82f6; }
.accent-swatch.purple { background: #a855f7; }
.accent-swatch:hover { transform: scale(1.15); }
.accent-swatch.selected {
  border-color: var(--text-primary);
  box-shadow: 0 0 0 2px var(--bg-base);
}
.accent-swatch.selected svg { stroke: white; }

/* -- Font Section with live preview -- */
.setting-group {
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg);
  padding: 14px;
  margin-top: 10px;
  background: var(--bg-elevated);
}
.setting-group-hdr {
  font-size: var(--fs-xs);
  color: var(--text-secondary);
  margin-bottom: 10px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  font-weight: 600;
}
.font-slider-wrap {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
}
.fs-val {
  font-size: var(--fs-md);
  font-weight: 600;
  color: var(--accent);
  min-width: 36px;
  text-align: center;
  font-family: var(--font-family);
}
.font-sample {
  margin-top: 10px;
  padding: 10px 12px;
  background: var(--bg-base);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  color: var(--text-primary);
  overflow: hidden;
}
.font-sample-line { white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.font-sample-line.secondary {
  margin-top: 4px;
  opacity: 0.5;
  font-size: 0.85em;
}
.section-panel :deep(.n-form-item) { margin-bottom: 0; }
.form-hint { font-size: var(--fs-2xs, 9px); color: var(--text-secondary); }

/* -- Data section -- */
.data-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.data-btn {
  height: auto !important;
  padding: 12px 14px !important;
  justify-content: flex-start !important;
}
.data-btn-content {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  margin-left: 8px;
}
.data-btn-title {
  font-weight: 600;
  font-size: var(--fs-sm);
}
.data-btn-desc {
  font-size: var(--fs-xs);
  color: var(--text-secondary);
  margin-top: 2px;
}

/* -- Footer with auto-save indicator -- */
.settings-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 10px;
}
.autosave-hint {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: var(--fs-xs);
  color: var(--text-muted);
  margin-right: auto;
  opacity: 0.7;
}
.autosave-hint svg { color: var(--accent); }
</style>
