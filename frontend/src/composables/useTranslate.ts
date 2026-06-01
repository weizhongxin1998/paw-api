import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

export function useTranslate() {
  const { t, locale } = useI18n()
  const $t = computed(() => {
    locale.value
    return (key: string) => t(key)
  })
  return { $t, locale, t }
}
