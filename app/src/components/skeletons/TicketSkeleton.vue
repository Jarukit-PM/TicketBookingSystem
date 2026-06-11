<script setup lang="ts">
import { Skeleton } from '@/components/ui'
</script>

<template>
  <div class="space-y-4" aria-busy="true" :aria-label="$t('booking.ticket.loading')">
    <Skeleton class="h-6 w-40" />
    <Skeleton class="h-4 w-full" />
    <Skeleton class="mx-auto h-48 w-48 rounded-lg" />
    <Skeleton class="h-4 w-3/4" />
    <Skeleton class="h-4 w-1/2" />
  </div>
</template>
