<script setup lang="ts">
import { ShoppingCart } from 'lucide-vue-next'
import AppHeader from '@/components/AppHeader.vue'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { addHolds, removeHolds } from '@/api/holds'
import { fetchSeatMap } from '@/api/seats'
import HoldCountdown from '@/components/seat-map/HoldCountdown.vue'
import SeatLegend from '@/components/seat-map/SeatLegend.vue'
import SeatMapGrid from '@/components/seat-map/SeatMapGrid.vue'
import SeatMapSkeleton from '@/components/skeletons/SeatMapSkeleton.vue'
import { Button, Card, CardContent, CardHeader, CardTitle, ErrorAlert } from '@/components/ui'
import { translateApiError } from '@/api/errors'
import { useLocaleFormat } from '@/composables/useLocaleFormat'
import { useShowtimeSocket } from '@/composables/useShowtimeSocket'
import { useAuthStore } from '@/stores/auth'
import { useBookingSessionStore } from '@/stores/bookingSession'
import { ApiError } from '@/api/client'
import type { SeatMapSnapshot } from '@/types/seats'

const { t } = useI18n()
const { formatDateTime, formatTHB } = useLocaleFormat()
const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const session = useBookingSessionStore()

const loading = ref(true)
const error = ref<string | null>(null)
const actionError = ref<string | null>(null)
const restSnapshot = ref<SeatMapSnapshot | null>(null)

const showtimeId = computed(() => route.params.showtimeId as string)

const { snapshot: wsSnapshot } = useShowtimeSocket(showtimeId)

const snapshot = computed(() => wsSnapshot.value ?? restSnapshot.value)

const selfHeldIds = computed(() => session.holdSet)
const pendingIds = computed(() => session.pendingSeatIds)

const selectedCount = computed(() => session.holds.length)

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

async function loadSeatMap(): Promise<void> {
  loading.value = true
  error.value = null
  try {
    restSnapshot.value = await fetchSeatMap(showtimeId.value)
  } catch {
    error.value = t('seatMap.loadError')
    restSnapshot.value = null
  } finally {
    loading.value = false
  }
}

async function restoreSessionHolds(): Promise<void> {
  if (!auth.isAuthenticated) {
    return
  }
  try {
    const result = await addHolds(showtimeId.value, [])
    session.applyHoldResult(result)
  } catch {
    // ignore — user may have no active holds
  }
}

async function toggleSeat(seatId: string): Promise<void> {
  actionError.value = null

  if (!auth.isAuthenticated) {
    await router.push({
      name: 'login',
      query: { redirect: route.fullPath },
    })
    return
  }

  const seat = snapshot.value?.seats.find((s) => s.seatId === seatId)
  if (!seat) {
    return
  }

  const isSelfHeld = session.isSelfHeld(seatId)
  if (!isSelfHeld && seat.status !== 'AVAILABLE') {
    return
  }

  session.markPending([seatId])

  try {
    if (isSelfHeld) {
      const result = await removeHolds(showtimeId.value, [seatId])
      session.applyHoldResult(result)
    } else {
      const result = await addHolds(showtimeId.value, [seatId])
      session.applyHoldResult(result)
    }
  } catch (err) {
    if (err instanceof ApiError) {
      actionError.value = translateApiError(err.code, err.message)
    } else {
      actionError.value = t('seatMap.actionError')
    }
    await loadSeatMap()
  } finally {
    session.clearPending([seatId])
  }
}

function goCheckout(): void {
  router.push({ name: 'checkout', params: { showtimeId: showtimeId.value } })
}

function goBack(): void {
  if (window.history.length > 1) {
    router.back()
    return
  }
  router.push({ name: 'home' })
}

watch(
  showtimeId,
  async (id) => {
    session.setShowtime(id)
    await loadSeatMap()
    await restoreSessionHolds()
  },
  { immediate: true },
)
</script>

<template>
  <div class="min-h-screen bg-base">
    <AppHeader show-back :subtitle="t('seatMap.selectSeats')" @back="goBack" />

    <main class="mx-auto max-w-4xl px-4 py-6 pb-28 sm:py-8 md:px-6">
      <SeatMapSkeleton v-if="loading && !snapshot" />
      <ErrorAlert v-else-if="error" :message="error" />

      <template v-else-if="snapshot">
        <div class="mb-6 space-y-4">
          <HoldCountdown :expires-at="session.expiresAt" />
          <ErrorAlert v-if="actionError" :message="actionError" />
        </div>

        <Card>
          <CardHeader>
            <CardTitle>{{ snapshot.screenName }}</CardTitle>
            <p class="text-sm text-copy-secondary">
              {{ formatDateTime(snapshot.startsAt) }}
            </p>
          </CardHeader>
          <CardContent class="space-y-6">
            <div
              class="rounded-lg border border-surface-border bg-elevated px-4 py-2 text-center text-xs uppercase tracking-widest text-copy-muted"
            >
              {{ t('seatMap.screen') }}
            </div>

            <SeatMapGrid
              :seats="snapshot.seats"
              :self-held-ids="selfHeldIds"
              :pending-ids="pendingIds"
              :interactive="true"
              @select="toggleSeat"
            />
            <SeatLegend :price-tiers="snapshot.priceTiers" />
          </CardContent>
        </Card>

        <div
          v-if="selectedCount > 0"
          class="fixed inset-x-0 bottom-0 z-20 border-t border-surface-border bg-elevated/95 p-4 shadow-elevation-2 backdrop-blur-md sm:sticky sm:inset-x-auto sm:bottom-4 sm:mt-6 sm:rounded-xl sm:border"
          style="padding-bottom: max(1rem, env(safe-area-inset-bottom))"
        >
          <div class="mx-auto flex max-w-4xl flex-col gap-3 sm:flex-row sm:items-center sm:justify-between sm:gap-4">
            <div>
              <p class="text-sm text-copy-secondary">
                {{ t('seatMap.seatsSelected', selectedCount) }}
              </p>
              <p class="text-lg font-semibold text-brand">
                {{ formatTHB(totalPrice) }}
              </p>
            </div>
            <Button type="button" class="w-full gap-1.5 sm:w-auto" @click="goCheckout">
              <ShoppingCart class="h-4 w-4" aria-hidden="true" />
              {{ t('seatMap.continueCheckout') }}
            </Button>
          </div>
        </div>
      </template>
    </main>
  </div>
</template>
