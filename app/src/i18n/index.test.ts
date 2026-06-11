import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import {
  LOCALE_STORAGE_KEY,
  applyDocumentLocale,
  detectBrowserLocale,
  resolveInitialLocale,
} from '@/i18n/locale'

describe('detectBrowserLocale', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('returns th when navigator prefers Thai', () => {
    vi.stubGlobal('navigator', { language: 'en-US', languages: ['th-TH', 'en-US'] })
    expect(detectBrowserLocale()).toBe('th')
  })

  it('returns en for non-Thai browsers', () => {
    vi.stubGlobal('navigator', { language: 'en-US', languages: ['en-US'] })
    expect(detectBrowserLocale()).toBe('en')
  })
})

describe('resolveInitialLocale', () => {
  const storage = new Map<string, string>()

  afterEach(() => {
    storage.clear()
    vi.unstubAllGlobals()
  })

  beforeEach(() => {
    vi.stubGlobal('localStorage', {
      getItem: (key: string) => storage.get(key) ?? null,
      setItem: (key: string, value: string) => {
        storage.set(key, value)
      },
      removeItem: (key: string) => {
        storage.delete(key)
      },
    })
  })

  it('prefers saved locale over browser', () => {
    storage.set(LOCALE_STORAGE_KEY, 'th')
    vi.stubGlobal('navigator', { language: 'en-US', languages: ['en-US'] })
    expect(resolveInitialLocale()).toBe('th')
  })

  it('falls back to browser detect when nothing saved', () => {
    vi.stubGlobal('navigator', { language: 'th-TH', languages: ['th-TH'] })
    expect(resolveInitialLocale()).toBe('th')
  })
})

describe('applyDocumentLocale', () => {
  it('sets html lang attribute', () => {
    vi.stubGlobal('document', { documentElement: { lang: 'en' } })
    applyDocumentLocale('th')
    expect(document.documentElement.lang).toBe('th')
  })
})
