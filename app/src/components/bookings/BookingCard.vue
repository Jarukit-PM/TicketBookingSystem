<script setup lang="ts">
import { ChevronRight, MapPin } from 'lucide-vue-next'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink } from 'vue-router'
import { Badge, Card, CardContent } from '@/components/ui'
import { useLocaleFormat } from '@/composables/useLocaleFormat'
import type { BookingListItem } from '@/types/bookings'

const props = defineProps<{
  booking: BookingListItem
  upcoming?: boolean
}>()

const { t } = useI18n()
const { formatDateTime } = useLocaleFormat()

const posterStyle = computed(() => ({
  backgroundImage: props.booking.movie.posterUrl ? `url(${props.booking.movie.posterUrl})` : undefined,
}))
</script>

<template>
  <RouterLink :to="{ name: 'booking-detail', params: { id: booking.id } }" class="block">
    <Card
      class="overflow-hidden transition-shadow hover:shadow-glow-brand/20"
      :class="upcoming ? 'border-l-2 border-l-brand' : ''"
    >
      <div class="flex gap-3 p-3 sm:gap-4 sm:p-4">
        <div
          class="h-24 w-16 shrink-0 rounded-lg bg-subtle bg-cover bg-center ring-1 ring-white/10 sm:h-28 sm:w-20"
          :style="posterStyle"
        />
        <CardContent class="flex min-w-0 flex-1 flex-col gap-1.5 p-0 sm:gap-2">
          <div class="flex items-start justify-between gap-2">
            <div class="min-w-0 flex-1">
              <h3 class="line-clamp-2 text-base font-semibold text-copy-primary sm:truncate sm:text-lg">
                {{ booking.movie.title }}
              </h3>
              <p class="mt-0.5 flex items-start gap-1 text-xs text-copy-secondary sm:text-sm">
                <MapPin class="mt-0.5 h-3.5 w-3.5 shrink-0" aria-hidden="true" />
                <span class="line-clamp-2">{{ booking.cinema.name }} · {{ booking.screen.name }}</span>
              </p>
            </div>
            <Badge variant="confirmed" class="shrink-0">
              {{ t('booking.status.confirmed') }}
            </Badge>
          </div>
          <p class="text-xs text-copy-secondary sm:text-sm">{{ formatDateTime(booking.startsAt) }}</p>
          <p class="line-clamp-2 text-xs text-copy-muted sm:text-sm">
            {{ t('common.ref') }} {{ booking.bookingRef }} · {{ booking.seats.join(', ') }}
          </p>
        </CardContent>
        <ChevronRight
          class="mt-1 hidden h-5 w-5 shrink-0 self-center text-copy-muted sm:block"
          aria-hidden="true"
        />
      </div>
    </Card>
  </RouterLink>
</template>
