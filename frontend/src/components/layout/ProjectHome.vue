<template>
  <div class="project-home">
    <div class="home-bg-grid"></div>

    <!-- ═══ Header ═══ -->
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
        <button class="btn-theme" @click="emit('toggle-theme')" :title="props.themeMode === 'dark' ? t('projectHome.dayMode') : t('projectHome.nightMode')">
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
            <circle v-if="props.themeMode === 'dark'" cx="12" cy="12" r="5"/><line v-if="props.themeMode === 'dark'" x1="12" y1="1" x2="12" y2="3"/><line v-if="props.themeMode === 'dark'" x1="12" y1="21" x2="12" y2="23"/><line v-if="props.themeMode === 'dark'" x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line v-if="props.themeMode === 'dark'" x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line v-if="props.themeMode === 'dark'" x1="1" y1="12" x2="3" y2="12"/><line v-if="props.themeMode === 'dark'" x1="21" y1="12" x2="23" y2="12"/><line v-if="props.themeMode === 'dark'" x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line v-if="props.themeMode === 'dark'" x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>
            <path v-else d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
          </svg>
        </button>
        <n-button size="tiny" @click="showImport = true" class="btn-import">
          <template #icon>
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
          </template>
          {{ t('projectHome.import') }}
        </n-button>
        <n-button type="primary" size="tiny" @click="showCreate = true">{{ t('projectHome.newProject') }}</n-button>
      </div>
    </div>

    <!-- ═══ Toolbar: search + count ═══ -->
    <div class="home-toolbar">
      <n-input
        ref="searchInputRef"
        v-model:value="searchQuery"
        :placeholder="t('projectHome.searchProjects')"
        clearable
        size="small"
        class="search-input"
        @keyup.enter="focusFirstCard"
      >
        <template #prefix>
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/>
          </svg>
        </template>
      </n-input>
      <span class="project-count" v-if="projectList.length > 0">
        {{ t('projectHome.projectCount', { count: filteredProjects.length + (searchQuery && filteredProjects.length !== projectList.length ? ' / ' + projectList.length : '') }) }}
      </span>
    </div>

    <!-- ═══ Project Grid ═══ -->
    <div class="project-grid" ref="gridRef">
      <TransitionGroup name="card-filter">
        <div
          v-for="(p, idx) in filteredProjects"
          :key="p.id"
          class="project-card"
          :style="{ '--card-delay': idx * 0.05 + 's' }"
          :ref="(el) => { if (el) cardRefs[p.id] = el as HTMLElement }"
          @click="enterProject(p.id)"
          @keydown.enter="enterProject(p.id)"
          @contextmenu="onCardContext($event, p)"
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

          <div class="card-footer">
            <div class="card-stats">
              <span class="stat">
                <span class="stat-num">{{ p.stats?.request_count ?? 0 }}</span>
                <span class="stat-label">{{ t('projectHome.statInterfaces') }}</span>
              </span>
              <span class="stat-divider"></span>
              <span class="stat">
                <span class="stat-num">{{ p.stats?.collection_count ?? 0 }}</span>
                <span class="stat-label">{{ t('projectHome.statCollections') }}</span>
              </span>
            </div>
            <div class="card-footer-right">
              <span class="card-time" v-if="p.last_opened">{{ relativeTime(p.last_opened) }}</span>
              <div class="card-actions">
                <button
                  class="card-action-btn"
                  :title="t('common.edit')"
                  @click.stop="onStartRename(p)"
                >
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M17 3a2.83 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z"/>
                  </svg>
                </button>
                <button
                  class="card-action-btn card-action-delete"
                  :title="t('common.delete')"
                  @click.stop="onStartDelete(p)"
                >
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <polyline points="3 6 5 6 21 6"/>
                    <path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/>
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </div>
      </TransitionGroup>

      <!-- ── Empty: no projects at all ── -->
      <div v-if="projectList.length === 0 && !loading" class="empty-state">
        <div class="empty-graphic">
          <svg class="paw-svg" width="88" height="88" viewBox="0 0 88 88" fill="none">
            <!-- Toe pads -->
            <ellipse cx="24" cy="24" rx="8.5" ry="10.5" fill="currentColor" fill-opacity="0.07" stroke="currentColor" stroke-opacity="0.18" stroke-width="1.4"/>
            <ellipse cx="40" cy="15" rx="7.5" ry="9.5" fill="currentColor" fill-opacity="0.07" stroke="currentColor" stroke-opacity="0.18" stroke-width="1.4"/>
            <ellipse cx="56" cy="17" rx="7.5" ry="9.5" fill="currentColor" fill-opacity="0.07" stroke="currentColor" stroke-opacity="0.18" stroke-width="1.4"/>
            <ellipse cx="70" cy="28" rx="7" ry="9" fill="currentColor" fill-opacity="0.07" stroke="currentColor" stroke-opacity="0.18" stroke-width="1.4"/>
            <!-- Main pad -->
            <path d="M24 52 C24 40, 34 33, 44 33 C54 33, 64 40, 64 52 C64 64, 56 72, 44 72 C32 72, 24 64, 24 52Z" fill="currentColor" fill-opacity="0.05" stroke="currentColor" stroke-opacity="0.18" stroke-width="1.4"/>
            <!-- Inner pad lines -->
            <path d="M34 46 Q44 38, 54 46" stroke="currentColor" stroke-opacity="0.1" stroke-width="1" fill="none"/>
            <path d="M30 56 Q44 64, 58 56" stroke="currentColor" stroke-opacity="0.1" stroke-width="1" fill="none"/>
          </svg>
          <div class="empty-ring"></div>
          <div class="empty-ring empty-ring-2"></div>
          <div class="paw-float"></div>
        </div>
        <h2 class="empty-title">{{ t('projectHome.noProjects') }}</h2>
        <p class="empty-desc">{{ t('projectHome.noProjectsDesc') }}</p>
        <n-button type="primary" @click="showCreate = true">{{ t('projectHome.newProject') }}</n-button>
      </div>

      <!-- ── Empty: search yields no results ── -->
      <div v-if="projectList.length > 0 && filteredProjects.length === 0 && !loading" class="empty-state empty-filter">
        <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" class="empty-filter-icon">
          <circle cx="11" cy="11" r="8" opacity="0.3"/><line x1="21" y1="21" x2="16.65" y2="16.65" opacity="0.3"/>
          <line x1="8" y1="8" x2="14" y2="14" opacity="0.5"/><line x1="14" y1="8" x2="8" y2="14" opacity="0.5"/>
        </svg>
        <p class="empty-filter-text">{{ t('projectHome.noMatch', { query: searchQuery }) }}</p>
      </div>
    </div>

    <!-- ═══ Right-click context menu ═══ -->
    <n-dropdown
      ref="dropdownRef"
      :show="dropdownVisible"
      :options="contextOptions"
      :x="dropdownX"
      :y="dropdownY"
      trigger="manual"
      placement="bottom-start"
      @select="onContextSelect"
      @update:show="onDropdownShowChange"
    />

    <!-- ═══ Create Project Modal ═══ -->
    <n-modal v-model:show="showCreate" preset="card" :title="t('projectHome.createProject')" :class="modalClass" style="width: 400px" :mask-closable="false">
      <n-form label-placement="top">
        <n-form-item :label="t('common.name')">
          <n-input v-model:value="newName" :placeholder="t('projectHome.namePlaceholder')" @keydown.enter="onCreate" />
        </n-form-item>
        <n-form-item :label="t('common.description')">
          <n-input v-model:value="newDesc" :placeholder="t('projectHome.descPlaceholder')" @keydown.enter="onCreate" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button @click="showCreate = false">{{ t('common.cancel') }}</n-button>
        <n-button type="primary" :disabled="!newName.trim()" @click="onCreate">{{ t('common.create') }}</n-button>
      </template>
    </n-modal>

    <!-- ═══ Edit Project Modal ═══ -->
    <n-modal v-model:show="showRename" preset="card" :title="t('projectHome.editProject')" :class="modalClass" style="width: 400px" :mask-closable="false">
      <n-form label-placement="top">
        <n-form-item :label="t('common.name')">
          <n-input
            ref="renameInputRef"
            v-model:value="renameValue"
            :placeholder="t('projectHome.namePlaceholder')"
            @keydown.enter="onRenameConfirm"
          />
        </n-form-item>
        <n-form-item :label="t('common.description')">
          <n-input
            v-model:value="renameDesc"
            :placeholder="t('projectHome.descPlaceholder')"
            @keydown.enter="onRenameConfirm"
          />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button @click="showRename = false">{{ t('common.cancel') }}</n-button>
        <n-button type="primary" :disabled="!renameValue.trim()" @click="onRenameConfirm">{{ t('common.save') }}</n-button>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick, watch, h, type VNode } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  NButton, NModal, NForm, NFormItem, NInput, NDropdown,
  useMessage, useDialog
} from 'naive-ui'
import {
  ListProjects, CreateProject, GetProjectStats, DeleteProject, UpdateProject
} from '../../../wailsjs/go/main/App'

