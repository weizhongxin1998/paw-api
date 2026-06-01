<script lang="ts" setup>
import { computed } from 'vue'
import { NConfigProvider, NLayout, NLayoutSider, NLayoutContent, NMessageProvider, NDialogProvider, zhCN, enUS, dateZhCN, dateEnUS } from 'naive-ui'
import { useI18n } from 'vue-i18n'
import { useTheme } from './composables/useTheme'
import AppSidebar from './components/AppSidebar.vue'

const { theme, themeOverrides } = useTheme()
const { locale } = useI18n()

const naiveLocale = computed(() => locale.value === 'zh-CN' ? zhCN : enUS)
const naiveDateLocale = computed(() => locale.value === 'zh-CN' ? dateZhCN : dateEnUS)
</script>

<template>
  <NConfigProvider :theme="theme" :theme-overrides="themeOverrides" :locale="naiveLocale" :date-locale="naiveDateLocale">
    <NMessageProvider>
      <NDialogProvider>
        <NLayout class="app-layout" has-sider :key="locale">
          <NLayoutSider bordered width="280" :native-scrollbar="false">
            <AppSidebar />
          </NLayoutSider>
          <NLayoutContent class="main-content">
            <router-view />
          </NLayoutContent>
        </NLayout>
      </NDialogProvider>
    </NMessageProvider>
  </NConfigProvider>
</template>

<style scoped>
.app-layout { height: 100vh; width: 100vw; }
.main-content { display: flex; flex-direction: column; height: 100vh; overflow: hidden; }
</style>
