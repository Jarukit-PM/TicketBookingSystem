<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

import { resolveAdminTicket } from '@/api/admin'
import { ApiError } from '@/api/client'
import QrScanner from '@/components/admin/QrScanner.vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui'
import { parseTicketScanUrl } from '@/lib/parseTicketScanUrl'

const router = useRouter()
const resolving = ref(false)
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

async function handleScan(raw: string): Promise<void> {
  if (resolving.value) {
    return
  }

  const parsed = parseTicketScanUrl(raw)
  if (!parsed) {
    showToast('QR code is not a valid ticket link.')
    return
  }

  resolving.value = true
  try {
    const result = await resolveAdminTicket(parsed.bookingRef, parsed.token)
    await router.push({ name: 'admin-user-bookings', params: { userId: result.userId } })
  } catch (error) {
    const message =
      error instanceof ApiError ? error.message : 'Could not resolve ticket. Try again.'
    showToast(message)
  } finally {
    resolving.value = false
  }
}
</script>

<template>
  <div class="space-y-8">
    <div>
      <h1 class="text-2xl font-semibold text-copy-primary">Scan ticket</h1>
      <p class="mt-1 text-sm text-copy-secondary">
        Scan a customer&apos;s digital ticket QR to open their booking history for support lookup.
      </p>
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
        <CardTitle>QR scanner</CardTitle>
        <p class="text-sm text-copy-secondary">
          Point the camera at the ticket QR code, or upload a screenshot if the camera is blocked.
        </p>
      </CardHeader>
      <CardContent>
        <p v-if="resolving" class="mb-4 text-sm text-copy-muted">Resolving ticket…</p>
        <QrScanner @scan="handleScan" @error="showToast($event)" />
      </CardContent>
    </Card>

    <Card>
      <CardHeader>
        <CardTitle>What you need to do</CardTitle>
      </CardHeader>
      <CardContent class="space-y-3 text-sm text-copy-secondary">
        <p>
          This flow requires a real device with camera permissions. Automated tests only cover URL
          parsing — not live scanning.
        </p>
        <ol class="list-decimal space-y-2 pl-5">
          <li>On a second phone or laptop, sign in as a customer and complete a booking.</li>
          <li>Open the digital ticket (QR on white pad) from My Bookings or the confirmation email.</li>
          <li>On this admin device, allow camera access when prompted.</li>
          <li>Scan the customer ticket QR — you should land on that user&apos;s booking history.</li>
          <li>
            Try an edited QR or random image to confirm you stay on this page with an error message.
          </li>
        </ol>
      </CardContent>
    </Card>
  </div>
</template>
