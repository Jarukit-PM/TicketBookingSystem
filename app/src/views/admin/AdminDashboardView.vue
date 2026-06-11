<script setup lang="ts">
import { onMounted, ref } from 'vue'

import { ApiError, api } from '@/api/client'
import BookingsTable from '@/components/admin/BookingsTable.vue'
import StatsCard from '@/components/admin/StatsCard.vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui'
import type { AdminDashboard } from '@/types/admin'

const dashboard = ref<AdminDashboard | null>(null)
const loading = ref(true)
const errorMessage = ref('')

async function loadDashboard() {
  loading.value = true
  errorMessage.value = ''
  try {
    dashboard.value = await api.get<AdminDashboard>('/admin/dashboard')
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : 'Failed to load dashboard'
  } finally {
    loading.value = false
  }
}

onMounted(loadDashboard)
</script>

<template>
  <div class="space-y-8">
    <div>
      <h1 class="text-2xl font-semibold text-copy-primary">Dashboard</h1>
      <p class="mt-1 text-sm text-copy-secondary">Today's bookings, showtimes, and occupancy.</p>
    </div>

    <p v-if="errorMessage" class="text-sm text-state-error" role="alert">{{ errorMessage }}</p>

    <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <StatsCard
        label="Bookings today"
        :value="loading ? '—' : String(dashboard?.bookingsToday ?? 0)"
      />
      <StatsCard
        label="Showtimes today"
        :value="loading ? '—' : String(dashboard?.showtimesToday ?? 0)"
      />
      <StatsCard
        label="Avg occupancy"
        :value="loading ? '—' : `${(dashboard?.avgOccupancyPct ?? 0).toFixed(1)}%`"
        hint="Sold seats / sellable capacity"
      />
    </div>

    <Card>
      <CardHeader>
        <CardTitle>Recent bookings</CardTitle>
      </CardHeader>
      <CardContent>
        <BookingsTable :bookings="dashboard?.recentBookings ?? []" :loading="loading" />
      </CardContent>
    </Card>
  </div>
</template>
