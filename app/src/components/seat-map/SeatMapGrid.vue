<script setup lang="ts">
import { computed } from 'vue'
import SeatCell from '@/components/seat-map/SeatCell.vue'
import { rowLabel } from '@/lib/seatLayoutEditor'
import type { Seat } from '@/types/seats'

const props = defineProps<{
  seats: Seat[]
  selfHeldIds?: Set<string>
  pendingIds?: Set<string>
  interactive?: boolean
}>()

const emit = defineEmits<{
  select: [seatId: string]
}>()

const grid = computed(() => {
  if (props.seats.length === 0) {
    return { rows: 0, cols: 0, cells: [] as Array<Array<Seat | null>> }
  }

  let maxRow = 0
  let maxCol = 0
  for (const seat of props.seats) {
    if (seat.row > maxRow) {
      maxRow = seat.row
    }
    if (seat.col > maxCol) {
      maxCol = seat.col
    }
  }

  const byPos = new Map<string, Seat>()
  for (const seat of props.seats) {
    byPos.set(`${seat.row}:${seat.col}`, seat)
  }

  const cells: Array<Array<Seat | null>> = []
  for (let row = 1; row <= maxRow; row++) {
    const line: Array<Seat | null> = []
    for (let col = 1; col <= maxCol; col++) {
      line.push(byPos.get(`${row}:${col}`) ?? null)
    }
    cells.push(line)
  }

  return { rows: maxRow, cols: maxCol, cells }
})
</script>

<template>
  <div class="-mx-2 overflow-x-auto px-2 sm:mx-0 sm:px-0">
    <div
      class="inline-block min-w-full p-2 sm:p-4"
      role="grid"
      :aria-rowcount="grid.rows + 1"
      :aria-colcount="grid.cols + 1"
    >
      <div
        class="mb-1.5 grid gap-1.5 sm:gap-2"
        :style="{ gridTemplateColumns: `2rem repeat(${grid.cols}, minmax(2.5rem, 1fr))` }"
      >
        <div aria-hidden="true" />
        <div
          v-for="col in grid.cols"
          :key="`col-${col}`"
          class="flex h-6 items-center justify-center text-[10px] font-medium text-copy-muted"
        >
          {{ col }}
        </div>
      </div>

      <div
        v-for="(row, rowIdx) in grid.cells"
        :key="rowIdx"
        class="mb-1.5 grid gap-1.5 sm:mb-2 sm:gap-2"
        :style="{ gridTemplateColumns: `2rem repeat(${grid.cols}, minmax(2.5rem, 1fr))` }"
        role="row"
      >
        <div
          class="flex h-10 items-center justify-center text-xs font-medium text-copy-muted"
          aria-hidden="true"
        >
          {{ rowLabel(rowIdx + 1) }}
        </div>
        <template v-for="(seat, colIdx) in row" :key="`${rowIdx}-${colIdx}`">
          <div v-if="seat" role="gridcell">
            <SeatCell
              :seat="seat"
              :self-held="selfHeldIds?.has(seat.seatId)"
              :pending="pendingIds?.has(seat.seatId)"
              :interactive="interactive"
              @select="emit('select', $event)"
            />
          </div>
          <div v-else class="h-10 w-10" aria-hidden="true" />
        </template>
      </div>
    </div>
  </div>
</template>
