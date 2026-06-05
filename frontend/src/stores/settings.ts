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
  root.style.setProperty('--fs-2xs', Math.round(10 * scale) + 'px')
  root.style.setProperty('--fs-xs', Math.round(11.5 * scale) + 'px')
  root.style.setProperty('--fs-sm', Math.round(12.5 * scale) + 'px')
  root.style.setProperty('--fs-base', base + 'px')
  root.style.setProperty('--fs-md', Math.round(14.5 * scale) + 'px')
  root.style.setProperty('--fs-lg', Math.round(16.5 * scale) + 'px')
  root.style.setProperty('--fs-xl', Math.round(19 * scale) + 'px')
  root.style.setProperty('--fs-2xl', Math.round(24 * scale) + 'px')
  root.style.setProperty('--fs-3xl', Math.round(32 * scale) + 'px')
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
      primaryColor: dark ? '#00e05a' : '#059669',
      primaryColorHover: dark ? '#00ff66' : '#047857',
      primaryColorPressed: dark ? '#00b84a' : '#065f46',
      primaryColorSuppl: dark ? '#00e05a' : '#059669',
      bodyColor: dark ? '#0a0a0b' : '#f8f8f7',
      cardColor: dark ? '#111113' : '#ffffff',
      modalColor: dark ? '#111113' : '#ffffff',
      popoverColor: dark ? '#18181b' : '#f4f4f3',
      borderColor: dark ? '#27272a' : '#d4d4d1',
      dividerColor: dark ? '#27272a' : '#d4d4d1',
      borderRadius: '8px',
      borderRadiusSmall: '5px',
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
      textColor1: dark ? '#ededf0' : '#18181b',
      textColor2: dark ? '#b4b4bc' : '#3f3f46',
      textColor3: dark ? '#8a8a95' : '#52525b',
      placeholderColor: dark ? '#6b6b76' : '#71717a',
      inputColor: dark ? '#111113' : '#ffffff',
      scrollbarColor: 'rgba(128,128,128,0.15)',
      scrollbarColorHover: 'rgba(128,128,128,0.25)',
    },
    Button: {
      textColor: dark ? '#ededf0' : '#18181b',
      textColorHover: dark ? '#ededf0' : '#18181b',
      textColorPressed: dark ? '#ededf0' : '#18181b',
      textColorPrimary: dark ? '#0a0a0b' : '#ffffff',
      textColorPrimaryHover: dark ? '#0a0a0b' : '#ffffff',
      textColorPrimaryPressed: dark ? '#0a0a0b' : '#ffffff',
      textColorGhost: dark ? '#00e05a' : '#059669',
      textColorGhostHover: dark ? '#00ff66' : '#047857',
      textColorGhostPressed: dark ? '#00b84a' : '#065f46',
      border: dark ? '1px solid #27272a' : '1px solid #d4d4d1',
      borderHover: dark ? '1px solid #3f3f46' : '1px solid #a8a8a4',
      borderFocus: dark ? '1px solid #00e05a' : '1px solid #059669',
      borderPrimary: dark ? '1px solid #00e05a' : '1px solid #059669',
      borderPrimaryHover: dark ? '1px solid #00ff66' : '1px solid #047857',
      borderPrimaryPressed: dark ? '1px solid #00b84a' : '1px solid #065f46',
      color: dark ? '#18181b' : '#f4f4f3',
      colorHover: dark ? '#1e1e22' : '#ededeb',
      colorPressed: dark ? '#26262b' : '#e4e4e2',
      colorPrimary: dark ? '#00e05a' : '#059669',
      colorPrimaryHover: dark ? '#00ff66' : '#047857',
      colorPrimaryPressed: dark ? '#00b84a' : '#065f46',
      borderRadiusSmall: '5px',
      borderRadiusMedium: '8px',
    },
    Input: {
      border: dark ? '1px solid #27272a' : '1px solid #d4d4d1',
      borderHover: dark ? '1px solid #3f3f46' : '1px solid #a8a8a4',
      borderFocus: dark ? '1px solid #00e05a' : '1px solid #059669',
      borderRadius: '8px',
      color: dark ? '#111113' : '#ffffff',
      textColor: dark ? '#ededf0' : '#18181b',
      placeholderColor: dark ? '#6b6b76' : '#71717a',
      lineHeight: '1.6',
    },
    Select: {
      peers: {
        InternalSelection: {
          textColor: dark ? '#ededf0' : '#18181b',
          placeholderColor: dark ? '#6b6b76' : '#71717a',
        },
        InternalSelectMenu: {
          optionTextColor: dark ? '#b4b4bc' : '#3f3f46',
          optionTextColorActive: dark ? '#00e05a' : '#059669',
        },
      },
    },
    Checkbox: {
      colorChecked: dark ? '#00e05a' : '#059669',
      borderChecked: dark ? '#00e05a' : '#059669',
      border: dark ? '1px solid #3f3f46' : '1px solid #bbb',
      checkMarkColor: dark ? '#0a0a0b' : '#fff',
    },
    Dropdown: {
      color: dark ? '#18181b' : '#ffffff',
      dividerColor: dark ? '#27272a' : '#e4e4e2',
      optionTextColor: dark ? '#b4b4bc' : '#3f3f46',
      optionColorActive: dark ? 'rgba(0,224,90,0.06)' : 'rgba(5,150,105,0.06)',
      optionTextColorActive: dark ? '#00e05a' : '#059669',
      borderRadius: '8px',
    },
    Modal: {
      color: dark ? '#111113' : '#ffffff',
      textColor: dark ? '#ededf0' : '#18181b',
      titleTextColor: dark ? '#ededf0' : '#18181b',
      borderRadius: '12px',
    },
    Tabs: {
      tabTextColorActiveLine: dark ? '#00e05a' : '#059669',
      tabTextColorActiveBar: dark ? '#00e05a' : '#059669',
      barColor: dark ? '#00e05a' : '#059669',
    },
    Spin: { color: dark ? '#00e05a' : '#059669' },
    Result: {
      titleTextColor: dark ? '#ededf0' : '#18181b',
      textColor: dark ? '#b4b4bc' : '#3f3f46',
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
