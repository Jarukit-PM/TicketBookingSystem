export interface BookingSummary {
  id: string
  bookingRef: string
  showtimeId: string
  movieTitle: string
  seats: string[]
  total: number
  confirmedAt: string
}

export interface AdminDashboard {
  bookingsToday: number
  showtimesToday: number
  avgOccupancyPct: number
  recentBookings: BookingSummary[]
}
