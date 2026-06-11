<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { confirmBooking } from '@/api/bookings'
import { ApiError } from '@/api/client'
import { fetchSeatMap } from '@/api/seats'
import HoldCountdown from '@/components/seat-map/HoldCountdown.vue'
import { Button, Card, CardContent, CardHeader, CardTitle } from '@/components/ui'
import { formatShowtime } from '@/lib/format'
import { useBookingSessionStore } from '@/stores/bookingSession'
import type { SeatMapSnapshot } from '@/types/seats'

const route = useRoute()
const router = useRouter()
const session = useBookingSessionStore()

const loading = ref(true)
const confirming = ref(false)
const error = ref<string | null>(null)
const snapshot = ref<SeatMapSnapshot | null>(null)

const showtimeId = computed(() => route.params.showtimeId as string)

const totalPrice = computed(() => {
  const map = snapshot.value?.priceTiers
  const seats = snapshot.value?.seats
  if (!map || !seats) {
    return 0
  }
  return session.holds.reduce((sum, seatId) => {
    const seat = seats.find((s) => s.seatId === seatId)
    if (!seat) {
      return sum
    }
    return sum + (map[seat.type] ?? 0)
  }, 0)
})

async function loadSnapshot(): Promise<void> {
  loading.value = true
  error.value = null
  try {
    snapshot.value = await fetchSeatMap(showtimeId.value)
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load checkout details'
  } finally {
    loading.value = false
  }
}

function backToSeats(): void {
  router.push({ name: 'book', params: { showtimeId: showtimeId.value } })
}

async function handleConfirm(): Promise<void> {
  if (!session.holds.length || confirming.value) {
    return
  }

  confirming.value = true
  error.value = null
  const idempotencyKey = session.ensureIdempotencyKey()

  try {
    const booking = await confirmBooking(showtimeId.value, idempotencyKey)
    session.setConfirmedBooking(booking)
    await router.push({
      name: 'booking-confirmation',
      params: { showtimeId: showtimeId.value, bookingId: booking.id },
    })
  } catch (err) {
    if (err instanceof ApiError) {
      if (err.code === 'NO_ACTIVE_HOLDS' || err.code === 'SEAT_CONFLICT') {
        session.resetCheckoutAttempt()
      }
      error.value = err.message
    } else {
      error.value = 'Failed to confirm booking'
    }
  } finally {
    confirming.value = false
  }
}

onMounted(() => {
  session.ensureIdempotencyKey()
})

watch(
  showtimeId,
  async (id) => {
    session.setShowtime(id)
    await loadSnapshot()
    if (!session.holds.length) {
      await router.replace({ name: 'book', params: { showtimeId: id } })
    }
  },
  { immediate: true },
)
</script>

<template>
  <div class="min-h-screen bg-base px-4 py-8 md:px-6">
    <div class="mx-auto max-w-lg space-y-6">
      <HoldCountdown :expires-at="session.expiresAt" />

      <Card>
        <CardHeader>
          <CardTitle>Checkout</CardTitle>
          <p v-if="snapshot" class="text-sm text-copy-secondary">
            {{ snapshot.screenName }} · {{ formatShowtime(snapshot.startsAt) }}
          </p>
        </CardHeader>
        <CardContent class="space-y-4">
          <p v-if="loading" class="text-sm text-copy-secondary">Loading order summary…</p>

          <template v-else>
            <p v-if="session.holds.length" class="text-sm text-copy-primary">
              Selected seats:
              <span class="font-medium">{{ session.holds.join(', ') }}</span>
            </p>
            <p v-else class="text-sm text-copy-secondary">No seats selected.</p>

            <p class="text-lg font-semibold text-brand">
              {{ totalPrice.toLocaleString() }} THB
            </p>

            <p v-if="error" class="text-sm text-state-error">{{ error }}</p>

            <div class="flex flex-wrap gap-3">
              <Button type="button" variant="ghost" @click="backToSeats">Back to seat map</Button>
              <Button
                type="button"
                :disabled="!session.holds.length || confirming"
                @click="handleConfirm"
              >
                {{ confirming ? 'Confirming…' : 'Confirm booking' }}
              </Button>
            </div>
          </template>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