const { t } = useI18n()

/* ──────────────────────────── Types ──────────────────────────── */

interface ProjectCard {
  id: number
  name: string
  description: string
  stats: { request_count: number; collection_count: number } | null
  last_opened?: number
}

/* ──────────────────────────── Props & Emits ──────────────────── */

const props = defineProps<{ themeMode: 'dark' | 'light' }>()

const emit = defineEmits<{
  (e: 'enter-project', id: number): void
  (e: 'toggle-theme'): void
}>()

/* ──────────────────────────── State ──────────────────────────── */

const projectList = ref<ProjectCard[]>([])
const loading = ref(true)
const searchQuery = ref('')
const showCreate = ref(false)
const showImport = ref(false)
const showRename = ref(false)
const newName = ref('')
const newDesc = ref('')
const renameValue = ref('')
const renameDesc = ref('')
const renameTargetId = ref<number | null>(null)

const message = useMessage()
const dialog = useDialog()

/* Refs */
const searchInputRef = ref<InstanceType<typeof NInput> | null>(null)
const renameInputRef = ref<InstanceType<typeof NInput> | null>(null)
const cardRefs = ref<Record<number, HTMLElement>>({})

const isLightMode = ref(false)
const modalClass = computed(() => isLightMode.value ? 'project-home-modal theme-light' : 'project-home-modal')

