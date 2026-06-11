package hold_test

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/hold"
)

func TestListHeldSeatIDs_readsRedisHolds(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis: %v", err)
	}
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	defer rdb.Close()

	showtimeID := "507f1f77bcf86cd799439011"
	ctx := context.Background()

	if err := rdb.Set(ctx, hold.SeatKey(showtimeID, "A-1"), `{"userId":"u1"}`, 0).Err(); err != nil {
		t.Fatalf("set hold: %v", err)
	}
	if err := rdb.Set(ctx, hold.SeatKey(showtimeID, "A-2"), `{"userId":"u2"}`, 0).Err(); err != nil {
		t.Fatalf("set hold: %v", err)
	}

	held, err := hold.ListHeldSeatIDs(ctx, rdb, showtimeID)
	if err != nil {
		t.Fatalf("ListHeldSeatIDs: %v", err)
	}

	for _, seatID := range []string{"A-1", "A-2"} {
		if _, ok := held[seatID]; !ok {
			t.Fatalf("expected held seat %s", seatID)
		}
	}
}
