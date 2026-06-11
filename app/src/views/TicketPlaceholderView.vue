<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { Button, Card, CardContent, CardHeader, CardTitle } from '@/components/ui'

const bookingRef = computed(() => useRoute().params.bookingRef as string)
</script>

<template>
  <div class="min-h-screen bg-base px-4 py-8 md:px-6">
    <div class="mx-auto max-w-lg space-y-6">
      <RouterLink to="/my-bookings" class="text-sm text-copy-secondary hover:text-copy-primary">
        ← Back to My Bookings
      </RouterLink>
      <Card>
        <CardHeader>
          <CardTitle>Digital ticket</CardTitle>
          <p class="text-sm text-copy-secondary">QR code and email delivery are coming in a future update.</p>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="rounded-lg border border-surface-border bg-elevated p-4 text-center">
            <p class="text-xs uppercase tracking-widest text-copy-muted">Booking reference</p>
            <p class="mt-1 text-2xl font-semibold text-brand">{{ bookingRef }}</p>
          </div>
          <RouterLink to="/my-bookings">
            <Button variant="secondary" class="w-full">Back to My Bookings</Button>
          </RouterLink>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