/* ──────────────────────────── Context menu ───────────────────── */

const dropdownVisible = ref(false)
const dropdownX = ref(0)
const dropdownY = ref(0)
const dropdownRef = ref<InstanceType<typeof NDropdown> | null>(null)
const contextTargetId = ref<number | null>(null)

const contextOptions = computed(() => [
  {
    label: t('common.edit'),
    key: 'rename',
    icon: () => h('svg', {
      width: 14, height: 14, viewBox: '0 0 24 24',
      fill: 'none', stroke: 'currentColor', 'stroke-width': 2,
      'stroke-linecap': 'round', 'stroke-linejoin': 'round'
    }, [
      h('path', { d: 'M17 3a2.83 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z' })
    ]) as VNode
  },
  { type: 'divider' as const, key: 'd1' },
  {
    label: t('common.delete'),
    key: 'delete',
    icon: () => h('svg', {
      width: 14, height: 14, viewBox: '0 0 24 24',
      fill: 'none', stroke: 'currentColor', 'stroke-width': 2,
      'stroke-linecap': 'round', 'stroke-linejoin': 'round'
    }, [
      h('polyline', { points: '3 6 5 6 21 6' }),
      h('path', { d: 'M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2' })
    ]) as VNode
  }
])

/* ──────────────────────────── Computed ───────────────────────── */

const filteredProjects = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  const list = q
    ? projectList.value.filter(p => p.name.toLowerCase().includes(q))
    : [...projectList.value]
  // Always show most recently opened first
  list.sort((a, b) => (b.last_opened ?? 0) - (a.last_opened ?? 0))
  return list
})

/* ──────────────────────────── Relative time ──────────────────── */

function relativeTime(ts: number): string {
  // Support both seconds and milliseconds
  const ms = ts > 1e12 ? ts : ts * 1000
  const diff = Math.max(0, Date.now() - ms)
  const seconds = Math.floor(diff / 1000)
  if (seconds < 60) return t('time.justNow')
  const minutes = Math.floor(seconds / 60)
  if (minutes < 60) return t('time.minutesAgo', { n: minutes })
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return t('time.hoursAgo', { n: hours })
  const days = Math.floor(hours / 24)
  if (days < 30) return t('time.daysAgo', { n: days })
  const months = Math.floor(days / 30)
  if (months < 12) return t('time.monthsAgo', { n: months })
  return t('time.yearsAgo', { n: Math.floor(months / 12) })
}

/* ──────────────────────────── Data loading ───────────────────── */

async function loadProjects() {
  loading.value = true
  try {
    const projects = await ListProjects()
    const items: ProjectCard[] = []
    for (const p of projects || []) {
      let stats = null
      try { stats = await GetProjectStats(p.id) } catch { /* ignore */ }
      items.push({
        id: p.id,
        name: p.name,
        description: p.description,
        stats,
        last_opened: (p as any).last_opened ?? (Date.parse(p.updated_at || '0') || null)
      })
    }
    projectList.value = items
  } catch {
    projectList.value = []
  } finally {
    loading.value = false
  }
}

