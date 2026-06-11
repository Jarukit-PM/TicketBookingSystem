package hold

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// ListHeldSeatIDs returns seat IDs currently held in Redis for a showtime.
func ListHeldSeatIDs(ctx context.Context, rdb *redis.Client, showtimeID string) (map[string]struct{}, error) {
	if rdb == nil {
		return map[string]struct{}{}, nil
	}

	held := make(map[string]struct{})
	pattern := SeatKeyPattern(showtimeID)
	var cursor uint64

	for {
		keys, next, err := rdb.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return nil, fmt.Errorf("scan hold keys: %w", err)
		}

		for _, key := range keys {
			seatID, ok := ParseSeatIDFromKey(key, showtimeID)
			if ok {
				held[seatID] = struct{}{}
			}
		}

		cursor = next
		if cursor == 0 {
			break
		}
	}

	return held, nil
}
