<script setup lang="ts">
import { ref } from 'vue'

import { ApiError, api } from '@/api/client'
import BookingsTable from '@/components/admin/BookingsTable.vue'
import { Button, Card, CardContent, CardHeader, CardTitle, Input } from '@/components/ui'
import type { BookingSummary } from '@/types/admin'

const bookingRef = ref('')
const email = ref('')
const userId = ref('')
const showtimeId = ref('')
const bookings = ref<BookingSummary[]>([])
const loading = ref(false)
const errorMessage = ref('')
const searched = ref(false)

async function search() {
  loading.value = true
  errorMessage.value = ''
  searched.value = true
  try {
    const params = new URLSearchParams()
    if (bookingRef.value.trim()) params.set('bookingRef', bookingRef.value.trim())
    if (email.value.trim()) params.set('email', email.value.trim())
    if (userId.value.trim()) params.set('userId', userId.value.trim())
    if (showtimeId.value.trim()) params.set('showtimeId', showtimeId.value.trim())
    const query = params.toString()
    const path = query ? `/admin/bookings?${query}` : '/admin/bookings'
    const data = await api.get<{ bookings: BookingSummary[] }>(path)
    bookings.value = data.bookings
  } catch (error) {
    bookings.value = []
    errorMessage.value = error instanceof ApiError ? error.message : 'Search failed'
  } finally {
    loading.value = false
  }
}

function clearFilters() {
  bookingRef.value = ''
  email.value = ''
  userId.value = ''
  showtimeId.value = ''
  bookings.value = []
  searched.value = false
  errorMessage.value = ''
}
</script>

<template>
  <div class="space-y-8">
    <div>
      <h1 class="text-2xl font-semibold text-copy-primary">Bookings</h1>
      <p class="mt-1 text-sm text-copy-secondary">Search confirmed bookings for support lookup.</p>
    </div>

    <Card>
      <CardHeader>
        <CardTitle>Search</CardTitle>
      </CardHeader>
      <CardContent class="space-y-4">
        <div class="grid gap-4 sm:grid-cols-2">
          <label class="block space-y-1.5">
            <span class="text-sm text-copy-secondary">Booking ref</span>
            <Input v-model="bookingRef" placeholder="TBS-…" autocomplete="off" />
          </label>
          <label class="block space-y-1.5">
            <span class="text-sm text-copy-secondary">Customer email</span>
            <Input v-model="email" type="email" placeholder="customer@example.com" autocomplete="off" />
          </label>
          <label class="block space-y-1.5">
            <span class="text-sm text-copy-secondary">User ID</span>
            <Input v-model="userId" placeholder="MongoDB ObjectId" autocomplete="off" />
          </label>
          <label class="block space-y-1.5">
            <span class="text-sm text-copy-secondary">Showtime ID</span>
            <Input v-model="showtimeId" placeholder="MongoDB ObjectId" autocomplete="off" />
          </label>
        </div>
        <div class="flex flex-wrap gap-3">
          <Button :disabled="loading" @click="search">Search</Button>
          <Button variant="secondary" :disabled="loading" @click="clearFilters">Clear</Button>
        </div>
        <p v-if="errorMessage" class="text-sm text-state-error" role="alert">{{ errorMessage }}</p>
      </CardContent>
    </Card>

    <Card>
      <CardHeader>
        <CardTitle>Results</CardTitle>
      </CardHeader>
      <CardContent>
        <BookingsTable
          :bookings="bookings"
          :loading="loading"
          show-customer
          :empty-message="searched ? 'No bookings match your search.' : 'Enter a filter and search.'"
        />
      </CardContent>
    </Card>
  </div>
</template>
