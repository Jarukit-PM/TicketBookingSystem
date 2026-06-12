<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { cn } from '@/lib/cn'
import {
  fillGrid,
  gridToLayout,
  layoutToGrid,
  resizeGrid,
  rowLabel,
  type LayoutCell,
  type LayoutGrid,
} from '@/lib/seatLayoutEditor'
import type { ScreenLayout } from '@/types/catalog'
import type { SeatType } from '@/types/seats'

const props = defineProps<{
  modelValue: ScreenLayout
}>()

const emit = defineEmits<{
  'update:modelValue': [layout: ScreenLayout]
}>()

const { t } = useI18n()

const grid = ref<LayoutGrid>(layoutToGrid(props.modelValue))
const brush = ref<SeatType | 'erase'>('standard')
const isPainting = ref(false)

const seatCount = computed(() => gridToLayout(grid.value).seats.length)

const brushOptions = computed(() => [
  { id: 'standard' as const, label: t('admin.screens.editor.types.standard'), class: 'bg-subtle border-surface-border' },
  { id: 'vip' as const, label: t('admin.screens.editor.types.vip'), class: 'bg-accent-dim border-brand/60' },
  {
    id: 'wheelchair' as const,
    label: t('admin.screens.editor.types.wheelchair'),
    class: 'bg-sky-500/20 border-sky-400/60',
  },
  {
    id: 'blocked' as const,
    label: t('admin.screens.editor.types.blocked'),
    class: 'bg-state-error-dim border-state-error/50 seat-blocked',
  },
  {
    id: 'erase' as const,
    label: t('admin.screens.editor.types.aisle'),
    class: 'border-dashed border-surface-border bg-transparent',
  },
])

const typeClasses: Record<SeatType, string> = {
  standard: 'bg-subtle border-surface-border text-copy-secondary',
  vip: 'bg-accent-dim border-brand/60 text-copy-primary',
  wheelchair: 'bg-sky-500/20 border-sky-400/60 text-copy-primary',
  blocked: 'bg-state-error-dim border-state-error/50 text-copy-faint seat-blocked',
}

function emitLayout() {
  emit('update:modelValue', gridToLayout(grid.value))
}

function setGrid(next: LayoutGrid) {
  grid.value = next
  emitLayout()
}

function onDimensionChange() {
  const rows = Math.max(1, Math.min(26, grid.value.rows))
  const cols = Math.max(1, Math.min(30, grid.value.cols))
  setGrid(resizeGrid(grid.value, rows, cols))
}

function onFillAll(type: SeatType) {
  setGrid(fillGrid(grid.value, type))
}

function paintCell(row: number, col: number) {
  const next = {
    ...grid.value,
    cells: grid.value.cells.map((line) => [...line]),
  }
  next.cells[row]![col] = brush.value === 'erase' ? null : brush.value
  setGrid(next)
}

function onCellPointerDown(row: number, col: number) {
  isPainting.value = true
  paintCell(row, col)
}

function onCellPointerEnter(row: number, col: number) {
  if (!isPainting.value) return
  paintCell(row, col)
}

function stopPainting() {
  isPainting.value = false
}

function cellClass(cell: LayoutCell) {
  if (!cell) {
    return 'border border-dashed border-surface-border/60 bg-transparent text-copy-faint hover:border-brand/40'
  }
  return typeClasses[cell]
}

function cellLabel(cell: LayoutCell, _row: number, col: number) {
  if (!cell) return ''
  return `${col + 1}`
}

watch(
  () => props.modelValue,
  (layout) => {
    const current = gridToLayout(grid.value)
    if (JSON.stringify(current) === JSON.stringify(layout)) return
    grid.value = layoutToGrid(layout)
  },
  { deep: true },
)
</script>

