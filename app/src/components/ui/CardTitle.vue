<script setup lang="ts">
import { cn } from '@/lib/cn'

defineProps<{
  class?: string
}>()
</script>

<template>
  <h3 :class="cn('text-lg font-semibold text-copy-primary', $props.class)">
    <slot />
  </h3>
</template>
