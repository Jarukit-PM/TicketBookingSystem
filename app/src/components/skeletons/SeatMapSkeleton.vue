<script setup lang="ts">
import { Card, CardContent, CardHeader, Skeleton } from '@/components/ui'
</script>

<template>
  <Card aria-busy="true" :aria-label="$t('seatMap.loading')">
    <CardHeader class="space-y-2">
      <Skeleton class="h-6 w-40" />
      <Skeleton class="h-4 w-56" />
    </CardHeader>
    <CardContent class="space-y-6">
      <Skeleton class="h-8 w-full rounded-lg" />
      <div class="flex flex-col items-center gap-2">
        <div v-for="row in 6" :key="row" class="flex gap-2">
          <Skeleton v-for="col in 10" :key="col" class="h-10 w-10 rounded-md" />
        </div>
      </div>
      <div class="flex flex-wrap justify-center gap-3">
        <Skeleton v-for="i in 4" :key="i" class="h-8 w-24 rounded-full" />
      </div>
    </CardContent>
  </Card>
</template>
