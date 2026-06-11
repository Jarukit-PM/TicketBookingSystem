<script setup lang="ts">
import { Skeleton } from '@/components/ui'

withDefaults(
  defineProps<{
    rows?: number
    columns?: number
  }>(),
  {
    rows: 5,
    columns: 6,
  },
)
</script>

<template>
  <div class="space-y-3" aria-busy="true">
    <div class="flex gap-4 border-b border-surface-border pb-3">
      <Skeleton v-for="col in columns" :key="`head-${col}`" class="h-4 flex-1" />
    </div>
    <div v-for="row in rows" :key="row" class="flex gap-4 py-2">
      <Skeleton v-for="col in columns" :key="`${row}-${col}`" class="h-4 flex-1" />
    </div>
  </div>
</template>
