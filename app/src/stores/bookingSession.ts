import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

export const useBookingSessionStore = defineStore('bookingSession', () => {
  const showtimeId = ref<string | null>(null)
  const holds = ref<string[]>([])
  const expiresAt = ref<string | null>(null)
  const pendingSeatIds = ref<Set<string>>(new Set())
  const idempotencyKey = ref<string | null>(null)
  const confirmedBooking = ref<import('@/types/bookings').ConfirmedBooking | null>(null)

  const hasHolds = computed(() => holds.value.length > 0)
  const holdSet = computed(() => new Set(holds.value))

  function setShowtime(id: string) {
    if (showtimeId.value !== id) {
      showtimeId.value = id
      holds.value = []
      expiresAt.value = null
      pendingSeatIds.value = new Set()
    }
  }

  function applyHoldResult(result: { holds: string[]; expiresAt?: string }) {
    holds.value = [...result.holds]
    expiresAt.value = result.expiresAt ?? null
  }

  function markPending(seatIds: string[]) {
    const next = new Set(pendingSeatIds.value)
    for (const id of seatIds) {
      next.add(id)
    }
    pendingSeatIds.value = next
  }

  function clearPending(seatIds: string[]) {
    const next = new Set(pendingSeatIds.value)
    for (const id of seatIds) {
      next.delete(id)
    }
    pendingSeatIds.value = next
  }

  function isPending(seatId: string): boolean {
    return pendingSeatIds.value.has(seatId)
  }

  function isSelfHeld(seatId: string): boolean {
    return holdSet.value.has(seatId)
  }

  function ensureIdempotencyKey(): string {
    if (!idempotencyKey.value) {
      idempotencyKey.value = crypto.randomUUID()
    }
    return idempotencyKey.value
  }

  function resetCheckoutAttempt(): void {
    idempotencyKey.value = crypto.randomUUID()
  }

  function setConfirmedBooking(booking: import('@/types/bookings').ConfirmedBooking): void {
    confirmedBooking.value = booking
    holds.value = []
    expiresAt.value = null
    idempotencyKey.value = null
  }

  function clear() {
    showtimeId.value = null
    holds.value = []
    expiresAt.value = null
    pendingSeatIds.value = new Set()
    idempotencyKey.value = null
    confirmedBooking.value = null
  }

  return {
    showtimeId,
    holds,
    expiresAt,
    pendingSeatIds,
    idempotencyKey,
    confirmedBooking,
    hasHolds,
    holdSet,
    setShowtime,
    applyHoldResult,
    markPending,
    clearPending,
    isPending,
    isSelfHeld,
    ensureIdempotencyKey,
    resetCheckoutAttempt,
    setConfirmedBooking,
    clear,
  }
})
