import { createI18n } from 'vue-i18n'
import zhCN from './zh-CN'
import en from './en'

const saved = localStorage.getItem('paw-locale')

export const i18n = createI18n({
  legacy: false,
  locale: (saved === 'en' || saved === 'zh-CN') ? saved : 'zh-CN',
  fallbackLocale: 'zh-CN',
  messages: { 'zh-CN': zhCN, en },
})
