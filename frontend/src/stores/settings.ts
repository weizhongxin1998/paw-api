import { defineStore } from 'pinia'
import { reactive, watch } from 'vue'
import type { GlobalThemeOverrides } from 'naive-ui'
import { setI18nLocale, type Locale } from '../i18n'

export interface Settings {
  timeout: number
  followRedirects: boolean
  maxRedirects: number
  sslVerify: boolean
  theme: string
  accentColor: string
  fontSize: number
  fontFamily: string
  locale: Locale
}

interface AccentPalette {
  primary: string
  hover: string
  pressed: string
  soft: string
  glow: string
  glowStrong: string
  text: string
  borderFocus: string
}

const ACCENT_COLOR_MAP: Record<string, { dark: AccentPalette; light: AccentPalette }> = {
  green: {
    dark: {
      primary: '#00e05a',
      hover: '#00ff66',
      pressed: '#00b84a',
      soft: 'rgba(0,224,90,0.07)',
      glow: 'rgba(0,224,90,0.14)',
      glowStrong: 'rgba(0,224,90,0.28)',
      text: '#00e05a',
      borderFocus: '#00e05a',
    },
    light: {
      primary: '#059669',
      hover: '#047857',
      pressed: '#065f46',
      soft: 'rgba(5,150,105,0.08)',
      glow: 'rgba(5,150,105,0.12)',
      glowStrong: 'rgba(5,150,105,0.2)',
      text: '#065f46',
      borderFocus: '#059669',
    },
  },
  blue: {
    dark: {
      primary: '#3b82f6',
      hover: '#60a5fa',
      pressed: '#2563eb',
      soft: 'rgba(59,130,246,0.07)',
      glow: 'rgba(59,130,246,0.14)',
      glowStrong: 'rgba(59,130,246,0.28)',
      text: '#3b82f6',
      borderFocus: '#3b82f6',
    },
    light: {
      primary: '#2563eb',
      hover: '#1d4ed8',
      pressed: '#1e40af',
      soft: 'rgba(37,99,235,0.08)',
      glow: 'rgba(37,99,235,0.12)',
      glowStrong: 'rgba(37,99,235,0.2)',
      text: '#1d4ed8',
      borderFocus: '#2563eb',
    },
  },
  purple: {
    dark: {
      primary: '#a855f7',
      hover: '#c084fc',
      pressed: '#9333ea',
      soft: 'rgba(168,85,247,0.07)',
      glow: 'rgba(168,85,247,0.14)',
      glowStrong: 'rgba(168,85,247,0.28)',
      text: '#a855f7',
      borderFocus: '#a855f7',
    },
    light: {
      primary: '#7c3aed',
      hover: '#6d28d9',
      pressed: '#5b21b6',
      soft: 'rgba(124,58,237,0.08)',
      glow: 'rgba(124,58,237,0.12)',
      glowStrong: 'rgba(124,58,237,0.2)',
      text: '#6d28d9',
      borderFocus: '#7c3aed',
    },
  },
}

function getAccentPalette(accentColor: string, dark: boolean): AccentPalette {
  const entry = ACCENT_COLOR_MAP[accentColor]
  if (!entry) return ACCENT_COLOR_MAP.green[dark ? 'dark' : 'light']
  return dark ? entry.dark : entry.light
}

