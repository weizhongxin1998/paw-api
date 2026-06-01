<script lang="ts" setup>
import { NTabs, NTab, NButton, NIcon } from 'naive-ui'
import { Add } from '@vicons/ionicons5'
import { useTabsStore } from '../stores/tabs'

const tabsStore = useTabsStore()

function handleTabChange(tabId: string) { tabsStore.setActiveTab(tabId) }
function addNewTab() { tabsStore.addTab() }
function closeTab(tabId: string) { tabsStore.removeTab(tabId) }
</script>

<template>
  <div class="tab-bar">
    <NTabs v-if="tabsStore.tabs.length > 0" :value="tabsStore.activeTabId ?? undefined" type="card" size="small" closable @update:value="handleTabChange" @close="closeTab" class="tab-tabs">
      <NTab v-for="tab in tabsStore.tabs" :key="tab.id" :name="tab.id" :tab="tab.title" />
    </NTabs>
    <span v-else class="tab-placeholder">{{ $t('request.noTabs') }}</span>
    <NButton quaternary size="tiny" @click="addNewTab" class="add-tab-btn">
      <template #icon><NIcon><Add /></NIcon></template>
    </NButton>
  </div>
</template>

<style scoped>
.tab-bar { display: flex; align-items: center; padding: 2px 4px; border-bottom: 1px solid var(--border-color); background: var(--tab-bar-bg); min-height: 34px; gap: 2px; }
.tab-tabs { flex: 1; overflow: hidden; }
.tab-placeholder { font-size: 12px; color: #999; padding: 0 8px; flex: 1; }
.add-tab-btn { flex-shrink: 0; }
</style>
