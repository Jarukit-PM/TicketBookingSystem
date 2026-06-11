<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useLocaleFormat } from '@/composables/useLocaleFormat'
import type { ShowtimeDateOption } from '@/composables/useShowtimeDates'

const WINDOW_SIZE = 7

const props = defineProps<{
  dates: ShowtimeDateOption[]
  modelValue: string | null
  timeZone?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const { t } = useI18n()
const { formatDayShort, formatDayNumber, formatMonthYear } = useLocaleFormat()

const windowStart = ref(0)

const selectedOption = computed(() =>
  props.dates.find((option) => option.key === props.modelValue),
)

const monthLabel = computed(() =>
  selectedOption.value
    ? formatMonthYear(selectedOption.value.sampleIso, props.timeZone)
    : '',
)

const visibleDates = computed(() =>
  props.dates.slice(windowStart.value, windowStart.value + WINDOW_SIZE),
)

const canGoPrev = computed(() => windowStart.value > 0)

const canGoNext = computed(() => windowStart.value + WINDOW_SIZE < props.dates.length)

function maxWindowStart(): number {
  return Math.max(0, props.dates.length - WINDOW_SIZE)
}

function ensureSelectedInWindow(): void {
  if (!props.modelValue) {
    return
  }
  const index = props.dates.findIndex((option) => option.key === props.modelValue)
  if (index === -1) {
    return
  }
  if (index < windowStart.value || index >= windowStart.value + WINDOW_SIZE) {
    const centered = index - Math.floor(WINDOW_SIZE / 2)
    windowStart.value = Math.min(Math.max(0, centered), maxWindowStart())
  }
}

function selectDate(dateKey: string): void {
  emit('update:modelValue', dateKey)
}

function goPrev(): void {
  windowStart.value = Math.max(0, windowStart.value - WINDOW_SIZE)
}

function goNext(): void {
  windowStart.value = Math.min(maxWindowStart(), windowStart.value + WINDOW_SIZE)
}

watch(
  () => props.modelValue,
  () => {
    ensureSelectedInWindow()
  },
  { immediate: true },
)

watch(
  () => props.dates,
  () => {
    windowStart.value = 0
    ensureSelectedInWindow()
  },
)
</script>

<template>
  <div class="mb-5">
    <p
      v-if="monthLabel"
      class="mb-3 text-sm font-medium text-copy-secondary"
      aria-live="polite"
    >
      {{ monthLabel }}
    </p>

    <div class="flex items-center gap-2">
      <button
        type="button"
        class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full border border-surface-border bg-surface text-copy-secondary transition-colors hover:border-brand/40 hover:bg-accent-dim hover:text-copy-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent-glow disabled:pointer-events-none disabled:opacity-40"
        :aria-label="t('movie.dateFilterPrevious')"
        :disabled="!canGoPrev"
        @click="goPrev"
      >
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="h-5 w-5" aria-hidden="true">
          <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 0 1 0 1.414L9.414 10l3.293 3.293a1 1 0 0 1-1.414 1.414l-4-4a1 1 0 0 1 0-1.414l4-4a1 1 0 0 1 1.414 0Z" clip-rule="evenodd" />
        </svg>
      </button>

      <div
        role="tablist"
        :aria-label="t('movie.dateFilterLabel')"
        class="grid min-w-0 flex-1 grid-cols-7 gap-1.5 sm:gap-2"
      >
        <button
          v-for="option in visibleDates"
          :key="option.key"
          type="button"
          role="tab"
          :aria-selected="option.key === modelValue"
          class="flex min-w-0 flex-col items-center rounded-xl border px-1 py-2 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent-glow sm:px-2 sm:py-2.5"
          :class="
            option.key === modelValue
              ? 'border-transparent bg-gradient-brand text-white shadow-glow-brand'
              : 'border-brand/30 bg-surface text-copy-primary hover:bg-accent-dim'
          "
          @click="selectDate(option.key)"
        >
          <span class="truncate text-[10px] font-medium uppercase tracking-wide sm:text-xs">
            {{ formatDayShort(option.sampleIso, timeZone) }}
          </span>
          <span class="text-lg font-semibold leading-tight sm:text-xl">
            {{ formatDayNumber(option.sampleIso, timeZone) }}
          </span>
        </button>
      </div>

      <button
        type="button"
        class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full border border-surface-border bg-surface text-copy-secondary transition-colors hover:border-brand/40 hover:bg-accent-dim hover:text-copy-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent-glow disabled:pointer-events-none disabled:opacity-40"
        :aria-label="t('movie.dateFilterNext')"
        :disabled="!canGoNext"
        @click="goNext"
      >
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="h-5 w-5" aria-hidden="true">
          <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 0 1 0-1.414L10.586 10 7.293 6.707a1 1 0 0 1 1.414-1.414l4 4a1 1 0 0 1 0 1.414l-4 4a1 1 0 0 1-1.414 0Z" clip-rule="evenodd" />
        </svg>
      </button>
    </div>
  </div>
</template>
