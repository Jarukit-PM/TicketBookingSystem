package booking

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	idempotencyKeyPrefix = "idempotency:confirm:"
	defaultIdempotencyTTL = 24 * time.Hour
)

// IdempotencyStore caches successful confirm results by client key.
type IdempotencyStore struct {
	redis *redis.Client
	ttl   time.Duration
}

// NewIdempotencyStore returns a Redis-backed idempotency cache.
func NewIdempotencyStore(rdb *redis.Client, ttl time.Duration) *IdempotencyStore {
	if ttl <= 0 {
		ttl = defaultIdempotencyTTL
	}
	return &IdempotencyStore{redis: rdb, ttl: ttl}
}

func idempotencyRedisKey(key string) string {
	return idempotencyKeyPrefix + key
}

// Get returns a cached booking for the idempotency key, or nil when absent.
func (s *IdempotencyStore) Get(ctx context.Context, key string) (*Booking, error) {
	if s.redis == nil || key == "" {
		return nil, nil
	}

	raw, err := s.redis.Get(ctx, idempotencyRedisKey(key)).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get idempotency: %w", err)
	}

	var b Booking
	if err := json.Unmarshal([]byte(raw), &b); err != nil {
		return nil, fmt.Errorf("decode idempotency: %w", err)
	}
	return &b, nil
}

// Set stores a successful booking for the idempotency key.
func (s *IdempotencyStore) Set(ctx context.Context, key string, b *Booking) error {
	if s.redis == nil || key == "" || b == nil {
		return nil
	}

	raw, err := json.Marshal(b)
	if err != nil {
		return fmt.Errorf("marshal idempotency: %w", err)
	}
	if err := s.redis.Set(ctx, idempotencyRedisKey(key), raw, s.ttl).Err(); err != nil {
		return fmt.Errorf("set idempotency: %w", err)
	}
	return nil
}
