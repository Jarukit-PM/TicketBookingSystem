export interface BookingSummary {
  id: string
  bookingRef: string
  userId?: string
  userEmail?: string
  showtimeId: string
  movieTitle: string
  seats: string[]
  total: number
  confirmedAt: string
}
