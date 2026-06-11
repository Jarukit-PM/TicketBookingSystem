<script setup lang="ts">
import { cn } from '@/lib/cn'

const props = defineProps<{
  class?: string
}>()
</script>

<template>
  <div :class="cn('skeleton-shimmer rounded-md', props.class)" aria-hidden="true" />
</template>
