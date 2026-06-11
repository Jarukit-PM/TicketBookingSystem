<script setup lang="ts">
import { Timer } from 'lucide-vue-next'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { cn } from '@/lib/cn'
import { useHoldCountdown } from '@/composables/useHoldCountdown'

const props = defineProps<{
  expiresAt: string | null
}>()

const { t } = useI18n()
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
    <p class="inline-flex items-center justify-center gap-1.5 text-xs font-medium uppercase tracking-wide">
      <Timer class="h-3.5 w-3.5" aria-hidden="true" />
      {{ t('seatMap.holdExpiresIn') }}
    </p>
    <p class="text-2xl font-semibold tabular-nums">{{ formatted }}</p>
  </div>
</template>
