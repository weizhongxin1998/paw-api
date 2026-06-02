<template>
  <n-config-provider :theme-overrides="themeOverrides" :theme="nTheme">
    <n-dialog-provider>
      <n-message-provider>
        <div class="app-container" :class="{ 'theme-light': themeMode === 'light' }">
          <div class="noise-overlay"></div>
          <AppHeader
            v-if="projectStore.currentId"
            :project-id="projectStore.currentId"
            @project-changed="onProjectChanged"
            @back-to-home="projectStore.currentId = null"
          />
          <ProjectHome
            v-if="!projectStore.currentId"
            @enter-project="onEnterProject"
          />
          <AppBody v-else :project-id="projectStore.currentId" />
        </div>
      </n-message-provider>
    </n-dialog-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { darkTheme } from 'naive-ui'
import { NConfigProvider, NMessageProvider, NDialogProvider } from 'naive-ui'
import AppHeader from './components/layout/AppHeader.vue'
import AppBody from './components/layout/AppBody.vue'
import ProjectHome from './components/layout/ProjectHome.vue'
import { useProjectStore } from './stores/project'
import { useCollectionStore } from './stores/collection'
import { useEnvStore } from './stores/env'

const projectStore = useProjectStore()
const collectionStore = useCollectionStore()
const envStore = useEnvStore()

const themeMode = ref<'dark' | 'light'>('dark')
const nTheme = ref(darkTheme)

async function onEnterProject(id: number) {
  await projectStore.switchProject(id)
  collectionStore.loadTree(id)
  envStore.loadEnvironments(id)
}

async function onProjectChanged(id: number) {
  await projectStore.switchProject(id)
  collectionStore.loadTree(id)
  envStore.loadEnvironments(id)
}

onMounted(async () => {
  await projectStore.loadProjects()
})

function toggleTheme() {
  themeMode.value = themeMode.value === 'dark' ? 'light' : 'dark'
  nTheme.value = themeMode.value === 'dark' ? darkTheme : null as any
}

const themeOverrides = {
  common: {
    primaryColor: '#00e05a',
    primaryColorHover: '#00ff66',
    primaryColorPressed: '#00b84a',
    primaryColorSuppl: '#00e05a',
    bodyColor: '#0d0d0d',
    cardColor: '#141414',
    modalColor: '#141414',
    popoverColor: '#1a1a1a',
    borderColor: '#2a2a2a',
    dividerColor: '#2a2a2a',
    borderRadius: '4px',
    fontSizeSmall: '12px',
    fontSizeMedium: '13px',
    fontSizeLarge: '15px',
    fontFamilyMono: "'JetBrains Mono','Cascadia Code','Fira Code','SF Mono','Consolas',monospace",
    fontFamily: "'JetBrains Mono','Cascadia Code','Fira Code','SF Mono','Consolas',monospace",
    fontWeightStrong: '600',
    textColor1: '#e0e0e0',
    textColor2: '#b0b0b0',
    textColor3: '#707070',
    placeholderColor: '#505050',
    inputColor: '#141414',
    inputBorderColor: '#2a2a2a',
    inputBorderColorHover: '#3a3a3a',
    inputBorderColorFocus: '#00e05a',
    buttonColor2: '#1a1a1a',
    buttonColor2Hover: '#252525',
    buttonColor2Pressed: '#303030',
    scrollbarColor: '#333',
    scrollbarColorHover: '#444',
    closeColor: '#707070',
    closeColorHover: '#e0e0e0',
    closeColorPressed: '#b0b0b0',
  },
  Button: {
    textColor: '#00e05a',
    textColorHover: '#00ff66',
    textColorPressed: '#00b84a',
    border: '1px solid #2a2a2a',
    borderHover: '1px solid #3a3a3a',
    borderFocus: '1px solid #00e05a',
    color: '#1a1a1a',
    colorHover: '#252525',
    colorPressed: '#303030',
    borderRadiusSmall: '3px',
    borderRadiusMedium: '4px',
  },
  Input: {
    border: '1px solid #2a2a2a',
    borderHover: '1px solid #3a3a3a',
    borderFocus: '1px solid #00e05a',
    borderRadius: '4px',
    color: '#141414',
    colorFocus: '#141414',
    textColor: '#e0e0e0',
    placeholderColor: '#505050',
    lineHeight: '1.6',
    fontSizeSmall: '12px',
    fontSizeMedium: '13px',
  },
  Select: {
    peers: { InternalSelection: { textColor: '#e0e0e0' } },
  },
  Checkbox: {
    colorChecked: '#00e05a',
    borderChecked: '#00e05a',
    border: '1px solid #3a3a3a',
    checkMarkColor: '#0d0d0d',
  },
  Dropdown: {
    color: '#1a1a1a',
    dividerColor: '#2a2a2a',
    optionColorActive: '#1a2a1a',
    optionTextColorActive: '#00e05a',
  },
  Modal: {
    color: '#141414',
    textColor: '#e0e0e0',
    titleTextColor: '#e0e0e0',
  },
  Tabs: {
    tabTextColorActiveLine: '#00e05a',
    tabTextColorActiveBar: '#00e05a',
    barColor: '#00e05a',
  },
  Spin: {
    color: '#00e05a',
  },
  Result: {
    titleTextColor: '#e0e0e0',
    textColor: '#b0b0b0',
  },
}

