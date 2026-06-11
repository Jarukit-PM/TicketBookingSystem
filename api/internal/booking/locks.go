package booking

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const confirmLockTTL = 10 * time.Second

// ConfirmLockKey returns the Redis key for a seat confirm lock.
func ConfirmLockKey(showtimeID, seatID string) string {
	return fmt.Sprintf("lock:confirm:%s:%s", showtimeID, seatID)
}

func acquireConfirmLocks(ctx context.Context, rdb *redis.Client, showtimeID string, seatIDs []string) (func(context.Context), error) {
	acquired := make([]string, 0, len(seatIDs))
	for _, seatID := range seatIDs {
		key := ConfirmLockKey(showtimeID, seatID)
		ok, err := rdb.SetNX(ctx, key, "1", confirmLockTTL).Result()
		if err != nil {
			releaseConfirmLocks(ctx, rdb, showtimeID, acquired)
			return nil, fmt.Errorf("acquire confirm lock: %w", err)
		}
		if !ok {
			releaseConfirmLocks(ctx, rdb, showtimeID, acquired)
			return nil, ErrSeatConflict
		}
		acquired = append(acquired, seatID)
	}

	return func(ctx context.Context) {
		releaseConfirmLocks(ctx, rdb, showtimeID, acquired)
	}, nil
}

func releaseConfirmLocks(ctx context.Context, rdb *redis.Client, showtimeID string, seatIDs []string) {
	if rdb == nil || len(seatIDs) == 0 {
		return
	}
	keys := make([]string, len(seatIDs))
	for i, seatID := range seatIDs {
		keys[i] = ConfirmLockKey(showtimeID, seatID)
	}
	_ = rdb.Del(ctx, keys...).Err()
}
