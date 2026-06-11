export type SeatStatus = 'AVAILABLE' | 'HELD' | 'SOLD' | 'BLOCKED'

export type SeatType = 'standard' | 'vip' | 'wheelchair' | 'blocked'

export interface Seat {
  seatId: string
  row: number
  col: number
  type: SeatType
  status: SeatStatus
}

export interface SeatMapSnapshot {
  showtimeId: string
  screenId: string
  screenName: string
  movieId: string
  startsAt: string
  priceTiers: Record<string, number>
  seats: Seat[]
}
