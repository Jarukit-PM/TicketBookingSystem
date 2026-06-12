<script setup lang="ts">
import { CheckCircle2, Home, Plus, QrCode } from 'lucide-vue-next'
import AppHeader from '@/components/AppHeader.vue'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { Button, Card, CardContent, CardHeader, CardTitle, EmptyState } from '@/components/ui'
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
  <div class="min-h-screen bg-base">
    <AppHeader />

    <div class="mx-auto max-w-lg space-y-6 px-4 py-8 md:px-6">
      <Card v-if="booking">
        <CardHeader class="text-center">
          <div class="mx-auto mb-3 flex h-14 w-14 items-center justify-center rounded-full bg-state-success-dim">
            <CheckCircle2 class="h-7 w-7 text-state-success" aria-hidden="true" />
          </div>
          <CardTitle>{{ t('booking.confirmation.title') }}</CardTitle>
          <p class="text-sm text-copy-secondary">
            {{ t('booking.confirmation.subtitle') }}
          </p>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="rounded-lg border border-surface-border bg-elevated p-4 text-center">
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

          <div class="flex flex-col gap-3 pt-2 sm:flex-row sm:flex-wrap">
            <Button type="button" class="w-full gap-1.5 sm:w-auto" @click="viewTicket">
              <QrCode class="h-4 w-4" aria-hidden="true" />
              {{ t('booking.confirmation.viewTicket') }}
            </Button>
            <Button type="button" variant="ghost" class="w-full gap-1.5 sm:w-auto" @click="goHome">
              <Home class="h-4 w-4" aria-hidden="true" />
              {{ t('nav.backToHome') }}
            </Button>
            <Button type="button" variant="ghost" class="w-full gap-1.5 sm:w-auto" @click="bookMore">
              <Plus class="h-4 w-4" aria-hidden="true" />
              {{ t('booking.confirmation.bookMore') }}
            </Button>
          </div>
        </CardContent>
      </Card>

      <EmptyState v-else :icon="CheckCircle2" :title="t('booking.confirmation.notFound')">
        <template #action>
          <Button type="button" class="gap-1.5" @click="goHome">
            <Home class="h-4 w-4" aria-hidden="true" />
            {{ t('nav.backToHome') }}
          </Button>
        </template>
      </EmptyState>
    </div>
  </div>
</template>
