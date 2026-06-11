<script setup lang="ts">
const items = [
  { label: 'Available', class: 'bg-subtle border-surface-border' },
  { label: 'Held', class: 'bg-state-warning-dim border-state-warning' },
  { label: 'Sold', class: 'bg-elevated border-surface-border-subtle' },
  { label: 'Blocked', class: 'bg-state-error-dim border-state-error/50' },
] as const
</script>

<template>
  <ul class="flex flex-wrap gap-4 text-sm text-copy-secondary">
    <li v-for="item in items" :key="item.label" class="flex items-center gap-2">
      <span
        class="inline-block h-4 w-4 rounded border"
        :class="item.class"
        aria-hidden="true"
      />
      {{ item.label }}
    </li>
  </ul>
</template>
