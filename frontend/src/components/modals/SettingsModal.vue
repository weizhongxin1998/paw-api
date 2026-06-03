<template>
  <n-modal
    :show="show"
    preset="card"
    title="设置"
    style="width: 680px"
    :mask-closable="false"
    :segmented="{ footer: true }"
    @update:show="onClose"
  >
    <div class="settings-body">
      <div class="settings-nav">
        <div
          v-for="tab in navTabs"
          :key="tab.key"
          class="nav-item"
          :class="{ active: activeNav === tab.key }"
          @click="activeNav = tab.key"
        >{{ tab.label }}</div>
      </div>
      <div class="settings-content">
        <div v-if="activeNav === 'general'">
          <n-form label-placement="left" label-width="140px">
            <n-form-item label="请求超时 (s)">
              <n-input-number v-model:value="timeout" :min="1" :max="300" style="width: 90px" />
            </n-form-item>
            <n-form-item label="跟随重定向">
              <n-switch v-model:value="followRedirects" />
            </n-form-item>
            <n-form-item label="最大重定向次数">
              <n-input-number v-model:value="maxRedirects" :min="0" :max="20" style="width: 90px" />
            </n-form-item>
            <n-form-item label="SSL 证书验证">
              <n-switch v-model:value="sslVerify" />
            </n-form-item>
          </n-form>
        </div>
        <div v-else-if="activeNav === 'proxy'" class="settings-muted">
          <span class="muted-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><circle cx="12" cy="12" r="10"/><path d="M2 12h20"/><path d="M12 2a15.3 15.3 0 014 10 15.3 15.3 0 01-4 10 15.3 15.3 0 01-4-10 15.3 15.3 0 014-10z"/></svg>
          </span>
          <span>尚未实现代理功能</span>
        </div>
        <div v-else-if="activeNav === 'cert'" class="settings-muted">
          <span class="muted-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0110 0v4"/></svg>
          </span>
          <span>尚未配置证书管理</span>
        </div>
        <div v-else-if="activeNav === 'appearance'" class="settings-appearance">
          <n-form label-placement="left" label-width="140px">
            <n-form-item label="主题">
              <n-select v-model:value="theme" :options="themeOptions" style="width: 140px" />
            </n-form-item>
            <n-form-item label="主色">
              <n-select v-model:value="accentColor" :options="accentOptions" style="width: 140px" />
            </n-form-item>
          </n-form>
          <div class="setting-group">
            <div class="setting-group-hdr">字体</div>
            <n-form label-placement="left" label-width="140px">
              <n-form-item label="字体大小">
                <div class="font-size-ctrl">
                  <n-button size="tiny" @click="changeFontSize(-1)" :disabled="fontSize <= 12">&minus;</n-button>
                  <span class="fs-val">{{ fontSize }}px</span>
                  <n-button size="tiny" @click="changeFontSize(+1)" :disabled="fontSize >= 16">+</n-button>
                </div>
              </n-form-item>
              <n-form-item label="字体族">
                <n-select v-model:value="fontFamily" :options="fontOptions" style="width: 200px" @update:value="onFontFamilyChange" />
              </n-form-item>
            </n-form>
            <div class="font-sample" :style="{ fontFamily: sampleFontStack, fontSize: fontSize + 'px' }">
              ABC abc 0123 &lt;tag&gt; {json}
            </div>
          </div>
        </div>
        <div v-else-if="activeNav === 'data'" class="settings-data">
          <n-button block @click="onBackup">
            <template #icon>
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
            </template>
            备份数据
          </n-button>
          <n-button block style="margin-top: 8px" @click="onRestore">
            <template #icon>
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="1 4 1 10 7 10"/><path d="M3.51 15a9 9 0 102.13-9.36L1 10"/></svg>
            </template>
            恢复数据
          </n-button>
        </div>
      </div>
    </div>
    <template #footer>
      <n-button @click="onClose">取消</n-button>
      <n-button type="primary" @click="onSave">保存</n-button>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { NModal, NButton, NForm, NFormItem, NInputNumber, NSelect, NSwitch } from 'naive-ui'
import { useSettingsStore } from '../../stores/settings'

defineProps<{ show: boolean }>()
const emit = defineEmits<{ 'update:show': [value: boolean] }>()
const settingsStore = useSettingsStore()

const navTabs = [
  { key: 'general', label: '通用' },
  { key: 'proxy', label: '代理' },
  { key: 'cert', label: '证书' },
  { key: 'appearance', label: '外观' },
  { key: 'data', label: '数据' },
]

const themeOptions = [
  { label: '夜间', value: 'dark' },
  { label: '日间', value: 'light' },
]

const accentOptions = [
  { label: '绿色', value: 'green' },
  { label: '蓝色', value: 'blue' },
  { label: '紫色', value: 'purple' },
]

const activeNav = ref('general')
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

function changeFontSize(delta: number) {
  const v = fontSize.value + delta
  if (v >= 12 && v <= 16) {
    fontSize.value = v
    document.documentElement.style.fontSize = v + 'px'
  }
}

function onFontFamilyChange() {
  const fam = `'${fontFamily.value}', 'Cascadia Code', 'Fira Code', 'SF Mono', 'Consolas', monospace`
  document.documentElement.style.fontFamily = fam
}

function onClose() { emit('update:show', false) }
function onSave() {
  const s = settingsStore.settings
  s.timeout = timeout.value
  s.followRedirects = followRedirects.value
  s.maxRedirects = maxRedirects.value
  s.sslVerify = sslVerify.value
  s.theme = theme.value
  s.accentColor = accentColor.value
  s.fontSize = fontSize.value
  s.fontFamily = fontFamily.value
  onClose()
}

function onBackup() {}
function onRestore() {}
</script>

<style scoped>
.settings-body { display: flex; min-height: 320px; }
.settings-nav { width: 90px; border-right: 1px solid var(--border-primary); padding: 4px 0; flex-shrink: 0; }
.nav-item { padding: 8px 14px; font-size: var(--fs-sm); cursor: pointer; color: var(--text-muted); border-right: 2px solid transparent; transition: all var(--transition); font-family: var(--font-family); }
.nav-item:hover { color: var(--text-secondary); background: var(--bg-hover); }
.nav-item.active { color: var(--accent); font-weight: 600; border-right-color: var(--accent); background: var(--accent-soft); }
.settings-content { flex: 1; padding: 14px 18px; overflow-y: auto; }
.setting-group { border: 1px solid var(--border-primary); border-radius: var(--radius-lg); padding: 14px; margin-top: 6px; background: var(--bg-elevated); }
.setting-group-hdr { font-size: var(--fs-xs); color: var(--text-muted); margin-bottom: 10px; text-transform: uppercase; letter-spacing: 0.5px; font-weight: 600; }
.settings-muted { display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 10px; height: 100%; color: var(--text-muted); font-size: var(--fs-sm); font-family: var(--font-family); }
.muted-icon { opacity: 0.3; }
.font-size-ctrl { display: flex; align-items: center; gap: 8px; }
.fs-val { font-size: var(--fs-md); font-weight: 600; color: var(--accent); min-width: 36px; text-align: center; font-family: var(--font-family); }
.font-sample { margin-top: 10px; padding: 10px 12px; background: var(--bg-base); border: 1px solid var(--border-primary); border-radius: var(--radius); color: var(--text-primary); white-space: nowrap; overflow: hidden; }
.settings-appearance :deep(.n-form-item) { margin-bottom: 0; }
.settings-appearance :deep(.n-form-item:last-child) { margin-bottom: 0; }
.settings-data { padding-top: 4px; }
</style>
