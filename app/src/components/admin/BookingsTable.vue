<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { RouterLink } from 'vue-router'

import TableSkeleton from '@/components/skeletons/TableSkeleton.vue'
import { EmptyState } from '@/components/ui'
import { useLocaleFormat } from '@/composables/useLocaleFormat'
import type { BookingSummary } from '@/types/admin'
import { Inbox } from 'lucide-vue-next'

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
  <TableSkeleton v-if="loading" :columns="showCustomer ? 7 : 6" :rows="4" />

  <EmptyState
    v-else-if="!bookings.length"
    :icon="Inbox"
    :title="emptyMessage ?? t('admin.bookings.table.empty')"
    class="py-10"
  />

  <div v-else class="overflow-x-auto">
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
        <tr
          v-for="booking in bookings"
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
