package inventory

import (
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/catalog"
)

// SoldSeatSet unions seat IDs from confirmed bookings.
func SoldSeatSet(bookings []booking.Booking) map[string]struct{} {
	sold := make(map[string]struct{})
	for _, b := range bookings {
		for _, seatID := range b.Seats {
			sold[seatID] = struct{}{}
		}
	}
	return sold
}

// ComputeSeats derives per-seat status from layout, sold seats, and active holds.
func ComputeSeats(layout catalog.ScreenLayout, sold, held map[string]struct{}) []Seat {
	seats := make([]Seat, 0, len(layout.Seats))
	for _, ls := range layout.Seats {
		status := StatusAvailable
		switch {
		case ls.Type == catalog.SeatTypeBlocked:
			status = StatusBlocked
		case contains(sold, ls.SeatID):
			status = StatusSold
		case contains(held, ls.SeatID):
			status = StatusHeld
		}

		seats = append(seats, Seat{
			SeatID: ls.SeatID,
			Row:    ls.Row,
			Col:    ls.Col,
			Type:   ls.Type,
			Status: status,
		})
	}
	return seats
}

func contains(set map[string]struct{}, key string) bool {
	_, ok := set[key]
	return ok
}
