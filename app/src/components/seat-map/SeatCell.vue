<script setup lang="ts">
import { Accessibility, Crown } from 'lucide-vue-next'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { cn } from '@/lib/cn'
import { rowLabel } from '@/lib/seatLayoutEditor'
import type { Seat, SeatStatus, SeatType } from '@/types/seats'

const props = defineProps<{
  seat: Seat
  selfHeld?: boolean
  pending?: boolean
  interactive?: boolean
}>()

const emit = defineEmits<{
  select: [seatId: string]
}>()

const { t } = useI18n()

const displayStatus = computed<SeatStatus | 'SELECTED'>(() => {
  if (props.selfHeld || props.pending) {
    return 'SELECTED'
  }
  return props.seat.status
})

const typeLabel = computed(() => {
  if (props.seat.type === 'blocked') {
    return ''
  }
  return t(`seatMap.types.${props.seat.type}`)
})

const statusLabel = computed(() => {
  const keyMap: Record<SeatStatus | 'SELECTED', string> = {
    AVAILABLE: 'seatMap.status.available',
    HELD: 'seatMap.status.held',
    SOLD: 'seatMap.status.sold',
    BLOCKED: 'seatMap.status.blocked',
    SELECTED: 'seatMap.status.selected',
  }
  return t(keyMap[displayStatus.value])
})

const ariaLabel = computed(() =>
  typeLabel.value
    ? t('seatMap.seatAriaWithType', {
        seatId: props.seat.seatId,
        type: typeLabel.value,
        status: statusLabel.value,
      })
    : t('seatMap.seatAria', { seatId: props.seat.seatId, status: statusLabel.value }),
)

const availableTypeClasses: Record<Exclude<SeatType, 'blocked'>, string> = {
  standard:
    'bg-subtle border-surface-border text-copy-secondary hover:border-brand/40 cursor-pointer',
  vip: 'bg-accent-dim border-brand/60 text-copy-primary hover:border-brand cursor-pointer',
  wheelchair:
    'bg-state-success-dim border-state-success/50 text-copy-primary hover:border-state-success cursor-pointer',
}

const statusClasses: Record<Exclude<SeatStatus, 'AVAILABLE'> | 'SELECTED', string> = {
  HELD: 'bg-state-warning-dim border-state-warning text-copy-muted cursor-not-allowed',
  SOLD: 'bg-elevated border-surface-border-subtle text-copy-faint cursor-not-allowed',
  BLOCKED:
    'bg-state-error-dim border-state-error/50 text-copy-faint cursor-not-allowed seat-blocked',
  SELECTED: 'bg-accent-dim border-transparent text-copy-primary seat-self-held cursor-pointer',
}

const classes = computed(() => {
  const status = displayStatus.value
  const statusClass =
    status === 'AVAILABLE'
      ? availableTypeClasses[props.seat.type === 'blocked' ? 'standard' : props.seat.type]
      : statusClasses[status]

  return cn(
    'relative flex h-10 w-10 items-center justify-center rounded-md border text-[9px] font-medium leading-none transition-colors',
    'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent-glow',
    statusClass,
    props.pending && 'opacity-70',
  )
})

const showVipBadge = computed(() => props.seat.type === 'vip')
const showWheelchairBadge = computed(() => props.seat.type === 'wheelchair')

const seatLabel = computed(() => `${rowLabel(props.seat.row)}${props.seat.col}`)

const isClickable = computed(() => {
  if (!props.interactive) {
    return false
  }
  if (props.selfHeld || props.pending) {
    return true
  }
  return props.seat.status === 'AVAILABLE'
})

function onClick() {
  if (!isClickable.value) {
    return
  }
  emit('select', props.seat.seatId)
}

function onKeydown(event: KeyboardEvent) {
  if (!isClickable.value) {
    return
  }
  if (event.key === 'Enter' || event.key === ' ') {
    event.preventDefault()
    emit('select', props.seat.seatId)
  }
}
</script>

<template>
  <button
    type="button"
    :class="classes"
    :disabled="!isClickable"
    :aria-label="ariaLabel"
    :aria-pressed="selfHeld"
    @click="onClick"
    @keydown="onKeydown"
  >
    <Crown
      v-if="showVipBadge"
      class="pointer-events-none absolute left-0.5 top-0.5 h-2.5 w-2.5 text-brand"
      aria-hidden="true"
    />
    <Accessibility
      v-if="showWheelchairBadge"
      class="pointer-events-none absolute left-0.5 top-0.5 h-2.5 w-2.5 text-state-success"
      aria-hidden="true"
    />
    {{ seatLabel }}
  </button>
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

.seat-self-held {
  border: 2px solid transparent;
  background:
    linear-gradient(var(--color-accent-dim), var(--color-accent-dim)) padding-box,
    var(--gradient-brand) border-box;
}
</style>
