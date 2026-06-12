<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import { translateApiError } from '@/api/errors'
import { ApiError, api } from '@/api/client'
import TableSkeleton from '@/components/skeletons/TableSkeleton.vue'
import { Button, Card, CardContent, CardHeader, CardTitle, EmptyState, ErrorAlert, Input } from '@/components/ui'
import { FileText } from 'lucide-vue-next'
import { useLocaleFormat } from '@/composables/useLocaleFormat'
import type { AuditLogEntry, EmailLogEntry } from '@/types/admin'

type LogTab = 'audit' | 'email'

const { t, te } = useI18n()
const { formatDateTime } = useLocaleFormat()

const activeTab = ref<LogTab>('audit')
const bookingId = ref('')
const auditLogs = ref<AuditLogEntry[]>([])
const emailLogs = ref<EmailLogEntry[]>([])
const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const resendingBookingId = ref('')
const page = ref(1)
const limit = 50

const emptyMessage = computed(() =>
  activeTab.value === 'audit' ? t('admin.logs.emptyAudit') : t('admin.logs.emptyEmail'),
)

function formatAuditAction(action: string) {
  const key = `admin.logs.auditActions.${action}`
  return te(key) ? t(key) : action
}

function formatAuditDetails(meta?: Record<string, unknown>) {
  if (!meta || !Object.keys(meta).length) return t('common.dash')
  const parts: string[] = []
  if (typeof meta.bookingRef === 'string') {
    parts.push(t('admin.logs.auditDetails.ref', { value: meta.bookingRef }))
  }
  if (Array.isArray(meta.seats)) parts.push(t('admin.logs.auditDetails.seats', { value: meta.seats.join(', ') }))
  if (Array.isArray(meta.seatIds)) parts.push(t('admin.logs.auditDetails.seats', { value: meta.seatIds.join(', ') }))
  if (typeof meta.seatId === 'string') parts.push(t('admin.logs.auditDetails.seat', { value: meta.seatId }))
  if (typeof meta.code === 'string') parts.push(meta.code)
  if (typeof meta.reason === 'string') parts.push(meta.reason)
  if (typeof meta.message === 'string') parts.push(meta.message)
  return parts.length ? parts.join(' · ') : JSON.stringify(meta)
}

async function resendEmail(bookingId: string) {
  if (resendingBookingId.value) return
  resendingBookingId.value = bookingId
  errorMessage.value = ''
  successMessage.value = ''
  try {
    await api.post(`/admin/bookings/${bookingId}/resend-email`, {})
    successMessage.value = t('admin.logs.resendQueued')
    await loadLogs()
  } catch (error) {
    errorMessage.value =
      error instanceof ApiError
        ? translateApiError(error.code, error.message)
        : t('admin.logs.resendFailed')
  } finally {
    resendingBookingId.value = ''
  }
}

async function loadLogs() {
  loading.value = true
  errorMessage.value = ''
  successMessage.value = ''
  try {
    if (activeTab.value === 'audit') {
      const data = await api.get<{ logs: AuditLogEntry[] }>(
        `/admin/audit-logs?page=${page.value}&limit=${limit}`,
      )
      auditLogs.value = data.logs ?? []
    } else {
      const params = new URLSearchParams({
        page: String(page.value),
        limit: String(limit),
      })
      if (bookingId.value.trim()) params.set('bookingId', bookingId.value.trim())
      const data = await api.get<{ logs: EmailLogEntry[] }>(`/admin/email-logs?${params}`)
      emailLogs.value = data.logs ?? []
    }
  } catch (error) {
    auditLogs.value = []
    emailLogs.value = []
    errorMessage.value =
      error instanceof ApiError
        ? translateApiError(error.code, error.message)
        : t('admin.logs.loadFailed')
  } finally {
    loading.value = false
  }
}

function setTab(tab: LogTab) {
  activeTab.value = tab
  page.value = 1
}

watch(activeTab, loadLogs, { immediate: true })
</script>

