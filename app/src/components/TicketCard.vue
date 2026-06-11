<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useLocaleFormat } from '@/composables/useLocaleFormat'
import type { TicketDetail } from '@/types/ticket'

const props = defineProps<{ ticket: TicketDetail }>()
const { t } = useI18n()
const { formatDateTime, formatTHB } = useLocaleFormat()
const qrSrc = computed(() => `data:image/png;base64,${props.ticket.qrPngBase64}`)
</script>

<template>
  <div class="space-y-6">
    <div class="rounded-lg border border-surface-border bg-elevated p-4">
      <p class="text-xs uppercase tracking-widest text-copy-muted">{{ t('booking.ticket.reference') }}</p>
      <p class="mt-1 text-2xl font-semibold text-brand">{{ ticket.bookingRef }}</p>
    </div>
    <div class="mx-auto w-fit rounded-lg bg-white p-4 shadow-1">
      <img
        :src="qrSrc"
        :alt="t('booking.ticket.qrAlt')"
        class="h-56 w-56"
        width="224"
        height="224"
      />
    </div>
    <dl class="space-y-2 text-sm">
      <div class="flex justify-between gap-4">
        <dt class="text-copy-secondary">{{ t('booking.ticket.movie') }}</dt>
        <dd class="text-right font-medium text-copy-primary">{{ ticket.movieTitle }}</dd>
      </div>
      <div class="flex justify-between gap-4">
        <dt class="text-copy-secondary">{{ t('booking.ticket.venue') }}</dt>
        <dd class="text-right text-copy-primary">{{ ticket.cinemaName }} · {{ ticket.screenName }}</dd>
      </div>
      <div class="flex justify-between gap-4">
        <dt class="text-copy-secondary">{{ t('booking.ticket.showtime') }}</dt>
        <dd class="text-copy-primary">{{ formatDateTime(ticket.startsAt) }}</dd>
      </div>
      <div class="flex justify-between gap-4">
        <dt class="text-copy-secondary">{{ t('booking.ticket.seats') }}</dt>
        <dd class="font-medium text-copy-primary">{{ ticket.seats.join(', ') }}</dd>
      </div>
      <div class="flex justify-between gap-4">
        <dt class="text-copy-secondary">{{ t('booking.ticket.total') }}</dt>
        <dd class="font-medium text-copy-primary">{{ formatTHB(ticket.total) }}</dd>
      </div>
    </dl>
  </div>
</template>