<template>
  <div class="space-y-4" @pointerup="stopPainting" @pointerleave="stopPainting">
    <div class="flex flex-wrap items-end gap-4">
      <div class="space-y-1">
        <label class="text-sm text-copy-secondary" for="layout-rows">
          {{ t('admin.screens.editor.rows') }}
        </label>
        <input
          id="layout-rows"
          v-model.number="grid.rows"
          type="number"
          min="1"
          max="26"
          class="w-20 rounded-lg border border-surface-border bg-surface px-3 py-2 text-sm text-copy-primary"
          @change="onDimensionChange"
        />
      </div>
      <div class="space-y-1">
        <label class="text-sm text-copy-secondary" for="layout-cols">
          {{ t('admin.screens.editor.cols') }}
        </label>
        <input
          id="layout-cols"
          v-model.number="grid.cols"
          type="number"
          min="1"
          max="30"
          class="w-20 rounded-lg border border-surface-border bg-surface px-3 py-2 text-sm text-copy-primary"
          @change="onDimensionChange"
        />
      </div>
      <button
        type="button"
        class="rounded-lg border border-surface-border px-3 py-2 text-sm text-copy-secondary hover:bg-subtle"
        @click="onFillAll('standard')"
      >
        {{ t('admin.screens.editor.fillStandard') }}
      </button>
      <p class="text-sm text-copy-muted">
        {{ t('admin.screens.editor.seatCount', { count: seatCount }) }}
      </p>
    </div>

    <div>
      <p class="mb-2 text-sm text-copy-secondary">{{ t('admin.screens.editor.brush') }}</p>
      <div class="flex flex-wrap gap-2">
        <button
          v-for="option in brushOptions"
          :key="option.id"
          type="button"
          :class="
            cn(
              'flex items-center gap-2 rounded-lg border px-3 py-2 text-sm transition-colors',
              brush === option.id
                ? 'border-brand bg-accent-dim text-brand'
                : 'border-surface-border text-copy-secondary hover:bg-subtle',
            )
          "
          :aria-pressed="brush === option.id"
          @click="brush = option.id"
        >
          <span class="inline-block h-4 w-4 rounded border" :class="option.class" aria-hidden="true" />
          {{ option.label }}
        </button>
      </div>
    </div>

    <div class="rounded-xl border border-surface-border bg-surface p-4">
      <p class="mb-4 text-center text-xs font-medium uppercase tracking-wider text-copy-muted">
        {{ t('admin.screens.editor.screenLabel') }}
      </p>

      <div class="overflow-x-auto">
        <div class="inline-block min-w-full">
          <div
            class="mb-2 grid gap-2"
            :style="{ gridTemplateColumns: `2.5rem repeat(${grid.cols}, minmax(2.5rem, 1fr))` }"
          >
            <div />
            <div
              v-for="col in grid.cols"
              :key="`col-${col}`"
              class="text-center text-[10px] font-medium text-copy-muted"
            >
              {{ col }}
            </div>
          </div>

          <div
            v-for="(rowCells, rowIdx) in grid.cells"
            :key="`row-${rowIdx}`"
            class="mb-2 grid gap-2"
            :style="{ gridTemplateColumns: `2.5rem repeat(${grid.cols}, minmax(2.5rem, 1fr))` }"
          >
            <div class="flex items-center justify-center text-xs font-medium text-copy-muted">
              {{ rowLabel(rowIdx + 1) }}
            </div>
            <button
              v-for="(cell, colIdx) in rowCells"
              :key="`${rowIdx}-${colIdx}`"
              type="button"
              :class="
                cn(
                  'flex h-10 w-10 items-center justify-center rounded-md border text-[10px] font-medium transition-colors select-none',
                  'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent-glow',
                  cellClass(cell),
                )
              "
              :aria-label="
                cell
                  ? t('admin.screens.editor.seatAria', {
                      seatId: `${rowLabel(rowIdx + 1)}-${colIdx + 1}`,
                      type: t(`admin.screens.editor.types.${cell}`),
                    })
                  : t('admin.screens.editor.aisleAria', { row: rowLabel(rowIdx + 1), col: colIdx + 1 })
              "
              @pointerdown.prevent="onCellPointerDown(rowIdx, colIdx)"
              @pointerenter="onCellPointerEnter(rowIdx, colIdx)"
            >
              {{ cellLabel(cell, rowIdx, colIdx) }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <p class="text-xs text-copy-muted">{{ t('admin.screens.editor.hint') }}</p>
  </div>
</template>

<style scoped>
.seat-blocked {
  background-image: repeating-linear-gradient(
    135deg,
    transparent,
    transparent 4px,
    rgba(248, 113, 113, 0.15) 4px,
    rgba(248, 113, 113, 0.15) 8px
  );
}
</style>