<template>
  <div class="space-y-8">
    <div>
      <h1 class="text-2xl font-semibold text-copy-primary">{{ t('admin.logs.title') }}</h1>
      <p class="mt-1 text-sm text-copy-secondary">{{ t('admin.logs.subtitle') }}</p>
    </div>

    <div class="flex flex-wrap gap-2">
      <Button :variant="activeTab === 'audit' ? 'primary' : 'secondary'" @click="setTab('audit')">
        {{ t('admin.logs.auditTab') }}
      </Button>
      <Button :variant="activeTab === 'email' ? 'primary' : 'secondary'" @click="setTab('email')">
        {{ t('admin.logs.emailTab') }}
      </Button>
    </div>

    <Card v-if="activeTab === 'email'">
      <CardHeader>
        <CardTitle>{{ t('admin.logs.filterTitle') }}</CardTitle>
      </CardHeader>
      <CardContent class="flex flex-wrap items-end gap-3">
        <label class="block min-w-[16rem] flex-1 space-y-1.5">
          <span class="text-sm text-copy-secondary">{{ t('admin.logs.bookingId') }}</span>
          <Input v-model="bookingId" :placeholder="t('admin.logs.bookingIdPlaceholder')" autocomplete="off" />
        </label>
        <Button :disabled="loading" @click="loadLogs">{{ t('common.apply') }}</Button>
      </CardContent>
    </Card>

    <ErrorAlert v-if="errorMessage" :message="errorMessage" />
    <p v-if="successMessage" class="text-sm text-state-success">{{ successMessage }}</p>

    <Card>
      <CardHeader>
        <CardTitle>{{ activeTab === 'audit' ? t('admin.logs.auditEntries') : t('admin.logs.emailDeliveries') }}</CardTitle>
      </CardHeader>
      <CardContent>
        <TableSkeleton v-if="loading" :columns="activeTab === 'email' ? 6 : 5" :rows="6" />
        <EmptyState
          v-else-if="(activeTab === 'audit' && !auditLogs.length) || (activeTab === 'email' && !emailLogs.length)"
          :icon="FileText"
          :title="emptyMessage"
          class="py-10"
        />
        <div v-else class="overflow-x-auto">
          <table class="w-full text-left text-sm">
            <thead class="sticky top-0 bg-surface text-copy-muted">
              <tr v-if="activeTab === 'audit'">
                <th class="pb-3 pr-4 font-medium">{{ t('common.when') }}</th>
                <th class="pb-3 pr-4 font-medium">{{ t('admin.logs.action') }}</th>
                <th class="pb-3 pr-4 font-medium">{{ t('admin.logs.entity') }}</th>
                <th class="pb-3 pr-4 font-medium">{{ t('admin.logs.entityId') }}</th>
                <th class="pb-3 font-medium">{{ t('admin.logs.details') }}</th>
              </tr>
              <tr v-else>
                <th class="pb-3 pr-4 font-medium">{{ t('common.when') }}</th>
                <th class="pb-3 pr-4 font-medium">{{ t('admin.logs.to') }}</th>
                <th class="pb-3 pr-4 font-medium">{{ t('admin.logs.type') }}</th>
                <th class="pb-3 pr-4 font-medium">{{ t('common.status') }}</th>
                <th class="pb-3 pr-4 font-medium">{{ t('admin.logs.booking') }}</th>
                <th class="pb-3 font-medium">{{ t('admin.logs.actions') }}</th>
              </tr>
            </thead>
            <tbody>
              <template v-if="activeTab === 'audit'">
                <tr
                  v-for="log in auditLogs"
                  :key="log.id"
                  class="border-t border-surface-border"
                >
                  <td class="py-3 pr-4 text-copy-secondary">
                    {{ formatDateTime(log.createdAt) }}
                  </td>
                  <td class="py-3 pr-4 font-medium text-copy-primary">
                    {{ formatAuditAction(log.action) }}
                  </td>
                  <td class="py-3 pr-4 text-copy-secondary">{{ log.entity }}</td>
                  <td class="py-3 pr-4 font-mono text-xs text-copy-muted">{{ log.entityId }}</td>
                  <td class="py-3 text-xs text-copy-secondary">{{ formatAuditDetails(log.meta) }}</td>
                </tr>
              </template>
              <template v-else>
                <tr
                  v-for="log in emailLogs"
                  :key="log.id"
                  class="border-t border-surface-border"
                >
                  <td class="py-3 pr-4 text-copy-secondary">
                    {{ formatDateTime(log.createdAt) }}
                  </td>
                  <td class="py-3 pr-4 text-copy-primary">{{ log.to }}</td>
                  <td class="py-3 pr-4 text-copy-secondary">{{ log.type }}</td>
                  <td class="py-3 pr-4 text-copy-primary">{{ log.status }}</td>
                  <td class="py-3 pr-4 font-mono text-xs text-copy-muted">{{ log.bookingId }}</td>
                  <td class="py-3">
                    <Button
                      variant="secondary"
                      :disabled="resendingBookingId === log.bookingId"
                      @click="resendEmail(log.bookingId)"
                    >
                      {{
                        resendingBookingId === log.bookingId
                          ? t('admin.logs.resending')
                          : t('admin.logs.resend')
                      }}
                    </Button>
                  </td>
                </tr>
              </template>
            </tbody>
          </table>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
