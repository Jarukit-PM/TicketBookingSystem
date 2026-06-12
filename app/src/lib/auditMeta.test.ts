import { describe, expect, it } from 'vitest'

import { buildAuditMetaRows, humanizeMetaKey } from './auditMeta'

describe('auditMeta', () => {
  it('formats known booking metadata in display order', () => {
    const rows = buildAuditMetaRows(
      {
        showtimeId: 'abc123',
        bookingRef: 'TBS-PV6779',
        seats: ['B-2'],
        total: 22000,
      },
      {
        formatTotal: (total) => `฿${total / 100}`,
        formatReason: (reason) => reason,
      },
    )

    expect(rows.map((row) => row.key)).toEqual(['bookingRef', 'seats', 'showtimeId', 'total'])
    expect(rows[0]?.value).toBe('TBS-PV6779')
    expect(rows[1]?.value).toBe('B-2')
    expect(rows[3]?.value).toBe('฿220')
  })

  it('humanizes unknown metadata keys', () => {
    expect(humanizeMetaKey('customField')).toBe('Custom Field')
  })
})
