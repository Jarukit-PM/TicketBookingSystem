<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/lib/cn'
import { useHoldCountdown } from '@/composables/useHoldCountdown'

const props = defineProps<{
  expiresAt: string | null
}>()

const { formatted, isActive, isUrgent } = useHoldCountdown(() => props.expiresAt)

const classes = computed(() =>
  cn(
    'rounded-xl border px-4 py-3 text-center transition-colors',
    isActive.value
      ? isUrgent.value
        ? 'border-state-error bg-state-error-dim text-state-error'
        : 'border-brand/40 bg-accent-dim text-brand'
      : 'border-surface-border bg-subtle text-copy-muted',
  ),
)
</script>

<template>
  <div v-if="isActive" :class="classes" role="timer" aria-live="polite">
    <p class="text-xs font-medium uppercase tracking-wide">Hold expires in</p>
    <p class="text-2xl font-semibold tabular-nums">{{ formatted }}</p>
  </div>
</template>
