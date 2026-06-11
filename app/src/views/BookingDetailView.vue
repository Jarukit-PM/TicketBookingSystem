<script setup lang="ts">
import { ArrowLeft, QrCode } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink, useRoute } from 'vue-router'
import { fetchBookingDetail } from '@/api/bookings'
import BookingDetailSkeleton from '@/components/skeletons/BookingDetailSkeleton.vue'
import { Badge, Button, Card, CardContent, CardHeader, CardTitle, ErrorAlert } from '@/components/ui'
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
      <RouterLink
        to="/my-bookings"
        class="inline-flex items-center gap-1.5 text-sm text-copy-secondary transition-colors hover:text-copy-primary"
      >
        <ArrowLeft class="h-4 w-4" aria-hidden="true" />
        {{ t('nav.backToMyBookings') }}
      </RouterLink>

      <BookingDetailSkeleton v-if="loading" />
      <ErrorAlert v-else-if="error" :message="error" />
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
            <Button variant="primary" class="w-full gap-1.5">
              <QrCode class="h-4 w-4" aria-hidden="true" />
              {{ t('booking.detail.viewTicket') }}
            </Button>
          </RouterLink>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
