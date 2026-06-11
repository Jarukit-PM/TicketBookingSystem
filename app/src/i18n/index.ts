import { createI18n } from 'vue-i18n'
import en from '@/locales/en.json'
import th from '@/locales/th.json'
import {
  applyDocumentLocale,
  resolveInitialLocale,
} from '@/i18n/locale'

export type { AppLocale } from '@/i18n/locale'
export {
  LOCALE_STORAGE_KEY,
  applyDocumentLocale,
  detectBrowserLocale,
  persistLocale,
  resolveInitialLocale,
  setLocaleCookie,
} from '@/i18n/locale'

const initialLocale = resolveInitialLocale()
applyDocumentLocale(initialLocale)

export const i18n = createI18n({
  legacy: false,
  locale: initialLocale,
  fallbackLocale: 'en',
  messages: { en, th },
})
