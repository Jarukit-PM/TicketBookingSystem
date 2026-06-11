<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/lib/cn'
import type { Seat, SeatStatus } from '@/types/seats'

const props = defineProps<{
  seat: Seat
}>()

const statusClasses: Record<SeatStatus, string> = {
  AVAILABLE: 'bg-subtle border-surface-border text-copy-secondary hover:border-brand/40',
  HELD: 'bg-state-warning-dim border-state-warning text-copy-muted cursor-not-allowed',
  SOLD: 'bg-elevated border-surface-border-subtle text-copy-faint cursor-not-allowed',
  BLOCKED:
    'bg-state-error-dim border-state-error/50 text-copy-faint cursor-not-allowed seat-blocked',
}

const classes = computed(() =>
  cn(
    'flex h-10 w-10 items-center justify-center rounded-md border text-[10px] font-medium transition-colors',
    'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent-glow',
    statusClasses[props.seat.status],
  ),
)
</script>

<template>
  <div
    :class="classes"
    :aria-label="`${seat.seatId}, ${seat.status.toLowerCase()}`"
    role="img"
  >
    {{ seat.seatId.split('-')[1] }}
  </div>
</template>

<style scoped>
.seat-blocked {
  background-image: repeating-linear-gradient(
    135deg,
    transparent,
    transparent 4px,
    rgba(248, 113, 113, 0.15) 4px,
    rgba(248, 113, 113, 0.15) 8px
  );
}
</style>
