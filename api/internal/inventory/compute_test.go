package inventory_test

import (
	"testing"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/catalog"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/inventory"
)

func TestComputeSeats_statuses(t *testing.T) {
	layout := catalog.ScreenLayout{
		Seats: []catalog.LayoutSeat{
			{SeatID: "A-1", Type: catalog.SeatTypeStandard},
			{SeatID: "A-2", Type: catalog.SeatTypeStandard},
			{SeatID: "A-3", Type: catalog.SeatTypeStandard},
			{SeatID: "B-1", Type: catalog.SeatTypeBlocked},
		},
	}

	sold := inventory.SoldSeatSet([]booking.Booking{{Seats: []string{"A-2"}}})
	held := map[string]struct{}{"A-3": {}}

	seats := inventory.ComputeSeats(layout, sold, held)
	byID := map[string]string{}
	for _, s := range seats {
		byID[s.SeatID] = s.Status
	}

	want := map[string]string{
		"A-1": inventory.StatusAvailable,
		"A-2": inventory.StatusSold,
		"A-3": inventory.StatusHeld,
		"B-1": inventory.StatusBlocked,
	}
	for seatID, status := range want {
		if byID[seatID] != status {
			t.Fatalf("seat %s = %q, want %q (all: %+v)", seatID, byID[seatID], status, byID)
		}
	}
}