defineExpose({ toggleTheme, themeMode })
</script>

<style>
:root {
  --bg-base: #0d0d0d;
  --bg-surface: #141414;
  --bg-elevated: #1a1a1a;
  --bg-hover: #1f1f1f;
  --bg-active: #252525;
  --border-primary: #2a2a2a;
  --border-hover: #3a3a3a;
  --border-focus: #00e05a;
  --border-active: #3a3a3a;
  --text-primary: #e0e0e0;
  --text-secondary: #b0b0b0;
  --text-muted: #707070;
  --text-placeholder: #505050;
  --accent: #00e05a;
  --accent-hover: #00ff66;
  --accent-pressed: #00b84a;
  --accent-soft: rgba(0,224,90,0.08);
  --accent-glow: rgba(0,224,90,0.15);
  --red: #ff4444;
  --red-hover: #ff6666;
  --red-soft: rgba(255,68,68,0.08);
  --amber: #ffaa00;
  --amber-soft: rgba(255,170,0,0.08);
  --blue: #4499ff;
  --blue-soft: rgba(68,153,255,0.08);
  --purple: #aa66ff;
  --purple-soft: rgba(170,102,255,0.08);
  --radius-sm: 3px;
  --radius: 4px;
  --radius-lg: 6px;
  --font-mono: 'JetBrains Mono', 'Cascadia Code', 'Fira Code', 'SF Mono', 'Consolas', monospace;
  --font-ui: 'JetBrains Mono', 'Cascadia Code', 'Fira Code', 'SF Mono', 'Consolas', monospace;
  --transition: 0.12s ease;
  --transition-slow: 0.2s ease;
}

.theme-light {
  --bg-base: #f5f5f4;
  --bg-surface: #ffffff;
  --bg-elevated: #fafaf9;
  --bg-hover: #f0f0ee;
  --bg-active: #e8e8e6;
  --border-primary: #e0e0dc;
  --border-hover: #d0d0cc;
  --border-focus: #009944;
  --border-active: #d0d0cc;
  --text-primary: #1a1a18;
  --text-secondary: #555552;
  --text-muted: #888884;
  --text-placeholder: #aaaaa8;
  --accent: #009944;
  --accent-hover: #007a33;
  --accent-pressed: #006b2a;
  --accent-soft: rgba(0,153,68,0.06);
  --accent-glow: rgba(0,153,68,0.12);
  --red: #cc3333;
  --red-soft: rgba(204,51,51,0.06);
  --amber: #cc7700;
  --amber-soft: rgba(204,119,0,0.06);
  --blue: #2266cc;
  --blue-soft: rgba(34,102,204,0.06);
  --purple: #7733cc;
  --purple-soft: rgba(119,51,204,0.06);
}

html, body, #app {
  margin: 0;
  padding: 0;
  height: 100%;
  font-family: var(--font-ui);
  font-size: 13px;
  color: var(--text-primary);
  background: var(--bg-base);
  -webkit-font-smoothing: antialiased;
  overflow: hidden;
}

.app-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
  position: relative;
  background: var(--bg-base);
}

.noise-overlay {
  position: fixed;
  inset: 0;
  pointer-events: none;
  z-index: 9999;
  opacity: 0.025;
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)' opacity='1'/%3E%3C/svg%3E");
  background-size: 256px 256px;
}

::-webkit-scrollbar { width: 6px; height: 6px; }
::-webkit-scrollbar-track { background: transparent; }
::-webkit-scrollbar-thumb { background: #333; border-radius: 3px; }
::-webkit-scrollbar-thumb:hover { background: #444; }
.theme-light ::-webkit-scrollbar-thumb { background: #ccc; }
.theme-light ::-webkit-scrollbar-thumb:hover { background: #aaa; }

::selection { background: var(--accent-soft); color: var(--accent); }
</style>
