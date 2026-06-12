package hold

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestFindUserIDForExpiredSeat(t *testing.T) {
	t.Parallel()

	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	ctx := context.Background()

	userID := primitive.NewObjectID()
	showtimeID := "6a2ae94e92c8e0854b5a8e68"
	seatID := "G-7"
	userKey := UserHoldsKey(userID.Hex(), showtimeID)

	if err := rdb.SAdd(ctx, userKey, seatID).Err(); err != nil {
		t.Fatal(err)
	}

	got := FindUserIDForExpiredSeat(ctx, rdb, showtimeID, seatID)
	if got != userID {
		t.Fatalf("FindUserIDForExpiredSeat() = %v, want %v", got, userID)
	}

	if got := FindUserIDForExpiredSeat(ctx, rdb, showtimeID, "Z-9"); got != primitive.NilObjectID {
		t.Fatalf("missing seat FindUserIDForExpiredSeat() = %v, want zero", got)
	}
}
