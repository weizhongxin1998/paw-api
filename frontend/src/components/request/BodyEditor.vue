<template>
  <div class="body-editor">
    <n-select
      v-model:value="bodyType"
      :options="bodyTypeOptions"
      size="small"
      class="type-select"
    />

    <div v-if="bodyType === 'none'" class="body-empty">
      <span class="hint">此请求没有 body</span>
    </div>

    <KeyValueTable
      v-else-if="bodyType === 'form-data'"
      :items="formDataItems"
      :show-type="true"
      @update:items="onFormDataChange"
    />

    <KeyValueTable
      v-else-if="bodyType === 'x-www-form-urlencoded'"
      :items="urlEncodedItems"
      @update:items="onUrlEncodedChange"
    />

    <div v-else-if="bodyType === 'raw'" class="raw-editor">
      <n-select
        v-model:value="rawSubType"
        :options="rawSubOptions"
        size="tiny"
        class="sub-type-select"
      />
      <n-input
        type="textarea"
        :value="rawContent"
        :rows="10"
        :placeholder="rawPlaceholder"
        class="raw-input"
        @update:value="onRawChange"
      />
      <n-button v-if="rawSubType === 'json'" size="tiny" @click="beautify">Beautify</n-button>
    </div>

    <div v-else-if="bodyType === 'binary'" class="body-empty">
      <n-button size="small" @click="onSelectFile">选择文件</n-button>
      <span v-if="binaryFileName" class="file-name">{{ binaryFileName }}</span>
      <span v-else class="hint">未选择文件</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { NSelect, NInput, NButton } from 'naive-ui'
import KeyValueTable from '../shared/KeyValueTable.vue'
import type { KvItem, RequestBody } from '../../types/request'

const props = defineProps<{
  bodyType: string
  bodyData: string
}>()

const emit = defineEmits<{
  (e: 'update:bodyType', v: string): void
  (e: 'update:bodyData', v: string): void
}>()

const bodyTypeOptions = [
  { label: 'none', value: 'none' },
  { label: 'form-data', value: 'form-data' },
  { label: 'x-www-form-urlencoded', value: 'x-www-form-urlencoded' },
  { label: 'raw', value: 'raw' },
  { label: 'binary', value: 'binary' },
]

const rawSubOptions = [
  { label: 'JSON', value: 'json' },
  { label: 'XML', value: 'xml' },
  { label: 'HTML', value: 'html' },
  { label: 'Text', value: 'text' },
  { label: 'JavaScript', value: 'javascript' },
]

const bodyType = ref(props.bodyType)
const formDataItems = ref<KvItem[]>([])
const urlEncodedItems = ref<KvItem[]>([])
const rawSubType = ref('json')
const rawContent = ref('')
const binaryFileName = ref('')

const rawPlaceholder = ref('输入请求体...')

watch(() => props.bodyType, (v) => { bodyType.value = v })
watch(bodyType, (v) => emit('update:bodyType', v))

function onFormDataChange(items: KvItem[]) {
  formDataItems.value = items
  syncBody('form-data', JSON.stringify(items))
}

function onUrlEncodedChange(items: KvItem[]) {
  urlEncodedItems.value = items
  syncBody('x-www-form-urlencoded', JSON.stringify(items))
}

function onRawChange(v: string) {
  rawContent.value = v
  syncBody('raw', JSON.stringify({ subType: rawSubType.value, content: v }))
}

function syncBody(type: string, data: string) {
  emit('update:bodyData', data)
}

function beautify() {
  try {
    const obj = JSON.parse(rawContent.value)
    rawContent.value = JSON.stringify(obj, null, 2)
    syncBody('raw', JSON.stringify({ subType: rawSubType.value, content: rawContent.value }))
  } catch { /* ignore */ }
}

function onSelectFile() {
  binaryFileName.value = '未实现'
}
</script>

<style scoped>
.body-editor {
  padding: 8px;
}
.type-select {
  margin-bottom: 8px;
}
.body-empty {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 20px;
}
.hint {
  color: #aaa;
  font-size: 13px;
}
.raw-input {
  font-family: 'SF Mono', 'Consolas', monospace;
  font-size: 13px;
}
.file-name {
  color: #2080f0;
  font-size: 13px;
}
.raw-editor {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.sub-type-select {
  width: 120px;
}
.raw-input {
  font-family: 'SF Mono', 'Consolas', monospace;
  font-size: 12px;
}
.file-name {
  color: #2080f0;
  font-size: 12px;
}
</style>
