<script setup lang="ts">
import { computed } from 'vue'
import SeatCell from '@/components/seat-map/SeatCell.vue'
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
  <div class="overflow-x-auto">
    <div
      class="inline-grid gap-2 p-4"
      :style="{ gridTemplateColumns: `repeat(${grid.cols}, minmax(2.5rem, 1fr))` }"
      role="grid"
      :aria-rowcount="grid.rows"
      :aria-colcount="grid.cols"
    >
      <template v-for="(row, rowIdx) in grid.cells" :key="rowIdx">
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
      </template>
    </div>
  </div>
</template>
