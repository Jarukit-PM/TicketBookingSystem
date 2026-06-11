<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'

import { ApiError, api } from '@/api/client'
import BookingsTable from '@/components/admin/BookingsTable.vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui'
import type { BookingSummary } from '@/types/admin'

const route = useRoute()
const bookings = ref<BookingSummary[]>([])
const loading = ref(true)
const errorMessage = ref('')

const userId = computed(() => String(route.params.userId ?? ''))
const customerEmail = computed(() => bookings.value[0]?.userEmail ?? '')

async function loadBookings() {
  if (!userId.value) return
  loading.value = true
  errorMessage.value = ''
  try {
    const data = await api.get<{ bookings: BookingSummary[] }>(
      `/admin/users/${userId.value}/bookings`,
    )
    bookings.value = data.bookings
  } catch (error) {
    bookings.value = []
    errorMessage.value = error instanceof ApiError ? error.message : 'Failed to load bookings'
  } finally {
    loading.value = false
  }
}

onMounted(loadBookings)
watch(userId, loadBookings)
</script>

<template>
  <div class="space-y-8">
    <div>
      <h1 class="text-2xl font-semibold text-copy-primary">User bookings</h1>
      <p class="mt-1 text-sm text-copy-secondary">
        Full confirmed booking history
        <span v-if="customerEmail"> for {{ customerEmail }}</span>.
      </p>
      <p class="mt-1 font-mono text-xs text-copy-muted">User ID: {{ userId }}</p>
    </div>

    <p v-if="errorMessage" class="text-sm text-state-error" role="alert">{{ errorMessage }}</p>

    <Card>
      <CardHeader>
        <CardTitle>Booking history</CardTitle>
      </CardHeader>
      <CardContent>
        <BookingsTable
          :bookings="bookings"
          :loading="loading"
          show-customer
          empty-message="No confirmed bookings for this user."
        />
      </CardContent>
    </Card>
  </div>
</template>
