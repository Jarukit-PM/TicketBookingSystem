import { apiRequest } from '@/api/client'
import type { TicketDetail } from '@/types/ticket'

export function fetchBookingTicket(bookingId: string): Promise<TicketDetail> {
  return apiRequest<TicketDetail>(`/bookings/${bookingId}/ticket`)
}
