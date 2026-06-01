<script lang="ts" setup>
import { ref, onMounted, computed } from 'vue'
import { NInput, NTag, NButton, NIcon, NEmpty, NSpin } from 'naive-ui'
import { Play } from '@vicons/ionicons5'
import { ListHistory } from '../../wailsjs/go/handlers/HistoryHandler'
import { useTabsStore } from '../stores/tabs'

const tabsStore = useTabsStore()
const history = ref<any[]>([])
const loading = ref(false)
const searchQuery = ref('')

const filteredHistory = computed(() => {
  if (!searchQuery.value.trim()) return history.value
  const q = searchQuery.value.toLowerCase()
  return history.value.filter((h: any) => h.url.toLowerCase().includes(q) || h.method.toLowerCase().includes(q))
})

function methodTagType(method: string): 'success' | 'info' | 'warning' | 'error' | 'default' {
  const map: Record<string, any> = { GET: 'success', POST: 'info', PUT: 'warning', DELETE: 'error', PATCH: 'info' }
  return map[method] ?? 'default'
}

async function loadHistory() {
  loading.value = true
  try { history.value = await ListHistory('', 100) } catch { history.value = [] }
  finally { loading.value = false }
}

function replay(item: any) {
  tabsStore.addHttpTab(undefined, `${item.method} ${item.url.slice(0, 40)}`)
  tabsStore.updateHttpData({ method: item.method, url: item.url, headers: parseJSON(item.headers, []), body: item.body || '', bodyType: item.body ? 'json' : 'none' })
}

function parseJSON(str: string, fallback: any): any { try { return JSON.parse(str) } catch { return fallback } }

onMounted(loadHistory)
</script>

<template>
  <div class="history-view">
    <div class="history-toolbar">
      <h2 class="history-title">{{ $t('history.title') }}</h2>
      <NInput v-model:value="searchQuery" :placeholder="$t('history.search')" size="small" clearable class="search-input" />
      <NButton size="small" quaternary @click="loadHistory">{{ $t('history.refresh') }}</NButton>
    </div>
    <NSpin :show="loading">
      <div v-if="filteredHistory.length === 0" class="empty"><NEmpty :description="$t('history.empty')" /></div>
      <div v-else class="history-list">
        <div v-for="item in filteredHistory" :key="item.id" class="history-item">
          <NTag :type="methodTagType(item.method)" size="small" class="method-tag">{{ item.method }}</NTag>
          <div class="history-url">{{ item.url }}</div>
          <div class="history-status">{{ item.response_status }}</div>
          <div class="history-duration">{{ item.duration_ms }}{{ $t('history.ms') }}</div>
          <NButton quaternary circle size="tiny" @click="replay(item)">
            <template #icon><NIcon><Play /></NIcon></template>
          </NButton>
        </div>
      </div>
    </NSpin>
  </div>
</template>

<style scoped>
.history-view { padding: 16px 20px; height: 100%; display: flex; flex-direction: column; }
.history-toolbar { display: flex; align-items: center; gap: 12px; margin-bottom: 16px; }
.history-title { font-size: 18px; font-weight: 600; white-space: nowrap; }
.search-input { max-width: 320px; }
.history-list { flex: 1; overflow-y: auto; }
.history-item { display: flex; align-items: center; gap: 8px; padding: 6px 8px; border-bottom: 1px solid var(--border-color); font-size: 13px; }
.history-item:hover { background: var(--hover-color); }
.method-tag { width: 60px; text-align: center; }
.history-url { flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; color: #333; }
.history-status { width: 40px; text-align: right; font-weight: 600; color: #18a058; }
.history-duration { width: 50px; text-align: right; color: #999; font-size: 12px; }
.empty { flex: 1; display: flex; align-items: center; justify-content: center; }
</style>
