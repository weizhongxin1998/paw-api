import { defineStore } from 'pinia'
import { reactive, watch } from 'vue'
import type { GlobalThemeOverrides } from 'naive-ui'
import { setI18nLocale, type Locale } from '../i18n'

export type Theme = 'dark' | 'light' | 'warm' | 'nord' | 'catppuccin' | 'neon'

export interface Settings {
  timeout: number
  followRedirects: boolean
  maxRedirects: number
  sslVerify: boolean
  theme: Theme
  accentColor: string
  fontSize: number
  fontFamily: string
  locale: Locale
}

const THEME_CLASSES: Record<string, string> = {
  light: 'theme-light',
  warm: 'theme-warm',
  nord: 'theme-nord',
  catppuccin: 'theme-catppuccin',
  neon: 'theme-neon',
}

export function isThemeDark(theme: string): boolean {
  return theme !== 'light'
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

function hexToRgba(hex: string, alpha: number): string {
  const r = parseInt(hex.slice(1, 3), 16)
  const g = parseInt(hex.slice(3, 5), 16)
  const b = parseInt(hex.slice(5, 7), 16)
  return `rgba(${r},${g},${b},${alpha})`
}

const ACCENT_COLOR_MAP: Record<string, { dark: AccentPalette; light: AccentPalette }> = {
  green: {
    dark: { primary: '#00e05a', hover: '#00ff66', pressed: '#00b84a', soft: 'rgba(0,224,90,0.07)', glow: 'rgba(0,224,90,0.14)', glowStrong: 'rgba(0,224,90,0.28)', text: '#00e05a', borderFocus: '#00e05a' },
    light: { primary: '#059669', hover: '#047857', pressed: '#065f46', soft: 'rgba(5,150,105,0.08)', glow: 'rgba(5,150,105,0.12)', glowStrong: 'rgba(5,150,105,0.2)', text: '#065f46', borderFocus: '#059669' },
  },
  blue: {
    dark: { primary: '#3b82f6', hover: '#60a5fa', pressed: '#2563eb', soft: 'rgba(59,130,246,0.07)', glow: 'rgba(59,130,246,0.14)', glowStrong: 'rgba(59,130,246,0.28)', text: '#3b82f6', borderFocus: '#3b82f6' },
    light: { primary: '#2563eb', hover: '#1d4ed8', pressed: '#1e40af', soft: 'rgba(37,99,235,0.08)', glow: 'rgba(37,99,235,0.12)', glowStrong: 'rgba(37,99,235,0.2)', text: '#1d4ed8', borderFocus: '#2563eb' },
  },
  purple: {
    dark: { primary: '#a855f7', hover: '#c084fc', pressed: '#9333ea', soft: 'rgba(168,85,247,0.07)', glow: 'rgba(168,85,247,0.14)', glowStrong: 'rgba(168,85,247,0.28)', text: '#a855f7', borderFocus: '#a855f7' },
    light: { primary: '#7c3aed', hover: '#6d28d9', pressed: '#5b21b6', soft: 'rgba(124,58,237,0.08)', glow: 'rgba(124,58,237,0.12)', glowStrong: 'rgba(124,58,237,0.2)', text: '#6d28d9', borderFocus: '#7c3aed' },
  },
  amber: {
    dark: { primary: '#f0a848', hover: '#f5b85e', pressed: '#d8922f', soft: 'rgba(240,168,72,0.07)', glow: 'rgba(240,168,72,0.14)', glowStrong: 'rgba(240,168,72,0.28)', text: '#f0a848', borderFocus: '#d97706' },
    light: { primary: '#d97706', hover: '#b85c00', pressed: '#92400e', soft: 'rgba(217,119,6,0.08)', glow: 'rgba(217,119,6,0.12)', glowStrong: 'rgba(217,119,6,0.2)', text: '#92400e', borderFocus: '#d97706' },
  },
  teal: {
    dark: { primary: '#5eead4', hover: '#7df0de', pressed: '#4ac0ac', soft: 'rgba(94,234,212,0.07)', glow: 'rgba(94,234,212,0.14)', glowStrong: 'rgba(94,234,212,0.28)', text: '#5eead4', borderFocus: '#5eead4' },
    light: { primary: '#0891b2', hover: '#06b6d4', pressed: '#0e7490', soft: 'rgba(8,145,178,0.08)', glow: 'rgba(8,145,178,0.12)', glowStrong: 'rgba(8,145,178,0.2)', text: '#0e7490', borderFocus: '#0891b2' },
  },
  mauve: {
    dark: { primary: '#c4a0f8', hover: '#d4b8ff', pressed: '#a880d8', soft: 'rgba(196,160,248,0.07)', glow: 'rgba(196,160,248,0.14)', glowStrong: 'rgba(196,160,248,0.28)', text: '#c4a0f8', borderFocus: '#c4a0f8' },
    light: { primary: '#8b5cf6', hover: '#7c3aed', pressed: '#6d28d9', soft: 'rgba(139,92,246,0.08)', glow: 'rgba(139,92,246,0.12)', glowStrong: 'rgba(139,92,246,0.2)', text: '#6d28d9', borderFocus: '#8b5cf6' },
  },
  magenta: {
    dark: { primary: '#ff0088', hover: '#ff3399', pressed: '#cc0066', soft: 'rgba(255,0,136,0.07)', glow: 'rgba(255,0,136,0.14)', glowStrong: 'rgba(255,0,136,0.28)', text: '#ff0088', borderFocus: '#ff0088' },
    light: { primary: '#e11d48', hover: '#f43f5e', pressed: '#be123c', soft: 'rgba(225,29,72,0.08)', glow: 'rgba(225,29,72,0.12)', glowStrong: 'rgba(225,29,72,0.2)', text: '#be123c', borderFocus: '#e11d48' },
  },
}

function getAccentPalette(accentColor: string, dark: boolean): AccentPalette {
  const entry = ACCENT_COLOR_MAP[accentColor]
  if (!entry) return ACCENT_COLOR_MAP.green[dark ? 'dark' : 'light']
  return dark ? entry.dark : entry.light
}

interface ThemeColors {
  bodyColor: string
  cardColor: string
  modalColor: string
  popoverColor: string
  borderColor: string
  dividerColor: string
  textColor1: string
  textColor2: string
  textColor3: string
  placeholderColor: string
  inputColor: string
  buttonColor: string
  buttonColorHover: string
  buttonColorPressed: string
  buttonBorder: string
  buttonBorderHover: string
  inputBorder: string
  inputBorderHover: string
  inputColorValue: string
  inputTextColor: string
  selectTextColor: string
  selectOptionTextColor: string
  checkboxBorder: string
  checkMarkColor: string
  dropdownColor: string
  dropdownDividerColor: string
  dropdownOptionTextColor: string
  modalTextColor: string
  modalTitleTextColor: string
  resultTitleTextColor: string
}

const THEME_COLORS: Record<string, ThemeColors> = {
  light: {
    bodyColor: '#f8f8f7',
    cardColor: '#ffffff',
    modalColor: '#ffffff',
    popoverColor: '#f4f4f3',
    borderColor: '#d4d4d1',
    dividerColor: '#d4d4d1',
    textColor1: '#18181b',
    textColor2: '#3f3f46',
    textColor3: '#52525b',
    placeholderColor: '#71717a',
    inputColor: '#ffffff',
    buttonColor: '#f4f4f3',
    buttonColorHover: '#ededeb',
    buttonColorPressed: '#e4e4e2',
    buttonBorder: '1px solid #d4d4d1',
    buttonBorderHover: '1px solid #a8a8a4',
    inputBorder: '1px solid #d4d4d1',
    inputBorderHover: '1px solid #a8a8a4',
    inputColorValue: '#ffffff',
    inputTextColor: '#18181b',
    selectTextColor: '#18181b',
    selectOptionTextColor: '#3f3f46',
    checkboxBorder: '1px solid #bbb',
    checkMarkColor: '#fff',
    dropdownColor: '#ffffff',
    dropdownDividerColor: '#e4e4e2',
    dropdownOptionTextColor: '#3f3f46',
    modalTextColor: '#18181b',
    modalTitleTextColor: '#18181b',
    resultTitleTextColor: '#18181b',
  },
  dark: {
    bodyColor: '#0a0a0b',
    cardColor: '#111113',
    modalColor: '#111113',
    popoverColor: '#18181b',
    borderColor: '#27272a',
    dividerColor: '#27272a',
    textColor1: '#ededf0',
    textColor2: '#b4b4bc',
    textColor3: '#8a8a95',
    placeholderColor: '#6b6b76',
    inputColor: '#111113',
    buttonColor: '#18181b',
    buttonColorHover: '#1e1e22',
    buttonColorPressed: '#26262b',
    buttonBorder: '1px solid #27272a',
    buttonBorderHover: '1px solid #3f3f46',
    inputBorder: '1px solid #27272a',
    inputBorderHover: '1px solid #3f3f46',
    inputColorValue: '#111113',
    inputTextColor: '#ededf0',
    selectTextColor: '#ededf0',
    selectOptionTextColor: '#b4b4bc',
    checkboxBorder: '1px solid #3f3f46',
    checkMarkColor: '#0a0a0b',
    dropdownColor: '#18181b',
    dropdownDividerColor: '#27272a',
    dropdownOptionTextColor: '#b4b4bc',
    modalTextColor: '#ededf0',
    modalTitleTextColor: '#ededf0',
    resultTitleTextColor: '#ededf0',
  },
  warm: {
    bodyColor: '#1f1a14',
    cardColor: '#262017',
    modalColor: '#262017',
    popoverColor: '#2d261c',
    borderColor: '#3d3228',
    dividerColor: '#3d3228',
    textColor1: '#e6d4b8',
    textColor2: '#b09878',
    textColor3: '#806850',
    placeholderColor: '#6b5a42',
    inputColor: '#262017',
    buttonColor: '#2d261c',
    buttonColorHover: '#332b20',
    buttonColorPressed: '#3a3125',
    buttonBorder: '1px solid #3d3228',
    buttonBorderHover: '1px solid #5c4d3b',
    inputBorder: '1px solid #3d3228',
    inputBorderHover: '1px solid #5c4d3b',
    inputColorValue: '#17130e',
    inputTextColor: '#e6d4b8',
    selectTextColor: '#e6d4b8',
    selectOptionTextColor: '#b09878',
    checkboxBorder: '1px solid #5c4d3b',
    checkMarkColor: '#1f1a14',
    dropdownColor: '#2d261c',
    dropdownDividerColor: '#3d3228',
    dropdownOptionTextColor: '#b09878',
    modalTextColor: '#e6d4b8',
    modalTitleTextColor: '#e6d4b8',
    resultTitleTextColor: '#e6d4b8',
  },
  nord: {
    bodyColor: '#1f2535',
    cardColor: '#242b3c',
    modalColor: '#242b3c',
    popoverColor: '#293148',
    borderColor: '#313b50',
    dividerColor: '#313b50',
    textColor1: '#dde4ee',
    textColor2: '#a0aec0',
    textColor3: '#6b7a8d',
    placeholderColor: '#52607a',
    inputColor: '#242b3c',
    buttonColor: '#293148',
    buttonColorHover: '#2e364d',
    buttonColorPressed: '#333c55',
    buttonBorder: '1px solid #313b50',
    buttonBorderHover: '1px solid #4a5670',
    inputBorder: '1px solid #313b50',
    inputBorderHover: '1px solid #4a5670',
    inputColorValue: '#191e2c',
    inputTextColor: '#dde4ee',
    selectTextColor: '#dde4ee',
    selectOptionTextColor: '#a0aec0',
    checkboxBorder: '1px solid #4a5670',
    checkMarkColor: '#1f2535',
    dropdownColor: '#293148',
    dropdownDividerColor: '#313b50',
    dropdownOptionTextColor: '#a0aec0',
    modalTextColor: '#dde4ee',
    modalTitleTextColor: '#dde4ee',
    resultTitleTextColor: '#dde4ee',
  },
  catppuccin: {
    bodyColor: '#1a1826',
    cardColor: '#201e30',
    modalColor: '#201e30',
    popoverColor: '#27243a',
    borderColor: '#353058',
    dividerColor: '#353058',
    textColor1: '#dedcf0',
    textColor2: '#b8b4d0',
    textColor3: '#7a7698',
    placeholderColor: '#656080',
    inputColor: '#201e30',
    buttonColor: '#27243a',
    buttonColorHover: '#2d2a42',
    buttonColorPressed: '#33304a',
    buttonBorder: '1px solid #353058',
    buttonBorderHover: '1px solid #524a75',
    inputBorder: '1px solid #353058',
    inputBorderHover: '1px solid #524a75',
    inputColorValue: '#13121e',
    inputTextColor: '#dedcf0',
    selectTextColor: '#dedcf0',
    selectOptionTextColor: '#b8b4d0',
    checkboxBorder: '1px solid #524a75',
    checkMarkColor: '#1a1826',
    dropdownColor: '#27243a',
    dropdownDividerColor: '#353058',
    dropdownOptionTextColor: '#b8b4d0',
    modalTextColor: '#dedcf0',
    modalTitleTextColor: '#dedcf0',
    resultTitleTextColor: '#dedcf0',
  },
  neon: {
    bodyColor: '#0a0a0f',
    cardColor: '#0f0f18',
    modalColor: '#0f0f18',
    popoverColor: '#151524',
    borderColor: '#1e1e33',
    dividerColor: '#1e1e33',
    textColor1: '#e6e6f2',
    textColor2: '#a8a8c0',
    textColor3: '#686880',
    placeholderColor: '#525268',
    inputColor: '#0f0f18',
    buttonColor: '#151524',
    buttonColorHover: '#1a1a2c',
    buttonColorPressed: '#202033',
    buttonBorder: '1px solid #1e1e33',
    buttonBorderHover: '1px solid #383858',
    inputBorder: '1px solid #1e1e33',
    inputBorderHover: '1px solid #383858',
    inputColorValue: '#080810',
    inputTextColor: '#e6e6f2',
    selectTextColor: '#e6e6f2',
    selectOptionTextColor: '#a8a8c0',
    checkboxBorder: '1px solid #383858',
    checkMarkColor: '#0a0a0f',
    dropdownColor: '#151524',
    dropdownDividerColor: '#1e1e33',
    dropdownOptionTextColor: '#a8a8c0',
    modalTextColor: '#e6e6f2',
    modalTitleTextColor: '#e6e6f2',
    resultTitleTextColor: '#e6e6f2',
  },
}

function getThemeColors(theme: string): ThemeColors {
  return THEME_COLORS[theme] || THEME_COLORS.dark
}

const FONT_FAMILIES = [
  { label: 'Microsoft YaHei', value: 'Microsoft YaHei' },
  { label: 'Noto Sans SC', value: 'Noto Sans SC' },
  { label: 'Inter', value: 'Inter' },
  { label: 'Segoe UI', value: 'Segoe UI' },
  { label: 'Calibri', value: 'Calibri' },
  { label: 'Trebuchet MS', value: 'Trebuchet MS' },
  { label: 'Georgia', value: 'Georgia' },
]

export function buildFontStack(family: string): string {
  return `'${family}', 'Microsoft YaHei', sans-serif`
}

export function applySettingsToDOM(settings: Settings) {
  const root = document.documentElement
  const base = settings.fontSize
  const scale = base / 13
  const fam = buildFontStack(settings.fontFamily)

  root.style.fontSize = base + 'px'
  root.style.fontFamily = fam
  document.body.style.fontSize = base + 'px'
  document.body.style.fontFamily = fam
  document.body.style.setProperty('--n-font-family', fam)
  document.body.style.setProperty('--n-font-family-mono', fam)
  root.style.setProperty('--n-font-family', fam)
  root.style.setProperty('--n-font-family-mono', fam)
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

const THEME_DEFAULT_ACCENT: Record<string, string> = {
  dark: 'green', light: 'green', warm: 'amber', nord: 'teal',
  catppuccin: 'mauve', neon: 'magenta',
}

export function buildNaiveOverrides(settings: Settings, dark: boolean): GlobalThemeOverrides {
  const base = settings.fontSize
  const fam = buildFontStack(settings.fontFamily)
  const ratio = base / 13
  const accent = getAccentPalette(settings.accentColor, dark)

  const colors = getThemeColors(settings.theme)

  const ns = (n: number) => Math.round(n * ratio) + 'px'
  const borderPrimary = `1px solid ${accent.primary}`
  const borderFocus = `1px solid ${accent.primary}`

  return {
    common: {
      primaryColor: accent.primary,
      primaryColorHover: accent.hover,
      primaryColorPressed: accent.pressed,
      primaryColorSuppl: accent.primary,
      bodyColor: colors.bodyColor,
      cardColor: colors.cardColor,
      modalColor: colors.modalColor,
      popoverColor: colors.popoverColor,
      borderColor: colors.borderColor,
      dividerColor: colors.dividerColor,
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
      textColor1: colors.textColor1,
      textColor2: colors.textColor2,
      textColor3: colors.textColor3,
      placeholderColor: colors.placeholderColor,
      inputColor: colors.inputColor,
      scrollbarColor: 'rgba(128,128,128,0.15)',
      scrollbarColorHover: 'rgba(128,128,128,0.25)',
    },
    Button: {
      textColor: colors.textColor1,
      textColorHover: colors.textColor1,
      textColorPressed: colors.textColor1,
      textColorPrimary: dark ? '#0a0a0b' : '#ffffff',
      textColorPrimaryHover: dark ? '#0a0a0b' : '#ffffff',
      textColorPrimaryPressed: dark ? '#0a0a0b' : '#ffffff',
      textColorGhost: accent.primary,
      textColorGhostHover: accent.hover,
      textColorGhostPressed: accent.pressed,
      border: colors.buttonBorder,
      borderHover: colors.buttonBorderHover,
      borderFocus,
      borderPrimary,
      borderPrimaryHover: `1px solid ${accent.hover}`,
      borderPrimaryPressed: `1px solid ${accent.pressed}`,
      color: colors.buttonColor,
      colorHover: colors.buttonColorHover,
      colorPressed: colors.buttonColorPressed,
      colorPrimary: accent.primary,
      colorPrimaryHover: accent.hover,
      colorPrimaryPressed: accent.pressed,
      borderRadiusSmall: '5px',
      borderRadiusMedium: '8px',
    },
    Input: {
      border: colors.inputBorder,
      borderHover: colors.inputBorderHover,
      borderFocus,
      borderRadius: '8px',
      color: colors.inputColorValue,
      textColor: colors.inputTextColor,
      placeholderColor: colors.placeholderColor,
      lineHeight: '1.6',
    },
    Select: {
      peers: {
        InternalSelection: {
          textColor: colors.selectTextColor,
          placeholderColor: colors.placeholderColor,
        },
        InternalSelectMenu: {
          optionTextColor: colors.selectOptionTextColor,
          optionTextColorActive: accent.primary,
        },
      },
    },
    Checkbox: {
      colorChecked: accent.primary,
      borderChecked: accent.primary,
      border: colors.checkboxBorder,
      checkMarkColor: colors.checkMarkColor,
    },
    Dropdown: {
      color: colors.dropdownColor,
      dividerColor: colors.dropdownDividerColor,
      optionTextColor: colors.dropdownOptionTextColor,
      optionColorActive: hexToRgba(accent.primary, 0.06),
      optionTextColorActive: accent.primary,
      borderRadius: '8px',
    },
    Modal: {
      color: colors.modalColor,
      textColor: colors.modalTextColor,
      titleTextColor: colors.modalTitleTextColor,
      borderRadius: '12px',
    },
    Tabs: {
      tabTextColorActiveLine: accent.primary,
      tabTextColorActiveBar: accent.primary,
      barColor: accent.primary,
    },
    Spin: { color: accent.primary },
    Result: {
      titleTextColor: colors.resultTitleTextColor,
      textColor: colors.textColor2,
    },
  }
}

function applyThemeClass(theme: string) {
  const html = document.documentElement
  for (const cls of Object.values(THEME_CLASSES)) {
    html.classList.remove(cls)
  }
  const themeCls = THEME_CLASSES[theme]
  if (themeCls) {
    html.classList.add(themeCls)
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
    fontFamily: 'Microsoft YaHei',
    locale: 'zh-CN',
  })

  watch(() => [settings.fontSize, settings.fontFamily] as const, () => {
    applySettingsToDOM(settings)
  }, { immediate: true })

  watch(() => settings.theme, (theme) => {
    applyThemeClass(theme)
    settings.accentColor = THEME_DEFAULT_ACCENT[theme] || 'green'
    applyAccentToDOM(settings.accentColor, isThemeDark(theme))
  }, { immediate: true })

  watch(() => settings.accentColor, (accentColor) => {
    applyAccentToDOM(accentColor, isThemeDark(settings.theme))
  })

  watch(() => settings.locale, (locale) => {
    setI18nLocale(locale)
  }, { immediate: true })

  return { settings, FONT_FAMILIES }
})
