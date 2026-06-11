export type ConfirmedBooking = {
  id: string
  bookingRef: string
  showtimeId: string
  seats: string[]
  total: number
  status: 'CONFIRMED'
  confirmedAt: string
}
