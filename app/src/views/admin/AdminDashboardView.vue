<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { translateApiError } from '@/api/errors'
import { ApiError, api } from '@/api/client'
import BookingsTable from '@/components/admin/BookingsTable.vue'
import StatsCard from '@/components/admin/StatsCard.vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui'
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

    <p v-if="errorMessage" class="text-sm text-state-error" role="alert">{{ errorMessage }}</p>

    <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <StatsCard
        :label="t('admin.dashboard.bookingsToday')"
        :value="loading ? t('common.dash') : String(dashboard?.bookingsToday ?? 0)"
      />
      <StatsCard
        :label="t('admin.dashboard.showtimesToday')"
        :value="loading ? t('common.dash') : String(dashboard?.showtimesToday ?? 0)"
      />
      <StatsCard
        :label="t('admin.dashboard.avgOccupancy')"
        :value="loading ? t('common.dash') : `${(dashboard?.avgOccupancyPct ?? 0).toFixed(1)}%`"
        :hint="t('admin.dashboard.occupancyHint')"
      />
    </div>

    <Card>
      <CardHeader>
        <CardTitle>{{ t('admin.dashboard.recentBookings') }}</CardTitle>
      </CardHeader>
      <CardContent>
        <BookingsTable :bookings="dashboard?.recentBookings ?? []" :loading="loading" />
      </CardContent>
    </Card>
  </div>
</template>
