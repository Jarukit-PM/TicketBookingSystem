import { i18n } from '@/i18n'

export function translateApiError(code: string, fallback?: string): string {
  const key = `errors.${code.toLowerCase()}`
  if (i18n.global.te(key)) return i18n.global.t(key)
  return fallback ?? i18n.global.t('errors.generic')
}
