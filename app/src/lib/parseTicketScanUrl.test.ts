import { describe, expect, it } from 'vitest'

import { parseTicketScanUrl } from './parseTicketScanUrl'

describe('parseTicketScanUrl', () => {
  it('parses absolute ticket URLs', () => {
    const result = parseTicketScanUrl(
      'https://cinema.example.com/ticket/TBS-ABC123?t=deadbeefcafebabe',
    )
    expect(result).toEqual({
      bookingRef: 'TBS-ABC123',
      token: 'deadbeefcafebabe',
    })
  })

  it('parses site-relative ticket paths', () => {
    const result = parseTicketScanUrl('/ticket/TBS-XYZ789?t=abc123')
    expect(result).toEqual({
      bookingRef: 'TBS-XYZ789',
      token: 'abc123',
    })
  })

  it('returns null for missing token', () => {
    expect(parseTicketScanUrl('https://cinema.example.com/ticket/TBS-ABC123')).toBeNull()
  })

  it('returns null for non-ticket URLs', () => {
    expect(parseTicketScanUrl('https://cinema.example.com/bookings/abc')).toBeNull()
  })

  it('returns null for empty input', () => {
    expect(parseTicketScanUrl('   ')).toBeNull()
  })
})
