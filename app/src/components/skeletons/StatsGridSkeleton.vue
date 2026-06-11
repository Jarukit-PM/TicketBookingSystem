<script setup lang="ts">
import { Skeleton } from '@/components/ui'

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
  <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3" aria-busy="true">
    <div
      v-for="i in count"
      :key="i"
      class="rounded-xl border border-surface-border bg-surface p-5 shadow-elevation-1"
    >
      <Skeleton class="h-3 w-24" />
      <Skeleton class="mt-3 h-9 w-16" />
      <Skeleton v-if="i === count" class="mt-2 h-4 w-32" />
    </div>
  </div>
</template>
