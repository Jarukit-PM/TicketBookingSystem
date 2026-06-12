<script setup lang="ts">
import { Calendar, History, Ticket } from 'lucide-vue-next'
import AppHeader from '@/components/AppHeader.vue'
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink } from 'vue-router'
import { fetchMyBookings } from '@/api/bookings'
import BookingCard from '@/components/bookings/BookingCard.vue'
import BookingListSkeleton from '@/components/skeletons/BookingListSkeleton.vue'
import { Button, EmptyState, ErrorAlert } from '@/components/ui'
import type { BookingListItem } from '@/types/bookings'

type Tab = 'upcoming' | 'history'

const { t } = useI18n()
const activeTab = ref<Tab>('upcoming')
const bookings = ref<BookingListItem[]>([])
const loading = ref(true)
const error = ref<string | null>(null)

const emptyMessage = computed(() =>
  activeTab.value === 'upcoming'
    ? t('booking.myBookings.emptyUpcoming')
    : t('booking.myBookings.emptyHistory'),
)

async function loadBookings(): Promise<void> {
  loading.value = true
  error.value = null
  try {
    bookings.value = await fetchMyBookings(activeTab.value === 'upcoming')
  } catch {
    error.value = t('booking.myBookings.loadError')
  } finally {
    loading.value = false
  }
}

function switchTab(tab: Tab): void {
  if (activeTab.value === tab) return
  activeTab.value = tab
  void loadBookings()
}

onMounted(() => void loadBookings())
</script>

<template>
  <div class="min-h-screen bg-base">
    <AppHeader />

    <div class="mx-auto max-w-3xl space-y-6 px-4 py-8 md:px-6">
      <div class="flex items-center gap-3">
        <div class="flex h-10 w-10 items-center justify-center rounded-full bg-accent-dim">
          <Ticket class="h-5 w-5 text-brand" aria-hidden="true" />
        </div>
        <h1 class="text-2xl font-semibold text-copy-primary">{{ t('booking.myBookings.title') }}</h1>
      </div>

      <div class="flex gap-2 rounded-lg border border-surface-border bg-surface p-1">
        <button
          type="button"
          class="inline-flex flex-1 items-center justify-center gap-2 rounded-md px-4 py-2 text-sm font-medium transition-colors"
          :class="activeTab === 'upcoming' ? 'bg-gradient-brand text-white' : 'text-copy-secondary'"
          @click="switchTab('upcoming')"
        >
          <Calendar class="h-4 w-4" aria-hidden="true" />
          {{ t('booking.myBookings.upcoming') }}
        </button>
        <button
          type="button"
          class="inline-flex flex-1 items-center justify-center gap-2 rounded-md px-4 py-2 text-sm font-medium transition-colors"
          :class="activeTab === 'history' ? 'bg-gradient-brand text-white' : 'text-copy-secondary'"
          @click="switchTab('history')"
        >
          <History class="h-4 w-4" aria-hidden="true" />
          {{ t('booking.myBookings.history') }}
        </button>
      </div>

      <BookingListSkeleton v-if="loading" />
      <ErrorAlert v-else-if="error" :message="error" />
      <EmptyState
        v-else-if="bookings.length === 0"
        :icon="Ticket"
        :title="emptyMessage"
      >
        <template #action>
          <RouterLink to="/">
            <Button variant="primary">{{ t('booking.myBookings.findShowtime') }}</Button>
          </RouterLink>
        </template>
      </EmptyState>
      <ul v-else class="space-y-4">
        <li v-for="booking in bookings" :key="booking.id">
          <BookingCard :booking="booking" :upcoming="activeTab === 'upcoming'" />
        </li>
      </ul>
    </div>
  </div>
</template>
