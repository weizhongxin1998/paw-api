<template>
  <n-config-provider :theme-overrides="themeOverrides">
    <n-message-provider>
      <div class="app-container">
        <AppHeader />
        <AppBody />
      </div>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { NConfigProvider, NMessageProvider } from 'naive-ui'
import AppHeader from './components/layout/AppHeader.vue'
import AppBody from './components/layout/AppBody.vue'

const themeOverrides = {
  common: {
    primaryColor: '#18a058',
    primaryColorHover: '#0c7a43',
  },
}
</script>

<style>
html, body, #app {
  margin: 0;
  padding: 0;
  height: 100%;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}
.app-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
}
</style>
