import { zhCN, dateZhCN, enUS, dateEnUS } from 'naive-ui'
import type { Locale } from './index'

interface NaiveLocalePair {
  locale: typeof zhCN | typeof enUS
  dateLocale: typeof dateZhCN | typeof dateEnUS
}

const localeMap: Record<Locale, NaiveLocalePair> = {
  'zh-CN': { locale: zhCN, dateLocale: dateZhCN },
  'en-US': { locale: enUS, dateLocale: dateEnUS },
}

export function getNaiveLocale(locale: Locale) {
  return localeMap[locale] ?? localeMap['zh-CN']
}
