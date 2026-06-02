<template>
  <div v-if="show" class="modal-overlay" @click.self="onClose">
    <div class="settings-box">
      <div class="settings-hdr">
        <h3>设置</h3>
        <button class="close-btn" @click="onClose">&times;</button>
      </div>
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
          <div v-if="activeNav === 'general'" class="settings-section">
            <div class="setting-row">
              <span>请求超时 (s)</span>
              <input v-model.number="timeout" type="number" min="1" max="300" class="setting-input-num" />
            </div>
            <div class="setting-row">
              <span>跟随重定向</span>
              <input type="checkbox" v-model="followRedirects" />
            </div>
            <div class="setting-row">
              <span>最大重定向次数</span>
              <input v-model.number="maxRedirects" type="number" min="0" max="20" class="setting-input-num" />
            </div>
            <div class="setting-row">
              <span>SSL 证书验证</span>
              <input type="checkbox" v-model="sslVerify" />
            </div>
          </div>
          <div v-else-if="activeNav === 'proxy'" class="settings-section">
            <div class="setting-muted">不使用代理</div>
          </div>
          <div v-else-if="activeNav === 'cert'" class="settings-section">
            <div class="setting-muted">未配置</div>
          </div>
          <div v-else-if="activeNav === 'appearance'" class="settings-section">
            <div class="setting-row">
              <span>主题</span>
              <select v-model="theme" class="setting-select">
                <option value="dark">夜间</option>
                <option value="light">日间</option>
              </select>
            </div>
            <div class="setting-row">
              <span>主色</span>
              <select v-model="accentColor" class="setting-select">
                <option value="green">绿色</option>
                <option value="blue">蓝色</option>
                <option value="purple">紫色</option>
              </select>
            </div>
            <div class="setting-row">
              <span>字号</span>
              <select v-model="fontSize" class="setting-select">
                <option :value="12">12px</option>
                <option :value="13">13px</option>
                <option :value="14">14px</option>
              </select>
            </div>
          </div>
          <div v-else-if="activeNav === 'data'" class="settings-section">
            <button class="data-btn" @click="onBackup">备份</button>
            <button class="data-btn" @click="onRestore">恢复</button>
          </div>
        </div>
      </div>
      <div class="settings-footer">
        <button class="btn-cancel" @click="onClose">取消</button>
        <button class="btn-save" @click="onSave">保存</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
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

const activeNav = ref('general')
const timeout = ref(settingsStore.settings.timeout ?? 30)
const followRedirects = ref(settingsStore.settings.followRedirects ?? true)
const maxRedirects = ref(settingsStore.settings.maxRedirects ?? 10)
const sslVerify = ref(settingsStore.settings.sslVerify ?? true)
const theme = ref(settingsStore.settings.theme ?? 'dark')
const accentColor = ref(settingsStore.settings.accentColor ?? 'green')
const fontSize = ref(settingsStore.settings.fontSize ?? 13)

function onClose() { emit('update:show', false) }
function onSave() {
  settingsStore.settings = { ...settingsStore.settings, timeout: timeout.value, followRedirects: followRedirects.value, maxRedirects: maxRedirects.value, sslVerify: sslVerify.value, theme: theme.value, accentColor: accentColor.value, fontSize: fontSize.value }
  onClose()
}
function onBackup() {}
function onRestore() {}
</script>

<style scoped>
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.6); display: flex; align-items: center; justify-content: center; z-index: 200; }
.settings-box { background: var(--bg-surface); border: 1px solid var(--border-primary); border-radius: var(--radius-lg); width: 660px; max-height: 80vh; box-shadow: 0 12px 40px rgba(0,0,0,0.5); overflow: hidden; display: flex; flex-direction: column; }
.settings-hdr { display: flex; align-items: center; justify-content: space-between; padding: 14px 18px; border-bottom: 1px solid var(--border-primary); }
.settings-hdr h3 { margin: 0; font-size: 14px; font-weight: 600; color: var(--text-primary); }
.close-btn { background: none; border: none; font-size: 20px; color: var(--text-muted); cursor: pointer; padding: 0 4px; line-height: 1; }
.close-btn:hover { color: var(--text-primary); }
.settings-body { display: flex; flex: 1; min-height: 280px; overflow: hidden; }
.settings-nav { width: 90px; border-right: 1px solid var(--border-primary); padding: 4px 0; flex-shrink: 0; }
.nav-item { padding: 7px 14px; font-size: 12px; cursor: pointer; color: var(--text-muted); border-right: 2px solid transparent; transition: all var(--transition); font-family: var(--font-mono); }
.nav-item:hover { color: var(--text-secondary); background: var(--bg-hover); }
.nav-item.active { color: var(--accent); font-weight: 600; border-right-color: var(--accent); background: var(--accent-soft); }
.settings-content { flex: 1; padding: 14px 18px; overflow-y: auto; }
.settings-section { display: flex; flex-direction: column; gap: 2px; }
.setting-row { display: flex; align-items: center; justify-content: space-between; padding: 3px 0; font-size: 12px; }
.setting-row span { color: var(--text-secondary); }
.setting-input-num { width: 70px; padding: 4px 6px; border: 1px solid var(--border-primary); border-radius: var(--radius-sm); text-align: right; font-size: 12px; outline: none; background: var(--bg-base); color: var(--text-primary); font-family: var(--font-mono); }
.setting-input-num:focus { border-color: var(--accent); }
.setting-select { padding: 4px 6px; border: 1px solid var(--border-primary); border-radius: var(--radius-sm); font-size: 12px; outline: none; cursor: pointer; background: var(--bg-base); color: var(--text-secondary); font-family: var(--font-mono); }
.setting-muted { color: var(--text-muted); font-size: 12px; padding: 10px 0; }
.data-btn { display: block; width: 100%; text-align: left; padding: 8px 12px; border: 1px solid var(--border-primary); border-radius: var(--radius-sm); background: var(--bg-base); cursor: pointer; font-size: 12px; margin-bottom: 8px; color: var(--text-secondary); font-family: var(--font-mono); transition: all var(--transition); }
.data-btn:hover { background: var(--bg-elevated); border-color: var(--border-hover); }
.settings-footer { display: flex; justify-content: flex-end; gap: 8px; padding: 12px 18px; border-top: 1px solid var(--border-primary); }
.btn-cancel { padding: 5px 16px; border: 1px solid var(--border-primary); border-radius: var(--radius); background: var(--bg-base); font-size: 11px; cursor: pointer; color: var(--text-secondary); font-family: var(--font-mono); transition: all var(--transition); }
.btn-cancel:hover { border-color: var(--border-hover); color: var(--text-primary); }
.btn-save { padding: 5px 16px; background: var(--accent); color: #000; border: 1px solid var(--accent); border-radius: var(--radius); font-size: 11px; cursor: pointer; font-weight: 600; font-family: var(--font-mono); }
.btn-save:hover { background: var(--accent-hover); }
</style>
