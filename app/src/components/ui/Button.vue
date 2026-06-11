<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/lib/cn'

type ButtonVariant = 'primary' | 'secondary' | 'ghost' | 'destructive'

const props = withDefaults(
  defineProps<{
    variant?: ButtonVariant
    type?: 'button' | 'submit' | 'reset'
    disabled?: boolean
  }>(),
  {
    variant: 'primary',
    type: 'button',
    disabled: false,
  },
)

const variantClasses: Record<ButtonVariant, string> = {
  primary:
    'rounded-full bg-gradient-brand text-white font-medium hover:bg-gradient-brand-hover shadow-glow-brand',
  secondary:
    'rounded-full border border-white/20 bg-transparent text-copy-primary hover:border-brand/50 hover:bg-accent-dim',
  ghost:
    'rounded-full text-brand hover:bg-accent-dim px-4',
  destructive:
    'rounded-full text-state-error hover:bg-state-error-dim',
}

const classes = computed(() =>
  cn(
    'inline-flex items-center justify-center px-5 py-2.5 text-sm font-medium transition-colors',
    'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent-glow focus-visible:ring-offset-2 focus-visible:ring-offset-base',
    'disabled:pointer-events-none disabled:opacity-50',
    variantClasses[props.variant],
  ),
)
</script>

<template>
  <button :type="type" :class="classes" :disabled="disabled">
    <slot />
  </button>
</template>
