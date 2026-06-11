<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { RouterLink } from 'vue-router'

import { useLocaleFormat } from '@/composables/useLocaleFormat'
import type { BookingSummary } from '@/types/admin'

const { t } = useI18n()
const { formatDateTime, formatTHB } = useLocaleFormat()

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
    emptyMessage: undefined,
  },
)
</script>

<template>
  <div class="overflow-x-auto">
    <table class="w-full text-left text-sm">
      <thead class="sticky top-0 bg-surface text-copy-muted">
        <tr>
          <th class="pb-3 pr-4 font-medium">{{ t('admin.bookings.table.ref') }}</th>
          <th v-if="showCustomer" class="pb-3 pr-4 font-medium">
            {{ t('admin.bookings.table.customer') }}
          </th>
          <th class="pb-3 pr-4 font-medium">{{ t('admin.bookings.table.movie') }}</th>
          <th class="pb-3 pr-4 font-medium">{{ t('admin.bookings.table.seats') }}</th>
          <th class="pb-3 pr-4 font-medium">{{ t('admin.bookings.table.total') }}</th>
          <th class="pb-3 pr-4 font-medium">{{ t('admin.bookings.table.locale') }}</th>
          <th class="pb-3 font-medium">{{ t('admin.bookings.table.confirmed') }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="loading">
          <td :colspan="showCustomer ? 7 : 6" class="py-6 text-copy-muted">
            {{ t('admin.bookings.table.loading') }}
          </td>
        </tr>
        <tr v-else-if="!bookings.length">
          <td :colspan="showCustomer ? 7 : 6" class="py-6 text-copy-muted">
            {{ emptyMessage ?? t('admin.bookings.table.empty') }}
          </td>
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
            <span v-else>{{ t('common.dash') }}</span>
          </td>
          <td class="py-3 pr-4 text-copy-primary">{{ booking.movieTitle }}</td>
          <td class="py-3 pr-4 text-copy-secondary">{{ booking.seats.join(', ') }}</td>
          <td class="py-3 pr-4 text-copy-primary">{{ formatTHB(booking.total) }}</td>
          <td class="py-3 pr-4 text-copy-secondary uppercase">{{ booking.locale || 'en' }}</td>
          <td class="py-3 text-copy-secondary">{{ formatDateTime(booking.confirmedAt) }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
