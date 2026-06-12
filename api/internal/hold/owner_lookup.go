package hold

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FindUserIDForExpiredSeat looks up the user who held a seat from user_holds sets.
// The seat hold key may already be expired; the user_holds set can still list the seat.
func FindUserIDForExpiredSeat(
	ctx context.Context,
	rdb *redis.Client,
	showtimeID, seatID string,
) primitive.ObjectID {
	if rdb == nil || showtimeID == "" || seatID == "" {
		return primitive.NilObjectID
	}

	pattern := UserHoldsKeyPattern(showtimeID)
	var cursor uint64
	for {
		keys, next, err := rdb.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return primitive.NilObjectID
		}
		for _, key := range keys {
			member, err := rdb.SIsMember(ctx, key, seatID).Result()
			if err != nil || !member {
				continue
			}
			userID, parsedShowtimeID, ok := ParseUserHoldsKey(key)
			if !ok || parsedShowtimeID != showtimeID {
				continue
			}
			oid, err := primitive.ObjectIDFromHex(userID)
			if err != nil {
				continue
			}
			return oid
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return primitive.NilObjectID
}
