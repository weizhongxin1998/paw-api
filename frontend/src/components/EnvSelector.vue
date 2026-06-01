<script lang="ts" setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { NSelect, NButton, NIcon, NSpace } from 'naive-ui'
import { Settings } from '@vicons/ionicons5'
import { useEnvironmentStore } from '../stores/environment'
import { useProjectStore } from '../stores/project'

const { t } = useI18n()
const envStore = useEnvironmentStore()
const projectStore = useProjectStore()

const emit = defineEmits<{
  (e: 'manage'): void
}>()

const envOptions = computed(() => envStore.environments.map(e => ({
  label: e.name,
  value: e.id,
})))

async function handleChange(envId: string | null) {
  if (!envId || !projectStore.currentProject) return
  try {
    const { SetActiveEnvironment } = await import('../../wailsjs/go/handlers/EnvironmentHandler')
    const env = await SetActiveEnvironment(envId, projectStore.currentProject.id)
    envStore.setActiveEnvironment(env)
  } catch (e: any) {
    console.error('Failed to set active environment', e)
  }
}
</script>

<template>
  <NSpace align="center" size="small">
    <NSelect
      :options="envOptions"
      :value="envStore.activeEnvironment?.id ?? null"
      :placeholder="t('env.noEnv')"
      size="tiny"
      style="width: 140px"
      clearable
      @update:value="handleChange"
    />
    <NButton quaternary size="tiny" @click="emit('manage')">
      <template #icon><NIcon><Settings /></NIcon></template>
    </NButton>
  </NSpace>
</template>
