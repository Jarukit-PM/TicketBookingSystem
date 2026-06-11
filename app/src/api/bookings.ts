import { apiRequest } from '@/api/client'
import type { BookingListItem, ConfirmedBooking } from '@/types/bookings'

export function confirmBooking(
  showtimeId: string,
  idempotencyKey: string,
  locale: string,
): Promise<ConfirmedBooking> {
  return apiRequest<ConfirmedBooking>('/bookings/confirm', {
    method: 'POST',
    headers: {
      'Idempotency-Key': idempotencyKey,
      'Content-Type': 'application/json',
      'X-Locale': locale,
    },
    body: JSON.stringify({ showtimeId }),
  })
}

type MyBookingsResponse = {
  bookings: BookingListItem[]
}

export function fetchMyBookings(upcoming: boolean): Promise<BookingListItem[]> {
  return apiRequest<MyBookingsResponse>(`/bookings/mine?upcoming=${upcoming}`).then(
    (res) => res.bookings,
  )
}

export function fetchBookingDetail(id: string): Promise<BookingListItem> {
  return apiRequest<BookingListItem>(`/bookings/${id}`)
}
