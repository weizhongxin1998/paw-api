<script lang="ts" setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { NTabs, NTab, NButton, NIcon, NDropdown, NInput, NModal, NList, NListItem, NTag } from 'naive-ui'
import { Add, Search } from '@vicons/ionicons5'
import { useTabsStore } from '../stores/tabs'
import { useProjectStore } from '../stores/project'
import { ListCollections } from '../../wailsjs/go/handlers/CollectionHandler'
import { ListRequests } from '../../wailsjs/go/handlers/RequestHandler'

const tabsStore = useTabsStore()
const projectStore = useProjectStore()
const showDropdown = ref(false)
const showSearch = ref(false)
const searchQueryGlobal = ref('')
const searchResults = ref<any[]>([])
const allRequests = ref<any[]>([])

const addOptions = [
  { label: 'HTTP Request', key: 'http' },
  { label: 'WebSocket', key: 'websocket' },
]

async function loadAllRequestsForSearch() {
  if (!projectStore.currentProject) return
  try {
    const cols = await ListCollections(projectStore.currentProject.id)
    const all: any[] = []
    for (const col of cols) {
      const reqs = await ListRequests(col.id)
      all.push(...reqs.map((r: any) => ({ ...r, collectionName: col.name })))
    }
    allRequests.value = all
  } catch { allRequests.value = [] }
}

const filteredResults = computed(() => {
  const q = searchQueryGlobal.value.toLowerCase().trim()
  if (!q) return []
  return allRequests.value.filter((r: any) =>
    (r.name && r.name.toLowerCase().includes(q)) ||
    (r.url && r.url.toLowerCase().includes(q))
  ).slice(0, 20)
})

function handleTabChange(tabId: string) { tabsStore.setActiveTab(tabId) }

function handleAdd(key: string) {
  showDropdown.value = false
  if (key === 'http') tabsStore.addHttpTab()
  else if (key === 'websocket') tabsStore.addWsTab()
}

function closeTab(tabId: string) { tabsStore.removeTab(tabId) }

function openSearch() {
  showSearch.value = true
  searchQueryGlobal.value = ''
  loadAllRequestsForSearch()
}

function selectSearchResult(req: any) {
  showSearch.value = false
  const tabId = tabsStore.addHttpTab(req.id, req.name)
  tabsStore.updateTabData({
    method: req.method,
    url: req.url,
    headers: safeParse(req.headers, [{ key: 'Content-Type', value: 'application/json', enabled: true }]),
    params: safeParse(req.params, []),
    body: safeParseBody(req.body),
    bodyType: safeParseBodyType(req.body),
  })
}

function safeParse(str: string, fallback: any): any {
  if (!str) return fallback
  try { return JSON.parse(str) } catch { return fallback }
}
function safeParseBody(body: string): string {
  if (!body) return ''
  try { const obj = JSON.parse(body); return obj.content || obj.body || '' } catch { return body }
}
function safeParseBodyType(body: string): string {
  if (!body) return 'none'
  try { const obj = JSON.parse(body); return obj.body_type || obj.type || 'none' } catch { return 'none' }
}

function onKeydown(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
    e.preventDefault()
    openSearch()
  }
}

onMounted(() => window.addEventListener('keydown', onKeydown))
onUnmounted(() => window.removeEventListener('keydown', onKeydown))
</script>

<template>
  <div class="tab-bar">
    <NTabs v-if="tabsStore.tabs.length > 0" :value="tabsStore.activeTabId ?? undefined" type="card" size="small" closable @update:value="handleTabChange" @close="closeTab" class="tab-tabs">
      <NTab v-for="tab in tabsStore.tabs" :key="tab.id" :name="tab.id" :tab="tab.title" />
    </NTabs>
    <span v-else class="tab-placeholder">No open tabs</span>
    <NButton quaternary size="tiny" class="search-btn" @click="openSearch">
      <template #icon><NIcon size="16"><Search /></NIcon></template>
    </NButton>
    <NDropdown :options="addOptions" @select="handleAdd">
      <NButton quaternary size="tiny" class="add-tab-btn">
        <template #icon><NIcon><Add /></NIcon></template>
      </NButton>
    </NDropdown>

    <NModal v-model:show="showSearch" title="Search (Ctrl+K)" preset="card" style="width:500px" :mask-closable="true">
      <NInput v-model:value="searchQueryGlobal" placeholder="Search by name or URL..." size="small" autofocus clearable />
      <NList v-if="filteredResults.length > 0" class="search-results">
        <NListItem v-for="r in filteredResults" :key="r.id" class="search-item" @click="selectSearchResult(r)">
          <div class="search-item-content">
            <NTag size="tiny" class="search-method">{{ r.method }}</NTag>
            <span class="search-name">{{ r.name }}</span>
            <span class="search-url">{{ r.url }}</span>
            <span class="search-col">{{ r.collectionName }}</span>
          </div>
        </NListItem>
      </NList>
      <div v-else-if="searchQueryGlobal" class="search-empty">No results found</div>
    </NModal>
  </div>
</template>

<style scoped>
.tab-bar { display: flex; align-items: center; padding: 2px 4px; border-bottom: 1px solid var(--border-color); background: var(--tab-bar-bg); min-height: 34px; gap: 2px; }
.tab-tabs { flex: 1; overflow: hidden; }
.tab-placeholder { font-size: 12px; color: #999; padding: 0 8px; flex: 1; }
.search-btn { flex-shrink: 0; }
.add-tab-btn { flex-shrink: 0; }
.search-results { max-height: 300px; overflow-y: auto; margin-top: 8px; }
.search-item { cursor: pointer; }
.search-item:hover { background: var(--hover-color); }
.search-item-content { display: flex; align-items: center; gap: 8px; font-size: 13px; }
.search-method { width: 50px; text-align: center; }
.search-name { font-weight: 500; min-width: 120px; }
.search-url { color: #666; flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.search-col { color: #999; font-size: 11px; }
.search-empty { padding: 20px; text-align: center; color: #999; font-size: 13px; }
</style>
