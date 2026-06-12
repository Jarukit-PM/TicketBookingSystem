<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/lib/cn'

type BadgeVariant = 'confirmed' | 'hold-active' | 'hold-expired'

const props = withDefaults(
  defineProps<{
    variant?: BadgeVariant
    class?: string
  }>(),
  {
    variant: 'confirmed',
  },
)

const variantClasses: Record<BadgeVariant, string> = {
  confirmed: 'bg-state-success-dim text-state-success',
  'hold-active': 'bg-accent-dim text-brand',
  'hold-expired': 'bg-subtle text-copy-muted',
}

const classes = computed(() =>
  cn(
    'inline-flex shrink-0 items-center gap-1.5 whitespace-nowrap rounded-full px-3 py-1 text-xs font-medium leading-none',
    variantClasses[props.variant],
    props.class,
  ),
)
</script>

<template>
  <span :class="classes">
    <span
      v-if="variant === 'hold-active'"
      class="h-1.5 w-1.5 rounded-full bg-gradient-brand"
      aria-hidden="true"
    />
    <slot />
  </span>
</template>
