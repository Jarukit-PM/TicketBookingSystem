export interface ParsedTicketScan {
  bookingRef: string
  token: string
}

/**
 * Extracts bookingRef and ticket token from a scanned QR payload.
 * Accepts absolute URLs, site-relative paths, or bare `/ticket/{ref}?t=...` strings.
 */
export function parseTicketScanUrl(input: string): ParsedTicketScan | null {
  const trimmed = input.trim()
  if (!trimmed) {
    return null
  }

  try {
    const url = new URL(trimmed, 'https://ticket.invalid')
    const match = url.pathname.match(/\/ticket\/([^/]+)\/?$/)
    if (!match?.[1]) {
      return null
    }

    const token = url.searchParams.get('t')
    if (!token) {
      return null
    }

    return {
      bookingRef: decodeURIComponent(match[1]),
      token,
    }
  } catch {
    return null
  }
}
