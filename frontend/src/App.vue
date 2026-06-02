<template>
  <n-config-provider :theme-overrides="themeOverrides" :theme="nTheme">
    <n-global-style />
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
import { ref, computed, onMounted, watch } from 'vue'
import { darkTheme, type GlobalThemeOverrides } from 'naive-ui'
import { NConfigProvider, NMessageProvider, NDialogProvider, NGlobalStyle } from 'naive-ui'
import AppHeader from './components/layout/AppHeader.vue'
import AppBody from './components/layout/AppBody.vue'
import ProjectHome from './components/layout/ProjectHome.vue'
import { useProjectStore } from './stores/project'
import { useCollectionStore } from './stores/collection'
import { useEnvStore } from './stores/env'
import { useSettingsStore, buildNaiveOverrides } from './stores/settings'

const projectStore = useProjectStore()
const collectionStore = useCollectionStore()
const envStore = useEnvStore()
const settingsStore = useSettingsStore()

const themeMode = ref<'dark' | 'light'>('dark')
const nTheme = computed(() => themeMode.value === 'dark' ? darkTheme : null)

const themeOverrides = ref<GlobalThemeOverrides>(
  buildNaiveOverrides(settingsStore.settings, true)
)

watch(
  () => [settingsStore.settings.fontSize, settingsStore.settings.fontFamily, themeMode.value] as const,
  () => {
    themeOverrides.value = buildNaiveOverrides(settingsStore.settings, themeMode.value === 'dark')
  }
)

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
  --transition: 0.12s ease;
  --transition-slow: 0.2s ease;

  /* Typography — set by settings store via JS */
  --fs-2xs: 9px;
  --fs-xs: 10px;
  --fs-sm: 11px;
  --fs-base: 13px;
  --fs-md: 14px;
  --fs-lg: 16px;
  --fs-xl: 18px;
  --fs-2xl: 22px;
  --font-family: 'JetBrains Mono', 'Cascadia Code', 'Fira Code', 'SF Mono', 'Consolas', monospace;
  --font-mono: 'JetBrains Mono', 'Cascadia Code', 'Fira Code', 'SF Mono', 'Consolas', monospace;
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
  font-family: var(--font-family); font-size: var(--fs-base);
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
