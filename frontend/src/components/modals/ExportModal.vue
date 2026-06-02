<template>
  <n-modal :show="show" preset="card" title="导出" style="width: 500px" :mask-closable="false" @update:show="emit('update:show', $event)">
    <template v-if="status === 'result'">
      <n-result :status="resultType" :title="resultTitle" :description="resultDesc">
        <template #footer>
          <n-button @click="onReset">继续导出</n-button>
          <n-button type="primary" @click="onClose">完成</n-button>
        </template>
      </n-result>
    </template>

    <template v-else>
      <n-form label-placement="left" label-width="80">
        <n-form-item label="导出范围">
          <n-select v-model:value="scope" :options="scopeOptions" placeholder="选择导出范围" />
        </n-form-item>
        <n-form-item label="导出格式">
          <n-select v-model:value="format" :options="formatOptions" placeholder="选择导出格式" />
        </n-form-item>
        <n-form-item label="保存路径">
          <n-input-group>
            <n-input v-model:value="filePath" placeholder="选择保存路径" :disabled="status === 'exporting'" />
            <n-button :disabled="status === 'exporting'" @click="onBrowse">浏览</n-button>
          </n-input-group>
        </n-form-item>
      </n-form>
      <div class="modal-footer">
        <n-button @click="onClose" :disabled="status === 'exporting'">取消</n-button>
        <n-button type="primary" :loading="status === 'exporting'" :disabled="!canExport" @click="onExport">导出</n-button>
      </div>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { NModal, NForm, NFormItem, NSelect, NInput, NInputGroup, NButton, NResult } from 'naive-ui'
import { ExportPostman } from '../../../wailsjs/go/main/App'

interface Props { show: boolean; projectId: number | null }
const props = defineProps<Props>()
const emit = defineEmits<{ 'update:show': [value: boolean]; exported: [] }>()

type Status = 'idle' | 'exporting' | 'result'
const status = ref<Status>('idle')
const scope = ref('project')
const format = ref('postman')
const filePath = ref('')
const resultType = ref<'success' | 'error'>('success')
const resultTitle = ref('')
const resultDesc = ref('')

const scopeOptions = [
  { label: '整个项目', value: 'project' },
  { label: '指定集合', value: 'collection', disabled: true },
]
const formatOptions = [
  { label: 'Postman Collection', value: 'postman' },
  { label: 'Paw 自有格式', value: 'paw', disabled: true },
]
const canExport = computed(() => props.projectId !== null)

function onBrowse() {}
function onReset() { status.value = 'idle'; filePath.value = ''; resultTitle.value = ''; resultDesc.value = '' }
function onClose() { emit('update:show', false); status.value = 'idle'; filePath.value = '' }

async function onExport() {
  if (!props.projectId) return
  status.value = 'exporting'
  try {
    const result = await ExportPostman(props.projectId)
    resultType.value = 'success'; resultTitle.value = '导出成功'
    resultDesc.value = '文件已保存至: ' + result
    status.value = 'result'; emit('exported')
  } catch (err: any) {
    resultType.value = 'error'; resultTitle.value = '导出失败'
    resultDesc.value = err?.message || err?.toString() || '未知错误'
    status.value = 'result'
  }
}
</script>

<style scoped>
.modal-footer { display: flex; justify-content: flex-end; gap: 10px; margin-top: 24px; }
</style>
