<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { fetchBookingDetail } from '@/api/bookings'
import { Badge, Button, Card, CardContent, CardHeader, CardTitle } from '@/components/ui'
import { formatShowtime } from '@/lib/format'
import type { BookingListItem } from '@/types/bookings'

const route = useRoute()
const booking = ref<BookingListItem | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)
const bookingId = computed(() => route.params.id as string)

onMounted(async () => {
  try {
    booking.value = await fetchBookingDetail(bookingId.value)
  } catch {
    error.value = 'Booking not found or you do not have access.'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="min-h-screen bg-base px-4 py-8 md:px-6">
    <div class="mx-auto max-w-lg space-y-6">
      <RouterLink to="/my-bookings" class="text-sm text-copy-secondary hover:text-copy-primary">
        ← Back to My Bookings
      </RouterLink>
      <p v-if="loading" class="text-sm text-copy-secondary">Loading booking…</p>
      <p v-else-if="error" class="text-sm text-state-error">{{ error }}</p>
      <Card v-else-if="booking">
        <CardHeader class="flex flex-row items-start justify-between gap-3">
          <CardTitle>{{ booking.movie.title }}</CardTitle>
          <Badge variant="confirmed">{{ booking.status }}</Badge>
        </CardHeader>
        <CardContent class="space-y-4">
          <p class="text-sm text-copy-secondary">{{ booking.cinema.name }} · {{ booking.screen.name }}</p>
          <dl class="space-y-2 text-sm">
            <div class="flex justify-between gap-4">
              <dt class="text-copy-secondary">Showtime</dt>
              <dd>{{ formatShowtime(booking.startsAt) }}</dd>
            </div>
            <div class="flex justify-between gap-4">
              <dt class="text-copy-secondary">Seats</dt>
              <dd>{{ booking.seats.join(', ') }}</dd>
            </div>
            <div class="flex justify-between gap-4">
              <dt class="text-copy-secondary">Total</dt>
              <dd>{{ booking.total.toLocaleString() }} THB</dd>
            </div>
            <div class="flex justify-between gap-4">
              <dt class="text-copy-secondary">Reference</dt>
              <dd class="text-brand">{{ booking.bookingRef }}</dd>
            </div>
          </dl>
          <RouterLink :to="{ name: 'ticket', params: { bookingRef: booking.bookingRef } }">
            <Button variant="primary" class="w-full">View ticket</Button>
          </RouterLink>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
