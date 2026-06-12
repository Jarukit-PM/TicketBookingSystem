export interface AuditLogEntry {
  id: string
  actorId?: string
  action: string
  entity: string
  entityId: string
  meta?: Record<string, unknown>
  createdAt: string
}

export interface EmailLogEntry {
  id: string
  bookingId: string
  type: string
  to: string
  providerId?: string
  status: string
  createdAt: string
}

export interface BookingSummary {
  id: string
  bookingRef: string
  userId?: string
  userEmail?: string
  showtimeId: string
  movieTitle: string
  seats: string[]
  total: number
  locale?: string
  confirmedAt: string
}

export interface AdminBookingDetail extends BookingSummary {
  startsAt: string
  cinemaId: string
  cinemaName: string
  screenId: string
  screenName: string
  posterUrl?: string
  status: string
}

export interface AdminDashboard {
  bookingsToday: number
  showtimesToday: number
  avgOccupancyPct: number
  recentBookings: BookingSummary[]
}
