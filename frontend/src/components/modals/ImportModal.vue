<template>
  <n-modal
    :show="show"
    preset="card"
    title="导入"
    style="width: 500px"
    :mask-closable="false"
    @update:show="emit('update:show', $event)"
  >
    <template v-if="status === 'result'">
      <n-result
        :status="resultType"
        :title="resultTitle"
        :description="resultDesc"
      >
        <template #footer>
          <n-button @click="onReset">继续导入</n-button>
          <n-button type="primary" @click="onClose">完成</n-button>
        </template>
      </n-result>
    </template>

    <template v-else>
      <n-form label-placement="left" label-width="80">
        <n-form-item label="导入格式">
          <n-select
            v-model:value="format"
            :options="formatOptions"
            placeholder="选择导入格式"
          />
        </n-form-item>

        <n-form-item label="文件路径">
          <n-input-group>
            <n-input
              v-model:value="filePath"
              placeholder="选择要导入的文件路径"
              :disabled="status === 'importing'"
            />
            <n-button
              :disabled="status === 'importing'"
              @click="onBrowse"
            >
              浏览
            </n-button>
          </n-input-group>
        </n-form-item>
      </n-form>

      <div class="modal-footer">
        <n-button @click="onClose" :disabled="status === 'importing'">取消</n-button>
        <n-button
          type="primary"
          :loading="status === 'importing'"
          :disabled="!canImport"
          @click="onImport"
        >
          导入
        </n-button>
      </div>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  NModal,
  NForm,
  NFormItem,
  NSelect,
  NInput,
  NInputGroup,
  NButton,
  NResult,
} from 'naive-ui'
import { ImportPostman } from '../../../wailsjs/go/main/App'

interface Props {
  show: boolean
  projectId: number | null
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:show': [value: boolean]
  imported: []
}>()

type Status = 'idle' | 'importing' | 'result'

const status = ref<Status>('idle')
const format = ref('postman')
const filePath = ref('')
const resultType = ref<'success' | 'error'>('success')
const resultTitle = ref('')
const resultDesc = ref('')

const formatOptions = [
  { label: 'Postman Collection', value: 'postman' },
  { label: 'OpenAPI 3.x', value: 'openapi3', disabled: true },
  { label: 'Swagger 2.x', value: 'swagger2', disabled: true },
  { label: 'Paw 自有格式', value: 'paw', disabled: true },
]

const canImport = computed(() => {
  return props.projectId !== null && filePath.value.trim() !== ''
})

function onBrowse() {
  // TODO: use Wails file dialog
}

function onReset() {
  status.value = 'idle'
  filePath.value = ''
  resultTitle.value = ''
  resultDesc.value = ''
}

function onClose() {
  emit('update:show', false)
  status.value = 'idle'
  filePath.value = ''
}

async function onImport() {
  if (!props.projectId || !filePath.value.trim()) return

  status.value = 'importing'
  try {
    const result = await ImportPostman(props.projectId, filePath.value.trim())
    resultType.value = 'success'
    resultTitle.value = '导入成功'
    resultDesc.value = `成功导入 ${result.collections} 个集合, ${result.requests} 个请求`
    status.value = 'result'
    emit('imported')
  } catch (err: any) {
    resultType.value = 'error'
    resultTitle.value = '导入失败'
    resultDesc.value = err?.message || err?.toString() || '未知错误'
    status.value = 'result'
  }
}
</script>

<style scoped>
.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 24px;
}
</style>
