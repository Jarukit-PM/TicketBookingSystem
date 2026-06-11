<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Button, Card, CardContent, CardHeader, CardTitle } from '@/components/ui'
import HoldCountdown from '@/components/seat-map/HoldCountdown.vue'
import { useBookingSessionStore } from '@/stores/bookingSession'

const route = useRoute()
const router = useRouter()
const session = useBookingSessionStore()

const showtimeId = computed(() => route.params.showtimeId as string)

function backToSeats(): void {
  router.push({ name: 'book', params: { showtimeId: showtimeId.value } })
}
</script>

<template>
  <div class="min-h-screen bg-base px-4 py-8 md:px-6">
    <div class="mx-auto max-w-lg space-y-6">
      <HoldCountdown :expires-at="session.expiresAt" />

      <Card>
        <CardHeader>
          <CardTitle>Checkout</CardTitle>
          <p class="text-sm text-copy-secondary">
            Booking confirmation will be implemented in the next release.
          </p>
        </CardHeader>
        <CardContent class="space-y-4">
          <p v-if="session.holds.length" class="text-sm text-copy-primary">
            Selected seats:
            <span class="font-medium">{{ session.holds.join(', ') }}</span>
          </p>
          <p v-else class="text-sm text-copy-secondary">No seats selected.</p>

          <div class="flex flex-wrap gap-3">
            <Button type="button" variant="ghost" @click="backToSeats">Back to seat map</Button>
            <Button type="button" disabled>Confirm booking (coming soon)</Button>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
