import { apiRequest } from '@/api/client'
import type { ConfirmedBooking } from '@/types/bookings'

export function confirmBooking(
  showtimeId: string,
  idempotencyKey: string,
): Promise<ConfirmedBooking> {
  return apiRequest<ConfirmedBooking>('/bookings/confirm', {
    method: 'POST',
    headers: {
      'Idempotency-Key': idempotencyKey,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ showtimeId }),
  })
}
