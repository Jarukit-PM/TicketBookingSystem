<script setup lang="ts">
import { Skeleton } from '@/components/ui'
</script>

<template>
  <div class="space-y-4" aria-busy="true" :aria-label="$t('booking.checkout.loading')">
    <Skeleton class="h-4 w-40" />
    <Skeleton class="h-4 w-full" />
    <Skeleton class="h-7 w-28" />
    <div class="flex gap-3 pt-2">
      <Skeleton class="h-10 w-32 rounded-full" />
      <Skeleton class="h-10 w-36 rounded-full" />
    </div>
  </div>
</template>
