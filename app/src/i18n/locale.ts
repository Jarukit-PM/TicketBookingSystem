export const SUPPORTED = ['en', 'th'] as const
export type AppLocale = (typeof SUPPORTED)[number]

export const LOCALE_STORAGE_KEY = 'tbs.locale'
export const LOCALE_COOKIE_NAME = 'Locale'

export function detectBrowserLocale(): AppLocale {
  const tags = navigator.languages?.length ? [...navigator.languages] : [navigator.language]
  for (const tag of tags) {
    if (tag.toLowerCase().startsWith('th')) return 'th'
  }
  return 'en'
}

export function resolveInitialLocale(): AppLocale {
  if (typeof localStorage === 'undefined') {
    return detectBrowserLocale()
  }
  const saved = localStorage.getItem(LOCALE_STORAGE_KEY)
  if (saved === 'en' || saved === 'th') return saved
  return detectBrowserLocale()
}

export function setLocaleCookie(locale: AppLocale): void {
  if (typeof document === 'undefined') return
  document.cookie = `${LOCALE_COOKIE_NAME}=${locale};path=/;max-age=31536000;samesite=lax`
}

export function applyDocumentLocale(locale: AppLocale): void {
  if (typeof document === 'undefined') return
  document.documentElement.lang = locale
}

export function persistLocale(locale: AppLocale): void {
  if (typeof localStorage !== 'undefined') {
    localStorage.setItem(LOCALE_STORAGE_KEY, locale)
  }
  setLocaleCookie(locale)
  applyDocumentLocale(locale)
}
