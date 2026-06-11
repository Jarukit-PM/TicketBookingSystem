<script setup lang="ts">
defineProps<{
  activeTab: 'now_showing' | 'coming_soon'
}>()

const emit = defineEmits<{
  change: [tab: 'now_showing' | 'coming_soon']
}>()

const tabs = [
  { id: 'now_showing' as const, label: 'Now Showing' },
  { id: 'coming_soon' as const, label: 'Coming Soon' },
]
</script>

<template>
  <div class="flex flex-wrap gap-2">
    <button
      v-for="tab in tabs"
      :key="tab.id"
      type="button"
      class="rounded-full border px-4 py-2 text-sm font-medium transition-colors"
      :class="
        activeTab === tab.id
          ? 'border-brand/30 bg-accent-dim text-brand'
          : 'border-surface-border bg-subtle text-copy-secondary hover:bg-elevated'
      "
      @click="emit('change', tab.id)"
    >
      {{ tab.label }}
    </button>
  </div>
</template>
