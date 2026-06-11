import { api } from '@/api/client'
import type { HoldResult } from '@/types/holds'

export function addHolds(showtimeId: string, seatIds: string[]): Promise<HoldResult> {
  return api.post<HoldResult>(`/showtimes/${showtimeId}/holds`, { seatIds })
}

export function removeHolds(showtimeId: string, seatIds: string[]): Promise<HoldResult> {
  return api.delete<HoldResult>(`/showtimes/${showtimeId}/holds`, { seatIds })
}
