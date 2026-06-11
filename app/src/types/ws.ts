import type { SeatMapSnapshot } from '@/types/seats'

export type WsEventType = 'snapshot' | 'seat_held' | 'seat_released' | 'seat_sold'

export interface WsMessage {
  type: WsEventType
  payload: unknown
}

export interface WsSnapshotPayload {
  snapshot: SeatMapSnapshot
}

export interface WsSeatHeldPayload {
  seatId: string
  expiresAt: string
}

export interface WsSeatReleasedPayload {
  seatId: string
}

export interface WsSeatSoldPayload {
  seatId: string
}
