<template>
  <n-config-provider :theme-overrides="themeOverrides" :theme="nTheme">
    <n-global-style />
    <n-dialog-provider>
      <n-message-provider>
        <div class="app-container" :class="{ 'theme-light': themeMode === 'light' }">
          <div class="noise-overlay"></div>
          <div class="vignette-overlay"></div>
          <AppHeader
            v-if="projectStore.currentId"
            :project-id="projectStore.currentId"
            :theme-mode="themeMode"
            @project-changed="onProjectChanged"
            @back-to-home="projectStore.currentId = null"
            @toggle-theme="toggleTheme"
          />
          <Transition name="view-fade" mode="out-in">
            <ProjectHome
              v-if="!projectStore.currentId"
              :key="'home'"
              :theme-mode="themeMode"
              @enter-project="onEnterProject"
              @toggle-theme="toggleTheme"
            />
            <AppBody v-else :key="'app'" :project-id="projectStore.currentId" />
          </Transition>
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
  /* ── Surface hierarchy ── */
  --bg-base: #0a0a0b;
  --bg-surface: #111113;
  --bg-elevated: #18181b;
  --bg-hover: #1e1e22;
  --bg-active: #26262b;
  --bg-inset: #0e0e10;

  /* ── Borders ── */
  --border-primary: #27272a;
  --border-hover: #3f3f46;
  --border-focus: #00e05a;
  --border-subtle: rgba(255,255,255,0.04);

  /* ── Typography ── */
  --text-primary: #e4e4e7;
  --text-secondary: #a1a1aa;
  --text-muted: #71717a;
  --text-placeholder: #52525b;
  --text-inverse: #0a0a0b;

  /* ── Accent — neon mint ── */
  --accent: #00e05a;
  --accent-hover: #00ff66;
  --accent-pressed: #00b84a;
  --accent-soft: rgba(0,224,90,0.07);
  --accent-glow: rgba(0,224,90,0.14);
  --accent-glow-strong: rgba(0,224,90,0.28);

  /* ── Semantic ── */
  --red: #ef4444;
  --red-hover: #f87171;
  --red-soft: rgba(239,68,68,0.1);
  --amber: #f59e0b;
  --amber-soft: rgba(245,158,11,0.1);
  --blue: #3b82f6;
  --blue-soft: rgba(59,130,246,0.1);
  --purple: #a855f7;
  --purple-soft: rgba(168,85,247,0.1);
  --cyan: #06b6d4;
  --cyan-soft: rgba(6,182,212,0.1);

  /* ── HTTP Methods ── */
  --method-get: #3b82f6;
  --method-post: #22c55e;
  --method-put: #f59e0b;
  --method-patch: #a855f7;
  --method-delete: #ef4444;
  --method-head: #06b6d4;
  --method-options: #71717a;

  /* ── Radii ── */
  --radius-xs: 3px;
  --radius-sm: 5px;
  --radius: 8px;
  --radius-lg: 12px;
  --radius-xl: 16px;

  /* ── Shadows ── */
  --shadow-sm: 0 1px 2px rgba(0,0,0,0.3);
  --shadow: 0 2px 8px rgba(0,0,0,0.35);
  --shadow-lg: 0 8px 24px rgba(0,0,0,0.45);
  --shadow-glow: 0 0 20px var(--accent-glow);

  /* ── Motion ── */
  --ease-out: cubic-bezier(0.16, 1, 0.3, 1);
  --ease-spring: cubic-bezier(0.34, 1.56, 0.64, 1);
  --transition-fast: 0.1s var(--ease-out);
  --transition: 0.18s var(--ease-out);
  --transition-slow: 0.3s var(--ease-out);

  /* ── Font sizes ── */
  --fs-2xs: 9px;
  --fs-xs: 10.5px;
  --fs-sm: 11.5px;
  --fs-base: 13px;
  --fs-md: 14px;
  --fs-lg: 16px;
  --fs-xl: 18px;
  --fs-2xl: 24px;
  --fs-3xl: 32px;

  /* ── Font families ── */
  --font-family: 'JetBrains Mono', 'Cascadia Code', 'Fira Code', 'SF Mono', 'Consolas', monospace;
  --font-mono: 'JetBrains Mono', 'Cascadia Code', 'Fira Code', 'SF Mono', 'Consolas', monospace;
  --font-ui: -apple-system, BlinkMacSystemFont, 'Segoe UI', system-ui, sans-serif;
}

