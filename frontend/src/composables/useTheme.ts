import { computed, ref, watch } from 'vue'
import { darkTheme } from 'naive-ui'
import type { GlobalTheme, GlobalThemeOverrides } from 'naive-ui'

type ColorMode = 'light' | 'dark'
type ThemeColor = 'green' | 'blue' | 'purple'

const themeOverrides: Record<ThemeColor, GlobalThemeOverrides> = {
  green: {
    common: { primaryColor: '#18a058', primaryColorHover: '#36ad6a', primaryColorPressed: '#0c7a43' }
  },
  blue: {
    common: { primaryColor: '#2080f0', primaryColorHover: '#4098f7', primaryColorPressed: '#1060c0' }
  },
  purple: {
    common: { primaryColor: '#8a63d2', primaryColorHover: '#b794f4', primaryColorPressed: '#6b46c0' }
  }
}

const colorMode = ref<ColorMode>('light')
const themeColor = ref<ThemeColor>('green')

export function useTheme() {
  const theme = computed<GlobalTheme | null>(() => colorMode.value === 'dark' ? darkTheme : null)

  const themeOverridesComputed = computed<GlobalThemeOverrides>(() => themeOverrides[themeColor.value])

  function toggleColorMode() {
    colorMode.value = colorMode.value === 'light' ? 'dark' : 'light'
  }

  function setThemeColor(color: ThemeColor) {
    themeColor.value = color
  }

  function applyBodyTheme(mode: ColorMode) {
    document.body.setAttribute('data-naive-ui-theme', mode)
  }

  watch(colorMode, (mode) => {
    localStorage.setItem('paw-color-mode', mode)
    applyBodyTheme(mode)
  })

  watch(themeColor, (color) => {
    localStorage.setItem('paw-theme-color', color)
  })

  applyBodyTheme(colorMode.value)

  return { theme, themeOverrides: themeOverridesComputed, colorMode, themeColor, toggleColorMode, setThemeColor }
}
