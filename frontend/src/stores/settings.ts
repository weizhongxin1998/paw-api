import { defineStore } from 'pinia'
import { reactive } from 'vue'

interface Settings {
  timeout: number
  followRedirects: boolean
  maxRedirects: number
  sslVerify: boolean
  theme: string
  accentColor: string
  fontSize: number
}

export const useSettingsStore = defineStore('settings', () => {
  const settings = reactive<Settings>({
    timeout: 30,
    followRedirects: true,
    maxRedirects: 10,
    sslVerify: true,
    theme: 'dark',
    accentColor: 'green',
    fontSize: 13,
  })

  return { settings }
})
