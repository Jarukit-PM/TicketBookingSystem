package catalog_test

import (
	"testing"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/catalog"
)

func TestSellableSeatCount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		layout catalog.ScreenLayout
		want   int
	}{
		{
			name: "excludes blocked seats",
			layout: catalog.ScreenLayout{
				Seats: []catalog.LayoutSeat{
					{SeatID: "A-1", Type: catalog.SeatTypeStandard},
					{SeatID: "A-2", Type: catalog.SeatTypeVIP},
					{SeatID: "A-X", Type: catalog.SeatTypeBlocked},
				},
			},
			want: 2,
		},
		{
			name:   "empty layout",
			layout: catalog.ScreenLayout{},
			want:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := catalog.SellableSeatCount(tt.layout); got != tt.want {
				t.Fatalf("SellableSeatCount() = %d, want %d", got, tt.want)
			}
		})
	}
}
