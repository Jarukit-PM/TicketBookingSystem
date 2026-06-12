<script setup lang="ts">
import { CheckCircle2 } from 'lucide-vue-next'

defineProps<{
  message: string
}>()
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="translate-y-2 opacity-0"
      enter-to-class="translate-y-0 opacity-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="translate-y-0 opacity-100"
      leave-to-class="translate-y-2 opacity-0"
    >
      <div
        v-if="message"
        class="pointer-events-none fixed right-4 top-4 z-50 flex max-w-sm items-center gap-2 rounded-lg border border-state-success/40 bg-state-success-dim px-4 py-3 text-sm text-state-success shadow-1"
        role="status"
        aria-live="polite"
      >
        <CheckCircle2 class="h-4 w-4 shrink-0" aria-hidden="true" />
        <span>{{ message }}</span>
      </div>
    </Transition>
  </Teleport>
</template>
