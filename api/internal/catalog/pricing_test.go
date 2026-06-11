package catalog_test

import (
	"testing"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/catalog"
)

func TestTotalForSeats(t *testing.T) {
	tiers := catalog.PriceTiers{
		Standard:   1200,
		VIP:        1800,
		Wheelchair: 1200,
	}
	layout := []catalog.LayoutSeat{
		{SeatID: "A-1", Type: catalog.SeatTypeStandard},
		{SeatID: "A-2", Type: catalog.SeatTypeVIP},
		{SeatID: "A-3", Type: catalog.SeatTypeWheelchair},
		{SeatID: "X-1", Type: catalog.SeatTypeBlocked},
	}

	tests := []struct {
		name    string
		seatIDs []string
		want    int64
		wantErr bool
	}{
		{
			name:    "single standard",
			seatIDs: []string{"A-1"},
			want:    1200,
		},
		{
			name:    "mixed tiers",
			seatIDs: []string{"A-1", "A-2", "A-3"},
			want:    4200,
		},
		{
			name:    "unknown seat",
			seatIDs: []string{"Z-9"},
			wantErr: true,
		},
		{
			name:    "blocked seat",
			seatIDs: []string{"X-1"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := catalog.TotalForSeats(tiers, layout, tt.seatIDs)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %d, want %d", got, tt.want)
			}
		})
	}
}
