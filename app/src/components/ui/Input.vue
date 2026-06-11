<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/lib/cn'

const props = defineProps<{
  modelValue?: string
  type?: string
  placeholder?: string
  disabled?: boolean
  id?: string
  class?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const classes = computed(() =>
  cn(
    'w-full rounded-lg border border-surface-border bg-surface px-3 py-2 text-sm text-copy-primary',
    'placeholder:text-copy-muted',
    'focus:outline-none focus:ring-2 focus:ring-accent-glow focus:border-brand/50',
    'disabled:cursor-not-allowed disabled:opacity-50',
    props.class,
  ),
)

function onInput(event: Event) {
  const target = event.target as HTMLInputElement
  emit('update:modelValue', target.value)
}
</script>

<template>
  <input
    :id="id"
    :type="type ?? 'text'"
    :value="modelValue"
    :placeholder="placeholder"
    :disabled="disabled"
    :class="classes"
    @input="onInput"
  />
</template>
