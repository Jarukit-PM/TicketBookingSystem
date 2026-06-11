<script setup lang="ts">
import type { Component } from 'vue'

defineProps<{
  label: string
  value: string
  hint?: string
  icon?: Component
}>()
</script>

<template>
  <div class="rounded-xl border border-surface-border bg-surface p-5 shadow-elevation-1">
    <div class="flex items-start justify-between gap-3">
      <p class="text-xs font-medium uppercase tracking-wide text-copy-muted">{{ label }}</p>
      <div
        v-if="icon"
        class="flex h-8 w-8 items-center justify-center rounded-full bg-accent-dim"
      >
        <component :is="icon" class="h-4 w-4 text-brand" aria-hidden="true" />
      </div>
    </div>
    <p class="mt-2 text-3xl font-semibold text-copy-primary">{{ value }}</p>
    <p v-if="hint" class="mt-1 text-sm text-copy-secondary">{{ hint }}</p>
  </div>
</template>
