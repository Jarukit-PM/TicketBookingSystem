<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { confirmBooking } from '@/api/bookings'
import { translateApiError } from '@/api/errors'
import { ApiError } from '@/api/client'
import { fetchSeatMap } from '@/api/seats'
import HoldCountdown from '@/components/seat-map/HoldCountdown.vue'
import { Button, Card, CardContent, CardHeader, CardTitle } from '@/components/ui'
import { useLocaleFormat } from '@/composables/useLocaleFormat'
import { useBookingSessionStore } from '@/stores/bookingSession'
import type { SeatMapSnapshot } from '@/types/seats'

const { t, locale } = useI18n()
const { formatDateTime, formatTHB } = useLocaleFormat()
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
  } catch {
    error.value = t('booking.checkout.loadError')
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
    const booking = await confirmBooking(showtimeId.value, idempotencyKey, locale.value)
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
      error.value = translateApiError(err.code, err.message)
    } else {
      error.value = t('booking.checkout.confirmError')
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
          <CardTitle>{{ t('booking.checkout.title') }}</CardTitle>
          <p v-if="snapshot" class="text-sm text-copy-secondary">
            {{ snapshot.screenName }} · {{ formatDateTime(snapshot.startsAt) }}
          </p>
        </CardHeader>
        <CardContent class="space-y-4">
          <p v-if="loading" class="text-sm text-copy-secondary">{{ t('booking.checkout.loading') }}</p>

          <template v-else>
            <p v-if="session.holds.length" class="text-sm text-copy-primary">
              {{ t('booking.checkout.selectedSeats') }}
              <span class="font-medium">{{ session.holds.join(', ') }}</span>
            </p>
            <p v-else class="text-sm text-copy-secondary">{{ t('booking.checkout.noSeats') }}</p>

            <p class="text-lg font-semibold text-brand">
              {{ formatTHB(totalPrice) }}
            </p>

            <p v-if="error" class="text-sm text-state-error">{{ error }}</p>

            <div class="flex flex-wrap gap-3">
              <Button type="button" variant="ghost" @click="backToSeats">
                {{ t('booking.checkout.backToSeatMap') }}
              </Button>
              <Button
                type="button"
                :disabled="!session.holds.length || confirming"
                @click="handleConfirm"
              >
                {{ confirming ? t('booking.checkout.confirming') : t('booking.checkout.confirm') }}
              </Button>
            </div>
          </template>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