/* ──────────────────────────── Actions ────────────────────────── */

async function onCreate() {
  const name = newName.value.trim()
  if (!name) return
  try {
    const p = await CreateProject(name, newDesc.value.trim())
    showCreate.value = false
    newName.value = ''
    newDesc.value = ''
    message.success(t('projectHome.projectCreated', { name }))
    await loadProjects()
    enterProject(p.id)
  } catch (e: any) {
    message.error(t('projectHome.createFailed', { error: e?.message || String(e) }))
  }
}

function enterProject(id: number) {
  emit('enter-project', id)
}

/* ── Card action button handlers ── */

function onStartRename(p: ProjectCard) {
  renameTargetId.value = p.id
  renameValue.value = p.name
  renameDesc.value = p.description || ''
  showRename.value = true
  nextTick(() => {
    const nativeInput = renameInputRef.value?.$el?.querySelector?.('input') as HTMLInputElement | null
    nativeInput?.select()
  })
}

function onStartDelete(p: ProjectCard) {
  dialog.warning({
    title: t('projectHome.confirmDeleteTitle'),
    content: t('projectHome.confirmDeleteContent', { name: p.name }),
    positiveText: t('common.delete'),
    negativeText: t('common.cancel'),
    onPositiveClick: async () => {
      try {
        await DeleteProject(p.id)
        message.success(t('projectHome.projectDeleted', { name: p.name }))
        await loadProjects()
      } catch (e: any) {
        message.error(t('projectHome.deleteFailed', { error: e?.message || String(e) }))
      }
    }
  })
}

/* ── Context menu handlers ── */

function onCardContext(e: MouseEvent, p: ProjectCard) {
  e.preventDefault()
  e.stopPropagation()
  contextTargetId.value = p.id
  dropdownX.value = e.clientX
  dropdownY.value = e.clientY
  dropdownVisible.value = true
}

function onContextSelect(key: string) {
  dropdownVisible.value = false
  const id = contextTargetId.value
  if (id == null) return
  const project = projectList.value.find(p => p.id === id)
  if (!project) return

  if (key === 'rename') {
    onStartRename(project)
  } else if (key === 'delete') {
    onStartDelete(project)
  }
}

function onDropdownShowChange(show: boolean) {
  dropdownVisible.value = show
}

async function onRenameConfirm() {
  const name = renameValue.value.trim()
  if (!name || renameTargetId.value == null) return
  try {
    await UpdateProject(renameTargetId.value, name, renameDesc.value.trim())
    showRename.value = false
    message.success(t('common.save'))
    await loadProjects()
  } catch (e: any) {
    message.error(t('projectHome.saveFailed', { error: e?.message || String(e) }))
  }
}

/* ── Keyboard: Enter in search focuses first matching card ── */

function focusFirstCard() {
  const first = filteredProjects.value[0]
  if (!first) return
  const el = cardRefs.value[first.id]
  if (el) {
    el.focus()
  }
}

/* ──────────────────────────── Lifecycle ──────────────────────── */

onMounted(async () => {
  await loadProjects()
  // Auto-focus search input after load
  nextTick(() => {
    searchInputRef.value?.focus()
  })
  const check = () => { isLightMode.value = !!document.querySelector('.theme-light') }
  check()
  const observer = new MutationObserver(check)
  observer.observe(document.documentElement, { attributes: true, attributeFilter: ['class'], subtree: true })
})

// Clear card refs when filtered list changes to avoid stale refs
watch(filteredProjects, () => {
  cardRefs.value = {}
})
</script>

<style scoped>
/* ═══════════════════════════════════════════════════════════════
   Layout
   ═══════════════════════════════════════════════════════════════ */
.project-home {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--bg-base);
  position: relative;
  overflow: hidden;
}

/* ── Background grid ── */
.home-bg-grid {
  position: absolute;
  inset: 0;
  pointer-events: none;
  background-image:
    linear-gradient(var(--border-subtle) 0.5px, transparent 0.5px),
    linear-gradient(90deg, var(--border-subtle) 0.5px, transparent 0.5px);
  background-size: 44px 44px;
  opacity: 0.5;
  mask-image: radial-gradient(ellipse 70% 60% at 50% 0%, black 20%, transparent 70%);
}

/* ═══════════════════════════════════════════════════════════════
   Header
   ═══════════════════════════════════════════════════════════════ */
