<script setup lang="ts">
import { RouterLink } from 'vue-router'

import type { BookingSummary } from '@/types/admin'

withDefaults(
  defineProps<{
    bookings: BookingSummary[]
    loading?: boolean
    showCustomer?: boolean
    emptyMessage?: string
  }>(),
  {
    loading: false,
    showCustomer: false,
    emptyMessage: 'No confirmed bookings yet.',
  },
)

function formatTotal(cents: number) {
  return new Intl.NumberFormat(undefined, { style: 'currency', currency: 'USD' }).format(cents / 100)
}

function formatWhen(iso: string) {
  return new Date(iso).toLocaleString(undefined, {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}
</script>

<template>
  <div class="overflow-x-auto">
    <table class="w-full text-left text-sm">
      <thead class="sticky top-0 bg-surface text-copy-muted">
        <tr>
          <th class="pb-3 pr-4 font-medium">Ref</th>
          <th v-if="showCustomer" class="pb-3 pr-4 font-medium">Customer</th>
          <th class="pb-3 pr-4 font-medium">Movie</th>
          <th class="pb-3 pr-4 font-medium">Seats</th>
          <th class="pb-3 pr-4 font-medium">Total</th>
          <th class="pb-3 font-medium">Confirmed</th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="loading">
          <td :colspan="showCustomer ? 6 : 5" class="py-6 text-copy-muted">Loading bookings…</td>
        </tr>
        <tr v-else-if="!bookings.length">
          <td :colspan="showCustomer ? 6 : 5" class="py-6 text-copy-muted">{{ emptyMessage }}</td>
        </tr>
        <tr
          v-for="booking in bookings"
          v-else
          :key="booking.id"
          class="border-t border-surface-border"
        >
          <td class="py-3 pr-4 font-medium text-brand">{{ booking.bookingRef }}</td>
          <td v-if="showCustomer" class="py-3 pr-4 text-copy-secondary">
            <RouterLink
              v-if="booking.userId"
              :to="`/admin/users/${booking.userId}/bookings`"
              class="text-brand hover:underline"
            >
              {{ booking.userEmail || booking.userId }}
            </RouterLink>
            <span v-else>—</span>
          </td>
          <td class="py-3 pr-4 text-copy-primary">{{ booking.movieTitle }}</td>
          <td class="py-3 pr-4 text-copy-secondary">{{ booking.seats.join(', ') }}</td>
          <td class="py-3 pr-4 text-copy-primary">{{ formatTotal(booking.total) }}</td>
          <td class="py-3 text-copy-secondary">{{ formatWhen(booking.confirmedAt) }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
