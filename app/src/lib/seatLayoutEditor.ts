import type { LayoutSeat, ScreenLayout } from '@/types/catalog'
import type { SeatType } from '@/types/seats'

export type LayoutCell = SeatType | null

export type LayoutGrid = {
  rows: number
  cols: number
  cells: LayoutCell[][]
}

export const SEAT_TYPES: SeatType[] = ['standard', 'vip', 'wheelchair', 'blocked']

export function rowLabel(row: number): string {
  if (row >= 1 && row <= 26) {
    return String.fromCharCode(64 + row)
  }
  return `R${row}`
}

export function createEmptyGrid(rows: number, cols: number, fill: LayoutCell = null): LayoutGrid {
  const safeRows = Math.max(1, rows)
  const safeCols = Math.max(1, cols)
  return {
    rows: safeRows,
    cols: safeCols,
    cells: Array.from({ length: safeRows }, () => Array.from({ length: safeCols }, () => fill)),
  }
}

export function layoutToGrid(layout: ScreenLayout): LayoutGrid {
  if (!layout.seats.length) {
    return createEmptyGrid(4, 8, 'standard')
  }

  let maxRow = 0
  let maxCol = 0
  for (const seat of layout.seats) {
    if (seat.row > maxRow) maxRow = seat.row
    if (seat.col > maxCol) maxCol = seat.col
  }

  const grid = createEmptyGrid(maxRow, maxCol)
  for (const seat of layout.seats) {
    grid.cells[seat.row - 1]![seat.col - 1] = seat.type as SeatType
  }
  return grid
}

export function gridToLayout(grid: LayoutGrid): ScreenLayout {
  const seats: LayoutSeat[] = []

  for (let row = 0; row < grid.rows; row++) {
    for (let col = 0; col < grid.cols; col++) {
      const type = grid.cells[row]?.[col]
      if (!type) continue
      seats.push({
        seatId: `${rowLabel(row + 1)}-${col + 1}`,
        row: row + 1,
        col: col + 1,
        type,
      })
    }
  }

  return { seats }
}

export function resizeGrid(grid: LayoutGrid, rows: number, cols: number): LayoutGrid {
  const next = createEmptyGrid(rows, cols)
  const copyRows = Math.min(grid.rows, rows)
  const copyCols = Math.min(grid.cols, cols)

  for (let row = 0; row < copyRows; row++) {
    for (let col = 0; col < copyCols; col++) {
      next.cells[row]![col] = grid.cells[row]?.[col] ?? null
    }
  }

  return next
}

export function fillGrid(grid: LayoutGrid, type: SeatType): LayoutGrid {
  return {
    ...grid,
    cells: grid.cells.map((row) => row.map(() => type)),
  }
}
