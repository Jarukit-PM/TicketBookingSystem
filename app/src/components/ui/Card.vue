<script setup lang="ts">
import { cn } from '@/lib/cn'

defineProps<{
  class?: string
}>()
</script>

<template>
  <div
    :class="
      cn(
        'rounded-xl border border-surface-border bg-surface shadow-elevation-1',
        $props.class,
      )
    "
  >
    <slot />
  </div>
</template>
