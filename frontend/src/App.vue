<template>
  <n-config-provider :theme-overrides="themeOverrides" :theme="nTheme">
    <n-dialog-provider>
      <n-message-provider>
        <div class="app-container" :class="{ 'theme-light': themeMode === 'light' }">
          <div class="noise-overlay"></div>
          <AppHeader
            v-if="projectStore.currentId"
            :project-id="projectStore.currentId"
            :theme-mode="themeMode"
            @project-changed="onProjectChanged"
            @back-to-home="projectStore.currentId = null"
            @toggle-theme="toggleTheme"
          />
          <ProjectHome
            v-if="!projectStore.currentId"
            :theme-mode="themeMode"
            @enter-project="onEnterProject"
            @toggle-theme="toggleTheme"
          />
          <AppBody v-else :project-id="projectStore.currentId" />
        </div>
      </n-message-provider>
    </n-dialog-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { darkTheme, type GlobalThemeOverrides } from 'naive-ui'
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
const nTheme = computed(() => themeMode.value === 'dark' ? darkTheme : null)

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

function toggleTheme() {
  themeMode.value = themeMode.value === 'dark' ? 'light' : 'dark'
}

onMounted(async () => {
  await projectStore.loadProjects()
})

const themeOverrides = computed<GlobalThemeOverrides>(() => {
  const dark = themeMode.value === 'dark'
  return {
    common: {
      primaryColor: dark ? '#00e05a' : '#009944',
      primaryColorHover: dark ? '#00ff66' : '#007a33',
      primaryColorPressed: dark ? '#00b84a' : '#006b2a',
      primaryColorSuppl: dark ? '#00e05a' : '#009944',
      bodyColor: dark ? '#0d0d0d' : '#f5f3f1',
      cardColor: dark ? '#141414' : '#ffffff',
      modalColor: dark ? '#141414' : '#ffffff',
      popoverColor: dark ? '#1a1a1a' : '#f2f2ef',
      borderColor: dark ? '#2a2a2a' : '#c8c8c2',
      dividerColor: dark ? '#2a2a2a' : '#c8c8c2',
      borderRadius: '4px',
      fontSizeSmall: '12px',
      fontSizeMedium: '13px',
      fontSizeLarge: '15px',
      fontFamilyMono: "'JetBrains Mono','Cascadia Code','Fira Code','SF Mono','Consolas',monospace",
      fontFamily: "'JetBrains Mono','Cascadia Code','Fira Code','SF Mono','Consolas',monospace",
      fontWeightStrong: '600',
      textColor1: dark ? '#e0e0e0' : '#1a1a18',
      textColor2: dark ? '#b0b0b0' : '#444442',
      textColor3: dark ? '#707070' : '#666662',
      placeholderColor: dark ? '#505050' : '#888884',
      inputColor: dark ? '#141414' : '#ffffff',
      inputColorFocus: dark ? '#141414' : '#ffffff',
      inputBorderColor: dark ? '#2a2a2a' : '#c8c8c2',
      inputBorderColorHover: dark ? '#3a3a3a' : '#a8a8a2',
      inputBorderColorFocus: dark ? '#00e05a' : '#009944',
      buttonColor2: dark ? '#1a1a1a' : '#fafaf9',
      buttonColor2Hover: dark ? '#252525' : '#f0f0ee',
      buttonColor2Pressed: dark ? '#303030' : '#e8e8e6',
      scrollbarColor: dark ? '#333' : '#ccc',
      scrollbarColorHover: dark ? '#444' : '#aaa',
      closeColor: dark ? '#707070' : '#888',
      closeColorHover: dark ? '#e0e0e0' : '#333',
      closeColorPressed: dark ? '#b0b0b0' : '#555',
    },
    Button: {
      textColor: dark ? '#00e05a' : '#009944',
      textColorHover: dark ? '#00ff66' : '#007a33',
      textColorPressed: dark ? '#00b84a' : '#006b2a',
      border: dark ? '1px solid #2a2a2a' : '1px solid #c8c8c2',
      borderHover: dark ? '1px solid #3a3a3a' : '1px solid #a8a8a2',
      borderFocus: dark ? '1px solid #00e05a' : '1px solid #009944',
      color: dark ? '#1a1a1a' : '#fafaf9',
      colorHover: dark ? '#252525' : '#f0f0ee',
      colorPressed: dark ? '#303030' : '#e8e8e6',
      borderRadiusSmall: '3px',
      borderRadiusMedium: '4px',
    },
    Input: {
      border: dark ? '1px solid #2a2a2a' : '1px solid #c8c8c2',
      borderHover: dark ? '1px solid #3a3a3a' : '1px solid #a8a8a2',
      borderFocus: dark ? '1px solid #00e05a' : '1px solid #009944',
      borderRadius: '4px',
      color: dark ? '#141414' : '#ffffff',
      colorFocus: dark ? '#141414' : '#ffffff',
      textColor: dark ? '#e0e0e0' : '#1a1a18',
      placeholderColor: dark ? '#505050' : '#999',
      lineHeight: '1.6',
      fontSizeSmall: '12px',
      fontSizeMedium: '13px',
    },
    Select: {
      peers: { InternalSelection: { textColor: dark ? '#e0e0e0' : '#1a1a18' } },
    },
    Checkbox: {
      colorChecked: dark ? '#00e05a' : '#009944',
      borderChecked: dark ? '#00e05a' : '#009944',
      border: dark ? '1px solid #3a3a3a' : '1px solid #bbb',
      checkMarkColor: dark ? '#0d0d0d' : '#fff',
    },
    Dropdown: {
      color: dark ? '#1a1a1a' : '#ffffff',
      dividerColor: dark ? '#2a2a2a' : '#e0e0da',
      optionColorActive: dark ? '#1a2a1a' : '#e6f7ec',
      optionTextColorActive: dark ? '#00e05a' : '#009944',
    },
    Modal: {
      color: dark ? '#141414' : '#ffffff',
      textColor: dark ? '#e0e0e0' : '#1a1a18',
      titleTextColor: dark ? '#e0e0e0' : '#1a1a18',
    },
    Tabs: {
      tabTextColorActiveLine: dark ? '#00e05a' : '#009944',
      tabTextColorActiveBar: dark ? '#00e05a' : '#009944',
      barColor: dark ? '#00e05a' : '#009944',
    },
    Spin: { color: dark ? '#00e05a' : '#009944' },
    Result: {
      titleTextColor: dark ? '#e0e0e0' : '#1a1a18',
      textColor: dark ? '#b0b0b0' : '#555552',
    },
  }
})

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
  --bg-base: #f5f3f1;
  --bg-surface: #ffffff;
  --bg-elevated: #f2f2ef;
  --bg-hover: #ebebe8;
  --bg-active: #e3e3e0;
  --border-primary: #c8c8c2;
  --border-hover: #a8a8a2;
  --border-focus: #009944;
  --text-primary: #1a1a18;
  --text-secondary: #444442;
  --text-muted: #666662;
  --text-placeholder: #888884;
  --accent: #009944;
  --accent-hover: #007a33;
  --accent-pressed: #006b2a;
  --accent-soft: rgba(0,153,68,0.08);
  --accent-glow: rgba(0,153,68,0.12);
  --red: #cc3333;
  --red-soft: rgba(204,51,51,0.08);
  --amber: #b36d00;
  --amber-soft: rgba(179,109,0,0.08);
  --blue: #2266cc;
  --blue-soft: rgba(34,102,204,0.08);
  --purple: #7733cc;
  --purple-soft: rgba(119,51,204,0.08);
}

html, body, #app {
  margin: 0; padding: 0; height: 100%;
  font-family: var(--font-ui); font-size: 13px;
  color: var(--text-primary); background: var(--bg-base);
  -webkit-font-smoothing: antialiased; overflow: hidden;
}

.app-container {
  display: flex; flex-direction: column; height: 100vh;
  overflow: hidden; position: relative; background: var(--bg-base);
}

.noise-overlay {
  position: fixed; inset: 0; pointer-events: none; z-index: 9999;
  opacity: 0.025;
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)' opacity='1'/%3E%3C/svg%3E");
  background-size: 256px 256px;
}
.theme-light .noise-overlay { opacity: 0.012; }

::-webkit-scrollbar { width: 6px; height: 6px; }
::-webkit-scrollbar-track { background: transparent; }
::-webkit-scrollbar-thumb { background: #333; border-radius: 3px; }
::-webkit-scrollbar-thumb:hover { background: #444; }
.theme-light ::-webkit-scrollbar-thumb { background: #ccc; }
.theme-light ::-webkit-scrollbar-thumb:hover { background: #aaa; }

::selection { background: var(--accent-soft); color: var(--accent); }
</style>
