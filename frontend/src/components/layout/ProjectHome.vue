<template>
  <div class="project-home">
    <div class="home-bg-grid"></div>
    <div class="home-header">
      <div class="home-brand">
        <svg class="brand-mark" width="32" height="32" viewBox="0 0 32 32" fill="none" stroke="currentColor" stroke-width="1.8">
          <rect x="2" y="5" width="28" height="22" rx="3" />
          <line x1="2" y1="12" x2="30" y2="12" />
          <circle cx="7" cy="8.5" r="1.2" fill="currentColor" stroke="none" />
          <circle cx="11" cy="8.5" r="1.2" fill="currentColor" stroke="none" />
          <circle cx="15" cy="8.5" r="1.2" fill="currentColor" stroke="none" />
          <path d="M10 19l3 3 6-6" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
        <h1 class="home-title">PAW API</h1>
        <span class="home-version">v1.0</span>
      </div>
      <div class="home-actions">
        <button class="btn-theme" @click="emit('toggle-theme')" :title="props.themeMode === 'dark' ? '日间模式' : '夜间模式'">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle v-if="props.themeMode === 'dark'" cx="12" cy="12" r="5"/><line v-if="props.themeMode === 'dark'" x1="12" y1="1" x2="12" y2="3"/><line v-if="props.themeMode === 'dark'" x1="12" y1="21" x2="12" y2="23"/><line v-if="props.themeMode === 'dark'" x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line v-if="props.themeMode === 'dark'" x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line v-if="props.themeMode === 'dark'" x1="1" y1="12" x2="3" y2="12"/><line v-if="props.themeMode === 'dark'" x1="21" y1="12" x2="23" y2="12"/><line v-if="props.themeMode === 'dark'" x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line v-if="props.themeMode === 'dark'" x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>
            <path v-else d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
          </svg>
        </button>
        <n-button size="tiny" @click="showImport = true" class="btn-import">
          <template #icon>
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
          </template>
          导入
        </n-button>
        <n-button type="primary" size="tiny" @click="showCreate = true">+ 新建项目</n-button>
      </div>
    </div>

    <div class="project-grid">
      <div
        v-for="(p, idx) in projectList"
        :key="p.id"
        class="project-card"
        :style="{ '--card-delay': idx * 0.04 + 's' }"
        @click="enterProject(p.id)"
        @keydown.enter="enterProject(p.id)"
        tabindex="0"
      >
        <div class="card-glow"></div>
        <div class="card-header">
          <span class="card-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M22 19a2 2 0 01-2 2H4a2 2 0 01-2-2V5a2 2 0 012-2h5l2 3h9a2 2 0 012 2z"/></svg>
          </span>
          <span class="card-arrow">&rarr;</span>
        </div>
        <div class="card-body">
          <div class="card-name">{{ p.name }}</div>
          <div class="card-desc" v-if="p.description">{{ p.description }}</div>
        </div>
        <div class="card-stats">
          <span class="stat">
            <span class="stat-num">{{ p.stats?.request_count ?? 0 }}</span>
            <span class="stat-label">接口</span>
          </span>
          <span class="stat-divider"></span>
          <span class="stat">
            <span class="stat-num">{{ p.stats?.collection_count ?? 0 }}</span>
            <span class="stat-label">集合</span>
          </span>
        </div>
      </div>

      <div v-if="projectList.length === 0" class="empty-state">
        <div class="empty-graphic">
          <svg width="64" height="64" viewBox="0 0 64 64" fill="none" stroke="currentColor" stroke-width="1" opacity="0.3">
            <rect x="8" y="14" width="48" height="36" rx="4" />
            <line x1="8" y1="24" x2="56" y2="24" />
            <circle cx="16" cy="19" r="2" fill="currentColor" stroke="none" />
            <circle cx="22" cy="19" r="2" fill="currentColor" stroke="none" />
            <circle cx="28" cy="19" r="2" fill="currentColor" stroke="none" />
            <path d="M24 34l4 4 8-8" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
          <div class="empty-ring"></div>
        </div>
        <h2 class="empty-title">还没有项目</h2>
        <p class="empty-desc">创建第一个项目，开始调试 API</p>
        <n-button type="primary" @click="showCreate = true">+ 新建项目</n-button>
      </div>
    </div>

    <n-modal v-model:show="showCreate" preset="card" title="新建项目" style="width: 400px" :mask-closable="false">
      <n-form label-placement="top">
        <n-form-item label="名称">
          <n-input v-model:value="newName" placeholder="项目名称" @keydown.enter="onCreate" />
        </n-form-item>
        <n-form-item label="描述">
          <n-input v-model:value="newDesc" placeholder="项目描述（可选）" @keydown.enter="onCreate" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button @click="showCreate = false">取消</n-button>
        <n-button type="primary" :disabled="!newName.trim()" @click="onCreate">创建</n-button>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NButton, NModal, NForm, NFormItem, NInput } from 'naive-ui'
import { GetProjectStats, CreateProject, ListProjects } from '../../../wailsjs/go/main/App'

interface ProjectCard {
  id: number
  name: string
  description: string
  stats: { request_count: number; collection_count: number } | null
}

const props = defineProps<{ themeMode: 'dark' | 'light' }>()

const emit = defineEmits<{
  (e: 'enter-project', id: number): void
  (e: 'toggle-theme'): void
}>()

const projectList = ref<ProjectCard[]>([])
const showCreate = ref(false)
const showImport = ref(false)
const newName = ref('')
const newDesc = ref('')

