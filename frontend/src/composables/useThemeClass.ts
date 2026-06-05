import { ref, computed, onMounted } from 'vue'

const THEME_CLASSES = ['theme-light', 'theme-warm', 'theme-nord', 'theme-catppuccin', 'theme-neon']

export function useThemeClass(modalPrefix: string) {
  const themeClass = ref('')

  const modalClass = computed(() => {
    const parts = [modalPrefix]
    if (themeClass.value) parts.push(themeClass.value)
    return parts.join(' ')
  })

  onMounted(() => {
    const check = () => {
      for (const cls of THEME_CLASSES) {
        if (document.documentElement.classList.contains(cls)) {
          themeClass.value = cls
          return
        }
      }
      themeClass.value = ''
    }
    check()
    const observer = new MutationObserver(check)
    observer.observe(document.documentElement, { attributes: true, attributeFilter: ['class'] })
  })

  return { themeClass, modalClass }
}
