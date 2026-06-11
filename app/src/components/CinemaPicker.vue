<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { Cinema } from '@/types/catalog'

defineProps<{
  cinemas: Cinema[]
  modelValue: string | null
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const { t } = useI18n()
</script>

<template>
  <div class="flex flex-col gap-2">
    <label class="text-xs font-medium uppercase tracking-wide text-copy-muted" for="cinema-picker">
      {{ t('catalog.cinema') }}
    </label>
    <select
      id="cinema-picker"
      class="rounded-lg border border-surface-border bg-surface px-3 py-2.5 text-sm text-copy-primary focus:outline-none focus:ring-2 focus:ring-accent-glow"
      :value="modelValue ?? ''"
      @change="emit('update:modelValue', ($event.target as HTMLSelectElement).value)"
    >
      <option value="" disabled>{{ t('catalog.selectCinemaOption') }}</option>
      <option v-for="cinema in cinemas" :key="cinema.id" :value="cinema.id">
        {{ cinema.name }}
      </option>
    </select>
  </div>
</template>
