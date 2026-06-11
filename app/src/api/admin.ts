import { apiGet } from '@/api/client'

export interface TicketResolveResult {
  userId: string
  bookingId: string
}

export function resolveAdminTicket(ref: string, token: string): Promise<TicketResolveResult> {
  return apiGet<TicketResolveResult>('/admin/tickets/resolve', { ref, t: token })
}
