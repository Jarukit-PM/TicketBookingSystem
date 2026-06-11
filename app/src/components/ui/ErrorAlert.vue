<script setup lang="ts">
import { AlertCircle } from 'lucide-vue-next'
import { cn } from '@/lib/cn'

defineProps<{
  message: string
  class?: string
}>()
</script>

<template>
  <p
    :class="
      cn(
        'flex items-start gap-2 rounded-lg border border-state-error/30 bg-state-error-dim px-4 py-3 text-sm text-state-error',
        $props.class,
      )
    "
    role="alert"
  >
    <AlertCircle class="mt-0.5 h-4 w-4 shrink-0" aria-hidden="true" />
    <span>{{ message }}</span>
  </p>
</template>
