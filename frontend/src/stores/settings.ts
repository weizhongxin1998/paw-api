import { defineStore } from 'pinia'
import { reactive, watch } from 'vue'
import type { GlobalThemeOverrides } from 'naive-ui'

export interface Settings {
  timeout: number
  followRedirects: boolean
  maxRedirects: number
  sslVerify: boolean
  theme: string
  accentColor: string
  fontSize: number
  fontFamily: string
}

const FONT_FAMILIES = [
  { label: 'JetBrains Mono', value: 'JetBrains Mono' },
  { label: 'Cascadia Code', value: 'Cascadia Code' },
  { label: 'Fira Code', value: 'Fira Code' },
  { label: 'Consolas', value: 'Consolas' },
  { label: 'SF Mono', value: 'SF Mono' },
  { label: 'Source Code Pro', value: 'Source Code Pro' },
  { label: 'IBM Plex Mono', value: 'IBM Plex Mono' },
]

function buildFontStack(family: string): string {
  return `'${family}', 'Cascadia Code', 'Fira Code', 'SF Mono', 'Consolas', monospace`
}

export function applySettingsToDOM(settings: Settings) {
  const root = document.documentElement
  const base = settings.fontSize
  const scale = base / 13
  const fam = buildFontStack(settings.fontFamily)

  root.style.fontSize = base + 'px'
  root.style.fontFamily = fam
  root.style.setProperty('--fs-2xs', Math.round(9 * scale) + 'px')
  root.style.setProperty('--fs-xs', Math.round(10 * scale) + 'px')
  root.style.setProperty('--fs-sm', Math.round(11 * scale) + 'px')
  root.style.setProperty('--fs-base', base + 'px')
  root.style.setProperty('--fs-md', Math.round(14 * scale) + 'px')
  root.style.setProperty('--fs-lg', Math.round(16 * scale) + 'px')
  root.style.setProperty('--fs-xl', Math.round(18 * scale) + 'px')
  root.style.setProperty('--fs-2xl', Math.round(22 * scale) + 'px')
  root.style.setProperty('--font-family', fam)
  root.style.setProperty('--font-mono', fam)
}

export function buildNaiveOverrides(settings: Settings, dark: boolean): GlobalThemeOverrides {
  const base = settings.fontSize
  const fam = buildFontStack(settings.fontFamily)
  const ratio = base / 13

  const ns = (n: number) => Math.round(n * ratio) + 'px'

  return {
    common: {
      primaryColor: dark ? '#00e05a' : '#009944',
      primaryColorHover: dark ? '#00ff66' : '#007a33',
      primaryColorPressed: dark ? '#00b84a' : '#006b2a',
      primaryColorSuppl: dark ? '#00e05a' : '#009944',
      bodyColor: dark ? '#0d0d0d' : '#f5f3f1',
      cardColor: dark ? '#141414' : '#ffffff',
      modalColor: dark ? '#141414' : '#ffffff',
      popoverColor: dark ? '#1a1a1a' : '#f2f2ef',
      borderColor: dark ? '#2a2a2a' : '#c8c8c2',
      dividerColor: dark ? '#2a2a2a' : '#c8c8c2',
      borderRadius: '4px',
      borderRadiusSmall: '3px',
      fontFamily: fam,
      fontFamilyMono: fam,
      fontWeight: '400',
      fontWeightStrong: '600',
      fontSize: ns(base),
      fontSizeMini: ns(11),
      fontSizeTiny: ns(11),
      fontSizeSmall: ns(13),
      fontSizeMedium: ns(14),
      fontSizeLarge: ns(16),
      fontSizeHuge: ns(18),
      lineHeight: '1.6',
      textColor1: dark ? '#e0e0e0' : '#1a1a18',
      textColor2: dark ? '#b0b0b0' : '#444442',
      textColor3: dark ? '#707070' : '#666662',
      placeholderColor: dark ? '#505050' : '#888884',
      inputColor: dark ? '#141414' : '#ffffff',
      scrollbarColor: dark ? '#333' : '#ccc',
      scrollbarColorHover: dark ? '#444' : '#aaa',
    },
    Button: {
      textColor: dark ? '#00e05a' : '#009944',
      textColorHover: dark ? '#00ff66' : '#007a33',
      textColorPressed: dark ? '#00b84a' : '#006b2a',
      border: dark ? '1px solid #2a2a2a' : '1px solid #c8c8c2',
      borderHover: dark ? '1px solid #3a3a3a' : '1px solid #a8a8a2',
      borderFocus: dark ? '1px solid #00e05a' : '1px solid #009944',
      color: dark ? '#1a1a1a' : '#f2f2ef',
      colorHover: dark ? '#252525' : '#ebebe8',
      colorPressed: dark ? '#303030' : '#e3e3e0',
      borderRadiusSmall: '3px',
      borderRadiusMedium: '4px',
    },
    Input: {
      border: dark ? '1px solid #2a2a2a' : '1px solid #c8c8c2',
      borderHover: dark ? '1px solid #3a3a3a' : '1px solid #a8a8a2',
      borderFocus: dark ? '1px solid #00e05a' : '1px solid #009944',
      borderRadius: '4px',
      color: dark ? '#141414' : '#ffffff',
      textColor: dark ? '#e0e0e0' : '#1a1a18',
      placeholderColor: dark ? '#505050' : '#888884',
      lineHeight: '1.6',
    },
    Select: {
      peers: { InternalSelection: { textColor: dark ? '#e0e0e0' : '#1a1a18' } },
    },
    Checkbox: {
      colorChecked: dark ? '#00e05a' : '#009944',
      borderChecked: dark ? '#00e05a' : '#009944',
      border: dark ? '1px solid #3a3a3a' : '1px solid #bbb',
      checkMarkColor: dark ? '#0d0d0d' : '#fff',
    },
    Dropdown: {
      color: dark ? '#1a1a1a' : '#ffffff',
      dividerColor: dark ? '#2a2a2a' : '#e0e0da',
      optionColorActive: dark ? '#1a2a1a' : '#e6f7ec',
      optionTextColorActive: dark ? '#00e05a' : '#009944',
    },
    Modal: {
      color: dark ? '#141414' : '#ffffff',
      textColor: dark ? '#e0e0e0' : '#1a1a18',
      titleTextColor: dark ? '#e0e0e0' : '#1a1a18',
    },
    Tabs: {
      tabTextColorActiveLine: dark ? '#00e05a' : '#009944',
      tabTextColorActiveBar: dark ? '#00e05a' : '#009944',
      barColor: dark ? '#00e05a' : '#009944',
    },
    Spin: { color: dark ? '#00e05a' : '#009944' },
    Result: {
      titleTextColor: dark ? '#e0e0e0' : '#1a1a18',
      textColor: dark ? '#b0b0b0' : '#444442',
    },
  }
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
    fontFamily: 'JetBrains Mono',
  })

  watch(() => [settings.fontSize, settings.fontFamily] as const, () => {
    applySettingsToDOM(settings)
  }, { immediate: true })

  return { settings, FONT_FAMILIES }
})
