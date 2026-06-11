import { computed, onScopeDispose, ref, toRef, watch } from 'vue'
import type { MaybeRefOrGetter } from 'vue'

function parseExpiresAt(value: string | null | undefined): number | null {
  if (!value) {
    return null
  }
  const ms = Date.parse(value)
  return Number.isNaN(ms) ? null : ms
}

export function useHoldCountdown(expiresAt: MaybeRefOrGetter<string | null | undefined>) {
  const expiresAtRef = toRef(expiresAt)
  const now = ref(Date.now())
  let timer: ReturnType<typeof setInterval> | null = null

  function startTicker() {
    stopTicker()
    timer = setInterval(() => {
      now.value = Date.now()
    }, 1000)
  }

  function stopTicker() {
    if (timer !== null) {
      clearInterval(timer)
      timer = null
    }
  }

  watch(
    expiresAtRef,
    (value) => {
      if (parseExpiresAt(value) !== null) {
        now.value = Date.now()
        startTicker()
      } else {
        stopTicker()
      }
    },
    { immediate: true },
  )

  onScopeDispose(stopTicker)

  const expiresAtMs = computed(() => parseExpiresAt(expiresAtRef.value))

  const remainingMs = computed(() => {
    const end = expiresAtMs.value
    if (end === null) {
      return 0
    }
    return Math.max(0, end - now.value)
  })

  const remainingSeconds = computed(() => Math.ceil(remainingMs.value / 1000))

  const isActive = computed(() => expiresAtMs.value !== null && remainingMs.value > 0)

  const isUrgent = computed(() => isActive.value && remainingSeconds.value <= 60)

  const formatted = computed(() => {
    const total = remainingSeconds.value
    const minutes = Math.floor(total / 60)
    const seconds = total % 60
    return `${minutes}:${seconds.toString().padStart(2, '0')}`
  })

  return {
    remainingMs,
    remainingSeconds,
    isActive,
    isUrgent,
    formatted,
  }
}
