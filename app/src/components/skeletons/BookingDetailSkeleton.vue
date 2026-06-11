<script setup lang="ts">
import { Card, CardContent, CardHeader, Skeleton } from '@/components/ui'
</script>

<template>
  <Card aria-busy="true" :aria-label="$t('booking.detail.loading')">
    <CardHeader class="flex flex-row items-start justify-between gap-3">
      <Skeleton class="h-7 w-48" />
      <Skeleton class="h-6 w-24 rounded-full" />
    </CardHeader>
    <CardContent class="space-y-4">
      <Skeleton class="h-4 w-56" />
      <div class="space-y-3">
        <Skeleton v-for="i in 5" :key="i" class="h-4 w-full" />
      </div>
      <Skeleton class="h-10 w-full rounded-full" />
    </CardContent>
  </Card>
</template>