.home-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 32px 44px 0;
  position: relative;
  animation: slideUp 0.4s var(--ease-out) both;
}
.home-brand {
  display: flex;
  align-items: center;
  gap: 12px;
}
.brand-mark {
  color: var(--accent);
  filter: drop-shadow(0 0 8px var(--accent-glow));
  animation: brandGradientPulse 4s ease-in-out infinite;
}

@keyframes brandGradientPulse {
  0%, 100% {
    filter: drop-shadow(0 0 8px var(--accent-glow));
  }
  50% {
    filter: drop-shadow(0 0 18px var(--accent-glow-strong))
           drop-shadow(0 0 4px var(--accent));
  }
}

.home-title {
  font-size: var(--fs-3xl);
  font-weight: 800;
  color: var(--text-primary);
  letter-spacing: 5px;
  margin: 0;
  font-family: var(--font-ui);
}
.home-version {
  font-size: var(--fs-xs);
  color: var(--text-secondary);
  background: var(--bg-elevated);
  padding: 2px 7px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border-primary);
  font-family: var(--font-mono);
  font-weight: 500;
}
.home-actions { display: flex; gap: 8px; align-items: center; }
.btn-theme {
  padding: 8px 11px;
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
.btn-theme:hover {
  border-color: var(--accent);
  color: var(--accent);
  background: var(--accent-soft);
}
.btn-import {
  border-radius: var(--radius) !important;
}

/* ═══════════════════════════════════════════════════════════════
   Toolbar (search + count)
   ═══════════════════════════════════════════════════════════════ */
.home-toolbar {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 20px 44px 0;
  animation: slideUp 0.4s var(--ease-out) 0.08s both;
}
.search-input {
  max-width: 260px;
}
.search-input :deep(.n-input__prefix) {
  color: var(--text-muted);
}
.project-count {
  font-size: var(--fs-sm);
  color: var(--text-muted);
  font-family: var(--font-ui);
  white-space: nowrap;
  letter-spacing: 0.02em;
}

/* ═══════════════════════════════════════════════════════════════
   Project Grid  (responsive with min-width)
   ═══════════════════════════════════════════════════════════════ */
.project-grid {
  flex: 1;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(232px, 1fr));
  align-content: start;
  gap: 18px;
  padding: 28px 44px;
  overflow-y: auto;
  position: relative;
}

/* ═══════════════════════════════════════════════════════════════
   Project Card
   ═══════════════════════════════════════════════════════════════ */
.project-card {
  min-height: 160px;
  background: var(--bg-surface);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg);
  cursor: pointer;
  display: flex;
  flex-direction: column;
  padding: 18px;
  position: relative;
  transition: all var(--transition-slow);
  animation: cardSlideIn 0.4s var(--ease-out) both;
  animation-delay: var(--card-delay, 0s);
  outline: none;
  overflow: hidden;
}
.project-card:focus-visible {
  border-color: var(--accent);
  box-shadow: 0 0 0 2px var(--accent-glow);
}
.project-card:hover {
  border-color: var(--accent);
  box-shadow:
    0 0 0 1px var(--accent),
    0 8px 28px var(--accent-glow),
    0 -12px 40px -8px var(--accent-glow-strong);
  transform: translateY(-4px);
}

/* ── Gradient glow overlay (top-down) ── */
.card-glow {
  position: absolute;
  inset: 0;
  pointer-events: none;
  background:
    linear-gradient(180deg, var(--accent-glow-strong) 0%, transparent 45%),
    radial-gradient(ellipse 130% 70% at 50% -15%, var(--accent-glow-strong) 0%, transparent 55%);
  opacity: 0;
  transition: opacity var(--transition-slow);
}
.project-card:hover .card-glow { opacity: 1; }

@keyframes cardSlideIn {
  from { opacity: 0; transform: translateY(16px) scale(0.96); }
  to   { opacity: 1; transform: translateY(0) scale(1); }
}

/* ── TransitionGroup stagger for filtered results ── */
.card-filter-enter-active {
  transition: all 0.35s var(--ease-out);
  transition-delay: var(--card-delay, 0s);
}
.card-filter-leave-active {
  transition: all 0.2s var(--ease-out);
}
.card-filter-enter-from {
  opacity: 0;
  transform: translateY(14px) scale(0.96);
}
.card-filter-leave-to {
  opacity: 0;
  transform: scale(0.94);
}
.card-filter-move {
  transition: transform 0.3s var(--ease-out);
}

