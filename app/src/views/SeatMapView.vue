<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { addHolds, removeHolds } from '@/api/holds'
import { fetchSeatMap } from '@/api/seats'
import HoldCountdown from '@/components/seat-map/HoldCountdown.vue'
import SeatLegend from '@/components/seat-map/SeatLegend.vue'
import SeatMapGrid from '@/components/seat-map/SeatMapGrid.vue'
import { Button, Card, CardContent, CardHeader, CardTitle } from '@/components/ui'
import { useShowtimeSocket } from '@/composables/useShowtimeSocket'
import { formatShowtime } from '@/lib/format'
import { useAuthStore } from '@/stores/auth'
import { useBookingSessionStore } from '@/stores/bookingSession'
import { ApiError } from '@/api/client'
import type { SeatMapSnapshot } from '@/types/seats'

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
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load seat map'
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
      actionError.value = err.message
    } else {
      actionError.value = 'Failed to update seat selection'
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
    <header
      class="sticky top-0 z-10 flex h-16 items-center border-b border-surface-border bg-base/80 px-4 backdrop-blur-md md:px-6"
    >
      <Button variant="ghost" type="button" @click="goBack">Back</Button>
      <span class="ml-4 text-sm text-copy-secondary">Select seats</span>
    </header>

    <main class="mx-auto max-w-4xl px-4 py-8 md:px-6">
      <p v-if="loading && !snapshot" class="text-copy-secondary">Loading seat map…</p>
      <p v-else-if="error" class="text-state-error">{{ error }}</p>

      <template v-else-if="snapshot">
        <div class="mb-6 space-y-4">
          <HoldCountdown :expires-at="session.expiresAt" />
          <p v-if="actionError" class="text-sm text-state-error">{{ actionError }}</p>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>{{ snapshot.screenName }}</CardTitle>
            <p class="text-sm text-copy-secondary">
              {{ formatShowtime(snapshot.startsAt) }}
            </p>
          </CardHeader>
          <CardContent class="space-y-6">
            <div
              class="rounded-lg border border-surface-border bg-elevated px-4 py-2 text-center text-xs uppercase tracking-widest text-copy-muted"
            >
              Screen
            </div>

            <SeatMapGrid
              :seats="snapshot.seats"
              :self-held-ids="selfHeldIds"
              :pending-ids="pendingIds"
              :interactive="true"
              @select="toggleSeat"
            />
            <SeatLegend />
          </CardContent>
        </Card>

        <div
          v-if="selectedCount > 0"
          class="sticky bottom-4 mt-6 rounded-xl border border-surface-border bg-elevated p-4 shadow-elevation-2"
        >
          <div class="flex flex-wrap items-center justify-between gap-4">
            <div>
              <p class="text-sm text-copy-secondary">
                {{ selectedCount }} seat{{ selectedCount === 1 ? '' : 's' }} selected
              </p>
              <p class="text-lg font-semibold text-brand">
                {{ totalPrice.toLocaleString() }} THB
              </p>
            </div>
            <Button type="button" @click="goCheckout">Continue to checkout</Button>
          </div>
        </div>
      </template>
    </main>
  </div>
</template>
