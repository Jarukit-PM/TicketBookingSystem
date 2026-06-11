<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

import { translateApiError } from '@/api/errors'
import { ApiError, api } from '@/api/client'
import BookingsTable from '@/components/admin/BookingsTable.vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui'
import type { BookingSummary } from '@/types/admin'

const { t } = useI18n()
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
    bookings.value = data.bookings ?? []
  } catch (error) {
    bookings.value = []
    errorMessage.value =
      error instanceof ApiError
        ? translateApiError(error.code, error.message)
        : t('admin.userBookings.loadFailed')
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
      <h1 class="text-2xl font-semibold text-copy-primary">{{ t('admin.userBookings.title') }}</h1>
      <p class="mt-1 text-sm text-copy-secondary">
        {{ t('admin.userBookings.subtitle') }}<span v-if="customerEmail">{{ t('admin.userBookings.subtitleFor', { email: customerEmail }) }}</span>.
      </p>
      <p class="mt-1 font-mono text-xs text-copy-muted">
        {{ t('admin.userBookings.userIdLabel', { id: userId }) }}
      </p>
    </div>

    <p v-if="errorMessage" class="text-sm text-state-error" role="alert">{{ errorMessage }}</p>

    <Card>
      <CardHeader>
        <CardTitle>{{ t('admin.userBookings.historyTitle') }}</CardTitle>
      </CardHeader>
      <CardContent>
        <BookingsTable
          :bookings="bookings"
          :loading="loading"
          show-customer
          :empty-message="t('admin.userBookings.empty')"
        />
      </CardContent>
    </Card>
  </div>
</template>
