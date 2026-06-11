import { apiGet } from '@/api/client'
import type { SeatMapSnapshot } from '@/types/seats'

export function fetchSeatMap(showtimeId: string): Promise<SeatMapSnapshot> {
  return apiGet<SeatMapSnapshot>(`/showtimes/${showtimeId}/seats`)
}
