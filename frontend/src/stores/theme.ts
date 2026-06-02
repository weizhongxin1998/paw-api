import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useThemeStore = defineStore('theme', () => {
  const mode = ref<'light' | 'dark'>('light')
  const color = ref('green')

  return { mode, color }
})