/* ════════════════════════════════════════════════════════════
   Light Theme
   ════════════════════════════════════════════════════════════ */
.theme-light {
  --bg-base: #f8f8f7;
  --bg-surface: #ffffff;
  --bg-elevated: #f4f4f3;
  --bg-hover: #ededeb;
  --bg-active: #e4e4e2;
  --bg-inset: #f0f0ee;
  --border-primary: #d4d4d1;
  --border-hover: #a8a8a4;
  --border-focus: #009944;
  --border-subtle: rgba(0,0,0,0.04);
  --text-primary: #18181b;
  --text-secondary: #3f3f46;
  --text-muted: #71717a;
  --text-placeholder: #a1a1aa;
  --text-inverse: #fafafa;
  --accent: #059669;
  --accent-hover: #047857;
  --accent-pressed: #065f46;
  --accent-soft: rgba(5,150,105,0.08);
  --accent-glow: rgba(5,150,105,0.12);
  --accent-glow-strong: rgba(5,150,105,0.2);
  --red: #dc2626;
  --red-hover: #ef4444;
  --red-soft: rgba(220,38,38,0.08);
  --amber: #d97706;
  --amber-soft: rgba(217,119,6,0.08);
  --blue: #2563eb;
  --blue-soft: rgba(37,99,235,0.08);
  --purple: #7c3aed;
  --purple-soft: rgba(124,58,237,0.08);
  --cyan: #0891b2;
  --cyan-soft: rgba(8,145,178,0.08);
  --shadow-sm: 0 1px 2px rgba(0,0,0,0.05);
  --shadow: 0 2px 8px rgba(0,0,0,0.08);
  --shadow-lg: 0 8px 24px rgba(0,0,0,0.12);
  --shadow-glow: 0 0 16px var(--accent-glow);
}

/* ════════════════════════════════════════════════════════════
   Base Reset
   ════════════════════════════════════════════════════════════ */
html, body, #app {
  margin: 0; padding: 0; height: 100%;
  font-family: var(--font-family); font-size: var(--fs-base);
  color: var(--text-primary); background: var(--bg-base);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  overflow: hidden;
}

.app-container {
  display: flex; flex-direction: column; height: 100vh;
  overflow: hidden; position: relative; background: var(--bg-base);
}

/* ════════════════════════════════════════════════════════════
   Texture Overlays
   ════════════════════════════════════════════════════════════ */
.noise-overlay {
  position: fixed; inset: 0; pointer-events: none; z-index: 9999;
  opacity: 0.025;
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.7' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)' opacity='1'/%3E%3C/svg%3E");
  background-size: 200px 200px;
  mix-blend-mode: overlay;
}
.theme-light .noise-overlay { opacity: 0.012; mix-blend-mode: multiply; }

.vignette-overlay {
  position: fixed; inset: 0; pointer-events: none; z-index: 9998;
  background: radial-gradient(ellipse 120% 100% at 50% 50%, transparent 55%, rgba(0,0,0,0.3) 100%);
}
.theme-light .vignette-overlay {
  background: radial-gradient(ellipse 120% 100% at 50% 50%, transparent 60%, rgba(0,0,0,0.04) 100%);
}

/* ════════════════════════════════════════════════════════════
   Transitions
   ════════════════════════════════════════════════════════════ */
.view-fade-enter-active, .view-fade-leave-active {
  transition: opacity 0.25s var(--ease-out), transform 0.25s var(--ease-out);
}
.view-fade-enter-from { opacity: 0; transform: translateY(8px) scale(0.995); }
.view-fade-leave-to { opacity: 0; transform: translateY(-4px) scale(1); }

