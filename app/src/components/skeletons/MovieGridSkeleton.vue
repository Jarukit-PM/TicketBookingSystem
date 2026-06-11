<script setup lang="ts">
import { Card, CardContent, Skeleton } from '@/components/ui'

withDefaults(
  defineProps<{
    count?: number
  }>(),
  {
    count: 8,
  },
)
</script>

<template>
  <div
    class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4"
    aria-busy="true"
    :aria-label="$t('catalog.loadingMovies')"
  >
    <Card v-for="i in count" :key="i" class="overflow-hidden">
      <Skeleton class="aspect-[2/3] w-full rounded-none rounded-t-xl" />
      <CardContent class="space-y-3 p-4">
        <Skeleton class="h-5 w-3/4" />
        <Skeleton class="h-4 w-1/2" />
        <Skeleton class="h-10 w-full rounded-full" />
      </CardContent>
    </Card>
  </div>
</template>
