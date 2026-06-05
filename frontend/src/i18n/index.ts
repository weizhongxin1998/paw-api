import { createI18n } from 'vue-i18n'
import zhCN from './locales/zh-CN.json'
import enUS from './locales/en-US.json'

export type Locale = 'zh-CN' | 'en-US'

const i18n = createI18n({
  legacy: false,
  locale: 'zh-CN',
  fallbackLocale: 'zh-CN',
  messages: {
    'zh-CN': zhCN,
    'en-US': enUS,
  },
})

export default i18n

/**
 * Switch locale at runtime.
 * Call this from the settings store or anywhere else.
 */
export function setI18nLocale(locale: Locale) {
  (i18n.global.locale as unknown as { value: Locale }).value = locale
}

/**
 * Get the current locale value.
 */
export function getI18nLocale(): Locale {
  return (i18n.global.locale as unknown as { value: Locale }).value
}

/**
 * Shorthand for i18n.global.t
 */
export function t(key: string, params?: Record<string, string | number>): string {
  return i18n.global.t(key, params ?? {})
}
