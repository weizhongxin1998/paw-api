<template>
  <div class="project-home">
    <div class="home-header">
      <h1 class="home-title">Paw API</h1>
      <div class="home-actions">
        <button class="btn-import" @click="showImport = true">导入</button>
        <button class="btn-new" @click="showCreate = true">+ 新建项目</button>
      </div>
    </div>

    <div class="project-grid">
      <div
        v-for="p in projectList"
        :key="p.id"
        class="project-card"
        @click="enterProject(p.id)"
      >
        <div class="card-icon">📁</div>
        <div class="card-body">
          <div class="card-name">{{ p.name }}</div>
          <div class="card-desc" v-if="p.description">{{ p.description }}</div>
          <div class="card-stats">
            <span>{{ p.stats?.request_count ?? 0 }} 接口</span>
            <span>{{ p.stats?.collection_count ?? 0 }} 集合</span>
          </div>
        </div>
      </div>

      <div v-if="projectList.length === 0" class="empty-state">
        <div class="empty-icon">📦</div>
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
      try {
        stats = await GetProjectStats(p.id)
      } catch { /* ignore */ }
      items.push({
        id: p.id,
        name: p.name,
        description: p.description,
        stats,
      })
    }
    projectList.value = items
  } catch {
    projectList.value = []
  }
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

function enterProject(id: number) {
  emit('enter-project', id)
}

onMounted(loadProjects)
</script>

<style scoped>
.project-home {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f5f6f8;
}
.home-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24px 40px 0;
}
.home-title {
  font-size: 22px;
  font-weight: 700;
  color: #333;
  margin: 0;
}
.home-actions {
  display: flex;
  gap: 10px;
}
.btn-new {
  padding: 8px 20px;
  background: #18a058;
  color: #fff;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
}
.btn-new:hover { background: #0c7a43; }
.btn-import {
  padding: 8px 20px;
  background: #fff;
  color: #555;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
}
.btn-import:hover { background: #f8f8f8; }
.project-grid {
  flex: 1;
  display: flex;
  flex-wrap: wrap;
  align-content: flex-start;
  gap: 20px;
  padding: 28px 40px;
  overflow-y: auto;
}
.project-card {
  width: 200px;
  height: 160px;
  background: #fff;
  border-radius: 10px;
  cursor: pointer;
  border: 1px solid #e8e8e8;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 16px;
  transition: box-shadow 0.15s, border-color 0.15s;
}
.project-card:hover {
  box-shadow: 0 3px 16px rgba(0,0,0,0.08);
  border-color: #18a058;
}
.card-icon {
  font-size: 36px;
  margin-bottom: 10px;
  opacity: 0.7;
}
.card-body {
  text-align: center;
  width: 100%;
}
.card-name {
  font-size: 14px;
  font-weight: 600;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.card-desc {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.card-stats {
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-top: 10px;
  font-size: 12px;
  color: #888;
}
.empty-state {
  width: 100%;
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #888;
}
.empty-icon { font-size: 56px; margin-bottom: 16px; opacity: 0.25; }
.empty-state h2 { font-size: 20px; color: #555; margin: 0 0 8px; }
.empty-state p { font-size: 14px; color: #999; margin: 0 0 24px; }

.modal-overlay {
  position: fixed; inset: 0;
  background: rgba(0,0,0,0.35);
  display: flex; align-items: center; justify-content: center;
  z-index: 300;
}
.modal-box {
  background: #fff; border-radius: 10px;
  padding: 26px; width: 400px;
  box-shadow: 0 8px 30px rgba(0,0,0,0.18);
}
.modal-box h3 { margin: 0 0 18px; font-size: 16px; }
.modal-box label { display: block; font-size: 13px; color: #888; margin-bottom: 4px; }
.modal-box input {
  width: 100%; padding: 8px 12px; border: 1px solid #ddd; border-radius: 6px;
  font-size: 14px; outline: none; margin-bottom: 12px; box-sizing: border-box;
}
.modal-box input:focus { border-color: #18a058; }
.modal-acts { display: flex; justify-content: flex-end; gap: 10px; margin-top: 6px; }
.btn-cancel {
  padding: 7px 20px; border: 1px solid #ddd; border-radius: 6px;
  font-size: 13px; cursor: pointer; background: #fff; color: #555;
}
.btn-save {
  padding: 7px 20px; background: #18a058; color: #fff;
  border: 1px solid #18a058; border-radius: 6px;
  font-size: 13px; cursor: pointer;
}
.btn-save:disabled { background: #aaa; border-color: #aaa; cursor: not-allowed; }
.btn-save:hover:not(:disabled) { background: #0c7a43; }
</style>
