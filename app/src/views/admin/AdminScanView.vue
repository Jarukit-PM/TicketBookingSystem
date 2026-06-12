<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { resolveAdminTicket } from '@/api/admin'
import { translateApiError } from '@/api/errors'
import { ApiError } from '@/api/client'
import AdminBookingDetailModal from '@/components/admin/AdminBookingDetailModal.vue'
import QrScanner from '@/components/admin/QrScanner.vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui'
import { parseTicketScanUrl } from '@/lib/parseTicketScanUrl'

const { t } = useI18n()
const resolving = ref(false)
const detailOpen = ref(false)
const resolvedBookingId = ref<string | null>(null)
const toastMessage = ref('')
const toastVariant = ref<'error' | 'info'>('error')

let toastTimer: ReturnType<typeof setTimeout> | null = null

function showToast(message: string, variant: 'error' | 'info' = 'error'): void {
  toastMessage.value = message
  toastVariant.value = variant
  if (toastTimer) {
    clearTimeout(toastTimer)
  }
  toastTimer = setTimeout(() => {
    toastMessage.value = ''
  }, 5000)
}

function onDetailOpenChange(open: boolean): void {
  detailOpen.value = open
  if (!open) {
    resolvedBookingId.value = null
  }
}

async function handleScan(raw: string): Promise<void> {
  if (resolving.value) {
    return
  }

  const parsed = parseTicketScanUrl(raw)
  if (!parsed) {
    showToast(t('admin.scan.invalidQr'))
    return
  }

  resolving.value = true
  try {
    const result = await resolveAdminTicket(parsed.bookingRef, parsed.token)
    resolvedBookingId.value = result.bookingId
    detailOpen.value = true
  } catch (error) {
    const message =
      error instanceof ApiError
        ? translateApiError(error.code, error.message)
        : t('admin.scan.resolveFailed')
    showToast(message)
  } finally {
    resolving.value = false
  }
}
</script>

<template>
  <div class="space-y-8">
    <div>
      <h1 class="text-2xl font-semibold text-copy-primary">{{ t('admin.scan.title') }}</h1>
      <p class="mt-1 text-sm text-copy-secondary">{{ t('admin.scan.subtitle') }}</p>
    </div>

    <div
      v-if="toastMessage"
      class="rounded-lg border px-4 py-3 text-sm"
      :class="
        toastVariant === 'error'
          ? 'border-state-error/40 bg-state-error-dim text-state-error'
          : 'border-surface-border bg-surface text-copy-secondary'
      "
      role="alert"
    >
      {{ toastMessage }}
    </div>

    <Card>
      <CardHeader>
        <CardTitle>{{ t('admin.scan.qrScanner') }}</CardTitle>
        <p class="text-sm text-copy-secondary">{{ t('admin.scan.scannerHint') }}</p>
      </CardHeader>
      <CardContent>
        <p v-if="resolving" class="mb-4 text-sm text-copy-muted">{{ t('admin.scan.resolving') }}</p>
        <QrScanner @scan="handleScan" @error="showToast($event)" />
      </CardContent>
    </Card>

    <Card>
      <CardHeader>
        <CardTitle>{{ t('admin.scan.instructionsTitle') }}</CardTitle>
      </CardHeader>
      <CardContent class="space-y-3 text-sm text-copy-secondary">
        <p>{{ t('admin.scan.instructionsIntro') }}</p>
        <ol class="list-decimal space-y-2 pl-5">
          <li>{{ t('admin.scan.step1') }}</li>
          <li>{{ t('admin.scan.step2') }}</li>
          <li>{{ t('admin.scan.step3') }}</li>
          <li>{{ t('admin.scan.step4') }}</li>
          <li>{{ t('admin.scan.step5') }}</li>
        </ol>
      </CardContent>
    </Card>

    <AdminBookingDetailModal
      :open="detailOpen"
      :booking-id="resolvedBookingId"
      @update:open="onDetailOpenChange"
    />
  </div>
</template>
