<script lang="ts" setup>
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { NLayoutSider, NTree, NButton, NIcon, NInput, NModal, NForm, NFormItem, NSpace, useMessage } from 'naive-ui'
import { Add } from '@vicons/ionicons5'
import { useRouter, useRoute } from 'vue-router'
import { useProjectStore } from '../stores/project'
import { CreateCollection } from '../../wailsjs/go/handlers/CollectionHandler'
import EnvSelector from './EnvSelector.vue'
import EnvManager from './EnvManager.vue'

const { t, locale } = useI18n()
const router = useRouter()
const route = useRoute()
const projectStore = useProjectStore()
const message = useMessage()

const showAddModal = ref(false)
const newCollectionName = ref('')
const showEnvManager = ref(false)

watch(locale, (val) => { localStorage.setItem('paw-locale', val) })

function toggleLocale() {
  locale.value = locale.value === 'zh-CN' ? 'en' : 'zh-CN'
}

const menuItems = [
  { labelKey: 'sidebar.workspace', key: '/workspace' },
  { labelKey: 'sidebar.history', key: '/history' },
  { labelKey: 'sidebar.docs', key: '/docs' },
  { labelKey: 'sidebar.testRunner', key: '/test-runner' },
  { labelKey: 'sidebar.websocket', key: '/websocket' },
]

const treeData = computed(() => {
  function buildTree(parentId: string | null): any[] {
    return projectStore.collections
      .filter(c => c.parent_id === parentId)
      .sort((a, b) => a.sort_order - b.sort_order)
      .map(c => ({ label: c.name, key: c.id, isLeaf: false, children: buildTree(c.id) }))
  }
  return [
    ...menuItems.map(m => ({ label: t(m.labelKey), key: m.key, isLeaf: true })),
    { label: t('sidebar.collections'), key: 'collections-header', isLeaf: true, disabled: true },
    ...buildTree(null),
  ]
})

function handleNodeSelect(keys: string[]) {
  if (keys.length === 0) return
  const key = keys[0]
  if (menuItems.find(m => m.key === key)) { router.push(key); return }
  projectStore.selectCollection(key)
}

function startAdd() { newCollectionName.value = ''; showAddModal.value = true }

async function confirmAdd() {
  if (!newCollectionName.value.trim()) return
  if (!projectStore.currentProject) { message.error(t('sidebar.noProject')); return }
  try {
    const col = await CreateCollection(projectStore.currentProject.id, '', newCollectionName.value.trim(), projectStore.collections.length)
    projectStore.addCollection(col)
    showAddModal.value = false
    message.success(t('sidebar.created'))
  } catch (e: any) { message.error(e.message || t('sidebar.failedCreate')) }
}
</script>

<template>
  <NLayoutSider bordered width="220" :native-scrollbar="false" class="app-sidebar">
    <div class="sidebar-header">
      <div class="sidebar-title-row">
        <span class="sidebar-title">{{ t('sidebar.title') }}</span>
        <NButton quaternary size="tiny" @click="toggleLocale">{{ locale === 'zh-CN' ? 'EN' : '中文' }}</NButton>
      </div>
      <EnvSelector @manage="showEnvManager = true" />
    </div>
    <NTree :data="treeData" :default-expand-all="true" block-line selectable @update:selected-keys="handleNodeSelect" />
    <div class="sidebar-actions">
      <NButton size="tiny" quaternary @click="startAdd">
        <template #icon><NIcon><Add /></NIcon></template>
        {{ t('sidebar.newCollection') }}
      </NButton>
    </div>
    <NModal v-model:show="showAddModal" :title="t('sidebar.newCollection')" preset="card" style="width: 360px">
      <NForm>
        <NFormItem :label="t('sidebar.collectionName')">
          <NInput v-model:value="newCollectionName" :placeholder="t('sidebar.collectionName')" />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace justify="end">
          <NButton @click="showAddModal = false">{{ t('sidebar.cancel') }}</NButton>
          <NButton type="primary" @click="confirmAdd">{{ t('sidebar.create') }}</NButton>
        </NSpace>
      </template>
    </NModal>
    <EnvManager v-model:show="showEnvManager" />
  </NLayoutSider>
</template>

<style scoped>
.app-sidebar { height: 100%; display: flex; flex-direction: column; }
.sidebar-header { padding: 12px 12px 8px; border-bottom: 1px solid var(--border-color); display: flex; flex-direction: column; gap: 6px; }
.sidebar-title-row { display: flex; align-items: center; justify-content: space-between; }
.sidebar-title { font-size: 16px; font-weight: 700; }
.sidebar-actions { padding: 8px; border-top: 1px solid var(--border-color); margin-top: auto; }
</style>
