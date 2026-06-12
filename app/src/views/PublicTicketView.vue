<script setup lang="ts">
import { ArrowLeft, Home } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink, useRoute } from 'vue-router'
import { fetchMyBookings } from '@/api/bookings'
import { fetchBookingTicket, fetchPublicTicket } from '@/api/tickets'
import TicketCard from '@/components/TicketCard.vue'
import TicketSkeleton from '@/components/skeletons/TicketSkeleton.vue'
import { Button, Card, CardContent, CardHeader, CardTitle, ErrorAlert } from '@/components/ui'
import { useAuthStore } from '@/stores/auth'
import type { TicketDetail } from '@/types/ticket'

const { t } = useI18n()
const route = useRoute()
const auth = useAuthStore()
const ticket = ref<TicketDetail | null>(null)
const error = ref<string | null>(null)
const loading = ref(true)

const bookingRef = computed(() => route.params.bookingRef as string)
const ticketToken = computed(() => {
  const value = route.query.t
  return typeof value === 'string' ? value : ''
})

async function loadOwnedTicketFallback(): Promise<TicketDetail | null> {
  if (!auth.isAuthenticated) {
    return null
  }
  const [upcoming, history] = await Promise.all([
    fetchMyBookings(true),
    fetchMyBookings(false),
  ])
  const owned = [...upcoming, ...history].find((b) => b.bookingRef === bookingRef.value)
  if (!owned) {
    return null
  }
  return fetchBookingTicket(owned.id)
}

onMounted(async () => {
  try {
    if (!ticketToken.value) {
      ticket.value = await loadOwnedTicketFallback()
      if (!ticket.value) {
        error.value = t('booking.ticket.missingToken')
      }
      return
    }

    try {
      ticket.value = await fetchPublicTicket(bookingRef.value, ticketToken.value)
    } catch {
      ticket.value = await loadOwnedTicketFallback()
      if (!ticket.value) {
        error.value = t('booking.ticket.loadError')
      }
    }
  } catch {
    error.value = t('booking.ticket.loadError')
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="min-h-screen bg-base px-4 py-8 md:px-6">
    <div class="mx-auto max-w-lg space-y-6">
      <RouterLink
        :to="auth.isAuthenticated ? '/my-bookings' : '/'"
        class="inline-flex items-center gap-1.5 text-sm text-copy-secondary transition-colors hover:text-copy-primary"
      >
        <ArrowLeft class="h-4 w-4" aria-hidden="true" />
        {{ auth.isAuthenticated ? t('nav.backToMyBookings') : t('nav.backToHome') }}
      </RouterLink>

      <Card>
        <CardHeader>
          <CardTitle>{{ t('booking.ticket.title') }}</CardTitle>
          <p class="text-sm text-copy-secondary">{{ t('booking.ticket.subtitle') }}</p>
        </CardHeader>
        <CardContent>
          <TicketSkeleton v-if="loading" />
          <ErrorAlert v-else-if="error" :message="error" />
          <TicketCard v-else-if="ticket" :ticket="ticket" />
          <div class="pt-6">
            <RouterLink to="/">
              <Button type="button" variant="ghost" class="gap-1.5">
                <Home class="h-4 w-4" aria-hidden="true" />
                {{ t('nav.backToHome') }}
              </Button>
            </RouterLink>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
