<script setup lang="ts">
defineProps<{
  label: string
  value: string
  hint?: string
}>()
</script>

<template>
  <div class="rounded-xl border border-surface-border bg-surface p-5 shadow-elevation-1">
    <p class="text-xs font-medium uppercase tracking-wide text-copy-muted">{{ label }}</p>
    <p class="mt-2 text-3xl font-semibold text-copy-primary">{{ value }}</p>
    <p v-if="hint" class="mt-1 text-sm text-copy-secondary">{{ hint }}</p>
  </div>
</template>