/* ── Card header ── */
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 14px;
  position: relative;
  z-index: 1;
}
.card-icon {
  color: var(--accent);
  opacity: 0.7;
}
.card-arrow {
  font-size: var(--fs-lg);
  color: var(--text-muted);
  opacity: 0;
  transform: translateX(-6px);
  transition: all var(--transition);
}
.project-card:hover .card-arrow {
  opacity: 0.6;
  transform: translateX(0);
}

/* ── Card action buttons (in footer, right side) ── */
.card-footer-right {
  display: flex;
  align-items: center;
  gap: 8px;
  min-height: 28px;
}
.card-actions {
  display: flex;
  align-items: center;
  gap: 4px;
  opacity: 0;
  transform: translateX(6px);
  transition: all 0.25s ease;
}
.project-card:hover .card-actions {
  opacity: 1;
  transform: translateX(0);
}
.card-action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border-primary);
  background: var(--bg-elevated);
  color: var(--text-muted);
  cursor: pointer;
  transition: all var(--transition);
}
.card-action-btn:hover {
  background: var(--accent-soft);
  color: var(--accent);
  border-color: var(--accent);
}
.card-action-delete:hover {
  background: var(--error-soft, rgba(224, 82, 82, 0.1));
  color: var(--error, #e05252);
  border-color: var(--error, #e05252);
}

/* ── Card body ── */
.card-body {
  flex: 1;
  position: relative;
  z-index: 1;
}
.card-name {
  font-size: var(--fs-md);
  font-weight: 600;
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-family: var(--font-mono);
}
.card-desc {
  font-size: var(--fs-sm);
  color: var(--text-secondary);
  margin-top: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-family: var(--font-ui);
}

/* ── Card footer (stats + time + actions) ── */
.card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 12px;
  border-top: 1px solid var(--border-subtle);
  position: relative;
  z-index: 1;
}
.card-stats {
  display: flex;
  align-items: center;
  gap: 10px;
}
.stat { display: flex; align-items: baseline; gap: 4px; }
.stat-num {
  font-size: var(--fs-lg);
  font-weight: 700;
  color: var(--accent);
  font-family: var(--font-mono);
}
.stat-label {
  font-size: var(--fs-sm);
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
  font-family: var(--font-ui);
}
.stat-divider { width: 1px; height: 18px; background: var(--border-primary); }

.card-time {
  font-size: var(--fs-xs);
  color: var(--text-muted);
  font-family: var(--font-ui);
  white-space: nowrap;
}

/* ═══════════════════════════════════════════════════════════════
   Empty State  (paw-print themed)
   ═══════════════════════════════════════════════════════════════ */
.empty-state {
  grid-column: 1 / -1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  animation: slideUp 0.5s var(--ease-out) both;
  padding: 48px 0;
}
.empty-graphic {
  position: relative;
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 120px;
  height: 120px;
}
.paw-svg {
  color: var(--accent);
  display: block;
  position: relative;
  z-index: 1;
}

/* Animated rings */
.empty-ring {
  position: absolute;
  inset: -6px;
  border: 1px solid var(--border-primary);
  border-radius: 50%;
  animation: spin 12s linear infinite;
}
.empty-ring-2 {
  inset: -18px;
  border-style: dashed;
  border-color: var(--border-subtle);
  animation: spin 20s linear infinite reverse;
}

/* Floating pulse behind paw */
.paw-float {
  position: absolute;
  inset: 10px;
  border-radius: 50%;
  background: radial-gradient(circle, var(--accent-glow) 0%, transparent 70%);
  animation: pawPulse 3s ease-in-out infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }

@keyframes pawPulse {
  0%, 100% { opacity: 0.3; transform: scale(0.9); }
  50%      { opacity: 0.6; transform: scale(1.08); }
}

.empty-title {
  font-size: var(--fs-xl);
  color: var(--text-secondary);
  margin: 0;
  font-weight: 600;
  font-family: var(--font-ui);
}
.empty-desc {
  font-size: var(--fs-base);
  color: var(--text-secondary);
  margin: 0 0 4px;
  font-family: var(--font-ui);
}

/* ── Empty filter state ── */
.empty-filter {
  padding: 32px 0;
}
.empty-filter-icon {
  color: var(--text-muted);
  margin-bottom: 4px;
}
.empty-filter-text {
  font-size: var(--fs-base);
  color: var(--text-secondary);
  margin: 0;
  font-family: var(--font-ui);
}
</style>
