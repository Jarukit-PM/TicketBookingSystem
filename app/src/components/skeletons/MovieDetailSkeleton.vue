<script setup lang="ts">
import { Skeleton } from '@/components/ui'
</script>

<template>
  <div
    class="grid gap-8 lg:grid-cols-[280px_1fr]"
    aria-busy="true"
    :aria-label="$t('movie.loading')"
  >
    <Skeleton class="aspect-[2/3] w-full max-w-xs rounded-xl" />
    <div class="flex flex-col gap-6">
      <div class="space-y-3">
        <Skeleton class="h-9 w-2/3 max-w-md" />
        <Skeleton class="h-4 w-40" />
        <Skeleton class="h-4 w-full max-w-2xl" />
        <Skeleton class="h-4 w-full max-w-xl" />
        <Skeleton class="h-4 w-3/4 max-w-lg" />
      </div>
      <div class="space-y-4">
        <Skeleton class="h-6 w-32" />
        <div class="flex gap-2">
          <Skeleton v-for="i in 5" :key="i" class="h-10 w-16 rounded-full" />
        </div>
        <div class="grid gap-3 sm:grid-cols-2">
          <Skeleton v-for="i in 4" :key="i" class="h-24 rounded-xl" />
        </div>
      </div>
    </div>
  </div>
</template>
