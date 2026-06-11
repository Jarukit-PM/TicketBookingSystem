<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { Button, Card, CardContent, CardHeader, CardTitle } from '@/components/ui'
import { useLocaleFormat } from '@/composables/useLocaleFormat'
import { useBookingSessionStore } from '@/stores/bookingSession'

const { t } = useI18n()
const { formatDateTime, formatTHB } = useLocaleFormat()
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
          <CardTitle>{{ t('booking.confirmation.title') }}</CardTitle>
          <p class="text-sm text-copy-secondary">
            {{ t('booking.confirmation.subtitle') }}
          </p>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="rounded-lg border border-surface-border bg-elevated p-4">
            <p class="text-xs uppercase tracking-widest text-copy-muted">
              {{ t('booking.confirmation.reference') }}
            </p>
            <p class="mt-1 text-2xl font-semibold text-brand">{{ booking.bookingRef }}</p>
          </div>

          <dl class="space-y-2 text-sm">
            <div class="flex justify-between gap-4">
              <dt class="text-copy-secondary">{{ t('booking.confirmation.seats') }}</dt>
              <dd class="font-medium text-copy-primary">{{ booking.seats.join(', ') }}</dd>
            </div>
            <div class="flex justify-between gap-4">
              <dt class="text-copy-secondary">{{ t('booking.confirmation.total') }}</dt>
              <dd class="font-medium text-copy-primary">{{ formatTHB(booking.total) }}</dd>
            </div>
            <div class="flex justify-between gap-4">
              <dt class="text-copy-secondary">{{ t('booking.confirmation.confirmed') }}</dt>
              <dd class="text-copy-primary">{{ formatDateTime(booking.confirmedAt) }}</dd>
            </div>
          </dl>

          <div class="flex flex-wrap gap-3 pt-2">
            <Button type="button" @click="viewTicket">{{ t('booking.confirmation.viewTicket') }}</Button>
            <Button type="button" variant="ghost" @click="goHome">{{ t('nav.backToHome') }}</Button>
            <Button type="button" variant="ghost" @click="bookMore">
              {{ t('booking.confirmation.bookMore') }}
            </Button>
          </div>
        </CardContent>
      </Card>

      <Card v-else>
        <CardContent class="space-y-4 py-8 text-center">
          <p class="text-copy-secondary">{{ t('booking.confirmation.notFound') }}</p>
          <Button type="button" @click="goHome">{{ t('nav.backToHome') }}</Button>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
