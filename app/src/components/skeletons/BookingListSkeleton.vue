<script setup lang="ts">
import { Card, Skeleton } from '@/components/ui'

withDefaults(
  defineProps<{
    count?: number
  }>(),
  {
    count: 3,
  },
)
</script>

<template>
  <ul class="space-y-4" aria-busy="true" :aria-label="$t('booking.myBookings.loading')">
    <li v-for="i in count" :key="i">
      <Card class="overflow-hidden">
        <div class="flex gap-4 p-4">
          <Skeleton class="h-28 w-20 shrink-0 rounded-lg" />
          <div class="flex min-w-0 flex-1 flex-col gap-3">
            <div class="flex items-start justify-between gap-2">
              <Skeleton class="h-5 w-2/3" />
              <Skeleton class="h-6 w-20 rounded-full" />
            </div>
            <Skeleton class="h-4 w-1/2" />
            <Skeleton class="h-4 w-3/4" />
          </div>
        </div>
      </Card>
    </li>
  </ul>
</template>
