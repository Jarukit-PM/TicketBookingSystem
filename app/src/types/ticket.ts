export interface TicketDetail {
  bookingRef: string
  ticketUrl: string
  qrPngBase64: string
  seats: string[]
  total: number
  movieTitle: string
  cinemaName: string
  screenName: string
  startsAt: string
}
