<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink, useRoute } from 'vue-router'
import { fetchBookingDetail } from '@/api/bookings'
import { Badge, Button, Card, CardContent, CardHeader, CardTitle } from '@/components/ui'
import { useLocaleFormat } from '@/composables/useLocaleFormat'
import type { BookingListItem } from '@/types/bookings'

const { t } = useI18n()
const { formatDateTime, formatTHB } = useLocaleFormat()
const route = useRoute()
const booking = ref<BookingListItem | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)
const bookingId = computed(() => route.params.id as string)

onMounted(async () => {
  try {
    booking.value = await fetchBookingDetail(bookingId.value)
  } catch {
    error.value = t('booking.detail.notFound')
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="min-h-screen bg-base px-4 py-8 md:px-6">
    <div class="mx-auto max-w-lg space-y-6">
      <RouterLink to="/my-bookings" class="text-sm text-copy-secondary hover:text-copy-primary">
        {{ t('nav.backToMyBookings') }}
      </RouterLink>
      <p v-if="loading" class="text-sm text-copy-secondary">{{ t('booking.detail.loading') }}</p>
      <p v-else-if="error" class="text-sm text-state-error">{{ error }}</p>
      <Card v-else-if="booking">
        <CardHeader class="flex flex-row items-start justify-between gap-3">
          <CardTitle>{{ booking.movie.title }}</CardTitle>
          <Badge variant="confirmed">{{ t('booking.status.confirmed') }}</Badge>
        </CardHeader>
        <CardContent class="space-y-4">
          <p class="text-sm text-copy-secondary">{{ booking.cinema.name }} · {{ booking.screen.name }}</p>
          <dl class="space-y-2 text-sm">
            <div class="flex justify-between gap-4">
              <dt class="text-copy-secondary">{{ t('booking.detail.showtime') }}</dt>
              <dd>{{ formatDateTime(booking.startsAt) }}</dd>
            </div>
            <div class="flex justify-between gap-4">
              <dt class="text-copy-secondary">{{ t('booking.detail.seats') }}</dt>
              <dd>{{ booking.seats.join(', ') }}</dd>
            </div>
            <div class="flex justify-between gap-4">
              <dt class="text-copy-secondary">{{ t('booking.detail.total') }}</dt>
              <dd>{{ formatTHB(booking.total) }}</dd>
            </div>
            <div class="flex justify-between gap-4">
              <dt class="text-copy-secondary">{{ t('booking.detail.reference') }}</dt>
              <dd class="text-brand">{{ booking.bookingRef }}</dd>
            </div>
            <div v-if="booking.locale" class="flex justify-between gap-4">
              <dt class="text-copy-secondary">{{ t('booking.detail.locale') }}</dt>
              <dd class="uppercase">{{ booking.locale }}</dd>
            </div>
          </dl>
          <RouterLink :to="{ name: 'ticket', params: { bookingRef: booking.bookingRef } }">
            <Button variant="primary" class="w-full">{{ t('booking.detail.viewTicket') }}</Button>
          </RouterLink>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
