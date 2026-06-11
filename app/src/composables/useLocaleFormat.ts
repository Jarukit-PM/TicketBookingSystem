import { useI18n } from 'vue-i18n'

export function useLocaleFormat() {
  const { locale } = useI18n()

  function intlLocale(): string {
    return locale.value === 'th' ? 'th-TH' : 'en-US'
  }

  function formatDateTime(iso: string, timeZone?: string): string {
    return new Intl.DateTimeFormat(intlLocale(), {
      weekday: 'short',
      month: 'short',
      day: 'numeric',
      hour: 'numeric',
      minute: '2-digit',
      timeZone,
    }).format(new Date(iso))
  }

  function formatTime(iso: string, timeZone?: string): string {
    return new Intl.DateTimeFormat(intlLocale(), {
      hour: 'numeric',
      minute: '2-digit',
      timeZone,
    }).format(new Date(iso))
  }

  function formatDayShort(iso: string, timeZone?: string): string {
    return new Intl.DateTimeFormat(intlLocale(), {
      weekday: 'short',
      timeZone,
    }).format(new Date(iso))
  }

  function formatDayNumber(iso: string, timeZone?: string): string {
    return new Intl.DateTimeFormat(intlLocale(), {
      day: 'numeric',
      timeZone,
    }).format(new Date(iso))
  }

  function formatMonthYear(iso: string, timeZone?: string): string {
    return new Intl.DateTimeFormat(intlLocale(), {
      month: 'short',
      year: 'numeric',
      timeZone,
    }).format(new Date(iso))
  }

  /** Format an amount stored in satang (minor units) as THB. */
  function formatTHB(satang: number): string {
    const baht = satang / 100
    return new Intl.NumberFormat(intlLocale(), {
      style: 'currency',
      currency: 'THB',
      minimumFractionDigits: 0,
      maximumFractionDigits: 2,
    }).format(baht)
  }

  return {
    formatDateTime,
    formatTime,
    formatDayShort,
    formatDayNumber,
    formatMonthYear,
    formatTHB,
  }
}

/** Calendar date in a timezone as YYYY-MM-DD (en-CA). */
export function toDateKey(iso: string, timeZone?: string): string {
  return new Intl.DateTimeFormat('en-CA', {
    timeZone,
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  }).format(new Date(iso))
}
