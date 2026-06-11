import { computed, ref, watch, type MaybeRefOrGetter, toValue } from 'vue'
import { toDateKey } from '@/composables/useLocaleFormat'
import type { Showtime } from '@/types/catalog'

export interface ShowtimeDateOption {
  key: string
  sampleIso: string
}

export function useShowtimeDates(
  showtimes: MaybeRefOrGetter<Showtime[]>,
  timeZone: MaybeRefOrGetter<string | undefined>,
) {
  const selectedDateKey = ref<string | null>(null)

  const dateOptions = computed<ShowtimeDateOption[]>(() => {
    const tz = toValue(timeZone)
    const byKey = new Map<string, string>()
    for (const showtime of toValue(showtimes)) {
      const key = toDateKey(showtime.startsAt, tz)
      if (!byKey.has(key)) {
        byKey.set(key, showtime.startsAt)
      }
    }
    return [...byKey.entries()]
      .sort(([a], [b]) => a.localeCompare(b))
      .map(([key, sampleIso]) => ({ key, sampleIso }))
  })

  watch(
    dateOptions,
    (options) => {
      if (options.length === 0) {
        selectedDateKey.value = null
        return
      }
      if (!selectedDateKey.value || !options.some((o) => o.key === selectedDateKey.value)) {
        selectedDateKey.value = options[0]!.key
      }
    },
    { immediate: true },
  )

  const selectedSampleIso = computed(
    () => dateOptions.value.find((o) => o.key === selectedDateKey.value)?.sampleIso,
  )

  const filteredShowtimes = computed(() => {
    if (!selectedDateKey.value) {
      return []
    }
    const tz = toValue(timeZone)
    return toValue(showtimes).filter(
      (showtime) => toDateKey(showtime.startsAt, tz) === selectedDateKey.value,
    )
  })

  return { selectedDateKey, dateOptions, selectedSampleIso, filteredShowtimes }
}
