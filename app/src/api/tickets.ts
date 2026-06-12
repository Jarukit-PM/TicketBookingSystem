import { apiGet, apiRequest } from '@/api/client'
import type { TicketDetail } from '@/types/ticket'

export function fetchBookingTicket(bookingId: string): Promise<TicketDetail> {
  return apiRequest<TicketDetail>(`/bookings/${bookingId}/ticket`)
}

export function fetchPublicTicket(bookingRef: string, token: string): Promise<TicketDetail> {
  return apiGet<TicketDetail>(`/tickets/${encodeURIComponent(bookingRef)}`, { t: token })
}
