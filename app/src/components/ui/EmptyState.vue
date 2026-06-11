<script setup lang="ts">
import type { Component } from 'vue'
import { cn } from '@/lib/cn'

defineProps<{
  icon?: Component
  title: string
  description?: string
  class?: string
}>()
</script>

<template>
  <div
    :class="cn('rounded-xl border border-surface-border bg-surface px-6 py-16 text-center', $props.class)"
    role="status"
  >
    <div
      v-if="icon"
      class="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-full bg-subtle"
    >
      <component :is="icon" class="h-6 w-6 text-copy-muted" aria-hidden="true" />
    </div>
    <p class="text-copy-primary">{{ title }}</p>
    <p v-if="description" class="mt-2 text-sm text-copy-secondary">{{ description }}</p>
    <div v-if="$slots.action" class="mt-6">
      <slot name="action" />
    </div>
  </div>
</template>
