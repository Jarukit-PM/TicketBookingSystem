import { describe, expect, it } from 'vitest'

import {
  createEmptyGrid,
  gridToLayout,
  layoutToGrid,
  resizeGrid,
  rowLabel,
} from '@/lib/seatLayoutEditor'
import type { ScreenLayout } from '@/types/catalog'

describe('seatLayoutEditor', () => {
  it('maps row numbers to letters', () => {
    expect(rowLabel(1)).toBe('A')
    expect(rowLabel(26)).toBe('Z')
    expect(rowLabel(27)).toBe('R27')
  })

  it('round-trips layout through grid', () => {
    const layout: ScreenLayout = {
      seats: [
        { seatId: 'A-1', row: 1, col: 1, type: 'standard' },
        { seatId: 'A-3', row: 1, col: 3, type: 'vip' },
        { seatId: 'B-2', row: 2, col: 2, type: 'wheelchair' },
        { seatId: 'B-4', row: 2, col: 4, type: 'blocked' },
      ],
    }

    const grid = layoutToGrid(layout)
    expect(grid.rows).toBe(2)
    expect(grid.cols).toBe(4)
    expect(grid.cells[0]?.[0]).toBe('standard')
    expect(grid.cells[0]?.[1]).toBeNull()
    expect(grid.cells[0]?.[2]).toBe('vip')
    expect(gridToLayout(grid)).toEqual(layout)
  })

  it('preserves seats when resizing larger', () => {
    const grid = layoutToGrid({
      seats: [{ seatId: 'A-1', row: 1, col: 1, type: 'standard' }],
    })
    const resized = resizeGrid(grid, 3, 3)
    expect(resized.rows).toBe(3)
    expect(resized.cols).toBe(3)
    expect(resized.cells[0]?.[0]).toBe('standard')
    expect(resized.cells[2]?.[2]).toBeNull()
  })

  it('creates a filled grid', () => {
    const layout = gridToLayout(createEmptyGrid(2, 2, 'standard'))
    expect(layout.seats).toHaveLength(4)
    expect(layout.seats[0]?.seatId).toBe('A-1')
    expect(layout.seats[3]?.seatId).toBe('B-2')
  })
})
