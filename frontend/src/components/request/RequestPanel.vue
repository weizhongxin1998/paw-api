<template>
  <div class="request-panel">
    <div class="sub-tabs">
      <n-button
        v-for="tab in tabs"
        :key="tab.key"
        :type="activeTab === tab.key ? 'primary' : 'default'"
        size="small"
        :ghost="activeTab !== tab.key"
        @click="activeTab = tab.key"
      >
        {{ tab.label }}
        <span v-if="tab.count != null" class="tab-count">{{ tab.count }}</span>
      </n-button>
    </div>

    <div class="sub-content">
      <KeyValueTable
        v-if="activeTab === 'params'"
        :items="paramsItems"
        @update:items="onParamsChange"
      />
      <KeyValueTable
        v-else-if="activeTab === 'headers'"
        :items="headersItems"
        :show-bulk-edit="true"
        @update:items="onHeadersChange"
      />
      <BodyEditor
        v-else-if="activeTab === 'body'"
        :body-type="bodyType"
        :body-data="bodyData"
        @update:body-type="v => { bodyType = v; sync() }"
        @update:body-data="v => { bodyData = v; sync() }"
      />
      <AuthEditor
        v-else-if="activeTab === 'auth'"
        :model-value="authData"
        @update:model-value="v => { authData = v; sync() }"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { NButton } from 'naive-ui'
import KeyValueTable from '../shared/KeyValueTable.vue'
import BodyEditor from './BodyEditor.vue'
import AuthEditor from './AuthEditor.vue'
import type { KvItem } from '../../types/request'

const props = defineProps<{
  headers: string
  params: string
  bodyType: string
  bodyData: string
  authData: string
}>()

const emit = defineEmits<{
  (e: 'update:headers', v: string): void
  (e: 'update:params', v: string): void
  (e: 'update:bodyType', v: string): void
  (e: 'update:bodyData', v: string): void
  (e: 'update:authData', v: string): void
}>()

const tabs = [
  { key: 'params', label: 'Params', count: undefined as number | undefined },
  { key: 'headers', label: 'Headers', count: undefined as number | undefined },
  { key: 'body', label: 'Body' },
  { key: 'auth', label: 'Auth' },
]

const activeTab = ref('params')

const bodyType = ref(props.bodyType)
const bodyData = ref(props.bodyData)
const authData = ref(props.authData)

const paramsItems = ref<KvItem[]>([])
const headersItems = ref<KvItem[]>([])

let iid = Date.now()

function parseKv(raw: string): KvItem[] {
  try {
    return JSON.parse(raw).map((i: any) => ({ ...i, id: i.id || String(++iid) }))
  } catch {
    return []
  }
}

watch(() => props.params, (v) => { paramsItems.value = parseKv(v) }, { immediate: true })
watch(() => props.headers, (v) => { headersItems.value = parseKv(v) }, { immediate: true })
watch(() => props.bodyType, (v) => { bodyType.value = v })
watch(() => props.bodyData, (v) => { bodyData.value = v })
watch(() => props.authData, (v) => { authData.value = v })

function sync() {
  emit('update:headers', JSON.stringify(headersItems.value))
  emit('update:params', JSON.stringify(paramsItems.value))
  emit('update:bodyType', bodyType.value)
  emit('update:bodyData', bodyData.value)
  emit('update:authData', authData.value)
}

function onHeadersChange(items: KvItem[]) {
  headersItems.value = items
  tabs[1].count = items.filter(i => i.enabled && i.key).length
  sync()
}

function onParamsChange(items: KvItem[]) {
  paramsItems.value = items
  tabs[0].count = items.filter(i => i.enabled && i.key).length
  sync()
}
</script>

<style scoped>
.request-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.sub-tabs {
  display: flex;
  padding: 0 10px;
  gap: 4px;
  border-bottom: 1px solid #eee;
  background: #fafafa;
}
.sub-content {
  flex: 1;
  overflow-y: auto;
}
.tab-count {
  font-size: 10px;
  background: #eee;
  padding: 0 5px;
  border-radius: 10px;
  margin-left: 3px;
}
</style>
