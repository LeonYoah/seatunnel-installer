/**
 * i18n 初始化
 */
import { createI18n } from 'vue-i18n'
import zhCN from '@/locales/zh-CN'
import enUS from '@/locales/en'

export type AppLocale = 'zh-CN' | 'en'

const saved = (localStorage.getItem('locale') as AppLocale) || 'zh-CN'

const i18n = createI18n({
  legacy: false,
  locale: saved,
  fallbackLocale: 'zh-CN',
  messages: {
    'zh-CN': zhCN,
    en: enUS
  }
})

export default i18n

