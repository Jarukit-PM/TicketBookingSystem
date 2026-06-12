const KNOWN_META_KEYS = new Set([
  'bookingRef',
  'seats',
  'seatIds',
  'seatId',
  'showtimeId',
  'total',
  'reason',
  'code',
  'message',
])

const META_FIELD_ORDER = [
  'bookingRef',
  'seats',
  'seatIds',
  'seatId',
  'showtimeId',
  'total',
  'reason',
  'code',
  'message',
] as const

export type AuditMetaFieldKey = (typeof META_FIELD_ORDER)[number]

export interface AuditMetaRow {
  key: AuditMetaFieldKey | string
  value: string
  mono?: boolean
}

function formatScalar(value: unknown): string {
  if (value === null || value === undefined) return ''
  if (typeof value === 'string') return value
  if (typeof value === 'number' || typeof value === 'boolean') return String(value)
  if (Array.isArray(value)) return value.map((item) => formatScalar(item)).filter(Boolean).join(', ')
  return ''
}

function seatsValue(meta: Record<string, unknown>): string | null {
  if (Array.isArray(meta.seats) && meta.seats.length) {
    return meta.seats.map((seat) => formatScalar(seat)).filter(Boolean).join(', ')
  }
  if (Array.isArray(meta.seatIds) && meta.seatIds.length) {
    return meta.seatIds.map((seat) => formatScalar(seat)).filter(Boolean).join(', ')
  }
  if (typeof meta.seatId === 'string' && meta.seatId) return meta.seatId
  return null
}

export function humanizeMetaKey(key: string): string {
  return key
    .replace(/([a-z0-9])([A-Z])/g, '$1 $2')
    .replace(/[_-]+/g, ' ')
    .replace(/\b\w/g, (char) => char.toUpperCase())
}

export function buildAuditMetaRows(
  meta: Record<string, unknown> | undefined,
  formatters: {
    formatTotal: (total: number) => string
    formatReason: (reason: string) => string
  },
): AuditMetaRow[] {
  if (!meta || !Object.keys(meta).length) return []

  const rows: AuditMetaRow[] = []

  for (const key of META_FIELD_ORDER) {
    if (key === 'seatIds' || key === 'seatId') continue

    if (key === 'seats') {
      const seats = seatsValue(meta)
      if (seats) rows.push({ key, value: seats })
      continue
    }

    const raw = meta[key]
    if (raw === null || raw === undefined || raw === '') continue

    if (key === 'total' && typeof raw === 'number') {
      rows.push({ key, value: formatters.formatTotal(raw) })
      continue
    }

    if (key === 'reason' && typeof raw === 'string') {
      rows.push({ key, value: formatters.formatReason(raw) })
      continue
    }

    const value = formatScalar(raw)
    if (!value) continue

    rows.push({
      key,
      value,
      mono: key === 'showtimeId' || key === 'bookingRef',
    })
  }

  for (const [key, raw] of Object.entries(meta)) {
    if (KNOWN_META_KEYS.has(key)) continue
    const value = formatScalar(raw)
    if (!value) continue
    rows.push({
      key,
      value,
      mono: typeof raw === 'string' && value.length > 24,
    })
  }

  return rows
}

export function auditMetaRowLabel(
  row: AuditMetaRow,
  labelFor: (key: AuditMetaFieldKey) => string,
): string {
  if (KNOWN_META_KEYS.has(row.key)) return labelFor(row.key as AuditMetaFieldKey)
  return humanizeMetaKey(row.key)
}
