export type ConfirmedBooking = {
  id: string
  bookingRef: string
  showtimeId: string
  seats: string[]
  total: number
  status: 'CONFIRMED'
  confirmedAt: string
}

export type BookingListItem = {
  id: string
  bookingRef: string
  showtimeId: string
  seats: string[]
  total: number
  status: 'CONFIRMED'
  confirmedAt: string
  startsAt: string
  movie: {
    id: string
    title: string
    posterUrl: string
  }
  cinema: {
    id: string
    name: string
  }
  screen: {
    id: string
    name: string
  }
}
