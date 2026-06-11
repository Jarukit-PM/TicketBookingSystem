<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { fetchMyBookings } from '@/api/bookings'
import BookingCard from '@/components/bookings/BookingCard.vue'
import { Button } from '@/components/ui'
import type { BookingListItem } from '@/types/bookings'

type Tab = 'upcoming' | 'history'
const activeTab = ref<Tab>('upcoming')
const bookings = ref<BookingListItem[]>([])
const loading = ref(true)
const error = ref<string | null>(null)
const emptyMessage = computed(() =>
  activeTab.value === 'upcoming' ? 'No upcoming bookings yet.' : 'No past bookings yet.',
)

async function loadBookings(): Promise<void> {
  loading.value = true
  error.value = null
  try {
    bookings.value = await fetchMyBookings(activeTab.value === 'upcoming')
  } catch {
    error.value = 'Could not load your bookings. Please try again.'
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
  <div class="min-h-screen bg-base px-4 py-8 md:px-6">
    <div class="mx-auto max-w-3xl space-y-6">
      <h1 class="text-2xl font-semibold text-copy-primary">My Bookings</h1>
      <div class="flex gap-2 rounded-lg border border-surface-border bg-surface p-1">
        <button
          type="button"
          class="flex-1 rounded-md px-4 py-2 text-sm font-medium"
          :class="activeTab === 'upcoming' ? 'bg-gradient-brand text-white' : 'text-copy-secondary'"
          @click="switchTab('upcoming')"
        >
          Upcoming
        </button>
        <button
          type="button"
          class="flex-1 rounded-md px-4 py-2 text-sm font-medium"
          :class="activeTab === 'history' ? 'bg-gradient-brand text-white' : 'text-copy-secondary'"
          @click="switchTab('history')"
        >
          History
        </button>
      </div>
      <p v-if="loading" class="text-sm text-copy-secondary">Loading bookings…</p>
      <p v-else-if="error" class="text-sm text-state-error">{{ error }}</p>
      <div
        v-else-if="bookings.length === 0"
        class="rounded-lg border border-surface-border bg-surface p-8 text-center"
      >
        <p class="text-copy-secondary">{{ emptyMessage }}</p>
        <RouterLink to="/" class="mt-4 inline-block">
          <Button variant="primary">Find a showtime</Button>
        </RouterLink>
      </div>
      <ul v-else class="space-y-4">
        <li v-for="booking in bookings" :key="booking.id">
          <BookingCard :booking="booking" />
        </li>
      </ul>
    </div>
  </div>
</template>
