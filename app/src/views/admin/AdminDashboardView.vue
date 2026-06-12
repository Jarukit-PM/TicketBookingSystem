<script setup lang="ts">
import { BarChart3, Calendar, Ticket } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink } from 'vue-router'

import { translateApiError } from '@/api/errors'
import { ApiError, api } from '@/api/client'
import BookingsTable from '@/components/admin/BookingsTable.vue'
import StatsCard from '@/components/admin/StatsCard.vue'
import StatsGridSkeleton from '@/components/skeletons/StatsGridSkeleton.vue'
import { Card, CardContent, CardHeader, CardTitle, ErrorAlert } from '@/components/ui'
import type { AdminDashboard } from '@/types/admin'

const { t } = useI18n()

const dashboard = ref<AdminDashboard | null>(null)
const loading = ref(true)
const errorMessage = ref('')

async function loadDashboard() {
  loading.value = true
  errorMessage.value = ''
  try {
    dashboard.value = await api.get<AdminDashboard>('/admin/dashboard')
  } catch (error) {
    errorMessage.value =
      error instanceof ApiError
        ? translateApiError(error.code, error.message)
        : t('admin.dashboard.loadFailed')
  } finally {
    loading.value = false
  }
}

onMounted(loadDashboard)
</script>

<template>
  <div class="space-y-8">
    <div>
      <h1 class="text-2xl font-semibold text-copy-primary">{{ t('admin.dashboard.title') }}</h1>
      <p class="mt-1 text-sm text-copy-secondary">{{ t('admin.dashboard.subtitle') }}</p>
    </div>

    <ErrorAlert v-if="errorMessage" :message="errorMessage" />

    <StatsGridSkeleton v-if="loading" />
    <div v-else class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <StatsCard
        :icon="Ticket"
        :label="t('admin.dashboard.bookingsToday')"
        :value="String(dashboard?.bookingsToday ?? 0)"
      />
      <StatsCard
        :icon="Calendar"
        :label="t('admin.dashboard.showtimesToday')"
        :value="String(dashboard?.showtimesToday ?? 0)"
      />
      <StatsCard
        :icon="BarChart3"
        :label="t('admin.dashboard.avgOccupancy')"
        :value="`${(dashboard?.avgOccupancyPct ?? 0).toFixed(1)}%`"
        :hint="t('admin.dashboard.occupancyHint')"
      />
    </div>

    <Card>
      <CardHeader class="flex flex-row flex-wrap items-center justify-between gap-3">
        <CardTitle>{{ t('admin.dashboard.recentBookings') }}</CardTitle>
        <RouterLink
          :to="{ name: 'admin-bookings' }"
          class="text-sm font-medium text-brand hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent-glow"
        >
          {{ t('admin.dashboard.viewAllBookings') }}
        </RouterLink>
      </CardHeader>
      <CardContent>
        <BookingsTable
          :bookings="dashboard?.recentBookings ?? []"
          :loading="loading"
          link-to-bookings
        />
      </CardContent>
    </Card>
  </div>
</template>
