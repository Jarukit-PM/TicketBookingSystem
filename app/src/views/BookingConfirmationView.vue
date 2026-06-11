<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Button, Card, CardContent, CardHeader, CardTitle } from '@/components/ui'
import { formatShowtime } from '@/lib/format'
import { useBookingSessionStore } from '@/stores/bookingSession'

const route = useRoute()
const router = useRouter()
const session = useBookingSessionStore()

const booking = computed(() => session.confirmedBooking)

const showtimeId = computed(() => route.params.showtimeId as string)
const bookingId = computed(() => route.params.bookingId as string)

function viewTicket(): void {
  router.push({ name: 'booking-ticket', params: { bookingId: bookingId.value } })
}

function goHome(): void {
  session.clear()
  router.push({ name: 'home' })
}

function bookMore(): void {
  session.resetCheckoutAttempt()
  session.confirmedBooking = null
  router.push({ name: 'book', params: { showtimeId: showtimeId.value } })
}
</script>

<template>
  <div class="min-h-screen bg-base px-4 py-8 md:px-6">
    <div class="mx-auto max-w-lg space-y-6">
      <Card v-if="booking">
        <CardHeader>
          <CardTitle>Booking confirmed</CardTitle>
          <p class="text-sm text-copy-secondary">
            Your tickets are reserved. A confirmation email will be sent shortly.
          </p>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="rounded-lg border border-surface-border bg-elevated p-4">
            <p class="text-xs uppercase tracking-widest text-copy-muted">Booking reference</p>
            <p class="mt-1 text-2xl font-semibold text-brand">{{ booking.bookingRef }}</p>
          </div>

          <dl class="space-y-2 text-sm">
            <div class="flex justify-between gap-4">
              <dt class="text-copy-secondary">Seats</dt>
              <dd class="font-medium text-copy-primary">{{ booking.seats.join(', ') }}</dd>
            </div>
            <div class="flex justify-between gap-4">
              <dt class="text-copy-secondary">Total</dt>
              <dd class="font-medium text-copy-primary">{{ booking.total.toLocaleString() }} THB</dd>
            </div>
            <div class="flex justify-between gap-4">
              <dt class="text-copy-secondary">Confirmed</dt>
              <dd class="text-copy-primary">{{ formatShowtime(booking.confirmedAt) }}</dd>
            </div>
          </dl>

          <div class="flex flex-wrap gap-3 pt-2">
            <Button type="button" @click="viewTicket">View ticket</Button>
            <Button type="button" variant="ghost" @click="goHome">Back to home</Button>
            <Button type="button" variant="ghost" @click="bookMore">Book more seats</Button>
          </div>
        </CardContent>
      </Card>

      <Card v-else>
        <CardContent class="space-y-4 py-8 text-center">
          <p class="text-copy-secondary">No booking details found.</p>
          <Button type="button" @click="goHome">Back to home</Button>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