function hexToRgba(hex: string, alpha: number): string {
  const r = parseInt(hex.slice(1, 3), 16)
  const g = parseInt(hex.slice(3, 5), 16)
  const b = parseInt(hex.slice(5, 7), 16)
  return `rgba(${r},${g},${b},${alpha})`
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
  const accent = getAccentPalette(settings.accentColor, dark)

  const ns = (n: number) => Math.round(n * ratio) + 'px'
  const borderPrimary = `1px solid ${accent.primary}`
  const borderFocus = `1px solid ${accent.primary}`

  return {
    common: {
      primaryColor: accent.primary,
      primaryColorHover: accent.hover,
      primaryColorPressed: accent.pressed,
      primaryColorSuppl: accent.primary,
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
      textColorGhost: accent.primary,
      textColorGhostHover: accent.hover,
      textColorGhostPressed: accent.pressed,
      border: dark ? '1px solid #27272a' : '1px solid #d4d4d1',
      borderHover: dark ? '1px solid #3f3f46' : '1px solid #a8a8a4',
      borderFocus,
      borderPrimary,
      borderPrimaryHover: `1px solid ${accent.hover}`,
      borderPrimaryPressed: `1px solid ${accent.pressed}`,
      color: dark ? '#18181b' : '#f4f4f3',
      colorHover: dark ? '#1e1e22' : '#ededeb',
      colorPressed: dark ? '#26262b' : '#e4e4e2',
      colorPrimary: accent.primary,
      colorPrimaryHover: accent.hover,
      colorPrimaryPressed: accent.pressed,
      borderRadiusSmall: '5px',
      borderRadiusMedium: '8px',
    },
    Input: {
      border: dark ? '1px solid #27272a' : '1px solid #d4d4d1',
      borderHover: dark ? '1px solid #3f3f46' : '1px solid #a8a8a4',
      borderFocus,
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
          optionTextColorActive: accent.primary,
        },
      },
    },
    Checkbox: {
      colorChecked: accent.primary,
      borderChecked: accent.primary,
      border: dark ? '1px solid #3f3f46' : '1px solid #bbb',
      checkMarkColor: dark ? '#0a0a0b' : '#fff',
    },
    Dropdown: {
      color: dark ? '#18181b' : '#ffffff',
      dividerColor: dark ? '#27272a' : '#e4e4e2',
      optionTextColor: dark ? '#b4b4bc' : '#3f3f46',
      optionColorActive: hexToRgba(accent.primary, 0.06),
      optionTextColorActive: accent.primary,
      borderRadius: '8px',
    },
    Modal: {
      color: dark ? '#111113' : '#ffffff',
      textColor: dark ? '#ededf0' : '#18181b',
      titleTextColor: dark ? '#ededf0' : '#18181b',
      borderRadius: '12px',
    },
    Tabs: {
      tabTextColorActiveLine: accent.primary,
      tabTextColorActiveBar: accent.primary,
      barColor: accent.primary,
    },
    Spin: { color: accent.primary },
    Result: {
      titleTextColor: dark ? '#ededf0' : '#18181b',
      textColor: dark ? '#b4b4bc' : '#3f3f46',
    },
  }
}

export function applyAccentToDOM(accentColor: string, dark: boolean) {
  const accent = getAccentPalette(accentColor, dark)
  const root = document.documentElement
  root.style.setProperty('--accent', accent.primary)
  root.style.setProperty('--accent-hover', accent.hover)
  root.style.setProperty('--accent-pressed', accent.pressed)
  root.style.setProperty('--accent-soft', accent.soft)
  root.style.setProperty('--accent-glow', accent.glow)
  root.style.setProperty('--accent-glow-strong', accent.glowStrong)
  root.style.setProperty('--accent-text', accent.text)
  root.style.setProperty('--border-focus', accent.borderFocus)
}

export const useSettingsStore = defineStore('settings', () => {
  const settings = reactive<Settings>({
    timeout: 30,
    followRedirects: true,
    maxRedirects: 10,
    sslVerify: true,
    theme: 'light',
    accentColor: 'green',
    fontSize: 14,
    fontFamily: 'JetBrains Mono',
    locale: 'zh-CN',
  })

  watch(() => [settings.fontSize, settings.fontFamily] as const, () => {
    applySettingsToDOM(settings)
  }, { immediate: true })

  watch(() => settings.theme, (theme) => {
    document.documentElement.classList.toggle('theme-light', theme === 'light')
    applyAccentToDOM(settings.accentColor, theme === 'dark')
  }, { immediate: true })

  watch(() => settings.accentColor, (accentColor) => {
    applyAccentToDOM(accentColor, settings.theme === 'dark')
  })

  watch(() => settings.locale, (locale) => {
    setI18nLocale(locale)
  }, { immediate: true })

  return { settings, FONT_FAMILIES }
})
