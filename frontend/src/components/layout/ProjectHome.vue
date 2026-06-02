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
        <button class="btn-import" @click="showImport = true">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
          导入
        </button>
        <button class="btn-new" @click="showCreate = true">+ 新建项目</button>
      </div>
    </div>

    <div class="project-grid">
      <div
        v-for="p in projectList"
        :key="p.id"
        class="project-card"
        @click="enterProject(p.id)"
        @keydown.enter="enterProject(p.id)"
      >
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
        <svg class="empty-icon" width="48" height="48" viewBox="0 0 48 48" fill="none" stroke="currentColor" stroke-width="1.2">
          <rect x="6" y="10" width="36" height="28" rx="3" />
          <line x1="6" y1="18" x2="42" y2="18" />
          <circle cx="13" cy="14" r="1.5" fill="currentColor" stroke="none" />
          <circle cx="18" cy="14" r="1.5" fill="currentColor" stroke="none" />
          <circle cx="23" cy="14" r="1.5" fill="currentColor" stroke="none" />
        </svg>
        <h2>还没有项目</h2>
        <p>创建第一个项目，开始调试 API</p>
        <button class="btn-new" @click="showCreate = true">+ 新建项目</button>
      </div>
    </div>

    <div v-if="showCreate" class="modal-overlay" @click.self="showCreate = false">
      <div class="modal-box">
        <h3>新建项目</h3>
        <label>名称</label>
        <input v-model="newName" placeholder="项目名称" @keydown.enter="onCreate" />
        <label>描述</label>
        <input v-model="newDesc" placeholder="描述（可选）" @keydown.enter="onCreate" />
        <div class="modal-acts">
          <button class="btn-cancel" @click="showCreate = false">取消</button>
          <button class="btn-save" @click="onCreate" :disabled="!newName.trim()">创建</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetProjectStats, CreateProject, ListProjects } from '../../../wailsjs/go/main/App'

interface ProjectCard {
  id: number
  name: string
  description: string
  stats: { request_count: number; collection_count: number } | null
}

const emit = defineEmits<{
  (e: 'enter-project', id: number): void
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
  opacity: 0.35;
  mask-image: radial-gradient(ellipse 80% 80% at 50% 0%, black 30%, transparent 70%);
}
.home-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 28px 40px 0;
  position: relative;
}
.home-brand {
  display: flex;
  align-items: center;
  gap: 10px;
}
.brand-mark { color: var(--accent); }
.home-title {
  font-size: 20px;
  font-weight: 700;
  color: var(--text-primary);
  letter-spacing: 3px;
  margin: 0;
  font-family: var(--font-mono);
}
.home-version {
  font-size: 10px;
  color: var(--text-muted);
  background: var(--bg-elevated);
  padding: 2px 6px;
  border-radius: 2px;
  border: 1px solid var(--border-primary);
  font-family: var(--font-mono);
}
.home-actions { display: flex; gap: 8px; }
.btn-new {
  padding: 7px 18px;
  background: var(--accent);
  color: #000;
  border: none;
  border-radius: var(--radius);
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  font-family: var(--font-mono);
  transition: all var(--transition);
}
.btn-new:hover { background: var(--accent-hover); }
.btn-import {
  padding: 7px 14px;
  background: transparent;
  color: var(--text-secondary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  font-size: 12px;
  cursor: pointer;
  font-family: var(--font-mono);
  display: flex;
  align-items: center;
  gap: 5px;
  transition: all var(--transition);
}
.btn-import:hover { border-color: var(--border-hover); color: var(--text-primary); background: var(--bg-elevated); }
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
  border-radius: var(--radius);
  cursor: pointer;
  display: flex;
  flex-direction: column;
  padding: 16px;
  position: relative;
  transition: all var(--transition-slow);
}
.project-card:hover {
  border-color: var(--accent);
  box-shadow: 0 0 0 1px var(--accent), 0 4px 20px var(--accent-glow);
  transform: translateY(-1px);
}
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}
.card-icon { color: var(--accent); opacity: 0.6; }
.card-arrow {
  font-size: 14px;
  color: var(--text-muted);
  opacity: 0;
  transform: translateX(-4px);
  transition: all var(--transition);
}
.project-card:hover .card-arrow { opacity: 0.5; transform: translateX(0); }
.card-body { flex: 1; }
.card-name {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-family: var(--font-mono);
}
.card-desc {
  font-size: 11px;
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
}
.stat { display: flex; align-items: baseline; gap: 3px; }
.stat-num { font-size: 16px; font-weight: 700; color: var(--accent); font-family: var(--font-mono); }
.stat-label { font-size: 10px; color: var(--text-muted); text-transform: uppercase; letter-spacing: 0.5px; }
.stat-divider { width: 1px; height: 16px; background: var(--border-primary); }
.empty-state {
  width: 100%;
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 6px;
}
.empty-icon { color: var(--text-muted); opacity: 0.3; margin-bottom: 8px; }
.empty-state h2 { font-size: 18px; color: var(--text-secondary); margin: 0; font-weight: 600; }
.empty-state p { font-size: 13px; color: var(--text-muted); margin: 0 0 16px; }

.modal-overlay {
  position: fixed; inset: 0;
  background: rgba(0,0,0,0.6);
  display: flex; align-items: center; justify-content: center;
  z-index: 300;
}
.modal-box {
  background: var(--bg-surface); border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg); padding: 22px; width: 380px;
  box-shadow: 0 12px 40px rgba(0,0,0,0.5);
}
.modal-box h3 { margin: 0 0 16px; font-size: 14px; font-weight: 600; color: var(--text-primary); }
.modal-box label { display: block; font-size: 10px; color: var(--text-muted); margin-bottom: 3px; letter-spacing: 0.5px; text-transform: uppercase; }
.modal-box input {
  width: 100%; padding: 7px 10px; border: 1px solid var(--border-primary); border-radius: var(--radius);
  font-size: 12px; outline: none; background: var(--bg-base); color: var(--text-primary);
  margin-bottom: 10px; box-sizing: border-box; font-family: var(--font-mono);
  transition: border-color var(--transition);
}
.modal-box input:focus { border-color: var(--accent); }
.modal-acts { display: flex; justify-content: flex-end; gap: 8px; margin-top: 18px; }
.btn-cancel {
  padding: 6px 16px; border: 1px solid var(--border-primary); border-radius: var(--radius);
  font-size: 11px; cursor: pointer; background: var(--bg-base); color: var(--text-secondary);
  font-family: var(--font-mono); transition: all var(--transition);
}
.btn-cancel:hover { border-color: var(--border-hover); color: var(--text-primary); }
.btn-save {
  padding: 6px 16px; background: var(--accent); color: #000;
  border: 1px solid var(--accent); border-radius: var(--radius);
  font-size: 11px; cursor: pointer; font-weight: 600; font-family: var(--font-mono);
  transition: all var(--transition);
}
.btn-save:disabled { background: var(--bg-elevated); border-color: var(--border-primary); color: var(--text-muted); cursor: not-allowed; }
.btn-save:hover:not(:disabled) { background: var(--accent-hover); }
</style>