/* ════════════════════════════════════════════════════════════
   Scrollbars
   ════════════════════════════════════════════════════════════ */
::-webkit-scrollbar { width: 7px; height: 7px; }
::-webkit-scrollbar-track { background: transparent; }
::-webkit-scrollbar-thumb {
  background: rgba(255,255,255,0.08); border-radius: 4px;
  border: 2px solid transparent; background-clip: padding-box;
  transition: background 0.15s;
}
::-webkit-scrollbar-thumb:hover { background: rgba(255,255,255,0.16); background-clip: padding-box; border: 2px solid transparent; }
::-webkit-scrollbar-corner { background: transparent; }
.theme-light ::-webkit-scrollbar-thumb { background: rgba(0,0,0,0.1); }
.theme-light ::-webkit-scrollbar-thumb:hover { background: rgba(0,0,0,0.2); }

/* ════════════════════════════════════════════════════════════
   Selection & Focus
   ════════════════════════════════════════════════════════════ */
::selection { background: var(--accent-soft); color: var(--accent); }
:focus-visible { outline: 1.5px solid var(--accent); outline-offset: 1px; border-radius: var(--radius-xs); }

/* ════════════════════════════════════════════════════════════
   Utility Classes (shared across components)
   ════════════════════════════════════════════════════════════ */

/* HTTP Method Badges */
.method-get    { color: var(--method-get) !important; }
.method-post   { color: var(--method-post) !important; }
.method-put    { color: var(--method-put) !important; }
.method-patch  { color: var(--method-patch) !important; }
.method-delete { color: var(--method-delete) !important; }
.method-head   { color: var(--method-head) !important; }
.method-options{ color: var(--method-options) !important; }

.method-pill {
  display: inline-flex; align-items: center; justify-content: center;
  font-size: var(--fs-2xs); font-weight: 700; letter-spacing: 0.04em;
  padding: 1px 6px; border-radius: var(--radius-xs);
  line-height: 1.6; text-transform: uppercase; font-family: var(--font-mono);
}
.method-pill.method-get    { background: var(--blue-soft); color: var(--method-get); }
.method-pill.method-post   { background: rgba(34,197,94,0.1); color: var(--method-post); }
.method-pill.method-put    { background: var(--amber-soft); color: var(--method-put); }
.method-pill.method-patch  { background: var(--purple-soft); color: var(--method-patch); }
.method-pill.method-delete { background: var(--red-soft); color: var(--method-delete); }
.method-pill.method-head   { background: var(--cyan-soft); color: var(--method-head); }
.method-pill.method-options{ background: rgba(113,113,122,0.1); color: var(--method-options); }

/* Status code pills */
.status-pill {
  display: inline-flex; align-items: center;
  font-size: var(--fs-xs); font-weight: 600; letter-spacing: 0.02em;
  padding: 2px 8px; border-radius: var(--radius-sm); font-family: var(--font-mono);
}
.status-pill.status-2xx { background: rgba(34,197,94,0.1); color: #22c55e; }
.status-pill.status-3xx { background: var(--blue-soft); color: var(--blue); }
.status-pill.status-4xx { background: var(--amber-soft); color: var(--amber); }
.status-pill.status-5xx { background: var(--red-soft); color: var(--red); }

/* ════════════════════════════════════════════════════════════
   Keyframes
   ════════════════════════════════════════════════════════════ */
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(8px); }
  to   { opacity: 1; transform: translateY(0); }
}
@keyframes fadeInScale {
  from { opacity: 0; transform: scale(0.96); }
  to   { opacity: 1; transform: scale(1); }
}
@keyframes slideUp {
  from { opacity: 0; transform: translateY(12px); }
  to   { opacity: 1; transform: translateY(0); }
}
@keyframes pulseGlow {
  0%, 100% { box-shadow: 0 0 0 0 var(--accent-glow); }
  50%      { box-shadow: 0 0 12px 2px var(--accent-glow); }
}
@keyframes shimmer {
  0%   { background-position: -200% 0; }
  100% { background-position: 200% 0; }
}
</style>
