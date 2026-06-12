<script setup lang="ts">
import { User } from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink } from 'vue-router'

import { translateApiError } from '@/api/errors'
import { ApiError, api } from '@/api/client'
import BookingDetailSkeleton from '@/components/skeletons/BookingDetailSkeleton.vue'
import { Badge, Button, ErrorAlert, Modal } from '@/components/ui'
import { useLocaleFormat } from '@/composables/useLocaleFormat'
import type { AdminBookingDetail, BookingSummary } from '@/types/admin'

const props = defineProps<{
  open: boolean
  bookingId: string | null
  summary?: BookingSummary | null
}>()

const emit = defineEmits<{
  'update:open': [value: boolean]
}>()

const { t } = useI18n()
const { formatDateTime, formatTHB } = useLocaleFormat()

const detail = ref<AdminBookingDetail | null>(null)
const loading = ref(false)
const errorMessage = ref('')

const display = computed(() => detail.value ?? props.summary ?? null)
const hasVenueContext = computed(
  () => Boolean(detail.value?.cinemaName && detail.value?.screenName),
)
const hasShowtime = computed(() => Boolean(detail.value?.startsAt))

function applySummary(summary: BookingSummary) {
  detail.value = {
    ...summary,
    startsAt: '',
    cinemaId: '',
    cinemaName: '',
    screenId: '',
    screenName: '',
    status: 'CONFIRMED',
  }
  errorMessage.value = ''
}

async function enrichBooking(id: string) {
  try {
    detail.value = await api.get<AdminBookingDetail>(`/admin/bookings/${id}`)
    errorMessage.value = ''
  } catch {
    // Keep summary-only view when the detail endpoint is unavailable.
  }
}

async function loadBooking(id: string, summary?: BookingSummary | null) {
  loading.value = true
  errorMessage.value = ''
  detail.value = null

  if (summary) {
    applySummary(summary)
    loading.value = false
    void enrichBooking(id)
    return
  }

  try {
    detail.value = await api.get<AdminBookingDetail>(`/admin/bookings/${id}`)
  } catch (error) {
    errorMessage.value =
      error instanceof ApiError
        ? translateApiError(error.code, error.message)
        : t('admin.bookings.detail.loadFailed')
  } finally {
    loading.value = false
  }
}

watch(
  () => [props.open, props.bookingId, props.summary] as const,
  ([isOpen, id, summary]) => {
    if (isOpen && id) {
      void loadBooking(id, summary)
      return
    }
    detail.value = null
    errorMessage.value = ''
  },
)
</script>

<template>
  <Modal
    :open="open"
    size="xl"
    :title="display?.bookingRef ?? t('admin.bookings.detail.title')"
    @update:open="emit('update:open', $event)"
  >
    <BookingDetailSkeleton v-if="loading" />
    <ErrorAlert v-else-if="errorMessage" :message="errorMessage" />
    <div v-else-if="display" class="space-y-5">
      <div class="flex flex-wrap items-start justify-between gap-3">
        <div class="space-y-1">
          <p class="font-semibold text-copy-primary">{{ display.movieTitle }}</p>
          <p v-if="hasVenueContext" class="text-sm text-copy-secondary">
            {{ detail?.cinemaName }} · {{ detail?.screenName }}
          </p>
        </div>
        <Badge variant="confirmed">{{ t('booking.status.confirmed') }}</Badge>
      </div>

      <dl class="grid gap-3 text-sm sm:grid-cols-2">
        <div class="space-y-1">
          <dt class="text-copy-secondary">{{ t('admin.bookings.detail.reference') }}</dt>
          <dd class="font-medium text-brand">{{ display.bookingRef }}</dd>
        </div>
        <div class="space-y-1">
          <dt class="text-copy-secondary">{{ t('admin.bookings.detail.customer') }}</dt>
          <dd>
            <RouterLink
              v-if="display.userId"
              :to="`/admin/users/${display.userId}/bookings`"
              class="inline-flex items-center gap-1 text-brand hover:underline"
              @click="emit('update:open', false)"
            >
              <User class="h-3.5 w-3.5" aria-hidden="true" />
              {{ display.userEmail || display.userId }}
            </RouterLink>
            <span v-else>{{ t('common.dash') }}</span>
          </dd>
        </div>
        <div v-if="hasShowtime" class="space-y-1">
          <dt class="text-copy-secondary">{{ t('booking.detail.showtime') }}</dt>
          <dd>{{ formatDateTime(detail!.startsAt) }}</dd>
        </div>
        <div class="space-y-1">
          <dt class="text-copy-secondary">{{ t('booking.detail.seats') }}</dt>
          <dd>{{ display.seats.join(', ') }}</dd>
        </div>
        <div class="space-y-1">
          <dt class="text-copy-secondary">{{ t('booking.detail.total') }}</dt>
          <dd>{{ formatTHB(display.total) }}</dd>
        </div>
        <div class="space-y-1">
          <dt class="text-copy-secondary">{{ t('booking.detail.emailLocale') }}</dt>
          <dd>{{ display.locale === 'th' ? t('locale.th') : t('locale.en') }}</dd>
        </div>
        <div class="space-y-1">
          <dt class="text-copy-secondary">{{ t('admin.bookings.table.confirmed') }}</dt>
          <dd>{{ formatDateTime(display.confirmedAt) }}</dd>
        </div>
        <div class="space-y-1 sm:col-span-2">
          <dt class="text-copy-secondary">{{ t('admin.bookings.showtimeId') }}</dt>
          <dd class="font-mono text-xs text-copy-muted">{{ display.showtimeId }}</dd>
        </div>
      </dl>

      <RouterLink
        v-if="display.userId"
        :to="`/admin/users/${display.userId}/bookings`"
        @click="emit('update:open', false)"
      >
        <Button variant="secondary" class="gap-1.5">
          <User class="h-4 w-4" aria-hidden="true" />
          {{ t('admin.bookings.detail.viewCustomerHistory') }}
        </Button>
      </RouterLink>
    </div>
  </Modal>
</template>