async function loadProjects() {
  try {
    const projects = await ListProjects()
    const items: ProjectCard[] = []
    for (const p of projects || []) {
      let stats = null
      try { stats = await GetProjectStats(p.id) } catch { /* ignore */ }
      items.push({ id: p.id, name: p.name, description: p.description, stats })
    }
    projectList.value = items
  } catch { projectList.value = [] }
}

async function onCreate() {
  const name = newName.value.trim()
  if (!name) return
  const p = await CreateProject(name, newDesc.value.trim())
  showCreate.value = false
  newName.value = ''
  newDesc.value = ''
  await loadProjects()
  enterProject(p.id)
}

function enterProject(id: number) { emit('enter-project', id) }

onMounted(loadProjects)
</script>

<style scoped>
.project-home {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--bg-base);
  position: relative;
  overflow: hidden;
}
.home-bg-grid {
  position: absolute;
  inset: 0;
  pointer-events: none;
  background-image:
    linear-gradient(var(--border-primary) 0.5px, transparent 0.5px),
    linear-gradient(90deg, var(--border-primary) 0.5px, transparent 0.5px);
  background-size: 40px 40px;
  opacity: 0.3;
  mask-image: radial-gradient(ellipse 80% 80% at 50% 0%, black 30%, transparent 70%);
}
.home-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 28px 40px 0;
  position: relative;
  animation: fadeIn 0.3s ease both;
}
.home-brand {
  display: flex;
  align-items: center;
  gap: 10px;
}
.brand-mark { color: var(--accent); filter: drop-shadow(0 0 6px var(--accent-glow)); }
.home-title {
  font-size: var(--fs-2xl);
  font-weight: 700;
  color: var(--text-primary);
  letter-spacing: 4px;
  margin: 0;
  font-family: var(--font-mono);
}
.home-version {
  font-size: var(--fs-2xs);
  color: var(--text-muted);
  background: var(--bg-elevated);
  padding: 2px 6px;
  border-radius: 3px;
  border: 1px solid var(--border-primary);
  font-family: var(--font-mono);
}
.home-actions { display: flex; gap: 8px; align-items: center; }
.btn-theme {
  padding: 7px 10px;
  background: transparent;
  color: var(--text-muted);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all var(--transition);
}
.btn-theme:hover { border-color: var(--accent); color: var(--accent); background: var(--accent-soft); }
.project-grid {
  flex: 1;
  display: flex;
  flex-wrap: wrap;
  align-content: flex-start;
  gap: 16px;
  padding: 32px 40px;
  overflow-y: auto;
  position: relative;
}
.project-card {
  width: 220px;
  min-height: 150px;
  background: var(--bg-surface);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg);
  cursor: pointer;
  display: flex;
  flex-direction: column;
  padding: 16px;
  position: relative;
  transition: all var(--transition-slow);
  animation: cardIn 0.35s ease both;
  animation-delay: var(--card-delay, 0s);
  outline: none;
  overflow: hidden;
}
.project-card:focus-visible { border-color: var(--accent); box-shadow: 0 0 0 2px var(--accent-glow); }
.project-card:hover {
  border-color: var(--accent);
  box-shadow: 0 0 0 1px var(--accent), 0 4px 20px var(--accent-glow);
  transform: translateY(-2px);
}
.card-glow {
  position: absolute;
  inset: 0;
  pointer-events: none;
  background: radial-gradient(ellipse 120% 80% at 50% -20%, var(--accent-glow-strong) 0%, transparent 60%);
  opacity: 0;
  transition: opacity var(--transition-slow);
}
.project-card:hover .card-glow { opacity: 1; }

@keyframes cardIn {
  from { opacity: 0; transform: translateY(12px) scale(0.97); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
  position: relative;
  z-index: 1;
}
.card-icon { color: var(--accent); opacity: 0.6; }
.card-arrow {
  font-size: var(--fs-md);
  color: var(--text-muted);
  opacity: 0;
  transform: translateX(-4px);
  transition: all var(--transition);
}
.project-card:hover .card-arrow { opacity: 0.5; transform: translateX(0); }
.card-body { flex: 1; position: relative; z-index: 1; }
.card-name {
  font-size: var(--fs-base);
  font-weight: 600;
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-family: var(--font-mono);
}
.card-desc {
  font-size: var(--fs-sm);
  color: var(--text-muted);
  margin-top: 3px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.card-stats {
  display: flex;
  align-items: center;
  gap: 8px;
  padding-top: 10px;
  border-top: 1px solid var(--border-primary);
  position: relative;
  z-index: 1;
}
.stat { display: flex; align-items: baseline; gap: 3px; }
.stat-num { font-size: var(--fs-lg); font-weight: 700; color: var(--accent); font-family: var(--font-mono); }
.stat-label { font-size: var(--fs-xs); color: var(--text-muted); text-transform: uppercase; letter-spacing: 0.5px; }
.stat-divider { width: 1px; height: 16px; background: var(--border-primary); }
.empty-state {
  width: 100%;
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  animation: fadeIn 0.4s ease both;
}
.empty-graphic { position: relative; margin-bottom: 4px; }
.empty-graphic svg { color: var(--text-muted); display: block; }
.empty-ring {
  position: absolute; inset: -8px;
  border: 1px solid var(--border-primary);
  border-radius: 50%;
  animation: spin 8s linear infinite;
}
@keyframes spin {
  to { transform: rotate(360deg); }
}
.empty-title { font-size: var(--fs-xl); color: var(--text-secondary); margin: 0; font-weight: 600; font-family: var(--font-mono); }
.empty-desc { font-size: var(--fs-base); color: var(--text-muted); margin: 0 0 4px; }
</style>
