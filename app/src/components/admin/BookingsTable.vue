<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink } from 'vue-router'

import AdminBookingDetailModal from '@/components/admin/AdminBookingDetailModal.vue'
import TableSkeleton from '@/components/skeletons/TableSkeleton.vue'
import { EmptyState } from '@/components/ui'
import { useLocaleFormat } from '@/composables/useLocaleFormat'
import type { BookingSummary } from '@/types/admin'
import { Inbox } from 'lucide-vue-next'

const props = withDefaults(
  defineProps<{
    bookings: BookingSummary[]
    loading?: boolean
    showCustomer?: boolean
    emptyMessage?: string
    focusBookingId?: string
    focusBookingRef?: string
    linkToBookings?: boolean
  }>(),
  {
    loading: false,
    showCustomer: false,
    emptyMessage: undefined,
    focusBookingId: undefined,
    focusBookingRef: undefined,
    linkToBookings: false,
  },
)

function bookingPageTo(booking: BookingSummary) {
  return {
    name: 'admin-bookings',
    query: { bookingRef: booking.bookingRef },
  }
}

const { t } = useI18n()
const { formatDateTime, formatTHB } = useLocaleFormat()

const detailOpen = ref(false)
const selectedBooking = ref<BookingSummary | null>(null)

function openBookingDetail(booking: BookingSummary) {
  selectedBooking.value = booking
  detailOpen.value = true
}

function tryFocusBooking() {
  if (props.loading || !props.bookings.length) return
  const match = props.bookings.find((booking) => {
    if (props.focusBookingId && booking.id === props.focusBookingId) return true
    if (props.focusBookingRef && booking.bookingRef === props.focusBookingRef) return true
    return false
  })
  if (match) openBookingDetail(match)
}

watch(
  () => [props.bookings, props.loading, props.focusBookingId, props.focusBookingRef] as const,
  tryFocusBooking,
  { immediate: true },
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

  <template v-else>
    <div class="space-y-3 md:hidden">
      <template v-for="booking in bookings" :key="booking.id">
        <RouterLink
          v-if="linkToBookings"
          :to="bookingPageTo(booking)"
          class="block w-full rounded-xl border border-surface-border bg-surface p-4 text-left transition-colors hover:bg-subtle"
        >
          <div class="flex items-start justify-between gap-3">
            <p class="font-medium text-brand">{{ booking.bookingRef }}</p>
            <span class="text-sm font-medium text-copy-primary">{{ formatTHB(booking.total) }}</span>
          </div>
          <p class="mt-2 text-sm text-copy-primary">{{ booking.movieTitle }}</p>
          <p class="mt-1 text-sm text-copy-secondary">{{ booking.seats.join(', ') }}</p>
          <p v-if="showCustomer" class="mt-1 truncate text-xs text-copy-muted">
            {{ booking.userEmail || booking.userId || t('common.dash') }}
          </p>
          <p class="mt-2 text-xs text-copy-muted">{{ formatDateTime(booking.confirmedAt) }}</p>
        </RouterLink>
        <button
          v-else
          type="button"
          class="w-full rounded-xl border border-surface-border bg-surface p-4 text-left transition-colors hover:bg-subtle"
          @click="openBookingDetail(booking)"
        >
          <div class="flex items-start justify-between gap-3">
            <p class="font-medium text-brand">{{ booking.bookingRef }}</p>
            <span class="text-sm font-medium text-copy-primary">{{ formatTHB(booking.total) }}</span>
          </div>
          <p class="mt-2 text-sm text-copy-primary">{{ booking.movieTitle }}</p>
          <p class="mt-1 text-sm text-copy-secondary">{{ booking.seats.join(', ') }}</p>
          <p v-if="showCustomer" class="mt-1 truncate text-xs text-copy-muted">
            {{ booking.userEmail || booking.userId || t('common.dash') }}
          </p>
          <p class="mt-2 text-xs text-copy-muted">{{ formatDateTime(booking.confirmedAt) }}</p>
        </button>
      </template>
    </div>

    <div class="hidden overflow-x-auto md:block">
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
          <td class="py-3 pr-4 font-medium">
            <RouterLink
              v-if="linkToBookings"
              :to="bookingPageTo(booking)"
              class="rounded-sm text-brand hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent-glow"
            >
              {{ booking.bookingRef }}
            </RouterLink>
            <button
              v-else
              type="button"
              class="rounded-sm text-brand hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent-glow"
              @click="openBookingDetail(booking)"
            >
              {{ booking.bookingRef }}
            </button>
          </td>
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
          <td class="py-3 pr-4 text-copy-secondary">
            {{ booking.locale === 'th' ? t('locale.th') : t('locale.en') }}
          </td>
          <td class="py-3 text-copy-secondary">{{ formatDateTime(booking.confirmedAt) }}</td>
        </tr>
      </tbody>
    </table>
    </div>
  </template>

  <AdminBookingDetailModal
    v-if="!linkToBookings"
    v-model:open="detailOpen"
    :booking-id="selectedBooking?.id ?? null"
    :summary="selectedBooking"
  />
</template>
