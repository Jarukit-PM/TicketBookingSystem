<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink } from 'vue-router'
import { Badge, Card, CardContent } from '@/components/ui'
import { useLocaleFormat } from '@/composables/useLocaleFormat'
import type { BookingListItem } from '@/types/bookings'

const props = defineProps<{ booking: BookingListItem }>()
const { t } = useI18n()
const { formatDateTime } = useLocaleFormat()

const posterStyle = computed(() => ({
  backgroundImage: props.booking.movie.posterUrl ? `url(${props.booking.movie.posterUrl})` : undefined,
}))
</script>

<template>
  <RouterLink :to="{ name: 'booking-detail', params: { id: booking.id } }" class="block">
    <Card class="overflow-hidden transition-shadow hover:shadow-glow-brand/20">
      <div class="flex gap-4 p-4">
        <div
          class="h-28 w-20 shrink-0 rounded-lg bg-subtle bg-cover bg-center ring-1 ring-white/10"
          :style="posterStyle"
        />
        <CardContent class="flex min-w-0 flex-1 flex-col gap-2 p-0">
          <div class="flex items-start justify-between gap-2">
            <div class="min-w-0">
              <h3 class="truncate text-lg font-semibold text-copy-primary">{{ booking.movie.title }}</h3>
              <p class="text-sm text-copy-secondary">{{ booking.cinema.name }} · {{ booking.screen.name }}</p>
            </div>
            <Badge variant="confirmed">{{ t('booking.status.confirmed') }}</Badge>
          </div>
          <p class="text-sm text-copy-secondary">{{ formatDateTime(booking.startsAt) }}</p>
          <p class="text-sm text-copy-muted">
            {{ t('common.ref') }} {{ booking.bookingRef }} · {{ booking.seats.join(', ') }}
          </p>
        </CardContent>
      </div>
    </Card>
  </RouterLink>
</template>
